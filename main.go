// main package is the entry point of the program
package main

import (
	"benchmarktools/dataset"
	"benchmarktools/input"
	"benchmarktools/runner"
	"benchmarktools/utils"
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

const version = "v1.2.0"

func init() {
	var versionPrint bool
	const usage = `Usage: %s [options] <input file>

Run HADDOCK benchmarking

Options:
`
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "WARNING")
	_ = flag.Set("v", "2")

	flag.BoolVar(&versionPrint, "version", false, "Print version and exit")
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Fprintf(flag.CommandLine.Output(), usage, "benchmarktools")
		for _, f := range []string{"version"} {
			flag := flagSet.Lookup(f)
			fmt.Printf("  -%s: %s\n", f, flag.Usage)
		}
		fmt.Println("")
	}

}

func main() {

	flag.Parse()

	if utils.IsFlagPassed("version") {
		fmt.Printf("benchmarktools version %s\n", version)
		os.Exit(0)
	}

	args := os.Args[1:]
	inputF := args[0]

	glog.Info("###########################################")
	glog.Info(" Starting benchmarktools " + version)
	glog.Info("###########################################")
	glog.Info("Loading input file: " + inputF)

	inp, errInp := input.LoadInput(inputF)
	if errInp != nil {
		glog.Error("Failed to load input file: " + errInp.Error())
		return
	}

	errExec := inp.ValidateExecutable()
	if errExec != nil {
		glog.Error("Failed to validate executable: " + errExec.Error())
		return
	}

	errPatt := inp.ValidatePatterns()
	if errPatt != nil {
		glog.Error("ERROR: " + errPatt.Error())
		return
	}

	// haddockVersion := inp.General.HaddockVersion
	var haddockVersion int
	if utils.IsHaddock24(inp.General.HaddockDir) {
		haddockVersion = 2
	} else {
		haddockVersion = 3
	}

	if haddockVersion == 2 {
		runCns, errFind := input.FindHaddock24RunCns(inp.General.HaddockDir)
		if errFind != nil {
			glog.Error("Failed to find run.cns-conf: " + errFind.Error())
			return
		}

		runCnsParams, errLoad := input.LoadHaddock24Params(runCns)
		if errLoad != nil {
			glog.Error("Failed to load run.cns-conf parameters" + errLoad.Error())
			return
		}

		for _, scenario := range inp.Scenarios {
			scenarioParams := scenario.Parameters.CnsParams
			errValidate := input.ValidateRunCNSParams(runCnsParams, scenarioParams)
			if errValidate != nil {
				glog.Error("Failed to validate scenario parameters: " + errValidate.Error())
				return
			}
		}
	} else if haddockVersion == 3 {

		moduleArr, errParams := input.LoadHaddock3Params(inp.General.HaddockDir)
		if errParams != nil {
			glog.Error("Failed to load HADDOCK3 parameters: " + errParams.Error())
			return
		}

		for _, scenario := range inp.Scenarios {
			scenarioModules := scenario.Parameters.Modules
			errValidate := input.ValidateHaddock3Params(moduleArr, scenarioModules)
			if errValidate != nil {
				glog.Error("Failed to validate scenario parameters: " + errValidate.Error())
				return
			}
		}

	}

	// Load the dataset
	data, errDataset := dataset.LoadDataset(inp.General.WorkDir, inp.General.InputList, inp.General.ReceptorSuffix, inp.General.LigandSuffix)
	if errDataset != nil {
		glog.Error("Failed to load dataset: " + errDataset.Error())
		return
	}

	// Organize the dataset
	orgData, errOrganize := dataset.OrganizeDataset(inp.General.WorkDir, data)
	if errOrganize != nil {
		glog.Error("Failed to organize dataset: " + errOrganize.Error())
		return
	}

	// Setup the scenarios & create the jobs
	jobArr := []runner.Job{}

	for _, target := range orgData {
		glog.Info("Setting up target " + target.ID)

		for _, scenario := range inp.Scenarios {
			if haddockVersion == 2 {
				job, errSetup := target.SetupHaddock24Scenario(inp.General.WorkDir, inp.General.HaddockDir, scenario)
				if errSetup != nil {
					glog.Error("Failed to setup scenario: " + errSetup.Error())
					return
				}
				jobArr = append(jobArr, job)
			} else if haddockVersion == 3 {
				job, errSetup := target.SetupHaddock3Scenario(inp.General.WorkDir, scenario)
				if errSetup != nil {
					glog.Error("Failed to setup scenario: " + errSetup.Error())
					return
				}
				jobArr = append(jobArr, job)
			}
		}
	}

	// glog.Exit("Exiting for now...")

	// Taken form:
	// `https://gist.github.com/AntoineAugusti/80e99edfe205baf7a094`
	maxConcurrent := inp.General.MaxConcurrent
	glog.Info("Running " + fmt.Sprint(len(jobArr)) + " jobs with " + fmt.Sprint(maxConcurrent) + " concurrent jobs")
	concurrentGoroutines := make(chan struct{}, maxConcurrent)
	for i := 0; i < maxConcurrent; i++ {
		concurrentGoroutines <- struct{}{}
	}
	done := make(chan bool)
	waitForAllJobs := make(chan bool)
	go func() {
		for i := 0; i < len(jobArr); i++ {
			<-done
			concurrentGoroutines <- struct{}{}
		}
		waitForAllJobs <- true
	}()

	for i, job := range jobArr {
		<-concurrentGoroutines
		glog.Info(" Starting goroutine " + fmt.Sprint(i+1) + " of " + fmt.Sprint(len(jobArr)) + " " + job.ID)
		go func(job runner.Job, counter int) {

			if haddockVersion == 2 {

				_, errSetup2 := job.SetupHaddock24(inp.General.HaddockExecutable)

				if errSetup2 != nil {
					glog.Error("Failed to setup HADDOCK: " + errSetup2.Error())
					return
				}
				_, errRun2 := job.RunHaddock24(inp.General.HaddockExecutable)
				if errRun2 != nil {
					glog.Error("Failed to run HADDOCK: " + errRun2.Error())
					return
				}
			} else if haddockVersion == 3 {
				_, errRun3 := job.RunHaddock3(inp.General.HaddockExecutable)
				if errRun3 != nil {
					glog.Error("Failed to run HADDOCK: " + errRun3.Error())
					return
				}
			}

			done <- true
		}(job, i)
	}

	// Wait until all the jobs are done.
	<-waitForAllJobs

	glog.Info("All jobs completed successfully")

}

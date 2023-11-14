// main package is the entry point of the program
package main

import (
	"flag"
	"fmt"
	"haddockrunner/dataset"
	"haddockrunner/input"
	"haddockrunner/runner"
	"haddockrunner/runner/status"
	"haddockrunner/utils"
	"haddockrunner/utils/checksum"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/golang/glog"
)

const version = "v1.6.5"

func init() {
	var versionPrint bool
	const usage = `Usage: %s [options] <input file>

Run HADDOCK on a dataset of complexes

Options:
`
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "WARNING")
	_ = flag.Set("v", "2")

	flag.BoolVar(&versionPrint, "version", false, "Print version and exit")
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Fprintf(flag.CommandLine.Output(), usage, "`executable`")
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
		fmt.Printf("haddockrunner version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		glog.Error("No arguments were provided\n\n")
		flag.Usage()
		os.Exit(1)
	}
	args := os.Args[1:]
	inputF := args[0]

	glog.Info("###########################################")
	glog.Info(" Starting haddockrunner " + version)
	glog.Info("###########################################")
	glog.Info("Loading input file: " + inputF)

	inp, errInp := input.LoadInput(inputF)
	if errInp != nil {
		glog.Exit("Failed to load input file: " + errInp.Error())
	}

	errExec := inp.ValidateExecutable()
	if errExec != nil {
		glog.Exit("Failed to validate executable: " + errExec.Error())
	}

	errPatt := inp.ValidatePatterns()
	if errPatt != nil {
		glog.Exit("ERROR: " + errPatt.Error())
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
			glog.Exit("Failed to find run.cns-conf: " + errFind.Error())
		}

		runCnsParams, errLoad := input.LoadHaddock24Params(runCns)
		if errLoad != nil {
			glog.Exit("Failed to load run.cns-conf parameters" + errLoad.Error())
		}

		for _, scenario := range inp.Scenarios {
			scenarioParams := scenario.Parameters.CnsParams
			errValidate := input.ValidateRunCNSParams(runCnsParams, scenarioParams)
			if errValidate != nil {
				glog.Exit("Failed to validate scenario parameters: " + errValidate.Error())
			}
		}
	} else if haddockVersion == 3 {

		moduleArr, errParams := input.LoadHaddock3Params(inp.General.HaddockDir)
		if errParams != nil {
			glog.Exit("Failed to load HADDOCK3 parameters: " + errParams.Error())
		}

		for _, scenario := range inp.Scenarios {
			scenarioModules := scenario.Parameters.Modules
			errValidate := input.ValidateHaddock3Params(moduleArr, scenarioModules)
			if errValidate != nil {
				glog.Exit("Failed to validate scenario parameters: " + errValidate.Error())
			}
		}

	}

	// Load the dataset
	data, errDataset := dataset.LoadDataset(inp.General.WorkDir, inp.General.InputList, inp.General.ReceptorSuffix, inp.General.LigandSuffix)
	if errDataset != nil {
		glog.Exit("Failed to load dataset: " + errDataset.Error())
	}

	// Validate the checksum
	checksumF := filepath.Join(inp.General.WorkDir, "checksum.txt")
	_, errValidateChecksum := checksum.ValidateChecksum(inputF, inp.General.InputList, checksumF)
	if errValidateChecksum != nil {
		glog.Exit("Failed to validate checksum: " + errValidateChecksum.Error())
	}

	// Organize the dataset
	orgData, errOrganize := dataset.OrganizeDataset(inp.General.WorkDir, data)
	if errOrganize != nil {
		glog.Exit("Failed to organize dataset: " + errOrganize.Error())
	}

	// Setup the scenarios & create the jobs
	jobArr := []runner.Job{}

	for _, target := range orgData {
		glog.Info("Setting up target " + target.ID)

		for _, scenario := range inp.Scenarios {
			if haddockVersion == 2 {
				job, errSetup := target.SetupHaddock24Scenario(inp.General.WorkDir, inp.General.HaddockDir, scenario)
				if errSetup != nil {
					glog.Exit("Failed to setup scenario: " + errSetup.Error())
				}
				jobArr = append(jobArr, job)
			} else if haddockVersion == 3 {
				job, errSetup := target.SetupHaddock3Scenario(inp.General.WorkDir, scenario)
				if errSetup != nil {
					glog.Exit("Failed to setup scenario: " + errSetup.Error())
				}
				jobArr = append(jobArr, job)
			}
		}
	}

	// Sort the job array
	sort.Slice(jobArr, func(i, j int) bool {
		return jobArr[i].ID < jobArr[j].ID
	})

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
	glog.Info("############################################")
	for i, job := range jobArr {
		<-concurrentGoroutines
		// glog.Info("> Running " + job.ID + " (" + fmt.Sprint(i+1) + "/" + fmt.Sprint(total) + ")")
		go func(job runner.Job, counter int) {

			job.GetStatus(haddockVersion)

			switch {
			case job.Status == status.DONE:
				glog.Info(job.ID + " - " + job.Status + " - skipping")

			case job.Status == status.FAILED || job.Status == status.INCOMPLETE:
				glog.Warning("+++ " + job.ID + " is " + job.Status + " - restarting +++")
				// --------------------------------------------
				// TODO: Add the cleaning logic here
				os.RemoveAll(filepath.Join(job.Path, "run1"))
				// --------------------------------------------
				fallthrough

			default:
				// glog.Info("Job " + job.ID + " - " + job.Status)
				now := time.Now()
				_, runErr := job.Run(haddockVersion, inp.General.HaddockExecutable)
				if runErr != nil {
					glog.Exit("Failed to run HADDOCK: " + runErr.Error())
				}

				job.GetStatus(haddockVersion)
				elapsed := time.Since(now)
				glog.Info(job.ID + " - " + job.Status + " in " + fmt.Sprintf("%.2f", elapsed.Seconds()) + " seconds")
			}

			done <- true
		}(job, i)
	}

	// Wait until all the jobs are done.
	<-waitForAllJobs
	glog.Info("############################################")

	glog.Info("haddock-runner finished successfully")

}

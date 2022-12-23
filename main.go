// main package is the entry point of the program
package main

import (
	"benchmarktools/dataset"
	"benchmarktools/input"
	"benchmarktools/runner"
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

func main() {
	// Set the glog flags
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "WARNING")
	flag.Set("v", "2")

	flag.Parse()

	args := os.Args[1:]

	if len(args) == 0 {
		glog.Error("No arguments provided")
		return
	}
	glog.Info("#######################")
	glog.Info("Starting benchmarktools")
	glog.Info("#######################")
	glog.Info("Loading input file: " + args[0])

	inp, errInp := input.LoadInput(args[0])
	if errInp != nil {
		glog.Error("Failed to load input file: " + errInp.Error())
		return
	}

	errExec := inp.ValidateExecutable()
	if errExec != nil {
		glog.Error("Failed to validate executable: " + errExec.Error())
		return
	}

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
			job, errSetup := target.SetupScenario(inp.General.WorkDir, inp.General.HaddockDir, scenario)
			if errSetup != nil {
				glog.Error("Failed to setup scenario: " + errSetup.Error())
				return
			}
			jobArr = append(jobArr, job)
		}
	}

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
			// glog.Info("[" + job.ID + "]" + " Job " + fmt.Sprint(counter+1) + " of " + fmt.Sprint(len(jobArr)))
			// glog.Info("[" + job.ID + "]" + " dir: " + job.Path)
			_, errSetup := job.SetupHaddock(inp.General.HaddockExecutable)

			if errSetup != nil {
				glog.Error("Failed to setup HADDOCK: " + errSetup.Error())
				return
			}
			_, errRun := job.RunHaddock(inp.General.HaddockExecutable)
			if errRun != nil {
				glog.Error("Failed to run HADDOCK: " + errRun.Error())
				return
			}
			done <- true
		}(job, i)
	}

	// Wait until all the jobs are done.
	<-waitForAllJobs

	glog.Info("All jobs completed successfully")

}

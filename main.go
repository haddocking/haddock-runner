// main package is the entry point of the program
package main

import (
	"flag"
	"fmt"
	"haddockrunner/constants"
	"haddockrunner/dataset"
	"haddockrunner/input"
	"haddockrunner/runner"
	"haddockrunner/runner/status"
	"haddockrunner/utils"
	"haddockrunner/utils/checksum"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang/glog"
)

const version = "v2.0.0"

func init() {
	var versionPrint bool
	var setupOnly bool
	const usage = `Usage: %s [options] <input file>

Run HADDOCK on a dataset of complexes

Options:
`
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "WARNING")
	_ = flag.Set("v", "2")

	flag.BoolVar(&versionPrint, "version", false, "Print version and exit")
	flag.BoolVar(&setupOnly, "setup", false, "Only perform the setup, do not execute the benchmark")
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Fprintf(flag.CommandLine.Output(), usage, "haddock-runner")
		for _, f := range []string{"version", "setup"} {
			flag := flagSet.Lookup(f)
			fmt.Printf("  -%s: %s\n", f, flag.Usage)
		}
		fmt.Println("")
	}

	if os.Getenv("SLURM_JOB_ID") != "" {
		glog.Exit("haddock-runner cannot be run inside a SLURM job")
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	args := os.Args[1:]

	if utils.IsFlagPassed("version") {
		fmt.Printf("haddock-runner version %s\n", version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		glog.Error("No arguments were provided\n\n")
		flag.Usage()
		os.Exit(1)
	}

	setupOnly := false
	if utils.IsFlagPassed("setup") {
		glog.Info("`setup` argument passed, the benchmark will not be executed")
		args = utils.RemoveString(args, "-setup")
		setupOnly = true
	}

	inputF := args[0]

	glog.Info("###########################################")
	glog.Info(" Starting haddock-runner " + version)
	glog.Info("###########################################")
	glog.Info("Loading input file: " + inputF)

	inp, errInp := input.LoadInput(inputF)
	if errInp != nil {
		glog.Exit("Failed to load input file: " + errInp.Error())
	}

	// Check if the workdir exists and confirm overwrite
	if !utils.ConfirmOverwriteIfExists(inp.General.WorkDir, os.Stdin, os.Stdout) {
		glog.Exit("terminating...")
	}

	errExec := inp.ValidateExecutable()
	if errExec != nil {
		glog.Exit("Failed to validate executable: " + errExec.Error())
	}

	errPatt := inp.ValidatePatterns()
	if errPatt != nil {
		glog.Exit("ERROR: " + errPatt.Error())
	}

	errExecutionModes := inp.ValidateExecutionModes()
	if errExecutionModes != nil {
		glog.Exit("ERROR: " + errExecutionModes.Error())
	}

	// haddockVersion := inp.General.HaddockVersion
	var haddockVersion int
	if utils.IsHaddock24(inp.General.HaddockDir) {
		haddockVersion = 2
	} else {
		haddockVersion = 3
	}

	switch haddockVersion {
	case 2:
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
	case 3:
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
	data, errDataset := dataset.LoadDataset(inp.General.WorkDir, inp.General.InputList, inp.General.ReceptorSuffix, inp.General.LigandSuffix, inp.General.ShapeSuffix)
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
			switch haddockVersion {
			case 2:
				job, errSetup := target.SetupHaddock24Scenario(inp.General.WorkDir, inp.General.HaddockDir, scenario)
				if errSetup != nil {
					glog.Exit("Failed to setup scenario: " + errSetup.Error())
				}
				jobArr = append(jobArr, job)
			case 3:
				job, errSetup := target.SetupHaddock3Scenario(inp.General.WorkDir, scenario)
				if errSetup != nil {
					glog.Exit("Failed to setup scenario: " + errSetup.Error())
				}
				jobArr = append(jobArr, job)
			}
		}
	}

	if setupOnly {
		glog.Info("Benchmark setup finished successfully, exiting")
		os.Exit(0)
	}

	// Sort the job array
	jobArr = runner.SortJobs(jobArr)

	glog.Info("############################################")
	glog.Info("Running " + fmt.Sprint(len(jobArr)) + " jobs, " + fmt.Sprint(inp.General.MaxConcurrent) + " concurrent")
	glog.Info("############################################")

	semaphore := make(chan struct{}, inp.General.MaxConcurrent)
	finishedStatues := []string{status.DONE, status.FAILED, status.INCOMPLETE}

	var wg sync.WaitGroup

	for _, job := range jobArr {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(j runner.Job) {
			defer wg.Done()
			defer func() { <-semaphore }()

			_ = j.UpdateStatus(utils.GetJobID, utils.CheckSlurmStatus, haddockVersion)

			switch j.Status {

			case status.DONE:
				glog.Info(j.ID + " - " + j.Status + " - skipping")

			case status.FAILED, status.INCOMPLETE:
				glog.Warning("+++ " + j.ID + " is " + j.Status + " - restarting +++")
				err := j.Clean()
				if err != nil {
					glog.Exit("Failed to clean job: " + err.Error())
				}
				err = j.Post(haddockVersion, inp.General.HaddockExecutable, inp.Slurm)
				if err != nil {
					glog.Exit("Failed to post job: " + err.Error())
				}
				err = j.WaitUntil(finishedStatues, constants.WAIT_TIMEOUT_COUNTER)
				if err != nil {
					glog.Exit("Failed to wait for job: " + err.Error())
				}

			case status.SUBMITTED:
				glog.Info(j.ID + " - " + j.Status + " - waiting")
				err := j.WaitUntil(finishedStatues, constants.WAIT_TIMEOUT_COUNTER)
				if err != nil {
					glog.Exit("Failed to wait for job: " + err.Error())
				}

			default:
				err := j.Post(haddockVersion, inp.General.HaddockExecutable, inp.Slurm)
				if err != nil {
					glog.Exit("Failed to post job: " + err.Error())
				}
				err = j.WaitUntil(finishedStatues, constants.WAIT_TIMEOUT_COUNTER)
				if err != nil {
					glog.Exit("Failed to wait for job: " + err.Error())
				}
			}
		}(job)
	}

	wg.Wait()

	glog.Info("############################################")
	glog.Info("haddock-runner finished successfully")
}

// main package is the entry point of the program
package main

import (
	"benchmarktools/dataset"
	"benchmarktools/input"
	"benchmarktools/runner"
	"fmt"
	"os"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	inp, errInp := input.LoadInput(args[0])
	if errInp != nil {
		fmt.Println("Failed to load input file")
		return
	}

	errExec := inp.ValidateExecutable()
	if errExec != nil {
		fmt.Println("Failed to validate executable" + errExec.Error())
		return
	}

	runCns, errFind := input.FindHaddock24RunCns(inp.General.HaddockDir)
	if errFind != nil {
		fmt.Println("Failed to find run.cns-conf")
		return
	}

	runCnsParams, errLoad := input.LoadHaddock24Params(runCns)
	if errLoad != nil {
		fmt.Println("Failed to load run.cns-conf parameters")
		return
	}

	for _, scenario := range inp.Scenarios {
		scenarioParams := scenario.Parameters.CnsParams
		errValidate := input.ValidateRunCNSParams(runCnsParams, scenarioParams)
		if errValidate != nil {
			fmt.Println("Failed to validate scenario parameters: " + errValidate.Error())
			return
		}
	}

	// Load the dataset
	data, errDataset := dataset.LoadDataset(inp.General.WorkDir, inp.General.InputList, inp.General.ReceptorSuffix, inp.General.LigandSuffix)
	if errDataset != nil {
		fmt.Println("Failed to load dataset: " + errDataset.Error())
		return
	}
	fmt.Println(data)

	// Organize the dataset
	orgData, errOrganize := dataset.OrganizeDataset(inp.General.WorkDir, data)
	if errOrganize != nil {
		fmt.Println("Failed to organize dataset: " + errOrganize.Error())
		return
	}

	// Setup the scenarios & create the jobs
	jobArr := []runner.Job{}

	for _, target := range orgData {
		for _, scenario := range inp.Scenarios {
			fmt.Println("Setting up scenario " + scenario.Name)
			job, errSetup := target.SetupScenario(inp.General.WorkDir, inp.General.HaddockDir, scenario)
			if errSetup != nil {
				fmt.Println("Failed to setup scenario: " + errSetup.Error())
				return
			}
			jobArr = append(jobArr, job)
		}
	}

	// Run the jobs
	for _, job := range jobArr {

		_, errSetup := job.SetupHaddock(inp.General.HaddockExecutable)

		if errSetup != nil {
			fmt.Println("Failed to setup HADDOCK: " + errSetup.Error())
			return
		}

		_, errRun := job.RunHaddock(inp.General.HaddockExecutable)
		if errRun != nil {
			fmt.Println("Failed to run HADDOCK: " + errRun.Error())
			return
		}

	}

	fmt.Println(jobArr)

}

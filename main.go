// main package is the entry point of the program
package main

import (
	"benchmarktools/dataset"
	"benchmarktools/input"
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
		errValidate := scenario.ValidateScenarioParams(runCnsParams)
		if errValidate != nil {
			fmt.Println("Failed to validate scenario parameters: " + errValidate.Error())
			return
		}
	}

	// Load the dataset
	data, errDataset := dataset.LoadDataset(inp.General.InputPDBList, inp.General.ReceptorSuffix, inp.General.LigandSuffix)
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

	// Setup the scenarios
	for _, target := range orgData {
		// for _, scenario := range inp.Scenarios {
		errSetup := target.SetupScenarios(inp)
		if errSetup != nil {
			fmt.Println("Failed to setup scenario: " + errSetup.Error())
			return
		}
		// }
	}

	// Dev Only -
	// defer os.RemoveAll(bmPath)
	// End Dev Only

	fmt.Println(orgData)

}

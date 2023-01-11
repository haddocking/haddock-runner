package input

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadInput(t *testing.T) {
	var err error
	// Pass by being able to load the input file
	// Create a OK input file
	inp := Input{
		General: GeneralStruct{
			HaddockExecutable: "haddock.sh",
			HaddockDir:        "haddock_dir",
			ReceptorSuffix:    "r_u",
			LigandSuffix:      "l_u",
		},
		Scenarios: []Scenario{
			{
				Name: "scenario1",
				Parameters: ScenarioParams{
					CnsParams: map[string]interface{}{
						"param1": false,
						"param2": "string",
						"param3": 1,
						"param4": 1.5,
					},
					Restraints: Restraints{
						Ambig:   "ambig",
						Unambig: "unambig",
					},
					Toppar: Toppar{
						Topology: "top",
						Param:    "param",
					},
				},
			},
		},
	}

	yamlData, err := yaml.Marshal(&inp)
	if err != nil {
		t.Errorf("Failed to marshal input file: %s", err)
	}

	err = os.WriteFile("input.yml", yamlData, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}
	defer os.Remove("input.yml")

	_, err = LoadInput("input.yml")
	if err != nil {
		t.Errorf("Failed to load input file: %s", err)
	}

	// Fail by not being able to load the input file
	err = os.WriteFile("wrong_input.yml", []byte("not a yml"), 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}
	defer os.Remove("wrong_input.yml")

	_, err = LoadInput("wrong_input.yml")
	if err == nil {
		t.Errorf("Failed to detect wrong input file")
	}

	// Fail by trying to load a file that does not exist
	_, err = LoadInput("does_not_exist.yml")
	if err == nil {
		t.Errorf("Failed to detect wrong input file")
	}

}

func TestValidatePatterns(t *testing.T) {

	inp := Input{
		General: GeneralStruct{
			ReceptorSuffix: "r_u",
			LigandSuffix:   "l_u",
		},
		Scenarios: []Scenario{
			{
				Name: "scenario1",
				Parameters: ScenarioParams{
					CnsParams:  map[string]interface{}{},
					Restraints: Restraints{},
					Toppar: Toppar{
						Topology: "top",
						Param:    "param",
					},
					Modules: ModuleParams{
						Rigidbody: map[string]interface{}{
							"param1_fname": "pattern1",
							"param2_fname": "pattern2",
						},
					},
					General: map[string]interface{}{},
				},
			},
		},
	}

	// Pass by validating the patterns
	err := inp.ValidatePatterns()
	if err != nil {
		t.Errorf("Failed to validate patterns: %s", err)
	}

	// Fail by not validating the patterns

	// Fail because the fname parameters are repeated
	badInp := Input{
		General: GeneralStruct{
			ReceptorSuffix: "r_u",
			LigandSuffix:   "l_u",
		},
		Scenarios: []Scenario{
			{
				Name: "scenario1",
				Parameters: ScenarioParams{
					CnsParams:  map[string]interface{}{},
					Restraints: Restraints{},
					Toppar:     Toppar{},
					Modules: ModuleParams{
						Order: []string{"rigidbody"},
						Rigidbody: map[string]interface{}{
							"param1_fname": "pattern1",
							"param2_fname": "pattern1",
						},
					},
					General: map[string]interface{}{},
				},
			},
		},
	}

	err = badInp.ValidatePatterns()
	if err == nil {
		t.Errorf("Failed to detect repeated patterns")
	}

	// fail because the receptor and ligand suffixes are not defined
	badInp.General.ReceptorSuffix = ""
	badInp.General.LigandSuffix = ""

	err = badInp.ValidatePatterns()
	if err == nil {
		t.Errorf("Failed to detect missing receptor and ligand suffixes")
	}
	badInp.General.ReceptorSuffix = "p1"
	badInp.General.LigandSuffix = "p2"

	// fail because the receptor/ligand suffixex are the same
	badInp.General.ReceptorSuffix = "p1"
	badInp.General.LigandSuffix = "p1"

	err = badInp.ValidatePatterns()
	if err == nil {
		t.Errorf("Failed to detect same receptor and ligand suffixes")
	}

	badInp.General.ReceptorSuffix = "p1"
	badInp.General.LigandSuffix = "p2"

	// Fail because the Ambig/Unambig patterns are the same
	badInp.Scenarios[0].Parameters.Restraints.Ambig = "p1"
	badInp.Scenarios[0].Parameters.Restraints.Unambig = "p1"

	err = badInp.ValidatePatterns()
	if err == nil {
		t.Errorf("Failed to detect same Ambig and Unambig patterns")
	}

	badInp.Scenarios[0].Parameters.Restraints.Ambig = "p1"
	badInp.Scenarios[0].Parameters.Restraints.Unambig = "p2"

	// Fail because Topology and Param are the same
	badInp.Scenarios[0].Parameters.Toppar.Topology = "p1"
	badInp.Scenarios[0].Parameters.Toppar.Param = "p1"

	err = badInp.ValidatePatterns()
	if err == nil {
		t.Errorf("Failed to detect same Topology and Param patterns")
	}

	badInp.Scenarios[0].Parameters.Toppar.Topology = "p1"
	badInp.Scenarios[0].Parameters.Toppar.Param = "p2"

}
func TestValidateHaddock3Executable(t *testing.T) {

}

func TestFindHaddock24RunCns(t *testing.T) {

	// Based on the executable, return the location of run.cns

	// Create an executable and place it two levels above run.cns
	haddockDir := "_test"
	protocolsDir := "_test/protocols"
	_ = os.MkdirAll(haddockDir, 0755)
	defer os.RemoveAll(haddockDir)

	_ = os.Mkdir(protocolsDir, 0755)
	runCnsF := "_test/protocols/run.cns-conf"
	err := os.WriteFile(runCnsF, []byte("{===>} parameter=\"value\";"), 0755)
	if err != nil {
		t.Errorf("Failed to write run.cns: %s", err)
	}

	// Pass by finding the run.cns file
	_, err = FindHaddock24RunCns(haddockDir)
	if err != nil {
		t.Errorf("Failed to find run.cns: %s", err)
	}

	// Fail by not finding the run.cns file
	_, err = FindHaddock24RunCns("does_not_exist")
	if err == nil {
		t.Errorf("Failed to detect wrong executable")
	}

}

func TestLoadHaddock24Params(t *testing.T) {
	// Parse the run.cns file and return the parameters as ParameterStruct
	params := []byte(
		"{===>} parameter1=\"value\";\n" +
			"{===>} parameter2=1;\n" +
			"{===>} parameter3=1.0;\n" +
			"{===>} parameter4=true;\n")
	err := os.WriteFile("_test_run.cns-conf", params, 0755)
	if err != nil {
		t.Errorf("Failed to write run.cns: %s", err)
	}
	defer os.Remove("_test_run.cns-conf")

	// Pass by finding the parameters
	p, err := LoadHaddock24Params("_test_run.cns-conf")
	if err != nil {
		t.Errorf("Failed to load parameters: %s", err)
	}

	if p["parameter1"] != "value" {
		t.Errorf("Failed to parse parameter1")
	}

	if p["parameter2"] != 1 {
		t.Errorf("Failed to parse parameter2")
	}

	if p["parameter3"] != 1.0 {
		t.Errorf("Failed to parse parameter3")
	}

	if p["parameter4"] != true {
		t.Errorf("Failed to parse parameter4")
	}

	// Fail by not finding the parameters
	_, err = LoadHaddock24Params("does_not_exist")
	if err == nil {
		t.Errorf("Failed to detect wrong executable")
	}

}

func TestValidateRunCNSParams(t *testing.T) {

	valid := map[string]interface{}{
		"param1": true,
		"param2": 1,
		"param3": 1.5,
		"param4": "string",
	}

	// Check if the input parameters of the scenario are valid
	params := map[string]any{
		"param1": true,
	}

	err := ValidateRunCNSParams(valid, params)
	if err != nil {
		t.Errorf("Failed to validate parameters: %s", err)
	}

	// Fail by not finding the parameters
	valid = map[string]any{
		"param1": true,
	}

	params = map[string]any{
		"param2": true,
	}

	err = ValidateRunCNSParams(valid, params)
	if err == nil {
		t.Errorf("Failed to detect wrong parameters")
	}

}

func TestLoadHaddock3DefaultParams(t *testing.T) {

	// Create a folder structure and fill it with dummy files

	rootPath := "_haddock3"
	modulePath := filepath.Join(rootPath, "/src/haddock/modules/")
	_ = os.MkdirAll(modulePath, 0755)
	defer os.RemoveAll(rootPath)

	type dummyParams struct {
		Default string
	}

	moduleNames := []string{"topoaa",
		"topocg",
		"exit",
		"emref",
		"flexref",
		"mdref",
		"gdock",
		"lightdock",
		"rigidbody",
		"emscoring",
		"mdscoring",
		"caprieval",
		"clustfcc",
		"clustrmsd",
		"rmsdmatrix",
		"seletop",
		"seletopclusts"}

	for _, mod := range moduleNames {
		_ = os.MkdirAll(filepath.Join(modulePath, mod), 0755)
		defaultsF := filepath.Join(modulePath, mod, "defaults.yaml")
		params := map[string]dummyParams{
			"param1": {"value1"},
		}
		data, err := yaml.Marshal(&params)
		if err != nil {
			t.Errorf("Failed to marshal parameters: %s", err)
		}
		err = os.WriteFile(defaultsF, data, 0755)
		if err != nil {
			t.Errorf("Failed to write defaults.yaml: %s", err)
		}
	}

	// Pass by finding the parameters
	_, err := LoadHaddock3Params(rootPath)
	if err != nil {
		t.Errorf("Failed to load parameters: %s", err)
	}

	// Fail by not finding the parameters
	_, err = LoadHaddock3Params("does_not_exist")
	if err == nil {
		t.Errorf("Failed to load parameters")
	}

	// Fail by trying to unmarshal a malformed file
	// defaultsF := filepath.Join(modulePath, "rigidbody", "defaults.yaml")
	wrongParams := filepath.Join(modulePath, "rigidbody", "wrong_params.yaml")
	err = os.WriteFile(wrongParams, []byte("not a yaml file"), 0755)
	if err != nil {
		t.Errorf("Failed to write defaults.yaml: %s", err)
	}
	_, err = LoadHaddock3Params(rootPath)
	if err == nil {
		t.Errorf("Failed to load parameters: %s", err)
	}

}

func TestValidateHaddock3Params(t *testing.T) {

	known := ModuleParams{}
	known.Rigidbody = map[string]any{
		"param1": "value1",
	}
	known.Topoaa = map[string]any{
		"param2": "value2",
	}

	test := ModuleParams{}
	test.Rigidbody = map[string]any{
		"param1": "value1",
	}
	test.Topoaa = map[string]any{
		"param2": "value2",
	}

	// Pass by finding the parameters
	err := ValidateHaddock3Params(known, test)
	if err != nil {
		t.Errorf("Failed to validate parameters: %s", err)
	}

	// Fail by not finding a parameter
	test.Rigidbody = map[string]any{
		"param10": "value",
	}

	err = ValidateHaddock3Params(known, test)
	if err == nil {
		t.Errorf("Failed to detect wrong parameters")
	}

}

func TestInput_ValidateExecutable(t *testing.T) {
	// Create a dummy executable
	wd, _ := os.Getwd()
	haddockDir := filepath.Join(wd, "_test")
	_ = os.MkdirAll(haddockDir, 0755)

	defer os.RemoveAll(haddockDir)

	haddockF := filepath.Join(haddockDir, "/haddock.sh")
	err := os.WriteFile(haddockF, []byte("#!/bin/bash"), 0755)
	if err != nil {
		t.Errorf("Failed to write executable: %s", err)
	}
	nonExecHaddockF := filepath.Join(haddockDir, "/nonExec-haddock.sh")
	err = os.WriteFile(nonExecHaddockF, []byte("#!/bin/bash"), 0644)
	if err != nil {
		t.Errorf("Failed to write executable: %s", err)
	}
	nonExistHaddockF := filepath.Join(haddockDir, "/does-not-exist.sh")

	type fields struct {
		General   GeneralStruct
		Scenarios []Scenario
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: haddockF,
				},
			},
			wantErr: false,
		},
		{
			name: "relative path",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: "haddock.sh",
				},
			},
			wantErr: true,
		},
		{
			name: "undefined",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: "",
				},
			},
			wantErr: true,
		},
		{
			name: "not executable",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: nonExecHaddockF,
				},
			},
			wantErr: true,
		},
		{
			name: "non-existing",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: nonExistHaddockF,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inp := &Input{
				General:   tt.fields.General,
				Scenarios: tt.fields.Scenarios,
			}
			if err := inp.ValidateExecutable(); (err != nil) != tt.wantErr {
				t.Errorf("Input.ValidateExecutable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

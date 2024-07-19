package input

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"gopkg.in/yaml.v3"
)

func TestLoadInput(t *testing.T) {
	// Create a OK input file
	d1 := []byte(`general:
  executable: /home/rodrigo/repos/haddock-runner/haddock3.sh
  max_concurrent: 999
  haddock_dir: ../haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: example/input_list.txt
  work_dir: bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
      restraints:
        ambig: ti
      custom_toppar:
        topology: _ligand.top
      general:
        ncores: 1
      modules:
        order: [rigidbody]
        rigidbody:
          param1: value1
`)

	d2 := []byte(`general:
  executable: /home/rodrigo/repos/haddock-runner/haddock3.sh
  max_concurrent: 999
  haddock_dir: ../haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: example/input_list.txt
  work_dir: bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
      restraints:
        ambig: ti
      custom_toppar:
        topology: _ligand.top
      general:
        ncores: 1
      modules:
        order: [rigidbody, rigidbody]
        rigidbody:
          param1: value1
        rigidbody:
          param1: value1
`)

	d3 := []byte(`general:
  executable: /home/rodrigo/repos/haddock-runner/haddock3.sh
  max_concurrent: 999
  haddock_dir: ../haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: example/input_list.txt
  work_dir: bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
      restraints:
        ambig: ti
      custom_toppar:
        topology: _ligand.top
      general:
        ncores: 1
      modules:
        order: [rigidbody, caprieval.1, caprieval.2]
        rigidbody:
          param1: value1
        caprieval.1:
          param1: value1
        caprieval.2:
          param1: value1
`)

	// Write a file with unknown modules
	// Fake a haddock3 directory structure with a defaults.yaml file
	temp_dir := "_testloadinput_haddock3"
	_ = os.MkdirAll(filepath.Join(temp_dir, "src/haddock/modules"), 0755)
	_ = os.WriteFile(filepath.Join(temp_dir, "src/haddock/modules/defaults.yaml"), []byte(""), 0755)
	// Create a known module
	// _ = os.MkdirAll(filepath.Join(temp_dir, "src/haddock/modules/rigidbody"), 0755)
	defer os.RemoveAll(temp_dir)

	d4 := []byte(`general:
  executable: /home/rodrigo/repos/haddock-runner/haddock3.sh
  max_concurrent: 999
  haddock_dir: _testloadinput_haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: example/input_list.txt
  work_dir: bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
      restraints:
        ambig: ti
      custom_toppar:
        topology: _ligand.top
      general:
        ncores: 1
      modules:
        order: [unknown, caprieval.1, caprieval.2]
        unknown:
          param1: value1
        caprieval.1:
          param1: value1
        caprieval.2:
          param1: value1
`)

	d5 := []byte(`general:
  executable: /home/rodrigo/repos/haddock-runner/haddock3.sh
  max_concurrent: 999
  haddock_dir: _testloadinput_haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: example/input_list.txt
  work_dir: bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
      restraints:
        ambig: ti
      custom_toppar:
        topology: _ligand.top
      general:
        ncores: 1
      modules:
        order: [rigidbody, unknown, unknownalso]
        rigidbody:
          param1: value1
        unknown:
          param1: value1
        unknownalso:
          param1: value1
`)

	err := os.WriteFile("test-input.yaml", d1, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}

	defer os.Remove("test-input.yaml")

	err = os.WriteFile("test-input-wrong.yaml", d2, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}

	defer os.Remove("test-input-wrong.yaml")

	err = os.WriteFile("test-input-repeated.yaml", d3, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}

	defer os.Remove("test-input-repeated.yaml")

	err = os.WriteFile("test-input-unknown.yaml", d4, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}

	defer os.Remove("test-input-unknown.yaml")

	err = os.WriteFile("test-input-unknown-twice.yaml", d5, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}

	defer os.Remove("test-input-unknown-twice.yaml")

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Input
		wantErr bool
	}{
		{
			name: "valid-input",
			args: args{
				filename: "test-input.yaml",
			},
			want: &Input{
				General: GeneralStruct{
					HaddockExecutable: "/home/rodrigo/repos/haddock-runner/haddock3.sh",
					MaxConcurrent:     999,
					HaddockDir:        "../haddock3",
					ReceptorSuffix:    "_r_u",
					LigandSuffix:      "_l_u",
					InputList:         "example/input_list.txt",
					WorkDir:           "bm-goes-here",
				},
				Scenarios: []Scenario{
					{
						Name: "true-interface",
						Parameters: ParametersStruct{
							General: map[string]any{
								"ncores": 1,
							},
							Restraints: Airs{
								Ambig: "ti",
							},
							Toppar: TopologyParams{
								Topology: "_ligand.top",
							},
							Modules: ModuleParams{
								Order: []string{"rigidbody"},
								Rigidbody: map[string]any{
									"param1": "value1",
								},
							},
							CnsParams: map[string]interface{}{
								"noecv": false,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "non-existing",
			args: args{
				filename: "does-not-exist.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong-input",
			args: args{
				filename: "test-input-wrong.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repeated-modules",
			args: args{
				filename: "test-input-repeated.yaml",
			},
			// [rigidbody, caprieval.1, caprieval.2]
			want: &Input{
				General: GeneralStruct{
					HaddockExecutable: "/home/rodrigo/repos/haddock-runner/haddock3.sh",
					MaxConcurrent:     999,
					HaddockDir:        "../haddock3",
					ReceptorSuffix:    "_r_u",
					LigandSuffix:      "_l_u",
					InputList:         "example/input_list.txt",
					WorkDir:           "bm-goes-here",
				},
				Scenarios: []Scenario{
					{
						Name: "true-interface",
						Parameters: ParametersStruct{
							General: map[string]any{
								"ncores": 1,
							},
							Restraints: Airs{
								Ambig: "ti",
							},
							Toppar: TopologyParams{
								Topology: "_ligand.top",
							},
							Modules: ModuleParams{
								Order: []string{"rigidbody", "caprieval.1", "caprieval.2"},
								Rigidbody: map[string]any{
									"param1": "value1",
								},
								Caprieval_1: map[string]any{
									"param1": "value1",
								},
								Caprieval_2: map[string]any{
									"param1": "value1",
								},
							},
							CnsParams: map[string]interface{}{
								"noecv": false,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "unknown module",
			args: args{
				filename: "test-input-unknown.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unknown modules",
			args: args{
				filename: "test-input-unknown-twice.yaml",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong type",
			args: args{
				filename: "test-input-wrong-type.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadInput(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error(diff)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadInput() \n\ngot:   %v\n, \n\n want: %v\n", got, tt.want)
			}
		})
	}
}

func TestInput_ValidatePatterns(t *testing.T) {
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
			name: "valid",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
					LigandSuffix:   "_l_u",
				},
				Scenarios: []Scenario{
					{
						Name: "true-interface",
						Parameters: ParametersStruct{
							Restraints: Airs{
								Ambig: "ti",
							},
							Modules: ModuleParams{
								Order: []string{"rigidbody"},
								Rigidbody: map[string]any{
									"something_fname": "pattern1",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid-empty-receptor-suffix",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "",
				},
				Scenarios: []Scenario{},
			},
			wantErr: true,
		},
		{
			name: "invalid-receptor-ligand-equal",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
					LigandSuffix:   "_r_u",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid-ambig-unambig-equal",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
				},
				Scenarios: []Scenario{
					{
						Name: "",
						Parameters: ParametersStruct{
							Restraints: Airs{
								Ambig:   "same",
								Unambig: "same",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid-topology-params-equal",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
				},
				Scenarios: []Scenario{
					{
						Name: "",
						Parameters: ParametersStruct{
							Toppar: TopologyParams{
								Topology: "same",
								Param:    "same",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid-multiple-fnames",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
				},
				Scenarios: []Scenario{
					{
						Name: "",
						Parameters: ParametersStruct{
							Modules: ModuleParams{
								Order: []string{"rigidbody"},
								Rigidbody: map[string]any{
									"something_fname":      "pattern1",
									"something_else_fname": "pattern1",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid-fnames-match-each-other",
			fields: fields{
				General: GeneralStruct{
					ReceptorSuffix: "_r_u",
				},
				Scenarios: []Scenario{
					{
						Name: "",
						Parameters: ParametersStruct{
							Modules: ModuleParams{
								Order: []string{"rigidbody"},
								Rigidbody: map[string]any{
									"something_fname": "_ti",
								},
							},
						},
					},
					{
						Name: "",
						Parameters: ParametersStruct{
							Modules: ModuleParams{
								Order: []string{"rigidbody"},
								Rigidbody: map[string]any{
									"something_fname": "_ti_something",
								},
							},
						},
					},
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
			if err := inp.ValidatePatterns(); (err != nil) != tt.wantErr {
				t.Errorf("Input.ValidatePatterns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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

func TestValidateExecutionModes(t *testing.T) {

	// Setup a haddock3 folder structure, it must have this subdirectories:
	wd, _ := os.Getwd()
	haddock3Dir := filepath.Join(wd, "TestValidateExecutionModes")
	_ = os.MkdirAll(haddock3Dir, 0755)
	defer os.RemoveAll("TestValidateExecutionModes")

	// Add the subdirectories
	_ = os.MkdirAll(filepath.Join(haddock3Dir, "src/haddock/modules"), 0755)

	// Add an empty defaults.yaml file
	defaultsF := filepath.Join(haddock3Dir, "src/haddock/modules/defaults.yaml")
	err := os.WriteFile(defaultsF, []byte(""), 0755)
	if err != nil {
		t.Errorf("Failed to write defaults.yaml: %s", err)
	}

	// Setup a haddock2 folder structure, it must have this subdirectories
	// protocols/run.cns-conf
	haddock2Dir := filepath.Join(wd, "TestValidateExecutionModes2")
	_ = os.MkdirAll(haddock2Dir, 0755)
	defer os.RemoveAll("TestValidateExecutionModes2")

	// Add the subdirectories
	_ = os.MkdirAll(filepath.Join(haddock2Dir, "protocols"), 0755)

	// Add an empty run.cns-conf file
	runCnsF := filepath.Join(haddock2Dir, "protocols/run.cns-conf")
	err = os.WriteFile(runCnsF, []byte(""), 0755)
	if err != nil {
		t.Errorf("Failed to write run.cns-conf: %s", err)
	}

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
			name: "valid",
			fields: fields{
				General: GeneralStruct{
					HaddockDir: haddock3Dir,
				},
				Scenarios: []Scenario{
					{
						Name: "true-interface",
						Parameters: ParametersStruct{
							General: map[string]any{
								"mode": "local",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid-haddock2",
			fields: fields{
				General: GeneralStruct{
					HaddockDir: haddock2Dir,
				},
				Scenarios: []Scenario{},
			},
			wantErr: true,
		},
		{
			name: "invalid-haddock3",
			fields: fields{
				General: GeneralStruct{
					HaddockDir: haddock3Dir,
				},
				Scenarios: []Scenario{
					{
						Name: "true-interface",
						Parameters: ParametersStruct{
							General: map[string]any{
								"mode": "anything",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		inp := &Input{
			General:   tt.fields.General,
			Scenarios: tt.fields.Scenarios,
		}
		if err := inp.ValidateExecutionModes(); (err != nil) != tt.wantErr {
			t.Errorf("Input.ValidateExecutionModes() error = %v, wantErr %v", err, tt.wantErr)
		}
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
	userExecHaddockF := filepath.Join(haddockDir, "/userExec-haddock.sh")
	err = os.WriteFile(userExecHaddockF, []byte("#!/bin/bash"), 0744)
	if err != nil {
		t.Errorf("Failed to write executable: %s", err)
	}
	groupExecHaddockF := filepath.Join(haddockDir, "/groupExec-haddock.sh")
	err = os.WriteFile(groupExecHaddockF, []byte("#!/bin/bash"), 0654)
	if err != nil {
		t.Errorf("Failed to write executable: %s", err)
	}
	publicExecHaddockF := filepath.Join(haddockDir, "/publicExec-haddock.sh")
	err = os.WriteFile(publicExecHaddockF, []byte("#!/bin/bash"), 0645)
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
			name: "valid-executable",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: haddockF,
				},
			},
			wantErr: false,
		},
		{
			name: "non-executable",
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
		{
			name: "empty",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: "",
				},
			},
			wantErr: true,
		},
		{
			name: "relative-path",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: "haddock.sh",
				},
			},
			wantErr: true,
		},
		{
			name: "user-executable",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: userExecHaddockF,
				},
			},
			wantErr: false,
		},
		{
			name: "group-executable",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: groupExecHaddockF,
				},
			},
			wantErr: false,
		},
		{
			name: "public-executable",
			fields: fields{
				General: GeneralStruct{
					HaddockExecutable: publicExecHaddockF,
				},
			},
			wantErr: false,
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

func TestValidateHaddock3Params(t *testing.T) {
	type args struct {
		known  ModuleParams
		loaded ModuleParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				known: ModuleParams{
					Rigidbody: map[string]any{
						"param1": "value1",
					},
					Topoaa: map[string]any{
						"param2": "value2",
					},
				},
				loaded: ModuleParams{
					Rigidbody: map[string]any{
						"param1": "value1",
					},
					Topoaa: map[string]any{
						"param2": "value2",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				known: ModuleParams{
					Rigidbody: map[string]any{
						"param1": "value1",
					},
					Topoaa: map[string]any{
						"param2": "value2",
					},
				},
				loaded: ModuleParams{
					Rigidbody: map[string]any{
						"param10": "value",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "pass-expandable-param",
			args: args{
				known: ModuleParams{
					Rigidbody: map[string]any{
						"param1_1": "value1",
					},
					Topoaa: map[string]any{
						"param2_": "value2",
					},
				},
				loaded: ModuleParams{
					Rigidbody: map[string]any{
						"param1_50": "value1",
					},
					Topoaa: map[string]any{
						"param2_1": "value2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateHaddock3Params(tt.args.known, tt.args.loaded); (err != nil) != tt.wantErr {
				t.Errorf("ValidateHaddock3Params() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

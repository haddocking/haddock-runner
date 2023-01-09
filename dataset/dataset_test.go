package dataset

import (
	"benchmarktools/input"
	"benchmarktools/runner"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteRunParamStub(t *testing.T) {

	_ = os.Mkdir("1abc", 0755)
	defer os.RemoveAll("1abc")

	cwd, _ := os.Getwd()
	projectDir := filepath.Join(cwd, "1abc")
	_ = os.MkdirAll(projectDir, 0755)
	defer os.RemoveAll(projectDir)
	haddockDir := "/home/haddock2.4"

	// Pass by writing a valid run parameter file
	target := Target{
		ID:           "1abc",
		Receptor:     []string{"1abc_r_1.pdb", "1abc_r_2.pdb"},
		ReceptorList: "some/path/receptor.list",
		Ligand:       []string{"1abc_l_1.pdb", "1abc_l_2.pdb"},
		LigandList:   "some/path/ligand.list",
	}

	_, err := target.WriteRunParamStub(projectDir, haddockDir)
	if err != nil {
		t.Error(err)
	}

	// Fail by writing a run parameter file with an empty target
	target = Target{}

	_, err = target.WriteRunParamStub(projectDir, haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty receptor
	target = Target{
		ID:       "1abc",
		Receptor: []string{},
		Ligand:   []string{"1abc.pdb"},
	}

	_, err = target.WriteRunParamStub(projectDir, haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty ligand
	target = Target{
		ID:       "1abc",
		Receptor: []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:   []string{},
	}

	_, err = target.WriteRunParamStub(projectDir, haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty haddock directory
	target = Target{
		ID:       "1abc",
		Receptor: []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:   []string{"1abc.pdb"},
	}

	_, err = target.WriteRunParamStub(projectDir, "")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty project directory
	target = Target{
		ID:       "1abc",
		Receptor: []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:   []string{"1abc.pdb"},
	}

	_, err = target.WriteRunParamStub("", haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by not being able to write to the runparam file
	target = Target{
		ID:       "1abc",
		Receptor: []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:   []string{"1abc.pdb"},
	}

	_, err = target.WriteRunParamStub("", haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by not being able to write to the run.param because the project directory does not exist
	target = Target{
		ID:       "1abc",
		Receptor: []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:   []string{"1abc.pdb"},
	}

	_, err = target.WriteRunParamStub("some/path/that/does/not/exist", haddockDir)
	if err == nil {
		t.Error("Expected error, got nil")
	}

}

func TestLoadDataset(t *testing.T) {

	var err error

	fileArr := []string{"structure1_r_u.pdb", "structure1_l_u.pdb"}

	for _, file := range fileArr {
		err := os.WriteFile(file, []byte(""), 0644)
		if err != nil {
			t.Errorf("Failed to write structure file: %s", err)
		}
		defer os.Remove(file)
	}
	projectDir := "some/path"

	defer os.RemoveAll("some")
	defer os.Remove("structure3_receptor.list")
	defer os.Remove("structure4_ligand.list")

	// Write a valid list
	err = os.WriteFile("pdb.list",
		[]byte(
			"some/path/structure1_r_u.pdb\n"+
				"some/path/structure1_l_u.pdb\n"+
				"some/path/structure1_ref.pdb\n"+
				"some/path/structure2_l_u.pdb\n"+
				"some/path/structure2_r_u.pdb\n"+
				"some/path/structure3_r_u_0.pdb\n"+
				"some/path/structure3_r_u_1.pdb\n"+
				"some/path/structure3_r_u_2.pdb\n"+
				"some/path/structure3_l_u.pdb\n"+
				"some/path/structure4_r_u.pdb\n"+
				"some/path/structure4_l_u_0.pdb\n"+
				"some/path/structure4_l_u_1.pdb\n"+
				"some/path/structure4_ambig.tbl\n"+
				"some/path/structure4_ATP.top\n"+
				"some/path/structure4_ATP.param\n"), 0644)
	defer os.Remove("pdb.list")

	// Pass by loading a valid dataset

	tArr, errData := LoadDataset(projectDir, "pdb.list", "_r_u", "_l_u")
	if errData != nil {
		t.Errorf("Failed to load dataset: %s", err.Error())
	}

	if len(tArr) != 4 {
		t.Errorf("Failed to load dataset: %d", len(tArr))
	}

	for _, v := range tArr {
		if v.ID == "structure4" {
			if len(v.Ligand) != 2 {
				t.Errorf("Failed: Not all ligands were loaded")
			}
			if len(v.Restraints) != 1 {
				t.Errorf("Failed: Not all restraints were loaded")
			}
			if len(v.Toppar) != 2 {
				t.Errorf("Failed: Not all toppar files were loaded")
			}
		}
		if v.ID == "structure3" {
			if len(v.Receptor) != 3 {
				t.Errorf("Failed: Not all receptors were loaded")
			}
		}
	}

	// Fail by loading a dataset with a wrong file
	_, errData = LoadDataset(projectDir, "wrong_file.pdb", "_r_u", "_l_u")
	if errData == nil {
		t.Errorf("Failed to detect wrong dataset file")
	}

	// Fail by loading a dataset that does not have a receptor suffix
	err = os.WriteFile("pdb2.list",
		[]byte(
			"root_r_u.pdb\n"+
				"root_something.pdb\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write dataset file: %s", err)
	}

	defer os.Remove("pdb2.list")
	_, errData = LoadDataset(projectDir, "pdb2.list", "_r_u", "_l_u")
	if errData == nil {
		t.Errorf("Failed to detect wrong dataset file")
	}

}

func TestValidateTarget(t *testing.T) {
	var err error
	err = os.WriteFile("receptor.pdb", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("receptor.pdb")

	err = os.WriteFile("ligand.pdb", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("ligand.pdb")

	err = os.WriteFile("receptor.pqr", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("receptor.pqr")

	err = os.WriteFile("ligand.pqr", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("ligand.pqr")

	// Pass by finding both a receptor and ligand in a Target
	target := Target{
		ID:       "1",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{"ligand.pdb"},
	}

	err = target.Validate()
	if err != nil {
		t.Errorf("Failed to validate target: %s", err)
	}

	// Fail by not finding a receptor in a Target
	target = Target{
		ID:       "1",
		Receptor: []string{"does_not_exist.pdb"},
		Ligand:   []string{"ligand.pdb"},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail by not finding a ligand in a Target
	target = Target{
		ID:       "1",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{"does_not_exist.pdb"},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand is not a PDB file
	target = Target{
		ID:       "1",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{"ligand.pqr"},
	}
	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor is not a PDB file
	target = Target{
		ID:       "1",
		Receptor: []string{"receptor.pqr"},
		Ligand:   []string{"ligand.pdb"},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor field is empty
	target = Target{
		ID:       "1",
		Receptor: []string{""},
		Ligand:   []string{"ligand.pdb"},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand field is empty
	target = Target{
		ID:       "1",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{""},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ID field is empty
	target = Target{
		ID:       "",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{"ligand.pdb"},
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

}

func TestCreateDatasetDir(t *testing.T) {

	// Create a temporary directory
	// testDir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll("test")

	err := CreateDatasetDir("test")

	if err != nil {
		t.Errorf("Failed to create directory: %s", err)
	}

	// Fail by creating a directory that already exists
	err = CreateDatasetDir("test")
	if err == nil {
		t.Errorf("Failed to detect existing directory")
	}

}

func TestOrganizeDataset(t *testing.T) {
	var err error

	fileArr := []string{
		"ref.pdb",
		"receptor.pdb",
		"ligand.pdb",
		"receptor1.pdb",
		"receptor2.pdb",
		"ligand1.pdb",
		"ligand2.pdb",
		"ambig.tbl",
		"unambig.tbl",
		"receptor_list.txt",
		"ligand_list.txt",
		"custom.top",
		"custom.param"}

	for _, file := range fileArr {
		err := os.WriteFile(file, []byte(""), 0644)
		if err != nil {
			t.Errorf("Failed to write structure file: %s", err)
		}
		defer os.Remove(file)
	}

	testBmPath := "test_bm"
	err = os.Mkdir(testBmPath, 0755)
	if err != nil {
		t.Errorf("Failed to create benchmark directory: %s", err)
	}
	defer os.RemoveAll(testBmPath)

	target := Target{
		ID:           "1",
		Receptor:     []string{"receptor1.pdb", "receptor2.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"ligand1.pdb", "ligand2.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"ambig.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
		MiscPDB:      []string{"ref.pdb"},
	}

	// Pass by organizing data in a valid directory
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err != nil {
		t.Errorf("Failed to organize data: %s", err)
	}

	// Check if the files were copied to the correct directory
	path := "test_bm/1/data"
	expectedFiles := []string{
		"receptor1.pdb", "receptor2.pdb", "ligand1.pdb",
		"ligand2.pdb", "ambig.tbl", "unambig.tbl",
		"custom.top",
		"custom.param"}

	for _, v := range expectedFiles {
		f := path + "/" + v
		_, err = os.Stat(f)
		if err != nil {
			t.Errorf("Failed to copy file: %s", err)
		}
	}

	// Fail by organizing inexisting data - receptor
	target = Target{
		ID:           "1",
		Receptor:     []string{"does-not-exist.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"ligand1.pdb", "ligand2.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"ambig.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
		MiscPDB:      []string{"ref.pdb"},
	}
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by organizing inexisting data - ligand
	target = Target{
		ID:           "1",
		Receptor:     []string{"receptor1.pdb", "receptor2.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"does-not-exist.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"ambig.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
		MiscPDB:      []string{"ref.pdb"},
	}
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by trying to copy unexisting restraints
	target = Target{
		ID:           "1",
		Receptor:     []string{"receptor1.pdb", "receptor2.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"ligand1.pdb", "ligand2.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"does-not-exist.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
		MiscPDB:      []string{"ref.pdb"},
	}
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by trying to copy unexisting toppar
	target = Target{
		ID:           "1",
		Receptor:     []string{"receptor1.pdb", "receptor2.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"ligand1.pdb", "ligand2.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"ambig.tbl", "unambig.tbl"},
		Toppar:       []string{"does-not-exist.top"},
		MiscPDB:      []string{"ref.pdb"},
	}
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by trying to copy unexisting miscpdbs
	target = Target{
		ID:           "1",
		Receptor:     []string{"receptor1.pdb", "receptor2.pdb"},
		ReceptorList: "receptor_list.txt",
		Ligand:       []string{"ligand1.pdb", "ligand2.pdb"},
		LigandList:   "ligand_list.txt",
		Restraints:   []string{"ambig.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
		MiscPDB:      []string{"does-not-exist.pdb"},
	}
	_, err = OrganizeDataset(testBmPath, []Target{target})
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

}

func TestSetupHaddock3Scenario(t *testing.T) {

	inp := input.Input{}
	inp.Scenarios = []input.Scenario{
		{
			Name: "scenario1",
			Parameters: input.ScenarioParams{
				Modules: input.ModuleParams{
					Order: []string{"topoaa", "rigidbody"},
					Topoaa: map[string]interface{}{
						"some-param": "some-value",
					},
					Rigidbody: map[string]interface{}{
						"ambig_fname": "_ti.tbl",
					},
				},
			},
		},
	}
	s := inp.Scenarios[0]

	// Create a dummy PDB
	err := os.WriteFile("dummy.pdb", []byte("ATOM      1  N   ALA A   1      10.000  10.000  10.000  1.00  0.00           N\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.RemoveAll("dummy.pdb")
	// Make a list of PDB files and save it to a file
	err = os.WriteFile("pdb-files.txt", []byte("\"dummy.pdb\"\n\"dummy.pdb\"\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.RemoveAll("pdb-files.txt")

	target := Target{
		ID:           "1abc",
		Receptor:     []string{"dummy.pdb", "dummy.pdb"},
		ReceptorList: "pdb-files.txt",
		Ligand:       []string{"dummy.pdb", "dummy.pdb"},
		LigandList:   "pdb-files.txt",
		Restraints:   []string{"1abc_ti.tbl", "unambig.tbl"},
		Toppar:       []string{"custom.top", "custom.param"},
	}

	wd := "some-workdir"
	// hdir := "haddock-dir"
	defer os.RemoveAll(wd)

	j, err := target.SetupHaddock3Scenario(wd, s)
	if err != nil {
		t.Errorf("Failed to setup scenario: %s", err)
	}

	if j.ID != target.ID+"_"+s.Name {
		t.Errorf("Wrong scenario name: %s", j.ID)
	}

	// check if the scenario was written to disk
	scenarioPath := filepath.Join(wd, target.ID, "scenario-"+s.Name)
	if _, err := os.Stat(scenarioPath); os.IsNotExist(err) {
		t.Errorf("Scenario was not written to disk")
	}

}

func TestWriteRunToml(t *testing.T) {

	// Create a temporary directory
	_ = os.MkdirAll("_some-workdir", 0755)
	defer os.RemoveAll("_some-workdir")

	target := Target{
		ID:         "1abc",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pdb"},
		Restraints: []string{"ambig.tbl", "unambig.tbl", "something.tbl"},
		Toppar:     []string{"custom1.top", "custom2.param"},
		MiscPDB:    []string{"ref.pdb"},
	}

	m := input.ModuleParams{
		Order: []string{"topoaa", "rigidbody", "flexref", "mdref"},
		Topoaa: map[string]interface{}{
			"some-param": "some-value",
		},
		Rigidbody: map[string]interface{}{
			"some-other-param":   10,
			"some_fname":         "ambig",
			"another_fname":      "unambig",
			"other_fname":        "custom1",
			"someother_fname":    "custom2",
			"thereference_fname": "ref",
		},
		Flexref: map[string]interface{}{
			"some-other-param": 3.5,
		},
		Mdref: map[string]interface{}{
			"some-other-param": false,
		},
	}

	g := make(map[string]interface{})
	g["general-param1"] = "general-value"
	g["general-param2"] = 2.5
	g["general-param3"] = false
	g["general-param4"] = 1

	_, err := target.WriteRunToml("_some-workdir", g, m)

	if err != nil {
		t.Errorf("Failed to write run.toml: %s", err)
	}

	// check if the run.toml was written to disk
	runTomlPath := filepath.Join("_some-workdir", "run.toml")
	if _, err := os.Stat(runTomlPath); os.IsNotExist(err) {
		t.Errorf("run.toml was not written to disk")
	}

	// Fail by trying to write to a directory that does not exist
	_, err = target.WriteRunToml("_some-workdir/does_not_exist", g, m)
	if err == nil {
		t.Errorf("Failed to detect wrong input")
	}

}

func TestTarget_SetupHaddock24Scenario(t *testing.T) {
	type fields struct {
		ID           string
		Receptor     []string
		ReceptorList string
		Ligand       []string
		LigandList   string
		Restraints   []string
		Toppar       []string
		MiscPDB      []string
	}
	type args struct {
		wd   string
		hdir string
		s    input.Scenario
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    runner.Job
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "pass",
			fields: fields{
				ID:         "1abc",
				Receptor:   []string{"receptor.pdb"},
				Ligand:     []string{"ligand.pdb"},
				Restraints: []string{"ambig_ti.tbl", "other_unambig.tbl", "something.tbl"},
				Toppar:     []string{"custom1.top", "custom2.param"},
			},
			args: args{
				wd:   "some-workdir",
				hdir: "haddock-dir",
				s: input.Scenario{
					Name: "some-scenario",
					Parameters: input.ScenarioParams{
						Restraints: input.Restraints{
							Ambig:   "_ti",
							Unambig: "_unambig",
						},
					},
				},
			},
			want: runner.Job{
				ID:   "1abc_some-scenario",
				Path: "some-workdir/1abc/scenario-some-scenario",
				Restraints: input.Restraints{
					Ambig:   "ambig_ti.tbl",
					Unambig: "other_unambig.tbl",
				},
				Toppar: input.Toppar{
					Topology: "custom1.top",
					Param:    "custom2.param",
				},
			},
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				ID: "1abc",
			},
			args: args{
				wd:   "some-workdir",
				hdir: "",
				s: input.Scenario{
					Name: "some-scenario",
				},
			},
			want:    runner.Job{},
			wantErr: true,
		},
	}

	defer os.RemoveAll("some-workdir")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Target{
				ID:           tt.fields.ID,
				Receptor:     tt.fields.Receptor,
				ReceptorList: tt.fields.ReceptorList,
				Ligand:       tt.fields.Ligand,
				LigandList:   tt.fields.LigandList,
				Restraints:   tt.fields.Restraints,
				Toppar:       tt.fields.Toppar,
				MiscPDB:      tt.fields.MiscPDB,
			}
			got, err := tr.SetupHaddock24Scenario(tt.args.wd, tt.args.hdir, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Target.SetupHaddock24Scenario() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Target.SetupHaddock24Scenario() =\ngot   %v,\n want %v", got, tt.want)
			}
		})
	}
}

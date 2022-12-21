package dataset

import (
	"benchmarktools/input"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteRunParam(t *testing.T) {

	_ = os.Mkdir("1abc", 0755)
	defer os.RemoveAll("1abc")

	cwd, _ := os.Getwd()

	pDir := filepath.Join(cwd, "1abc")

	// Pass by writing a valid run parameter file
	target := Target{
		ID:         "1abc",
		Receptor:   []string{"1abc_r_1.pdb", "1abc_r_2.pdb"},
		Ligand:     []string{"1abc_l_1.pdb", "1abc_l_2.pdb"},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "/home/haddock2.4",
		ProjectDir: pDir,
	}

	_, err := target.WriteRunParam()
	if err != nil {
		t.Error(err)
	}

	// Fail by writing a run parameter file with an empty target
	target = Target{}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty receptor
	target = Target{
		ID:         "1abc",
		Receptor:   []string{},
		Ligand:     []string{"1abc.pdb"},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "/home/haddock2.4",
		ProjectDir: pDir,
	}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty ligand
	target = Target{
		ID:         "1abc",
		Receptor:   []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:     []string{},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "/home/haddock2.4",
		ProjectDir: pDir,
	}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty haddock directory
	target = Target{
		ID:         "1abc",
		Receptor:   []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:     []string{"1abc.pdb"},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "",
		ProjectDir: pDir,
	}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by writing a run parameter file with an empty project directory
	target = Target{
		ID:         "1abc",
		Receptor:   []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:     []string{"1abc.pdb"},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "/home/haddock2.4",
		ProjectDir: "",
	}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by not being able to write to the runparam file
	target = Target{
		ID:         "1abc",
		Receptor:   []string{"1abcA.pdb", "1abcB.pdb"},
		Ligand:     []string{"1abc.pdb"},
		Ambig:      "ambig.tbl",
		Unambig:    "unambig.tbl",
		HaddockDir: "/home/haddock2.4",
		ProjectDir: "/home/haddock2.4/projects/1abc",
	}

	_, err = target.WriteRunParam()
	if err == nil {
		t.Error("Expected error, got nil")
	}

}

func TestLoadDataset(t *testing.T) {

	var err error

	err = os.WriteFile("structure1_r_u.pdb", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("structure1_r_u.pdb")
	err = os.WriteFile("structure1_l_u.pdb", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write structure file: %s", err)
	}
	defer os.Remove("structure1_l_u.pdb")

	// Write a valid list
	err = os.WriteFile("pdb.list",
		[]byte(
			"some/path/structure1_r_u.pdb\n"+
				"some/path/structure1_l_u.pdb\n"+
				"some/path/structure2_l_u.pdb\n"+
				"some/path/structure2_r_u.pdb\n"+
				"some/path/structure3_r_u_0.pdb\n"+
				"some/path/structure3_r_u_1.pdb\n"+
				"some/path/structure3_l_u.pdb\n"), 0644)
	defer os.Remove("pdb.list")

	// Pass by loading a valid dataset

	tArr, errData := LoadDataset("pdb.list", "_r_u", "_l_u")
	if errData != nil {
		t.Errorf("Failed to load dataset: %s", err)
	}

	if len(tArr) != 3 {
		t.Errorf("Failed to load dataset: %d", len(tArr))
	}

	if len(tArr[2].Receptor) != 2 {
		t.Errorf("Failed to load dataset")
	}

	// Fail by loading a dataset with a wrong file
	_, errData = LoadDataset("wrong_file.pdb", "_r_u", "_l_u")
	if errData == nil {
		t.Errorf("Failed to detect wrong dataset file")
	}

	// Fail by loading a dataset with files that do not exist
	err = os.WriteFile("pdb2.list",
		[]byte(
			"not-found.pdb\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write dataset file: %s", err)
	}
	defer os.Remove("pdb2.list")
	_, errData = LoadDataset("pdb2.list", "_r_u", "_l_u")
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
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err != nil {
		t.Errorf("Failed to validate target: %s", err)
	}

	// Fail by not finding a receptor in a Target
	target = Target{
		ID:         "1",
		Receptor:   []string{"does_not_exist.pdb"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail by not finding a ligand in a Target
	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"does_not_exist.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand is not a PDB file
	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pqr"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}
	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor is not a PDB file
	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pqr"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor field is empty
	target = Target{
		ID:         "1",
		Receptor:   []string{""},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand field is empty
	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{""},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ID field is empty
	target = Target{
		ID:         "",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the HaddockDir field is empty

	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "",
		ProjectDir: "project",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ProjectDir field is empty
	target = Target{
		ID:         "1",
		Receptor:   []string{"receptor.pdb"},
		Ligand:     []string{"ligand.pdb"},
		HaddockDir: "haddock",
		ProjectDir: "",
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

	fileArr := []string{"receptor.pdb", "ligand.pdb", "ambig.tbl", "unambig.tbl"}

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

	arr := []Target{

		{
			ID:         "1",
			Receptor:   []string{"receptor.pdb"},
			Ligand:     []string{"ligand.pdb"},
			Ambig:      "ambig.tbl",
			Unambig:    "unambig.tbl",
			HaddockDir: "haddock",
			ProjectDir: "project",
		},
		{
			ID:         "2",
			Receptor:   []string{"receptor.pdb"},
			Ligand:     []string{"ligand.pdb"},
			HaddockDir: "haddock",
			ProjectDir: "project",
		},
	}

	// Pass by organizing data in a valid directory
	_, err = OrganizeDataset(testBmPath, arr)
	if err != nil {
		t.Errorf("Failed to organize data: %s", err)
	}

	// Fail by organizing inexisting data - receptor
	arr[0].Receptor = []string{"does_not_exist.pdb"}
	_, err = OrganizeDataset(testBmPath, arr)
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by organizing inexisting data - ligand
	arr[0].Receptor = []string{"receptor.pdb"}
	arr[0].Ligand = []string{"does_not_exist.pdb"}
	_, err = OrganizeDataset(testBmPath, arr)
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by organizing inexisting data - ambig
	arr[0].Ligand = []string{"ligand.pdb"}
	arr[0].Ambig = "does_not_exist.tbl"
	_, err = OrganizeDataset(testBmPath, arr)
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

	// Fail by organizing inexisting data - unambig
	arr[0].Ambig = "ambig.tbl"
	arr[0].Unambig = "does_not_exist.tbl"
	_, err = OrganizeDataset(testBmPath, arr)
	if err == nil {
		t.Errorf("Failed to detect wrong directory")
	}

}

func TestSetupScenarios(t *testing.T) {

	// Write the different run.params files in different directories

	inp := input.Input{
		General: input.GeneralStruct{
			HaddockDir: "location-of-haddock",
			WorkDir:    "location-of-workdir",
		},
		Scenarios: []input.ScenarioStruct{
			{
				Name: "scenario1",
			},
		},
	}
	defer os.RemoveAll("location-of-workdir")

	target := Target{
		ID:       "1",
		Receptor: []string{"receptor.pdb"},
		Ligand:   []string{"ligand.pdb"},
	}

	// Write the different run.params files in different directories

	err := target.SetupScenarios(&inp)
	fmt.Println(target.ProjectDir)
	if err != nil {
		t.Errorf("Failed to setup scenarios: %s", err)
	}

	// Fail by trying to setup a scenario without defining the HaddockDir field
	inp.General.HaddockDir = ""
	err = target.SetupScenarios(&inp)
	if err == nil {
		t.Errorf("Failed to detect wrong input")
	}

}

package dataset

import (
	"os"
	"testing"
)

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
				"some/path/structure2_r_u.pdb\n"), 0644)
	defer os.Remove("pdb.list")

	// Pass by loading a valid dataset

	_, errData := LoadDataset("pdb.list", "_r_u", "_l_u")
	if errData != nil {
		t.Errorf("Failed to load dataset: %s", err)
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
		ID:       "1",
		Receptor: "receptor.pdb",
		Ligand:   "ligand.pdb",
	}

	err = target.Validate()
	if err != nil {
		t.Errorf("Failed to validate target: %s", err)
	}

	// Fail by not finding a receptor in a Target
	target = Target{
		ID:       "1",
		Receptor: "does_not_exist.pdb",
		Ligand:   "ligand.pdb",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail by not finding a ligand in a Target
	target = Target{
		ID:       "1",
		Receptor: "receptor.pdb",
		Ligand:   "does_not_exist.pdb",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand is not a PDB file
	target = Target{
		ID:       "1",
		Receptor: "receptor.pdb",
		Ligand:   "ligand.pqr",
	}
	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor is not a PDB file
	target = Target{
		ID:       "1",
		Receptor: "receptor.pqr",
		Ligand:   "ligand.pdb",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the receptor field is empty
	target = Target{
		ID:       "1",
		Receptor: "",
		Ligand:   "ligand.pdb",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ligand field is empty
	target = Target{
		ID:       "1",
		Receptor: "receptor.pdb",
		Ligand:   "",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}

	// Fail when the ID field is empty
	target = Target{
		ID:       "",
		Receptor: "receptor.pdb",
		Ligand:   "ligand.pdb",
	}

	err = target.Validate()
	if err == nil {
		t.Errorf("Failed to detect wrong target")
	}
}

package runner

import (
	"benchmarktools/input"
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestRunHaddock24(t *testing.T) {

	_ = os.MkdirAll("cmd-test/run1", 0755)
	defer os.RemoveAll("cmd-test")

	j := Job{
		ID:   "test",
		Path: "cmd-test",
		Params: map[string]interface{}{
			"cmrest": "false",
		},
	}

	cmd := "echo test"
	logF, err := j.RunHaddock24(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

	// fail by passing a non existing command
	cmdNon := "non_existing_command"
	logF, err = j.RunHaddock24(cmdNon)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

}

func TestSetupHaddock24(t *testing.T) {

	_ = os.MkdirAll("cmd-test/run1/structures/it0", 0755)
	_ = os.MkdirAll("cmd-test/run1/structures/it1/water", 0755)
	_ = os.MkdirAll("cmd-test/run1/data/distances", 0755)
	_ = os.MkdirAll("cmd-test/run1/toppar", 0755)
	defer os.RemoveAll("cmd-test")

	_ = os.WriteFile("cmd-test/run1/run.cns", []byte("{===>} param1=true;"), 0644)

	for _, f := range []string{"ambig.tbl", "unambig.tbl", "gdp.top", "gdp.param"} {
		_ = os.WriteFile(f, []byte(""), 0644)
		defer os.Remove(f)
	}

	j := Job{
		ID:   "test",
		Path: "cmd-test",
		Params: map[string]interface{}{
			"param1": false,
		},
		Restraints: input.Restraints{
			Ambig:   "ambig.tbl",
			Unambig: "unambig.tbl",
		},
		Toppar: input.Toppar{
			Top:   "gdp.top",
			Param: "gdp.param",
		},
	}

	cmd := "echo test"
	logF, err := j.SetupHaddock24(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

	// Check if run.cns was edited
	cnsF, err := os.Open("cmd-test/run1/run.cns")
	if err != nil {
		t.Errorf("Error reading run.cns: %v", err)
	}
	scanner := bufio.NewScanner(cnsF)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "{===>} param1=") {
			if line != "{===>} param1=false;" {
				t.Errorf("Error editing run.cns: %v", err)
			}
			break
		}
	}

	// Check if restraints were copied
	_, err = os.Stat("cmd-test/run1/data/distances/ambig.tbl")
	if err != nil {
		t.Errorf("Error copying restraints: %v", err)
	}
	_, err = os.Stat("cmd-test/run1/data/distances/unambig.tbl")
	if err != nil {
		t.Errorf("Error copying restraints: %v", err)
	}

	// Check if toppar was copied
	_, err = os.Stat("cmd-test/run1/toppar/ligand.top")
	if err != nil {
		t.Errorf("Error copying toppar: %v", err)
	}
	_, err = os.Stat("cmd-test/run1/toppar/ligand.param")
	if err != nil {
		t.Errorf("Error copying toppar: %v", err)
	}

	// fail by passing a non existing command
	cmdNon := "non_existing_command"
	logF, err = j.SetupHaddock24(cmdNon)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

	// Fail by not being able to edit run.cns
	os.Remove("cmd-test/run1/run.cns")
	_, err = j.SetupHaddock24(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	_, _ = os.Create("cmd-test/run1/run.cns")

	// Fail by not being able to copy restraints - ambig
	j.Restraints.Ambig = "non_existing_file"
	_, err = j.SetupHaddock24(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	j.Restraints.Ambig = "ambig.tbl"

	// Pass by not having ambig restraints
	j.Restraints.Ambig = ""
	_, err = j.SetupHaddock24(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}
	j.Restraints.Ambig = "ambig.tbl"

	// Fail by not being able to copy restraints - unambig
	j.Restraints.Unambig = "non_existing_file"
	_, err = j.SetupHaddock24(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	j.Restraints.Unambig = "unambig.tbl"

	// Pass by not having unambig restraints
	j.Restraints.Unambig = ""
	_, err = j.SetupHaddock24(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}
	j.Restraints.Unambig = "unambig.tbl"

	// Fail by not being able to copy toppar - top
	j.Toppar.Top = "non_existing_file"
	_, err = j.SetupHaddock24(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	j.Toppar.Top = "gdp.top"

	// Fail by not being able to copy toppar - param
	j.Toppar.Param = "non_existing_file"
	_, err = j.SetupHaddock24(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}

}

func TestRunHaddock3(t *testing.T) {

	// Create a directory
	_ = os.MkdirAll("_run-test", 0755)
	defer os.RemoveAll("_run-test")

	// Create a Job
	j := Job{
		ID:   "test",
		Path: "_run-test",
	}

	// define the cmd
	cmd := "echo test"

	// Pass by running
	logF, err := j.RunHaddock3(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}

	// Check if log file was created
	_, err = os.Stat(logF)
	if err != nil {
		t.Errorf("Error creating log file: %v", err)
	}

	// Fail by running a non existing command
	cmdNon := "non_existing_command"
	_, err = j.RunHaddock3(cmdNon)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}

}

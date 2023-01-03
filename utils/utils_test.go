package utils

import (
	"flag"
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {

	var err error

	err = os.WriteFile("some-file", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.Remove("some-file")
	// Pass by copying a file

	err = CopyFile("some-file", "copied-file")
	if err != nil {
		t.Errorf("Failed to copy file: %s", err)
	}
	defer os.Remove("copied-file")

	// Fail by copying a file that does not exist
	err = CopyFile("does-not-exist", "some-file-copy")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

	// Fail by copying a file to a directory that does not exist
	err = CopyFile("some-file", "does-not-exist/some-file-copy")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

}

func TestIsFlagPassed(t *testing.T) {

	// Pass by passing a flag
	os.Args = []string{"benchmarktools", "-option1"}
	var option1 bool
	flag.BoolVar(&option1, "option1", false, "")
	flag.Parse()

	if !IsFlagPassed("option1") {
		t.Errorf("Failed to detect flag")
	}

	// Pass by passing a flag that is not set
	if IsFlagPassed("option2") {
		t.Errorf("Failed to detect flag")
	}

}

func TestIsHaddock3(t *testing.T) {

	// Create a folder structure that is the same as haddock3's
	err := os.MkdirAll("_test_haddock3/src/haddock/modules", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_haddock3")
	_, err = os.Create("_test_haddock3/src/haddock/modules/defaults.yaml")
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}

	// Pass by finding the defaults.yaml file
	if !IsHaddock3("_test_haddock3") {
		t.Errorf("Failed to detect haddock3")
	}

	// Fail by not finding the defaults.yaml file
	if IsHaddock3("_test_haddock3/src") {
		t.Errorf("Failed to detect haddock3")
	}

}

func TestIsHaddock24(t *testing.T) {

	// Create a folder structure that is the same as haddock2.4's
	err := os.MkdirAll("_test_haddock24/protocols", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_haddock24")
	_, err = os.Create("_test_haddock24/protocols/run.cns-conf")
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}

	// Pass by finding the run.cns-conf file
	if !IsHaddock24("_test_haddock24") {
		t.Errorf("Failed to detect haddock2.4")
	}

	// Fail by not finding the run.cns-conf file
	if IsHaddock24("_test_haddock24/protocols") {
		t.Errorf("Failed to detect haddock2.4")
	}

}

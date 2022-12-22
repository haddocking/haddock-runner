package runner

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {

	// Pass running a command without arguments
	cmdNoArg := "ls"
	cwd, _ := os.Getwd()
	_, errNoArg := Run(cmdNoArg, cwd)
	if errNoArg != nil {
		t.Error(errNoArg)
	}

	// Pass by running a command with argument
	cmdArg := "ls -l"
	logF, errArg := Run(cmdArg, cwd)
	if errArg != nil {
		t.Error(errArg)
	}
	os.Remove(logF)

	// Fail by passing an empty string
	_, errEmpty := Run("", cwd)
	if errEmpty == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by passing a non-existing command
	cmdNonExisting := "non-existing-command"
	logF, errNonExisting := Run(cmdNonExisting, cwd)
	if errNonExisting == nil {
		t.Error("Expected error, got nil")
	}
	os.Remove(logF)

}

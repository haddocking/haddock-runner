package fs

import "testing"

func TestRun(t *testing.T) {

	// Pass running a command without arguments
	cmdNoArg := "ls"
	_, errNoArg := Run(cmdNoArg)
	if errNoArg != nil {
		t.Error(errNoArg)
	}

	// Pass by running a command with argument
	cmdArg := "ls -l"
	_, errArg := Run(cmdArg)
	if errArg != nil {
		t.Error(errArg)
	}

	// Fail by passing an empty string
	_, errEmpty := Run("")
	if errEmpty == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by passing a non-existing command
	cmdNonExisting := "non-existing-command"
	_, errNonExisting := Run(cmdNonExisting)
	if errNonExisting == nil {
		t.Error("Expected error, got nil")
	}

}

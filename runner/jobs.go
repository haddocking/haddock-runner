package runner

import (
	"errors"
	"path/filepath"

	"benchmarktools/wrapper/haddock2"
)

type Job struct {
	ID     string
	Path   string
	Params map[string]interface{}
}

func (j Job) SetupHaddock(cmd string) (string, error) {

	logF, err := Run(cmd, j.Path)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return logF, err
	}

	// Edit run.cns
	runCns := filepath.Join(j.Path, "run1", "run.cns")
	if err := haddock2.EditRunCns(runCns, j.Params); err != nil {
		err := errors.New("Error editing run.cns: " + err.Error())
		return logF, err
	}

	return logF, nil
}

func (j Job) RunHaddock(cmd string) (string, error) {

	// Run HADDOCK
	run1Path := filepath.Join(j.Path, "run1")
	logF, err := Run(cmd, run1Path)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return logF, err
	}

	return logF, nil

}

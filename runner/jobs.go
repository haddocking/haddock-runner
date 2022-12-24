// Package runner provides a set of functions to run commands
package runner

import (
	"errors"
	"os"
	"path/filepath"

	"benchmarktools/input"
	"benchmarktools/utils"
	"benchmarktools/wrapper/haddock2"
)

// Job is the HADDOCK job
type Job struct {
	ID         string
	Path       string
	Params     map[string]interface{}
	Restraints input.Restraints
	Toppar     input.Toppar
}

// SetupHaddock sets up the HADDOCK job
// - Setup the `run1“ directory by running the haddock executable
// - Edit the `run.cns` file
// - Copy the restraints
// - Copy the custom toppar
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

	// Copy restraints
	targetPaths := []string{
		filepath.Join(j.Path, "run1", "data", "distances"),
		filepath.Join(j.Path, "run1", "structures", "it0"),
		filepath.Join(j.Path, "run1", "structures", "it1"),
		filepath.Join(j.Path, "run1", "structures", "it1", "water"),
	}
	for _, target := range targetPaths {

		ambigF := filepath.Join(target, "ambig.tbl")
		if j.Restraints.Ambig != "" {
			// Copy ambig file
			if err := utils.CopyFile(j.Restraints.Ambig, ambigF); err != nil {
				err := errors.New("Error copying ambig.tbl: " + err.Error())
				return logF, err
			}
		} else {
			// Create empty ambig.tbl
			_, _ = os.Create(ambigF)

		}

		unambigF := filepath.Join(target, "unambig.tbl")
		if j.Restraints.Unambig != "" {
			// Copy unambig file
			if err := utils.CopyFile(j.Restraints.Unambig, unambigF); err != nil {
				err := errors.New("Error copying unambig.tbl: " + err.Error())
				return logF, err
			}
		} else {
			// Create empty unambig.tbl
			_, _ = os.Create(unambigF)

		}
	}

	// Copy toppar
	topparPath := filepath.Join(j.Path, "run1", "toppar")
	if j.Toppar.Top != "" {
		dest := filepath.Join(topparPath, "ligand.top")
		if err := utils.CopyFile(j.Toppar.Top, dest); err != nil {
			err := errors.New("Error copying custom topology: " + err.Error())
			return logF, err
		}
	}
	if j.Toppar.Param != "" {
		dest := filepath.Join(topparPath, "ligand.param")
		if err := utils.CopyFile(j.Toppar.Param, dest); err != nil {
			err := errors.New("Error copying custom param: " + err.Error())
			return logF, err
		}
	}

	return logF, nil
}

// RunHaddock runs the HADDOCK job in run1 directory
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

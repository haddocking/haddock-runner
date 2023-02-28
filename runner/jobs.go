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
	Restraints input.Airs
	Toppar     input.TopologyParams
}

// SetupHaddock24 sets up the HADDOCK job
// - Setup the `run1“ directory by running the haddock executable
// - Edit the `run.cns` file
// - Copy the restraints
// - Copy the custom toppar
func (j Job) SetupHaddock24(cmd string) error {

	// TODO: Refactor this as a separate function
	// Append the restraints to run.param
	m := map[string]string{
		j.Restraints.Ambig:     "AMBIG_TBL",
		j.Restraints.Unambig:   "UNAMBIG_TBL",
		j.Restraints.Hbonds:    "HBOND_FILE",
		j.Restraints.Dihedrals: "DIHED_FILE",
		j.Restraints.Tensor:    "TENSOR_FILE",
		j.Restraints.Cryoem:    "CRYO-EM_FILE",
		j.Restraints.Rdc1:      "RDC1_FILE",
		j.Restraints.Rdc2:      "RDC2_FILE",
		j.Restraints.Rdc3:      "RDC3_FILE",
		j.Restraints.Rdc4:      "RDC4_FILE",
		j.Restraints.Rdc5:      "RDC5_FILE",
		j.Restraints.Rdc6:      "RDC6_FILE",
		j.Restraints.Rdc7:      "RDC7_FILE",
		j.Restraints.Rdc8:      "RDC8_FILE",
		j.Restraints.Rdc9:      "RDC9_FILE",
		j.Restraints.Rdc10:     "RDC10_FILE",
		j.Restraints.Dani1:     "DANI1_FILE",
		j.Restraints.Dani2:     "DANI2_FILE",
		j.Restraints.Dani3:     "DANI3_FILE",
		j.Restraints.Dani4:     "DANI4_FILE",
		j.Restraints.Dani5:     "DANI5_FILE",
		j.Restraints.Pcs1:      "PCS1_FILE",
		j.Restraints.Pcs2:      "PCS2_FILE",
		j.Restraints.Pcs3:      "PCS3_FILE",
		j.Restraints.Pcs4:      "PCS4_FILE",
		j.Restraints.Pcs5:      "PCS5_FILE",
		j.Restraints.Pcs6:      "PCS6_FILE",
		j.Restraints.Pcs7:      "PCS7_FILE",
		j.Restraints.Pcs8:      "PCS8_FILE",
		j.Restraints.Pcs9:      "PCS9_FILE",
		j.Restraints.Pcs10:     "PCS10_FILE",
	}
	runParam := filepath.Join(j.Path, "run.param")
	f, err := os.OpenFile(runParam, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		err := errors.New("Error opening run.param: " + err.Error())
		return err
	}

	for k, v := range m {
		if k != "" {
			_, _ = f.WriteString(v + "=../data/" + filepath.Base(k) + "\n")
		}
	}

	f.Close()

	_, err = Run(cmd, j.Path)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return err
	}

	// Edit run.cns
	runCns := filepath.Join(j.Path, "run1", "run.cns")
	_ = haddock2.EditRunCns(runCns, j.Params)

	// Copy toppar
	topparPath := filepath.Join(j.Path, "run1", "toppar")
	if j.Toppar.Topology != "" {
		dest := filepath.Join(topparPath, "ligand.top")
		src := j.Toppar.Topology
		if err := utils.CopyFile(src, dest); err != nil {
			err := errors.New("Error copying custom topology: " + err.Error())
			return err
		}
	}
	if j.Toppar.Param != "" {
		dest := filepath.Join(topparPath, "ligand.param")
		src := j.Toppar.Param
		if err := utils.CopyFile(src, dest); err != nil {
			err := errors.New("Error copying custom param: " + err.Error())
			return err
		}
	}

	return nil
}

// RunHaddock24 runs the HADDOCK job in run1 directory
func (j Job) RunHaddock24(cmd string) (string, error) {

	// Run HADDOCK24
	run1Path := filepath.Join(j.Path, "run1")
	logF, err := Run(cmd, run1Path)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return logF, err
	}

	return logF, nil

}

// RunHaddock3 runs the HADDOCK3 job in run directory
func (j Job) RunHaddock3(cmd string) (string, error) {

	// Run HADDOCK3
	runWD := filepath.Join(j.Path)
	cmd = cmd + " run.toml"
	logF, err := Run(cmd, runWD)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return logF, err
	}

	return logF, nil

}

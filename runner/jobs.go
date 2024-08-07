// Package runner provides a set of functions to run commands
package runner

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"haddockrunner/constants"
	"haddockrunner/input"
	"haddockrunner/runner/status"
	"haddockrunner/utils"
	"haddockrunner/wrapper/haddock2"

	"github.com/golang/glog"
)

var sleepFunc = time.Sleep

// Job is the HADDOCK job
type Job struct {
	ID         string
	Path       string
	Params     map[string]interface{}
	Restraints input.Airs
	Toppar     input.TopologyParams
	Status     string
}

// SetupHaddock24 sets up the HADDOCK job
// - Setup the `run1` directory by running the haddock executable
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

	// Check if there is a `job.sh` file in the run directory
	jobF := filepath.Join(runWD, "job.sh")
	_, err := os.Stat(jobF)
	if err == nil {
		cmd = utils.Sbatch_cmd + " " + jobF
	} else {
		cmd = cmd + " run.toml"
	}
	logF, err := Run(cmd, runWD)
	if err != nil {
		err := errors.New("Error running HADDOCK: " + err.Error())
		return logF, err
	}

	return logF, nil

}

// WaitUntil waits for the job to be in a given status
func (j *Job) WaitUntil(s []string, timeoutcounter int) error {
	while := true
	c := 0
	for while {
		_ = j.UpdateStatus(utils.GetJobID, utils.CheckSlurmStatus, 3)
		// if err != nil {
		// 	err := errors.New("Error getting job status: " + err.Error())
		// 	return err
		// }
		if slices.Contains(s, j.Status) {
			glog.Info("Job " + j.ID + " is " + j.Status)
			while = false
		} else {
			sleepFunc(time.Duration(constants.WAIT_FOR_SLURM) * time.Second)
		}

		// Check if the job is running for too long
		c++
		if c > timeoutcounter {
			totalWait := strconv.Itoa(constants.WAIT_FOR_SLURM * timeoutcounter)
			glog.Warning("Job " + j.ID + " is took too long (" + totalWait + "s), cancelling it")
			while = false
			j.Status = status.FAILED
		}
	}

	return nil

}

func (j Job) Run(version int, cmd string) (string, error) {

	var logF string
	var err error

	switch version {
	case 2:
		errSetup2 := j.SetupHaddock24(cmd)
		if errSetup2 != nil {
			err := errors.New("Failed to setup HADDOCK: " + errSetup2.Error())
			return logF, err
		}
		logF, _ = j.RunHaddock24(cmd)
		// if err != nil {
		// 	err := errors.New("Failed to run HADDOCK: " + err.Error())
		// 	return logF, err
		// }
	case 3:
		logF, err = j.RunHaddock3(cmd)
		if err != nil {
			err := errors.New("Failed to run HADDOCK: " + err.Error())
			return logF, err
		}
	}

	return logF, nil

}

// PrepareJobFile prepares the job file, returns the path to the job file
func (j *Job) PrepareJobFile(executable string, slurm input.SlurmParams) error {

	var header string
	var body string

	// Create the JobFile
	header = utils.CreateJobHeader(
		slurm.Partition,
		slurm.Account,
		slurm.Mail_user,
		slurm.Time,
		slurm.Cpus_per_task,
		slurm.Nodes,
		slurm.Ntasks_per_node,
	)

	// Create the JobBody
	body = utils.CreateJobBody(executable, j.Path)

	// Create the JobFile
	jobFile := filepath.Join(j.Path, "job.sh")
	f, err := os.Create(jobFile)
	if err != nil {
		err := errors.New("Error creating job file: " + err.Error())
		return err
	}

	// Write the JobFile
	_, _ = f.WriteString(header + body)

	_ = f.Close()

	return nil

}

type GetJobIDFunc func(file string) (string, error)
type GetSlurmStatusFunc func(jobID string) (string, error)

// func (j *Job) UpdateStatus(getJobID GetJobIDFunc, getSlurmStatus GetSlurmStatusFunc, version int) error {
func (j *Job) UpdateStatus(getJobID GetJobIDFunc, getSlurmStatus GetSlurmStatusFunc, version int) error {

	var logF string
	var positiveKeys []string
	var negativeKeys []string

	if version == 2 {
		logF = filepath.Join(j.Path, "run1", "haddock.out")
		positiveKeys = []string{"Finishing HADDOCK on:"}
		negativeKeys = []string{"An error has occurred"}

	} else if version == 3 {
		logF = filepath.Join(j.Path, "run1", "log")
		negativeKeys = []string{"ERROR"}
		positiveKeys = []string{"This HADDOCK3 run took"}

	} else {
		err := errors.New("invalid HADDOCK version")
		return err
	}

	// Check if the log file exists
	_, err := os.Stat(logF)
	if os.IsNotExist(err) {
		j.Status = status.QUEUED
		return nil
	}

	// Check if the log file contains any of the negative keys
	for _, k := range negativeKeys {
		found, _ := utils.SearchInLog(logF, k)
		// if err != nil {
		// 	return err
		// }
		if found {
			j.Status = status.FAILED
			return nil
		}
	}

	// Check if the log file contains any of the positive keys
	for _, k := range positiveKeys {
		found, _ := utils.SearchInLog(logF, k)
		// if err != nil {
		// 	return err
		// }
		if found {
			j.Status = status.DONE
			return nil
		}
	}

	// Before saying that the job is incomplete, check if its running on slurm
	newestFile := utils.FindNewestLogFile(j.Path)

	found, _ := utils.SearchInLog(newestFile, "Submitted batch job")
	// if err != nil {
	// 	return err
	// }

	if found {

		// Get what is the SLURM status of the job
		slurmJobID, err := getJobID(newestFile)
		if err != nil {
			err := errors.New("Error getting job ID: " + err.Error())
			return err
		}

		jobStatus, err := getSlurmStatus(slurmJobID)
		if err != nil {
			return err
		}

		// VERY IMPORTANT: The status from SLURM are different than the
		//  internal states of the `haddock-runner`!!
		switch jobStatus {
		case "COMPLETED":
			j.Status = status.DONE
			return nil
		case "RUNNING":
			j.Status = status.QUEUED
			return nil
		case "PENDING":
			j.Status = status.SUBMITTED
			return nil
		case "CANCELLED", "FAILED", "TIMEOUT":
			j.Status = status.FAILED
			return nil

		}

	}

	j.Status = status.INCOMPLETE
	return nil

}

func (j Job) Clean() error {
	_ = os.RemoveAll(filepath.Join(j.Path, "run1"))
	// if err != nil {
	// 	err := errors.New("Error cleaning job: " + err.Error())
	// 	return err
	// }

	extensions := []string{"*.out", "*.err", "*.txt"}
	for _, e := range extensions {
		files, _ := filepath.Glob(filepath.Join(j.Path, e))
		for _, f := range files {
			_ = os.Remove(f)
		}
	}

	return nil

}

func (j Job) Post(haddockVersion int, executable string, slurm input.SlurmParams) error {

	if slurm != (input.SlurmParams{}) {
		err := j.PrepareJobFile(executable, slurm)
		if err != nil {
			glog.Error("Failed to prepare job file: " + err.Error())
			return err
		}
	}

	_, runErr := j.Run(haddockVersion, executable)
	if runErr != nil {
		glog.Error("Failed to run HADDOCK: " + runErr.Error())
		return runErr
	}

	return nil

}

// SortJobs sorts an array of Job structs based on their ID field.
// It first sorts by the part after the underscore in the ID, then by the part before.
// If an ID doesn't contain an underscore, it's treated as a single part for comparison.
func SortJobs(arr []Job) []Job {
	sort.Slice(arr, func(i, j int) bool {
		return compareJobs(arr[i].ID, arr[j].ID)
	})
	return arr
}

// compareJobs compares two job IDs for sorting.
// It splits each ID at the underscore and compares the parts.
// The comparison prioritizes the part after the underscore, then the part before.
// If either ID lacks an underscore, the whole strings are compared.
func compareJobs(a, b string) bool {
	partsA := strings.SplitN(a, "_", 2)
	partsB := strings.SplitN(b, "_", 2)

	// If either string doesn't have a second part, compare the whole strings
	if len(partsA) != 2 || len(partsB) != 2 {
		return a < b
	}

	// First, compare the parts after the underscore
	if partsA[1] != partsB[1] {
		return partsA[1] < partsB[1]
	}

	// If parts after underscore are equal, compare the parts before the underscore
	return partsA[0] < partsB[0]
}

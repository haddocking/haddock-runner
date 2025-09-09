// Package utils contains utility functions
package utils

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	Sbatch_cmd = "sbatch"
	Sacct_cmd  = "sacct"
	Sacct_args = []string{
		"--format=JobID,State",
		"-n",
		"-j",
	}
)

// CopyFile copies a file from src to dst
// If the file does not exist or cannot be created, an error is returned
func CopyFile(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		err := errors.New("file does not exist: " + src)
		return err
	}

	source, _ := os.Open(src)
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, _ = io.Copy(destination, source)

	return nil
}

// IsFlagPassed returns true if the flag was passed
func IsFlagPassed(name string) bool {
	var found bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// IsHaddock3 returns true if the path is a Haddock 3 project
func IsHaddock3(p string) bool {
	rootPath, _ := filepath.Abs(p)
	DefaultsLoc := filepath.Clean(filepath.Join(rootPath, "src/haddock/modules/defaults.yaml"))

	if _, err := os.Stat(DefaultsLoc); os.IsNotExist(err) {
		return false
	}

	return true
}

// IsHaddock24 returns true if the path is a Haddock 2.4 project
func IsHaddock24(p string) bool {
	rootPath, _ := filepath.Abs(p)
	runCnsLoc := filepath.Clean(filepath.Join(rootPath, "protocols/run.cns-conf"))

	if _, err := os.Stat(runCnsLoc); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateEnsemble creates an ensemble file from a list of PDB files
func CreateEnsemble(p string, out string) error {
	// Read the list file and check how many models there should be
	file, err := os.Open(p)
	if err != nil {
		return err
	}
	defer file.Close()

	s := bufio.NewScanner(file)

	s.Split(bufio.ScanLines)
	nbModels := 1
	ens := ""

	for s.Scan() {
		// modelF := s.Text()
		path := strings.Trim(s.Text(), "\"")

		modelF, err := os.Open(path)
		if err != nil {
			return err
		}

		// Keep only ATOM records
		modelScanner := bufio.NewScanner(modelF)
		modelScanner.Split(bufio.ScanLines)
		modelStr := ""
		for modelScanner.Scan() {
			line := modelScanner.Text()
			if strings.HasPrefix(line, "ATOM") || strings.HasPrefix(line, "HETATM") {
				modelStr += line + "\n"
			}
		}
		if modelStr == "" {
			err := errors.New("empty file: " + path)
			return err
		}

		header := formatModelHeader(nbModels) + "\n"
		footer := "ENDMDL\n"
		ens += header + modelStr + footer

		nbModels++
	}
	ens += "END\n"

	_ = os.WriteFile(out, []byte(ens), 0644)

	return nil
}

// formatModelHeader formats the header of a model, see https://www.wwpdb.org/documentation/file-format-content/format33/sect9.html#MODEL
func formatModelHeader(model_id int) string {
	return fmt.Sprintf("MODEL     %-4d", model_id)
}

// IsUnique returns true if the slice contains unique elements
func IsUnique(s []string) bool {
	seen := make(map[string]struct{}, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			return false
		}
		seen[v] = struct{}{}
	}
	return true
}

// CopyFileArrTo copy files from an array to a location
func CopyFileArrTo(files []string, dst string) error {
	for _, f := range files {
		_, file := filepath.Split(f)
		err := CopyFile(f, filepath.Join(dst, file))
		if err != nil {
			return err
		}
	}

	return nil
}

// IntSliceToStringSlice converts an int slice to a string slice
func IntSliceToStringSlice(intSlice []int) []string {
	var stringSlice []string
	for _, v := range intSlice {
		stringSlice = append(stringSlice, strconv.Itoa(v))
	}
	return stringSlice
}

// InterfaceSliceToStringSlice converts an interface slice to a string slice
func InterfaceSliceToStringSlice(slice []interface{}) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		// this can be multiple things, so we need to convert it to a string
		switch v := v.(type) {
		case string:
			s[i] = v
		case int:
			s[i] = strconv.Itoa(v)
		case float64:
			s[i] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			s[i] = strconv.FormatBool(v)
		case nil:
			s[i] = ""
		}
	}
	return s
}

// FloatSliceToStringSlice converts a float slice to a string slice
func FloatSliceToStringSlice(slice []float64) []string {
	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return s
}

// Helper function to check if "cg" is present in the string
func ContainsCG(s string) bool {
	// Make it lower case
	s = strings.ToLower(s)
	return regexp.MustCompile(`cg`).MatchString(s)
}

// FindFname checks the array of strings for a pattern, returns an error if multiple files match
func FindFname(arr []string, pattern *regexp.Regexp) (string, error) {
	var fname string
	for _, f := range arr {
		if pattern.MatchString(f) {
			if fname != "" {
				err := errors.New("multiple files match the pattern: `" + pattern.String() + "` please use a more specific pattern")
				return "", err
			}
			fname = f
		}
	}

	return fname, nil
}

func SearchInLog(filePath, searchString string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchString) {
			return true, nil
		}
	}

	return false, nil
}

func CreateJobHeader(partition, account, mail_user, runtime string, cpus_per_task, nodes, ntasks_per_node int) string {
	header := "#!/bin/bash\n"
	header += "#SBATCH --job-name=haddock\n"
	header += "#SBATCH --output=haddock-%j.out\n"
	header += "#SBATCH --error=haddock-%j.err\n"
	if partition != "" {
		header += "#SBATCH --partition=" + partition + "\n"
	}
	if account != "" {
		header += "#SBATCH --account=" + account + "\n"
	}
	if mail_user != "" {
		header += "#SBATCH --mail-user=" + mail_user + "\n"
		header += "#SBATCH --mail-type=ALL\n"
	}
	if cpus_per_task != 0 {
		header += "#SBATCH --cpus-per-task=" + strconv.Itoa(cpus_per_task) + "\n"
	}
	if nodes != 0 {
		header += "#SBATCH --nodes=" + strconv.Itoa(nodes) + "\n"
	}
	if ntasks_per_node != 0 {
		header += "#SBATCH --ntasks-per-node=" + strconv.Itoa(ntasks_per_node) + "\n"
	}
	if runtime != "" {
		header += "#SBATCH --time=" + runtime + "\n"
	}
	header += "\n"

	return header
}

func CreateJobBody(cmd, path string) string {
	body := "cd " + path + "\n"
	body += cmd + " run.toml\n"

	return body
}

func FindNewestLogFile(path string) string {
	files, _ := filepath.Glob(filepath.Join(path, "*.txt"))
	// if err != nil {
	// 	return ""
	// }

	// Find what is the newest file
	var newestFile string
	var newestTime int64

	for _, f := range files {
		fi, _ := os.Stat(f)
		// if err != nil {
		// 	return ""
		// }
		if fi.ModTime().Unix() > newestTime {
			newestTime = fi.ModTime().Unix()
			newestFile = f
		}
	}
	return newestFile
}

// GetJobID recieves a log file and returns the job ID
//
// The log file should contain: "Submitted batch job XXXXXXX"
func GetJobID(logF string) (string, error) {
	file, err := os.Open(logF)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Submitted batch job") {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				return "", errors.New("job ID not found in " + logF)
			}
			jobID := fields[3]
			return jobID, nil
		}
	}

	return "", errors.New("job ID not found in " + logF)
}

func CheckSlurmStatus(jobID string) (string, error) {
	// Add the job ID to the arguments
	Sacct_args = append(Sacct_args, jobID)

	cmd := exec.Command(Sacct_cmd, Sacct_args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(string(out))
	jobStatus := ""
	if len(fields) >= 2 {
		jobStatus = fields[1]
	} else {
		return "", errors.New("error: Could not find enough fields in the string")
	}

	return jobStatus, nil
}

func CreateRootRegex(rsuf, lsuf, ssuf string) *regexp.Regexp {
	suffixes := []string{}

	if rsuf != "" {
		suffixes = append(suffixes, rsuf)
	}
	if lsuf != "" {
		suffixes = append(suffixes, lsuf)
	}
	if ssuf != "" {
		suffixes = append(suffixes, ssuf)
	}

	if len(suffixes) == 0 {
		return nil // or handle this case as needed
	}

	pattern := `(.*)(?:` + strings.Join(suffixes, "|") + `)`
	return regexp.MustCompile(pattern)
}

// RemoveString removes a string from a slice of strings
func RemoveString(slice []string, s string) []string {
	result := []string{}
	for _, v := range slice {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}

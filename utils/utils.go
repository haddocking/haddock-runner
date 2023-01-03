// Package utils contains utility functions
package utils

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
			if strings.HasPrefix(line, "ATOM") {
				modelStr += line + "\n"
			}
		}
		if modelStr == "" {
			err := errors.New("empty file: " + path)
			return err
		}

		header := fmt.Sprintf("MODEL    %-5d\n", nbModels)
		footer := "ENDMDL\n"
		ens += header + modelStr + footer

		nbModels++
	}
	ens += "END\n"

	_ = os.WriteFile(out, []byte(ens), 0644)

	return nil

}

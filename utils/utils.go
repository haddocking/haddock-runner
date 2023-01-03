// Package utils contains utility functions
package utils

import (
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

// MapInterfaceToString converts a map[string]interface{} to a string representation
//
// This is useful when writing to HADDOCK3's `run.toml`
func MapInterfaceToString(m map[string]interface{}) string {

	s := ""

	for k, v := range m {
		switch v := v.(type) {
		case string:
			s += k + " = \"" + v + "\"\n"
		case int:
			s += k + " = " + strconv.Itoa(v) + "\n"
		case float64:
			s += k + " = " + strconv.FormatFloat(v, 'f', -1, 64) + "\n"
		case bool:
			s += k + " = " + strconv.FormatBool(v) + "\n"
		}
	}

	return s

}

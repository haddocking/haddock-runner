// Package utils contains utility functions
package utils

import (
	"errors"
	"flag"
	"io"
	"os"
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

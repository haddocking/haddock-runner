// Package runner provides a set of functions to run commands
package runner

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Run executes a command in a given path
func Run(command string, path string) (string, error) {
	var app string
	var args []string

	if command == "" {
		err := errors.New("no command passed")
		return "", err
	}

	w := strings.Fields(command)

	if len(w) > 1 {
		app = w[0]
		args = w[1:]
	} else {
		app = w[0]
	}

	cmd := exec.Command(app, args...)
	cmd.Dir = path

	timestamp := time.Unix(time.Now().Unix(), 0).Format("2006-01-02_15:04:05")

	out, err := cmd.CombinedOutput()

	logFname := path + "/log_" + timestamp + ".txt"
	_ = os.WriteFile(logFname, out, 0644)

	if err != nil {
		err := errors.New("Error running command: " + command + " check: " + logFname)
		return logFname, err
	}

	return logFname, nil

}

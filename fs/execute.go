// Package fs provides a set of functions to interact with the filesystem
package fs

import (
	"errors"
	"os/exec"
	"strings"
)

// Path: fs/execute/execute.go

func Run(cmd string) (string, error) {
	var app string
	var args []string

	if cmd == "" {
		err := errors.New("no command passed")
		return "", err
	}

	w := strings.Fields(cmd)

	if len(w) > 1 {
		app = w[0]
		args = w[1:]
	} else {
		app = w[0]
	}

	output, err := exec.Command(app, args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil

}

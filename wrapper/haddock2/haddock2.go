// Package haddock2 provides a set of functions to interact with HADDOCK2+
package haddock2

import (
	"benchmarktools/input"
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

func EditRunCns(f string, s input.ScenarioStruct) error {

	if f == "" {
		return errors.New("run.cns file not defined")
	}

	if s.Parameters == nil {
		return errors.New("scenario parameters not defined")
	}

	if s.Name == "" {
		return errors.New("scenario name not defined")
	}

	// Read the run.cns file
	runCnsString, err := os.Open(f)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(runCnsString)

	scanner.Split(bufio.ScanLines)

	subTag := " ! This line was edited by the benchmarking tool"

	var newLines string

	for scanner.Scan() {
		line := scanner.Text()
		for key, data := range s.Parameters {

			paramRegex := regexp.MustCompile(`(?m){===>}\s(` + key + `)=.*;`)
			match := paramRegex.MatchString(line)

			if match {
				var subs string
				switch v := data.(type) {
				case string:
					subs = "{===>} $1=\"" + v + "\";" + subTag
				case int:
					subs = "{===>} $1=" + fmt.Sprint(v) + ";" + subTag
				case float64:
					subs = "{===>} $1=" + fmt.Sprint(v) + ";" + subTag
				case bool:
					subs = "{===>} $1=" + fmt.Sprint(v) + ";" + subTag
				}

				line = paramRegex.ReplaceAllString(line, subs)
				line += ""
				break
			}
		}
		newLines += line + "\n"

	}

	_ = os.WriteFile(f, []byte(newLines), 0644)

	return nil
}

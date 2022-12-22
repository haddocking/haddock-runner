// Package haddock2 provides a set of functions to interact with HADDOCK2+
package haddock2

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

func EditRunCns(runCns string, params map[string]interface{}) error {

	if runCns == "" {
		return errors.New("run.cns file not defined")
	}

	if len(params) == 0 {
		return errors.New("scenario parameters not defined")
	}

	// Read the run.cns file
	runCnsString, err := os.Open(runCns)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(runCnsString)

	scanner.Split(bufio.ScanLines)

	subTag := " ! This line was edited by the benchmarking tool"

	var newLines string

	for scanner.Scan() {
		line := scanner.Text()
		for key, data := range params {

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

	_ = os.WriteFile(runCns, []byte(newLines), 0644)

	return nil
}

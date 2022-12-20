// Input file handling
package input

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Input is the input structure
type Input struct {
	General   GeneralStruct    `yaml:"general"`
	Scenarios []ScenarioStruct `yaml:"scenarios"`
}

// GeneralStruct is the general structure
type GeneralStruct struct {
	HaddockExecutable string `yaml:"haddock_executable"`
	ReceptorSuffix    string `yaml:"receptor_suffix"`
	LigandSuffix      string `yaml:"ligand_suffix"`
}

// ScenarioStruct is the scenario structure
type ScenarioStruct struct {
	Name       string                 `yaml:"name"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

// ValidateExecutable checks if the executable is defined in PATH
func (inp *Input) ValidateExecutable() error {
	if inp.General.HaddockExecutable == "" {
		err := errors.New("executable not defined")
		return err
	}

	_, err := exec.Command("which", inp.General.HaddockExecutable).Output()
	if err != nil {
		return err
	}

	return nil
}

// LoadInput loads the input file
func LoadInput(filename string) (*Input, error) {

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	inp := &Input{}
	err = yaml.Unmarshal(yamlFile, inp)
	if err != nil {
		return nil, err
	}

	return inp, nil
}

// FindHaddock24RunCns finds the run.cns-conf file based on the executable location
func FindHaddock24RunCns(exec string) (string, error) {

	execPath, _ := filepath.Abs(exec)
	rootPath := filepath.Dir(execPath)
	runCnsLoc := filepath.Clean(filepath.Join(rootPath, "../protocols", "run.cns-conf"))

	if _, err := os.Stat(runCnsLoc); os.IsNotExist(err) {
		return "", err
	}

	return runCnsLoc, nil

}

// LoadHaddock24Params loads the parameters from the run.cns-conf file
func LoadHaddock24Params(filename string) (map[string]interface{}, error) {

	runCnsFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// params := make(map[string]interface{})
	scanner := bufio.NewScanner(runCnsFile)

	scanner.Split(bufio.ScanLines)

	stringRegex := regexp.MustCompile(`^\{===>\}\s(?P<key>\w+)=\"(?P<value>\w*)";`)
	intRegex := regexp.MustCompile(`^\{===>\}\s(?P<key>\w+)=(?P<value>\d+);`)
	floatRegex := regexp.MustCompile(`^\{===>\}\s(?P<key>\w+)=(?P<value>\d+\.\d+);`)
	boolRegex := regexp.MustCompile(`^\{===>\}\s(?P<key>\w+)=(?P<value>true|false);`)

	m := make(map[string]interface{})
	for scanner.Scan() {
		var res [][]string

		line := scanner.Text()

		res = stringRegex.FindAllStringSubmatch(line, -1)
		for i := range res {
			m[res[i][1]] = res[i][2]
		}

		res = intRegex.FindAllStringSubmatch(line, -1)
		for i := range res {
			valueInt, _ := strconv.Atoi(res[i][2])
			m[res[i][1]] = valueInt
		}

		res = floatRegex.FindAllStringSubmatch(line, -1)
		for i := range res {
			valueFloat, _ := strconv.ParseFloat(res[i][2], 64)
			m[res[i][1]] = valueFloat
		}

		res = boolRegex.FindAllStringSubmatch(line, -1)
		for i := range res {
			valueBool, _ := strconv.ParseBool(res[i][2])
			m[res[i][1]] = valueBool
		}
	}

	return m, nil

}

// Package input handles input parameters
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
	General   GeneralStruct `yaml:"general"`
	Scenarios []Scenario    `yaml:"scenarios"`
}

// GeneralStruct is the general structure
type GeneralStruct struct {
	HaddockExecutable string `yaml:"executable"`
	HaddockDir        string `yaml:"haddock_dir"`
	ReceptorSuffix    string `yaml:"receptor_suffix"`
	LigandSuffix      string `yaml:"ligand_suffix"`
	InputList         string `yaml:"input_list"`
	WorkDir           string `yaml:"work_dir"`
	MaxConcurrent     int    `yaml:"max_concurrent"`
}

// ScenarioStruct is the scenario structure
type Scenario struct {
	Name       string         `yaml:"name"`
	Parameters ScenarioParams `yaml:"parameters"`
}

type ScenarioParams struct {
	CnsParams  map[string]interface{} `yaml:"run_cns"`
	Restraints Restraints             `yaml:"restraints"`
	Toppar     Toppar                 `yaml:"custom_toppar"`
}

type Restraints struct {
	Ambig   string
	Unambig string
}

type Toppar struct {
	Top   string
	Param string
}

// ValidateExecutable checks if the executable is defined in PATH
func (inp *Input) ValidateExecutable() error {
	if inp.General.HaddockExecutable == "" {
		err := errors.New("executable not defined")
		return err
	}

	cmd := exec.Command(inp.General.HaddockExecutable)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// ValidateScenarioParams checks if the parameters names are valid
func ValidateRunCNSParams(known map[string]interface{}, params map[string]interface{}) error {

	for key := range params {
		if known[key] == nil {
			err := errors.New("`" + key + "` not valid")
			return err
		}

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
func FindHaddock24RunCns(p string) (string, error) {

	rootPath, _ := filepath.Abs(p)
	runCnsLoc := filepath.Clean(filepath.Join(rootPath, "protocols", "run.cns-conf"))

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

	m["ambig"] = ""
	m["unambig"] = ""

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

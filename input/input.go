package input

import (
	"errors"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Input struct {
	General   GeneralStruct    `yaml:"general"`
	Scenarios []ScenarioStruct `yaml:"scenarios"`
}

type GeneralStruct struct {
	HaddockExecutable string `yaml:"haddock_executable"`
	ReceptorSuffix    string `yaml:"receptor_suffix"`
	LigandSuffix      string `yaml:"ligand_suffix"`
}

type ScenarioStruct struct {
	Name       string                 `yaml:"name"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

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

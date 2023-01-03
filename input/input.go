// Package input handles input parameters
package input

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// ---------------------------------------------------------------------

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

// Scenario is the scenario structure
type Scenario struct {
	Name       string         `yaml:"name"`
	Parameters ScenarioParams `yaml:"parameters"`
}

// ScenarioParams is the scenario parameters structure
type ScenarioParams struct {
	CnsParams  map[string]interface{} `yaml:"run_cns"`
	Restraints Restraints             `yaml:"restraints"`
	Toppar     Toppar                 `yaml:"custom_toppar"`
	Modules    ModuleParams           `yaml:"modules"`
	General    map[string]interface{} `yaml:"general"`
}

// Restraints is the restraints structure
type Restraints struct {
	Ambig   string
	Unambig string
}

// Toppar is the toppar structure
type Toppar struct {
	Top   string
	Param string
}

// ModuleParams is the module parameters structure
type ModuleParams struct {
	Order         []string               `yaml:"order"`
	Topoaa        map[string]interface{} `yaml:"topoaa"`
	Topocg        map[string]interface{} `yaml:"topocg"`
	Exit          map[string]interface{} `yaml:"exit"`
	Emref         map[string]interface{} `yaml:"emref"`
	Flexref       map[string]interface{} `yaml:"flexref"`
	Mdref         map[string]interface{} `yaml:"mdref"`
	Gdock         map[string]interface{} `yaml:"gdock"`
	Lightdock     map[string]interface{} `yaml:"lightdock"`
	Rigidbody     map[string]interface{} `yaml:"rigidbody"`
	Emscoring     map[string]interface{} `yaml:"emscoring"`
	Mdscoring     map[string]interface{} `yaml:"mdscoring"`
	Caprieval     map[string]interface{} `yaml:"caprieval"`
	Clustfcc      map[string]interface{} `yaml:"clustfcc"`
	Clustrmsd     map[string]interface{} `yaml:"clustrmsd"`
	Rmsdmatrix    map[string]interface{} `yaml:"rmsdmatrix"`
	Seletop       map[string]interface{} `yaml:"seletop"`
	Seletopclusts map[string]interface{} `yaml:"seletopclusts"`
}

// ---------------------------------------------------------------------

// ValidateExecutable checks if the executable script has the correct permissions
func (inp *Input) ValidateExecutable() error {
	if inp.General.HaddockExecutable == "" {
		err := errors.New("executable not defined")
		return err
	}

	info, err := os.Stat(inp.General.HaddockExecutable)
	if err != nil {
		return err
	}
	mode := info.Mode()
	if mode&0111 != 0 && mode&0011 != 0 {
		return nil
	}
	return errors.New("executable not executable")

}

// ValidateRunCNSParams checks if the parameters names are valid
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

// LoadHaddock3Params reads the defaults.yaml files recursively and returns a list of modules
//
//	It returns an array of `Module` structs
func LoadHaddock3Params(p string) (ModuleParams, error) {

	// Check if path exists
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return ModuleParams{}, err
	}

	m := ModuleParams{}
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {

			moduleName := filepath.Base(filepath.Dir(path))
			name := cases.Title(language.Und, cases.NoLower).String(moduleName)
			yamlFile, _ := os.ReadFile(path)

			data := make(map[string]interface{})
			errMarshal := yaml.Unmarshal(yamlFile, &data)
			if errMarshal != nil {
				return errMarshal
			}

			// Add the data to the correct module
			v := reflect.ValueOf(&m).Elem()
			if v.FieldByName(name).IsValid() {
				v.FieldByName(name).Set(reflect.ValueOf(data))
			}

		}
		return nil
	})
	if err != nil {
		return ModuleParams{}, err
	}

	// TODO: Load the mandatory/optional parameters

	return m, nil

}

// ValidateHaddock3Params checks if the parameters names are valid
func ValidateHaddock3Params(known ModuleParams, loaded ModuleParams) error {

	v := reflect.ValueOf(loaded)
	k := reflect.ValueOf(known)

	types := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Map {
			for key := range field.Interface().(map[string]interface{}) {
				if !k.Field(i).MapIndex(reflect.ValueOf(key)).IsValid() {
					err := errors.New("`" + key + "` not valid for " + types.Field(i).Name)
					return err
				}
			}
		}
	}

	return nil

}

// Package input handles input parameters
package input

import (
	"bufio"
	"errors"
	"haddockrunner/utils"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------------

// Input is the input structure
type Input struct {
	General   GeneralStruct
	Slurm     SlurmParams
	Scenarios []Scenario
}

// GeneralStruct is the general structure
type GeneralStruct struct {
	HaddockExecutable string `yaml:"executable"`
	HaddockDir        string `yaml:"haddock_dir"`
	ReceptorSuffix    string `yaml:"receptor_suffix"`
	LigandSuffix      string `yaml:"ligand_suffix"`
	ShapeSuffix       string `yaml:"shape_suffix"`
	InputList         string `yaml:"input_list"`
	WorkDir           string `yaml:"work_dir"`
	MaxConcurrent     int    `yaml:"max_concurrent"`
}

type SlurmParams struct {
	Partition       string `yaml:"partition"`
	Cpus_per_task   int    `yaml:"cpus_per_task"`
	Ntasks_per_node int    `yaml:"ntasks_per_node"`
	Nodes           int    `yaml:"nodes"`
	Time            string `yaml:"time"`
	Account         string `yaml:"account"`
	Mail_user       string `yaml:"mail_user"`
}

// Scenario is the scenario structure
type Scenario struct {
	Name       string `yaml:"name"`
	Parameters ParametersStruct
}

// ParametersStruct is the parameters structure
type ParametersStruct struct {
	General    map[string]interface{} `yaml:"general"`
	Restraints Airs                   `yaml:"restraints"`
	Toppar     TopologyParams         `yaml:"custom_toppar"`
	CnsParams  map[string]interface{} `yaml:"run_cns"`
	Modules    ModuleParams
}

// Airs is the restraint structure
type Airs struct {
	Ambig     string
	Unambig   string
	Dihedrals string
	Hbonds    string
	Tensor    string
	Cryoem    string
	Rdc1      string
	Rdc2      string
	Rdc3      string
	Rdc4      string
	Rdc5      string
	Rdc6      string
	Rdc7      string
	Rdc8      string
	Rdc9      string
	Rdc10     string
	Dani1     string
	Dani2     string
	Dani3     string
	Dani4     string
	Dani5     string
	Pcs1      string
	Pcs2      string
	Pcs3      string
	Pcs4      string
	Pcs5      string
	Pcs6      string
	Pcs7      string
	Pcs8      string
	Pcs9      string
	Pcs10     string
}

// TopologyParams is the topology parameters structure
type TopologyParams struct {
	Topology string
	Param    string
}

// ---------------------------------------------------------------------

// ValidateExecutable checks if the executable script has the correct permissions
func (inp *Input) ValidateExecutable() error {
	if inp.General.HaddockExecutable == "" {
		err := errors.New("executable not defined")
		return err
	}

	if !filepath.IsAbs(inp.General.HaddockExecutable) {
		err := errors.New("`" + inp.General.HaddockExecutable + "` is not an absolute path")
		return err
	}

	info, err := os.Stat(inp.General.HaddockExecutable)
	if err != nil {
		return err
	}
	mode := info.Mode()
	if mode&0111 != 0 {
		return nil
	}
	return errors.New("executable not executable")
}

// ValidatePatterns checks if there are duplicated patterns in the input struct
func (inp *Input) ValidatePatterns() error {
	// ReceptorSuffix and LigandSuffix
	if inp.General.ReceptorSuffix == "" {
		err := errors.New("receptor_suffix not defined in `general` section")
		return err
	} else if inp.General.ReceptorSuffix == inp.General.LigandSuffix {
		err := errors.New("receptor_suffix and ligand_suffix are the same at `general` section")
		return err
	}

	scenarioFnameArray := []string{}

	for _, scenario := range inp.Scenarios {
		// Ambig and Unambig
		if scenario.Parameters.Restraints.Ambig != "" || scenario.Parameters.Restraints.Unambig != "" {
			if scenario.Parameters.Restraints.Ambig == scenario.Parameters.Restraints.Unambig {
				err := errors.New("ambig and unambig patterns are the same at`" + scenario.Name + "`scenario")
				return err
			}
		}

		// Topology and Param
		if scenario.Parameters.Toppar.Topology != "" || scenario.Parameters.Toppar.Param != "" {
			if scenario.Parameters.Toppar.Topology == scenario.Parameters.Toppar.Param {
				err := errors.New("topology and param patterns are the same at `" + scenario.Name + "` scenario")
				return err
			}
		}

		v := reflect.ValueOf(scenario.Parameters.Modules)
		types := v.Type()

		for _, m := range scenario.Parameters.Modules.Order {
			fnameArr := []string{}
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				fieldName := types.Field(i).Name
				if m == strings.ToLower(fieldName) {
					if field.Kind() == reflect.Map {
						for key, value := range field.Interface().(map[string]interface{}) {
							if strings.Contains(key, "_fname") {
								if value != nil {
									fnameArr = append(fnameArr, value.(string))
								}
							}
						}
					}
				}
			}
			// Check if there are duplicated patterns
			if !utils.IsUnique(fnameArr) {
				err := errors.New("duplicated patterns in `" + m + "` modules' `fname` parameters")
				return err
			}

			// Check if there are patterns that match each other
			scenarioFnameArray = append(scenarioFnameArray, fnameArr...)
		}
		// Check if there are patterns that match each other
		for i := 0; i < len(scenarioFnameArray); i++ {
			for j := i + 1; j < len(scenarioFnameArray); j++ {

				patternA := regexp.MustCompile(scenarioFnameArray[i])
				patternB := regexp.MustCompile(scenarioFnameArray[j])

				if scenarioFnameArray[i] != scenarioFnameArray[j] && (patternA.MatchString(scenarioFnameArray[j]) || patternB.MatchString(scenarioFnameArray[i])) {
					err := errors.New("patterns `" + scenarioFnameArray[i] + "` and `" + scenarioFnameArray[j] + "` match each other, please rename them")
					return err
				}
			}
		}

	}

	return nil
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

// ValidateExecutionModes checks if the execution modes are valid
func (inp *Input) ValidateExecutionModes() error {
	if inp.Slurm != (SlurmParams{}) {
		// Check if the executable is HADDOCK3
		if utils.IsHaddock24(inp.General.HaddockDir) {
			err := errors.New("cannot use SLURM with HADDOCK2")
			return err
		} else if utils.IsHaddock3(inp.General.HaddockDir) {
			// We need to check if the Scenarios are using the correct execution modes
			for _, scenario := range inp.Scenarios {
				if scenario.Parameters.General["mode"] != "local" {
					err := errors.New("SLURM can only be used with `mode: local`")
					return err
				}
			}
		}
	}
	return nil
}

// checkUnmappedModules compares the YAML map to the ModuleParams fields
func checkUnmappedModules(inp *Input, yamlMap map[string]interface{}) []string {
	var unmappedModules []string

	scenarios, _ := yamlMap["scenarios"].([]interface{})
	// if !ok {
	// 	return unmappedModules
	// }

	for _, scenario := range scenarios {
		scenarioMap, _ := scenario.(map[string]interface{})
		// if !ok {
		// 	continue
		// }

		parameters, _ := scenarioMap["parameters"].(map[string]interface{})
		// if !ok {
		// 	continue
		// }

		modules, _ := parameters["modules"].(map[string]interface{})
		// if !ok {
		// 	continue
		// }

		v := reflect.ValueOf(inp.Scenarios[0].Parameters.Modules)
		t := v.Type()

		for key := range modules {
			if key == "order" {
				continue // Skip the "order" field as it's handled differently
			}

			found := false
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				yamlTag := field.Tag.Get("yaml")
				if yamlTag == key || (yamlTag == "" && strings.ToLower(field.Name) == key) {
					found = true
					break
				}
			}
			if !found {
				unmappedModules = append(unmappedModules, key)
			}
		}
	}

	return unmappedModules
}

// LoadInput loads the input file and checks for unmapped modules
func LoadInput(filename string) (*Input, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a map to get all fields from YAML
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &yamlMap)
	if err != nil {
		return nil, err
	}

	// Unmarshal into the Input struct
	inp := &Input{}
	_ = yaml.Unmarshal(yamlFile, inp)
	// No checking the errors here since the unmarshalling above will catch it
	// if _ != nil {
	// 	return nil, err
	// }

	if utils.IsHaddock3(inp.General.HaddockDir) {

		// Check for unmapped modules
		unmappedModules := checkUnmappedModules(inp, yamlMap)
		if len(unmappedModules) > 0 {

			_s := "Module"
			_a := "was"
			if len(unmappedModules) > 1 {
				_s = "Modules"
				_a = "were"
			}
			unknownModules := strings.Join(unmappedModules, ", ")
			return nil, errors.New(_s + " `" + unknownModules + "` " + _a + " not found in the HADDOCK installation or supported in this version.")

		}
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
	paramPath := filepath.Join(p, "src/")
	// Check if path exists
	if _, err := os.Stat(paramPath); os.IsNotExist(err) {
		err := errors.New("path `" + paramPath + "` does not exist, is the `haddock_dir` correct?")
		return ModuleParams{}, err
	}

	m := ModuleParams{}
	err := filepath.Walk(paramPath, func(path string, info os.FileInfo, err error) error {
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
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				fieldName := v.Type().Field(i).Name
				if strings.Contains(strings.ToLower(fieldName), strings.ToLower(name)) {
					if field.IsValid() {
						if field.CanSet() {
							field.Set(reflect.ValueOf(data))
						}
					}
				}
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

	expandableRe := regexp.MustCompile(`(.)_\d?`)

	types := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Map {
			for key := range field.Interface().(map[string]interface{}) {
				if !k.Field(i).MapIndex(reflect.ValueOf(key)).IsValid() {
					// Check if the key is an expandable parameter
					match := expandableRe.MatchString(key)
					if !match {
						err := errors.New("`" + key + "` not valid for " + types.Field(i).Name)
						return err
					}
				}
			}
		}
	}

	return nil
}

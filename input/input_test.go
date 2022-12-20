package input

import (
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadInput(t *testing.T) {
	var err error
	// Pass by being able to load the input file
	// Create a OK input file
	inp := Input{
		General: GeneralStruct{
			HaddockExecutable: "haddockX",
			ReceptorSuffix:    "r_u",
			LigandSuffix:      "l_u",
		},
		Scenarios: []ScenarioStruct{
			{
				Name: "scenario1",
				Parameters: map[string]interface{}{
					"param1": false,
					"param2": "string",
					"param3": 1,
					"param4": 1.5,
				},
			},
		},
	}

	yamlData, err := yaml.Marshal(&inp)
	if err != nil {
		t.Errorf("Failed to marshal input file: %s", err)
	}

	err = os.WriteFile("input.yml", yamlData, 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}
	defer os.Remove("input.yml")

	_, err = LoadInput("input.yml")
	if err != nil {
		t.Errorf("Failed to load input file: %s", err)
	}

	// Fail by not being able to load the input file
	err = os.WriteFile("wrong_input.yml", []byte("not a yml"), 0644)
	if err != nil {
		t.Errorf("Failed to write input file: %s", err)
	}
	defer os.Remove("wrong_input.yml")

	_, err = LoadInput("wrong_input.yml")
	if err == nil {
		t.Errorf("Failed to detect wrong input file")
	}

	// Fail by trying to load a file that does not exist
	_, err = LoadInput("does_not_exist.yml")
	if err == nil {
		t.Errorf("Failed to detect wrong input file")
	}

}

func TestValidateExecutable(t *testing.T) {

	// Pass by finding the executable in the PATH
	inp := Input{
		General: GeneralStruct{
			HaddockExecutable: "ls",
		},
	}

	err := inp.ValidateExecutable()
	if err != nil {
		t.Errorf("Failed to validate executable: %s", err)
	}

	// Fail by not finding the executable in the PATH
	inp = Input{
		General: GeneralStruct{
			HaddockExecutable: "does_not_exist",
		},
	}

	err = inp.ValidateExecutable()
	if err == nil {
		t.Errorf("Failed to detect wrong executable")
	}

	// Fail because no executable is defined
	inp = Input{
		General: GeneralStruct{
			HaddockExecutable: "",
		},
	}

	err = inp.ValidateExecutable()
	if err == nil {
		t.Errorf("Failed to detect empty executable")
	}

}

func TestFindHaddock24RunCns(t *testing.T) {

	// Based on the executable, return the location of run.cns

	// Create an executable and place it two levels above run.cns
	os.MkdirAll("_test/Haddock", 0755)
	defer os.RemoveAll("_test")

	err := os.WriteFile("_test/Haddock/RunHaddock.py", []byte(""), 0755)
	if err != nil {
		t.Errorf("Failed to write executable: %s", err)
	}

	os.Mkdir("_test/protocols", 0755)
	err = os.WriteFile("_test/protocols/run.cns-conf", []byte("{===>} parameter=\"value\";"), 0755)
	if err != nil {
		t.Errorf("Failed to write run.cns: %s", err)
	}

	// Pass by finding the run.cns file
	_, err = FindHaddock24RunCns("_test/Haddock/RunHaddock.py")
	if err != nil {
		t.Errorf("Failed to find run.cns: %s", err)
	}

	// Fail by not finding the run.cns file
	_, err = FindHaddock24RunCns("does_not_exist")
	if err == nil {
		t.Errorf("Failed to detect wrong executable")
	}

}

func TestLoadHaddock24Params(t *testing.T) {
	// Parse the run.cns file and return the parameters as ParameterStruct
	params := []byte(
		"{===>} parameter1=\"value\";\n" +
			"{===>} parameter2=1;\n" +
			"{===>} parameter3=1.0;\n" +
			"{===>} parameter4=true;\n")
	err := os.WriteFile("_test_run.cns-conf", params, 0755)
	if err != nil {
		t.Errorf("Failed to write run.cns: %s", err)
	}
	defer os.Remove("_test_run.cns-conf")

	// Pass by finding the parameters
	p, err := LoadHaddock24Params("_test_run.cns-conf")
	if err != nil {
		t.Errorf("Failed to load parameters: %s", err)
	}

	if p["parameter1"] != "value" {
		t.Errorf("Failed to parse parameter1")
	}

	if p["parameter2"] != 1 {
		t.Errorf("Failed to parse parameter2")
	}

	if p["parameter3"] != 1.0 {
		t.Errorf("Failed to parse parameter3")
	}

	if p["parameter4"] != true {
		t.Errorf("Failed to parse parameter4")
	}

	// Fail by not finding the parameters
	_, err = LoadHaddock24Params("does_not_exist")
	if err == nil {
		t.Errorf("Failed to detect wrong executable")
	}

}

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

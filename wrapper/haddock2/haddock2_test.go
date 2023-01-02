package haddock2

import (
	"os"
	"testing"
)

func TestEditRunCns(t *testing.T) {

	dummyRunCNs := "{+ Dummy run.cns for testing +}\n"
	dummyRunCNs += "{===>} parameter_1=\"value\";\n"
	dummyRunCNs += "{===>} parameter_2=true;\n"
	dummyRunCNs += "{===>} parameter_3=1;\n"
	dummyRunCNs += "{===>} parameter_4=1.5;\n"
	dummyRunCNs += "{===>} parameter_5=\"must-remain\";\n"

	_ = os.WriteFile("test_run.cns", []byte(dummyRunCNs), 0644)
	defer os.Remove("test_run.cns")

	// Pass by editing a valid run.cns file
	params := map[string]interface{}{
		"parameter_1": "new_value",
		"parameter_2": false,
		"parameter_3": 2,
		"parameter_4": 2.5,
	}

	// s := input.ScenarioStruct{
	// 	Name:       "test",
	// 	Parameters: params,
	// }

	err := EditRunCns("test_run.cns", params)
	if err != nil {
		t.Error(err)
	}

	// Fail by editing a run.cns file with an empty scenario
	// s = input.ScenarioStruct{}
	emptyParams := map[string]interface{}{}

	err = EditRunCns("test_run.cns", emptyParams)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by editing an empty run.cns file
	err = EditRunCns("", params)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Fail by trying to edit a run.cns that does not exist
	err = EditRunCns("does_not_exist.cns", params)

	if err == nil {
		t.Error("Expected error, got nil")
	}

}

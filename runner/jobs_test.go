package runner

import (
	"benchmarktools/input"
	"os"
	"testing"
)

func TestRunHaddock24(t *testing.T) {

	_ = os.MkdirAll("cmd-test/run1", 0755)
	defer os.RemoveAll("cmd-test")

	j := Job{
		ID:   "test",
		Path: "cmd-test",
		Params: map[string]interface{}{
			"cmrest": "false",
		},
	}

	cmd := "echo test"
	logF, err := j.RunHaddock24(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

	// fail by passing a non existing command
	cmdNon := "non_existing_command"
	logF, err = j.RunHaddock24(cmdNon)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}
	os.Remove(logF)

}

func TestRunHaddock3(t *testing.T) {

	// Create a directory
	_ = os.MkdirAll("_run-test", 0755)
	defer os.RemoveAll("_run-test")

	// Create a Job
	j := Job{
		ID:   "test",
		Path: "_run-test",
	}

	// define the cmd
	cmd := "echo test"

	// Pass by running
	logF, err := j.RunHaddock3(cmd)
	if err != nil {
		t.Errorf("Error running haddock: %v", err)
	}

	// Check if log file was created
	_, err = os.Stat(logF)
	if err != nil {
		t.Errorf("Error creating log file: %v", err)
	}

	// Fail by running a non existing command
	cmdNon := "non_existing_command"
	_, err = j.RunHaddock3(cmdNon)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}

}

func TestJob_SetupHaddock24(t *testing.T) {
	type fields struct {
		ID         string
		Path       string
		Params     map[string]interface{}
		Restraints input.Airs
		Toppar     input.TopologyParams
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				ID:   "test",
				Path: "cmd-test",
				Params: map[string]interface{}{
					"param1": true,
				},
				Restraints: input.Airs{
					Ambig: "ambig.tbl",
				},
				Toppar: input.TopologyParams{},
			},
			args:    args{"echo test"},
			wantErr: false,
		},
		{
			name: "missing run.param",
			fields: fields{
				Path: "does-not-exist",
			},
			args:    args{"echo test"},
			wantErr: true,
		},
		{
			name: "cannot run",
			fields: fields{
				Path: "cmd-test",
			},
			args:    args{"non_existing_command"},
			wantErr: true,
		},
		{
			name: "cannot copy topology",
			fields: fields{
				Path: "cmd-test",
				Toppar: input.TopologyParams{
					Topology: "non_existing_file",
				},
			},
			args:    args{"echo test"},
			wantErr: true,
		},
		{
			name: "cannot copy param",
			fields: fields{
				Path: "cmd-test",
				Toppar: input.TopologyParams{
					Param: "non_existing_file",
				},
			},
			args:    args{"echo test"},
			wantErr: true,
		},
	}
	_ = os.MkdirAll("cmd-test/run1/structures/it0", 0755)
	_ = os.MkdirAll("cmd-test/run1/structures/it1/water", 0755)
	_ = os.MkdirAll("cmd-test/run1/data/distances", 0755)
	_ = os.MkdirAll("cmd-test/run1/toppar", 0755)
	defer os.RemoveAll("cmd-test")

	_ = os.WriteFile("cmd-test/run1/run.cns", []byte("{===>} param1=true;"), 0644)

	for _, f := range []string{"cmd-test/run.param", "ambig.tbl", "unambig.tbl", "gdp.top", "gdp.param"} {
		_ = os.WriteFile(f, []byte(""), 0644)
		defer os.Remove(f)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Job{
				ID:         tt.fields.ID,
				Path:       tt.fields.Path,
				Params:     tt.fields.Params,
				Restraints: tt.fields.Restraints,
				Toppar:     tt.fields.Toppar,
			}
			if err := j.SetupHaddock24(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Job.SetupHaddock24() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

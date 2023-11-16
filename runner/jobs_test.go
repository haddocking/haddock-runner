package runner

import (
	"haddockrunner/input"
	"os"
	"path/filepath"
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

	// Create a test directory
	testD := "cmd-test"
	_ = os.MkdirAll(testD, 0755)
	defer os.RemoveAll(testD)

	setupHaddock24ForTest(testD)

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

func TestJobRun(t *testing.T) {

	temptestP, _ := os.MkdirTemp("", "test-job-run")
	setupHaddock24ForTest(temptestP)
	defer os.RemoveAll(temptestP)

	type fields struct {
		j Job
	}
	type args struct {
		version int
		cmd     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass v2",
			fields: fields{
				j: Job{
					Path: temptestP,
				},
			},
			args: args{
				version: 2,
				cmd:     "echo test",
			},
			wantErr: false,
		},
		{
			name: "pass v3",
			fields: fields{
				j: Job{},
			},
			args: args{
				version: 3,
				cmd:     "echo test",
			},
			wantErr: false,
		},
		{
			name: "fail v2",
			fields: fields{
				j: Job{
					Path: "non-existing-path",
				},
			},
			args: args{
				version: 2,
				cmd:     "echo test",
			},
			wantErr: true,
		},
		{
			name: "fail v3",
			fields: fields{
				j: Job{
					Path: "non-existing-path",
				},
			},
			args: args{
				version: 3,
				cmd:     "echo test",
			},
			wantErr: true,
		},
		{
			name: "fail v2 with non-existing command",
			fields: fields{
				j: Job{
					Path: temptestP,
				},
			},
			args: args{
				version: 2,
				cmd:     "non_existing_command",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := tt.fields.j.Run(tt.args.version, tt.args.cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("Job.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestJobGetStatus(t *testing.T) {

	// Setup the positive test for v2
	v2PositiveTempD, _ := os.MkdirTemp("", "v2-positive")
	defer os.RemoveAll(v2PositiveTempD)
	os.MkdirAll(filepath.Join(v2PositiveTempD, "run1"), 0755)
	logF := filepath.Join(v2PositiveTempD, "run1", "haddock.out")
	_ = os.WriteFile(logF, []byte("Finishing HADDOCK on:"), 0644)

	// Setup the negative test for v2
	v2NegativeTempD, _ := os.MkdirTemp("", "v2-negative")
	defer os.RemoveAll(v2NegativeTempD)
	os.MkdirAll(filepath.Join(v2NegativeTempD, "run1"), 0755)
	logF = filepath.Join(v2NegativeTempD, "run1", "haddock.out")
	_ = os.WriteFile(logF, []byte("An error has occurred"), 0644)

	// Setup the positive test for v3
	v3PositiveTempD, _ := os.MkdirTemp("", "v3-positive")
	defer os.RemoveAll(v3PositiveTempD)
	os.MkdirAll(filepath.Join(v3PositiveTempD, "run1"), 0755)
	logF = filepath.Join(v3PositiveTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte("This HADDOCK3 run took"), 0644)

	// Setup the incomplete scenario
	incompleteTempD, _ := os.MkdirTemp("", "incomplete")
	defer os.RemoveAll(incompleteTempD)
	os.MkdirAll(filepath.Join(incompleteTempD, "run1"), 0755)
	logF = filepath.Join(incompleteTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte(""), 0644)

	type fields struct {
		j Job
	}
	type args struct {
		version int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass v2 positive",
			fields: fields{
				j: Job{
					Path: v2PositiveTempD,
				},
			},
			args: args{
				version: 2,
			},
			wantErr: false,
		},
		{
			name: "pass v2 negative",
			fields: fields{
				j: Job{
					Path: v2NegativeTempD,
				},
			},
			args: args{
				version: 2,
			},
			wantErr: false,
		},
		{
			name: "pass v3 positive",
			fields: fields{
				j: Job{
					Path: v3PositiveTempD,
				},
			},
			args: args{
				version: 3,
			},
			wantErr: false,
		},
		{
			name: "fail with invalid version",
			fields: fields{
				j: Job{
					Path: v2PositiveTempD,
				},
			},
			args: args{
				version: 0,
			},
			wantErr: true,
		},
		{
			name: "pass with incomplete v3",
			fields: fields{
				j: Job{
					Path: incompleteTempD,
				},
			},
			args: args{
				version: 3,
			},
			wantErr: false,
		},
		{
			name: "pass without log file",
			fields: fields{
				j: Job{
					Path: "non-existing-path",
				},
			},
			args: args{
				version: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.j.GetStatus(tt.args.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("Job.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

// setupHaddock24ForTest is a helper function that creates a directory structure for testing
//
// Note: This could be a Mock- feel free to change it (:
func setupHaddock24ForTest(p string) error {
	_ = os.MkdirAll(filepath.Join(p, "run1/structures/it0"), 0755)
	_ = os.MkdirAll(filepath.Join(p, "run1/structures/it1/water"), 0755)
	_ = os.MkdirAll(filepath.Join(p, "run1/data/distances"), 0755)
	_ = os.MkdirAll(filepath.Join(p, "run1/toppar"), 0755)
	_ = os.WriteFile(filepath.Join(p, "run1/run.cns"), []byte("{===>} param1=true;"), 0644)
	inputF := []string{
		filepath.Join(p, "run.param"),
		filepath.Join(p, "ambig.tbl"),
		filepath.Join(p, "unambig.tbl"),
		filepath.Join(p, "gdp.top"),
		filepath.Join(p, "gdp.param"),
	}

	for _, f := range inputF {
		_ = os.WriteFile(f, []byte(""), 0644)
	}
	return nil
}

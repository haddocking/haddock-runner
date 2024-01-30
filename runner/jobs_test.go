package runner

import (
	"errors"
	"haddockrunner/input"
	"haddockrunner/runner/status"
	"os"
	"path/filepath"
	"testing"
	"time"
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

	// Create a directory that contains a job.sh file
	_ = os.MkdirAll("_run-test-with-job-file", 0755)
	jobF := filepath.Join("_run-test-with-job-file", "job.sh")
	_ = os.WriteFile(jobF, []byte(""), 0644)
	defer os.RemoveAll("_run-test-with-job-file")

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

	// Fail by running a command in a directory that contains a job.sh file
	j.Path = "_run-test-with-job-file"
	_, err = j.RunHaddock3(cmd)
	if err == nil {
		t.Errorf("Error running haddock: %v", err)
	}

}

func TestJob_SetupHaddock24(t *testing.T) {

	// Create a test directory
	testD := "cmd-test"
	err := os.MkdirAll(testD, 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}

	defer os.RemoveAll(testD)

	err = setupHaddock24ForTest(testD)
	if err != nil {
		t.Errorf("Error setting up test directory: %v", err)
	}

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
	defer os.RemoveAll(temptestP)

	err := setupHaddock24ForTest(temptestP)
	if err != nil {
		t.Errorf("Error setting up test directory: %v", err)
	}

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
				j: Job{
					Path: temptestP,
				},
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
	err := os.MkdirAll(filepath.Join(v2PositiveTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF := filepath.Join(v2PositiveTempD, "run1", "haddock.out")
	_ = os.WriteFile(logF, []byte("Finishing HADDOCK on:"), 0644)

	// Setup the negative test for v2
	v2NegativeTempD, _ := os.MkdirTemp("", "v2-negative")
	defer os.RemoveAll(v2NegativeTempD)
	err = os.MkdirAll(filepath.Join(v2NegativeTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF = filepath.Join(v2NegativeTempD, "run1", "haddock.out")
	_ = os.WriteFile(logF, []byte("An error has occurred"), 0644)

	// Setup the positive test for v3
	v3PositiveTempD, _ := os.MkdirTemp("", "v3-positive")
	defer os.RemoveAll(v3PositiveTempD)
	err = os.MkdirAll(filepath.Join(v3PositiveTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF = filepath.Join(v3PositiveTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte("This HADDOCK3 run took"), 0644)

	// Setup a submitted test for v3
	v3SubmittedTempD, _ := os.MkdirTemp("", "v3-submitted")
	defer os.RemoveAll(v3SubmittedTempD)
	err = os.MkdirAll(filepath.Join(v3SubmittedTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF = filepath.Join(v3SubmittedTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte(""), 0644)
	subLogF := filepath.Join(v3SubmittedTempD, "something.txt")
	_ = os.WriteFile(subLogF, []byte("Submitted batch job 42"), 0644)

	// Setup a path which is using slurm but has the incorrect log file
	slurmTempD, _ := os.MkdirTemp("", "slurm")
	defer os.RemoveAll(slurmTempD)
	err = os.MkdirAll(filepath.Join(slurmTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF = filepath.Join(slurmTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte(""), 0644)
	subLogWrongF := filepath.Join(slurmTempD, "something.txt")
	_ = os.WriteFile(subLogWrongF, []byte("Submitted batch job 42"), 0644)

	// Setup the incomplete scenario
	incompleteTempD, _ := os.MkdirTemp("", "incomplete")
	defer os.RemoveAll(incompleteTempD)
	err = os.MkdirAll(filepath.Join(incompleteTempD, "run1"), 0755)
	if err != nil {
		t.Errorf("Error creating test directory: %v", err)
	}
	logF = filepath.Join(incompleteTempD, "run1", "log")
	_ = os.WriteFile(logF, []byte(""), 0644)

	mockGetJobID := func(logF string) (string, error) {
		return "42", nil
	}

	mockGetJobIDFail := func(logF string) (string, error) {
		return "", errors.New("")
	}

	mockCheckSlurmStatus := func(jobID string) (string, error) {
		return "RUNNING", nil
	}

	mockCheckSlurmStatusFail := func(jobID string) (string, error) {
		return "", errors.New("")
	}

	mockCheckSlurmStatusDone := func(jobID string) (string, error) {
		return "COMPLETED", nil
	}

	mockCheckSlurmStatusQueue := func(jobID string) (string, error) {
		return "RUNNING", nil
	}

	mockCheckSlurmStatusSubmitted := func(jobID string) (string, error) {
		return "PENDING", nil
	}

	mockCheckSlurmStatusFailed := func(jobID string) (string, error) {
		return "FAILED", nil
	}

	type fields struct {
		j Job
	}
	type args struct {
		getJobID       GetJobIDFunc
		getSlurmStatus GetSlurmStatusFunc
		version        int
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        2,
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        2,
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        3,
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        0,
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        3,
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
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        3,
			},
			wantErr: false,
		},
		{
			name: "pass with submitted v3",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        3,
			},
			wantErr: false,
		},
		{
			name: "fail with mockGetJobIDFail",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobIDFail,
				getSlurmStatus: mockCheckSlurmStatus,
				version:        3,
			},
			wantErr: true,
		},
		{
			name: "fail with mockCheckSlurmStatusFail",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatusFail,
				version:        3,
			},
			wantErr: true,
		},
		{
			name: "pass with mockCheckSlurmStatusDone",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatusDone,
				version:        3,
			},
			wantErr: false,
		},
		{
			name: "pass with mockCheckSlurmStatusQueue",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatusQueue,
				version:        3,
			},
			wantErr: false,
		},
		{
			name: "pass with mockCheckSlurmStatusSubmitted",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatusSubmitted,
				version:        3,
			},
			wantErr: false,
		},
		{
			name: "pass with mockCheckSlurmStatusFailed",
			fields: fields{
				j: Job{
					Path: v3SubmittedTempD,
				},
			},
			args: args{
				getJobID:       mockGetJobID,
				getSlurmStatus: mockCheckSlurmStatusFailed,
				version:        3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.j.UpdateStatus(tt.args.getJobID, tt.args.getSlurmStatus, tt.args.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("Job.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestJobPrepareJobFile(t *testing.T) {

	// Create a valid Path
	err := os.MkdirAll("_test_prepare_job_file", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_prepare_job_file")

	type fields struct {
		j Job
	}
	type args struct {
		executable string
		slurm      input.SlurmParams
	}

	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "pass by creating a job file",
			args: args{
				executable: "echo test",
				slurm:      input.SlurmParams{},
			},
			fields: fields{
				j: Job{
					Path: "_test_prepare_job_file",
				},
			},
		},
		{
			name: "fail by passing a non-existing path",
			args: args{
				executable: "echo test",
				slurm:      input.SlurmParams{},
			},
			fields: fields{
				j: Job{
					Path: "does-not-exist",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.j.PrepareJobFile(tt.args.executable, tt.args.slurm)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareJobFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestJobClean(t *testing.T) {

	// Create a valid Path
	cleanPath := "_test_clean"
	err := os.MkdirAll(cleanPath, 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll(cleanPath)

	runPath := filepath.Join(cleanPath, "run1")
	err = os.MkdirAll(runPath, 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}

	// Fill it with some files
	txtF := filepath.Join(cleanPath, "test.txt")
	_ = os.WriteFile(txtF, []byte(""), 0644)
	errF := filepath.Join(cleanPath, "test.err")
	_ = os.WriteFile(errF, []byte(""), 0644)
	outF := filepath.Join(cleanPath, "test.out")
	_ = os.WriteFile(outF, []byte(""), 0644)

	type fields struct {
		j Job
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass by cleaning",
			fields: fields{
				j: Job{
					Path: cleanPath,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.j.Clean()
			if (err != nil) != tt.wantErr {
				t.Errorf("Clean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestWaitUntil(t *testing.T) {

	sleepCallCount := 0
	originalSleepFunc := sleepFunc
	sleepFunc = func(d time.Duration) {
		sleepCallCount++
	}
	defer func() { sleepFunc = originalSleepFunc }()

	type fields struct {
		j Job
	}

	type args struct {
		s              []string
		timeoutcounter int
	}

	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name: "pass by waiting",
			fields: fields{
				j: Job{
					Status: status.QUEUED,
				},
			},
			args: args{
				s:              []string{status.QUEUED},
				timeoutcounter: 1,
			},
			wantErr: false,
		},
		{
			name: "test with expected sleep call",
			fields: fields{
				j: Job{
					Status: status.QUEUED,
				},
			},
			args: args{
				s:              []string{status.UNKNOWN},
				timeoutcounter: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.j.WaitUntil(tt.args.s, tt.args.timeoutcounter)
			if (err != nil) != tt.wantErr {
				t.Errorf("WaitUntil() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}

}

func TestJobPost(t *testing.T) {

	// Create a valid Path
	err := os.MkdirAll("_test_post", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_post")

	_ = setupHaddock24ForTest("_test_post")

	type field struct {
		j Job
	}

	type args struct {
		haddockVersion int
		executable     string
		slurm          input.SlurmParams
	}

	tests := []struct {
		name    string
		fields  field
		args    args
		wantErr bool
	}{
		{
			name: "pass by posting",
			fields: field{
				j: Job{
					Path: "_test_post",
				},
			},
			args: args{
				haddockVersion: 2,
				executable:     "echo test",
				slurm: input.SlurmParams{
					Partition: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "fail trying to run in a folder that does not exist",
			fields: field{
				j: Job{
					Path: "does-not-exist",
				},
			},
			args: args{
				haddockVersion: 3,
				executable:     "echo test",
				slurm: input.SlurmParams{
					Partition: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "fail trying to setup a slurm job in a folder that does not exist",
			fields: field{
				j: Job{
					Path: "does-not-exist",
				},
			},
			args: args{
				haddockVersion: 3,
				executable:     "echo test",
				slurm:          input.SlurmParams{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.fields.j.Post(tt.args.haddockVersion, tt.args.executable, tt.args.slurm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
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

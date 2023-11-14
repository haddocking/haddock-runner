package checksum

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {

	// Generate a dummyfile to test
	tempF, err := os.CreateTemp("", "test_generate")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	defer os.Remove(tempF.Name())

	type fields struct {
		Input     string
		InputList string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				filePath: tempF.Name(),
			},
			want:    "d41d8cd98f00b204e9800998ecf8427e",
			wantErr: false,
		},
		{
			name: "fail with empty file",
			args: args{
				filePath: "empty_file",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checksum.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Checksum.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRead tests the Read function with various scenarios
func TestRead(t *testing.T) {

	// Write a file for the test
	tmpfile, err := os.CreateTemp("", "test_read")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test case content to the file
	_, err = tmpfile.WriteString("first line\nsecond line")
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Write a file with only one line
	tmpfile2, err := os.CreateTemp("", "test_read")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	defer os.Remove(tmpfile2.Name())

	type args struct {
		filePath string
	}

	testCases := []struct {
		name    string
		args    args
		want    Checksum
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				filePath: tmpfile.Name(),
			},
			want: Checksum{
				Input:     "first line",
				InputList: "second line",
			},
			wantErr: false,
		},
		{
			name: "fail non-existing file",
			args: args{
				filePath: "non-existing-file",
			},
			want:    Checksum{},
			wantErr: true,
		},
		{
			name: "test with empty file",
			args: args{
				filePath: tmpfile2.Name(),
			},
			want:    Checksum{},
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			// Run the Read function
			result, err := Read(tt.args.filePath)

			// Check the result
			if (err != nil) != tt.wantErr {
				t.Error("Expected an error, but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check the result content
			if result != tt.want {
				t.Errorf("Expected result %v, got %v", tt.want, result)
			}
		})
	}
}

func TestAreEqual(t *testing.T) {
	type args struct {
		a Checksum
		b Checksum
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "pass",
			args: args{
				a: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
				b: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
			},
			want: true,
		},
		{
			name: "fail",
			args: args{
				a: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
				b: Checksum{
					Input:     "first line",
					InputList: "third line",
				},
			},
			want: false,
		},
		{
			name: "fail with empty",
			args: args{
				a: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
				b: Checksum{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Run the AreEqual function
			if got := AreEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("AreEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tempF, err := os.CreateTemp("", "test_write")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	defer os.Remove(tempF.Name())

	type args struct {
		filePath string
		c        Checksum
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				filePath: tempF.Name(),
				c: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
			},
			wantErr: false,
		},
		{
			name: "fail writing to a non-existing directory",
			args: args{
				filePath: "non-existing-dir/test_write",
				c: Checksum{
					Input:     "first line",
					InputList: "second line",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Run the Write function
			if err := Write(tt.args.filePath, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestValidateChecksum(t *testing.T) {

	// Create a tempdir
	tempDir, err := os.MkdirTemp("", "test_validate_checksum")
	if err != nil {
		t.Errorf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dummy file
	tempF, err := os.CreateTemp(tempDir, "test_validate_checksum")
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempF.Name())

	// Write test case content to the file
	_, err = tempF.WriteString("first line\nsecond line")
	if err != nil {
		t.Fatal(err)
	}

	// Write a valid checksum file
	checksumFile := tempDir + "/checksum.txt"
	c := Checksum{
		Input:     "311f4a819681297456bb96840218f676",
		InputList: "311f4a819681297456bb96840218f676",
	}
	errWriteChecksum := Write(checksumFile, c)
	if errWriteChecksum != nil {
		t.Errorf("Error writing checksum file: %v", errWriteChecksum)
	}

	// Write a new checksum file
	newChecksumF := tempDir + "/new_checksum.txt"
	defer os.Remove(newChecksumF)

	// Write a different checksum file
	differentChecksumF := tempDir + "/different_checksum.txt"
	c2 := Checksum{
		Input:     "different",
		InputList: "different",
	}
	errWriteDifferentChecksum := Write(differentChecksumF, c2)
	if errWriteDifferentChecksum != nil {
		t.Errorf("Error writing checksum file: %v", errWriteDifferentChecksum)
	}
	defer os.Remove(differentChecksumF)

	type args struct {
		inputF    string
		inputL    string
		checksumF string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				inputF:    tempF.Name(),
				inputL:    tempF.Name(),
				checksumF: checksumFile,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "fail without passing inputF",
			args: args{
				inputF:    "",
				inputL:    tempF.Name(),
				checksumF: checksumFile,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail without passing inputL",
			args: args{
				inputF:    tempF.Name(),
				inputL:    "",
				checksumF: checksumFile,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "fail trying to write checksum to a non-existing directory",
			args: args{
				inputF:    tempF.Name(),
				inputL:    tempF.Name(),
				checksumF: "non-existing-dir/checksum.txt",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "pass writing a new checksum",
			args: args{
				inputF:    tempF.Name(),
				inputL:    tempF.Name(),
				checksumF: newChecksumF,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "fail with a different checksum",
			args: args{
				inputF:    tempF.Name(),
				inputL:    tempF.Name(),
				checksumF: differentChecksumF,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Run the ValidateChecksum function
			got, err := ValidateChecksum(tt.args.inputF, tt.args.inputL, tt.args.checksumF)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateChecksum() = %v, want %v", got, tt.want)
			}
		})
	}

}

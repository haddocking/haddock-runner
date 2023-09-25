package utils

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestCopyFile(t *testing.T) {

	var err error

	err = os.WriteFile("some-file", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.Remove("some-file")
	// Pass by copying a file

	err = CopyFile("some-file", "copied-file")
	if err != nil {
		t.Errorf("Failed to copy file: %s", err)
	}
	defer os.Remove("copied-file")

	// Fail by copying a file that does not exist
	err = CopyFile("does-not-exist", "some-file-copy")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

	// Fail by copying a file to a directory that does not exist
	err = CopyFile("some-file", "does-not-exist/some-file-copy")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

}

func TestIsFlagPassed(t *testing.T) {

	// Pass by passing a flag
	os.Args = []string{"haddockrunner", "-option1"}
	var option1 bool
	flag.BoolVar(&option1, "option1", false, "")
	flag.Parse()

	if !IsFlagPassed("option1") {
		t.Errorf("Failed to detect flag")
	}

	// Pass by passing a flag that is not set
	if IsFlagPassed("option2") {
		t.Errorf("Failed to detect flag")
	}

}

func TestIsHaddock3(t *testing.T) {

	// Create a folder structure that is the same as haddock3's
	err := os.MkdirAll("_test_haddock3/src/haddock/modules", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_haddock3")
	_, err = os.Create("_test_haddock3/src/haddock/modules/defaults.yaml")
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}

	// Pass by finding the defaults.yaml file
	if !IsHaddock3("_test_haddock3") {
		t.Errorf("Failed to detect haddock3")
	}

	// Fail by not finding the defaults.yaml file
	if IsHaddock3("_test_haddock3/src") {
		t.Errorf("Failed to detect haddock3")
	}

}

func TestIsHaddock24(t *testing.T) {

	// Create a folder structure that is the same as haddock2.4's
	err := os.MkdirAll("_test_haddock24/protocols", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("_test_haddock24")
	_, err = os.Create("_test_haddock24/protocols/run.cns-conf")
	if err != nil {
		t.Errorf("Failed to create file: %s", err)
	}

	// Pass by finding the run.cns-conf file
	if !IsHaddock24("_test_haddock24") {
		t.Errorf("Failed to detect haddock2.4")
	}

	// Fail by not finding the run.cns-conf file
	if IsHaddock24("_test_haddock24/protocols") {
		t.Errorf("Failed to detect haddock2.4")
	}

}

func TestCreateEnsemble(t *testing.T) {

	// Write a dummy PDB file
	err := os.WriteFile("dummy.pdb", []byte("ATOM      1  N   ALA A   1      10.000  10.000  10.000  1.00  0.00           N\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.Remove("dummy.pdb")

	// Make a list of PDB files and save it to a file
	err = os.WriteFile("pdb-files.txt", []byte("\"dummy.pdb\"\n\"dummy.pdb\"\n\"dummy.pdb\"\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.RemoveAll("pdb-files.txt")

	// Create an ensemble
	outF := "ensemble.pdb"
	err = CreateEnsemble("pdb-files.txt", "ensemble.pdb")
	if err != nil {
		t.Errorf("Failed to create ensemble: %s", err)
	}

	// Check if the ensemble file exists
	if _, err := os.Stat(outF); os.IsNotExist(err) {
		t.Errorf("Failed to create ensemble file")
	}
	defer os.Remove(outF)

	// Fail by passing a file that does not exist
	err = CreateEnsemble("does-not-exist.txt", "ensemble.pdb")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

	// Fail by passing a file that does not point to a PDB file
	err = os.WriteFile("pdb-files.txt", []byte("\"i-dont-exist.pdb\"\n\"dummy.pdb\"\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	err = CreateEnsemble("pdb-files.txt", "ensemble.pdb")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}

	// Fail by passing a file with empty pdb files
	err = os.WriteFile("dummy.pdb", []byte(""), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	err = os.WriteFile("pdb-files.txt", []byte("\"dummy.pdb\"\n\"dummy.pdb\"\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	err = CreateEnsemble("pdb-files.txt", "ensemble.pdb")
	if err == nil {
		t.Errorf("Failed to detect wrong file")
	}
}

func TestIsUnique(t *testing.T) {

	// Pass by passing a list of unique elements
	if !IsUnique([]string{"a", "b", "c"}) {
		t.Errorf("Failed to detect unique elements")
	}

	// Fail by passing a list of non-unique elements
	if IsUnique([]string{"a", "b", "a"}) {
		t.Errorf("Failed to detect non-unique elements")
	}

}

func TestCopyFileArrTo(t *testing.T) {

	var err error

	arr := []string{"dummy.pdb", "dummy.pdb", "dummy.pdb"}
	err = os.WriteFile(arr[0], []byte("ATOM      1  N   ALA A   1      10.000  10.000  10.000  1.00  0.00           N\n"), 0644)
	if err != nil {
		t.Errorf("Failed to write file: %s", err)
	}
	defer os.Remove("dummy.pdb")

	err = os.MkdirAll("dummy", 0755)
	if err != nil {
		t.Errorf("Failed to create folder: %s", err)
	}
	defer os.RemoveAll("dummy")

	// Pass by copying the files
	err = CopyFileArrTo(arr, "dummy")
	if err != nil {
		t.Errorf("Failed to copy files: %s", err)
	}

	// Fail by not being able to copy the files
	err = CopyFileArrTo(arr, "does-not-exist")
	if err == nil {
		t.Errorf("Failed to detect wrong folder")
	}

}

func TestIntSliceToStringSlice(t *testing.T) {
	type args struct {
		intSlice []int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestIntSliceToStringSlice",
			args: args{
				intSlice: []int{1, 2, 3},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntSliceToStringSlice(tt.args.intSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntSliceToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceSliceToStringSlice(t *testing.T) {
	type args struct {
		slice []interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestInterfaceSliceToStringSlice",
			args: args{
				slice: []interface{}{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestInterfaceSliceToStringSlice-mutliple-types",
			args: args{
				slice: []interface{}{"1", 2, 3.5, true, nil},
			},
			want: []string{"1", "2", "3.5", "true", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceSliceToStringSlice(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceSliceToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatSliceToStringSlice(t *testing.T) {
	type args struct {
		slice []float64
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestFloatSliceToStringSlice",
			args: args{
				slice: []float64{1.0, 2.0, 3.0},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "TestFloatSliceToStringSlice",
			args: args{
				slice: []float64{1.6, 2.1, 3.3},
			},
			want: []string{"1.6", "2.1", "3.3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FloatSliceToStringSlice(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FloatSliceToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsCG(t *testing.T) {
	result := ContainsCG("cg")
	if !result {
		t.Errorf("Failed to detect cg")
	}
	result = ContainsCG("c")
	if result {
		t.Errorf("Failed to detect cg")
	}
}

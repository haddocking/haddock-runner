// Package dataset handles dataset parameters and files
package dataset

import (
	"haddockrunner/input"
	"haddockrunner/runner"
	"haddockrunner/utils"
	"io"
	"reflect"
	"strconv"
	"strings"

	// "haddockrunner/wrapper/haddock2"
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/golang/glog"
)

// Target is the target structure
type Target struct {
	ID           string
	Receptor     []string
	ReceptorList string
	Ligand       []string
	LigandList   string
	Restraints   []string
	Toppar       []string
	MiscPDB      []string
	Shape        []string
}

// Validate validates the Target checking if
//   - Fields are not empty
//   - Files exist
//   - Files are PDB files
func (t *Target) Validate() error {
	if t.ID == "" {
		return errors.New("Target ID not defined")
	}

	for _, r := range t.Receptor {
		if r == "" {
			return errors.New("Target receptor not defined")
		}
		if _, err := os.Stat(r); err != nil {
			return errors.New("Target receptor file not found" + r)
		}
		if filepath.Ext(r) != ".pdb" {
			return errors.New("Target receptor file not a PDB file" + r)
		}
	}

	for _, l := range t.Ligand {
		if l == "" {
			return errors.New("Target ligand not defined")
		}
		if _, err := os.Stat(l); err != nil {
			return errors.New("Target ligand file not found" + l)
		}
		if filepath.Ext(l) != ".pdb" {
			return errors.New("Target ligand file not a PDB file" + l)
		}
	}

	return nil
}

// SetupHaddock24Scenario method prepares the scenario
//   - Creates the scenario directory
//   - Creates the run.params file
func (t *Target) SetupHaddock24Scenario(wd string, hdir string, s input.Scenario) (runner.Job, error) {
	sPath := filepath.Join(wd, t.ID, "scenario-"+s.Name)
	glog.Info("Preparing : " + s.Name)
	_ = os.MkdirAll(sPath, 0755)

	// Generate the run.params file
	_, err := t.WriteRunParamStub(sPath, hdir)
	if err != nil {
		return runner.Job{}, err
	}

	// Find which restraints need to be used
	// FIXME: there's probably a better way to do this
	restraints := input.Airs{}

	v := reflect.ValueOf(&s.Parameters.Restraints).Elem()
	k := reflect.ValueOf(&restraints).Elem()

	types := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).String()
		name := types.Field(i).Name
		for _, r := range t.Restraints {
			if field != "" {
				if strings.Contains(r, field) {
					k.FieldByName(name).SetString(r)
				}
			}
		}
	}

	toppar := input.TopologyParams{}
	for _, t := range t.Toppar {
		if filepath.Ext(t) == ".top" {
			toppar.Topology = t
		}
		if filepath.Ext(t) == ".param" {
			toppar.Param = t
		}

	}

	j := runner.Job{
		ID:         t.ID + "_" + s.Name,
		Path:       sPath,
		Params:     s.Parameters.CnsParams,
		Restraints: restraints,
		Toppar:     toppar,
	}

	return j, nil
}

// WriteRunParamStub writes the run.param file
func (t *Target) WriteRunParamStub(projectDir string, haddockDir string) (string, error) {
	var runParamString string
	nMol := 1

	if haddockDir == "" {
		err := errors.New("haddock directory not defined")
		return "", err
	}

	if projectDir == "" {
		err := errors.New("project directory not defined")
		return "", err
	}

	if len(t.Receptor) == 0 {
		err := errors.New("receptor not defined")
		return "", err
	}

	runParamString += "RUN_NUMBER=1\n"
	runParamString += "PROJECT_DIR=./\n"
	runParamString += "HADDOCK_DIR=" + haddockDir + "\n"

	// Write receptor files
	runParamString += "PDB_FILE1=../data/" + filepath.Base(t.Receptor[0]) + "\n"

	// Write receptor list file
	if t.ReceptorList != "" {
		runParamString += "PDB_LIST1=../data" + filepath.Base(t.ReceptorList) + "\n"
	}

	// Write ligand files
	if len(t.Ligand) >= 1 {
		runParamString += "PDB_FILE2=../data/" + filepath.Base(t.Ligand[0]) + "\n"
		nMol++

		// write ligand list files
		if t.LigandList != "" {
			runParamString += "PDB_LIST2=../data/" + filepath.Base(t.LigandList) + "\n"
		}
	}

	runParamString += "N_COMP=" + strconv.Itoa(nMol) + "\n"

	runParamF := filepath.Join(projectDir, "/run.param")
	err := os.WriteFile(runParamF, []byte(runParamString), 0644)
	if err != nil {
		return "", err
	}

	return runParamF, nil
}

// SetupHaddock3Scenario method prepares the scenario for HADDOCK3
//   - Creates the scenario directory
//   - Creates the `run.toml` file
func (t *Target) SetupHaddock3Scenario(wd string, s input.Scenario) (runner.Job, error) {
	glog.Info("Preparing : " + s.Name)
	sPath := filepath.Join(wd, t.ID, "scenario-"+s.Name)
	dataPath := filepath.Join(wd, t.ID, "data")
	_ = os.MkdirAll(sPath, 0755)

	// Handle the ensembles
	if t.ReceptorList != "" {
		ensembleF := filepath.Join(dataPath, t.ID+"-receptor_ens.pdb")
		_ = utils.CreateEnsemble(t.ReceptorList, ensembleF)
		t.Receptor = []string{ensembleF}
	}
	if t.LigandList != "" {
		ensembleF := filepath.Join(dataPath, t.ID+"-ligand_ens.pdb")
		_ = utils.CreateEnsemble(t.LigandList, ensembleF)
		t.Ligand = []string{ensembleF}
	}

	// Generate the run.toml file - it will handle the restraints
	_, err := t.WriteRunToml(sPath, s.Parameters.General, s.Parameters.Modules)
	if err != nil {
		return runner.Job{}, err
	}

	j := runner.Job{
		ID:   t.ID + "_" + s.Name,
		Path: sPath,
	}

	return j, nil
}

// WriteRunToml writes the run.toml file
func (t *Target) WriteRunToml(projectDir string, general map[string]interface{}, mod input.ModuleParams) (string, error) {
	runTomlString := ""
	for k, v := range general {
		switch v := v.(type) {
		case string:
			runTomlString += k + " = \"" + v + "\"\n"
		case int:
			runTomlString += k + " = " + strconv.Itoa(v) + "\n"
		case float64:
			runTomlString += k + " = " + strconv.FormatFloat(v, 'f', -1, 64) + "\n"
		case bool:
			runTomlString += k + " = " + strconv.FormatBool(v) + "\n"
		case []string:
			runTomlString += k + " = [" + strings.Join(v, ",") + "]\n"
		case []int:
			runTomlString += k + " = [" + strings.Join(utils.IntSliceToStringSlice(v), ",") + "]\n"
		case []float64:
			runTomlString += k + " = [" + strings.Join(utils.FloatSliceToStringSlice(v), ",") + "]\n"
		}
	}

	runTomlString += "run_dir = \"run1\"\n"
	runTomlString += "molecules = [\n"
	for _, r := range t.Receptor {
		runTomlString += "    \"../data/" + filepath.Base(r) + "\",\n"
	}
	for _, l := range t.Ligand {
		runTomlString += "    \"../data/" + filepath.Base(l) + "\",\n"
	}
	for _, l := range t.Shape {
		runTomlString += "    \"../data/" + filepath.Base(l) + "\",\n"
	}
	runTomlString += "]\n\n"

	fnameArray := [][]string{
		t.Restraints,
		t.Toppar,
		t.MiscPDB,
	}

	// NOTE: THE ORDER OF THE MODULES IS IMPORTANT!!
	// Range over the modules in the order they are defined
	v := reflect.ValueOf(mod)
	types := v.Type()

	for _, m := range mod.Order {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			name := types.Field(i).Name
			// ------------------------------------------------------------------
			// Important!!
			// Here we match the module string inside the Order field,
			//  with the module name in the struct
			// The Struct does not support `.`, and this is the way the module
			//  will be written in the TOML file.
			// For simplicity, users input the module name with `.`,
			//   but here we convert it to `_` so we can match it with the struct.
			if strings.Contains(m, ".") {
				m = strings.ReplaceAll(m, ".", "_")
			}
			// ------------------------------------------------------------------
			if m == strings.ToLower(name) {
				// Now, change back the module name to the TOML format
				// and add single quotes
				if strings.Contains(m, "_") {
					m = strings.ReplaceAll(m, "_", ".")
					m = "'" + m + "'"
				}
				runTomlString += "[" + m + "]\n"
				for k, v := range field.Interface().(map[string]interface{}) {
					// Find the file to be used as fname
					if strings.Contains(k, "_fname") {
						pattern := regexp.MustCompile(v.(string))

						for _, fArr := range fnameArray {
							fname, err := utils.FindFname(fArr, pattern)
							if err != nil {
								return "", err
							}
							// If fname is not empty, add it to the TOML file
							if fname != "" {
								runTomlString += k + " = \"../data/" + filepath.Base(fname) + "\"\n"
							}
						}

					} else {
						switch v := v.(type) {
						case string:
							runTomlString += k + " = \"" + v + "\"\n"
						case int:
							runTomlString += k + " = " + strconv.Itoa(v) + "\n"
						case float64:
							runTomlString += k + " = " + strconv.FormatFloat(v, 'f', -1, 64) + "\n"
						case bool:
							runTomlString += k + " = " + strconv.FormatBool(v) + "\n"
						case []string:
							runTomlString += k + " = [" + strings.Join(v, ",") + "]\n"
						case []int:
							runTomlString += k + " = [" + strings.Join(utils.IntSliceToStringSlice(v), ",") + "]\n"
						case []float64:
							runTomlString += k + " = [" + strings.Join(utils.FloatSliceToStringSlice(v), ",") + "]\n"
						case []interface{}:
							runTomlString += k + " = [" + strings.Join(utils.InterfaceSliceToStringSlice(v), ",") + "]\n"
						}
					}
				}

			}
		}
	}

	runTomlF := filepath.Join(projectDir, "/run.toml")
	err := os.WriteFile(runTomlF, []byte(runTomlString), 0644)
	if err != nil {
		return "", err
	}

	return runTomlF, nil
}

// LoadDataset loads a dataset from a list file
func LoadDataset(projectDir string, pdbList string, rsuf string, lsuf string, ssuf string) ([]Target, error) {
	var rootRegex *regexp.Regexp
	rootRegex = utils.CreateRootRegex(rsuf, lsuf, ssuf)

	recRegex := regexp.MustCompile(`(.*)` + rsuf)
	ligRegex := regexp.MustCompile(`(.*)` + lsuf)
	shapeRegex := regexp.MustCompile(`(.*)` + ssuf)
	_ = os.MkdirAll(projectDir, 0755)

	file, err := os.Open(pdbList)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	m := make(map[string]Target)
	pdbArr := []string{}
	for s.Scan() {
		line := s.Text()
		if !strings.HasSuffix(line, ".pdb") {
			// This is not a PDB file, ignore
			continue
		}

		var receptor, ligand, root, shape string
		fullPath := line
		basePath := filepath.Base(fullPath)
		// Find root and receptor/ligand names
		match := rootRegex.FindStringSubmatch(basePath)
		if len(match) == 0 {
			// Neither receptor nor ligand, add to a list of PDBs
			pdbArr = append(pdbArr, fullPath)
			continue
		}
		root = match[1]

		basePathCG := utils.ContainsCG(basePath)

		RecMatch := recRegex.MatchString(basePath)
		if RecMatch && !basePathCG {
			receptor = fullPath
		}

		LigMatch := ligRegex.MatchString(basePath)
		if LigMatch && !basePathCG {
			ligand = fullPath
		}

		if ssuf != "" {
			ShapeMatch := shapeRegex.MatchString(basePath)
			if ShapeMatch {
				shape = fullPath
			}
		}

		if entry, ok := m[root]; !ok {
			// create new target

			if receptor != "" {
				m[root] = Target{
					ID:       root,
					Receptor: []string{receptor},
					Ligand:   []string{},
				}
			} else if ligand != "" {
				m[root] = Target{
					ID:       root,
					Receptor: []string{},
					Ligand:   []string{ligand},
				}
			}
		} else {
			// update existing target
			if receptor != "" {
				entry.Receptor = append(entry.Receptor, receptor)
			}
			if ligand != "" {
				entry.Ligand = append(entry.Ligand, ligand)
			}
			if shape != "" {
				entry.Shape = append(entry.Shape, shape)
			}
			m[root] = entry
		}
	}

	// Check if Targets have both receptor and ligand
	for _, v := range m {
		if len(v.Receptor) == 0 || len(v.Ligand) == 0 {
			glog.Warning("Target " + v.ID + " does not have both receptor and ligand")
			// err := errors.New("Target " + v.ID + " does not have both receptor and ligand")
			// return nil, err
		}
	}

	// Read the file again, now looking for restraints and toppars
	// TODO: Optimize this
	_, _ = file.Seek(0, io.SeekStart)
	s = bufio.NewScanner(file)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		for k, v := range m {
			// Handle the restraints
			tblRegex := regexp.MustCompile(`(` + k + `)_.*tbl`)
			tblMatch := tblRegex.FindStringSubmatch(filepath.Base(line))
			if len(tblMatch) != 0 {
				v.Restraints = append(v.Restraints, s.Text())
			}

			// Handle the Toppar
			topparRegex := regexp.MustCompile(`(` + k + `)_.*(top|param)`)
			topparMatch := topparRegex.FindStringSubmatch(line)
			if len(topparMatch) != 0 {
				v.Toppar = append(v.Toppar, s.Text())
			}

			m[k] = v
		}
	}

	// Check if there is a target without a restraint and display a warning
	for _, v := range m {
		if len(v.Restraints) == 0 {
			glog.Warning("Target " + v.ID + " does not have restraints!")
		}
	}

	// Add the misc PDBs
	for _, pdb := range pdbArr {
		for k, v := range m {
			rootRegex := regexp.MustCompile(k + `_`)
			if rootRegex.MatchString(pdb) {
				v.MiscPDB = append(v.MiscPDB, pdb)
			}
			m[k] = v
		}
	}

	arr := []Target{}
	for _, v := range m {
		arr = append(arr, v)
	}
	return arr, nil
}

// CreateDatasetDir creates the dataset folder
func CreateDatasetDir(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_ = os.Mkdir(p, 0755)
	} else {
		return errors.New("Dataset folder already exists: " + p)
	}

	return nil
}

// OrganizeDataset creates the folder structure
//   - Create a ID/data folder
//   - Copy the receptor and ligand files to the data folder
//   - Update the paths in the Target struct
//   - Copy the restraints and toppars to the data folder
func OrganizeDataset(bmPath string, bm []Target) ([]Target, error) {
	var tArr []Target

	for _, t := range bm {
		_ = os.MkdirAll(filepath.Join(bmPath, t.ID, "data"), 0755)

		newT := Target{
			ID: t.ID,
		}

		// Update the paths in the Target struct
		for _, r := range t.Receptor {
			rdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(r))
			newT.Receptor = append(newT.Receptor, rdest)
		}

		for _, l := range t.Ligand {
			ldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(l))
			newT.Ligand = append(newT.Ligand, ldest)
		}

		for _, p := range t.MiscPDB {
			mdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(p))
			newT.MiscPDB = append(newT.MiscPDB, mdest)
		}

		for _, r := range t.Restraints {
			rdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(r))
			newT.Restraints = append(newT.Restraints, rdest)
		}

		for _, toppar := range t.Toppar {
			tdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(toppar))
			newT.Toppar = append(newT.Toppar, tdest)
		}

		for _, shape := range t.Shape {
			sdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(shape))
			newT.Shape = append(newT.Shape, sdest)
		}

		// Copy to the data folder
		dataDir := filepath.Join(bmPath, t.ID, "data")
		objArr := [][]string{t.Receptor, t.Ligand, t.MiscPDB, t.Restraints, t.Toppar, t.Shape}

		for _, obj := range objArr {
			err := utils.CopyFileArrTo(obj, dataDir)
			if err != nil {
				return nil, err
			}
		}

		// Create lists
		if len(newT.Receptor) > 1 {
			l := ""
			for _, r := range newT.Receptor {
				l += "\"" + r + "\"" + "\n"
			}
			receptListFile := filepath.Join(bmPath, t.ID, "data", t.ID+"_receptor.list")
			_ = os.WriteFile(receptListFile, []byte(l), 0644)
			newT.ReceptorList = receptListFile
		}

		if len(newT.Ligand) > 1 {
			l := ""
			for _, r := range newT.Ligand {
				l += "\"" + r + "\"" + "\n"
			}
			ligandListFile := filepath.Join(bmPath, t.ID, "data", t.ID+"_ligand.list")
			_ = os.WriteFile(ligandListFile, []byte(l), 0644)
			newT.LigandList = ligandListFile
		}

		tArr = append(tArr, newT)

	}

	return tArr, nil
}

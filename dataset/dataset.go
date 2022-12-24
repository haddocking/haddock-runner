// Package dataset handles dataset parameters and files
package dataset

import (
	"benchmarktools/input"
	"benchmarktools/runner"
	"benchmarktools/utils"
	"io"
	"strings"

	// "benchmarktools/wrapper/haddock2"
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

// SetupScenario method prepares the scenario
//   - Creates the scenario directory
//   - Creates the run.params file
func (t *Target) SetupScenario(wd string, hdir string, s input.Scenario) (runner.Job, error) {

	sPath := filepath.Join(wd, t.ID, "scenario-"+s.Name)
	glog.Info("Preparing : " + s.Name)
	_ = os.MkdirAll(sPath, 0755)

	// Generate the run.params file
	_, err := t.WriteRunParam(sPath, hdir)
	if err != nil {
		return runner.Job{}, err
	}

	// Find which restraints need to be used
	restraints := input.Restraints{}
	for _, r := range t.Restraints {
		if strings.Contains(r, s.Parameters.Restraints.Ambig) {
			restraints.Ambig = r
		}
		if strings.Contains(r, s.Parameters.Restraints.Unambig) {
			restraints.Unambig = r
		}
	}

	toppar := input.Toppar{}
	for _, t := range t.Toppar {
		if filepath.Ext(t) == ".top" {
			toppar.Top = t
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

func (t *Target) WriteRunParam(projectDir string, haddockDir string) (string, error) {

	var runParamString string

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

	if len(t.Ligand) == 0 {
		err := errors.New("ligand not defined")
		return "", err
	}

	runParamString += "N_COMP=2\n"
	runParamString += "RUN_NUMBER=1\n"
	runParamString += "PROJECT_DIR=" + projectDir + "\n"
	runParamString += "HADDOCK_DIR=" + haddockDir + "\n"

	// Write receptor files
	runParamString += "PDB_FILE1=" + t.Receptor[0] + "\n"

	// Write receptor list file
	if t.ReceptorList != "" {
		runParamString += "PDB_LIST1=" + t.ReceptorList + "\n"
	}

	// Write ligand files
	runParamString += "PDB_FILE2=" + t.Ligand[0] + "\n"

	// write ligand list files
	if t.LigandList != "" {
		runParamString += "PDB_LIST2=" + t.LigandList + "\n"
	}

	runParamF := filepath.Join(projectDir, "/run.param")
	err := os.WriteFile(runParamF, []byte(runParamString), 0644)
	if err != nil {
		return "", err
	}

	return runParamF, nil

}

// LoadDataset loads a dataset from a list file
func LoadDataset(projectDir string, pdbList string, rsuf string, lsuf string) ([]Target, error) {

	rootRegex := regexp.MustCompile(`(.*)(?:` + rsuf + `|` + lsuf + `)`)
	recRegex := regexp.MustCompile(`(.*)` + rsuf)
	ligRegex := regexp.MustCompile(`(.*)` + lsuf)
	_ = os.MkdirAll(projectDir, 0755)

	file, err := os.Open(pdbList)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	m := make(map[string]Target)
	for s.Scan() {
		line := s.Text()
		if !strings.HasSuffix(line, ".pdb") {
			// This is not a PDB file, ignore
			continue
		}

		var receptor, ligand, root string
		fullPath := line
		basePath := filepath.Base(fullPath)
		// Find root and receptor/ligand names
		match := rootRegex.FindStringSubmatch(basePath)
		if len(match) == 0 {
			err := errors.New("root name not found with suffixes " + rsuf + " and " + lsuf)
			return nil, err
		} else {
			root = match[1]
		}

		RecMatch := recRegex.FindStringSubmatch(basePath)
		if len(RecMatch) != 0 {
			receptor = fullPath
		}

		LigMatch := ligRegex.FindStringSubmatch(basePath)
		if len(LigMatch) != 0 {
			ligand = fullPath
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
			m[root] = entry
		}
	}

	// Handle the lists
	for k, v := range m {
		if len(v.Receptor) > 1 {
			l := ""
			for _, r := range v.Receptor {
				l += "\"" + r + "\"" + "\n"
			}
			receptListFile := filepath.Join(projectDir, v.ID+"_receptor.list")
			_ = os.WriteFile(receptListFile, []byte(l), 0644)
			v.ReceptorList = receptListFile
		}
		if len(v.Ligand) > 1 {
			l := ""
			for _, r := range v.Ligand {
				l += "\"" + r + "\"" + "\n"
			}
			ligandListFile := filepath.Join(projectDir, v.ID+"_ligand.list")
			_ = os.WriteFile(ligandListFile, []byte(l), 0644)
			v.LigandList = ligandListFile
		}
		m[k] = v
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
			tblRegex := regexp.MustCompile(`(` + k + `)\w+\.tbl`)
			tblMatch := tblRegex.FindStringSubmatch(line)
			if len(tblMatch) != 0 {
				v.Restraints = append(v.Restraints, s.Text())
			}

			// Handle the Toppar
			topparRegex := regexp.MustCompile(`(` + k + `)\w+\.(top|param)`)
			topparMatch := topparRegex.FindStringSubmatch(line)
			if len(topparMatch) != 0 {
				v.Toppar = append(v.Toppar, s.Text())
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

		for _, r := range t.Receptor {
			rdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(r))
			err := utils.CopyFile(r, rdest)
			if err != nil {
				// os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Receptor = append(newT.Receptor, rdest)

		}
		for _, l := range t.Ligand {
			ldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(l))
			err := utils.CopyFile(l, ldest)
			if err != nil {
				// os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Ligand = append(newT.Ligand, ldest)
		}

		if t.ReceptorList != "" {
			rldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.ReceptorList))
			err := utils.CopyFile(t.ReceptorList, rldest)
			if err != nil {
				// os.RemoveAll(bmPath)
				return nil, err
			}
			os.Remove(t.ReceptorList)
			newT.ReceptorList = rldest
		}

		if t.LigandList != "" {
			lldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.LigandList))
			err := utils.CopyFile(t.LigandList, lldest)
			if err != nil {
				// os.RemoveAll(bmPath)
				return nil, err
			}
			os.Remove(t.LigandList)
			newT.LigandList = lldest
		}

		if len(t.Restraints) > 0 {
			for _, r := range t.Restraints {
				rdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(r))
				err := utils.CopyFile(r, rdest)
				if err != nil {
					// os.RemoveAll(bmPath)
					return nil, err
				}
				newT.Restraints = append(newT.Restraints, rdest)
			}
		}

		if len(t.Toppar) > 0 {
			for _, toppar := range t.Toppar {
				tdest := filepath.Join(bmPath, t.ID, "data", filepath.Base(toppar))
				err := utils.CopyFile(toppar, tdest)
				if err != nil {
					// os.RemoveAll(bmPath)
					return nil, err
				}
				newT.Toppar = append(newT.Toppar, tdest)
			}
		}

		tArr = append(tArr, newT)

	}

	return tArr, nil

}

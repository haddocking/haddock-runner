// Package dataset handles dataset parameters and files
package dataset

import (
	"benchmarktools/input"
	"benchmarktools/utils"

	// "benchmarktools/wrapper/haddock2"
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Target is the target structure
type Target struct {
	ID           string
	Receptor     []string
	ReceptorList string
	Ligand       []string
	LigandList   string
	Ambig        string
	Unambig      string
	HaddockDir   string
	ProjectDir   string
}

// Validate validates the Target checking if
//   - Fields are not empty
//   - Files exist
//   - Files are PDB files
func (t *Target) Validate() error {

	if t.ID == "" {
		return errors.New("Target ID not defined")
	}

	if t.HaddockDir == "" {
		return errors.New("Target Haddock directory not defined")
	}

	if t.ProjectDir == "" {
		return errors.New("Target Project directory not defined")
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

func (t *Target) SetupScenarios(inp *input.Input) error {

	for _, s := range inp.Scenarios {
		fmt.Println("Setting up scenario " + s.Name)

		scenarioPath := filepath.Join(inp.General.WorkDir, t.ID, s.Name)
		_ = os.MkdirAll(scenarioPath, 0755)

		t.ProjectDir = scenarioPath
		t.HaddockDir = inp.General.HaddockDir

		// Generate the run.params file
		_, err := t.WriteRunParam()
		if err != nil {
			return err
		}

	}
	return nil
}

func (t *Target) WriteRunParam() (string, error) {

	var runParamString string

	if t.HaddockDir == "" {
		err := errors.New("haddock directory not defined")
		return "", err
	}

	if t.ProjectDir == "" {
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

	if t.Ambig != "" {
		runParamString += "AMBIG_TBL=" + t.Ambig + "\n"
	}

	if t.Unambig != "" {
		runParamString += "UNAMBIG_TBL=" + t.Unambig + "\n"
	}

	runParamString += "N_COMP=2\n"
	runParamString += "PROJECT_DIR=" + t.ProjectDir + "\n"
	runParamString += "HADDOCK_DIR=" + t.HaddockDir + "\n"

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

	runParamF := filepath.Join(t.ProjectDir, "/run.param")
	err := os.WriteFile(runParamF, []byte(runParamString), 0644)
	if err != nil {
		return "", err
	}

	return runParamF, nil

}

// LoadDataset loads a dataset from a list file
func LoadDataset(l string, rsuf string, lsuf string) ([]Target, error) {

	listFile, err := os.Open(l)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(listFile)

	s.Split(bufio.ScanLines)

	rootRegex := regexp.MustCompile(`(.*)(?:` + rsuf + `|` + lsuf + `)`)
	recRegex := regexp.MustCompile(`(.*)` + rsuf)
	ligRegex := regexp.MustCompile(`(.*)` + lsuf)

	m := make(map[string]Target)
	for s.Scan() {

		var receptor, ligand, root string
		fullPath := s.Text()
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

	for k, v := range m {
		if len(v.Receptor) > 1 {
			l := ""
			for _, r := range v.Receptor {
				l += "\"" + r + "\"" + "\n"
			}
			receptListFile := filepath.Join(v.ProjectDir, v.ID+"_receptor.list")
			_ = os.WriteFile(receptListFile, []byte(l), 0644)
			v.ReceptorList = receptListFile
		}
		if len(v.Ligand) > 1 {
			l := ""
			for _, r := range v.Ligand {
				l += "\"" + r + "\"" + "\n"
			}
			ligandListFile := filepath.Join(v.ProjectDir, v.ID+"_ligand.list")
			_ = os.WriteFile(ligandListFile, []byte(l), 0644)
			v.LigandList = ligandListFile
		}
		m[k] = v
	}

	arr := []Target{}
	for _, v := range m {
		arr = append(arr, v)
	}

	return arr, nil

}

func CreateDatasetDir(p string) error {

	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.Mkdir(p, 0755)
	} else {
		return errors.New("Dataset folder already exists: " + p)
	}

	return nil

}

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
				os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Receptor = append(newT.Receptor, rdest)

		}
		for _, l := range t.Ligand {
			ldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(l))
			err := utils.CopyFile(l, ldest)
			if err != nil {
				os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Ligand = append(newT.Ligand, ldest)
		}

		if t.Ambig != "" {
			adest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.Ambig))
			err := utils.CopyFile(t.Ambig, adest)
			if err != nil {
				os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Ambig = adest
		}

		if t.Unambig != "" {
			udest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.Unambig))
			err := utils.CopyFile(t.Unambig, udest)
			if err != nil {
				os.RemoveAll(bmPath)
				return nil, err
			}
			newT.Unambig = udest
		}

		if t.ReceptorList != "" {
			rldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.ReceptorList))
			err := utils.CopyFile(t.ReceptorList, rldest)
			if err != nil {
				os.RemoveAll(bmPath)
				return nil, err
			}
			os.Remove(t.ReceptorList)
			newT.ReceptorList = rldest
		}

		if t.LigandList != "" {
			lldest := filepath.Join(bmPath, t.ID, "data", filepath.Base(t.LigandList))
			err := utils.CopyFile(t.LigandList, lldest)
			if err != nil {
				os.RemoveAll(bmPath)
				return nil, err
			}
			os.Remove(t.LigandList)
			newT.LigandList = lldest
		}

		tArr = append(tArr, newT)

	}

	return tArr, nil

}

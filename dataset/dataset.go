// Package dataset handles dataset parameters and files
package dataset

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

// Benchmark is the benchmark structure
type Benchmark struct {
	Targets []Target
}

// Target is the target structure
type Target struct {
	ID       string
	Receptor string
	Ligand   string
}

// Validate validates the Target checking if
//   - Fields are not empty
//   - Files exist
//   - Files are PDB files
func (t *Target) Validate() error {

	if t.ID == "" {
		return errors.New("Target ID not defined")
	}

	if t.Receptor == "" {
		return errors.New("Target receptor not defined")
	}

	if t.Ligand == "" {
		return errors.New("Target ligand not defined")
	}

	if _, err := os.Stat(t.Receptor); err != nil {
		return errors.New("Target receptor file not found")
	}

	if _, err := os.Stat(t.Ligand); err != nil {
		return errors.New("Target ligand file not found")
	}

	if filepath.Ext(t.Receptor) != ".pdb" {
		return errors.New("Target receptor file not a PDB file")
	}

	if filepath.Ext(t.Ligand) != ".pdb" {
		return errors.New("Target ligand file not a PDB file")
	}

	return nil

}

// LoadDataset loads a dataset from a list file
func LoadDataset(l string, rsuf string, lsuf string) (*Benchmark, error) {

	listFile, err := os.Open(l)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(listFile)

	s.Split(bufio.ScanLines)

	// (.*)(?:_r_u|_l_u)
	rootRegex := regexp.MustCompile(`(.*)(?:` + rsuf + `|` + lsuf + `)`)
	recRegex := regexp.MustCompile(`(.*)` + rsuf)
	ligRegex := regexp.MustCompile(`(.*)` + lsuf)

	// m := make(map[string]map[string]string)
	// data := make(map[string]Dataset)
	// benchmark := &Benchmark{}
	m := make(map[string]Target)
	for s.Scan() {

		var receptor, ligand, root string

		line := filepath.Base(s.Text())
		// Find root and receptor/ligand names
		match := rootRegex.FindStringSubmatch(line)
		if len(match) == 0 {
			err := errors.New("root name not found with suffixes " + rsuf + " and " + lsuf)
			return nil, err
		} else {
			root = match[1]
		}

		RecMatch := recRegex.FindStringSubmatch(line)
		if len(RecMatch) != 0 {
			receptor = line
		}

		LigMatch := ligRegex.FindStringSubmatch(line)
		if len(LigMatch) != 0 {
			ligand = line
		}

		if entry, ok := m[root]; !ok {
			// create new target
			m[root] = Target{
				ID:       root,
				Receptor: receptor,
				Ligand:   ligand,
			}
		} else {
			// update existing target
			if receptor != "" {
				entry.Receptor = receptor
			}
			if ligand != "" {
				entry.Ligand = ligand
			}
			m[root] = entry
		}
	}

	b := &Benchmark{}
	for _, v := range m {
		b.Targets = append(b.Targets, v)
	}

	return b, nil

}

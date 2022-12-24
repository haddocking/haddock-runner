# `benchmark-tools` for [HADDOCK](https://www.bonvinlab.org/software/haddock2.4/)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

---

## Table of contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Development](#development)

---

### [Introduction](#introduction)

#### What is this repository for?

This repository contains a set of tools to benchmark the performance of HADDOCK2.4.

These can be used to compare the perfomance of HADDOCK against other software
packages, to compare the performance of different versions of HADDOCK, or to
compare the performance of HADDOCK on different hardware.

Additionally it can be used to perform large-scale docking experiments in
different scenarios (parameters), for example:

- You have obtained experimental data for a set of proteins and you want to
  dock them against a set of targets. You want to test different parameters
  to see which one gives the best results.
  - Scenario 1: Use all information
  - Scenario 2: Use only 50% of the information
  - Scenario 3: _ab initio_ docking (without information)

#### How do I get set up?

**To run the benchmarking tools, you need to have a working (local)
installation of HADDOCK2.4.** The software is free for academic use and can be
obtained via [registration](https://www.bonvinlab.org/software/haddock2.4/download/).
More information can be obtained from the [HADDOCK website](https://www.bonvinlab.org/software/haddock2.4/).

For more information on how to install HADDOCK, please refer to the
[documentation](https://www.bonvinlab.org/software/haddock2.4/installation/)
and also to the [HADDOCK.md](HADDOCK.md) file in this repository.

#### Previous version

"Hey, what happened to the previous version of this repository? Where is the
python code?!" - you might ask.

The previous vesion of this repository was indeed written in Python but have been
migrated to Go. The main reason for this is that the python version was slow, not
very efficient (or well designed) also it had no tests...! The Go version is faster,
efficient and easier to maintain - see the code coverage.

However you can still find the Python version as the v0.2.1 tag in this
repository [HERE](https://github.com/haddocking/benchmark-tools/tree/v0.2.1).

### [Installation](#installation)

#### Requirements

- [Go 1.19](https://go.dev/doc/install) or higher

#### Installation

Clone the repository

```bash
git clone https://github.com/haddocking/benchmark-tools.git
cd benchmark-tools
go build -o benchmark-tools
./benchmark-tools --help
```

OR

Use the pre-compiled binaries

\<pending>

#### Usage

- `input.yml`

The input of `benchmark-tools` is a `.yml` file; YAML is a human-readable
data-serialization language. It is commonly used for configuration files and
in applications where data is being stored or transmitted. For more information,
please refer to the [YAML website](https://yaml.org/).

An example file is provided in the `examples` folder (`example_input.yml`) and
also below. Its composed of two main sections, `general` which defines the
general parameters of the benchmarking experiment, and `scenarios` which defines
the different scenarios to be tested.

```yaml
# General parameters
general:
  # Location of the HADDOCK script (see more below)
  executable: /trinity/login/rodrigo/projects/benchmarking/haddock24.sh
  # How many jobs should be executed at a given time
  max_concurrent: 2
  # Location where HADDOCK is installed
  haddock_dir: /trinity/login/abonvin/haddock_git/haddock2.4
  # Pattern used to identify the receptor files
  receptor_suffix: _r_u
  # Pattern used to identify the ligand files
  ligand_suffix: _l_u
  # Location of the input list
  input_list: /trinity/login/rodrigo/projects/benchmarking/input.txt
  # Location of the benchmark output
  work_dir: /trinity/login/rodrigo/projects/benchmarking

# Scenarios of the benchmarking experiment
scenarios:
  # Name can be anything you want to identify the scenario
  - name: true-interface
    # Parameters to be used in the scenario
    parameters:
      # The parameters below are the same as the
      #  ones used in the HADDOCK input file `run.cns`
      run_cns:
        noecv: false
        structures_0: 2
        structures_1: 2
        waterrefine: 2
      # Patterns used to identify the restraints files
      restraints:
        ambig: ti
        unambig: unambig
        hbonds: hb
      # Patterns used to identify the custom topology files
      custom_toppar:
        topology: _ligand.top
        param: _ligand.param

  - name: center-of-mass
    parameters:
      run_cns:
        cmrest: true
        structures_0: 2
        structures_1: 2
        waterrefine: 2
```

- `haddock24.sh`

The `haddock24.sh` script is a wrapper around the HADDOCK2.4 executable.
It is used to run HADDOCK in a given folder, and it is called by
`benchmark-tools` for each scenario. The script is provided in the `examples` folder (`haddock24.sh`) and also below.

**Important: Keep in mind that HADDOCK2.4 runs on Python2.7, which is likely not present in recent systems. For tips on how to install it, please refer to the [PYTHON2.md](PYTHON2.md) file in this repository.**

```bash
#!/bin/bash
#===============================================================
# HADDOCK2.4 wrapper script
#===============================================================

# Export the required environment variables
export HADDOCK="/trinity/login/abonvin/haddock_git/haddock2.4"
export HADDOCKTOOLS="$HADDOCK/tools"
export PYTHONPATH="${PYTHONPATH}:$HADDOCK"

# Command to run HADDOCK
$(which python2.7) $HADDOCK/Haddock/RunHaddock.py
#===============================================================
```

- `input_list.txt`

The input file is a list of the input files to be used in the benchmarking experiment.
This is a simple text file with one line per input. Each line contains the path to one
of the input files. The input files must be:

- `.pdb`: for receptor and ligand files
- `.top`: for custom topology files (used for small-molecules)
- `.param` : for custom parameter files (used for small-molecules)
- `.tbl`: a table file containing the restraints to be used in the docking experiment

This list is parsed by `benchmark-tools` and are identified according to the patterns
 set in `input.yml`. Lines begining with `#` are ignored and can be used to document
 the input list for future reference - in-line comments are not supported.

An example file is provided in the `examples` folder (`example_input_list.txt`)
and also below.

```text
#  Lines starting with # are comments
# ------------------------------------------------------------ #
# Input list
# ------------------------------------------------------------ #
# 1A2K
example/1A2K/1A2K_r_u.pdb
example/1A2K/1A2K_l_u.pdb
example/1A2K/1A2K_ligand.top
example/1A2K/1A2K_ligand.param
example/1A2K/1A2K_ti.tbl
example/1A2K/1A2K_unambig.tbl
# 1GGR
example/1GGR/1GGR_r_u.pdb
example/1GGR/1GGR_l_u_1.pdb
example/1GGR/1GGR_l_u_2.pdb
example/1GGR/1GGR_l_u_3.pdb
example/1GGR/1GGR_l_u_4.pdb
example/1GGR/1GGR_l_u_5.pdb
example/1GGR/1GGR_ti.tbl
# 1PPE
example/1PPE/1PPE_l_u.pdb
example/1PPE/1PPE_r_u.pdb
example/1PPE/1PPE_ti.tbl
example/1PPE/1PPE_hb.tbl
example/1PPE/1PPE_unambig.tbl
# 2OOB
example/2OOB/2OOB_l_u.pdb
example/2OOB/2OOB_r_u.pdb
example/2OOB/2OOB_ti.tbl
example/2OOB/2OOB_hb.tbl
# ------------------------------------------------------------ #
```

**Ensembles**: multiple conformations of a receptor or ligand are also supported,
they need to follow the naming convention: `<root>_<ligand|receptor suffix>_N.pdb`,
where `N` is the ensemble number. For example, if the ligand suffix is `_l_u`,
the ligand files for the first ensemble would be:

```text
<root>_l_u_1.pdb
<root>_l_u_2.pdb
<root>_l_u_3.pdb
```

**Important: Do not provide a multi-model ensemble file,
instead provide the individual models.**

#### [Development](#development)

pending

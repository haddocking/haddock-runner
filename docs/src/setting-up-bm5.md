# Setting up BM5

The Protein-Protein docking benchmark v5 ([Vreven, 2015](https://doi.org/10.1016/j.jmb.2015.07.016)), namely BM5, contains a is a large set of non-redundat high-quality structures, check [here](https://zlab.umassmed.edu/benchmark/) the full set.

The BonvinLab provides a HADDOCK-ready sub-version of the BM5 which can be easily used as input for `haddock-runner`. This version is available the following repository; [github.com/haddocking/BM5-clean](https://github.com/haddocking/BM5-clean). Below we will go over step-by-step instructions on how to use it as input.

## Create a working directory

Create a working directory and change to it;

```bash
mkdir -p ~/projects/benchmarking && cd ~/projects/benchmarking
```

## Download the BM5-clean and create a `bm5-input.list` file

Clone the repository and checkout a version. Note that its always recomended to use a specific version, as the main branch might change and for reproducibility.

As previously mentioned, the `BM5-clean` repository is already an organized sub-version, thus its very simple to create the `bm5-input.list` file with a few bash commands;

```bash
git clone https://github.com/haddocking/BM5-clean.git ~/projects/benchmarking/BM5-clean && \
  cd ~/projects/benchmarking/BM5-clean && \
  git checkout v1.1 && \
  ls ~/projects/benchmarking/BM5-clean/HADDOCK-ready/**/*.{pdb,tbl} | grep -v "ana_scripts\|matched\|cg" | sort > bm5-input.list && \
  cp bm5-input.list ~/projects/benchmarking/ && \
  cd ~/projects/benchmarking
```

## Prepare a `haddock3.sh` script

See below an example of a `haddock3.sh` script that can be used to run HADDOCK3.0 locally;

```bash
#!/bin/bash
source /opt/conda/etc/profile.d/conda.sh
conda activate env
haddock3 "$@"
```

Make sure to make this script executable;

```bash
chmod +x ~/projects/benchmarking/haddock3.sh
```

## Prepare the `bm5.yaml` configuration file

Below is a template for the `bm5.yaml` configuration file using haddock3; keep in mind that this must be adapted to your specific setup!

```yaml
general:
  executable: /home/dev/projects/benchmarking/haddock3.sh
  max_concurrent: 100
  haddock_dir: /opt/haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: /home/dev/projects/benchmarking/bm5-input.list
  work_dir: /home/dev/projects/benchmarking/my-benchmarking

slurm:
  cpus_per_task: 8

scenarios:
  - name: true-interface
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order: [topoaa, rigidbody, seletop, flexref, emref, caprieval]
        topoaa:
          autohis: true
        rigidbody:
          sampling: 1000
          ambig_fname: _ti.tbl
          unambig_fname: _unambig.tbl
          ligand_top_fname: _ligand.top
          ligand_param_fname: _ligand.param
        seletop:
          select: 200
        flexref:
          ambig_fname: _ti.tbl
          unambig_fname: _unambig.tbl
          ligand_top_fname: _ligand.top
          ligand_param_fname: _ligand.param
        emref: ~
        caprieval:
          reference_fname: _ref.pdb

  - name: center-of-mass
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order: [topoaa, rigidbody, seletop, caprieval]
        topoaa:
          autohis: true
        rigidbody:
          sampling: 10000
          cmrest: true
        seletop:
          select: 400
        caprieval:
          reference_fname: _ref.pdb

  - name: random-restraints
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order: [topoaa, rigidbody, seletop, caprieval]
        topoaa:
          autohis: true
        rigidbody:
          sampling: 10000
          ranair: true
        seletop:
          select: 400
        caprieval:
          reference_fname: _ref.pdb
```

## Run the benchmarking

Finally, run the benchmarking with the following command;

```bash
haddock-runner bm5.yaml
```

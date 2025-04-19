# Writing a benchmark.yaml file

The `benchmark.yaml` file is a configuration file in [`YAML`](https://yaml.org/) format that will be used by `haddock-runner` to run the benchmark. The central idea is that one configuration file can define multiple scenarios, each scenario being a set of parameters that will be used to run HADDOCK.

This file should be the replicable part of the benchmark, i.e. the part that you want to share with others. It should contain all the information needed to run the benchmark, alongside the input list.

This file is divided into 3 main sections: [`general`](#general-section), [`slurm`](#slurm-section) and [`scenarios`](#scenario-section).

## General section

Here you must define the following parameters:

- `executable`: Path to the `run-haddock.sh` script (see [here](./writing-a-run-haddock.sh-script.md) for more details)
- `max_concurrent`: Maximum number of jobs that can be executed at a given time
- `haddock_dir`: Path to the HADDOCK installation, this is used to validate the parameters of the [`scenarios`](#scenario-section) section
- `receptor_suffix`: This pattern will identify what is the receptor file
- `ligand_suffix`: This will be used to identify the ligand files
- `shape_suffix`: This will be used to identify shape files
- `input_list`: The path to the input list (see [here](./writing-a-input.list-file.md) for more details)
- `work_dir`: The path where the results will be stored

See below an example:

```yaml
general:
  executable: /workspaces/haddock-runner/example/haddock3.sh
  max_concurrent: 4
  haddock_dir: /opt/haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: /workspaces/haddock-runner/example/input_list.txt
  work_dir: /workspaces/haddock-runner/bm-goes-here
```

## Slurm section

This section is optional but highly recomended! For these to take effect you must be running the benchmark in an HPC environment. These will be used internally by the runner to compose the `.job` file. Here you can define the following parameters or, if left blank, SLURM will pick up the default values:

- `partition`: The name of the partition to be used
- `cpus_per_task`: Number of CPUs per task
- `ntasks_per_node`: Number of tasks per node
- `nodes`: Number of nodes
- `time`: Maximum time for the job
- `account`: Account to be used
- `mail_user`: Email to be notified when the job starts and ends

See below an example:

```yaml
slurm:
  partition: short # use the short partition
  cpus_per_task: 8 # use 8 cores per task
```

## Scenario section

Here you must define the scenarios that you want to run, which are slightly different for [HADDOCK2.4](#haddock24) and [HADDOCK3.0](#haddock30)

### HADDOCK2.4

For HADDOCK2.4 you must define the following:

- `name`: The name of the scenario
- `parameters`: The parameters to be used in the scenario (requires only those that differ from the default)
  - `run_cns`: Parameters that will be used in the `run.cns` file
  - `restraints`: Patterns used to identify the restraints files
    - `ambig`: Pattern used to identify the ambiguous restraints file
    - `unambig`: Pattern used to identify the unambiguous restraints file
    - `hbonds`: Pattern used to identify the hydrogen bonds restraints file
  - `custom_toppar`: Patterns used to identify the custom topology files
    - `topology`: Pattern used to identify the topology file
    - `param`: Pattern used to identify the parameter file

```yaml
# HADDOCK2.4
scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
        structures_0: 1000
        structures_1: 200
        waterrefine: 200
      restraints:
        ambig: ti
```

### HADDOCK3.0

> Note: HADDOCK3.0 is still under development and is not meant to be used for production runs! Please use HADDOCK2.4 instead.
> For information about the available modules, please refer to the [HADDOCK3 tutorial](/education/HADDOCK3/HADDOCK3-antibody-antigen/#a-brief-introduction-to-haddock3) and the [documentation](https://www.bonvinlab.org/haddock3).

For HADDOCK3.0 you must define the following:

- `name`: The name of the scenario
- `parameters`: The parameters to be used in the scenario
  - `general`: General parameters; those are the ones defined in the "top" section of the `run.toml` script
  - `modules`: This subsection is related to the parameters of each module in HADDOCK3.0
    - `order`: The order of the modules to be used in HADDOCK3.0
    - `<module-name>`: Parameters for the module

```yaml
# HADDOCK3.0
scenarios:
  - name: true-interface
    parameters:
      general:
        mode: local
        ncores: 4

      modules:
        order: [topoaa, rigidbody, seletop, flexref, emref]
        topoaa:
          autohis: true
        rigidbody:
          ambig_fname: _ti.tbl
        seletop:
          select: 200
        flexref:
        emref:
```

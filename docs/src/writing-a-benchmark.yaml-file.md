# Writing a benchmark.yaml file

The `benchmark.yaml` file is a configuration file in [`YAML`](https://yaml.org/) format that will be used by `haddock-runner` to run the benchmark. This file is divided in 2 main sections; `general` and `scenarios`

## General section

Here you must define the following:

- `executable`: the path to the `run-haddock.sh` script (see above for more details)
- `max_concurrent`: the maximum number of runs that can be executed at a given time (a run is a target in a given scenario)
- `haddock_dir`: the path to the HADDOCK installation
- `receptor_suffix`: the suffix used to identify the receptor files
- `ligand_suffix`: the suffix used to identify the ligand files
- `input_list`: the path to the input list (see above for more details)
- `work_dir`: the path to the benchmark output

```yaml
general:
  executable: /trinity/login/rodrigo/projects/benchmarking/haddock24.sh
  max_concurrent: 2
  haddock_dir: /trinity/login/abonvin/haddock_git/haddock2.4
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: /trinity/login/rodrigo/projects/benchmarking/input.txt
  work_dir: /trinity/login/rodrigo/projects/benchmarking
```

## Slurm section

Soon to come...

## Scenario section

Here you must define the scenarios that you want to run, it is slightly different for HADDOCK2.4 and HADDOCK3.0.

### HADDOCK2.4

For HADDOCK2.4 you must define the following:

- `name`: the name of the scenario
- `parameters`: the parameters to be used in the scenario
  - `run_cns`: parameters that will be used in the `run.cns` file
  - `restraints`: patterns used to identify the restraints files
    - `ambig`: pattern used to identify the ambiguous restraints file
    - `unambig`: pattern used to identify the unambiguous restraints file
    - `hbonds`: pattern used to identify the hydrogen bonds restraints file
  - `custom_toppar`: patterns used to identify the custom topology files
    - `topology`: pattern used to identify the topology file
    - `param`: pattern used to identify the parameter file

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

- `name`: the name of the scenario
- `parameters`: the parameters to be used in the scenario
  - `general`: general parameters; those are the ones defined in the "top" section of the `run.toml` script
  - `modules`: this subsection is related to the parameters of each module in HADDOCK3.0
    - `order`: the order of the modules to be used in HADDOCK3.0
    - `<module-name>`: parameters for the module

```yaml
# HADDOCK3.0
scenarios:
  - name: true-interface
    parameters:
      general:
        # execution mode using a batch system
        mode: hpc
        # batch queue name to use
        queue: short
        # number of jobs to submit to the batch system
        queue_limit: 100
        # number of models to concatenate within one job
        concat: 5

      modules:
        order: [topoaa, rigidbody, seletop, flexref, emref]
        topoaa:
          autohis: true
        rigidbody:
          ambig_fname: "_ti.tbl"
        seletop:
          select: 200
        flexref:
        emref:
```

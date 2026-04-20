# Writing a benchmark.yaml file

The `benchmark.yaml` file is a configuration file in [`YAML`](https://yaml.org/) format that will be used by `haddock-runner` to run the benchmark. The central idea is that one configuration file can define multiple scenarios, each scenario being a set of parameters that will be used to run HADDOCK.

This file should be the replicable part of the benchmark, i.e. the part that you want to share with others. It should contain all the information needed to run the benchmark, alongside the input list.

This file is divided into 3 main sections: [`general`](#general-section), [`slurm`](#slurm-section) and [`scenarios`](#scenario-section).

## General section

Here you must define the following parameters:

- `mol_suffixes`: Array of suffix patterns to identify molecule files (receptor, ligand, etc.)
- `shape_suffix`: (Optional) Suffix pattern to identify shape files
- `input_list`: The path to the input list (see [here](./writing-a-input.list-file.md) for more details)
- `work_dir`: The path where the results will be stored
- `max_concurrent`: Maximum number of jobs that can be executed at a given time
- `ncores`: Number of CPU cores to use per job
- `execution`: Execution mode, either `local` or `slurm`

> **Important**: haddock-runner requires `haddock3` to be available in your system PATH. 

See below an example:

```yaml
general:
  mol_suffixes: [_r_u, _l_u, _x_u]
  shape_suffix: _shape.pdb  # optional
  input_list: docking/input_list.txt
  work_dir: ../bm-goes-here
  max_concurrent: 100
  ncores: 1
  execution: local # slurm
```

## Scenario section

- `name`: The name of the scenario
- `workflow`: Workflow configuration containing HADDOCK3 modules and their parameters

The workflow section uses a flattened structure where each module is a key, and its parameters are nested under it.

```yaml
scenarios:
  - name: true-interface
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 10
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
        ligand_top_fname: _ligand.top
        ligand_param_fname: _ligand.param
      seletop:
        select: 2
      flexref:
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
        ligand_top_fname: _ligand.top
        ligand_param_fname: _ligand.param
      emref:
      caprieval:
        reference_fname: _ref.pdb

  - name: ab-initio
    workflow:
      topoaa:
      rigidbody:
        sampling: 1
        cmrest: true
```


## Migration from v2 to v3



**v2+:**
```yaml
general:
  executable: /path/to/haddock3.sh
  haddock_dir: /opt/haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  # ... other parameters
```

**v3+:**
```yaml
general:
  mol_suffixes: [_r_u, _l_u]  # Array instead of separate fields
  # haddock3 must be in PATH
  # ... other parameters
```

### Scenario Changes

**v2+:**

```yaml
scenarios:
  - name: true-interface
    parameters:
      general:
        mode: local
        ncores: 4
      modules:
        order: [topoaa, rigidbody, seletop]
        topoaa:
          autohis: true
        rigidbody:
          ambig_fname: _ti.tbl
```

**v3+:**

```yaml
scenarios:
  - name: true-interface
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        ambig_fname: _ti.tbl
      seletop:
```

### Key Differences

1. **No separate executable configuration**: Rust version requires `haddock3` in PATH
2. **Simplified workflow structure**: No need to specify `order` or `general` parameters
3. **Array-based molecule suffixes**: Use `mol_suffixes` array instead of separate receptor/ligand suffixes
4. **Automatic parameter handling**: General parameters like `ncores` are handled automatically

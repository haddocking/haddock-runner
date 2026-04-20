# Examples

Here are examples of the `benchmark.yaml` file 

> **Note**: The current version supports HADDOCK3.0 only. For HADDOCK2.4 examples, please refer to the `v2.x` version .

```yaml
general:
  mol_suffixes: [_r_u, _l_u]
  input_list: example/docking/input_list.txt
  work_dir: ../bm-goes-here
  max_concurrent: 100
  ncores: 1
  execution: slurm

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

  - name: center-of-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1
        cmrest: true

  - name: random-restraints
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 5
        ranair: true

  #-----------------------------------------------
```


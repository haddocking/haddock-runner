# Examples

Here is a full example of the `benchmark.yaml` file for both [HADDOCK2.4](#haddock24) and [HADDOCK3.0](#haddock30).

## `HADDOCK2.4`

```yaml
general:
  executable: /workspaces/haddock-runner/haddock24.sh
  max_concurrent: 2
  haddock_dir: /Users/rodrigo/repos/haddock
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: /workspaces/haddock-runner/example/input_list.txt
  work_dir: /workspaces/haddock-runner/bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      run_cns:
        noecv: false
        structures_0: 1000
        structures_1: 200
        waterrefine: 200
      restraints:
        ambig: ambig
        unambig: restraint-bodies
        hbonds: hbonds
      custom_toppar:
        topology: _ligand.top
        param: _ligand.param

  - name: center-of-mass
    parameters:
      run_cns:
        cmrest: true
        structures_0: 10000
        structures_1: 400
        waterrefine: 400
        anastruc_1: 400
      custom_toppar:
        topology: _ligand.top
        param: _ligand.param

  - name: random-restraints
    parameters:
      run_cns:
        ranair: true
        structures_0: 10000
        structures_1: 400
        waterrefine: 400
        anastruc_1: 400
      custom_toppar:
        topology: _ligand.top
        param: _ligand.param

  #-----------------------------------------------
```

## `HADDOCK3.0`

```yaml
general:
  executable: /workspaces/haddock-runner/example/haddock3.sh
  max_concurrent: 4
  haddock_dir: /opt/haddock3
  receptor_suffix: _r_u
  ligand_suffix: _l_u
  input_list: /workspaces/haddock-runner/example/input_list.txt
  work_dir: /workspaces/haddock-runner/bm-goes-here

scenarios:
  - name: true-interface
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order:
          [topoaa, rigidbody, seletop, flexref, emref, clustfcc, seletopclusts]
        topoaa:
          autohis: true
        rigidbody:
          ambig_fname: _ambig.tbl
          unambig_fname: _restraint-bodies.tbl
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        seletop:
          select: 200
        flexref:
          ambig_fname: _ambig.tbl
          unambig_fname: _restraint-bodies.tbl
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        emref:
          ambig_fname: _ambig
        clustfcc:
        seletopclusts:

  - name: center-of-mass
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order:
          [topoaa, rigidbody, seletop, flexref, emref, clustfcc, seletopclusts]
        topoaa:
          autohis: true
        rigidbody:
          sampling: 10000
          cmrest: true
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        seletop:
          select: 400
        flexref:
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        emref:
        clustfcc:
        seletopclusts:

  - name: random-restraints
    parameters:
      general:
        mode: local
        ncores: 8

      modules:
        order:
          [topoaa, rigidbody, seletop, flexref, emref, clustfcc, seletopclusts]
        topoaa:
          autohis: true
        rigidbody:
          sampling: 10000
          ranair: true
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        seletop:
          select: 400
        flexref:
          contactairs: true
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        emref:
          contactairs: true
          ligand_param_fname: _ligand.param
          ligand_top_fname: _ligand.top
        clustfcc:
        seletopclusts:

  #-----------------------------------------------
```

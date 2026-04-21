# Examples

This page provides comprehensive examples of `haddock-runner` configurations for various benchmarking scenarios. These examples demonstrate the tool's flexibility and help you design your own benchmarks.

## Basic Examples

### Restraint Strategy Comparison

Compare different restraint approaches for the same targets:

```yaml
general:
  mol_suffixes: [_r_u, _l_u]
  input_list: input_list.txt
  work_dir: ./results/restraint-comparison
  max_concurrent: 2
  ncores: 4
  execution: local

scenarios:
  - name: true-interface-restraints
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        ambig_fname: _ti.tbl
      flexref:
        ambig_fname: _ti.tbl
      caprieval:
        reference_fname: _ref.pdb

  - name: hbond-only-restraints
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        ambig_fname: _hb.tbl
      flexref:
        ambig_fname: _hb.tbl
      caprieval:
        reference_fname: _ref.pdb

  - name: center-of-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        cmrest: true
      flexref:
        cmrest: true
      caprieval:
        reference_fname: _ref.pdb

  - name: random-air-restraints
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        ranair: true
      flexref:
        ranair: true
      caprieval:
        reference_fname: _ref.pdb
```

## Advanced Examples

> The purpose of this scenario is to sample antibody-peptide complexes, re-docking experimental structures. Rigid docking, flexible refinement and em refinement. Unambiguous restrains to keep Ab heavy and light chain together and ambiguous for CDR loops and whole peptide. <Jahaciel Villaverde Nogues>

```yaml
general:
  input_list: /trinity/csbdevel/jvillave/deeprank-ab-pep/haddock3/2_real/2_config/data_2_cutoff_043/input_test_scenario_0.list
  mol_suffixes: [_antibody, _antigen]
  work_dir: /trinity/csbdevel/jvillave/deeprank-ab-pep/haddock3/2_real/3_2_results/scenario_0_results
  execution: slurm 
  max_concurrent: 100
  ncores: 24

scenarios:
  # ------------------------------------------------------------
  # 1) scenario 0, ab initio ground truth
  # ------------------------------------------------------------
  - name: ground-truth
    workflow:
      topoaa:
        tolerance : 20
      rigidbody:
        tolerance : 20
        crossdock: false
        sampling: 10000
        ambig_fname: _ti.tbl
        unambig_fname: _antibody-unambig.tbl
      clustfcc:
        plot_matrix: true
      # select up to 100 clusters per target,
      # keeping 5 top models each (max 500 models)
      seletopclusts:
        top_clusters: 100
        top_models: 5
      flexref:
        tolerance : 20
        ambig_fname: _ti.tbl
        unambig_fname: _antibody-unambig.tbl
      # final energy minimisation
      emref:
        tolerance : 20
        ambig_fname: _ti.tbl
        unambig_fname: _antibody-unambig.tbl
      caprieval:
        reference_fname: _matched.pdb
        fnat_cutoff: 4.0
        irmsd_cutoff: 8.0
      emscoring:
        tolerance : 20
        per_interface_scoring: true
```

...

## Configuration Variations

### HPC Cluster Configuration

Optimized for SLURM workload manager:

```yaml
general:
  mol_suffixes: [_r_u, _l_u]
  input_list: large_input_list.txt
  work_dir: /scratch/results/large-benchmark
  max_concurrent: 20
  ncores: 8
  execution: slurm

scenarios:
  - name: hpc-optimized
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 2000
        ambig_fname: _ti.tbl
      flexref:
        ambig_fname: _ti.tbl
      caprieval:
        reference_fname: _ref.pdb
```

### Minimal Configuration

Simple setup for quick testing:

```yaml
general:
  mol_suffixes: [_r_u, _l_u]
  input_list: test_input.txt
  work_dir: ./test-results
  max_concurrent: 2
  ncores: 2
  execution: local

scenarios:
  - name: quick-test
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 100
        cmrest: true
```

## Real-world Example: BM5 Benchmark

A configuration similar to the BM5 benchmark setup:

```yaml
general:
  mol_suffixes: [_r_u, _l_u]
  input_list: bm5_input_list.txt
  work_dir: ./results/bm5-style
  max_concurrent: 10
  ncores: 4
  execution: local

scenarios:
  - name: bm5-true-interface
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
      seletop:
        select: 200
        sort_by: score
      semiflexref:
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
      emref:
        mdsteps: 500
      caprieval:
        reference_fname: _ref.pdb
        clusters: 4

  - name: bm5-center-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        cmrest: true
      seletop:
        select: 200
        sort_by: score
      semiflexref:
        cmrest: true
      emref:
        mdsteps: 500
      caprieval:
        reference_fname: _ref.pdb
        clusters: 4
```

## Input List Examples

### Simple Input List

```text
# Target 1A2K
structures/1A2K/1A2K_r_u.pdb
structures/1A2K/1A2K_l_u.pdb
structures/1A2K/1A2K_ti.tbl
structures/1A2K/1A2K_ref.pdb

# Target 1GGR
structures/1GGR/1GGR_r_u.pdb
structures/1GGR/1GGR_l_u.pdb
structures/1GGR/1GGR_ti.tbl
structures/1GGR/1GGR_ref.pdb
```

### Complex Input List with Multiple File Types

```text
# Target 1PPE - Protein-protein with multiple restraint types
structures/1PPE/1PPE_r_u.pdb
structures/1PPE/1PPE_l_u.pdb
structures/1PPE/1PPE_ti.tbl
structures/1PPE/1PPE_hb.tbl
structures/1PPE/1PPE_unambig.tbl
structures/1PPE/1PPE_ref.pdb

# Target 2OOB - With ligand files
structures/2OOB/2OOB_r_u.pdb
structures/2OOB/2OOB_l_u.pdb
structures/2OOB/2OOB_x_u.pdb
structures/2OOB/2OOB_ti.tbl
structures/2OOB/2OOB_hb.tbl
structures/2OOB/2OOB_ligand.top
structures/2OOB/2OOB_ligand.param
structures/2OOB/2OOB_ref.pdb
```

## Best Practices for Examples

### Starting with Examples

1. **Begin with simple configurations** and gradually add complexity
2. **Test with small datasets** before scaling up
3. **Use `--setup` mode** to validate configurations before full runs
4. **Start with low sampling** numbers for initial testing

### Adapting Examples

- **Modify scenarios** to match your research questions
- **Adjust resource settings** based on your hardware
- **Customize workflows** for your specific docking needs
- **Scale parameters** appropriately for your system size

### Creating Your Own

Use these examples as templates and:

1. Replace file paths with your actual data
2. Adjust sampling parameters for your needs
3. Add or remove workflow steps as required
4. Configure resource limits for your environment

## Troubleshooting Examples

### Common Configuration Issues

**Problem**: Jobs fail with missing file errors
**Solution**: Verify all files in input list exist and paths are correct

**Problem**: Out of memory errors
**Solution**: Reduce `max_concurrent` or increase `ncores` per job

**Problem**: HADDOCK module not found
**Solution**: Ensure HADDOCK3 is properly installed and in PATH

**Problem**: Slow execution
**Solution**: Adjust `max_concurrent` and `ncores` for optimal resource usage

## Additional Resources

- **Complete configuration reference**: [Writing a Benchmark YAML File](./usage.md#step-3-write-the-benchmark-configuration)
- **Input list format guide**: [Writing an Input List File](./usage.md#step-2-create-the-input-list-file)
- **Running benchmarks**: [Running Haddock Runner](./usage.md#step-4-run-the-benchmark)
- **Real-world setup**: [Setting Up a Benchmark](./setting-up-bm5.md)

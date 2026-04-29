# Usage Guide

This guide provides a comprehensive, step-by-step introduction to using `haddock-runner` for running large-scale HADDOCK docking benchmarks. No prior experience with previous versions is assumed.

## Quick Start Workflow

Using `haddock-runner` involves three main steps:

1. **Prepare your input files**
2. **Configure your benchmark**
3. **Run the benchmark**

## Complete Usage Guide

### Step 1: Prepare Your Molecular Data

Before using `haddock-runner`, you need:

- **Protein structures**: PDB files for your docking targets
- **Restraint files** (optional): TBL files for guided docking
- **Topology/parameter files** (optional): For ligands or special molecules

**Organize your files**:

```text
your_project/
├── structures/
│   ├── target1_r_u.pdb    # Receptor structure
│   ├── target1_l_u.pdb    # Ligand structure
│   ├── target1_ti.tbl     # True interface restraints
│   └── target1_ref.pdb    # Reference structure (for evaluation)
└── ...
```

### Step 2: Create the Input List File

The input list file specifies all files needed for each docking target.

**Key points**:

- One target per section (separated by comments)
- List all required files for each target
- Paths can be relative or absolute
- Use consistent naming conventions

**Example** (`input_list.txt`):

```text
# Target 1A2K - Protein-protein complex
structures/1A2K/1A2K_r_u.pdb
structures/1A2K/1A2K_l_u.pdb
structures/1A2K/1A2K_ti.tbl
structures/1A2K/1A2K_unambig.tbl
structures/1A2K/1A2K_ref.pdb

# Target 1GGR - Another complex
structures/1GGR/1GGR_r_u.pdb
structures/1GGR/1GGR_l_u.pdb
structures/1GGR/1GGR_ti.tbl
```

### Step 3: Write the Benchmark Configuration

The YAML configuration file defines your benchmark scenarios and settings.

**Main sections**:

- `general`: Global settings (concurrency, resources, directories)
- `scenarios`: Different docking workflows to test
- Each scenario defines a complete HADDOCK workflow

**Example** (`benchmark.yaml`):

```yaml
general:
  max_concurrent: 4        # How many jobs to run simultaneously
  ncores: 2               # CPU cores per job
  execution: local        # Execution mode (local, slurm, etc.)
  mol_suffixes: [_r_u, _l_u]  # File name suffixes for molecules
  input_list: input_list.txt  # Path to your input list file
  work_dir: ./results     # Where to store results

scenarios:
  - name: true-interface
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        ambig_fname: _ti.tbl
      flexref:
        ambig_fname: _ti.tbl
      caprieval:
        reference_fname: _ref.pdb

  - name: center-of-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        cmrest: true
```

See [Configuration Reference](./reference.md) for complete configuration options.

### Step 4: Run the Benchmark

Execute `haddock-runner` with your configuration:

```bash
# Basic execution
haddock-runner benchmark.yaml

# Setup mode (validate without running)
haddock-runner --setup benchmark.yaml

# Debug mode (verbose logging)
haddock-runner --debug benchmark.yaml
```

**What happens during execution**:

1. Input validation and checksum verification
2. Job creation for each target-scenario combination
3. Concurrent execution according to resource limits
4. Results organization in the working directory
5. Progress logging and error handling

See [Running the Benchmark](#step-4-run-the-benchmark) for runtime details.

### Step 5: Analyze Results

After completion, results are organized by scenario and target:

```text
results/
├── true-interface/
│   ├── 1A2K/
│   │   ├── haddock3.cfg
│   │   ├── run1/
│   │   └── ...
│   └── 1GGR/
│       └── ...
└── center-of-mass/
    ├── 1A2K/
    └── 1GGR/
        └── ...
```

**Result analysis tips**:

- Compare docking success rates between scenarios
- Analyze CAPRI metrics for quality assessment
- Examine computation times and resource usage
- Use HADDOCK analysis tools for detailed evaluation

## Practical Tips

### Starting Small

For your first benchmark:

- Use 2-3 well-characterized targets
- Test 2 different scenarios
- Start with small sampling numbers (100-500)
- Use `--setup` mode to validate before full execution

### Resource Management

- **Memory**: Each job needs ~2-4GB RAM
- **CPU**: Allocate cores based on your system capacity
- **Storage**: Results can be large (1-10GB per target)
- **Time**: Docking runs can take hours to days

### Common Workflows

**Parameter optimization**:

```yaml
scenarios:
  - name: sampling-500
    workflow:
      rigidbody:
        sampling: 500
  - name: sampling-1000
    workflow:
      rigidbody:
        sampling: 1000
  - name: sampling-2000
    workflow:
      rigidbody:
        sampling: 2000
```

**Restraint strategy comparison**:

```yaml
scenarios:
  - name: true-interface
    workflow:
      rigidbody:
        ambig_fname: _ti.tbl
  - name: hbond-only
    workflow:
      rigidbody:
        ambig_fname: _hb.tbl
  - name: center-of-mass
    workflow:
      rigidbody:
        cmrest: true
```

## Troubleshooting

**Common issues and solutions**:

**Input file errors**:

- Verify all files exist and are readable
- Check file paths in your input list
- Use absolute paths if relative paths don't work

**HADDOCK module errors**:

- Ensure HADDOCK3 is properly installed
- Verify all required modules are available
- Check your HADDOCK3 configuration

**Resource limitations**:

- Reduce `max_concurrent` if running out of memory
- Lower sampling numbers for faster testing
- Use `--setup` to validate before full runs

**Permission issues**:

- Ensure write access to working directory
- Check execution permissions for the binary
- Verify HADDOCK3 has proper file access

## Best Practices

### File Organization

```text
benchmark_project/
├── configs/
│   ├── benchmark.yaml
│   └── input_list.txt
├── structures/
│   ├── target1/
│   ├── target2/
│   └── ...
├── results/
│   └── (auto-generated)
└── analysis/
    └── (your analysis scripts)
```

### Version Control

- Keep configuration files in Git
- Store input structures separately (large files)
- Document changes between benchmark runs
- Use meaningful commit messages

### Reproducibility

- Fix random seeds when comparing methods
- Document exact HADDOCK3 version used
- Record system specifications
- Archive complete configuration files

## Next Steps

Now that you understand the basic workflow:

1. **Set up your first benchmark** → [Setting Up a Benchmark](./setting-up-bm5.md)
2. **Explore example configurations** → [Examples](./examples.md)
3. **Learn about advanced features** → [Development](./development.md)
4. **Get help with specific issues** → [Getting Help](#getting-help)

## Getting Help

If you encounter any issues:

- Check the [Troubleshooting](#troubleshooting) section above
- Consult the [GitHub Issues](https://github.com/haddocking/haddock-runner/issues)
- Review the [HADDOCK3 documentation](https://github.com/haddocking/haddock3)
- Contact the support team via the channels mentioned in the main documentation

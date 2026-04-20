# Setting Up a BM5 Benchmark: Step-by-Step Guide

This guide provides comprehensive, up-to-date instructions for setting up and running a BM5 (Protein-Protein Docking Benchmark v5) benchmark using `haddock-runner`. The BM5 benchmark ([Vreven, 2015](https://doi.org/10.1016/j.jmb.2015.07.016)) is a widely-used set of 144 non-redundant, high-quality protein-protein complexes for evaluating docking methods.

## Prerequisites

Before starting, ensure you have:

- `haddock-runner` installed (see [Installation](../installation.md))
- HADDOCK3 properly installed and configured
- Access to a computing environment with sufficient resources
- Basic familiarity with command-line tools

## Step 1: Set Up Your Project Directory

Create a dedicated directory structure for your benchmark:

```bash
# Create project directory
mkdir -p ~/bm5-benchmark && cd ~/bm5-benchmark

# Create subdirectories
mkdir -p {data,configs,results,scripts}
```

Your project structure will look like:

```
bm5-benchmark/
├── data/          # BM5 dataset files
├── configs/       # Configuration files
├── results/       # Benchmark results (auto-created)
├── scripts/       # Custom scripts
└── README.md      # Your notes and documentation
```

## Step 2: Download and Prepare BM5 Dataset

The BonvinLab provides a HADDOCK-ready version of BM5:

```bash
# Clone the BM5-clean repository
git clone https://github.com/haddocking/BM5-clean.git ~/bm5-benchmark/data/BM5-clean

# Check out a specific version for reproducibility
git checkout v1.1

# Create input list file
find ~/bm5-benchmark/data/BM5-clean/HADDOCK-ready -name "*.pdb" -o -name "*.tbl" \
  | grep -E "(r_u|_l_u|_ti|_unambig|_ref)" \
  | sort > ~/bm5-benchmark/configs/bm5-input.list
```

## Step 3: Create the Benchmark Configuration

Create a modern `bm5-benchmark.yaml` configuration file:

```yaml
# File: ~/bm5-benchmark/configs/bm5-benchmark.yaml
general:
  # File patterns and locations
  mol_suffixes: [_r_u, _l_u]          # Standard BM5 naming
  input_list: configs/bm5-input.list   # Path to input list
  work_dir: results/bm5-results       # Where to store results
 
  # Resource management
  max_concurrent: 8                  # Adjust based on your system
  ncores: 4                          # Cores per HADDOCK job
  execution: local                   # Use 'slurm' for HPC clusters

scenarios:
  # Scenario 1: True Interface
  - name: true-interface
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
      flexref:
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
      emref:
        mdsteps: 500
      caprieval:
        reference_fname: _ref.pdb
        clusters: 4

  # Scenario 2: Center of Mass
  - name: center-of-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 2000
        cmrest: true
      seletop:
        select: 200
        sort_by: score
      flexref:
        cmrest: true
      emref:
        mdsteps: 500
      caprieval:
        reference_fname: _ref.pdb
        clusters: 4

  # Scenario 3: Random Air Restraints
  - name: random-restraints
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 2000
        ranair: true
      seletop:
        select: 200
        sort_by: score
      flexref:
        ranair: true
      emref:
        mdsteps: 500
      caprieval:
        reference_fname: _ref.pdb
```

## Step 4: Validate Your Setup

Before running the full benchmark, validate your configuration:

```bash
# Check that haddock-runner is working
haddock-runner --version

# Validate configuration without execution
haddock-runner --setup configs/bm5-benchmark.yaml

# Check input file count
wc -l configs/bm5-input.list
# Should show ~1000-1500 files for full BM5
```

## Step 5: Run the Benchmark

```bash
# Run with progress monitoring
nohup haddock-runner configs/bm5-benchmark.yaml > benchmark.log 2>&1 &

# Monitor progress
tail -f benchmark.log

# Check resource usage
htop  # or your preferred system monitor
```

### HPC Cluster Execution

For SLURM clusters, modify your config:

```yaml
general:
  execution: slurm
  cpus_per_task: 4
```

## Step 6: Monitor and Manage the Benchmark

### Monitoring Progress

```bash
# Check running jobs
ps aux | grep haddock

# For SLURM
squeue -u $USER

# Check disk usage
du -sh results/
```

### Handling Interruptions

If the benchmark is interrupted:

```bash
# Check what completed
find results/ -name "*.done" | wc -l

# Resume from where it left off
haddock-runner configs/bm5-benchmark.yaml
```

## Step 7: Analyze Results

To be added soon.

## Troubleshooting

### Common Issues and Solutions

**Problem**: "File not found" errors

- **Solution**: Verify all paths in `bm5-input.list` are correct
- **Check**: `head configs/bm5-input.list` and verify files exist

**Problem**: HADDOCK3 module errors

- **Solution**: Ensure HADDOCK3 is properly installed and in PATH
- **Check**: `haddock3 --version` works from command line

**Problem**: Out of memory errors

- **Solution**: Reduce `max_concurrent` or increase system memory
- **Check**: Monitor memory with `free -h` or `htop`

**Problem**: Slow progress

- **Solution**: Adjust `max_concurrent` and `ncores` for optimal balance
- **Check**: Monitor CPU usage with `htop`

## Best Practices

### Reproducibility

```bash
# Record exact versions
 echo "haddock-runner $(haddock-runner --version)" > VERSION.txt
 echo "HADDOCK3 $(haddock3 --version)" >> VERSION.txt
 echo "Date: $(date)" >> VERSION.txt

# Save complete configuration
cp configs/bm5-benchmark.yaml results/config-used.yaml
```

### Data Management

```bash
# Compress completed results
 tar -czvf bm5-results-$(date +%Y%m%d).tar.gz results/

# Clean up intermediate files (if needed)
 find results/ -name "*.tmp" -delete
```

### Documentation

```markdown
# BM5 Benchmark Notes

## Setup
- Date: YYYY-MM-DD
- System: Describe your hardware
- HADDOCK3 version: X.X.X
- haddock-runner version: X.X.X

## Configuration
- Scenarios: true-interface, center-of-mass, random-restraints
- Sampling: 1000-2000 per scenario
- Targets: Full BM5 set (144 complexes)

## Results Summary
- Start time: 
- End time: 
- Total runtime: 
- Success rate: 
```

## Complete Example: From Start to Finish

Here's a complete workflow example:

```bash
# 1. Set up project
mkdir -p ~/bm5-benchmark/{data,configs,results,scripts} && cd ~/bm5-benchmark

# 2. Get data
git clone https://github.com/haddocking/BM5-clean.git data/BM5-clean
find data/BM5-clean/HADDOCK-ready -name "*.pdb" -o -name "*.tbl" \
  | grep -E "(r_u|_l_u|_ti|_unambig|_ref)" \
  | sort > configs/bm5-input.list

# 3. Create configuration (use the YAML example above)

# 4. Run full benchmark
haddock-runner configs/bm5-benchmark.yaml

# 5. Analyze results
# to be updated
```

## Additional Resources

- **BM5 Original Publication**: [Vreven et al. (2015)](https://doi.org/10.1016/j.jmb.2015.07.016)
- **BM5 Dataset**: [https://zlab.umassmed.edu/benchmark/](https://zlab.umassmed.edu/benchmark/)
- **BM5-clean Repository**: [https://github.com/haddocking/BM5-clean](https://github.com/haddocking/BM5-clean)
- **HADDOCK3 Documentation**: [https://github.com/haddocking/haddock3](https://github.com/haddocking/haddock3)

## Getting Help

If you encounter issues specific to BM5 setup:

- Check the [BM5-clean issues](https://github.com/haddocking/BM5-clean/issues)
- Consult the [HADDOCK forum](https://ask.bioexcel.eu)
- Review the [haddock-runner issues](https://github.com/haddocking/haddock-runner/issues)
- Contact the BonvinLab support team

This guide provides a complete, up-to-date approach to setting up BM5 benchmarks with the current version of `haddock-runner`, focusing on clarity, reproducibility, and practical execution.

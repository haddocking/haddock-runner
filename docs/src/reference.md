# Configuration Reference

This document provides a comprehensive reference for the `haddock-runner` configuration YAML file format. It describes all available options, their purpose, valid values, and examples.

> Keep in mind that YAML format is indentation-sensitive!

## Configuration File Structure

The configuration file is a YAML document with two main sections:

```yaml
general:
  # Global configuration options

scenarios:
  # Benchmark scenarios to execute
```

---

## General Configuration

The `general` section contains global settings that apply to all scenarios and targets.

### Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `max_concurrent` | integer | Yes | Maximum number of jobs to run simultaneously. Controls how many target-scenario combinations execute in parallel. |
| `ncores` | integer | Yes | Number of CPU cores to allocate per job. |
| `execution` | string | Yes | Execution backend. Valid values: `local`, `slurm`. |
| `mol_suffixes` | array of strings | Yes | File suffixes used to identify molecule files. Must contain at least 2 suffixes (typically receptor and ligand). |
| `input_list` | string | Yes | Path to the input list file containing file paths for all targets. |
| `work_dir` | string | Yes | Directory where benchmark results will be stored. Created automatically if it doesn't exist. |

### Example

```yaml
general:
  max_concurrent: 4
  ncores: 2
  execution: local
  mol_suffixes: [_r_u, _l_u, _x_u]
  input_list: docking/input_list.txt
  work_dir: ./results
```

### Notes

- **Local execution**: When using `execution: local`, the total number of CPU cores required is `max_concurrent * ncores`. Ensure your system has enough cores.
- **SLURM execution**: When using `execution: slurm`, ensure SLURM is installed and configured. The `sbatch` and `sacct` commands must be available in your PATH.
- **File suffixes**: The `mol_suffixes` array defines patterns used to identify molecule files in the input list. Files matching these patterns are grouped together as molecules for each target.

---

## Scenarios Configuration

The `scenarios` section defines the different docking workflows to test. Each scenario is executed for every target specified in the input list.

### Scenario Options

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | string | Yes | Unique identifier for this scenario. Used as directory name in the results. |
| `workflow` | mapping | Yes | HADDOCK3 workflow configuration defining modules and their parameters. |

### Example

```yaml
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

---

## Input List File Format

The input list file (specified by `general.input_list`) contains paths to all files required for each docking target. Files are automatically grouped into targets by a shared identifier derived from the filename: for molecule files, the identifier is the part of the filename before the configured `mol_suffixes` match; for restraints, topology/parameter, shape, and miscellaneous files, grouping typically uses the part before the first underscore.

### File Classification

Files in the input list are automatically categorized based on their extensions and patterns:

| File Type | Pattern | Description |
|-----------|---------|-------------|
| Molecules | Matches `mol_suffixes` patterns |  Structure files (PDB format) |
| Restraints | `_*.tbl` | Distance restraint files |
| Topology/Parameters | `.top`, `.param` | Topology and parameter files for ligands |
| Shape | `_shape*` or configured pattern | Shape files for shape-based docking |
| Miscellaneous | All other files | Any additional files (reference structures, etc.) |

### Example Input List

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

### Notes

- Lines starting with `#` are treated as comments and ignored.
- Empty lines are ignored.
- Paths can be relative to the configuration file location or absolute.
- Files are grouped by their root identifier which is extracted by splitting on underscore, taking the first part.

---

## HADDOCK3 Workflow Modules

The `workflow` section within each scenario defines the HADDOCK3 modules to execute and their parameters. Each module is specified as a YAML key, with its parameters as nested key-value pairs.

> IMPORTANT! You should not set any of haddock's "General Default Parameters" - these are handled by `haddock-runner` internally!

### Haddock Module Patterns

Look in the [haddock repository](https://github.com/haddocking/haddock3) for information about modules/parameters for each module.

### Module Parameter Patterns

Many parameters accept **filename patterns** instead of explicit paths. These patterns are matched against the files available for each target. The pattern matching uses regular expressions.

Common filename patterns:

| Pattern | Matches |
|---------|---------|
| `_ti.tbl` | Files ending with `_ti.tbl` |
| `_unambig.tbl` | Files ending with `_unambig.tbl` |
| `_ref.pdb` | Files ending with `_ref.pdb` |
| `_ligand.top` | Files ending with `_ligand.top` |
| `_ligand.param` | Files ending with `_ligand.param` |

**Note:** When using filename patterns, ensure the corresponding files are listed in the input list and have consistent naming conventions across all targets.

---

## Complete Configuration Examples

### Basic Benchmark

```yaml
general:
  max_concurrent: 2
  ncores: 2
  execution: local
  mol_suffixes: [_r_u, _l_u]
  input_list: input_list.txt
  work_dir: ./results

scenarios:
  - name: standard
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
      flexref:
      emref:
```

### Parameter Optimization with Shape Docking

```yaml
general:
  max_concurrent: 4
  ncores: 4
  execution: slurm
  mol_suffixes: [_r_u, _l_u, _shape]
  input_list: shape/input.txt
  work_dir: shape-results

scenarios:
  - name: sampling-500
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 500
        mol_shape_3: true

  - name: sampling-1000
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        mol_shape_3: true

  - name: sampling-2000
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 2000
        mol_shape_3: true
```

### Restraint Strategy Comparison

```yaml
general:
  max_concurrent: 2
  ncores: 2
  execution: local
  mol_suffixes: [_r_u, _l_u]
  input_list: input_list.txt
  work_dir: restraint-comparison

scenarios:
  - name: true-interface
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        ambig_fname: _ti.tbl
        unambig_fname: _unambig.tbl
      flexref:
        ambig_fname: _ti.tbl
      caprieval:
        reference_fname: _ref.pdb

  - name: center-of-mass
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        cmrest: true
      flexref:
      caprieval:
        reference_fname: _ref.pdb

  - name: random-restart
    workflow:
      topoaa:
        autohis: true
      rigidbody:
        sampling: 1000
        ranair: true
      flexref:
      caprieval:
        reference_fname: _ref.pdb
```

---

## Validation Rules

The configuration file is validated before execution. The following rules apply:

### General Section

1. **`mol_suffixes`**: Must be a non-empty array with at least 2 entries.
2. **`mol_suffixes`**: Must contain unique values (no duplicates).
3. **`work_dir`**: Must not be an empty string.
4. **`input_list`**: Must not be an empty string, and the file must exist.
5. **`max_concurrent`**: Must be greater than 0.
6. **`ncores`**: Must be greater than 0.

### Local Execution

When `execution: local`:

- `max_concurrent * ncores` must not exceed the available CPU cores on the system.

### SLURM Execution

When `execution: slurm`:

- The `sbatch` and `sacct` commands must be available in the system PATH.

---

## File Resolution

### Path Resolution

- **Relative paths** in the configuration file (for `input_list` and `work_dir`) are resolved relative to the current working directory.

### Filename Pattern Resolution

When a module parameter ends with `_fname` and contains a pattern (e.g., `_ti.tbl`), the pattern is matched against all files available for the target. The matching is done using regular expressions.

**IMPORTANT: If multiple files match the pattern, the match is treated as ambiguous and the resolver returns `None`, so the parameter is omitted from the generated run TOML.**

---

## Directory Structure

After running a benchmark, results are organized as follows:

```text
work_dir/
├── scenario1/
│   ├── target1/
│   │   ├── run1/
│   │   │   └── ... (HADDOCK3 output)
│   │   └── job.sh (only for SLURM execution)
│   └── target2/
│       ├── run1/
│       └── job.sh
└── scenario2/
    ├── target1/
    └── target2/
```

---

## Tips for Configuration

1. **Start small**: Begin with a few targets and simple scenarios to validate your setup.
2. **Use `--setup` mode**: Always run with `haddock-runner --setup configuration.yaml` first to validate the configuration before full execution.
3. **Check file patterns**: Ensure your `mol_suffixes` patterns correctly match your molecule filenames.
4. **Resource planning**: Calculate total CPU requirements as `max_concurrent * ncores` and ensure your system can handle it.
5. **Consistent naming**: Use consistent file naming conventions across all targets for filename patterns to work correctly.

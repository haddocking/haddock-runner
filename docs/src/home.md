# Welcome to the `haddock-runner` docs

![image](./banner_home-mini.jpg)

The `haddock-runner` is a powerful tool for running large-scale HADDOCK docking experiments. It automates the execution of HADDOCK3 workflows across multiple protein complexes, enabling comprehensive benchmarking and performance evaluation.

HADDOCK (High Ambiguity Driven protein-protein DOCKing) is a widely-used software suite for flexible docking of biomolecular complexes, particularly useful for studying protein-protein interactions.

## Key Features

- **Large-scale Benchmarking**: Execute HADDOCK workflows on multiple molecular complexes simultaneously
- **Scenario Testing**: Run different docking scenarios (workflows, parameters) on the same datasets
- **Concurrent Execution**: Process multiple targets concurrently for efficient resource utilization
- **Input Validation**: Automatic checksum validation to ensure data integrity
- **Flexible Configuration**: YAML-based configuration for complex benchmarking setups

## How It Works

`haddock-runner` takes a YAML configuration file that defines:

- **General settings**: Maximum concurrent jobs, core allocation, working directory
- **Input datasets**: List of molecular structures and associated files
- **Docking scenarios**: Different HADDOCK workflows and parameters to test

The tool then automatically:

1. Validates all input files using checksums
2. Creates individual HADDOCK jobs for each target-scenario combination
3. Executes jobs concurrently according to resource constraints
4. Organizes results in a structured working directory

## Quick Start

### Prerequisites

- HADDOCK3 installed and properly configured
- Input molecular structures in PDB format
- Optional restraint files (TBL format) for guided docking

### Basic Usage

```bash
haddock-runner benchmark_config.yaml
```

### Common Options

Setup mode (validate and prepare without execution):

```bash
haddock-runner --setup benchmark_config.yaml
```

Debug mode (verbose logging):

```bash
haddock-runner --debug benchmark_config.yaml
```

## Typical Use Cases

When running benchmarks, researchers typically investigate:

- **Parameter Optimization**: How different sampling parameters affect docking quality
- **Workflow Comparison**: Performance of different docking protocols
- **Method Validation**: Testing new restraint strategies or scoring functions
- **Performance Benchmarking**: Execution time and resource usage patterns
- **Reproducibility Studies**: Consistent results across different computational environments

## Example Workflow

A typical benchmark might include:

- 5-10 different protein complexes
- 3-5 different docking scenarios (true interface, center-of-mass, random restraints)
- 100-1000 docking runs per scenario
- Concurrent execution on 4-8 CPU cores

Results are organized by scenario and target, making it easy to compare performance across different conditions.

## Getting Started with Your Own Benchmark

1. **Prepare your molecular structures** in PDB format
2. **Create restraint files** if using guided docking
3. **Write a configuration file** defining your scenarios
4. **List your input files** in the required format
5. **Run the benchmark** and analyze results

See the [Setting Up a Benchmark](/setting-up-bm5) and [Writing a Benchmark YAML File](/usage#step-3-write-the-benchmark-configuration) sections for detailed instructions.

## Getting Help

If you encounter any issues or have questions:

- Open an issue on the [GitHub repository](https://github.com/haddocking/haddock-runner)
- Contact us at _bonvinlab.support@uu.nl_
- Join the [BioExcel forum](https://ask.bioexcel.eu) and post your question

The HADDOCK team and community are available to help with setup, configuration, and analysis of your benchmarks.

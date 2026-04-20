# Development Guide

This guide provides comprehensive information for developers contributing to `haddock-runner`, including setup instructions, architecture overview, coding standards, and CI/CD workflows.

## Project Overview

`haddock-runner` is a Rust application for running large-scale HADDOCK docking benchmarks. It features:

- **Modern Rust stack**: Using Cargo, clap, and serde
- **Concurrent execution**: Multi-threaded job processing
- **Flexible configuration**: YAML-based benchmark definitions
- **Multiple backends**: Local, SLURM, and other HPC integrations

## Getting Started with Development

### Prerequisites

- **Rust toolchain**: Latest stable version
- **HADDOCK3**: For testing and development
- **Git**: For version control
- **Slurm**: To develop HPC integration

### Development Environment Setup

```bash
# Install Rust toolchain
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env

# Clone the repository
git clone https://github.com/haddocking/haddock-runner.git
cd haddock-runner

# Install development tools
rustup component add rustfmt clippy rust-analysis
cargo install cargo-edit cargo-audit
```

## Project Structure

```text
src/
├── main.rs              # Main application entry point
├── input.rs             # Input file parsing and validation
├── dataset.rs           # Dataset loading and management
├── job.rs               # Job creation and management
├── queue.rs             # Job queue and scheduling
├── runner/              # Execution backends
│   ├── mod.rs           # Runner interface
│   ├── local.rs         # Local execution backend
│   ├── slurm.rs         # SLURM backend
│   └── status.rs       # Job status tracking
├── logging.rs           # Logging configuration
├── checksum.rs          # File integrity checking
└── utils.rs             # Utility functions

Cargo.toml              # Rust package configuration
Cargo.lock              # Dependency versions
.example/               # Example configurations
.docs/                  # Documentation
.github/workflows/      # CI/CD workflows
```

## Development Workflow

### Building the Project

```bash
# Build in development mode
cargo build

# Build with optimizations
cargo build --release

# Build with all features
cargo build --all-features
```

### Running Tests

```bash
# Run all tests
cargo test

# Run tests with coverage (requires tarpaulin)
cargo tarpaulin

# Run specific test
cargo test test_specific_function
```

### Code Quality

```bash
# Format code
cargo fmt

# Check for linting issues
cargo clippy

# Audit dependencies for vulnerabilities
cargo audit

# Check for outdated dependencies
cargo outdated
```

## Architecture Overview

### Core Components

1. **Input System**: Parses YAML configuration and validates inputs
2. **Dataset Manager**: Loads and organizes molecular data
3. **Job Creator**: Generates HADDOCK jobs from scenarios
4. **Queue System**: Manages job scheduling and execution
5. **Runner Backends**: Local, SLURM, and other execution environments
6. **Monitoring**: Tracks job progress and status

### Data Flow

```text
Input Files → Configuration Parsing → Dataset Loading → Job Creation → Queue Scheduling → Job Execution → Result Collection
```

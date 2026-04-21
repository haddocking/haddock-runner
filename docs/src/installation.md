# Installation

The `haddock-runner` is designed for researchers, developers, and advanced users who are familiar with HADDOCK and command-line computing. It is particularly suited for those with access to HPC infrastructure for running large-scale docking experiments.

## Prerequisites

### HADDOCK3 Installation

> **IMPORTANT**: `haddock-runner` requires HADDOCK3 to be installed on your system.
>
> This tool is **not** a replacement for HADDOCK itself, but rather a benchmarking framework that automates the execution of multiple HADDOCK runs.

If you are new to HADDOCK, we recommend:

- Completing the basic [HADDOCK3 tutorials](/education/HADDOCK3/index.md)
- Familiarizing yourself with HADDOCK3 workflows and configuration

For single target docking or small-scale experiments, consider using:

- [HADDOCK2.4 web server](https://wenmr.science.uu.nl/haddock2.4/) for interactive use
- HADDOCK3 command-line interface for small batches

### System Requirements

- **Operating System**: Linux (recommended), macOS, or Windows with WSL
- **Memory**: Minimum 8GB RAM (16GB+ recommended for concurrent execution)
- **Storage**: Sufficient disk space for input structures and results
- **HPC Access**: Recommended for large-scale benchmarks

## Installation Methods

### Method 1: Install via crates.io (Recommended)

The easiest way to install `haddock-runner` is through cargo, Rust's package manager:

```bash
# Install directly from crates.io
cargo install haddock-runner

# This will install the binary to ~/.cargo/bin/haddock-runner
```

> **Note**: If you don't have cargo installed, you can install Rust from [https://www.rust-lang.org/tools/install](https://www.rust-lang.org/tools/install)

After installation, ensure the cargo bin directory is in your PATH:

```bash
# Add cargo bin to your PATH (add this to your ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.cargo/bin:$PATH"

# Verify installation
source $HOME/.cargo/env
haddock-runner --version
```

### Method 2: Install Pre-built Binary from GitHub Releases (Coming Soon)

Pre-compiled binaries will be available for each release on GitHub:

```bash
# Download the latest release for your platform
# Check https://github.com/haddocking/haddock-runner/releases for the latest version
VERSION="v3.0.0"  # Update to latest version
OS_ARCH="x86_64-unknown-linux-gnu"  # Choose your platform

wget https://github.com/haddocking/haddock-runner/releases/download/${VERSION}/haddock-runner-${OS_ARCH}

# Make it executable
chmod +x haddock-runner-${OS_ARCH}

# Move to your PATH (optional)
sudo mv haddock-runner-${OS_ARCH} /usr/local/bin/haddock-runner

# Verify installation
haddock-runner --version
```

Available platforms will include:

- `x86_64-unknown-linux-gnu` (Linux 64-bit)
- `x86_64-apple-darwin` (macOS Intel)
- `aarch64-apple-darwin` (macOS Apple Silicon)

> **Note**: Pre-built binaries are coming soon. For now, please use Method 1 (crates.io) or see the [Development](/development) section for building from source.

## Post-Installation Setup

### Add to PATH (Optional)

To make `haddock-runner` available system-wide:

```bash
# Create a symlink or copy the binary to a directory in your PATH
sudo ln -s $(pwd)/target/release/haddock-runner /usr/local/bin/haddock-runner

# Verify it's accessible
which haddock-runner
haddock-runner --version
```

### Verify HADDOCK3 Integration

Before running benchmarks, ensure HADDOCK3 is properly installed and accessible:

```bash
# Check HADDOCK3 installation
haddock3 --version

# Verify required modules are available
haddock3 --list-modules
```

## Troubleshooting

### Common Issues

**Rust installation problems**:

- Ensure you have proper internet connectivity
- Check that you have required system dependencies (`build-essential`, `curl`, etc.)
- Try `rustup update` if you already have Rust installed

**Missing HADDOCK3**:

- Ensure HADDOCK3 is installed and in your PATH
- Check that all required HADDOCK modules are available
- Verify your HADDOCK3 configuration files are properly set up

**Permission issues**:

- Ensure you have read/write access to the working directory
- Check that input files are readable
- Verify you have execution permissions for the binary

### Getting Help

If you encounter installation issues:

- Check the [GitHub Issues](https://github.com/haddocking/haddock-runner/issues) for known problems
- Consult the [HADDOCK3 documentation](https://github.com/haddocking/haddock3) for HADDOCK-specific requirements
- Contact the development team via the support channels mentioned in the [Getting Help](/getting-help) section

## Next Steps

Now that you have `haddock-runner` installed, you're ready to:

1. **Set up your first benchmark** - See [Setting Up a Benchmark](/setting-up-bm5)
2. **Write a configuration file** - See [Writing a Benchmark YAML File](/usage#step-3-write-the-benchmark-configuration)
3. **Prepare your input files** - See [Writing an Input List File](/usage#step-2-create-the-input-list-file)
4. **Run your benchmark** - See [Running Haddock Runner](/usage#step-4-run-the-benchmark)

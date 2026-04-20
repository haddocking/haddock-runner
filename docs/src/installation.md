# Installation

The tool is designed for users/students/developers that are familiar with
HADDOCK and command-line scripting, and have access to an HPC infrastructure.

If this is the first time you are using HADDOCK, please first familiarize
yourself with the software by running the basic [HADDOCK3](/education/HADDOCK3/index.md) tutorials.

This tool is not meant to be used by end-users who want to run a single target,
or a small set of targets; for that purpose we recommend instead using
the [HADDOCK2.4 web server](https://wenmr.science.uu.nl/haddock2.4/).

> **VERY IMPORTANT**: You need to have HADDOCK installed on your system.
> This is not covered in this documentation.
> Please refer to the [HADDOCK3.0 repository](https://github.com/haddocking/haddock3) for more information.

`haddock-runner` is a standalone open-source software licensed under Apache 2.0
and freely available from the following repository: [github.com/haddocking/haddock-runner](https://github.com/haddocking/haddock-runner).

The Rust version requires Rust to be installed. First, install Rust if you haven't already:

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env

# Clone the repository and build
$ git clone https://github.com/haddocking/haddock-runner.git
$ cd haddock-runner
$ cargo build --release
$ ./target/release/haddock-runner --version
haddock-runner 3.0.0
```


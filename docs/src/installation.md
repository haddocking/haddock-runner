# Installation

The tool is designed for users/students/developers that are familiar with HADDOCK, command-line scripting and with access to a HPC infrastructure.

If this is the first time you are using HADDOCK, please familiarize first yourself with the software by running the basic [HADDOCK2.4](/education/HADDOCK24/index.md) or [HADDOCK3](/education/HADDOCK3/index.md) tutorials.

This tool is not meant to be used by end-users who want to run a single target, or a small set of targets; for that purpose we recommend instead using the [HADDOCK2.4 web server](https://wenmr.science.uu.nl/haddock2.4/).

> **VERY IMPORTANT**: You need to have HADDOCK installed on your system.
> This is not covered in this documentation.
> Please refer to the [HADDOCK2.4 installation instructions](/software/haddock2.4/installation) or [HADDOCK3.0 repository](https://github.com/haddocking/haddock3) for more information.

`haddock-runner` is a standalone open-source software licensed under Apache 2.0 and freely available from the following repository: [github.com/haddocking/haddock-runner](https://github.com/haddocking/haddock-runner).

To us it simply download the latest binary from the [releases page](https://github.com/haddocking/haddock-runner/releases):

```bash
$ wget https://github.com/haddocking/haddock-runner/releases/download/v1.8.1/haddock-runner_1.8.1_linux_386.tar.gz
$ tar -zxvf haddock-runner_1.8.1_linux_386.tar.gz
$ ./haddock-runner -version
haddock-runner version v1.8.1
```

Alternatively, you can build the latest version from source (you probably don't need to do that), make sure [go](https://go.dev/doc/install) is installed and run the following commands:

```bash
$ git clone https://github.com/haddocking/haddock-runner.git
$ cd haddock-runner
$ go build -o haddock-runner
$ ./haddock-runner -version
haddock-runner version v1.8.1
```

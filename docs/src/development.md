# Development

The [code repository](https://github.com/haddocking/haddock-runner) contains a DevContainer configuration that can be used to set up a development environment, have a look at [Developing inside a Container](https://code.visualstudio.com/docs/devcontainers/containers) for more information.

The only caveat is that a `cns` binary must be in the `.devcontainer` path before building the container.

The development container comes pre-configured with Go, Haddock3 and Slurm.

This is the recommended way to develop the tool, as it will ensure that the development environment is consistent across different platforms with minimal setup.

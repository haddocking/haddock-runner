# How does it work?

The execution of a HADDOCK benchmark consists of a few steps:

1. Setup the benchmark

   - Copy the target structures to the location where the HADDOCK run will be executed

2. Setup the HADDOCK run
   - For HADDOCK2.4, writing the `run.param` file and executing the `haddock2.4` program once to setup the folder structure
   - For HADDOCK3, writing the `run.toml`
3. Distribute several HADDOCK runs in a HPC-friendly manner

`haddock-runner` aim to automate all these steps, additionally giving the user the possibility of setting up various _scenarios_. A scenario is a set of parameters that will be used to run HADDOCK. For example, a user may want to run HADDOCK against a set of targets with different sampling values, different restraints, different parameters, etc.

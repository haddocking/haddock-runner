# Running haddock-runner

Considering the config input file and the config `.yaml` file have been properly set, you can run the benchmark by executing the `haddock-runner` simply with:

```bash
haddock-runner my-benchmark-config-file.yml
```

`haddock-runner` will read the input file, create the working directory, copy the input files to a `data/` directory and start the benchmark. Make sure you have enough space in your disk to store the input files and the results.

**VERY IMPORTANT:** In the current version, `haddock-runner` does not submit jobs to the queue, instead it leverages the internal scheduling routines of HADDOCK2.4/HADDOCK3.0. This means that the number of concurrent runs is related to the number of docking runs at a given time, not to the total number of processors being used by HADDOCK! The actual number of processors being used depends on how HADDOCK was configured. For HADDOCK2.4 this depends on parameters defined in the `run.cns` (`queue_N`/`cpunumber_N`) and for HADDOCK3, the number of processors (or queue slots) to use and the running mode is defined in the config file under the `general` section (see examples above).

**Example:** `max_concurrent: 10` with `scenarios.parameters.mode: local` and `scenarios.parameters.ncores: 10` means 10x10 processors will be required!

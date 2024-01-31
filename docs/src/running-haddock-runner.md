# Running haddock-runner

Assuming the config input file and the config `.yaml` file have been properly set, you can run the benchmark by executing the `haddock-runner` simply with:

```bash
haddock-runner my-benchmark-config-file.yaml
```

`haddock-runner` will read the input file, create the working directory, copy the input files to a `data/` directory and start the benchmark. Make sure you have enough space in your disk to store the input files and the results.

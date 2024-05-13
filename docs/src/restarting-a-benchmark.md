# Restarting a benchmark

In [v1.7.0](https://github.com/haddocking/haddock-runner/releases/tag/v1.7.0) we introduced the possibility to restart a benchmark. This is useful when you want to continue a benchmark that was interrupted for some reason. To restart a benchmark you must have the `benchmark.yaml` file and the `input.list` file used in the original benchmark. The `benchmark.yaml` file must have the `work_dir` parameter set to the directory where the original benchmark was run.

Just run it again without the need of any special flags or parameters:

```bash
haddock-runner my-benchmark-config-file.yaml
```

`haddock-runner` should automagically detect which runs are completed and which are not. It does this by searching the log produced by Haddock (both v2 and v3) and based on keywords it will assign a status to it during runtime.

```text
I1116 13:28:09.721754   58085 main.go:192] ############################################
W1116 13:28:09.721797   58085 main.go:207] +++ 2OOB_true-interface is INCOMPLETE - restarting +++
I1116 13:28:09.721810   58085 main.go:204] 1GGR_center-of-mass - DONE - skipping
I1116 13:28:09.721823   58085 main.go:204] 1A2K_random-restraints - DONE - skipping
W1116 13:28:09.721988   58085 main.go:207] +++ 1GGR_true-interface is INCOMPLETE - restarting +++
I1116 13:28:09.721999   58085 main.go:204] 1GGR_random-restraints - DONE - skipping
I1116 13:28:09.722030   58085 main.go:204] 1A2K_center-of-mass - DONE - skipping
I1116 13:28:09.722010   58085 main.go:204] 1PPE_random-restraints - DONE - skipping
W1116 13:28:09.722072   58085 main.go:207] +++ 1PPE_true-interface is INCOMPLETE - restarting +++
W1116 13:28:09.722087   58085 main.go:207] +++ 1A2K_true-interface is INCOMPLETE - restarting +++
I1116 13:28:09.722165   58085 main.go:204] 1PPE_center-of-mass - DONE - skipping
I1116 13:28:09.722041   58085 main.go:204] 2OOB_center-of-mass - DONE - skipping
I1116 13:28:09.722483   58085 main.go:204] 2OOB_random-restraints - DONE - skipping
I1116 13:28:57.531951   58085 main.go:226] 2OOB_true-interface - DONE in 47.81 seconds
I1116 13:29:46.939726   58085 main.go:226] 1GGR_true-interface - DONE in 97.22 seconds
I1116 13:29:56.830500   58085 main.go:226] 1PPE_true-interface - DONE in 107.11 seconds
I1116 13:30:40.741859   58085 main.go:226] 1A2K_true-interface - DONE in 151.02 seconds
I1116 13:30:40.741907   58085 main.go:235] ############################################
```

To make sure the results are consistent, it will create a checksum of both the configuration yaml and of the input txt and show you a warning. This ensures that parameters and input has not changed mid-execution.

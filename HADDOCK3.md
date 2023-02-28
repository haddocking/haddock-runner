# `benchmark-tools` and [HADDOCK3](https://github.com/haddocking/haddock3)

IMPORTANT: HADDOCK v3.0 is still under development and is
**not ready for production use**!! Please use HADDOCK v2.4.

The input file for `haddock3` is slightly different than the one for
`haddock2.4`. This is needed to accomodate for the new features of HADDOCK3.
More information about its new features and how to use it can be found in the
[HADDOCK3 repository](https://github.com/haddocking/haddock3).

## Input file

The input file for `haddock3` is also a YAML file; see the example below:

```yaml
general:
  # ...
  # same as haddock2.4
  # ...

scenarios:
  - name: true-interface

    parameters:
      general:
        # these are the general parameters defined in haddock3's
        #  `run.toml` header - there's currently limited support for this
        mode: local
        ncores: 1

      modules:
        # these are the modules to be executed in the order they are defined
        order: [topoaa, rigidbody, seletop, flexref, emref]
        # these are the parameters for each module
        topoaa:
          autohis: true

        rigidbody:
          sampling: 10
          ambig_fname: "_ti.tbl"

        seletop:
          select: 2

        # if no custom parameters are being used, it can be empty
        flexref:

        emref:
```

## Executable script

Here you also need to create an executable script that will
call `haddock3`, see the example below:

```bash
#!/bin/bash
#===============================================================================

source $HOME/software/miniconda3/etc/profile.d/conda.sh
conda activate haddock3

# mind the $@ at the end, it passes all the arguments to haddock3
haddock3 "$@"
```

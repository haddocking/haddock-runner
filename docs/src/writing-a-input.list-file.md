# Writing a `input.list` file

The input list is a flat text file with the paths of the targets;

```text
# input.list
/home/rodrigo/projects/haddock-benchmark/data/complex1_r_u.pdb
/home/rodrigo/projects/haddock-benchmark/data/complex1_l_u.pdb
/home/rodrigo/projects/haddock-benchmark/data/complex1_ti.tbl
#
# comments are allowed, use it to organize your file
#
/home/rodrigo/projects/haddock-benchmark/data/complex2_r_u.pdb
/home/rodrigo/projects/haddock-benchmark/data/complex2_l_u.pdb
/home/rodrigo/projects/haddock-benchmark/data/complex2_ti.tbl
/home/rodrigo/projects/haddock-benchmark/data/complex2_ligand.top
/home/rodrigo/projects/haddock-benchmark/data/complex2_ligand.param
```

Note that this file **must** follow the pattern:

```text
path/to/the/structure/NAME_receptor_suffix.pdb
path/to/the/structure/NAME_ligand_suffix.pdb
```

In the above example, `complex1` and `complex2` correspond thus to `NAME`, identifying the complex which is modelled.
Each PDB file (indicated by the `.pdb` extension) has a **suffix**, this is extremely important as it will be used to organize the data. For example, the file `complex1_r_u.pdb` is the receptor of the target `complex1` and `complex1_l_u` is the ligand of the same target.

In this example the suffixes are:

- `receptor_suffix: _r_u`
- `ligand_suffix: _l_u`

These suffixes are defined in the `benchmark.yaml` file, see [here](./writing-a-benchmark.yaml-file.md) for more details.

The same logic applies to the restraints files, in the example above the pattern for the ambiguous restraint can be defined as `ambig: "ti"`, so the file `complex1_ti.tbl` will be used as the ambiguous restraint for the target `complex1`, `complex2_ti.tbl` for the target `complex2`, etc. See section 3.2.2 for information specific to the definition of restraints when setting up a HADDOCK3.0 run.

HADDOCK supports many modified amino acids/bases/glycans/ions (check the [full list](https://wenmr.science.uu.nl/haddock2.4/library)). However if your target molecule is not present in this library, you can also provide it following the same logic; `topology: "_ligand.top"` and `param: "_ligand.param"` will use the files `protein2_ligand.top` and `protein2_ligand.param` for the target `protein2`.

> **IMPORTANT**: For ensembles, **provide each model individually** and append a number to the suffix, for example: `complex1_l_u_1.pdb`, `complex1_l_u_2.pdb`, etc.

See below a full example of the `input.list` file

```text
# -------------------------------- #
# 1A2K
./example/1A2K/1A2K_r_u.pdb
./example/1A2K/1A2K_l_u.pdb
./example/1A2K/1A2K_ligand.top
./example/1A2K/1A2K_ligand.param
./example/1A2K/1A2K_ti.tbl
./example/1A2K/1A2K_unambig.tbl
# 1GGR
./example/1GGR/1GGR_r_u.pdb
./example/1GGR/1GGR_l_u_1.pdb
./example/1GGR/1GGR_l_u_2.pdb
./example/1GGR/1GGR_l_u_3.pdb
./example/1GGR/1GGR_l_u_4.pdb
./example/1GGR/1GGR_l_u_5.pdb
./example/1GGR/1GGR_ti.tbl
# 1PPE
./example/1PPE/1PPE_l_u.pdb
./example/1PPE/1PPE_r_u.pdb
./example/1PPE/1PPE_ti.tbl
./example/1PPE/1PPE_hb.tbl
./example/1PPE/1PPE_unambig.tbl
# 2OOB
./example/2OOB/2OOB_l_u.pdb
./example/2OOB/2OOB_r_u.pdb
./example/2OOB/2OOB_ti.tbl
./example/2OOB/2OOB_hb.tbl
# -------------------------------- #
```

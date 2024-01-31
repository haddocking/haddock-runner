# Setting up BM5

The Protein-Protein docking benchmark v5 ([Vreven, 2015](https://doi.org/10.1016/j.jmb.2015.07.016)), namely BM5, contains a is a large set of non-redundat high-quality structures, check [here](https://zlab.umassmed.edu/benchmark/) the full set.

The BonvinLab provides a HADDOCK-ready sub-version of the BM5 which can be easily used as input for `haddock-runner`. This version is available the following repository; [github.com/haddocking/BM5-clean](https://github.com/haddocking/BM5-clean). Below we will go over step-by-step instructions on how to use it as input.

## Downloading the BM5-clean

Clone the repository and checkout a version. Note that its always recomended to use a specific version, as the main branch might change and for reproducibility.

```text
git clone https://github.com/haddocking/BM5-clean.git && cd BM5-clean && git checkout v1.1
```

## Create `bm5-input.list`

As previously mentioned, the `BM5-clean` repository is already an organized sub-version, thus its very simple to create the `bm5-input.list` file with a few bash commands;

```bash
$ ls `pwd`/HADDOCK-ready/**/*.{pdb,tbl} | grep -v "ana_scripts\|matched\|cg" | sort > bm5-input.list

# The file should look like this
$ head bm5-input.list && echo "..." && tail bm5-input.list
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_ambig.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_ambig5.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_hbonds.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_l_u.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_l_u_GDP.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_r_u.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_restraint-bodies.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1A2K/1A2K_target.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1ACB/1ACB_ambig.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/1ACB/1ACB_ambig5.tbl
...
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/BP57/BP57_r_u.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/BP57/BP57_restraint-bodies.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/BP57/BP57_target.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_ambig.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_ambig5.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_hbonds.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_l_u.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_r_u.pdb
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_restraint-bodies.tbl
/Users/rvhonorato/projects/benchmarking/BM5-clean/HADDOCK-ready/CP57/CP57_target.pdb
```

### 3. Add `bm5-input.list` to the `haddock-runner` configuration file

```yaml
general:
  # ...
  input_list: /Users/rvhonorato/projects/benchmarking/BM5-clean/bm5-input.list
  # ...
```

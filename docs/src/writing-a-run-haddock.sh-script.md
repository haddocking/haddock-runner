# Writing a run-haddock.sh script

The `run-haddock.sh` script is a bash script that will be executed by `haddock-runner` for each target. The purpose of this script is to provide an "adapter" to account for different HADDOCK versions and/or different python versions and even different operating systems and configurations on your cluster.

This script should contain all the commands necessary to run HADDOCK and it must be customized for your installation, for example:

`haddock24.sh`

```bash
#!/bin/bash
#===============================================================================
# HADDOCK2.4 runs on python2.7, which is EOL.
# This script is a workaround to run HADDOCK with a custom python2 installation

## With pyenv
export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
pyenv shell 2.7.18

## With Anaconda
# source $HOME/miniconda3/etc/profile.d/conda.sh
# conda create -n haddock24_env python=2.7
# conda activate haddock24_env

#===============================================================================
# Configure HADDOCK2.4
export HADDOCK="$HOME/repos/haddock24"
export HADDOCKTOOLS="$HADDOCK/tools"
export PYTHONPATH="${PYTHONPATH}:$HADDOCK"

python "$HADDOCK/Haddock/RunHaddock.py"
#===============================================================================
```

`haddock3.sh`

```bash
#!/bin/bash
#===============================================================================
HADDOCK3_DIR="$HOME/repos/haddock3"

# Activate the virtual environment
source "$HADDOCK3_DIR/venv/bin/activate" || exit

### Or if installed with conda
## source $HOME/miniconda3/etc/profile.d/conda.sh
## conda activate haddock3

# Mind the "$@" at the end, this is necessary to pass the arguments to the script
haddock3 "$@"
#===============================================================================
```

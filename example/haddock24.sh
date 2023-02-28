#!/bin/bash
#===============================================================================
HADDOCK24_DIR="$HOME/repos/haddock24"

### Activate the virtual environment
## if your haddock24 installation uses venv
# source "$HOME/repos/haddock24/venv/bin/activate" || exit
## if your haddock24 installation uses conda
conda activate haddock24

### Set the environment variables
export HADDOCK="$HADDOCK24_DIR"
export HADDOCKTOOLS="$HADDOCK/tools"
export PYTHONPATH="${PYTHONPATH}:$HADDOCK"

#######################################################################
# IMPORTANT: HADDOCK2.4 might fail with exit code 0 even if it fails
#######################################################################

python "$HADDOCK/Haddock/RunHaddock.py"
#===============================================================================

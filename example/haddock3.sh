#!/bin/bash
#===============================================================================
HADDOCK3_DIR="$HOME/repos/haddock3"

### Activate the virtual environment
## if your haddock3 installation uses venv
source "$HADDOCK3_DIR/venv/bin/activate" || exit
## if your haddock3  installation uses conda
# conda activate haddock3

haddock3 "$@"
#===============================================================================

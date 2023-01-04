#!/bin/bash
#===============================================================================

HADDOCK3_DIR="$HOME/repos/haddock3"

# shellcheck source=/dev/null
source "$HADDOCK3_DIR/venv/bin/activate" || exit
# if your haddock3 installation uses conda
# conda init
# source $HOME/anaconda3/bin/activate
# source activate haddock3
if [ $# -eq 0 ]; then
	echo "No arguments supplied"
	echo "Usage: haddock3.sh <input_file>"
	exit 0
fi

haddock3 "$@"

#===============================================================================

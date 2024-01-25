#!/bin/bash
#===============================================================================

# shellcheck source=/dev/null
source /opt/conda/etc/profile.d/conda.sh
conda activate env

haddock3 "$@"
#===============================================================================

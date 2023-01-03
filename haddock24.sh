#!/bin/bash
#===============================================================================
# HADDOCK2.4 runs on python2.7, which is EOL.
# This script is a workaround to run HADDOCK with a custom python2 installation
export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
pyenv shell 2.7.18

#===============================================================================
# Configure HADDOCK2.4
export HADDOCK="$HOME/repos/haddock24"
export HADDOCKTOOLS="$HADDOCK/tools"
export PYTHONPATH="${PYTHONPATH}:$HADDOCK"

########
# IMPORTANT: HADDOCK2.4 might fail if there's no `run.cns` or `run.param` files
#   in the current directory. However it will still exit with a 0 status code.
########

python "$HADDOCK/Haddock/RunHaddock.py"

#===============================================================================

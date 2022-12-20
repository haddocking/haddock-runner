#!/bin/bash
export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"

pyenv shell 2.7.18
python /Users/rodrigo/repos/haddock/Haddock/RunHaddock.py

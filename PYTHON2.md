# Python 2.7.18 installation tips and tricks

---

Python2 has reached its [end-of-life](https://www.python.org/doc/sunset-python-2/)
in 2020 but some software, such as HADDOCK still requires it. Here are some
tips on how to install it.

## MacOs

### out-of-the-box

Depending on the version of your system, python2 might be installed already.
Check it with:

```bash
python --version
```

or

```bash
python2 --version
```

### Homebrew

If its not, you might be able to install it with brew, check the full
instructions [here](https://docs.python-guide.org/starting/install/osx/).

```bash
brew install python@2
```

### Pyenv

If that also does not work you can use [`pyenv`](https://github.com/pyenv/pyenv).
After installing pyenv you need to install the python2.7.18 version;

```bash
pyenv install 2.7.18
```

And to activate/use this version you need to add the following to your
`haddock24.sh` script:

```bash
export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
pyenv shell 2.7.18

# $ python --version
# Python 2.7.18
```

### From source

If you want to install python2 from source, you can try your luck with
the instructions [here](https://www.python.org/downloads/release/python-2718/).

---

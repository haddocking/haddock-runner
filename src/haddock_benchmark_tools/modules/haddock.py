import ast
import logging
import os
import pathlib
import shlex
import subprocess
import sys
import tempfile
from pathlib import Path

from haddock_benchmark_tools.modules.errors import HaddockError

haddocklog = logging.getLogger("setuplog")


def is_py3(file_path):
    """Check if code is Python3 compatible."""
    # https://stackoverflow.com/a/40886697
    code_data = open(file_path, "rb").read()
    try:
        ast.parse(code_data)
    except SyntaxError:
        return False
    return True


class HaddockWrapper:
    """Wrapper for HADDOCK."""

    def __init__(self, haddock_path, py2):
        self.path = pathlib.Path(haddock_path)
        self.py2 = pathlib.Path(py2)
        try:
            self.run_haddock = list(Path(self.path).glob("*/*addock.py"))[0]
        except KeyError:
            raise HaddockError(f"{self.path} does not contain Haddock")

        if is_py3(self.run_haddock):
            self.haddock_exec = shlex.split(f"{sys.executable} {self.run_haddock}")
        else:
            self.haddock_exec = shlex.split(f"{self.py2} {self.run_haddock}")
        self.env = os.environ.copy()
        self.env["PYTHONPATH"] = self.path
        self.pid = None

    def check_if_executable(self):
        """Check if Haddock can be executed."""

        out_f = tempfile.NamedTemporaryFile(delete=False, suffix=".out")
        err_f = tempfile.NamedTemporaryFile(delete=False, suffix=".err")

        p = subprocess.Popen(
            self.haddock_exec, env=self.env, stdout=out_f, stderr=err_f
        )
        p.communicate()

        # close the file so we can access it via the system,
        #  we might need to show this to the user
        out_f.close()

        if "run.cns OR run.param" not in open(out_f.name).read():
            raise HaddockError(out_f.name)
        else:
            # All good, make the temporary file disapear
            os.unlink(out_f.name)
            os.unlink(err_f.name)

    def setup(self, target_dir, identifier):
        """Execute Haddock in 'setup' mode."""

        prev_loc = pathlib.Path.cwd()
        os.chdir(target_dir)

        output_f = target_dir / f"haddock.out-{identifier}"
        with open(output_f, "w") as out:
            p = subprocess.Popen(self.haddock_exec, env=self.env, stdout=out)
            p.communicate()
        out.close()

        os.chdir(prev_loc)

        # check for errors
        with open(output_f, "r") as out_fh:
            for line in out_fh.readlines():
                if "already exists => HADDOCK stopped" in line:
                    raise HaddockError(line)
                if "could not" in line:
                    raise HaddockError(line)
                if "does not contain an END statement" in line:
                    raise HaddockError(line)

        return output_f


class HaddockJob(HaddockWrapper):
    def __init__(self, haddock_path, py2, run_path):
        HaddockWrapper.__init__(self, haddock_path, py2)
        self.path = run_path
        self.status = None
        self.process = None
        self.output = run_path / "haddock.out"
        self.error = run_path / "haddock.err"

    def update_status(self):
        """Check the status of the process."""
        try:
            if self.process.poll() is None:
                self.status = "running"
            else:
                if (self.path / "structures/it1/water/file.list").exists():
                    self.status = "complete"
                else:
                    self.status = "failed"

        except AttributeError:
            # this job has not been initiated
            self.status = "null"

        return self.status

    def run(self):
        """Execute the Haddock job."""
        os.chdir(self.path)
        self.process = subprocess.Popen(
            self.haddock_exec,
            env=self.env,
            stdout=open(self.output, "w"),
            stderr=open(self.error, "w"),
        )

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
        self.process = None
        self.output = Path(run_path, "haddock.out")
        self.error = Path(run_path, "haddock.err")
        self.size = self._get_size()
        self.name = run_path

    def status(self):
        """Check the status of the process."""
        current_status = ""
        try:
            if self.process.poll() is None:
                current_status = "running"
            else:
                if Path(self.path, "structures/it1/water/file.list").exists():
                    current_status = "complete"
                else:
                    current_status = "failed"

        except AttributeError:
            # this job has not been initiated
            current_status = "null"

        return current_status

    def run(self):
        """Execute the Haddock job."""
        os.chdir(self.path)
        self.process = subprocess.Popen(
            self.haddock_exec,
            env=self.env,
            stdout=open(self.output, "w"),
            stderr=open(self.error, "w"),
        )

    def _get_size(self):
        """Assign a size to the job."""
        # The size of the input
        input_size = sum(file.stat().st_size for file in Path(self.path).rglob("*"))
        # Here we can also consider the "restrain complexity",
        #  but that is not simply a measure of how many lines there are in the file
        #  for example assign ( segid B ) (segid A) 2.0 2.0 0.0
        #  will calculate ALL the distances between A and B and its only one line.
        return input_size

    def __lt__(self, other):
        return self.size < other.size

    def __gt__(self, other):
        return self.size > other.size

    def __eq__(self, other):
        return self.size == other.size

    def __hash__(self):
        return id(self)

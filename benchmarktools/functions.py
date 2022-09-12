import copy
import logging
import os
import re
from pathlib import Path
from typing import Optional

from benchmarktools.modules.errors import InvalidParameter

log = logging.getLogger("bmlog")


def glob_pdb_re(target_path, regex):
    """Apply a regex in Path.glob function"""
    found_list = []
    for pdb_path in target_path.glob("*pdb"):
        matched = re.findall(regex, str(pdb_path))
        if matched:
            found_list.append(pdb_path)
    return found_list


def get_haddock_version(cmd: str) -> Optional[int]:
    """Validate if a path is a HADDOCK 2 or 3 installation."""
    cmd = cmd.split()[-1]
    haddock_script = Path(cmd)
    with open(haddock_script, "r") as f:
        for line in f:
            if "__HaddockVersion__" in line:
                try:
                    version = float(line.split()[2][1:])
                except ValueError:
                    raise InvalidParameter("HADDOCK version is not a float number")
                if version == 2.5:
                    return 2
                if version == 2.4:
                    return 2
    return None


def load_cns_params(run_cns_f):
    """Read a run.cns file and return all parameter keys."""
    param_regex = r"{===>}\s(\w*)=.*\;"
    param_l = []
    with open(run_cns_f, "r") as fh:
        for line in fh.readlines():
            try:
                parameter = re.findall(param_regex, line)[0]
            except IndexError:
                # this line does not contain a parameter
                continue
            param_l.append(parameter)
    return param_l


def edit_cns(run_cns, parameters):
    """Edit the run.cns parameter file."""
    param_regex = r"{===>}\s(\w*)=(.*)\;"
    edited_cns = run_cns.parent / "run.cns-edit"

    with open(edited_cns, "w") as edited_cns_fh:
        with open(run_cns, "r") as cns_fh:

            for line in cns_fh.readlines():
                if line.startswith("{===>}"):

                    param, value = re.findall(param_regex, line)[0]

                    if param in parameters:
                        custom_value = parameters[param]

                        # workaround to handle booleans
                        custom_value = str(custom_value).lower()

                        if custom_value != value:
                            log.debug(
                                f"Changing {param} from {value} to {custom_value}"
                            )

                            line = f"{{===>}} {param}={custom_value};" + os.linesep

                edited_cns_fh.write(line)

    edited_cns_fh.close()

    return edited_cns


def remove_cg(path_list):
    """Removes _cg.pdb from a list of Posix paths."""
    clean_list = copy.deepcopy(path_list)
    for path in clean_list:
        if "cg" in str(path):
            clean_list.remove(path)
    return clean_list

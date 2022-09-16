import logging
import shutil
import time
from pathlib import Path
from typing import Tuple

import toml

from benchmarktools.functions import get_haddock_version, load_cns_params
from benchmarktools.modules.errors import (
    ConfigKeyEmptyError,
    ConfigKeyUndefinedError,
    HeaderUndefined,
    InvalidParameter,
    InvalidRunName,
    ScenarioUndefined,
)

configlog = logging.getLogger("bmlog")


def read_config_file(config_file: str) -> Tuple[int, Path, str, str, list[str]]:
    """Check if a configuration file is valid."""
    concurrent: int
    dataset_path: Path
    receptor_suffix: str
    ligand_suffix: str

    conf = toml.load(config_file)

    if "general" not in conf:
        raise HeaderUndefined("general")

    # check the general path fields
    obrigatory_fields = ["dataset_path", "haddock_exec"]
    for field in obrigatory_fields:
        if field not in conf["general"]:
            raise ConfigKeyUndefinedError(field)
        elif not conf["general"][field]:
            raise ConfigKeyEmptyError(field)

    if "concurrent_jobs" in conf["general"]:
        try:
            concurrent = int(conf["general"]["concurrent_jobs"])
        except ValueError:
            raise InvalidParameter("concurrent_jobs must be an integer")
    else:
        concurrent = 10

    dataset_path = Path(conf["general"]["dataset_path"])

    # check the receptor/ligand suffix
    suffix_fields = ["receptor_suffix", "ligand_suffix"]
    for field in suffix_fields:
        if field not in conf["general"]:
            raise ConfigKeyUndefinedError(field)
        elif not conf["general"][field]:
            # its defined but its empty
            raise ConfigKeyEmptyError(field)

    receptor_suffix = conf["general"]["receptor_suffix"]
    ligand_suffix = conf["general"]["ligand_suffix"]

    receptor_suffix = receptor_suffix.replace(".pdb", "")
    ligand_suffix = ligand_suffix.replace(".pdb", "")

    # check if there are any scenarios
    scenario_name_list = [s for s in conf if "scenario" in s]
    if not scenario_name_list:
        raise ScenarioUndefined()

    return concurrent, dataset_path, receptor_suffix, ligand_suffix, scenario_name_list


#     pass


# class ConfigFile:
#     def __init__(self, toml_config_f):
#         conf = toml.load(toml_config_f)

#     def validate(self):
#         """Validate the fields in the config file."""

#         if "general" not in conf:
#             raise HeaderUndefined("general")

#         # check the general path fields
#         obrigatory_fields = ["dataset_path", "haddock_exec"]
#         for field in obrigatory_fields:
#             if field not in conf["general"]:
#                 raise ConfigKeyUndefinedError(field)
#             elif not conf["general"][field]:
#                 # its defined but its empty
#                 raise ConfigKeyEmptyError(field)

#         if "concurrent_jobs" in conf["general"]:
#             try:
#                 concurrent = int(conf["general"]["concurrent_jobs"])
#             except ValueError:
#                 raise InvalidParameter("concurrent_jobs must be an integer")
#         else:
#             concurrent = 10

#         dataset_path = Path(conf["general"]["dataset_path"])

#         # check the receptor/ligand suffix
#         suffix_fields = ["receptor_suffix", "ligand_suffix"]
#         for field in suffix_fields:
#             if field not in conf["general"]:
#                 raise ConfigKeyUndefinedError(field)
#             elif not conf["general"][field]:
#                 # its defined but its empty
#                 raise ConfigKeyEmptyError(field)

#         receptor_suffix = conf["general"]["receptor_suffix"]
#         ligand_suffix = conf["general"]["ligand_suffix"]

#         receptor_suffix = receptor_suffix.replace(".pdb", "")
#         ligand_suffix = ligand_suffix.replace(".pdb", "")

#         # check if there are any scenarios
#         scenario_name_list = [s for s in conf if "scenario" in s]
#         if not scenario_name_list:
#             raise ScenarioUndefined()
#         # TODO: Into the HaddockJob
#         # else:
#         #     scenarios = []
#         #     if haddock_version == 2:
#         #         run_cns_f = haddock_path / "protocols/run.cns-conf"
#         #         configlog.info(f"HADDOCK version 2 detected, lookig for {run_cns_f}")
#         #         cns_params = load_cns_params(run_cns_f)
#         #         run_name_l = []
#         #         for scenario_name in scenario_name_list:
#         #             scenarios.append(conf[scenario_name])
#         #             for param in conf[scenario_name]:
#         #                 if param == "run_name":
#         #                     name = conf[scenario_name][param]
#         #                     if name in run_name_l:
#         #                         raise InvalidRunName(name, message="duplicated")
#         #                     else:
#         #                         run_name_l.append(name)
#         #                 elif param == "ambig_tbl":
#         #                     # TODO: implement a tbl validator
#         #                     pass
#         #                 elif param not in cns_params:
#         #                     raise InvalidParameter(param)

#         # if not shutil.which("ssub"):
#         #     # this is specific for execution in the cluster
#         #     configlog.warning(
#         #         "ssub not in PATH, HADDOCK will fail if you are"
#         #         " running in the cluster!"
#         #     )
#         #     time.sleep(5)

#         return True

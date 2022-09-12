import logging
import shutil
import time
from pathlib import Path

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

configlog = logging.getLogger("setuplog")


class ConfigFile:
    def __init__(self, toml_config_f):
        self.conf = toml.load(toml_config_f)

    def validate(self):
        """Validate the fields in the config file."""

        if "general" not in self.conf:
            raise HeaderUndefined("general")

        # check the general path fields
        obrigatory_fields = ["dataset_path", "haddock_exec"]
        for field in obrigatory_fields:
            if field not in self.conf["general"]:
                raise ConfigKeyUndefinedError(field)
            elif not self.conf["general"][field]:
                # its defined but its empty
                raise ConfigKeyEmptyError(field)

        # check if the executable is executable
        self.executable = self.conf["general"]["haddock_exec"]
        self.haddock_version = get_haddock_version(self.executable)
        if self.haddock_version == 2:
            self.haddock_path = Path(
                self.conf["general"]["haddock_exec"].split()[-1]
            ).parent.parent
        else:
            self.haddock_path = ""
        # self.conf["general"]["haddock_path"] = haddock_path

        if "concurrent_jobs" in self.conf["general"]:
            try:
                self.concurrent = int(self.conf["general"]["concurrent_jobs"])
            except ValueError:
                raise InvalidParameter("concurrent_jobs must be an integer")
        else:
            self.concurrent = 10

        # self.haddock_path = Path(self.conf["general"]["haddock_path"])
        self.dataset_path = Path(self.conf["general"]["dataset_path"])

        # check the receptor/ligand suffix
        suffix_fields = ["receptor_suffix", "ligand_suffix"]
        for field in suffix_fields:
            if field not in self.conf["general"]:
                raise ConfigKeyUndefinedError(field)
            elif not self.conf["general"][field]:
                # its defined but its empty
                raise ConfigKeyEmptyError(field)

        self.receptor_suffix = self.conf["general"]["receptor_suffix"]
        self.ligand_suffix = self.conf["general"]["ligand_suffix"]

        self.receptor_suffix = self.receptor_suffix.replace(".pdb", "")
        self.ligand_suffix = self.ligand_suffix.replace(".pdb", "")

        # check if there are any scenarios
        scenario_name_list = [s for s in self.conf if "scenario" in s]
        if not scenario_name_list:
            raise ScenarioUndefined()
        else:
            self.scenarios = []
            if self.haddock_version == 2:
                run_cns_f = self.haddock_path / "protocols/run.cns-conf"
                configlog.info(f"HADDOCK version 2 detected, lookig for {run_cns_f}")
                cns_params = load_cns_params(run_cns_f)
                run_name_l = []
                for scenario_name in scenario_name_list:
                    self.scenarios.append(self.conf[scenario_name])
                    for param in self.conf[scenario_name]:
                        if param == "run_name":
                            name = self.conf[scenario_name][param]
                            if name in run_name_l:
                                raise InvalidRunName(name, message="duplicated")
                            else:
                                run_name_l.append(name)
                        elif param == "ambig_tbl":
                            # TODO: implement a tbl validator
                            pass
                        elif param not in cns_params:
                            raise InvalidParameter(param)

        if not shutil.which("ssub"):
            # this is specific for execution in the cluster
            configlog.warning(
                "ssub not in PATH, HADDOCK will fail if you are"
                " running in the cluster!"
            )
            time.sleep(5)

        return True

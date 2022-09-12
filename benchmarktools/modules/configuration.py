import logging
import pathlib
import re
import shutil
import time

import toml  # type: ignore

from benchmarktools.modules.errors import (
    ConfigKeyEmptyError,
    ConfigKeyUndefinedError,
    HeaderUndefined,
    InvalidParameter,
    InvalidRunName,
    PathNotFound,
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
        obrigatory_fields = ["dataset_path", "haddock_path", "python2"]
        for field in obrigatory_fields:
            if field not in self.conf["general"]:
                raise ConfigKeyUndefinedError(field)
            elif not self.conf["general"][field]:
                # its defined but its empty
                raise ConfigKeyEmptyError(field)
            else:
                obrigatory_path = pathlib.Path(self.conf["general"][field])
                if not obrigatory_path.exists():
                    raise PathNotFound(obrigatory_path)

        if "concurrent_jobs" in self.conf["general"]:
            try:
                self.concurrent = int(self.conf["general"]["concurrent_jobs"])
            except ValueError:
                raise InvalidParameter("concurrent_jobs must be an integer")
        else:
            self.concurrent = 10

        self.haddock_path = pathlib.Path(self.conf["general"]["haddock_path"])
        self.dataset_path = pathlib.Path(self.conf["general"]["dataset_path"])
        self.py2_path = pathlib.Path(self.conf["general"]["python2"])

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
            run_cns_f = self.haddock_path / "protocols/run.cns-conf"
            cns_params = self.load_cns_params(run_cns_f)
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

    @staticmethod
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

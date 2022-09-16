import logging
import os
import shutil
import sys
import time
from pathlib import Path
from typing import Optional, Tuple

import toml

from benchmarktools.functions import (
    get_haddock_version,
    glob_pdb_re,
    load_cns_params,
    remove_cg,
)
from benchmarktools.modules.dataset import Dataset
from benchmarktools.modules.errors import (
    ConfigKeyEmptyError,
    ConfigKeyUndefinedError,
    HeaderUndefined,
    InvalidParameter,
    InvalidRunName,
    ScenarioUndefined,
    SuffixError,
)
from benchmarktools.modules.haddock import HaddockJob

log = logging.getLogger("bmlog")


def check_fields(config_file: str) -> Optional[dict]:
    """Check if a configuration file is valid."""
    conf = toml.load(config_file)

    if "general" not in conf:
        raise HeaderUndefined("general")

    # check the general path fields
    obrigatory_fields = [
        "dataset_path",
        "haddock_exec",
        "receptor_suffix",
        "ligand_suffix",
    ]
    for field in obrigatory_fields:
        if field not in conf["general"]:
            raise ConfigKeyUndefinedError(field)
        elif not conf["general"][field]:
            raise ConfigKeyEmptyError(field)

    if "concurrent_jobs" in conf["general"]:
        try:
            int(conf["general"]["concurrent_jobs"])
        except ValueError:
            raise InvalidParameter("concurrent_jobs must be an integer")

    # check the receptor/ligand suffix
    suffix_fields = ["receptor_suffix", "ligand_suffix"]
    for field in suffix_fields:
        if field not in conf["general"]:
            raise ConfigKeyUndefinedError(field)
        elif not conf["general"][field]:
            raise ConfigKeyEmptyError(field)

    # check if there are any scenarios
    scenario_name_list = [s for s in conf if "scenario" in s]
    if not scenario_name_list:
        raise ScenarioUndefined()

    return conf


def parse_dataset_path(
    dataset_path: Path | str, receptor_suffix: str, ligand_suffix: str
) -> None:
    """Check if the files inside the dataset are formatted correctly."""
    # glob returns a abritrary ordered list
    target_directory_list = Path(dataset_path).glob("*")

    for target in target_directory_list:
        if os.path.isfile(target):
            # this is a file, dont check it
            continue

        # Demo of this regex using receptor_suffix = r_u
        # https://regex101.com/r/8MviAN/1
        receptor_regex = rf"(.*{receptor_suffix}_?)(\d.pdb|.pdb)"
        receptor_l = glob_pdb_re(target, receptor_regex)

        # FIXME: for CG mode to work we need to rewrite this part
        receptor_l = remove_cg(receptor_l)

        if len(receptor_l) == 0:
            raise SuffixError(target.name, receptor_suffix)
        elif not receptor_l:
            raise SuffixError(target.name, receptor_suffix)

        # Demo of this regex using ligand_suffix = l_u
        # https://regex101.com/r/8n3Bg3/1
        ligand_regex = rf"(.*{ligand_suffix}_?)(\d.pdb|.pdb)"
        ligand_l = glob_pdb_re(target, ligand_regex)

        # FIXME: for CG mode to work we need to rewrite this part
        ligand_l = remove_cg(ligand_l)

        if len(ligand_l) == 0:
            raise SuffixError(target.name, ligand_suffix)
        elif not ligand_l:
            raise SuffixError(target.name, ligand_suffix)


def initialize(config_file):
    """Initialize the Setup script and do the validations."""
    log.info("Initializing Setup")

    # Configuration
    try:
        log.info("Reading configuration file")
        (
            concurrent,
            dataset_path,
            receptor_suffix,
            ligand_suffix,
            scenario_name_list,
        ) = read_config_file(config_file)
    except Exception as e:
        log.error(e)
        sys.exit()
    log.info("Configuration file OK")

    # Dataset
    log.info("Loading Dataset")
    dataset = Dataset(dataset_path=dataset_path)

    try:
        log.info("Checking if receptor and ligands match suffix")
        dataset.check_input_files(
            receptor_suffix=receptor_suffix, ligand_suffix=ligand_suffix
        )
    except Exception as e:
        log.error(e)
        sys.exit()

    log.info("Input files are OK")

    # # Haddock
    # log.info("Initializing HADDOCK Wrapper")
    # if conf.haddock_version == 2:
    #     haddock = Haddock2Wrapper(haddock_cmd=conf.executable)
    # elif conf.haddock_version == 3:
    #     pass

    # try:
    #     log.info("Checking if HADDOCK is executable")
    #     haddock.check_if_executable()
    # except Exception as e:
    #     log.error(e)
    #     sys.exit()

    log.info("HADDOCK execution OK")

    # All checks ok!
    log.info("Initialization done!")
    return concurrent, receptor_suffix, ligand_suffix  # , dataset_obj, haddock


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

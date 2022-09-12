import logging
import sys

from benchmarktools.modules.configuration import ConfigFile
from benchmarktools.modules.dataset import Dataset
from benchmarktools.modules.haddock import Haddock2Wrapper

log = logging.getLogger("bmlog")


def initialize(config_file):
    """Initialize the Setup script and do the validations."""
    log.info("Initializing Setup")

    # Configuration
    log.info("Reading configuration file")
    conf = ConfigFile(config_file)

    try:
        log.info("Validating configuration file")
        conf.validate()
    except Exception as e:
        log.error(e)
        sys.exit()

    log.info("Configuration file OK")

    # Dataset
    log.info("Loading Dataset")
    dataset = Dataset(dataset_path=conf.dataset_path)

    try:
        log.info("Checking if receptor and ligands match suffix")
        dataset.check_input_files(
            receptor_suffix=conf.receptor_suffix, ligand_suffix=conf.ligand_suffix
        )
    except Exception as e:
        log.error(e)
        sys.exit()

    log.info("Input files are OK")

    # Haddock
    log.info("Initializing HADDOCK Wrapper")
    if conf.haddock_version == 2:
        haddock = Haddock2Wrapper(haddock_cmd=conf.executable)
    elif conf.haddock_version == 3:
        pass

    try:
        log.info("Checking if HADDOCK is executable")
        haddock.check_if_executable()
    except Exception as e:
        log.error(e)
        sys.exit()

    log.info("HADDOCK execution OK")

    # All checks ok!
    log.info("Initialization done!")
    return conf, dataset, haddock

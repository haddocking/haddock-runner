import logging
import sys

from benchmarktools.modules.configuration import ConfigFile
from benchmarktools.modules.dataset import Dataset
from benchmarktools.modules.haddock import HaddockWrapper

setuplog = logging.getLogger("setuplog")
ch = logging.StreamHandler()
formatter = logging.Formatter(
    " %(asctime)s %(module)s:%(lineno)d %(levelname)s - %(message)s"
)
ch.setFormatter(formatter)
setuplog.addHandler(ch)


def init(config_file):
    """Initialize the Setup script and do the validations."""
    setuplog.info("Initializing Setup")

    # Configuration
    setuplog.info("Reading configuration file")
    conf = ConfigFile(config_file)

    try:
        setuplog.info("Validating configuration file")
        conf.validate()
    except Exception as e:
        setuplog.error(e)
        sys.exit()

    setuplog.info("Configuration file OK")

    # Dataset
    setuplog.info("Loading Dataset")
    dataset = Dataset(dataset_path=conf.dataset_path)

    try:
        setuplog.info("Checking if receptor and ligands match suffix")
        dataset.check_input_files(
            receptor_suffix=conf.receptor_suffix, ligand_suffix=conf.ligand_suffix
        )
    except Exception as e:
        setuplog.error(e)
        sys.exit()

    setuplog.info("Input files are OK")

    # Haddock
    setuplog.info("Initializing HADDOCK Wrapper")
    haddock = HaddockWrapper(haddock_path=conf.haddock_path, py2=conf.py2_path)

    try:
        setuplog.info("Checking if HADDOCK is executable")
        haddock.check_if_executable()
    except Exception as e:
        setuplog.error(e)
        sys.exit()

    setuplog.info("HADDOCK execution OK")

    # All checks ok!
    setuplog.info("Initialization done!")
    return conf, dataset, haddock

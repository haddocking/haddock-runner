import argparse
import logging
import sys

from benchmarktools.modules.haddock import HaddockJob
from benchmarktools.modules.initialization import initialize
from benchmarktools.modules.queue import Queue
from benchmarktools.version import version

log = logging.getLogger("bmlog")
ch = logging.StreamHandler()
formatter = logging.Formatter(
    " %(asctime)s %(module)s:%(lineno)d %(levelname)s - %(message)s"
)
ch.setFormatter(formatter)
log.addHandler(ch)


def main(log_level="INFO"):

    ap = argparse.ArgumentParser(description="Run a HADDOCK Benchmark")
    ap.add_argument("config_file", help="Configuration file, toml format")
    ap.add_argument(
        "--force",
        dest="force",
        action="store_true",
        default=False,
        help="DEV only, forcefully removeinitiated runs",
    )

    ap.add_argument(
        "-v",
        "--version",
        help="show version",
        action="version",
        version=f"Running {ap.prog} v{version}",
    )
    args = ap.parse_args()

    log.setLevel(log_level)
    log.info("###############################################")
    log.info("")
    log.info(f"      Welcome to benchmark-tools v{version}")
    log.info("")
    log.info("###############################################")
    log.warning(
        "If this is not running in the background, your "
        "benchmarking will stop when you close the terminal!"
    )
    log.warning(
        f"To run it in the background, run with: "
        f'"nohup python {" ".join(sys.argv)} &"'
    )
    # time.sleep(2)

    config, dataset, haddock = initialize(args.config_file)

    prepared_run_l = []
    for i, run_scenario in enumerate(config.scenarios):
        log.info(f"Setting up Scenario {i+1}")
        log.info(run_scenario)
        ready_runs = dataset.setup(
            haddock=haddock,
            parameters=run_scenario,
            receptor_suffix=config.receptor_suffix,
            ligand_suffix=config.ligand_suffix,
            force=args.force,
        )
        prepared_run_l.extend(ready_runs)

    log.info("Generating Job list")
    job_list = []
    for target in prepared_run_l:
        job = HaddockJob(
            haddock_cmd=config.executable,
            job_path=target,
        )
        job_list.append(job)

    # Execute!
    total_jobs = len(job_list)
    log.info(f"Executing jobs n={total_jobs}")

    queue = Queue(job_list, concurrent=config.concurrent)
    queue.execute()

    log.info("Done!")


if __name__ == "__main__":
    sys.exit(main())

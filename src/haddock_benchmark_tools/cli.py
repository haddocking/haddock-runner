import argparse
import logging
import sys
import time

from haddock_benchmark_tools.modules.haddock import HaddockJob
from haddock_benchmark_tools.modules.initialization import init
from haddock_benchmark_tools.modules.queue import Queue
from haddock_benchmark_tools.version import version

setuplog = logging.getLogger("setuplog")
ch = logging.StreamHandler()
formatter = logging.Formatter(
    " %(asctime)s %(module)s:%(lineno)d %(levelname)s - %(message)s"
)
ch.setFormatter(formatter)
setuplog.addHandler(ch)


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

    setuplog.setLevel(log_level)
    setuplog.info("###############################################")
    setuplog.info("")
    setuplog.info(f"      Welcome to benchmark-tools v{version}")
    setuplog.info("")
    setuplog.info("###############################################")
    setuplog.warning(
        "If this is not running in the background, your "
        "benchmarking will stop when you close the terminal!"
    )
    setuplog.warning(
        f"To run it in the background, run with: "
        f'"nohup python {" ".join(sys.argv)} &"'
    )
    time.sleep(2)

    config, dataset, haddock = init(args.config_file)

    prepared_run_l = []
    for i, run_scenario in enumerate(config.scenarios):
        setuplog.info(f"Setting up Scenario {i+1}")
        setuplog.info(run_scenario)
        ready_runs = dataset.setup(
            haddock=haddock,
            parameters=run_scenario,
            receptor_suffix=config.receptor_suffix,
            ligand_suffix=config.ligand_suffix,
            force=args.force,
        )
        prepared_run_l.extend(ready_runs)

    setuplog.info("Generating Job list")
    job_list = []
    for target in prepared_run_l:
        job = HaddockJob(
            haddock_path=config.haddock_path, py2=config.py2_path, run_path=target
        )
        job_list.append(job)

    # Execute!
    total_jobs = len(job_list)
    setuplog.info(f"Executing jobs n={total_jobs}")

    queue = Queue(job_list, concurrent=config.concurrent)
    queue.execute()

    setuplog.info("Done!")


if __name__ == "__main__":
    sys.exit(main())

from setuptools import find_packages, setup

from src.haddock_benchmark_tools.version import version

with open("requirements.txt") as f:
    required = f.read().splitlines()

setup(
    name="haddock-benchmark-tools",
    license="Apache License 2.0",
    version=version,
    author="BonvinLab",
    description="benchmarking framework for HADDOCK v2.4+",
    author_email="software.csb@gmail.com",
    packages=find_packages("src"),
    package_dir={"": "src"},
    classifiers=[],
    python_requires=">=3.6, <4",
    install_requires=required,
    entry_points={
        "console_scripts": [
            "haddock_bm=haddock_benchmark_tools.cli:maincli",
        ],
    },
)

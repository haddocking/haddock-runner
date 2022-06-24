from setuptools import find_packages, setup

with open("requirements.txt") as f:
    required = f.read().splitlines()

setup(
    name="haddock-benchmark-tools",
    license="Apache License 2.0",
    version="0.3.0",
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

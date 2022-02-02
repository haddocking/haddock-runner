from setuptools import setup, find_packages

setup(
    name="haddock-benchmark-tools",
    license="Apache License 2.0",
    version="0.2.2",
    author="HADDOCK",
    description="benchmarking framework for HADDOCK v2.4",
    author_email="A.M.J.J.Bonvin@uu.nl",
    packages=find_packages("src"),
    package_dir={"": "src"},
    classifiers=[],
    python_requires=">=3.6, <4",
    install_requires=["toml"],
    entry_points={
        "console_scripts": [
            "haddock_bm=haddock_benchmark_tools.cli:maincli",
        ],
    },
)

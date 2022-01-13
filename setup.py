from setuptools import setup, find_packages

setup(
    name="haddock-benchmark-tools",
    license="Apache License 2.0",
    version="0.1.0",
    author="HADDOCK",
    description="benchmarking framework for HADDOCK v2.4",
    author_email="A.M.J.J.Bonvin@uu.nl",
    packages=find_packages("src"),
    package_dir={"": "src"},
    classifiers=[
        # complete classifier list:
        # http://pypi.python.org/pypi?%3Aaction=list_classifiers
        "Development Status :: 4 - Beta",
        "License :: OSI Approved :: Apache Software License",
        "Natural Language :: English",
        "Operating System :: POSIX",
        "Operating System :: POSIX :: Linux",
        "Operating System :: MacOS",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3 :: Only",
    ],
    python_requires=">=3.6, <4",
    install_requires=["toml"],
    entry_points={
        "console_scripts": [
            "haddock_bm=haddock_benchmark_tools:maincli",
        ],
    },
    setup_requires=["pytest-runner"],
    tests_require=["pytest"],
)

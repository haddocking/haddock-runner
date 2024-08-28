#========================================================================================================#
FROM golang:1.21 AS golang

#========================================================================================================#
FROM xenonmiddleware/slurm:20

#==============================================================================================
# Define ARGs
ARG HADDOCK_VERSION=v3.0.0-beta.5
ARG USERNAME=dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID

#==============================================================================================
# Install system dependencies
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  build-essential=12.1ubuntu2 \
  sudo=1.8.16-0ubuntu1.10 \
  wget=1.17.1-1ubuntu1.5 \
  git=1:2.7.4-0ubuntu1.10 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

#==============================================================================================
# Configure User
RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME

#==============================================================================================
# Install miniconda
ENV CONDA_DIR /opt/conda
RUN wget --quiet --no-check-certificate https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O ~/miniconda.sh  \
  && /bin/bash ~/miniconda.sh -b -p /opt/conda

ENV PATH=$CONDA_DIR/bin:$PATH
RUN conda install python=3.9

#==============================================================================================
# Install HADDOCK3
WORKDIR /opt

RUN git clone --recursive https://github.com/haddocking/haddock3.git
WORKDIR /opt/haddock3
RUN git checkout ${HADDOCK_VERSION}

WORKDIR /opt/haddock3/src/fcc/src
RUN make

WORKDIR /opt/haddock3
RUN pip install --no-cache-dir -r requirements.txt \
  && python setup.py develop

WORKDIR /opt/haddock3/bin
COPY cns /opt/haddock3/bin/cns

#==============================================================================================
# Copy Go
COPY --from=golang /usr/local/go/ /usr/local/go/

#========================================================================================================#
# Configure container startup
WORKDIR /app
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

USER $USERNAME
WORKDIR $HOME
#========================================================================================================#

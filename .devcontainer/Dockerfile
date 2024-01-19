#========================================================================================================#
FROM golang:1.21 AS golang

#========================================================================================================#
FROM xenonmiddleware/slurm:20

ARG USERNAME=dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
  && apt-get update \
  && apt-get install -y --no-install-recommends sudo=1.8.16-0ubuntu1.10 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* \
  && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
  && chmod 0440 /etc/sudoers.d/$USERNAME

WORKDIR /app
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

COPY --from=golang /usr/local/go/ /usr/local/go/

#========================================================================================================#

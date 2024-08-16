FROM debian:12-slim AS dev-base

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Make apt family commands to not ask any user input
ENV DEBIAN_FRONTEND=noninteractive

# Install dev dependecies.
RUN echo \
  && echo "update system packages manager cache" \
  && apt-get update \
  \
  && echo "installing base packages" \
  && apt-get install --no-install-recommends -y \
    bash-completion \
    ca-certificates \
    curl \
    file \
    git \
    locales \
    make \
    nano \
    openssh-client \
    procps \
    rsync \
    sshpass \
    sudo \
    unzip \
    wget \
  \
  && echo "installing X11 packages" \
  && apt-get install --no-install-recommends -y \
    fontconfig \
    libfreetype6 \
    libxi6 \
    libxrender1 \
    libxtst6 \
    xauth \
  \
  && echo "clean system packages manager cache" \
  && apt-get clean -y \
  \
  && echo "configuring base image" \
  && echo "root:root" | chpasswd


FROM dev-base AS dev-base-with-sshd

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install openssh server.
RUN echo \
  && echo "update system packages manager cache" \
  && apt-get update \
  \
  && echo "installing SSH server packages" \
  && apt-get install --no-install-recommends -y \
    openssh-server \
  \
  && echo "clean system packages manager cache" \
  && apt-get clean -y \
  \
  && echo "configuring SSH server packages" \
  && mkdir -p /var/run/sshd \
  && sed -i 's/^#\(PermitRootLogin\) .*/\1 yes/' /etc/ssh/sshd_config \
  && sed -i 's/^\(UsePAM yes\)/# \1/' /etc/ssh/sshd_config \
  && service ssh start


FROM dev-base-with-sshd AS dev-base-with-user

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Create an user.
RUN echo \
  && echo "creating user" \
  && useradd -rm -d /home/user -s /bin/bash -g root -G sudo -u 1000 user \
  && echo "user:user" | chpasswd \
  \
  && echo "configuring sudo for user" \
  && mkdir -p /etc/sudoers.d \
  && echo "user ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/user \
  && chmod 0440 /etc/sudoers.d/user \
  \
  && echo "setup bashrc.d" \
  && mkdir -p /home/user/.bashrc.d \
  && echo 'for f in /home/user/.bashrc.d/* ; do . "${f}" ; done' >> /home/user/.bashrc

COPY 00_sensible.bash /home/user/.bashrc.d
COPY 01_x11.bash /home/user/.bashrc.d


FROM dev-base-with-user AS dev-base-with-dev-packags

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN echo \
  && echo "dev-DataGrip-debian12-slim" >> /etc/debian_chroot


FROM dev-base-with-dev-packags AS dev-base-with-ide

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install Datagrip.
ARG DATAGRIP_VERSION
RUN echo \
  && echo "installing DataGrip IDE" \
  && mkdir -p /opt/jetbrains/DataGrip \
  && curl --progress-bar -sSLo \
    /opt/jetbrains/DataGrip.tar.gz \
    https://download-cdn.jetbrains.com/datagrip/datagrip-${DATAGRIP_VERSION}.tar.gz \
  && tar -xzvf \
    /opt/jetbrains/DataGrip.tar.gz\
    --strip-components=1 \
    -C /opt/jetbrains/DataGrip \
  && ln -s \
    /opt/jetbrains/DataGrip/bin/datagrip.sh \
    /usr/bin/ide.sh \
  && rm /opt/jetbrains/DataGrip.tar.gz


FROM dev-base-with-ide AS dev-base-final

# Make commands to fail due to an error at any stage in the pipe,
# prepend set -o pipefail && to ensure that an unexpected error prevents
# the build from inadvertently succeeding.
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.source="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.authors="Stanislav Lebedev <kugui@yandex.ru>"
LABEL org.opencontainers.image.url="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.vendor="dev-containers"
LABEL org.opencontainers.image.licenses="BSD"
LABEL org.opencontainers.image.description="DataGrip dev container"
LABEL org.opencontainers.image.title="DataGrip dev container"

LABEL dev.containers.compilers=""
LABEL dev.containers.distro.variant="slim"
LABEL dev.containers.distro.version="debian12"
LABEL dev.containers.ide="DataGrip"

# Run SSH server
ENV DEV_CONTAINER_SSH_PORT=2221
CMD ["bash", "-c", "/usr/sbin/sshd -De -p${DEV_CONTAINER_SSH_PORT}"]

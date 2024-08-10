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
    git \
    file \
    locales \
    nano \
    make \
    openssh-client \
    procps \
    rsync \
    sshpass \
    sudo \
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
  && echo "setup environment variables for user" \
  && echo "DISPLAY=:0" >> /home/user/.bashrc \
  && echo "LIBGL_ALWAYS_INDIRECT=1" >> /home/user/.bashrc


FROM dev-base-with-user AS dev-base-with-dev-packags

# Applications versions.
ARG GOLANG_VERSION
ARG GOFUMPT_VERSION
ARG BAZELISK_VERSION
ARG BUILDIFIER_VERSION
ARG BUILDOZER_VERSION

RUN echo \
  && echo "installing go" \
  && mkdir -p /opt \
  && curl -sL "https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" | tar -C /opt -xzvf - \
  \
  && echo "installing gofumpt" \
  && curl -sLo \
    /usr/bin/gofumpt \
    "https://github.com/mvdan/gofumpt/releases/download/${GOFUMPT_VERSION}/gofumpt_${GOFUMPT_VERSION}_linux_amd64" \
  && chmod +x /usr/bin/gofumpt \
  \
  && echo "installing bazelisk" \
  && curl -sLo \
    /usr/bin/bazel \
    "https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64" \
  && chmod +x /usr/bin/bazel \
  \
  && echo "installing buildifier" \
  && curl -sLo \
    /usr/bin/buildifier \
    "https://github.com/bazelbuild/buildtools/releases/download/${BUILDIFIER_VERSION}/buildifier-linux-amd64" \
  && chmod +x /usr/bin/buildifier \
  \
  && echo "installing buildozer" \
  && curl -sLo \
    /usr/bin/buildozer \
    "https://github.com/bazelbuild/buildtools/releases/download/${BUILDOZER_VERSION}/buildozer-linux-amd64" \
  && chmod +x /usr/bin/buildozer \
  \
  && echo "dev-ideaIU-debian12-slim" >> /etc/debian_chroot


FROM dev-base-with-dev-packags AS dev-base-with-ide

# Install IDEA Intellij Ultimate.
ARG IDEA_IU_VERSION
RUN echo \
  && echo "installing ideaIU IDE" \
  && mkdir -p /opt/jetbrains/ideaIU \
  && curl --progress-bar -sSLo \
    /opt/jetbrains/ideaIU.tar.gz \
    https://download-cdn.jetbrains.com/idea/ideaIU-${IDEA_IU_VERSION}.tar.gz \
  && tar -xzvf \
    /opt/jetbrains/ideaIU.tar.gz\
    --strip-components=1 \
    -C /opt/jetbrains/ideaIU \
  && ln -s \
    /opt/jetbrains/ideaIU/bin/idea.sh \
    /usr/bin/ide.sh \
  && rm -f /opt/jetbrains/ideaIU.tar.gz


FROM dev-base-with-ide AS dev-base-final

LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.source="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.authors="Stanislav Lebedev <kugui@yandex.ru>"
LABEL org.opencontainers.image.url="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.vendor="dev-containers"
LABEL org.opencontainers.image.licenses="BSD"
LABEL org.opencontainers.image.description="IDEA Intellij Ultimate dev container with Go compiler"
LABEL org.opencontainers.image.title="IDEA Intellij Ultimate dev container with Go compiler"

LABEL dev.containers.compilers="go"
LABEL dev.containers.distro.variant="slim"
LABEL dev.containers.distro.version="debian12"
LABEL dev.containers.ide="ideaIU"

# Run SSH server
ENV DEV_CONTAINER_SSH_PORT=2221
CMD ["bash", "-c", "/usr/sbin/sshd -De -p$DEV_CONTAINER_SSH_PORT{}"]

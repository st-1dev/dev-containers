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
ARG CMAKE_VERSION
ARG NINJA_VERSION
ARG BAZELISK_VERSION
ARG BUILDIFIER_VERSION
ARG BUILDOZER_VERSION

RUN echo \
  && echo "update system packages manager cache" \
  && apt-get update \
  \
  && echo "installing C/C++/Python packages" \
  && apt-get install --no-install-recommends -y \
    clang-15 \
    clang-format-15 \
    clang-tidy-15 \
    lld-15 \
    lldb-15 \
    python3 \
  \
  && echo "clean system packages manager cache" \
  && apt-get clean -y \
  \
  && echo "installing cmake" \
  && curl -sLo \
    /opt/cmake.tar.gz \
    "https://github.com/Kitware/CMake/releases/download/v${CMAKE_VERSION}/cmake-${CMAKE_VERSION}-linux-x86_64.tar.gz" \
  && mkdir -p /opt/cmake \
  && tar --strip-components=1  -xf /opt/cmake.tar.gz -C /opt/cmake \
  && rm /opt/cmake.tar.gz \
  \
  && echo "installing ninja" \
  && curl -sLo \
    /opt/ninja.zip \
   "https://github.com/ninja-build/ninja/releases/download/v${NINJA_VERSION}/ninja-linux.zip" \
  && mkdir -p /opt/ninja/bin \
  && unzip -x /opt/ninja.zip -d /opt/ninja/bin \
  && rm /opt/ninja.zip \
  \
  && echo "installing bazelisk" \
  && curl -sLo \
    /usr/bin/bazel \
    "https://github.com/bazelbuild/bazelisk/releases/download/v${BAZELISK_VERSION}/bazelisk-linux-amd64" \
  && chmod +x /usr/bin/bazel \
  \
  && echo "installing buildifier" \
  && curl -sLo \
    /usr/bin/buildifier \
    "https://github.com/bazelbuild/buildtools/releases/download/v${BUILDIFIER_VERSION}/buildifier-linux-amd64" \
  && chmod +x /usr/bin/buildifier \
  \
  && echo "installing buildozer" \
  && curl -sLo \
    /usr/bin/buildozer \
    "https://github.com/bazelbuild/buildtools/releases/download/v${BUILDOZER_VERSION}/buildozer-linux-amd64" \
  && chmod +x /usr/bin/buildozer \
  \
  && echo "dev-CLion-debian12-slim-clang13" >> /etc/debian_chroot


FROM dev-base-with-dev-packags AS dev-base-with-ide

# Install CLion.
ARG CLION_VERSION
RUN echo \
  && echo "installing CLion IDE" \
  && mkdir -p /opt/jetbrains/CLion \
  && curl --progress-bar -sSLo \
    /opt/jetbrains/CLion.tar.gz \
    https://download-cdn.jetbrains.com/cpp/CLion-${CLION_VERSION}.tar.gz \
  && tar -xzvf \
    /opt/jetbrains/CLion.tar.gz\
    --strip-components=1 \
    -C /opt/jetbrains/CLion \
  && ln -s \
    /opt/jetbrains/CLion/bin/clion.sh \
    /usr/bin/ide.sh \
  && rm /opt/jetbrains/CLion.tar.gz


FROM dev-base-with-ide AS dev-base-final

LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.source="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.authors="Stanislav Lebedev <kugui@yandex.ru>"
LABEL org.opencontainers.image.url="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.documentation="https://github.com/stle85/dev-containers"
LABEL org.opencontainers.image.vendor="dev-containers"
LABEL org.opencontainers.image.licenses="BSD"
LABEL org.opencontainers.image.description="CLion dev container with clang compiler"
LABEL org.opencontainers.image.title="CLion dev container with clang compiler"

LABEL dev.containers.compilers="clang15,python3"
LABEL dev.containers.distro.variant="slim"
LABEL dev.containers.distro.version="debian12"
LABEL dev.containers.ide="CLion"

# Run SSH server
ENV DEV_CONTAINER_SSH_PORT=2221
CMD ["bash", "-c", "/usr/sbin/sshd -De -p${DEV_CONTAINER_SSH_PORT}"]

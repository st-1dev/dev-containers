# ----------------------------------------------------------------------------------

CLION_VERSION=2024.1.4
DATAGRIP_VERSION=2024.1.4
IDEA_IU_VERSION=2024.1.4
PYCHARM_COMMUNITY_VERSION=2024.1.4
PYCHARM_PROFESSIONAL_VERSION=2024.1.4

# ----------------------------------------------------------------------------------

GOLANG_VERSION=1.22.6
GOFUMPT_VERSION=v0.6.0
BAZELISK_VERSION=v1.20.0
BUILDIFIER_VERSION=v7.1.2
BUILDOZER_VERSION=v7.1.2

# ----------------------------------------------------------------------------------
# Debian 12

# CLION

build_dev_clion_debina12-slim_dummy:
	docker build --rm \
	  --tag "clion:${CLION_VERSION}-debian12-slim-dummy" \
	  --file ./dockerfiles/dev_CLion_debian-12-slim_dummy.Dockerfile \
	  --label "dev.containers.version=${CLION_VERSION}" \
	  --build-arg "CLION_VERSION=${CLION_VERSION}" \
	  --build-arg "BAZELISK_VERSION=${BAZELISK_VERSION}" \
	  --build-arg "BUILDIFIER_VERSION=${BUILDIFIER_VERSION}" \
	  --build-arg "BUILDOZER_VERSION=${BUILDOZER_VERSION}" \
	  ./build

build_dev_clion_debina12-slim_clang13:
	docker build --rm \
	  --tag "clion:${CLION_VERSION}-debian12-slim-clang13" \
	  --file ./dockerfiles/dev_CLion_debian-12-slim_clang-13.Dockerfile \
	  --label "dev.containers.version=${CLION_VERSION}" \
	  --build-arg "CLION_VERSION=${CLION_VERSION}" \
	  --build-arg "BAZELISK_VERSION=${BAZELISK_VERSION}" \
	  --build-arg "BUILDIFIER_VERSION=${BUILDIFIER_VERSION}" \
	  --build-arg "BUILDOZER_VERSION=${BUILDOZER_VERSION}" \
	  ./build

# DATAGRIP

build_dev_datagrip_debina12-slim:
	docker build --rm \
	  --tag "datagrip:${DATAGRIP_VERSION}-debian12-slim" \
	  --file ./dockerfiles/dev_datagrip_debian-12-slim.Dockerfile \
	  --label "dev.containers.version=${DATAGRIP_VERSION}" \
	  --build-arg "DATAGRIP_VERSION=${DATAGRIP_VERSION}" \
	  ./build

# IDEA INTELIJ ULTIMATE

build_dev_idea-iu_debina12-slim:
	docker build --rm \
	  --tag "idea-iu:${IDEA_IU_VERSION}-debian12-slim" \
	  --file ./dockerfiles/dev_ideaIU_debian-12-slim.Dockerfile \
	  --label "dev.containers.version=${IDEA_IU_VERSION}" \
	  --build-arg "IDEA_IU_VERSION=${IDEA_IU_VERSION}" \
	  --build-arg "GOLANG_VERSION=${GOLANG_VERSION}" \
	  --build-arg "GOFUMPT_VERSION=${GOFUMPT_VERSION}" \
	  --build-arg "BAZELISK_VERSION=${BAZELISK_VERSION}" \
	  --build-arg "BUILDIFIER_VERSION=${BUILDIFIER_VERSION}" \
	  --build-arg "BUILDOZER_VERSION=${BUILDOZER_VERSION}" \
	  ./build

# PYCHARM COMMUNITY EDITION

build_dev_pycharm-community_debina12-slim:
	docker build --rm \
	  --tag "pycharm-community:${PYCHARM_COMMUNITY_VERSION}-debian12-slim" \
	  --file ./dockerfiles/dev_pycharm-community_debian-12-slim.Dockerfile \
	  --label "dev.containers.version=${PYCHARM_COMMUNITY_VERSION}" \
	  --build-arg "PYCHARM_COMMUNITY_VERSION=${PYCHARM_COMMUNITY_VERSION}" \
	  ./build

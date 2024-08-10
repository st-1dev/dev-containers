.PHONY:
default: build

# ----------------------------------------------------------------------------------

# Go compiler path and flags
GO?=go
GOOS?=linux
GOARCH?=amd64

CGO_ENABLED=0
GO_FLAGS=-tags=containers_image_openpgp

# Code formatting tool
GOFUMPT=gofumpt

# ----------------------------------------------------------------------------------

# App
APP=dev-runner

# Build dir
BUILD_DIR=out

# Source code
GO_SRC=$(shell find . -type f -name '*.go')

# ----------------------------------------------------------------------------------

# build apps
build: "$(BUILD_DIR)" "$(BUILD_DIR)/$(APP)_$(GOOS)"

"$(BUILD_DIR)":
	mkdir -p "$(BUILD_DIR)"
	touch "$(BUILD_DIR)/.gitkeep"

"$(BUILD_DIR)/$(APP)_$(GOOS)": $(GO_SRC)
	env "GOOS=$(GOOS)" "GOARCH=$(GOARCH)" "CGO_ENABLED=$(CGO_ENABLED)" \
	"$(GO)" build "$(GO_FLAGS)" -o "$(BUILD_DIR)/$(APP)_$(GOOS)" "./cmd/$(APP)"

# ----------------------------------------------------------------------------------

# clean build cache
clean:
	rm -fr "$(BUILD_DIR)"
	mkdir -p "$(BUILD_DIR)"
	touch "$(BUILD_DIR)/.gitkeep"

# ----------------------------------------------------------------------------------

# format code
format:
	"$(GOFUMPT)" -l -w ./

# update dependency
tidy:
	"$(GO)" mod tidy

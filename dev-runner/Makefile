.PHONY:
default: build

GO=go
GO_FLAGS=
GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

APP=dev-runner
GO_SRC=$(shell find . -type f -name '*.go')

build: $(APP)

$(APP): $(GO_SRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GO_FLAGS) -o $(APP) ./cmd/$(APP)

format:
	gofumpt -l -w ./

clean:
	rm -f $(APP)

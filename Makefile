.DEFAULT_GOAL := all
VERSION := $(shell git describe --tags --abbrev=0)

# https://www.forkingbytes.com/blog/dynamic-versioning-your-go-application/
GOOPTS := "-ldflags=-s -w -X main.Version=$(VERSION)"
MAKE := $(MAKE) --no-print-directory --jobs

# Not used here, but this is fascinating:
# https://stackoverflow.com/a/12110773/11133327
OPERATING_SYSTEMS = darwin linux windows freebsd openbsd
$(OPERATING_SYSTEMS):; -@mkdir -p ./bin
	GOARCH=$(ARCH) GOOS=$(@) go build $(GOOPTS) -o ./bin/$(CMD)-$(@)-$(ARCH) ./cmd/$(CMD)

ARCHITECTURES = amd64 arm64
$(ARCHITECTURES): ; @CMD=$(CMD) ARCH=$(@) $(MAKE) $(OPERATING_SYSTEMS)

CMDS = alw-api alw-cli
$(CMDS):; @CMD=$(@) $(MAKE) $(ARCHITECTURES)

lint:; @golangci-lint run

test:
	@go vet ./...
	@go test ./... -covermode=count -coverprofile=c.out
	@go tool cover -func=c.out

all: $(CMDS) test lint

clean:; @rm -rfv ./bin

run:
	@go build $(GOOPTS) -o ./bin/alw-api ./cmd/alw-api
	@./bin/alw-api -p 80

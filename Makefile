VERSION = $(shell git describe --tags --abbrev=0)

# https://www.forkingbytes.com/blog/dynamic-versioning-your-go-application/
FLAGS = "-ldflags=-s -w -X main.Version=$(VERSION)"

# Not used here, but this is fascinating:
# https://stackoverflow.com/a/12110773/11133327
OPERATING_SYSTEMS = darwin linux windows freebsd openbsd
$(OPERATING_SYSTEMS):
	GOARCH=$(ARCH) GOOS=$(@) go build $(FLAGS) -o ./$(CMD)-$(@)-$(ARCH) ./cmd/$(CMD)

ARCHITECTURES = amd64 arm64
$(ARCHITECTURES): ; @CMD=$(CMD) ARCH=$(@) $(MAKE) $(OPERATING_SYSTEMS)

CMDS = alw-api alw-cli
$(CMDS): ; @CMD=$(@) $(MAKE) -j $(ARCHITECTURES)

all: $(CMDS)

clean: ; @rm -fv ./alw-*-*-*

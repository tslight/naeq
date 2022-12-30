# Not used here, but this is fascinating:
# https://stackoverflow.com/a/12110773/11133327

FLAGS = '-ldflags=-s -w'
VERSION = $(shell git describe --tags --abbrev=0)

OPERATING_SYSTEMS = darwin linux windows freebsd openbsd
$(OPERATING_SYSTEMS):
	GOARCH=$(ARCH) GOOS=$(@) go build $(FLAGS) -o ./$(CMD)-$(@)-$(ARCH) ./cmd/$(CMD)

ARCHITECTURES = amd64 arm64
$(ARCHITECTURES): ; @CMD=$(CMD) ARCH=$(@) $(MAKE) $(OPERATING_SYSTEMS)

CMDS = alw-api alw-cli
$(CMDS): ; @CMD=$(@) $(MAKE) -j $(ARCHITECTURES)

all: $(CMDS)

clean: ; @rm -fv ./alw-*-*-*

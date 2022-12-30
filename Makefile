FLAGS = '-ldflags=-s -w'
VERSION = $(shell git describe --tags --abbrev=0)

OPERATING_SYSTEMS = darwin linux windows freebsd openbsd
$(OPERATING_SYSTEMS):
	GOARCH=$(ARCH) GOOS=$(@) go build $(FLAGS) -o ./alw-$(TYPE)-$(@)-$(ARCH) ./cmd/alw_$(TYPE)

ARCHITECTURES = amd64 arm64
$(ARCHITECTURES):
	@TYPE=$(TYPE) ARCH=$(@) $(MAKE) $(OPERATING_SYSTEMS)

TYPES = api cli
$(TYPES):
	@TYPE=$(@) $(MAKE) $(ARCHITECTURES)

all: $(TYPES)

clean:
	@rm -fv ./alw-*-*-*

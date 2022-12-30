ARCH ?= amd64
FLAGS = '-ldflags=-s -w'
VERSION = $(shell git describe --tags --abbrev=0)

OPERATING_SYSTEMS = darwin linux windows
$(OPERATING_SYSTEMS):
	GOARCH=${ARCH} GOOS=$(@) go build ${FLAGS} -o ./alw-${TYPE}-$(@)-${ARCH} ./cmd/alw_${TYPE}

ARCHITECTURES = amd64 arm64
$(ARCHITECTURES):
	@TYPE=api ARCH=$(@) make $(OPERATING_SYSTEMS)
	@TYPE=cli ARCH=$(@) make $(OPERATING_SYSTEMS)

all: $(ARCHITECTURES)

clean:
	@rm -fv ./alw-*-*-*
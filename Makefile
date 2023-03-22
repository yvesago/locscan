
VERSION="1.0.0-pub"
#VERSION=$(shell git describe --abbrev=0 --tags)

#BUILD=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%FT%T%z)

# Binaries to be build
PLATFORMS = linux/locscan windows/locscan.exe darwin/locscan-app
BINS = $(wildcard build/*/*)

# functions
temp = $(subst /, ,$@)
os = $(word 1, $(temp))

# Setup the -ldflags option for go building, interpolate the variable values
#LDFLAGS=-trimpath -ldflags "-w -s -X 'main.Version=${VERSION}, git: ${BUILD}, build: ${DATE}'"
LDFLAGS=-trimpath -ldflags "-w -s -X 'main.Version=${VERSION}, build: ${DATE}'"

# Build binaries
#  first build : linux/locscan
$(PLATFORMS):
	@mkdir -p build/${os}
	CGO_ENABLED=0 GOOS=${os} go build ${LDFLAGS} -o build/$@
	@echo " => bin builded: build/$@"

build: $(PLATFORMS)

# List binaries
$(BINS):
	@sha256sum $@ 

sha: $(BINS)

# Cleans our project: deletes binaries
clean:
	rm -rf build/
	@echo "Build cleaned"


all: build

.PHONY: clean build sha $(BINS) $(PLATFORMS)


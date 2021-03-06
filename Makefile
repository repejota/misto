BINARY=misto
VERSION=`cat VERSION`
BUILD=`git symbolic-ref HEAD 2> /dev/null | cut -b 12-`-`git log --pretty=format:%h -1`
PACKAGES = $(shell go list ./...)

# Setup the -ldflags option for go build here, interpolate the variable
# values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Build & Install

install:	## Build and install package on your system
	go install $(LDFLAGS) -v $(PACKAGES)

.PHONY: version
version:	## Show version information
	@echo $(VERSION)-$(BUILD)

# Testing

.PHONY: test
test:	## Execute package tests 
	go test -v $(PACKAGES)

.PHONY: cover-profile
cover-profile:	## Compile tests coverage data
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	rm -rf coverage.out

.PHONY: cover
cover: cover-profile	## Generate test coverage data
	go tool cover -func=coverage-all.out

.PHONY: cover-html
cover-html: cover-profile	## Generate coverage report
	go tool cover -html=coverage-all.out

.PHONY: coveralls
coveralls:	## Send coverage report to https://coveralls.io/github/repejota/misto
	goveralls -service circle-ci -repotoken 9EmpV6j34d3itxKKXJCjTYicQPZhgzwj3

# Lint

lint:	## Lint source code
	gometalinter --disable-all --enable=errcheck --enable=vet --enable=vetshadow

# Dependencies

deps:	## Install package dependencies
	go get -v -t -d -u github.com/sirupsen/logrus
	go get -v -t -d -u github.com/docker/docker/client
	go get -v -t -d -u github.com/fatih/color
	go get -v -t -d -u github.com/spf13/cobra
	go get -v -t -d -u github.com/repejota/cscanner

dev-deps:	## Install devellpment dependencies
	go get -v -t -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -v -t -u github.com/mattn/goveralls

# Cleaning up

.PHONY: clean
clean:	## Delete generated development environment
	go clean
	rm -rf ${BINARY}
	rm -rf coverage-all.out
	rm -rf ${BINARY}-*

# Docs

godoc-serve:	## Serve documentation (godoc format) for this package at port HTTP 9090
	godoc -http=":9090"

.PHONY: help
help:	## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
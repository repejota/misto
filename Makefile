BINARY=misto
VERSION=`cat VERSION`
BUILD=`git symbolic-ref HEAD 2> /dev/null | cut -b 12-`-`git log --pretty=format:%h -1`
PACKAGES = $(shell go list ./...)

# Setup the -ldflags option for go build here, interpolate the variable
# values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Build & Install

install:
	go install $(LDFLAGS) -v $(PACKAGES)

.PHONY: version
version:
	@echo $(VERSION)-$(BUILD)

# Testing

.PHONY: test
test:
	go test -v $(PACKAGES)

.PHONY: cover-profile
cover-profile:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	rm -rf coverage.out

.PHONY: cover
cover: cover-profile
	go tool cover -func=coverage-all.out

.PHONY: cover-html
cover-html: cover-profile
	go tool cover -html=coverage-all.out

.PHONY: coveralls
coveralls:
	goveralls -repotoken 9EmpV6j34d3itxKKXJCjTYicQPZhgzwj3

# Lint

lint:
	gometalinter --tests ./... --disable=gas

# Dependencies

deps:
	go get -v -t -d -u github.com/docker/docker/client
	go get -v -t -d -u github.com/fatih/color
	go get -v -t -d -u github.com/spf13/cobra
	go get -v -t -d -u github.com/repejota/cscanner

dev-deps:
	go get -v -t -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -v -t -u github.com/mattn/goveralls

# Cleaning up

.PHONY: clean
clean:
	go clean
	rm -rf ${BINARY}
	rm -rf coverage-all.out
	rm -rf ${BINARY}-*

# Docs

godoc-serve:
	godoc -http=":9090"

# Logs

logs:
	docker run -t --name logs --rm alpine echo "foooo"
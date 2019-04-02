# #VERSION=$(shell git describe --tags)
# LDFLAGS=-ldflags "-X main.Version=${VERSION} -s" # -w" (https://golang.org/cmd/link/)
# #.PHONY: build
# build:
# 	GO111MODULE=on GOOS=windows go build $(LDFLAGS) -o dist/lychee.exe *.go
# 	GO111MODULE=on GOOS=linux go build $(LDFLAGS) -o dist/lychee-linux *.go
# 	GO111MODULE=on GOOS=darwin go build $(LDFLAGS) -o dist/lychee-darwin *.go


# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY=lychee
LDFLAGS=-ldflags "-X main.Version=${VERSION} -s" # -w" (https://golang.org/cmd/link/)

.PHONY: release
release: windows linux darwin

.PHONY: linux
linux:
	mkdir -p release
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o release/$(BINARY)-v1.0.0-linux-amd64

.PHONY: darwin
darwin:
	mkdir -p release
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o release/$(BINARY)-v1.0.0-darwin-amd64

.PHONY: windows
windows:
	mkdir -p release
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o release/$(BINARY)-v1.0.0-windows-amd64

test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

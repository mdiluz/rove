VERSION := $(shell git describe --always --long --dirty --tags)

build:
	go build -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

install:
	go install -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

.PHONY: install
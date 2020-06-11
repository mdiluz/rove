VERSION := $(shell git describe --always --long --dirty --tags)

build:
	go mod download
	go build -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

install:
	go mod download
	go install -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

gen:
	protoc --proto_path pkg/accounts --go_out=plugins=grpc:pkg/accounts/ --go_opt=paths=source_relative  pkg/accounts/accounts.proto

test:
	docker-compose up --build --exit-code-from=rove-tests --abort-on-container-exit rove-tests
	go tool cover -html=/tmp/coverage-data/c.out -o /tmp/coverage.html
	@echo Done, coverage data can be found in /tmp/coverage.html

.PHONY: build install test gen

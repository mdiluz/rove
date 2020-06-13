VERSION := $(shell git describe --always --long --dirty --tags)

build:
	go mod download
	go build -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

install:
	go mod download
	go install -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

gen:
	@echo Generating accountant gRPC
	protoc --proto_path proto --go_out=plugins=grpc:pkg/ --go_opt=paths=source_relative  proto/accounts/accounts.proto
	@echo Generating rove server gRPC
	protoc --proto_path proto --go_out=plugins=grpc:pkg/ --go_opt=paths=source_relative  proto/rove/rove.proto
	protoc --proto_path proto --grpc-gateway_out=paths=source_relative:pkg/ proto/rove/rove.proto
	protoc --proto_path proto --swagger_out=logtostderr=true:pkg/ proto/rove/rove.proto

test:
	@echo Unit tests
	go test -v ./...

	@echo Integration tests
	docker-compose up --build --exit-code-from=rove-tests --abort-on-container-exit rove-tests
	go tool cover -html=/tmp/coverage-data/c.out -o /tmp/coverage.html
	
	@echo Done, coverage data can be found in /tmp/coverage.html

.PHONY: build install test gen

VERSION := $(shell git describe --always --long --dirty --tags)

build:
	@echo Running no-output build
	go mod download
	go build -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

install:
	@echo Installing to GOPATH
	go mod download
	go install -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

gen:
	@echo Installing go dependencies
	go install \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go
	go mod download
	@echo Generating rove server gRPC and gateway
	protoc --proto_path proto --go_out=plugins=grpc,paths=source_relative:pkg/ proto/rove/rove.proto
	protoc --proto_path proto --grpc-gateway_out=paths=source_relative:pkg/ proto/rove/rove.proto
	@echo Generating rove server swagger
	protoc --proto_path proto --swagger_out=logtostderr=true:pkg/ proto/rove/rove.proto

test:
	@echo Unit tests
	go test -v ./...

	@echo Integration tests
	docker-compose up --build --exit-code-from=rove-tests --abort-on-container-exit rove-tests
	docker-compose down
	go tool cover -html=/tmp/coverage-data/c.out -o /tmp/coverage.html
	
	@echo Done, coverage data can be found in /tmp/coverage.html

.PHONY: build install test gen

VERSION := $(shell git describe --always --long --dirty --tags)

build:
	go mod download
	go build -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

install:
	go mod download
	go install -ldflags="-X 'github.com/mdiluz/rove/pkg/version.Version=${VERSION}'" ./...

test:
	go mod download
	go build ./...

	# Run the server and shut it down again to ensure our docker-compose works
	ROVE_ARGS="--quit 1" docker-compose up --build --exit-code-from=rove-server --abort-on-container-exit

	# Run the server and shut it down again to ensure our docker-compose works
	docker-compose up -d

	# Run tests with coverage and integration tags
	go test -v ./... --tags=integration -cover -coverprofile=/tmp/c.out -count 1

	# 
	docker-compose down

	# Convert the coverage data to html
	go tool cover -html=/tmp/c.out -o /tmp/coverage.html

.PHONY: install test
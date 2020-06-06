#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Check we can build everything
go mod download
go build ./...

# Run the server and shut it down again to ensure our docker-compose works
ROVE_ARGS="--quit 1" docker-compose up --build --exit-code-from=rove-server --abort-on-container-exit

# Run tests with coverage
go test -v ./... -cover -coverprofile=/tmp/c.out

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
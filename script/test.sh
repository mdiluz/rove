#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Build and start rove-server
docker-compose up --detach --build

# Run tests, including integration tests
go mod download
go test -v ./... -tags integration -cover -coverprofile=/tmp/c.out

# Take down the service
docker-compose down

# Check that the cmdline client builds
docker build -f "cmd/rove/Dockerfile" .

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
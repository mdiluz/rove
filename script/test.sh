#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Test the build
go build -v .

# Run unit tests
go test -v ./... -cover

# Verify docker build
docker build .

# Run the integration tests with docker-compose
docker-compose up --build --detach
go test -v ./... -tags integration -cover -coverprofile=/tmp/c.out
docker-compose down

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
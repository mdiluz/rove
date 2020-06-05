#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Check that the cmdline client builds
docker build -f "cmd/rove/Dockerfile" .

# Build and start rove-server
docker-compose -f docker-compose-test.yml up --detach --build

# Run tests, including integration tests
go mod download
go test -v ./... -tags integration -cover -coverprofile=/tmp/c.out

# Take down the service
docker-compose down

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
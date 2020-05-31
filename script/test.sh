#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Build and start the service
docker-compose up --build --detach

# Run tests, including integration tests
go test -v ./... -tags integration -cover -coverprofile=/tmp/c.out

# Take down the service
docker-compose down

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
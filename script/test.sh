#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Build the image
bash script/build.sh

# Build and start the service
docker-compose up --detach

# Run tests, including integration tests
go mod download
go test -v ./... -tags integration -cover -coverprofile=/tmp/c.out

# Take down the service
docker-compose down

# Convert the coverage data to html
go tool cover -html=/tmp/c.out -o /tmp/coverage.html
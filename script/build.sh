#!/bin/bash
set -e
cd "$(dirname "$0")"
cd ..
set -x

# Generate a version string
export VERSION=$(git describe --always --long --dirty --tags)

# Build and tag as latest and version
docker build -t "rove:latest" -t "rove:${VERSION}" -f "cmd/rove/Dockerfile" .
docker build -t "rove-server:latest" -t "rove-server:${VERSION}" -f "cmd/rove-server/Dockerfile" .
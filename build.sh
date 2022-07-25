#!/bin/bash -

set -o nounset # Treat unset variables as an error
# for mac use below
#env GOOS=darwin GOARCH=arm64 go build -o halogp_linux ./...
#for linux ubuntu use below
env GOOS=linux GOARCH=amd64 go build -o halogp_linux ./...

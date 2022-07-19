#!/bin/bash -

set -o nounset # Treat unset variables as an error

env GOOS=linux GOARCH=amd64 go build -o halogp_linux ./...

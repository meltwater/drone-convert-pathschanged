#!/bin/sh

# disable go modules
export GOPATH=""

# disable cgo
export CGO_ENABLED=0

set -e
set -x

# linux
GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/drone-convert-pathschanged
GOOS=linux GOARCH=arm64 go build -o release/linux/arm64/drone-convert-pathschanged

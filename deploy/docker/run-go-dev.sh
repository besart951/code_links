#!/bin/sh
set -eu

if [ ! -d /go/pkg/mod/cache/download ]; then
	go mod download
fi

exec go run .

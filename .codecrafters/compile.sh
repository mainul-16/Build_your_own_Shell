#!/bin/sh
set -e

cd "$(dirname "$0")/.."
go build -o /tmp/codecrafters-build-shell-go .

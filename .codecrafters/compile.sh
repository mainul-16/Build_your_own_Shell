#!/bin/sh
#
# This script is used to compile your program on CodeCrafters
#
# This runs before .codecrafters/run.sh
#

set -e

go build -o /tmp/codecrafters-build-shell-go .

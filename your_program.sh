#!/bin/sh
#
# Use this script to run your program LOCALLY.
#
# Note: Changing this script WILL NOT affect how CodeCrafters runs your program.
#
# Learn more: https://codecrafters.io/program-interface

set -e # Exit early if any commands fail

(
  cd "$(dirname "$0")"
  go build -o /tmp/codecrafters-build-http-server-go ./app
)

exec /tmp/codecrafters-build-http-server-go "$@"

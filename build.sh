#!/bin/sh

set -e

mkdir build 2>/dev/null || true
cd build
GOOS=linux CGO_ENABLED=0 go build ../bin/nightbot
zip nightbot.zip nightbot

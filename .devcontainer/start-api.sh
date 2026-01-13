#!/bin/sh
set -e

export GOPATH=/root/go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH

echo "Starting API server..."
cd /workspace/api

# Use full path to air
AIR_BIN=$GOBIN/air

# Verify air exists, install if not
if [ ! -f "$AIR_BIN" ]; then
  echo "Installing air to $AIR_BIN..."
  go install github.com/cosmtrek/air@latest
fi

echo "Running air from $AIR_BIN..."
exec $AIR_BIN -c .air.toml

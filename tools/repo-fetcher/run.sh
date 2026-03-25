#!/usr/bin/env bash

set -e

cd "$(dirname "$0")"

echo "Building repo-fetcher..."
go build -o ../bin/repo-fetcher

echo "Running repo-fetcher..."
../bin/repo-fetcher

echo "Done."

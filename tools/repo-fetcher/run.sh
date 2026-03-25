#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"
if [ ! -f ../bin/repo-fetcher ] || [ main.go -nt ../bin/repo-fetcher ]; then
  go build -o ../bin/repo-fetcher
fi
../bin/repo-fetcher
cp projects.json /var/www/html/
chmod 644 /var/www/html/projects.json

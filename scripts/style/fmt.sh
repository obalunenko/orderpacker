#!/bin/bash

set -Eeuo pipefail

SCRIPT_NAME="$(basename "$0")"

echo "${SCRIPT_NAME} is running... "

echo "Making filelist"
GO_FILES=($(find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./.git/*"))

for f in "${GO_FILES[@]}"; do
  echo "Fixing fmt at ${f}"
    gofmt -s -w "$f"
done

echo "${SCRIPT_NAME} done."

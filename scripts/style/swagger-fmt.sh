#!/bin/bash

SCRIPT_NAME="$(basename "$0")"
SCRIPT_DIR="$(dirname "$0")"
REPO_ROOT="$(cd "${SCRIPT_DIR}" && git rev-parse --show-toplevel)"
SCRIPTS_DIR="${REPO_ROOT}/scripts"
source "${SCRIPTS_DIR}/helpers-source.sh"

DOCS_DIR="${REPO_ROOT}/docs"

echo "${SCRIPT_NAME} is running... "

echo "Installing swaggo/swag"
raw_version=$(go list -m github.com/swaggo/swag)
echo "Raw version: $raw_version"
version=$(echo $raw_version | awk '{print $2}')
echo "Version: $version"
echo "Installing github.com/swaggo/swag/cmd/swag@$version"
go install github.com/swaggo/swag/cmd/swag@$version

echo "Formatting annotations"

swag fmt --dir ./cmd/orderpacker,./internal/service


echo "${SCRIPT_NAME} is done."
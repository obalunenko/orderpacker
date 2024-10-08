APP_NAME?=orderpacker
SHELL := env APP_NAME=$(APP_NAME) $(SHELL)

BIN_DIR?=$(CURDIR)/bin

GOVERSION:=1.23

TEST_DISCARD_LOG?=false
SHELL := env TEST_DISCARD_LOG=$(TEST_DISCARD_LOG) $(SHELL)

format-code: swagger-fmt fmt goimports
.PHONY: format-code

fmt:
	@echo "Formatting code..."
	./scripts/style/fmt.sh
.PHONY: fmt

goimports:
	@echo "Formatting code..."
	./scripts/style/fix-imports.sh
.PHONY: goimports

vet:
	@echo "Vetting code..."
	@go vet ./...
	@echo "Done"
.PHONY: vet

test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Done"
.PHONY: test

build: swagger-gen
	@echo "Building..."
	@./scripts/build/app.sh
	@echo "Done"
.PHONY: build

run:
	@echo "Running..."
	@${BIN_DIR}/$(APP_NAME)
	@echo "Done"
.PHONY: run

vendor:
	@echo "Vendoring..."
	@go mod tidy && go mod vendor
	@echo "Done"
.PHONY: vendor

docker-build:
	@echo "Building docker image..."
	@docker build -t $(APP_NAME):latest -f Dockerfile .
	@echo "Done"
.PHONY: docker-build

docker-run: docker-build
	@echo "Running docker image..."
	@docker compose -f compose.yaml up
	@echo "Done"
.PHONY: docker-run

docker-stop:
	@echo "Stopping docker image..."
	@docker compose -f compose.yaml down
	@echo "Done"
.PHONY: docker-stop

## Release
release:
	./scripts/release/release.sh
.PHONY: release

## Release local snapshot
release-local-snapshot:
	./scripts/release/local-snapshot-release.sh
.PHONY: release-local-snapshot

## Check goreleaser config.
check-releaser:
	./scripts/release/check.sh
.PHONY: check-releaser

## Issue new release.
new-version: vet test build docker-build
	./scripts/release/new-version.sh
.PHONY: new-release

## Bump go version
bump-go-version:
	./scripts/bump-go.sh $(GOVERSION)
.PHONY: bump-go-version

## Generate swagger docs
swagger-gen:
	./scripts/swagger-docs.sh
.PHONY: swagger-gen

## Format swagger annotations
swagger-fmt:
	./scripts/style/swagger-fmt.sh
.PHONY: swagger-fmt


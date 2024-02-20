APP_NAME?=orderpacker
SHELL := env APP_NAME=$(APP_NAME) $(SHELL)

GOVERSION:=1.22

TEST_DISCARD_LOG?=false
SHELL := env TEST_DISCARD_LOG=$(TEST_DISCARD_LOG) $(SHELL)

format-code: fmt goimports
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

build:
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
new-version: vet test build
	./scripts/release/new-version.sh
.PHONY: new-release

bump-go-version:
	./scripts/bump-go.sh $(GOVERSION)
.PHONY: bump-go-version

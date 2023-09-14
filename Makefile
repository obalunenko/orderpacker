APP_NAME := "orderpacker"

format-code: fmt goimports
.PHONY: format-code

fmt:
	@echo "Formatting code..."
	./scripts/fmt.sh
.PHONY: fmt

goimports:
	@echo "Formatting code..."
	./scripts/fix-imports.sh
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
	@go build -o bin/$(APP_NAME) -v ./cmd/$(APP_NAME)
	@echo "Done"
.PHONY: build

vendor:
	@echo "Vendoring..."
	@go mod tidy && go mod vendor
	@echo "Done"
.PHONY: vendor

docker-build:
	@echo "Building docker image..."
	@docker build -t $(APP_NAME) .
	@echo "Done"
.PHONY: docker-build

docker-run:
	@echo "Running docker image..."
	@docker compose -f compose.yaml up --detach --build
	@echo "Done"
.PHONY: docker-run

docker-stop:
	@echo "Stopping docker image..."
	@docker compose -f compose.yaml down
	@echo "Done"
.PHONY: docker-stop

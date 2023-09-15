![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/obalunenko/orderpacker)
[![Latest release artifacts](https://img.shields.io/github/v/release/obalunenko/instadiff-cli)](https://github.com/obalunenko/instadiff-cli/releases/latest)
[![Go [lint, test]](https://github.com/obalunenko/orderpacker/actions/workflows/go.yml/badge.svg)](https://github.com/obalunenko/orderpacker/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/obalunenko/orderpacker)](https://goreportcard.com/report/github.com/obalunenko/orderpacker)

# OrderPacker Service

## What is OrderPacker?

OrderPacker is a Golang based application that calculates the number of packs needed to ship to a customer.

## How does OrderPacker work?

The application exposes its functionality through an HTTP API and accepts a JSON payload with the following structure:

```json
{
  "items": 501
}
```

The `items` field is a positive integer that represents the number of items that need to be packed.

The application responds with a JSON payload with the following structure:

```json
{
  "packs": [
    {
      "pack": 250,
      "quantity": 2
    },
    {
      "pack": 1,
      "quantity": 1
    }
  ]
}
```

It primarily runs on `localhost` port `8080` and acts upon `GET` requests to the `/pack` endpoint.

Below is a Curl command snippet demonstrating how to call this endpoint

```bash
curl --location --request GET 'localhost:8080/pack' \
--header 'Content-Type: application/json' \
--data '{
    "items": 501
}'
```

## Configuration

The application can be configured using the following environment variables:

    -`ORDERPACKER_CONFIG_PATH` - path to the configuration file. Default value is empty.

If the `ORDERPACKER_CONFIG_PATH` environment variable is not set, the application will use default configuration values.

The configuration file is a JSON file with the following structure:

```json
{
  "http": {
    "port": "8080"
  },
  "pack": {
    "boxes": [
      1,2,4,8,16,32
    ]
  }
}
```

YAML configuration file is also supported.

```yaml 
http:
  port: 8080
pack:
  - 1
  - 2
  - 4
  - 8
  - 16
  - 32
```


## Development

### Prerequisites

- [Go](https://golang.org/doc/install) 1.21 or higher
- [Docker](https://docs.docker.com/get-docker/) 24.0 or higher
- [Docker Compose](https://docs.docker.com/compose/install/) 2.21 or higher

### Running the application

For development purposes, the application can be run locally using the following command:

```bash
make build && make run
```

To run the application in a Docker container, use the following command:

```bash
make docker-build && make docker-run
```

### Running tests

To run the tests, use the following command:

```bash
make test
```

### Linting

To run the linter, use the following command:

```bash
make vet
```

### Code formatting

To format the code, use the following command:

```bash
make format-code
```

### Vendoring

To vendor the dependencies, use the following command:

```bash
make vendor
```

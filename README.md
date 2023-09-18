![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/obalunenko/orderpacker)
[![Latest release artifacts](https://img.shields.io/github/v/release/obalunenko/orderpacker)](https://github.com/obalunenko/orderpacker/releases/latest)
[![Go [lint, test]](https://github.com/obalunenko/orderpacker/actions/workflows/go.yml/badge.svg)](https://github.com/obalunenko/orderpacker/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/obalunenko/orderpacker)](https://goreportcard.com/report/github.com/obalunenko/orderpacker)

# OrderPacker Service

## What is OrderPacker?

OrderPacker is a Golang based application that calculates the number of packs needed to ship to a customer.

## How does OrderPacker work?

### Frontend
The OrderPacker service also provides a user-friendly frontend, from which the aforementioned API can be conveniently accessed and tested. 
You can reach the frontend from your browser at:

`http://localhost:8080`

The frontend itself is quite minimalistic - it contains an input field for submitting the number of items to be packed, 
and upon submission, it presents neatly formatted API responses. 
The responses are conveniently displayed, showing each pack and the corresponding quantity.

### API

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

It primarily runs on `localhost` port `8080` and acts upon `POST` requests to the `/pack` endpoint.

Below is a Curl command snippet demonstrating how to call this endpoint

```bash
curl --location --request POST 'localhost:8080/pack' \
--header 'Content-Type: application/json' \
--data '{
    "items": 501
}'
```

## Configuration

Application follows the [12-factor app](https://12factor.net/) methodology and can be configured using environment variables.

Following environment variables are supported:

| Name         | Description                                                          | Default value             |
|--------------|----------------------------------------------------------------------|---------------------------|
| `PORT`       | The port on which the application will listen for incoming requests. | `8080`                    |
| `HOST`       | The host on which the application will listen for incoming requests. | `0.0.0.0`                 |
| `LOG_LEVEL`  | The log level of the application.                                    | `info`                    |
| `LOG_FORMAT` | The log format of the application.                                   | `text`                    |
| `PACK_BOXES` | The pack boxes for packing orders. Values should be separated by `,` | `250,500,1000,2000,5000,` |


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

To run tests without logs, use the following command:

```bash
TEST_DISCARD_LOGS=true make test
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

# Cold Chain Fleet Dispatch Service

A B2B service that assigns refrigerated shipments to compatible vehicles and
drivers, records in-transit temperature readings, and raises alerts when cargo
leaves its allowed temperature range.

## Run

```bash
go run ./cmd/server
```

The HTTP API listens on `HTTP_ADDR` (default `127.0.0.1:18080`).

## Test

```bash
go test ./...
```

## API

- `GET /health`
- `POST /v1/shipments`
- `GET /v1/shipments/{id}`
- `POST /v1/dispatch/{id}`
- `POST /v1/shipments/{id}/readings`
- `POST /v1/shipments/{id}/start-transit`
- `POST /v1/shipments/{id}/complete-delivery`
- `GET /v1/alerts`

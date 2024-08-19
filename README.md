# Tracing

Neste vídeo é apresentado como integrar tracing a sua aplicação Go.

## YouTube

https://www.youtube.com/watch?v=f6F5zoEEbTw

## Docker

```sh
docker compose up -d
```

UI -> http://localhost:16686

## Setup

```sh
export OTEL_EXPORTER_OTLP_INSECURE=true
cd checkouts
go mod tidy
```

## Run

```sh
go run ./app/cmd/server
```

FROM golang:1.24.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

ENV CGO_ENABLED=1

RUN go build -o server ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /root

RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server .

COPY config/local.yaml ./config.yaml

EXPOSE 8080

CMD ["./server"]
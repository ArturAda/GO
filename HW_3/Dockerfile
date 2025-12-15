FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o balance_server ./cmd/balance_server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /workspace/balance_server /app/balance_server

RUN mkdir -p /files

ENV HTTP_ADDRESS=:8080 STORAGE_TYPE=file JSON_PATH=/files/balances.json

EXPOSE 8080

ENTRYPOINT ["/app/balance_server"]

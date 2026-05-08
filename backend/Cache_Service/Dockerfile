# ===== Stage 1: Builder =====
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/cache-service ./cmd/cacheService/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o /app/bin/migrator ./cmd/migrator/main.go

# ===== Stage 2: Final image =====
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata postgresql-client

RUN adduser -D -g '' appuser

COPY --from=builder --chown=appuser:appuser /app/bin/cache-service /app/cache-service
COPY --from=builder --chown=appuser:appuser /app/bin/migrator /app/migrator
COPY --from=builder --chown=appuser:appuser /app/config /app/config
COPY --from=builder --chown=appuser:appuser /app/migrations /app/migrations

USER appuser

WORKDIR /app

EXPOSE 44044 44045

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/cache-service", "health"] || exit 1

ENTRYPOINT ["/app/cache-service"]
CMD ["--config=./config/local.yaml"]
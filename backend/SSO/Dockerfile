
FROM golang:1.25-alpine

RUN apk add --no-cache git ca-certificates postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY ./config ./config
COPY ./migrations ./migrations

RUN go build -o main ./cmd/sso
RUN go build -o migrator ./cmd/migrator

RUN chmod +x main migrator


RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 44044

CMD ["sh", "-c", "\
  echo 'Waiting for database...' && \
  until pg_isready -h ${DB_HOST:-aurora-db} -p ${DB_PORT:-5432} -U ${DB_USER:-postgres}; do sleep 2; done && \
  echo 'Running migrations...' && \
  ./migrator -db-dsn \"$DATABASE_URL\" -migrations-path ./migrations -migrations-table schema_migrations up && \
  echo 'Starting app...' && \
  ./main --config ./config/local.yaml \
"]

FROM golang:1.25-alpine

USER root

RUN sed -i 's/dl-cdn.alpinelinux.org/ftp.halifax.rwth-aachen.de/g' /etc/apk/repositories

RUN apk add --no-cache git ca-certificates postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY ./migrations ./migrations

RUN go build -o sso ./cmd/sso
RUN go build -o migrator ./cmd/migrator

RUN chmod +x sso migrator

# RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 44044

CMD ["sh", "-c", "\
  echo 'Waiting for database...' && \
  until pg_isready -h postgres -p 5432 -U postgres; do sleep 2; done && \
  echo 'Running migrations...' && \
  ./migrator -db-dsn=$STORAGE_PASS -migrations-path=./migrations -migrations-table=schema_migrations up && \
  echo 'Starting app...' && \
  ./sso"]

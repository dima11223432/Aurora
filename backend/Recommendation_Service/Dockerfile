FROM golang:1.25-alpine

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o recommendation_service ./cmd/recomendation_service/main.go

RUN chmod +x recommendation_service

RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 44000

CMD ["./recommendation_service", "--config", "./config/local.yaml"]

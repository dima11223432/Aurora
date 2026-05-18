FROM golang:1.25-alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/ftp.halifax.rwth-aachen.de/g' /etc/apk/repositories

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY ./config ./config

RUN go build -o recommendation_service ./cmd/recomendation_service/main.go

RUN chmod +x recommendation_service

RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 44040

CMD ["./recommendation_service"]

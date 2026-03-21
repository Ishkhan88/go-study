FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o concert-service ./concert/cmd/server

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/concert-service .

EXPOSE 8082
CMD ["./concert-service"]
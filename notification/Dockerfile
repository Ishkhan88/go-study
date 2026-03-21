FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o notification-service ./notification/cmd/server

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/notification-service .

CMD ["./notification-service"]
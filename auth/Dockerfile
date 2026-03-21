FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-service ./auth/cmd/server

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/auth-service .

EXPOSE 8081
CMD ["./auth-service"]
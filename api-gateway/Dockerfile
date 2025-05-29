FROM golang:1.23.5-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./api-gateway ./api-gateway
COPY ./proto ./proto

WORKDIR /app/api-gateway
RUN go build -o api-gateway_service ./cmd/api-gateway/main.go

EXPOSE 8080
CMD ["./api-gateway_service"]
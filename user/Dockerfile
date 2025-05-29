FROM golang:1.23.5-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./user ./user
COPY ./proto ./proto

WORKDIR /app/user
RUN go build -o user_service ./cmd/user/main.go

EXPOSE 4000
CMD ["./user_service"]
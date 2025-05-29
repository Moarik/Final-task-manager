FROM golang:1.23.5-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./task ./task
COPY ./proto ./proto

WORKDIR /app/task
RUN go build -o task_service ./cmd/task/main.go

EXPOSE 4001
CMD ["./task_service"]
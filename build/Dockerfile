FROM golang:1.22-alpine AS builder
WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "/"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/main ./cmd/main.go
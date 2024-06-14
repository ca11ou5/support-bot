FROM registry.sovcombank.group/web-ecom/golang-1.22-alpine-certs@sha256:59ac4e11705439e6b17a1ab9b45b28cd0352f5fa81257abfdc46d82e2c7c8962 AS builder
WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "/"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/main ./cmd/main.go

FROM alpine:latest AS runner

COPY --from=builder /usr/local/src/bin/main main
COPY ./config/envs/dev.env /envs/dev.env

CMD ["/main"]
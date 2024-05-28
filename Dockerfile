FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o main ./cmd/main.go

CMD ./main




# Use Go 1.23.6 as base image 
FROM golang:1.23.6 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o oasis ./cmd/server
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/oasis .

CMD ["./oasis"]

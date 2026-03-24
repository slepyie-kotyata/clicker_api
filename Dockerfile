FROM golang:1.25.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o clicker_api server.go
FROM ubuntu:latest
WORKDIR /root/
RUN apt update && apt install -y libc6
COPY --from=builder /app/clicker_api .

ARG DB_CONNECTION
ARG REDIS_CONNECTION
ARG ACCESS_TOKEN_SECRET
ARG REFRESH_TOKEN_SECRET

RUN echo "DB_CONNECTION=${DB_CONNECTION}" > .env && \
    echo "REDIS_CONNECTION=${REDIS_CONNECTION}" >> .env && \
    echo "ACCESS_TOKEN_SECRET=${ACCESS_TOKEN_SECRET}" >> .env && \
    echo "REFRESH_TOKEN_SECRET=${REFRESH_TOKEN_SECRET}" >> .env

EXPOSE 1323
CMD ["./clicker_api"]

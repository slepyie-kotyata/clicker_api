FROM golang:1.25.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/clicker_api ./services/main_api
FROM ubuntu:latest
WORKDIR /root/
RUN apt update && apt install -y libc6
COPY --from=builder /app/clicker_api ./clicker_api
COPY --from=builder /app/.env .env
EXPOSE 1323
CMD ["./clicker_api"]

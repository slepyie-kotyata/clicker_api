FROM golang:1.24.1
WORKDIR /clicker
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /clicker
RUN chmod +x /clicker
EXPOSE 1323
CMD ["/clicker"]

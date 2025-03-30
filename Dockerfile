FROM golang:1.24.1
WORKDIR /clicker_api
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /clicker
EXPOSE 1323
CMD ["/clicker_api"]
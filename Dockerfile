# syntax=docker/dockerfile:1
FROM golang:1.22.0
WORKDIR /
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o ./cmd/server/main ./cmd/server/.
EXPOSE 8080
CMD ["./cmd/server/main"]
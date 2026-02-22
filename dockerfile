# Stage 1: Build
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.io

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o peyk ./main.go

# Stage 2: Run 
FROM alpine:latest

WORKDIR /app

# Copy the binary
COPY --from=builder /app/peyk .

EXPOSE 8080

ENTRYPOINT ["./peyk"]

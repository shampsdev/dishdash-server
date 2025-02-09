# This Dockerfile must be built with the context of root directory of repository
# e.g. in this folder: docker build -t my-image -f ./Dockerfile ..

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY dashboard/go.mod dashboard/go.sum ./
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY dashboard .
COPY server ../server

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

CMD ["./server"]

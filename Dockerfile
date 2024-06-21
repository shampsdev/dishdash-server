# Build stage
FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o dishdash-server ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/dishdash-server .

CMD ["./dishdash-server"]

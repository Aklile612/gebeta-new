# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go.mod and go.sum from backend folder
COPY backend/go.mod backend/go.sum ./

RUN go mod download

# Copy the entire backend folder
COPY backend/. .

# 👇 Build from the correct main package directory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/server

# Stage 2: Run
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]

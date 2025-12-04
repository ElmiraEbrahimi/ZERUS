# syntax=docker/dockerfile:1

# Build the Go binary using a minimal Go image
FROM golang:1.20-alpine AS builder

WORKDIR /src

# Install build dependencies and prepare modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Final runtime image
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/api ./api

# Expose the default port (can be overridden)
EXPOSE 8080

# Run the service
CMD ["./api"]
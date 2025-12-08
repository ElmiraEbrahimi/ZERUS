# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copy go module files and download dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source.
COPY . .

# Build the server binary. Disable cgo to simplify the runtime image.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /counter-api ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /counter-api .
EXPOSE 8080
ENTRYPOINT ["/counter-api"]
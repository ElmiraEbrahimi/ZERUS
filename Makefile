.PHONY: build test up down run

# Build the Go binary using local environment
build:
	go build -v -o bin/api ./cmd/api

# Run all unit tests
test:
	go test ./...

# Start the full development stack using docker-compose
up:
	docker-compose up -d

# Stop and remove the docker-compose services
down:
	docker-compose down

# Run the API directly (requires local Postgres and Ethereum nodes)
run:
	go run ./cmd/api
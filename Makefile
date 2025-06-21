.PHONY: build run dev test clean migrate-up migrate-down install-dev fmt lint tidy swagger

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Run in development mode with auto-reload (install air first: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations up
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

# Install development dependencies
install-dev:
	go install github.com/cosmtrek/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest # Add swag to dev dependencies install

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go

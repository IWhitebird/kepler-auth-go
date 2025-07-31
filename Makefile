.PHONY: build run test clean docker-build docker-run swagger deps migrate

APP_NAME=kepler-auth-go
BUILD_DIR=bin
DOCKER_IMAGE=skylarklabs/kepler-auth-go

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs/

migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go -action=up

migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -action=status

migrate-rollback:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go -action=down

migrate-fresh:
	@echo "Running fresh migrations..."
	@go run cmd/migrate/main.go -action=fresh

migrate-force:
	@echo "Force running migrations..."
	@go run cmd/migrate/main.go -action=up -force

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):latest .

docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose up --build

docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

dev:
	@echo "Starting development environment..."
	@docker-compose up -d db
	@sleep 5
	@make run

install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest

lint:
	@echo "Running linter..."
	@golangci-lint run

format:
	@echo "Formatting code..."
	@go fmt ./...

security:
	@echo "Running security scan..."
	@gosec ./...

all: clean deps build test

production: clean deps build docker-build
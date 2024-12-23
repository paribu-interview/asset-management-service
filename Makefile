# Variables
APP_NAME := service
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_IMAGE := $(APP_NAME):latest
BINARY_NAME := $(APP_NAME)

# Commands
.PHONY: all build run test clean docker-build docker-run docker-clean

# Default target
all: build

# Build the Go application
build:
	@echo "Building the Go application..."
	go build -o $(BINARY_NAME) ./cmd/main.go

# Run the application locally
run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Clean build files
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run the application in Docker
docker-run: docker-build
	@echo "Running the application in Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

# Stop and clean Docker containers and images
docker-clean:
	@echo "Stopping and removing Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes
	@echo "Removing Docker images..."
	docker rmi $(DOCKER_IMAGE) || true

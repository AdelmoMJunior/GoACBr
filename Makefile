.PHONY: build run test lint swagger migrate-up migrate-down docker-up docker-down clean help

APP_NAME := goacbr-api
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/api

# Go parameters
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED := 1

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application binary
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application locally
	go run $(MAIN_PATH)/main.go

test: ## Run unit tests
	go test -v -race -count=1 ./...

test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

swagger: ## Generate swagger documentation
	swag init -g $(MAIN_PATH)/main.go -o ./docs --parseInternal --parseDependency

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Roll back the last migration
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_xyz)
	migrate create -ext sql -dir migrations -seq $(NAME)

docker-up: ## Start all Docker containers
	docker-compose up -d --build

docker-down: ## Stop all Docker containers
	docker-compose down

docker-logs: ## Tail Docker container logs
	docker-compose logs -f

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR) coverage.out coverage.html

tidy: ## Tidy Go modules
	go mod tidy

fmt: ## Format Go code
	gofmt -s -w .

.PHONY: build run clean test

# Build the application
build:
	go build -o bin/go-web cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run with development settings
dev:
	PORT=8080 go run cmd/main.go

# Stop the running server
stop:
	@echo "Stopping go-web server..."
	@-lsof -i :8080 -t | xargs kill -9 2>/dev/null || echo "No process on port 8080"
	@-pkill -f "go-web" 2>/dev/null || echo "No server binary running"
	@-pkill -f "go run cmd/main.go" 2>/dev/null || echo "No development server running"

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/go-web cmd/main.go

# Docker build
docker-build:
	docker build -t go-fork-web .

# Install dependencies
deps:
	go mod tidy
	go mod download

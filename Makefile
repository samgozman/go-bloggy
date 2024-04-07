# Download go dependencies
download:
	go mod download

# Lint go files
lint:
	golangci-lint run -v

# Test go files
test:
	go test -v ./... -race

# Test go files with coverage report
test-coverage:
	go test -v ./... -race -coverprofile=coverage.out -covermode=atomic

# Generate go files (api)
generate:
	go generate ./...

# Create copy .env.sample to .env if not exists
env:
	cp -n .env.sample .env || true

# Initialize project
init : download env generate

# Run local server without docker
run:
	go run cmd/server/main.go

# Run via docker compose
docker-run:
	docker compose up --build

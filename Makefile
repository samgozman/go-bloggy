# Lint go files
lint:
	golangci-lint run

# Test go files
test:
	go test -v ./... -race

# Test go files with coverage report
test-coverage:
	go test -v ./... -race -coverprofile=coverage.txt -covermode=atomic

# Generate go files (api)
generate:
	go generate ./...

# Run local server without docker
run:
	go run cmd/server/main.go

# Lint go files
lint:
	golangci-lint run

# Test go files
test:
	go test -v ./... -race

# Generate go files (api)
generate:
	go generate ./...

# Run local server without docker
run:
	go run cmd/server/main.go

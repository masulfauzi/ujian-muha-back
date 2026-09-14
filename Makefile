.PHONY: help build run test clean migrate docker-build docker-run

BINARY_NAME=backend
GO=go
GOFLAGS=-v

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make install-deps - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linter"

install-deps:
	$(GO) mod download
	$(GO) mod tidy

build:
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

run:
	$(GO) run $(GOFLAGS) cmd/server/main.go

test:
	$(GO) test $(GOFLAGS) ./...

test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

clean:
	$(GO) clean
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	@which golangci-lint > /dev/null || echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin"
	golangci-lint run ./...

docker-build:
	docker build -t $(BINARY_NAME):latest .

docker-run:
	docker run -p 3000:3000 --env-file .env $(BINARY_NAME):latest

dev:
	$(GO) run cmd/server/main.go

.PHONY: dev

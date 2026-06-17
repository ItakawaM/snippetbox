include .env

export

BINARY ?= server.exe
MAIN_PATH := ./cmd/web/

.DEFAULT_GOAL := help

.PHONY: run build build-run clean test help

run:
	go run $(MAIN_PATH) -addr=$(HOST_ADDRESS) -dsn=$(POSTGRESQL_DSN)

build:
	go build -o $(BINARY) $(MAIN_PATH)

build-run: build
	./$(BINARY) -addr=$(HOST_ADDRESS) -dsn=$(POSTGRESQL_DSN)

test:
	go test ./... -v

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY)
	go clean
	@echo "Done!"

help:
	@echo "Commands:"
	@echo "  make run         - Runs the Web Server"
	@echo "  make build       - Builds the Web Server into an executable"
	@echo "  make build-run   - Runs the Web Server executable"
	@echo "  make clean       - Removes the compiled binary"
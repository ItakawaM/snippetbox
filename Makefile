BINARY ?= server.exe

.PHONY: run help

run:
	go run ./cmd/web/

build:
	go build -o $(BINARY) ./cmd/web/

help:
	@echo "Commands:"
	@echo "  make run     - Runs the Web Server"
	@echo "  make build   - Builds the Web Server into an executable"
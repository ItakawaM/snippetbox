BINARY ?= server.exe
ADDRESS ?= :4000

.PHONY: run help

run:
	go run ./cmd/web/ -addr=$(ADDRESS)

build:
	go build -o $(BINARY) ./cmd/web/ -addr=$(ADDRESS)

help:
	@echo "Commands:"
	@echo "  make run     - Runs the Web Server"
	@echo "  make build   - Builds the Web Server into an executable"
include .env

export

BINARY ?= server.exe

.PHONY: run help

run:
	go run ./cmd/web/ -addr=$(HOST_ADDRESS) -dsn=$(POSTGRESQL_DSN)

build:
	go build -o $(BINARY) ./cmd/web/ -addr=$(HOST_ADDRESS)

help:
	@echo "Commands:"
	@echo "  make run     - Runs the Web Server"
	@echo "  make build   - Builds the Web Server into an executable"
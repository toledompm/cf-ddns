MAIN = main.go
BIN ?= ./dist/cf-ddns

.PHONY: build
build:
	@echo "Building..."
	@go build -o $(BIN) $(MAIN)

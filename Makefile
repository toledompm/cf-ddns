MAIN = main.go
BIN ?= ./dist/cf-ddns

.PHONY: build
build:
	@echo "Building..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN)-darwin-arm $(MAIN)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN)-linux-amd64 $(MAIN)

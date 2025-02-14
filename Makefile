.PHONY: all build install clean test run help

all: build

build:
	go build -o redir cmd/redir/main.go redir.go

install:
	go install github.com/4thel00z/redir@latest

clean:
	rm -f redir

test:
	go test ./...

run:
	./redir -url="https://example.com" -output=table

help:
	@echo "Makefile targets:"
	@echo "  all     - Build the project (alias for build)"
	@echo "  build   - Build the CLI executable 'redir'"
	@echo "  install - Install the CLI using go install"
	@echo "  clean   - Remove the built executable"
	@echo "  test    - Run tests"
	@echo "  run     - Run the CLI with a sample URL"
	@echo "  help    - Display this help message"

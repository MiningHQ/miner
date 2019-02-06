#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# This makes the APP_NAME be the name of the current directory
# Ex. in path /home/dev/app/awesome-app the APP_NAME will be set to awesome-app
APP_NAME := $(notdir $(CURDIR))

default: build ## Build the binary

build: build_linux ## Build binaries for Windows and Linux
	make build_windows

build_linux: ## Build the binary for Linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/${APP_NAME} ./src/*.go

build_windows: ## Build the binary for Windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/${APP_NAME}.exe ./src/*.go

clean: ## Remove compiled binaries from bin/
	rm ./bin/*

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Name of the app
APP_NAME := 'MiningHQ\ Miner\ Manager'

default: build ## Build the binary

build: ## Run the bundler to build the GUI
	cd src/; astilectron-bundler -v

run: build ## Build and run the GUI for Linux
	./bin/linux-amd64/'${APP_NAME}'

run_debug: build ## Build and run the GUI for Linux in debug mode
	./bin/linux-amd64/'${APP_NAME}' -d

run_only_debug: ## Run the GUI for Linux without building
	./bin/linux-amd64/'${APP_NAME}' -d

webdev: ## Run the web interface as a stanalone site (for development)
	./scripts/run_browsersync.sh

fmt: ## Format the code using `go fmt`
	go fmt ./...

test: ## Run the tests
	go test ./...

test_cover: ## Run tests with a coverage report
	go test ./... -v -cover -covermode=count -coverprofile=./coverage.out

clean: ## Remove compiled binaries from bin/
	rm -Rf bin/
	rm src/windows.syso
	rm src/bind*

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

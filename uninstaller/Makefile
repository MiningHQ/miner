#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Name of the app
APP_NAME := 'uninstall-mininghq'

build_linux: ## Build the binary for linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/${APP_NAME} ./src/*.go

build_windows: ## Build the binary for Windows
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -o ./bin/${APP_NAME}.exe ./src/*.go

run: build_linux ## Build and run the binary in CLI mode
	./bin/${APP_NAME}

clean: ## Remove compiled binaries from bin/
	rm -Rf bin/

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

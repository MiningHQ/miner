#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Name of the app
APP_NAME := 'MiningHQ\ Miner\ Manager'

build_linux: ## Build the binary for Linux
	cd src/; astilectron-bundler -v -c bundler_linux.json

build_windows: ## Build the binary for Windows
	cd src/; astilectron-bundler -v -c bundler_windows.json

run: build_linux ## Build and run the binary in GUI mode for Linux
	./bin/linux-amd64/'${APP_NAME}'

clean: ## Remove compiled binaries from bin/
	rm -Rf bin/
	rm src/windows.syso
	rm src/bind*

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

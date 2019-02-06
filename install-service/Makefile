#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

build: build_linux ## Build binaries for Windows and Linux
	make build_windows

build_linux: ## Build the binary for Linux
	# Build the standalone installer
	GOOS=linux GOARCH=amd64 go build -o bin/install-service ./src

build_windows: ## Build the binary for Windows
	# Embed the manifest to request admin rights using https://github.com/akavel/rsrc
	rsrc -manifest src/install-service.exe.manifest -o bin/manifest.syso
	# Build the standalone installer
	GOOS=windows GOARCH=amd64 go build -o bin/install-service.exe ./src
	rm bin/manifest.syso

clean: ## Remove compiled output
	rm -Rf bin/

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Name of the app
APP_NAME := 'MiningHQ\ Miner\ Manager'

build: ## Build the binary
	cd src/; astilectron-bundler -v

run: build ## Build and run the binary in GUI mode
	./bin/linux-amd64/'${APP_NAME}'

embed_linux: ## Embed the MininHQ miner service into this binary for building on Linux
	rm -Rf ../miner-service
	mkdir ../miner-service
	rm src/embedded/miner_service.go
	# Build the standalone installer
	go build -o ../install-service/bin/install-service ../install-service/src/main.go
	mv ../install-service/bin/install-service ../miner-service/install-service
	cp ${GOPATH}/src/github.com/donovansolms/mininghq-miner/bin/mininghq-miner ../miner-service/mininghq-miner
	esc -o src/embedded/miner_service.go -pkg embedded ../miner-service

embed_windows: ## Embed the MininHQ miner service into this binary for building on Windows
	rm -Rf ../miner-service
	mkdir ../miner-service
	rm src/embedded/miner_service.go
	# Embed the manifest to request admin rights using https://github.com/akavel/rsrc
	# rsrc -manifest install-service/install-service.exe.manifest -o install-service/manifest.syso
	# Build the standalone installer
	# GOOS=windows GOARCH=amd64 go build -o install-service/install-service.exe install-service/main.go
	mv install-service/install-service.exe miner-service/install-service.exe
	cp install-service/install-service.exe.manifest miner-service/install-service.exe.manifest
	cp ${GOPATH}/src/github.com/donovansolms/mininghq-miner/bin/mininghq-miner.exe miner-service/mininghq-miner.exe
	esc -o src/embedded/miner_service.go -pkg embedded miner-service

clean: ## Remove compiled binaries from bin/
	rm -Rf bin/
	rm src/windows.syso
	rm src/bind*

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

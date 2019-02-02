#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

# Name of the app
APP_NAME := 'MiningHQ\ Miner\ Manager'

build: ## Build the binary
	cd src/; astilectron-bundler -v

run: build ## Build and run the binary in CLI mode
	./bin/linux-amd64/'${APP_NAME}' --no-gui

run_cli:  ## Build and run the binary in CLI mode
	go run -race ./src/*.go --no-gui

run_gui: build ## Build and run the binary in GUI mode
	./bin/linux-amd64/'${APP_NAME}'

# default: build ## Build the binary
#
# build: ## Run the bundler to build the GUI
# 	cd src/; astilectron-bundler -v
#
# run: build ## Build and run the GUI for Linux
# 	./bin/linux-amd64/'${APP_NAME}'
#
# run_debug: build ## Build and run the GUI for Linux in debug mode
# 	./bin/linux-amd64/'${APP_NAME}' -d
#
# run_only_debug: ## Run the GUI for Linux without building
# 	./bin/linux-amd64/'${APP_NAME}' -d
#
embed_linux: ## Embed the MininHQ miner service into this binary for building on Linux
	rm -Rf miner-service
	mkdir miner-service
	# Build the standalone installer
	go build -o install-service/install-service install-service/main.go
	mv install-service/install-service miner-service/install-service
	cp ${GOPATH}/src/github.com/donovansolms/mininghq-miner/bin/mininghq-miner miner-service/mininghq-miner
	esc -o src/embedded/miner_service.go -pkg embedded miner-service

embed_windows: ## Embed the MininHQ miner service into this binary for building on Windows
	rm -Rf miner-service
	mkdir miner-service
	# Build the standalone installer
	GOOS=windows GOARCH=amd64 go build -o install-service/install-service.exe install-service/main.go
	mv install-service/install-service.exe miner-service/install-service.exe
	cp ${GOPATH}/src/github.com/donovansolms/mininghq-miner/bin/mininghq-miner.exe miner-service/mininghq-miner.exe
	esc -o src/embedded/miner_service.go -pkg embedded miner-service


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

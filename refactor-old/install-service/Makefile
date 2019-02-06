#
# A Makefile to build, run and test Go code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

build_windows: ## Build the binary
	# Embed the manifest to request admin rights using https://github.com/akavel/rsrc
	rsrc -manifest install-service.exe.manifest -o manifest.syso
	# Build the standalone installer
	GOOS=windows GOARCH=amd64 go build -o install-service.exe .

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// apiEndpoint is the MiningHQ API endpoint. Defined as a constant since
	// we don't ship any config files
	apiEndpoint = "http://mininghq.local/api/v1"
)

// AppName is injected by the Astilectron packager
var AppName string

// main is the main runnable of the application
func main() {

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Unable to get user home directory: %s\n", err)
	}

	if isInstalled() {
		// Installed, run manager
	}

	// Not installed, run installer
	// AppName, Asset and RestoreAssets are injected by the bundler
	gui, err := NewInstaller(
		AppName,
		Asset,
		RestoreAssets,
		homeDir,
		runtime.GOOS,
		apiEndpoint,
		false, // read from flag
	)
	if err != nil {
		// Setting the output to stdout so the user can see the error
		log.SetOutput(os.Stdout)
		log.Fatalf("Unable to set up miner: %s", err)
	}

	err = gui.Run()
	if err != nil {
		// Setting the output to stdout so the user can see the error
		log.SetOutput(os.Stdout)
		log.Fatalf("Unable to run miner: %s", err)
	}

}

// isInstalled checks if the Miner Manager has been installed already.
//
// The Miner Manager acts as both installer and manager. We need to decide
// which one to execute based on the installed services
func isInstalled() bool {
	// We created a file $USERHOME/.mhqpath containing the installation dir
	// If the path exists in .mhqpath then the services are installed
	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf(
			`Unable to find your home directory: %s\n

Please contact our support via our help channels listed at https://www.mininghq.io/help
`, err)
		os.Exit(1)
	}

	installedCheckfilePath := filepath.Join(homeDir, ".mhqpath")
	installedPath, err := ioutil.ReadFile(installedCheckfilePath)
	if err != nil {
		// if the mhqpath file doesn't exist, nothing is installed
		return false
	}

	// Check if installedPath exists
	info, err := os.Stat(string(installedPath))
	if err != nil {
		return false
	}

	// If it is a directory, then the services should be installed
	if info.IsDir() {
		return true
	}

	return false
}

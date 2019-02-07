// NOTE
//
// This app was initially the text based installer. We've moved to
// a GUI based installer since, but uninstall is still text based.
//
// We are evaluating the need for keeping the text-based installer for Linux
// servers.
//
// END NOTE

package main

import (
	"fmt"
	"io/ioutil"
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

// main is the main runnable of the application
func main() {

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Unable to get user home directory: %s\n", err)
	}

	mhqInstaller, err := NewInstaller(homeDir, runtime.GOOS, apiEndpoint)
	if err != nil {
		fmt.Printf("Unable to create installer: %s\n", err)
		return
	}

	if isInstalled() {
		fmt.Println("MiningHQ is already installed.")
		return

	}
	err = mhqInstaller.Install()
	if err != nil {
		fmt.Println("ERR", err)
	}
	return

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

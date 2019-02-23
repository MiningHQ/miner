/*
  MiningHQ Miner Manager - The MiningHQ Miner Manager GUI
  https://mininghq.io

	Copyright (C) 2018  Donovan Solms     <https://github.com/donovansolms>
                                        <https://github.com/mininghq>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mininghq/rpcproto/rpcproto"
	homedir "github.com/mitchellh/go-homedir"
	"google.golang.org/grpc"
)

const (
	// apiEndpoint is the MiningHQ API endpoint. Defined as a constant since
	// we don't ship any config files
	apiEndpoint = "https://www.mininghq.io/api/v1"
)

// AppName is injected by the Astilectron packager
var AppName string

// main is the main runnable of the application
func main() {

	debug := flag.Bool("d", false, "Enable debug mode")
	flag.Parse()

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Unable to get user home directory: %s\n", err)
	}

	if isInstalled() {
		// Installed, run manager
		conn, err := grpc.Dial("localhost:64630", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		client := rpcproto.NewManagerServiceClient(conn)

		// Start the Electron interface
		// AppName, Asset and RestoreAssets are injected by the bundler
		gui, err := NewManager(
			client,
			AppName,
			Asset,
			RestoreAssets,
			*debug,
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
		return
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
		*debug,
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
		fmt.Printf(`
Unable to find your home directory: %s\n

Please contact our support via our help channels listed at https://www.mininghq.io/connect
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

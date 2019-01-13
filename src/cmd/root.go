// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var noGUI bool
var mustUninstall bool
var apiEndpoint string
var debug bool

// Execute the main command
func Execute() {

	// The way we decide to run the installer or the manager doesn't work
	// cleanly with cobra's flags since we don't have a rootCmd. Maybe I'll
	// change this in future, but for now we'll just check if any uninstall
	// flag or command was passed
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if strings.Contains(strings.ToLower(arg), "uninstall") {
				mustUninstall = true
				break
			}
		}
	}

	if isInstalled() && mustUninstall == false{
		// Installed, run manager
		if err := manageCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if mustUninstall == false {
			fmt.Println("No installation found, installing MiningHQ")
		}
		if err := installCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
// TODO: REMOVE?


	// TODO: I don't think we need reading of config for the manager?
	//
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		homeDir, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".mhq")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

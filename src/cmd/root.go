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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var noGUI bool
var apiEndpoint string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if isInstalled() {
		fmt.Println("Installed already, run manager")
		if err := manageCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
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
	// TODO: Check homedir, if the user has certain things installed we can
	// assume miner manager has been installed already
	//
	// TODO: Can we use service.install?
	return false
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

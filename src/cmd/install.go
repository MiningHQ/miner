package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/donovansolms/mininghq-miner-manager/src/installer"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// AppName is injected by the Astilectron bundler
var AppName string

// Asset is injected by the Astilectron bundler
var Asset bootstrap.Asset

// RestoreAssets is injected by the Astilectron bundler
var RestoreAssets bootstrap.RestoreAssets

// installCmd represents a fresh installation
var installCmd = &cobra.Command{
	Use:   "MinerManager",
	Short: "The MiningHQ Miner Manager GUI",
	Long: `
   __  ____      _           __ ______
  /  |/  (_)__  (_)__  ___ _/ // / __ \
 / /|_/ / / _ \/ / _ \/ _ '/ _  / /_/ /
/_/  /_/_/_//_/_/_//_/\_, /_//_/\___\_\
                    /___/ Miner Manager


The MiningHQ Manager installs and configures the MiningHQ
services required for managing your rigs. It can be run as a GUI or command line.

Once installed, the Miner Manager can be used to view system and miner
stats on the local machine, however, the MiningHQ Dashboard
(https://www.mininghq.io/dashboard) is the best place to monitor miners from.`,
	Run: func(cmd *cobra.Command, args []string) {
		if noGUI {

			homeDir, err := homedir.Dir()
			if err != nil {
				fmt.Println("ERR", err)
			}

			installer, err := installer.New(homeDir, runtime.GOOS, apiEndpoint)
			if err != nil {
				fmt.Println("ERR", err)
				return
			}

			err = installer.InstallSync()
			if err != nil {
				fmt.Println("ERR", err)
				return
			}
			return
		}

		fmt.Println("Run as GUI")
		// If the '--no-gui' flag wasn't specified, we'll start the Electron
		// interface
		// AppName, Asset and RestoreAssets are injected by the bundler
		gui, err := installer.NewGUI(
			AppName,
			Asset,
			RestoreAssets,
			false, // TODO: Debug should come from somewhere else
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

	},
}

func init() {
	//installCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.mhq.yaml)")
	installCmd.Flags().BoolVar(&noGUI, "no-gui", false, "Run the manager without GUI")
	//installCmd.Flags().StringVar(&apiEndpoint, "api-endpoint", "http://mininghq.local/api/v1", "The base API endpoint for MiningHQ")
}

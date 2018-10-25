package cmd

import (
	"fmt"
	"runtime"

	"github.com/donovansolms/mininghq-miner-manager/src/installer"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

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
	},
}

func init() {
	installCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.mhq.yaml)")
	installCmd.Flags().BoolVar(&noGUI, "no-gui", false, "Run the manager without GUI")
	installCmd.Flags().StringVar(&apiEndpoint, "api-endpoint", "http://mininghq.local/api/v1", "The base API endpoint for MiningHQ")
}

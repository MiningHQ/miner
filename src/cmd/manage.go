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

package cmd

import (
	"log"
	"os"

	"github.com/donovansolms/mininghq-miner-manager/src/manager"
	"github.com/donovansolms/mininghq-rpcproto/rpcproto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// manageCmd represents running the GUI manager
var manageCmd = &cobra.Command{
	Use:   "MiningHQ Miner Manager",
	Short: "The MiningHQ Miner Manager GUI",
	Long: `
   __  ____      _           __ ______
  /  |/  (_)__  (_)__  ___ _/ // / __ \
 / /|_/ / / _ \/ / _ \/ _ '/ _  / /_/ /
/_/  /_/_/_//_/_/_//_/\_, /_//_/\___\_\
                    /___/ Miner Manager


The MiningHQ Manager shows you stats from your current rig and allows for
basic mining functionality.

The rig must be managed from your MiningHQ dashboard available at
https://www.mininghq.io`,
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := grpc.Dial("localhost:64630", grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		client := rpcproto.NewManagerServiceClient(conn)

		// Start the Electron interface
		// AppName, Asset and RestoreAssets are injected by the bundler
		gui, err := manager.NewGUI(
			client,
			AppName,
			Asset,
			RestoreAssets,
			debug,
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
	manageCmd.Flags().BoolVar(&noGUI, "no-gui", false, "Run the manager without GUI")
	manageCmd.Flags().BoolVar(&debug, "debug", false, "Run the manager in debug mode, a log file will be created")
	manageCmd.Flags().BoolVar(&mustUninstall, "uninstall", false, "Completely remove MiningHQ services from this system")
}

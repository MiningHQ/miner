package cmd

import (
	"fmt"
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

		if noGUI {
			// If we run with noGUI, we'll only print the stats once

			// homeDir, err := homedir.Dir()
			// if err != nil {
			// 	fmt.Println("ERR", err)
			// }
			//
			// installer, err := installer.New(homeDir, runtime.GOOS, apiEndpoint)
			// if err != nil {
			// 	fmt.Println("ERR", err)
			// 	return
			// }
			//
			// fmt.Println("Run manager!")
			//
			// conn, err := grpc.Dial("localhost:64630", grpc.WithInsecure())
			// if err != nil {
			// 	panic(err)
			// }
			// defer conn.Close()
			//
			// client := rpcproto.NewManagerServiceClient(conn)
			//
			// fmt.Println("Current state")
			// stateResponse, err := client.GetState(context.Background(), &rpcproto.StateRequest{})
			// if err != nil {
			// 	panic(err)
			// }
			//
			// fmt.Println(stateResponse.Status)
			// fmt.Println("currentState")
			// fmt.Println(stateResponse.State.String())

			//
			//
			//

			// fmt.Println("Start mining call")
			// //fmt.Println("Stop mining call")
			// response, err := client.SetState(context.Background(), &rpcproto.StateRequest{
			// 	//State: rpcproto.MinerState_StopMining,
			// 	State: rpcproto.MinerState_StartMining,
			// })
			// if err != nil {
			// 	panic(err)
			// }
			//
			// fmt.Println(response.Status)

			//
			// fmt.Println("GetLogs call")
			// response, err := client.GetLogs(context.Background(), &rpcproto.LogsRequest{
			// 	MaxLines: 100,
			// })
			// if err != nil {
			// 	panic(err)
			// }
			//
			// fmt.Println("Printing logs")
			// for _, logs := range response.MinerLogs {
			// 	fmt.Println(logs.Logs)
			// }
			//
			// fmt.Println("Getstats call")
			// statsResponse, err := client.GetStats(context.Background(), &rpcproto.StatsRequest{})
			// if err != nil {
			// 	panic(err)
			// }
			//
			// fmt.Println("Printing stats")
			// for _, stats := range statsResponse.Stats {
			// 	fmt.Println(stats.Hashrate)
			// }

			// err = installer.InstallSync()
			// if err != nil {
			// 	fmt.Println("ERR", err)
			// 	return
			// }
			return
		}

		fmt.Println("Run as GUI")

		// If the '--no-gui' flag wasn't specified, we'll start the Electron
		// interface
		// AppName, Asset and RestoreAssets are injected by the bundler
		gui, err := manager.NewGUI(
			client,
			AppName,
			Asset,
			RestoreAssets,
			true, // TODO: Debug should come from somewhere else
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
	manageCmd.Flags().BoolVar(&mustUninstall, "uninstall", false, "Completely remove MiningHQ services from this system")
}

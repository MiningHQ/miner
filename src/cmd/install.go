package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/donovansolms/mininghq-miner-manager/src/installer"
	"github.com/fatih/color"
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
		homeDir, err := homedir.Dir()
		if err != nil {
			fmt.Printf("Unable to get user home directory: %s\n", err)
		}

		installFiller := "install"
		if mustUninstall {
			installFiller = "uninstall"
		}

		currentUser, err := user.Current()
		if err != nil {
			fmt.Printf(`
We were unable to determine the current user on this system. The installer
requires Administrator (or sudo) rights to %s. Please run the installer
as the Administrator or with sudo.

If you are sure you have the permissions, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
	`, installFiller)
			fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		}

		// TODO: Find a better way to sudo install the service
		// 	TODO: PErhaps a standalone executable
		_ = currentUser
		// 		if currentUser.Uid != "0" {
		// 			if strings.ToLower(runtime.GOOS) == "linux" ||
		// 				strings.ToLower(runtime.GOOS) == "darwin" {
		// 				fmt.Printf(`
		// The installer requires sudo rights to %s the required MiningHQ services.
		// Please run the installer with sudo.
		//
		// If you are sure you have the permissions, please contact support to resolve
		// the issue. Support can be contacted via our help channels listed at
		// https://www.mininghq.io/help
		//
		// `, installFiller)
		// 			} else {
		// 				fmt.Printf(`
		// The installer requires Administrator rights to %s the required MiningHQ services.
		// Please run the installer as the Administrator.
		//
		// To do that, right click on the installer and select 'Run as Administrator'.
		//
		// If you are sure you have the permissions, please contact support to resolve
		// the issue. Support can be contacted via our help channels listed at
		// https://www.mininghq.io/help
		//
		// `, installFiller)
		// 			}
		// 			os.Exit(0)
		//
		// 		}

		mhqInstaller, err := installer.New(homeDir, runtime.GOOS, apiEndpoint)
		if err != nil {
			fmt.Printf("Unable to create installer: %s\n", err)
			return
		}

		if mustUninstall {

			// Get the current installed path
			installedCheckfilePath := filepath.Join(homeDir, ".mhqpath")
			installedPath, err := ioutil.ReadFile(installedCheckfilePath)
			if err != nil {
				fmt.Printf(`
We were unable to find the installed location for the MiningHQ services. Please
remove the files manually where you installed the services.
				`)
				fmt.Println()
				os.Exit(0)
			}

			// We're not checking for noGUI here since --uninstall is a command-line
			// only operation. GUI uninstall is triggered from the Miner Manager GUI
			err = mhqInstaller.UninstallSync(strings.TrimSpace(string(installedPath)), installedCheckfilePath)
			if err != nil {
				fmt.Println("ERR", err)
				return
			}

			os.Exit(0)
		}

		if noGUI {
			err = mhqInstaller.InstallSync()
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
	installCmd.Flags().BoolVar(&noGUI, "no-gui", false, "Run the manager without GUI")
	installCmd.Flags().StringVar(&apiEndpoint, "api-endpoint", "http://mininghq.local/api/v1", "The base API endpoint for MiningHQ")
	installCmd.Flags().BoolVar(&mustUninstall, "uninstall", false, "Completely remove MiningHQ services from this system")
}

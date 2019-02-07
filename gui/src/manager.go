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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/buildkite/terminal"
	"github.com/donovansolms/mininghq-rpcproto/rpcproto"
	"github.com/sirupsen/logrus"
)

// GUIManager implements the manager GUI
type GUIManager struct {
	// window is the main Astilectron window
	window *astilectron.Window
	// astilectronOptions holds the Astilectron options
	astilectronOptions bootstrap.Options
	// managerClient is the client to the miner controller's manager API
	managerClient rpcproto.ManagerServiceClient
	// logger logs to stdout
	logger   *logrus.Entry
	debugLog *os.File
}

// NewManager creates a new instance of the graphical installer
func NewManager(
	client rpcproto.ManagerServiceClient,
	appName string,
	asset bootstrap.Asset,
	restoreAssets bootstrap.RestoreAssets,
	isDebug bool) (*GUIManager, error) {

	gui := GUIManager{
		managerClient: client,
	}

	// If no config is specified then this is the first run
	startPage := "manager.html"

	var menu []*astilectron.MenuItemOptions

	// Setup the logging, by default we log to stdout
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetOutput(os.Stdout)

	// Create the window options
	windowOptions := astilectron.WindowOptions{
		// If frame is false, the window frame is removed. If isDebug is true,
		// we show the frame to have debugging options available
		Frame:           astilectron.PtrBool(isDebug),
		BackgroundColor: astilectron.PtrStr("#0B0C22"),
		Center:          astilectron.PtrBool(true),
		Height:          astilectron.PtrInt(500),
		Width:           astilectron.PtrInt(980),
	}

	if isDebug {
		// Get current path
		debugLogPath := filepath.Join(os.TempDir(), "mininghq_debug.log")
		executable, err := os.Executable()
		if err == nil {
			debugLogPath = filepath.Join(filepath.Dir(executable), "mininghq_debug.log")
		}

		gui.debugLog, err = os.OpenFile(
			debugLogPath,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0644)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(gui.debugLog)

		// We only show the menu bar in debug mode
		menu = append(menu, &astilectron.MenuItemOptions{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Role: astilectron.MenuItemRoleClose,
				},
			},
		})
	}
	// To make copy and paste work on Mac, the copy and paste entries need to
	// be defined, the alternative is to implement the clipboard API
	// https://github.com/electron/electron/blob/master/docs/api/clipboard.md
	if runtime.GOOS == "darwin" {
		menu = append(menu, &astilectron.MenuItemOptions{
			Label: astilectron.PtrStr("Edit"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Role: astilectron.MenuItemRoleCut,
				},
				{
					Role: astilectron.MenuItemRoleCopy,
				},
				{
					Role: astilectron.MenuItemRolePaste,
				},
				{
					Role: astilectron.MenuItemRoleSelectAll,
				},
			},
		})

		windowOptions.Frame = astilectron.PtrBool(isDebug)
		windowOptions.TitleBarStyle = astilectron.PtrStr("hidden")
	}

	// Setting the WithFields now will ensure all log entries from this point
	// includes the fields
	gui.logger = logrus.WithFields(logrus.Fields{
		"service": "mininghq-manager",
	})

	gui.astilectronOptions = bootstrap.Options{
		Debug:         isDebug,
		Asset:         asset,
		RestoreAssets: restoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       startPage,
			MessageHandler: gui.handleElectronCommands,
			Options:        &windowOptions,
		}},
		AstilectronOptions: astilectron.Options{
			AppName:            appName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		// TODO: Fix this tray to display nicely
		/*TrayOptions: &astilectron.TrayOptions{
			Image:   astilectron.PtrStr("/static/i/miner-logo.png"),
			Tooltip: astilectron.PtrStr(appName),
		},*/
		MenuOptions: menu,
		// OnWait is triggered as soon as the electron window is ready and running
		OnWait: func(
			_ *astilectron.Astilectron,
			windows []*astilectron.Window,
			_ *astilectron.Menu,
			_ *astilectron.Tray,
			_ *astilectron.Menu) error {
			gui.window = windows[0]
			return nil
		},
	}

	gui.logger.Info("Setup complete")
	return &gui, nil
}

// Run the miner!
func (gui *GUIManager) Run() error {
	gui.logger.Info("Starting manager")

	err := bootstrap.Run(gui.astilectronOptions)
	if err != nil {
		return err
	}
	gui.debugLog.Close()
	return nil
}

// updateLoop is executed every X seconds, it fetches the latest state, stats
// and logs from the miner controller and sends it to the Electron.
func (gui *GUIManager) updateLoop() {

	var managerUpdate rpcproto.ManagerUpdate

	for {
		gui.logger.Debug("Fetching update information")

		managerUpdate = rpcproto.ManagerUpdate{
			Stats: &rpcproto.MinerStats{},
		}

		// Get the miner's stats
		statsResponse, err := gui.managerClient.GetStats(context.Background(), &rpcproto.StatsRequest{})
		if err != nil {
			gui.logger.WithField(
				"op", "GetStats",
			).Errorf("Unable to get stats from controller: %s", err)
		} else if statsResponse != nil {
			// Combine all the miner stats into one
			for _, stats := range statsResponse.Stats {
				managerUpdate.Stats.Hashrate += stats.Hashrate
				managerUpdate.Stats.TotalShares += stats.TotalShares
				managerUpdate.Stats.AcceptedShares += stats.AcceptedShares
				managerUpdate.Stats.RejectedShares += stats.RejectedShares
			}
		}

		// Get the miner's state
		stateResponse, err := gui.managerClient.GetState(context.Background(), &rpcproto.StateRequest{})
		if err != nil {
			gui.logger.WithField(
				"op", "GetState",
			).Errorf("Unable to get state from controller: %s", err)
		} else if stateResponse != nil {
			managerUpdate.State = stateResponse.State
		}

		// Get the miner's logs
		logsResponse, err := gui.managerClient.GetLogs(context.Background(), &rpcproto.LogsRequest{
			MaxLines: 500,
		})
		if err != nil {
			gui.logger.WithField(
				"op", "GetLogs",
			).Errorf("Unable to get logs from controller: %s", err)
		} else if logsResponse != nil {
			var logs []string
			// We need to format the logs to HTML for display
			for _, log := range logsResponse.MinerLogs {
				for _, line := range log.Logs {
					htmlLine := terminal.Render([]byte(line))
					logs = append(logs, string(htmlLine))
				}
			}
			managerUpdate.HTMLLogs = strings.Join(logs, "<br/>")
		}

		err = gui.sendElectronCommand("update", managerUpdate)
		if err != nil {
			gui.logger.WithField(
				"method", "update",
			).Errorf("Unable to send update to Electron: %s", err)
		}

		time.Sleep(time.Second * 5)
	}
}

// handleElectronCommands handles the messages sent by the Electron front-end
func (gui *GUIManager) handleElectronCommands(
	_ *astilectron.Window,
	command bootstrap.MessageIn) (interface{}, error) {

	gui.logger.WithField(
		"method", command.Name,
	).Debug("Received command from Electron")

	// Every Electron command has a name together with a payload containing the
	// actual message
	switch command.Name {

	case "ready":
		// Check if the miner controller is available by fetching the
		// base information

		response, err := gui.managerClient.GetInfo(
			context.Background(),
			&rpcproto.RigInfoRequest{})
		if err != nil {
			gui.logger.WithField(
				"method", "setup",
			).Errorf("Unable to query miner controller: %s", err)

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
Unable to query miner controller, please ensure the MiningHQ Miner service is running.
</p>
<p>
%s
</p>`, err),
			}, nil
		}

		// Send the initial setup packet
		// This includes a link to the rig on MiningHQ as well as the
		// rig name. This is injected into the frontend for display purposes
		err = gui.sendElectronCommand("setup", map[string]string{
			"name": response.Name,
			"link": response.Link,
		})
		if err != nil {
			gui.logger.WithField(
				"method", "setup",
			).Errorf("Unable to send setup to Electron: %s", err)
		}

		// TODO: Do this loop better
		go gui.updateLoop()

	case "Cancel":

	}
	return nil, fmt.Errorf("'%s' is an unknown command", command.Name)
}

// sendElectronCommand sends the given data to Electron under the command name
func (gui *GUIManager) sendElectronCommand(
	name string,
	data interface{}) error {
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return bootstrap.SendMessage(gui.window, name, string(dataBytes))
}

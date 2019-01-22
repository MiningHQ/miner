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

package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/donovansolms/mininghq-miner-controller/src/mhq"
	"github.com/donovansolms/mininghq-miner-manager/src/embedded"
	"github.com/donovansolms/mininghq-spec/spec/caps"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

// GUIInstaller implements a graphical installer
type GUIInstaller struct {
	// window is the main Astilectron window
	window *astilectron.Window
	// astilectronOptions holds the Astilectron options
	astilectronOptions bootstrap.Options

	// homeDir is the user's home directory
	homeDir string
	// os is the system operating system
	os string
	// mhqEndpoint is the MiningHQ API endpoint to use
	mhqEndpoint string

	serviceName        string
	serviceDisplayName string
	serviceDescription string

	// logger logs to stdout
	logger *logrus.Entry

	// helper functions
	helper   Helper
	debugLog *os.File

	// Rig related information
	rigName     string
	installPath string
}

// NewGUI creates a new instance of the graphical installer
func NewGUI(
	appName string,
	asset bootstrap.Asset,
	restoreAssets bootstrap.RestoreAssets,
	homeDir string,
	systemOS string,
	apiEndpoint string,
	isDebug bool) (*GUIInstaller, error) {

	gui := GUIInstaller{
		serviceName:        serviceName,
		serviceDisplayName: serviceDisplayName,
		serviceDescription: serviceDescription,
		helper:             Helper{},
		homeDir:            homeDir,
		os:                 systemOS,
		mhqEndpoint:        apiEndpoint,
	}

	// If no config is specified then this is the first run
	startPage := "installer.html"

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
		Width:           astilectron.PtrInt(900),
		Height:          astilectron.PtrInt(600),
	}

	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)

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
		"service": "mininghq-installer",
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
func (gui *GUIInstaller) Run() error {
	gui.logger.Info("Starting installer")
	err := bootstrap.Run(gui.astilectronOptions)
	if err != nil {
		return err
	}
	gui.debugLog.Close()
	return nil
}

// handleElectronCommands handles the messages sent by the Electron front-end
func (gui *GUIInstaller) handleElectronCommands(
	_ *astilectron.Window,
	command bootstrap.MessageIn) (interface{}, error) {

	gui.logger.WithField(
		"method", command.Name,
	).Debug("Received command from Electron")

	// Every Electron command has a name together with a payload containing the
	// actual message
	switch command.Name {

	case "get-defaults":
		return map[string]string{
			"status":  "ok",
			"message": filepath.Join(gui.homeDir, "MiningHQ"),
		}, nil

	case "install":

		var payload map[string]string
		err := json.Unmarshal(command.Payload, &payload)
		if err != nil {
			// TODO: Send error back to electron
			return nil, err
		}

		if _, ok := payload["rigName"]; !ok {
			return map[string]string{
				"status":  "error",
				"message": "A Rig Name must be set to install",
			}, nil
		}

		if _, ok := payload["installPath"]; !ok {
			return map[string]string{
				"status":  "error",
				"message": "The install path must be set to install",
			}, nil
		}

		gui.rigName = strings.TrimSpace(payload["rigName"])
		gui.installPath = strings.TrimSpace(payload["installPath"])

		// TODO: Send message to electron we're installing

		avExcludeDirectory, err := gui.helper.CreateInstallDirectories(gui.installPath)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
We could not create one or more of the installation directories. Please
ensure you have sufficient permissions (like Administrator or root) access
to create directories in '%s'.

Include the following error in your report: %s
`, gui.installPath, err.Error()),
			}, nil
		}
		// Return the exclude directory for antivirus
		// and then wait for the confirmation to be sent to us to continue
		return map[string]string{
			"status":  "confirm-av",
			"message": avExcludeDirectory,
		}, nil

	// Sent after the user confirmed the exclude of the miner path, we can
	// continue with the install
	case "confirmed-av":

		// We need to know about the base system specs
		systemInfo, err := caps.GetSystemInfo()
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to determine the capabilities of this rig. Please ensure you
have sufficient permissions to check installed hardware on this system.
</p>
<p>
If you are sure you have the permissions, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
</p>
<p>
Include the following error in your report '%s'"), err.Error())
</p>`, err.Error()),
			}, nil
		}

		_ = gui.sendElectronCommand("install_progress", map[string]string{
			"status":  "ok",
			"message": "Gather rig capabilities",
		})

		miningKeyPath := "mining_key"
		apiCreateError := fmt.Sprintf(`
We were unable to connect to the MiningHQ API to register your rig.
Please check that the file '%s' is present in the same directory you are
running the installer from. If not, please download the Miner Manager again
from <a href="https://www.mininghq.io/rigs">https://www.mininghq.io/rigs</a>
		`,
			miningKeyPath)

		// Get the mining key for the user
		miningKey, err := gui.helper.GetMiningKeyFromFile(miningKeyPath)
		if err != nil {
			return map[string]string{
				"status":  "error",
				"message": fmt.Sprintf("%s<br/>Include the following error in your report: %s", apiCreateError, err.Error()),
			}, nil
		}

		apiClient, err := mhq.NewClient(miningKey, gui.mhqEndpoint)
		if err != nil {
			return map[string]string{
				"status":  "error",
				"message": fmt.Sprintf("%s<br/>Include the following error in your report: %s", apiCreateError, err.Error()),
			}, nil
		}

		registerRequest := mhq.RegisterRigRequest{
			Name: gui.rigName,
			Caps: systemInfo,
		}
		rigID, err := apiClient.RegisterRig(registerRequest)
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to register your rig with MiningHQ. Please ensure that
you are connected to the internet and that the file '%s' contains the same
mining key that you can find under 'Mining' in your settings available at
https://www.mininghq.io/user/settings
</p>
<p>
If you are sure everything is in order, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
<a href="https://www.mininghq.io/help">https://www.mininghq.io/help</a>
</p>
<p>
Include the following error in your report '%s'
</p>
				`, miningKeyPath, err.Error()),
			}, nil
		}

		_ = gui.sendElectronCommand("install_progress", map[string]string{
			"status":  "ok",
			"message": "Register rig with MiningHQ",
		})

		// Creating files and copying
		err = copy.Copy(
			miningKeyPath,
			filepath.Join(gui.installPath, "miner-controller", filepath.Base(miningKeyPath)))
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to copy your mining key to your installation.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		err = ioutil.WriteFile(
			filepath.Join(gui.installPath, "miner-controller", "rig_id"),
			[]byte(rigID),
			0644)
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to create the new rig files for your installation.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		_ = gui.sendElectronCommand("install_progress", map[string]string{
			"status":  "ok",
			"message": "Create config files",
		})

		// Install the miner and service
		embeddedFilename := "mininghq-miner"
		embeddedServiceInstallerFilename := "install-service"
		if strings.ToLower(runtime.GOOS) == "windows" {
			embeddedFilename = "mininghq-miner.exe"
			embeddedServiceInstallerFilename = "install-service.exe"
		}
		embeddedFS := embedded.FS(false)
		// Extract the miner service
		embeddedFile, err := embeddedFS.Open("/miner-service/" + embeddedFilename)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to extract the miner from the installer.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		installFile, err := os.OpenFile(
			filepath.Join(gui.installPath, embeddedFilename),
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to create the miner in the correct location. Please check that
you have sufficient space on your harddrive.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		_, err = io.Copy(installFile, embeddedFile)
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to install the miner to the correct location.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}
		installFile.Close()
		embeddedFile.Close()

		// Extract the service installer
		embeddedFile, err = embeddedFS.Open("/miner-service/" + embeddedServiceInstallerFilename)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to extract the service installer from the installer.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		installFile, err = os.OpenFile(
			filepath.Join(gui.installPath, embeddedServiceInstallerFilename),
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to create the service installer in the correct location. Please check that
you have sufficient space on your harddrive.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		_, err = io.Copy(installFile, embeddedFile)
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to install the service installer to the correct location.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}
		installFile.Close()
		embeddedFile.Close()

		// Install mininghq-miner as a service
		// We do this using a separate executable so that only the service install
		// requires Administrator/sudo rights and not the entire installer
		out, err := exec.Command(
			"pkexec",
			filepath.Join(gui.installPath, embeddedServiceInstallerFilename),
			"-op", "install",
			"-serviceName", gui.serviceName,
			"-serviceDisplayName", gui.serviceDisplayName,
			"-serviceDescription", gui.serviceDescription,
			"-installedPath", gui.installPath,
			"-serviceFilename", embeddedFilename,
		).CombinedOutput()
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to install the miner service.
</p>
<p>
Include the following error in your report '%s', %s
</p>
				`, err.Error(), out),
			}, nil
		}

		installedCheckfilePath := filepath.Join(gui.homeDir, ".mhqpath")
		installedCheckfile, err := os.OpenFile(
			installedCheckfilePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0644)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to create the installer check file in ~/.mhqpath. This
will cause MiningHQ services to be unable to detect the installation.
</p>
<p>
Please ensure you have the correct permissions to write to your home directory.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}
		defer installedCheckfile.Close()

		_, err = installedCheckfile.WriteString(gui.installPath)
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to write to the installer check file in ~/.mhqpath. This
will cause MiningHQ services to be unable to detect the installation.
</p>
<p>
Please ensure you have the correct permissions to write to your home directory.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		// Copy the manager
		managerBinaryPath, err := os.Executable()
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to copy the miner manager to your installation path.
</p>
<p>
Please ensure you have the correct permissions to write to your install directory.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		managerName := filepath.Base(managerBinaryPath)
		err = os.Rename(managerBinaryPath, filepath.Join(gui.installPath, managerName))
		if err != nil {

			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to copy the miner manager to your installation path.
</p>
<p>
Please ensure you have the correct permissions to write to your install directory.
</p>
<p>
Include the following error in your report '%s'
</p>
				`, err.Error()),
			}, nil
		}

		_ = gui.sendElectronCommand("install_progress", map[string]string{
			"status":  "ok",
			"message": "Installing MiningHQ Miner",
		})

		// Start the mininghq-miner service
		// We do this using a separate executable so that only the service install
		// requires Administrator/sudo rights and not the entire installer
		out, err = exec.Command(
			"pkexec",
			filepath.Join(gui.installPath, embeddedServiceInstallerFilename),
			"-op", "start",
			"-serviceName", gui.serviceName,
			"-serviceDisplayName", gui.serviceDisplayName,
			"-serviceDescription", gui.serviceDescription,
			"-installedPath", gui.installPath,
			"-serviceFilename", embeddedFilename,
		).CombinedOutput()
		if err != nil {
			return map[string]string{
				"status": "error",
				"message": fmt.Sprintf(`
<p>
We were unable to copy the miner manager to your installation path.
</p>
<p>
Please ensure you have the correct permissions to write to your install directory.
</p>
<p>
Include the following error in your report '%s', %s
</p>
				`, err.Error(), out),
			}, nil
		}

		return map[string]string{
			"status":  "ok",
			"message": "",
		}, nil

	// Firstrun is received on the first run of the miner. We return the current
	// logged in username
	case "setup":
		var username string
		currentUser, err := user.Current()
		if err == nil {
			if currentUser.Name != "" {
				username = currentUser.Name
			} else if currentUser.Username != "" {
				username = currentUser.Username
			}
		}
		return username, nil

		// Cancel the installation
	case "Cancel":
		// err := gui.stopMiner()
		// if err != nil {
		// 	// _ = gui.sendElectronCommand("fatal_error", ElectronMessage{
		// 	// 	Data: fmt.Sprintf("Unable to stop miner backend."+
		// 	// 		"Please close the miner and open it again."+
		// 	// 		"<br/>The error was '%s'", err),
		// 	// })
		// 	// // Give the UI some time to display the message
		// 	// time.Sleep(time.Second * 15)
		// 	// gui.logger.Fatalf("Unable to reconfigure miner: '%s'", err)
		// }
	}
	return nil, fmt.Errorf("'%s' is an unknown command", command.Name)
}

// sendElectronCommand sends the given data to Electron under the command name
func (gui *GUIInstaller) sendElectronCommand(
	name string,
	data map[string]string) error {
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return bootstrap.SendMessage(gui.window, name, string(dataBytes))
}

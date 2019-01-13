/*
  MiningHQ Miner Manager - The MiningHQ Miner Manager GUI and installer
  https://mininghq.io

  Copyright (C) 2018  Donovan Solms     <https://github.com/donovansolms>

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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/donovansolms/mininghq-miner-controller/src/mhq"
	"github.com/donovansolms/mininghq-miner-manager/src/embedded"
	"github.com/donovansolms/mininghq-spec/spec/caps"
	"github.com/fatih/color"
	"github.com/otiai10/copy"
	input "github.com/tcnksm/go-input"
)

// CLIInstaller install the Miner Manager from the terminal
type CLIInstaller struct {
	// homeDir is the user's home directory
	homeDir string
	// os is the system operating system
	os string
	// mhqEndpoint is the MiningHQ API endpoint to use
	mhqEndpoint string

	serviceName        string
	serviceDisplayName string
	serviceDescription string

	// helper continas helper install functions
	helper Helper
}

// New creates a new installer instance
func New(homeDir string, os string, mhqEndpoint string) (*CLIInstaller, error) {
	if strings.TrimSpace(homeDir) == "" {
		return nil, errors.New("A home directory must be set")
	}

	os = strings.ToLower(os)
	if strings.TrimSpace(os) != Windows && strings.TrimSpace(os) != MacOS &&
		strings.TrimSpace(os) != Linux {
		return nil, fmt.Errorf("OS may only be %s, %s or %s", Windows, MacOS, Linux)
	}

	installer := CLIInstaller{
		homeDir:            homeDir,
		os:                 os,
		mhqEndpoint:        mhqEndpoint,
		serviceName:        "GoServiceExampleLogging",
		serviceDisplayName: "Go Service Example for Logging",
		serviceDescription: "This is an example Go service that outputs log messages.",
		helper:             Helper{},
	}
	return &installer, nil
}

// UninstallSync uninstalls the miner manager and services using
// a synchronous process
func (installer *CLIInstaller) UninstallSync(
	installedPath string,
	installedPathFilepath string) error {

	// Note: This will not be the prettiest code you'll ever see :)
	// If anyone has some good advice in controlling the output for this process,
	// feel free to let me know

	fmt.Printf(`
    __  ____      _           __ ______
   /  |/  (_)__  (_)__  ___ _/ // / __ \
  / /|_/ / / _ \/ / _ \/ _ '/ _  / /_/ /
 /_/  /_/_/_//_/_/_//_/\_, /_//_/\___\_\
                     /___/ Miner Uninstaller
                           www.mininghq.io

This will remove the MiningHQ Miner Manager and all related services.
We detected the installation in '%s'

`, installedPath)

	ui := &input.UI{}
	question := "\nAre you sure you wish to remove the MiningHQ Miner and all MiningHQ services? [Y/yes/N/no]"
	response, _ := ui.Ask(question, &input.Options{
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			validConfirmations := map[string]bool{
				"y":   true,
				"yes": true,
				"n":   true,
				"no":  true,
			}
			answer := strings.ToLower(s)
			if _, ok := validConfirmations[answer]; !ok {
				return fmt.Errorf(
					"Answer '%s' is invalid. Must be 'y', 'yes', 'n' or 'no'", s)
			}
			return nil
		},
	})
	allowContinue := strings.ToLower(response)
	if allowContinue == "n" || allowContinue == "no" {
		color.HiRed("********************************")
		color.HiRed("* Uninstall has been cancelled *")
		color.HiRed("********************************")
		color.HiYellow(`
Something wrong? If so, please let us know by getting in contact
via our help channels listed at https://www.mininghq.io/help
`)
		os.Exit(0)
	}

	// Remove the service
	fmt.Print("Deregister rig\t\t\t\t")
	miningKeyPath := filepath.Join(installedPath, "miner-controller", "mining_key")
	rigIDPath := filepath.Join(installedPath, "miner-controller", "rig_id")
	apiCreateError := fmt.Sprintf(`
We were unable to connect to the MiningHQ API to deregister your rig.
Please check that the file '%s' and '%s' is present in your installation directory.
`,
		miningKeyPath,
		rigIDPath)
	miningKeyBytes, miningKeyErr := ioutil.ReadFile(miningKeyPath)
	rigIDBytes, rigIDErr := ioutil.ReadFile(rigIDPath)
	if miningKeyErr != nil || rigIDErr != nil {
		color.HiRed("FAIL")
		fmt.Println(apiCreateError)
		if miningKeyErr != nil {
			fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), miningKeyErr.Error())
		} else if rigIDErr != nil {
			fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), rigIDErr.Error())
		}
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	miningKey := strings.TrimSpace(string(miningKeyBytes))
	rigID := strings.TrimSpace(string(rigIDBytes))
	apiClient, err := mhq.NewClient(miningKey, installer.mhqEndpoint)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Println(apiCreateError)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	err = apiClient.DeregisterRig(mhq.DeregisterRigRequest{
		RigID: rigID,
	})
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to deregister your rig with MiningHQ. Please ensure that
you are connected to the internet and that the file '%s' contains the same
mining key that you can find under 'Mining' in your settings available at
https://www.mininghq.io/user/settings

If you are sure everything is in order, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
`,
			miningKeyPath)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()

		// If we can't deregister the rig, continue with the rest of the removal
		// anyways
	} else {
		// Rig removed
		color.HiGreen("OK")
	}

	// Remove the service
	fmt.Print("Removing the MiningHQ Miner service\n")

	serviceFilename := "mininghq-miner"
	if strings.ToLower(runtime.GOOS) == "windows" {
		serviceFilename = "mininghq-miner.exe"
	}

	// Uninstall mininghq-miner as a service
	// We do this using a separate executable so that only the service uninstall
	// requires Administrator/sudo rights and not the entire installer
	out, err := exec.Command(
		"sudo",
		"/home/donovan/Development/Go/code/src/github.com/donovansolms/mininghq-miner-manager/install-service/install-service",
		"-op", "uninstall",
		"-serviceName", installer.serviceName,
		"-serviceDisplayName", installer.serviceDisplayName,
		"-serviceDescription", installer.serviceDescription,
		"-installedPath", installedPath,
		"-serviceFilename", serviceFilename,
	).CombinedOutput()
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to uninstall the miner service (it might already be uninstalled).
`)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s': %s"), err.Error(), out)
		fmt.Println()
		fmt.Println()
		color.Unset()

		// If we can't remove the service, continue with the rest of the removal
		// anyways
	} else {
		// Service uninstalled
		color.HiGreen("OK")
	}

	// Remove files
	fmt.Print("Remove the files\t\t\t")
	err = os.RemoveAll(installedPath)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to remove the MiningHQ files from '%s'.
`, installedPath)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(0)
	}
	err = os.Remove(installedPathFilepath)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to remove the MiningHQ file from '%s'.
`, installedPathFilepath)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(0)
	}
	// Files removed
	color.HiGreen("OK")

	fmt.Printf(`


***************************
*  MiningHQ uninstalled!  *
***************************

The MiningHQ Miner Manager and related services have been uninstalled. If you
wish to add this rig back, visit the rigs page and click 'add rig'

https://www.mininghq.io/rigs

Please join the MiningHQ community on Discord, Twitter and elsewhere, you can find
all our channels at https://www.mininghq.io/connect

We hope we see you again,
The MiningHQ Team
	`)

	fmt.Println()
	fmt.Println()

	return nil
}

// InstallSync installs the miner manager using a synchronous process,
// no feedback is given to the caller via channels
func (installer *CLIInstaller) InstallSync() error {

	// Note: This will not be the prettiest code you'll ever see :)
	// If anyone has some good advice in controlling the output for this process,
	// feel free to let me know

	fmt.Println(`
    __  ____      _           __ ______
   /  |/  (_)__  (_)__  ___ _/ // / __ \
  / /|_/ / / _ \/ / _ \/ _ '/ _  / /_/ /
 /_/  /_/_/_//_/_/_//_/\_, /_//_/\___\_\
                     /___/ Miner Installer
                           www.mininghq.io

This installer will install the MiningHQ Miner Manager.
The Miner Manager connects to your MiningHQ account to allow you to control
all your mining rigs easily.

The setup will guide you through the steps now,
the installation will take less than 5 minutes`)

	ui := &input.UI{}
	installDir := filepath.Join(installer.homeDir, "MiningHQ")
	rigName := "My first rig"

	hostName, err := os.Hostname()
	if err == nil {
		rigName = fmt.Sprintf("%s Rig", hostName)
	}
	rigName = strings.Title(strings.ToLower(rigName))

	question := "\nWhere should the services be installed to? "
	response, _ := ui.Ask(question, &input.Options{
		Default:  installDir,
		Required: true,
		Loop:     true,
	})
	installDir = response

	question = "\nWhat would you like to name this rig? "
	response, _ = ui.Ask(question, &input.Options{
		Default:  rigName,
		Required: true,
		Loop:     true,
	})
	rigName = response

	color.Yellow(`
The MiningHQ Miner Manager will now download and install the
required services, please verify the installation details below
`)

	fmt.Println()
	fmt.Printf("Installation directory: '%s'\n", installDir)
	fmt.Printf("Rig name: \t\t'%s'\n", rigName)
	fmt.Println()

	question = "\nDo you want to continue? [Y/yes/N/no]"
	response, _ = ui.Ask(question, &input.Options{
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			validConfirmations := map[string]bool{
				"y":   true,
				"yes": true,
				"n":   true,
				"no":  true,
			}
			answer := strings.ToLower(s)
			if _, ok := validConfirmations[answer]; !ok {
				return fmt.Errorf(
					"Answer '%s' is invalid. Must be 'y', 'yes', 'n' or 'no'", s)
			}
			return nil
		},
	})
	allowContinue := strings.ToLower(response)
	if allowContinue == "n" || allowContinue == "no" {
		color.HiRed("***************************************")
		color.HiRed("* The installation has been cancelled *")
		color.HiRed("***************************************")
		color.HiYellow(`
Something wrong? If so, please let us know by getting in contact
via our help channels listed at https://www.mininghq.io/help
`)
		os.Exit(0)
	}

	// Create the installation directory
	fmt.Print("Creating installation directory\t\t")
	avExcludeDirectory, err := installer.helper.CreateInstallDirectories(installDir)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We could not create one or more of the installation directories. Please
ensure you have sufficient permissions (like Administrator or root) access
to create directories in '%s'.
`,
			installDir)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	// Installation directory created
	color.HiGreen("OK")

	blinking := color.New(color.BlinkSlow, color.FgHiYellow)
	fmt.Println()
	fmt.Println("***************************************")
	fmt.Println("*   A note about antivirus software   *")
	fmt.Print("*")
	blinking.Print("         action required             ")
	fmt.Print("*\n")
	fmt.Println("***************************************")
	fmt.Printf(`
Cryptocurrency miners are detected by most antivirus
software and removed even though they don't contain any
viruses. To use this computer to mine cryptocurrencies you
need to exclude the miner installation path from being scanned.
To protect you, the MiningHQ Mining Manager will regularly scan
the miners to verify that they have not been tampered with and
only contain our official releases.


The directory to exclude is: '%s'


You can follow the following instructions on how to exclude the directory: %s
`,
		color.HiYellowString(avExcludeDirectory),
		installer.helper.GetOSAVGuides(),
	)

	fmt.Println()
	bold := color.New(color.Bold, color.Underline)
	bold.Println("Please exclude the directory from your antivirus now")
	question = "\nHave you excluded the directory? [Y/yes/N/no]"
	response, _ = ui.Ask(question, &input.Options{
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			validConfirmations := map[string]bool{
				"y":   true,
				"yes": true,
				"n":   true,
				"no":  true,
			}
			answer := strings.ToLower(s)
			if _, ok := validConfirmations[answer]; !ok {
				return fmt.Errorf(
					"Answer '%s' is invalid. Must be 'y', 'yes', 'n' or 'no'", s)
			}
			return nil
		},
	})
	allowContinue = strings.ToLower(response)
	if allowContinue == "n" || allowContinue == "no" {
		color.HiRed("****************************************")
		color.HiRed("* You must exclude the miner directory *")
		color.HiRed("****************************************")
		color.HiYellow(`
Something wrong? If so, please let us know by getting in contact
via our help channels listed at https://www.mininghq.io/help
`)
		os.Exit(0)
	}

	fmt.Print("Gather rig capabilities\t\t\t")
	systemInfo, err := caps.GetSystemInfo()
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to determine the capabilities of this rig. Please ensure you
have sufficient permissions to check installed hardware on this system.

If you are sure you have the permissions, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
`)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	// Installation directory created
	color.HiGreen("OK")

	// Register this rig with MiningHQ
	miningKeyPath := "mining_key"
	fmt.Print("Register rig with MiningHQ\t\t")
	apiCreateError := fmt.Sprintf(`
We were unable to connect to the MiningHQ API to register your rig.
Please check that the file '%s' is present in the same directory you are
running the installer from. If not, please download the Miner Manager again
from https://www.mininghq.io/rigs
`,
		miningKeyPath)

	miningKey, err := installer.helper.GetMiningKeyFromFile(miningKeyPath)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Println(apiCreateError)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	apiClient, err := mhq.NewClient(miningKey, installer.mhqEndpoint)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Println(apiCreateError)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	registerRequest := mhq.RegisterRigRequest{
		Name: rigName,
		Caps: systemInfo,
	}
	rigID, err := apiClient.RegisterRig(registerRequest)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to register your rig with MiningHQ. Please ensure that
you are connected to the internet and that the file '%s' contains the same
mining key that you can find under 'Mining' in your settings available at
https://www.mininghq.io/user/settings

If you are sure everything is in order, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
`,
			miningKeyPath)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	// Rig registered
	color.HiGreen("OK")

	// To create the config files we need to do two things
	// 1. Copy the mining_key to the installation directory
	// 2. Create a rig_id file in the installation directory
	fmt.Print("Create config files\t\t\t")
	err = copy.Copy(
		miningKeyPath,
		filepath.Join(installDir, "miner-controller", filepath.Base(miningKeyPath)))
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to copy your mining key to your installation.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	err = ioutil.WriteFile(
		filepath.Join(installDir, "miner-controller", "rig_id"),
		[]byte(rigID),
		0644)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to create the new rig files for your installation.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	// Config files created
	color.HiGreen("OK")

	// The MiningHQ miner service is embedded into this installer
	// It needs to be extracted into the installation directory
	fmt.Print("Installing MiningHQ Miner\n")

	embeddedFilename := "mininghq-miner"
	if strings.ToLower(runtime.GOOS) == "windows" {
		embeddedFilename = "mininghq-miner.exe"
	}
	embeddedFS := embedded.FS(false)
	embeddedFile, err := embeddedFS.Open("/miner-service/" + embeddedFilename)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to extract the miner from the installer.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	installFile, err := os.OpenFile(
		filepath.Join(installDir, embeddedFilename),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to create the miner in the correct location. Please check that
you have sufficient space on your harddrive.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	_, err = io.Copy(installFile, embeddedFile)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to install the miner to the correct location.
		`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	installFile.Close()
	embeddedFile.Close()

	// Install mininghq-miner as a service
	// We do this using a separate executable so that only the service install
	// requires Administrator/sudo rights and not the entire installer
	out, err := exec.Command(
		"sudo",
		"/home/donovan/Development/Go/code/src/github.com/donovansolms/mininghq-miner-manager/install-service/install-service",
		"-op", "install",
		"-serviceName", installer.serviceName,
		"-serviceDisplayName", installer.serviceDisplayName,
		"-serviceDescription", installer.serviceDescription,
		"-installedPath", installDir,
		"-serviceFilename", embeddedFilename,
	).CombinedOutput()
	if err != nil {
		color.HiRed("FAIL")
		fmt.Println("We were unable to install the miner service.")
		fmt.Printf(color.HiRedString("Include the following error in your report '%s', %s"), err.Error(), out)
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	// 	serviceConfig := &service.Config{
	// 		Name:             installer.serviceName,
	// 		DisplayName:      installer.serviceDisplayName,
	// 		Description:      installer.serviceDescription,
	// 		WorkingDirectory: installDir,
	// 		Executable:       filepath.Join(installDir, embeddedFilename),
	// 	}
	// 	svc, err := service.New(nil, serviceConfig)
	// 	if err != nil {
	// 		color.HiRed("FAIL")
	// 		fmt.Printf(`
	// We were unable to create the miner service.
	// `)
	// 		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
	// 		fmt.Println()
	// 		fmt.Println()
	// 		color.Unset()
	// 		os.Exit(1)
	// 	}
	// 	err = svc.Install()
	// 	if err != nil {
	// 		color.HiRed("FAIL")
	// 		fmt.Printf(`
	// We were unable to install the miner service.
	// `)
	// 		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
	// 		fmt.Println()
	// 		fmt.Println()
	// 		color.Unset()
	// 		os.Exit(1)
	// 	}

	installedCheckfilePath := filepath.Join(installer.homeDir, ".mhqpath")
	installedCheckfile, err := os.OpenFile(
		installedCheckfilePath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644)
	if err != nil {
		fmt.Printf(`
We were unable to create the installer check file in ~/.mhqpath. This
will cause MiningHQ services to be unable to detect the installation.

Please ensure you have the correct permissions to write to your home directory.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
	}
	defer installedCheckfile.Close()

	_, err = installedCheckfile.WriteString(installDir)
	if err != nil {
		fmt.Printf(`
We were unable to write to the installer check file in ~/.mhqpath. This
will cause MiningHQ services to be unable to detect the installation.

Please ensure you have the correct permissions to write to your home directory.
`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		os.Exit(0)
	}

	// Copy the manager
	managerBinaryPath, err := os.Executable()
	if err != nil {
		fmt.Printf(`
We were unable to copy the miner manager to your installation path.

Please ensure you have the correct permissions to write to your home directory.
	`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		os.Exit(0)
	}

	managerName := filepath.Base(managerBinaryPath)
	err = os.Rename(managerBinaryPath, filepath.Join(installDir, managerName))
	if err != nil {
		fmt.Printf(`
We were unable to copy the miner manager to your installation path.

Please ensure you have the correct permissions to write to your home directory.
	`)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		os.Exit(0)
	}

	// Service installed
	color.HiGreen("OK")

	// Start the mininghq-miner service
	// We do this using a separate executable so that only the service install
	// requires Administrator/sudo rights and not the entire installer
	out, err = exec.Command(
		"sudo",
		"/home/donovan/Development/Go/code/src/github.com/donovansolms/mininghq-miner-manager/install-service/install-service",
		"-op", "start",
		"-serviceName", installer.serviceName,
		"-serviceDisplayName", installer.serviceDisplayName,
		"-serviceDescription", installer.serviceDescription,
		"-installedPath", installDir,
		"-serviceFilename", embeddedFilename,
	).CombinedOutput()
	if err != nil {
		fmt.Printf(`
Unable to start the MiningHQ service, please start the 'MiningHQ-Miner' service manually. Reason: %s, %s
		`, err.Error(), out)
	}

	// 	err = svc.Start() // TODO: Start doesn't start it, problems
	// 	if err != nil {
	// 		fmt.Printf(`
	// Unable to start the MiningHQ service, please start the 'MiningHQ-Miner' service manually.
	// `)
	// 	}

	fmt.Printf(`


*************************
*  MiningHQ installed!  *
*************************

The MiningHQ Miner Manager and related services have been installed. You
can now open you MiningHQ dashboard to manage this rig.

https://www.mininghq.io/dashboard

The MiningHQ Manager is also now available in the installed path at
'%s'. You'll need to restart your rig for the automatic start to take effect.

Please join the MiningHQ community on Discord, Twitter and elsewhere, you can find
all our channels at https://www.mininghq.io/connect

Let's mine!
The MiningHQ Team
	`, installDir)

	fmt.Println()
	fmt.Println()
	return nil
}

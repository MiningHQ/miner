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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/donovansolms/mininghq-spec/spec/caps"
	"github.com/fatih/color"
	"github.com/mininghq/miner-controller/src/mhq"
	"github.com/mininghq/miner/helper"
	"github.com/otiai10/copy"
	input "github.com/tcnksm/go-input"
)

// Installer install the Miner Manager from the terminal
type Installer struct {
	// homeDir is the user's home directory
	homeDir string
	// os is the system operating system
	os string
	// mhqEndpoint is the MiningHQ API endpoint to use
	mhqEndpoint string

	serviceName        string
	serviceDisplayName string
	serviceDescription string
}

const (
	// Windows operating system
	Windows = "windows"
	// Linux operating system
	Linux = "linux"
	// MacOS operating system
	MacOS = "darwin"
)

// New creates a new installer instance
func NewInstaller(homeDir string, os string, mhqEndpoint string) (*Installer, error) {
	if strings.TrimSpace(homeDir) == "" {
		return nil, errors.New("A home directory must be set")
	}

	os = strings.ToLower(os)
	if strings.TrimSpace(os) != Windows && strings.TrimSpace(os) != MacOS &&
		strings.TrimSpace(os) != Linux {
		return nil, fmt.Errorf("OS may only be %s, %s or %s", Windows, MacOS, Linux)
	}

	installer := Installer{
		homeDir:            homeDir,
		os:                 os,
		mhqEndpoint:        mhqEndpoint,
		serviceName:        helper.ServiceName,
		serviceDisplayName: helper.ServiceDisplayName,
		serviceDescription: helper.ServiceDescription,
	}
	return &installer, nil
}

// Install the miner manager using a synchronous process,
// no feedback is given to the caller via channels
func (installer *Installer) Install() error {

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

You are about to install the MiningHQ Miner on this rig. This will
enable you to manage your mining, tweak your performance settings
and see this rig's stats - all from your dashboard.

Let's set up this rig.

We refer to any computer used to mine cryptocurrencies as a rig.
`)

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
via our help channels listed at https://www.mininghq.io/connect
`)
		os.Exit(0)
	}

	// Create the installation directory
	fmt.Print("Creating installation directory\t\t")
	avExcludeDirectory, err := helper.CreateInstallDirectories(installDir)
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
		helper.GetOSAVGuides(),
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
via our help channels listed at https://www.mininghq.io/connect
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
https://www.mininghq.io/connect
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

	miningKey, err := helper.GetMiningKeyFromFile(miningKeyPath)
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
https://www.mininghq.io/connect
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

	// We package a bunch of helper applications with the download. It needs
	// to be moved to the installation directory
	fmt.Print("Installing MiningHQ Miner\n")

	installFiles := map[string]string{
		"miner-service":     "miner-service",
		"service-installer": "install-service",
		"uninstaller":       "uninstall-mininghq",
	}

	for _, src := range installFiles {
		err = helper.CopyFile(filepath.Join("tools", src), filepath.Join(installDir, src))
		if err != nil {
			color.HiRed("FAIL")
			fmt.Printf(`
We were unable to install the service files. Please ensure you have write
permissions to the directory '%s'
		`, installDir)
			fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
			fmt.Println()
			fmt.Println()
			color.Unset()
			os.Exit(1)
		}
	}

	// Install mininghq-miner as a service
	// We do this using a separate executable so that only the service install
	// requires Administrator/sudo rights and not the entire installer
	out, err := exec.Command(
		"sudo",
		filepath.Join(installDir, installFiles["service-installer"]),
		"-op", "install",
		"-serviceName", installer.serviceName,
		"-serviceDisplayName", installer.serviceDisplayName,
		"-serviceDescription", installer.serviceDescription,
		"-installedPath", installDir,
		"-serviceFilename", installFiles["miner-service"],
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

Please ensure you have the correct permissions to write to your install directory.
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
		filepath.Join(installDir, installFiles["service-installer"]),
		"-op", "start",
		"-serviceName", installer.serviceName,
		"-serviceDisplayName", installer.serviceDisplayName,
		"-serviceDescription", installer.serviceDescription,
		"-installedPath", installDir,
		"-serviceFilename", installFiles["miner-service"],
	).CombinedOutput()
	if err != nil {
		fmt.Printf(`
Unable to start the MiningHQ service, please start the 'MiningHQ-Miner' service manually. Reason: %s, %s
		`, err.Error(), out)
	}

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

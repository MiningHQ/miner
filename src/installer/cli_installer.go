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
		homeDir:     homeDir,
		os:          os,
		mhqEndpoint: mhqEndpoint,
	}
	return &installer, nil
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
	avExcludeDirectory, err := installer.CreateInstallDirectories(installDir)
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
		installer.GetOSAVGuides(),
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
	miningKeyBytes, err := ioutil.ReadFile(miningKeyPath)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Println(apiCreateError)
		fmt.Printf(color.HiRedString("Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	miningKey := strings.TrimSpace(string(miningKeyBytes))
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
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
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
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
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
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	// Config files created
	color.HiGreen("OK")

	// The MiningHQ miner service is embedded into this installer
	// It needs to be extracted into the installation directory
	fmt.Print("Installing MiningHQ Miner\t\t")

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
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	defer embeddedFile.Close()

	installFile, err := os.OpenFile(
		filepath.Join(installDir, embeddedFilename),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to create the miner in the correct location. Please check that
you have sufficient space on your harddrive.
`)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}
	defer installFile.Close()

	_, err = io.Copy(installFile, embeddedFile)
	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
		We were unable to install the miner to the correct location.
		`)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
		os.Exit(1)
	}

	// TODO: Install mininghq-miner as a service

	// Config files created
	color.HiGreen("OK")

	fmt.Println()
	fmt.Println()
	return nil
}

// CreateInstallDirectories creates the directories needed for installation
//
// It returns the path where miners will be installed, users need to exclude
// this path from antivirus scanning.
func (installer *CLIInstaller) CreateInstallDirectories(
	installDirectory string) (string, error) {

	paths := []string{
		"miner-controller",
		filepath.Join("miner-controller", "miners"),
	}
	avExcludePath := "miners"
	for _, path := range paths {
		path = filepath.Join(installDirectory, path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return avExcludePath, fmt.Errorf(
				"Unable to create installation directory: '%s': %s",
				path,
				err)
		}
		if strings.Contains(path, "miners") {
			avExcludePath = path
		}
	}
	return avExcludePath, nil
}

// GetOSAVGuides returns a list of links and descriptions for antivirus
// directory exclude guides
func (installer *CLIInstaller) GetOSAVGuides() string {
	return fmt.Sprintf(`https://www.mininghq.io/help/antivirus`)
}

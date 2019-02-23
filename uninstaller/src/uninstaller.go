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
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ProtonMail/go-autostart"
	"github.com/fatih/color"
	"github.com/mininghq/miner-controller/src/mhq"
	"github.com/mininghq/miner/helper"
	ps "github.com/mitchellh/go-ps"
)

const (
	// Windows operating system
	Windows = "windows"
	// Linux operating system
	Linux = "linux"
	// MacOS operating system
	MacOS = "darwin"
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

// NewInstaller creates a new installer instance
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

// Uninstall uninstalls the miner manager and services using
// a synchronous process
func (installer *Installer) Uninstall(
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

	// TODO: Doesn't work on Windows, keeps returning " ' is invalid "
	// ui := &input.UI{}
	// question := "\nAre you sure you wish to remove the MiningHQ Miner and all MiningHQ services? [Y/yes/N/no]"
	// response, _ := ui.Ask(question, &input.Options{
	// 	Required: true,
	// 	Loop:     true,
	// 	ValidateFunc: func(s string) error {
	// 		validConfirmations := map[string]bool{
	// 			"y":   true,
	// 			"yes": true,
	// 			"n":   true,
	// 			"no":  true,
	// 		}
	// 		answer := strings.ToLower(s)
	// 		if _, ok := validConfirmations[answer]; !ok {
	// 			return fmt.Errorf(
	// 				"Answer '%s' is invalid. Must be 'y', 'yes', 'n' or 'no'", s)
	// 		}
	// 		return nil
	// 	},
	// })
	// 	allowContinue := strings.ToLower(response)
	// 	if allowContinue == "n" || allowContinue == "no" {
	// 		color.HiRed("********************************")
	// 		color.HiRed("* Uninstall has been cancelled *")
	// 		color.HiRed("********************************")
	// 		color.HiYellow(`
	// Something wrong? If so, please let us know by getting in contact
	// via our help channels listed at https://www.mininghq.io/help
	// `)
	// 		os.Exit(0)
	// 	}

	// Stop the service
	fmt.Print("Stopping services\t\t\t")
	minerServiceName := "miner-service"
	if strings.ToLower(runtime.GOOS) == "windows" {
		minerServiceName = "miner-service.exe"
	}

	stopError := fmt.Sprintf(`
We were unable to stop the MiningHQ services. Please stop the
services manually.

If you need help, please contact support to resolve
the issue. Support can be contacted via our help channels listed at
https://www.mininghq.io/help
`)

	processReference, err := installer.findProcessByName(minerServiceName)
	if err != nil {
		color.HiYellow("NOTICE")
		fmt.Printf(`
We were unable to find running MiningHQ services, they might be stopped already.
Uninstall will continue...
`)
		fmt.Printf(color.HiYellowString(
			"Include the following message if you contact support '%s'"), err.Error())
		fmt.Println()
		fmt.Println()
		color.Unset()
	} else {
		// Different OSs have different killing methods
		if strings.ToLower(runtime.GOOS) == "windows" {
			// Windows
			err = exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(processReference.Pid())).Run()
		} else {
			// -PID (minus PID) to kill the process and all their children
			err = syscall.Kill(-processReference.Pid(), syscall.SIGKILL)
		}
		if err != nil {
			color.HiRed("FAIL")
			fmt.Println(stopError)
			fmt.Printf(color.HiRedString(
				"Include the following error in your report '%s'"), err.Error())
			fmt.Println()
			fmt.Println()
			color.Unset()
		} else {
			color.HiGreen("OK")
		}
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
	fmt.Print("Removing the MiningHQ Miner service\t")

	// NOTE We no longer run as a service
	// serviceFilename := "mininghq-miner"
	// serviceInstallerFilename := "install-service"
	// if strings.ToLower(runtime.GOOS) == "windows" {
	// 	serviceFilename = "mininghq-miner.exe"
	// 	serviceInstallerFilename = "install-service.exe"
	// }
	//
	// var out []byte
	// if strings.ToLower(runtime.GOOS) == "windows" {
	// 	// Uninstall mininghq-miner as a service
	// 	// We do this using a separate executable so that only the service uninstall
	// 	// requires Administrator/sudo rights and not the entire installer
	// 	out, err = exec.Command(
	// 		"cmd.exe", "/C",
	// 		"sc.exe", "delete",
	// 		installer.serviceName,
	// 	).CombinedOutput()
	// } else {
	// 	// Uninstall mininghq-miner as a service
	// 	// We do this using a separate executable so that only the service uninstall
	// 	// requires Administrator/sudo rights and not the entire installer
	// 	out, err = exec.Command(
	// 		"sudo",
	// 		filepath.Join(installedPath, serviceInstallerFilename),
	// 		"-op", "uninstall",
	// 		"-serviceName", installer.serviceName,
	// 		"-serviceDisplayName", installer.serviceDisplayName,
	// 		"-serviceDescription", installer.serviceDescription,
	// 		"-installedPath", installedPath,
	// 		"-serviceFilename", serviceFilename,
	// 	).CombinedOutput()
	// }
	// END NOTE

	// We now use autorun
	app := &autostart.App{
		Name:        installer.serviceName,
		DisplayName: installer.serviceDisplayName,
		Exec:        []string{filepath.Join(installedPath, installer.serviceName)},
	}
	if app.IsEnabled() == true {
		err = app.Disable()
	}

	if err != nil {
		color.HiRed("FAIL")
		fmt.Printf(`
We were unable to uninstall the miner service (it might already be uninstalled).
`)
		fmt.Printf(color.HiRedString(
			"Include the following error in your report '%s'"), err.Error())
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

	fmt.Println("Uninstaller will exit in 10 seconds...")
	time.Sleep(time.Second * 10)

	fmt.Println("Press Enter to exit")
	os.Stdin.Read([]byte{0})

	return nil
}

// findProcessByName finds and returns the process information by name
func (installer *Installer) findProcessByName(name string) (ps.Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return process, nil
		}
	}
	return nil, fmt.Errorf("Unable to find process with name '%s'", name)
}

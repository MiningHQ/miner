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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
)

// main runs one of three operations, install, uninstall or start
// This is a standalone tool so that the GUI and CLI installers don't have
// to be run as sudo/administrator but rather only sudo the service install
func main() {
	var operation string
	var serviceName string
	var serviceDisplayName string
	var serviceDescription string
	var installedPath string
	var serviceFilename string

	flag.StringVar(&operation, "op", "", "The operation to perform")
	flag.StringVar(&serviceName, "serviceName", "", "The serviceName for the service")
	flag.StringVar(&serviceDisplayName, "serviceDisplayName", "", "The serviceDisplayName for the service")
	flag.StringVar(&serviceDescription, "serviceDescription", "", "The serviceDescription for the service")
	flag.StringVar(&installedPath, "installedPath", "", "The installedPath for the service")
	flag.StringVar(&serviceFilename, "serviceFilename", "", "The serviceFilename for the service")

	flag.Parse()

	serviceConfig := &service.Config{
		Name:             serviceName,
		DisplayName:      serviceDisplayName,
		Description:      serviceDescription,
		WorkingDirectory: installedPath,
		Executable:       filepath.Join(installedPath, serviceFilename),
	}
	svc, err := service.New(nil, serviceConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	operation = strings.ToLower(operation)
	// Install mininghq-miner as a service
	if operation == "install" {
		err = svc.Install()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Uninstall mininghq-miner service
	if operation == "uninstall" {
		err = svc.Stop()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = svc.Uninstall()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if operation == "start" {
		err = svc.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
)

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

/*
  MiningHQ Miner - The MiningHQ Miner service
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
	"log"

	"github.com/donovansolms/mininghq-miner-manager/helper"
	"github.com/donovansolms/mininghq-miner-manager/miner-service/src/miner"
	"github.com/kardianos/service"
)

var logger service.Logger

// Service setup
// This is mainly taken from the sample at github.com/kardianos/service
//
// Apart from being a simple service it downloads and runs the MiningHQ
// Miner Controller. The Controller runs all the mining logic.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	serviceConfig := &service.Config{
		Name:        helper.ServiceName,
		DisplayName: helper.ServiceDisplayName,
		Description: helper.ServiceDescription,
	}

	// Set up the new miner
	minerService, err := miner.New()
	if err != nil {
		log.Fatal(err)
	}

	// Create a service from the miner instance
	svc, err := service.New(minerService, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the error log handling for the service
	errs := make(chan error, 5)
	logger, err = svc.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}
	// We set the miner logger to be the service-created logger for the
	// OS we are running on
	// Windows would be event log
	// Linux would be syslog
	// MacOS would be system log
	minerService.SetLogger(logger)

	go func() {
		for {
			err = <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err = service.Control(svc, *svcFlag)
		if err != nil {
			log.Printf("Valid actions are: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = svc.Run()
	if err != nil {
		logger.Error(err)
	}
}

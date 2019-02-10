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

package miner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	unattended "github.com/ProjectLimitless/go-unattended"
	"github.com/kardianos/service"
	logrus "github.com/sirupsen/logrus"
)

// Miner is the primary MiningHQ service on a user's rig that manages and
// updates the Miner controller
type Miner struct {
	// log provides the proper OS based logging
	log service.Logger
	// updateWrapper handles the automatic updates of the miner controller
	updateWrapper *unattended.Unattended
	// exit channel signals when we should shut down
	exit chan struct{}
}

// New creates a new instance of the Miner
func New() (*Miner, error) {
	miner := Miner{}
	return &miner, nil
}

// Run should be called shortly after the program entry point.
// After Interface.Stop has finished running, Run will stop blocking.
// After Run stops blocking, the program must exit shortly after.
func (miner *Miner) run() error {

	// Start the miner controller with updates enabled
	err := miner.updateWrapper.Run()
	if err != nil {
		miner.log.Errorf("Unable to run miner controller: %s", err)
		return err
	}

	miner.log.Infof("Running miner service '%v'.", service.Platform())
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			miner.log.Infof("Still running at %v...", tm)
		case <-miner.exit:
			ticker.Stop()
			return nil
		}
	}
}

// Start signals to the OS service manager the given service should start.
func (miner *Miner) Start(s service.Service) error {
	if service.Interactive() {
		miner.log.Info("Running in terminal.")
	} else {
		miner.log.Info("Running under service manager.")
	}
	miner.exit = make(chan struct{})

	// Set up unattended updates
	miner.log.Info("Setting up Unattended updates")

	// TODO Unattended wants a Logrus log, it should rather take a standard
	// Go log interface
	// For now, let's give it a log if we're running in interactive mode
	logOutputFormat := logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	}
	logrus.SetFormatter(&logOutputFormat)
	logrus.SetLevel(logrus.DebugLevel)
	logger := logrus.WithFields(logrus.Fields{
		"service_class": "mininghq-miner",
	})
	logrus.SetOutput(ioutil.Discard)
	if service.Interactive() {
		logrus.SetOutput(os.Stdout)
	}

	executablePath, err := os.Executable()
	if err != nil {
		miner.log.Errorf("Unable to get executable path: %s", err)
		logger.Fatalf("Unable to get executable path: %s", err)
	}

	basePath := filepath.Dir(executablePath)

	miner.updateWrapper, err = unattended.New(
		"TEST001", // TODO clientID
		unattended.Target{ // target
			VersionsPath:   filepath.Join(basePath, "miner-controller"),
			AppID:          fmt.Sprintf("miner-controller-%s", strings.ToLower(runtime.GOOS)),
			UpdateEndpoint: "https://unattended.mininighq.io",
			//UpdateEndpoint:        "https://unattended-old.local",
			UpdateChannel:         "stable",
			ApplicationName:       "mininghq-miner-controller",
			ApplicationParameters: []string{},
		},
		time.Hour, // UpdateCheckInterval
		logger,
	)
	if err != nil {
		miner.log.Errorf("Unable to create Unattended update manager: %s", err)
		return err
	}

	// During construction we check for any updates as well, this has the
	// side effect that *if* the software isn't available, it will be downloaded
	hasUpdate, err := miner.updateWrapper.ApplyUpdates()
	if err != nil {
		miner.log.Errorf("Unable to apply controller updates: %s", err)
		return err
	}
	if hasUpdate == false {
		miner.log.Infof("No updates available for miner-controller")
	}

	// Start should not block. Do the actual work async.
	go miner.run()
	return nil
}

// Stop signals to the OS service manager the given service should stop.
func (miner *Miner) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	miner.log.Info("Stopping miner service")
	close(miner.exit)
	// Stop the miner controller when the service stops
	return miner.updateWrapper.Stop()
}

// SetLogger sets the logger for the service
func (miner *Miner) SetLogger(log service.Logger) {
	miner.log = log
}

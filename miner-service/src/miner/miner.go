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
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	unattended "github.com/ProjectLimitless/go-unattended"
	logrus "github.com/sirupsen/logrus"
)

// Miner is the primary MiningHQ service on a user's rig that manages and
// updates the Miner controller
type Miner struct {
	// log is the service log
	log *logrus.Entry
	// updateWrapper handles the automatic updates of the miner controller
	updateWrapper *unattended.Unattended
}

// New creates a new instance of the Miner
func New() (*Miner, error) {
	miner := Miner{}
	return &miner, nil
}

// Run starts the miner download and runner
func (miner *Miner) Run() error {

	// TODO Unattended wants a Logrus log, it should rather take a standard
	// Go log interface
	// For now, let's give it a log as if we're running in interactive mode
	logOutputFormat := logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05",
	}
	logrus.SetFormatter(&logOutputFormat)
	logrus.SetLevel(logrus.DebugLevel)
	miner.log = logrus.WithFields(logrus.Fields{
		"service_class": "mininghq-miner",
	})
	logrus.SetOutput(os.Stdout)

	// Set up unattended updates
	miner.log.Info("Setting up Unattended updates")

	executablePath, err := os.Executable()
	if err != nil {
		miner.log.Fatalf("Unable to get executable path: %s", err)
	}

	basePath := filepath.Dir(executablePath)

	miner.updateWrapper, err = unattended.New(
		"TEST001", // TODO clientID
		unattended.Target{ // target
			VersionsPath:   filepath.Join(basePath, "miner-controller"),
			AppID:          fmt.Sprintf("miner-controller-%s", strings.ToLower(runtime.GOOS)),
			UpdateEndpoint: "https://unattended.mininghq.io",
			//UpdateEndpoint:        "https://unattended-old.local",
			UpdateChannel:         "stable",
			ApplicationName:       "miner-controller",
			ApplicationParameters: []string{},
		},
		time.Hour, // UpdateCheckInterval
		miner.log,
	)
	if err != nil {
		miner.log.Fatalf("Unable to create Unattended update manager: %s", err)
	}

	// During construction we check for any updates as well, this has the
	// side effect that *if* the software isn't available, it will be downloaded
	hasUpdate, err := miner.updateWrapper.ApplyUpdates()
	if err != nil {
		// If unattended updates can't be applied, it's a real problem
		miner.log.Errorf("Unable to apply controller updates: %s", err)
		panic(err)
	}
	if hasUpdate == false {
		miner.log.Infof("No updates available for miner-controller")
	}

	// Start the miner controller with updates enabled
	err = miner.updateWrapper.Run()
	if err != nil {
		miner.log.Errorf("Unable to run miner controller: %s", err)
	}
	return err
}

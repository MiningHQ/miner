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

package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	ServiceName        = "GoServiceExampleLogging"
	ServiceDisplayName = "Go Service Example for Logging"
	ServiceDescription = "This is an example Go service that outputs log messages."
)

// Helper contains shares helper install functions
type Helper struct {
}

// CreateInstallDirectories creates the directories needed for installation
//
// It returns the path where miners will be installed, users need to exclude
// this path from antivirus scanning.
func (helper *Helper) CreateInstallDirectories(
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
func (helper *Helper) GetOSAVGuides() string {
	return fmt.Sprintf(`https://www.mininghq.io/help/antivirus`)
}

// GetMiningKeyFromFile reads the user's mining key from a file
func (helper *Helper) GetMiningKeyFromFile(path string) (string, error) {
	miningKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(miningKeyBytes)), nil
}

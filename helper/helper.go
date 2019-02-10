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

// Package helper implements various helper functions
package helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ServiceName for the service
	ServiceName = "mininghq-miner"
	// ServiceDisplayName for the service
	ServiceDisplayName = "MiningHQ Miner"
	// ServiceDescription for the service
	ServiceDescription = "The MiningHQ.io Miner service for controlling mining with this rig"
)

// CreateInstallDirectories creates the directories needed for installation
//
// It returns the path where miners will be installed, users need to exclude
// this path from antivirus scanning.
func CreateInstallDirectories(
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
func GetOSAVGuides() string {
	return fmt.Sprintf(`https://www.mininghq.io/help/antivirus`)
}

// GetMiningKeyFromFile reads the user's mining key from a file
func GetMiningKeyFromFile(path string) (string, error) {
	miningKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(miningKeyBytes)), nil
}

// CopyFile copies a file from src to dst
func CopyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Modified to get executable file
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0554)
	//out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

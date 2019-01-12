/*
  MiningHQ Miner Manager - The MiningHQ Miner Manager GUI
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

// Package main implements the main runnable for the miner manager.
// It constructs and launches the Electron front-end
// package main
//
// import (
// 	"flag"
// 	"fmt"
// )
//
// // AppName is injected by the Astilectron packager
// var AppName string
//
// // BuiltAt is injected by the Astilectron packager
// var BuiltAt string
//
// // main implements the main runnable of the application
// func main() {
// 	// Grab the command-line flags
// 	debug := flag.Bool("d", false, "Enable debug mode")
// 	flag.Parse()
//
// 	if *debug {
// 		fmt.Println("RUNNING DEBUG")
// 	}
// 	fmt.Println("Running")
// }

// HACK Testing code below
//
package main

import (
	"github.com/donovansolms/mininghq-miner-manager/src/cmd"
)

// AppName is injected by the Astilectron bundler
var AppName string

func main() {
	cmd.AppName = AppName
	cmd.Asset = Asset
	cmd.RestoreAssets = RestoreAssets
	cmd.Execute()
}

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
	"log"

	"github.com/mininghq/miner/miner-service/src/miner"
)

// Rownloads and runs the MiningHQ Miner Controller.
// The Controller runs all the mining logic.
func main() {

	// Set up the new miner
	minerService, err := miner.New()
	if err != nil {
		log.Fatal(err)
	}

	// Run
	err = minerService.Run()
	if err != nil {
		log.Fatal(err)
	}

}

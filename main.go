package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 6
	numUsers := 3
	numMal := 4
	pos.Run(runType, numValidators, numUsers, numMal)
}

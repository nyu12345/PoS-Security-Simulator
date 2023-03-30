package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "manual"
	numValidators := 100
	numUsers := 20
	numMal := 0
	pos.Run(runType, numValidators, numUsers, numMal)
}

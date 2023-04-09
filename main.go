package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 6
	numUsers := 3
	numMal := 0
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, attack)
}

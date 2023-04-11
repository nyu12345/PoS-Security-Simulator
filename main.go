package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 10
	numUsers := 3
	numMal := 0
	committeeSize := numValidators / 3
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, attack)
}

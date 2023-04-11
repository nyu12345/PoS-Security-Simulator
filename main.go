package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	attack := "balance"
	numValidators := 6
	numUsers := 3
	numMal := 2
	pos.Run(runType, numValidators, numUsers, numMal, attack)
}

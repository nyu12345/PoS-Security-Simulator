package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 100
	numUsers := 3
	numMal := 0
	committeeSize := 50
	//pos, delegated, or reputation
	blockchainType := "pos"
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, blockchainType, attack)
}

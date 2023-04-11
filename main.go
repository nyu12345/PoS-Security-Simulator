package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 30
	numUsers := 5
	numMal := 0
	committeeSize := 3
	delegateSize := 5
	//pos or reputation
	blockchainType := "reputation"
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

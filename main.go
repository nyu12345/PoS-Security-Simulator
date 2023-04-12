package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 10
	numUsers := 5
	numMal := 7
	committeeSize := 6
	delegateSize := 5
	//pos or reputation
	blockchainType := "pos"
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

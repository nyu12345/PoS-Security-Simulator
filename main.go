package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 30
	numUsers := 5
	numMal := 8
	committeeSize := 4
	delegateSize := 5
	//pos or reputation
	blockchainType := "reputation"
	attack := "balance"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

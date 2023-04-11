package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 100
	numUsers := 5
	numMal := 0
	committeeSize := 20
	delegateSize := 5
	//pos or reputation
	blockchainType := "pos"
	attack := "none"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

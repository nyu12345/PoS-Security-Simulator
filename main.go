package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 10
	numUsers := 5
	numMal := 3
	committeeSize := 6
	delegateSize := 5
	//pos or reputation
	blockchainType := "reputation"
	attack := "network_partition"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	numValidators := 6
	numUsers := 3
	numMal := 2
	committeeSize := 3
	//pos, delegated, or reputation
	blockchainType := "pos"
	//To do short range attack, type "network_partition", else ""
	attack := "network_partition"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, blockchainType, attack)
}

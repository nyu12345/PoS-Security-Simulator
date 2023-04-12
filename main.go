package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "auto"
	//100
	numValidators := 10
	//10
	numUsers := 5
	//20, 50, 70
	numMal := 3
	//20
	committeeSize := 6
	//5
	delegateSize := 5
	//pos, slashing, or reputation
	blockchainType := "reputation"
	//network_partition, balance
	attack := "network_partition"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

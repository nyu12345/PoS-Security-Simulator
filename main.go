package main

import (
	"PoS-Security-Simulator/pos"
)

func main() {
	//manual or auto
	runType := "manual"
	//100
	numValidators := 100
	//10
	numUsers := 10
	//20, 50, 70
	numMal := 20
	//20
	committeeSize := 20
	//5
	delegateSize := 5
	//pos, slashing, or reputation
	blockchainType := "pos"
	//network_partition, balance
	attack := "network_partition"
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

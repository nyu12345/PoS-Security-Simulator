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
<<<<<<< HEAD
	attack := "network_partition"
=======
	attack := "balance"
>>>>>>> 3691bb82bca5c6da7b94259cd232c60ef57353f8
	pos.Run(runType, numValidators, numUsers, numMal, committeeSize, delegateSize, blockchainType, attack)
}

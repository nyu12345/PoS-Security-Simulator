package pos

type GenesisBlockMessage struct {
	genesisBlock Block
}

type ValidateBlockMessage struct {
	newBlock Block
	oldBlock Block
	fork0Length int
	fork1Length int
	proposerView int
}

type ValidationStatusMessage struct {
	isValid bool
}

type NewTransactionMessage struct {
	transaction Transaction
}

type VerifiedBlockMessage struct {
	transactions []Transaction
	newBlock     Block
}

type ConsensusMessage struct {
	blockchain              []Block
	unconfirmedTransactions map[int]Transaction
	confirmedTransactions   map[int]bool
}

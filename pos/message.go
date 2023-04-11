package pos

type GenesisBlockMessage struct {
	genesisBlock Block
}

type ValidateBlockMessage struct {
	newBlock Block
}

type ValidateShortAttackBlockMessage struct {
	newBlock    Block
	newBlockTwo Block
	index       int
}

type ValidationStatusMessage struct {
	isValid bool
}

type ValidationShortAttackStatusMessage struct {
	isValid    bool
	isValidTwo bool
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

package pos

type GenesisBlockMessage struct {
	genesisBlock Block
}

type ValidateBlockMessage struct {
	newBlock Block
}

type ValidationStatusMessage struct {
	isValid bool
}

type NewTransactionMessage struct {
	transaction Transaction
}

type VerifiedBlockMessage struct {
	transactions []Transaction
	newBlock Block
}

type TransactionConsensusMessage struct {
	unconfirmedTransactions map[int]Transaction
	confirmedTransactions map[int]bool
}


package pos

type GenesisBlockMessage struct {
	genesisBlock Block
}

type ValidateBlockMessage struct {
	newBlock Block
	malVote bool
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

type DelegateVoteRequestMessage struct {
	delegateSize int
}

type DelegateVoteMessage struct {
	delegateVotes []*Validator
}

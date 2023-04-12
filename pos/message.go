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
}

type ValidationStatusMessage struct {
	isValid bool
}

type ValidationShortAttackStatusMessage struct {
	isValid    bool
	isValidTwo bool
}

type ValidationForkedChainStatusMessage struct {
	isValid bool
}

type NewTransactionMessage struct {
	transaction Transaction
}

type VerifiedBlockMessage struct {
	transactions []Transaction
	newBlock     Block
}

type VerifiedShortAttackBlockMessage struct {
	transactions []Transaction
	newBlock     Block
}

type VerifiedShortAttackBlockTwoMessage struct {
	transactions []Transaction
	newBlockTwo  Block
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

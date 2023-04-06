package pos

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

type VerifiedTransactionMessage struct {
	transactions []Transaction
}

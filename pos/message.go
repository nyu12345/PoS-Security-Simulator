package pos

type ValidateBlockMessage struct {
	newBlock Block
	oldBlock Block
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

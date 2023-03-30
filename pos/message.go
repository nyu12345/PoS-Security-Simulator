package pos

import (
	"sync"
)

type ValidateBlockMessage struct {
	newBlock Block
	oldBlock Block
	wg       *sync.WaitGroup
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

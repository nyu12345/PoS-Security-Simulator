package pos

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Hash         string
	PrevHash     string
	Validator    string
	IsMalicious  bool
	PrevBlock    *Block
	NextBlock    *Block
	CurBlock     *Block
}

// SHA256 hasing
// calculateHash is a simple SHA256 hashing function
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// calculateBlockHash returns the hash of all block information
func calculateBlockHash(block Block) string {
	record := fmt.Sprintf("%d%s%s", block.Index, block.Timestamp, block.PrevHash)
	for _, transaction := range block.Transactions {
		record += fmt.Sprintf("%d%p%p%s%f", transaction.ID, transaction.Sender, transaction.Receiver, transaction.Signature, transaction.Reward)
	}
	return calculateHash(record)
}

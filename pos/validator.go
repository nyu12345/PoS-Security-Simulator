package pos

import (
	"bufio"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Validator struct {
	conn                    net.Conn
	incomingChannel         chan interface{}
	outgoingChannel         chan interface{}
	transactionChannel      chan NewTransactionMessage
	Address                 string
	Stake                   float64
	unconfirmedTransactions map[int]Transaction
	confirmedTransactions   map[int]bool
	IsMalicious             bool
	validatorLock           sync.Mutex
	committeeCount          int
	proposerCount           int
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, proposer *Validator) (Block, error) {

	var newBlock Block

	//read transactions from local mempool if there are enough
	transactions := []Transaction{}
	if len(proposer.unconfirmedTransactions) > 0 {
		transactionsSize := len(proposer.unconfirmedTransactions)
		if transactionsSize > 5 {
			transactionsSize = 5
		}
		for id := range proposer.unconfirmedTransactions {
			transactions = append(transactions, proposer.unconfirmedTransactions[id])
			if transactionsSize == len(transactions) {
				break
			}
		}
	} else {
		//else return an error
		err := errors.New("No transactions to validate")
		return newBlock, err
	}

	//set block information

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = proposer.Address
	newBlock.Transactions = transactions
	newBlock.Hash = calculateBlockHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println("old block is not the previous block")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println("old block hash does not match with the previous hash")
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		fmt.Println("Recomputation of the hash is incorrect")
		return false
	}

	return true
}

func isTransactionValid(transaction Transaction, validator *Validator) bool {
	//Sender and receiver are both real users
	if transaction.Sender == nil || transaction.Receiver == nil {
		io.WriteString(validator.conn, "Transaction sender or receiver is not an active user\n")
		return false
	}

	//Public key verifies transaction
	signatureBytes, _ := hex.DecodeString(transaction.Signature)

	// Compute the transaction hash
	data := fmt.Sprintf("%d%p%p%f%f", transaction.ID, transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Reward)
	hash := sha256.Sum256([]byte(data))

	err := rsa.VerifyPKCS1v15(transaction.Sender.PublicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		io.WriteString(validator.conn, "Transaction could not be verified with public key\n")
		return false
	}

	//Transaction was already spent
	if validator.confirmedTransactions[transaction.ID] == true {
		io.WriteString(validator.conn, "Transaction was already spent\n")
		return false
	}
	//User has insufficient funds
	transaction.Sender.userLock.Lock()
	if (transaction.Amount + transaction.Reward) > transaction.Sender.Balance {
		io.WriteString(validator.conn, "Sender has insufficient funds\n")
		transaction.Sender.userLock.Unlock()
		return false
	}
	transaction.Sender.userLock.Unlock()
	io.WriteString(validator.conn, "Transaction is valid\n")
	return true
}

func handleValidatorConnection(conn net.Conn, runType string, malString string) {
	defer conn.Close()

	//Enter initial stake and whether or not validator is malicious
	io.WriteString(conn, "Enter token stake:\n")
	scannedBalance := bufio.NewScanner(conn)
	if runType == "auto" {
		randomStake := 0.0
		rand.Seed(time.Now().UnixNano())
		randomStake = rand.Float64()*700 + 300
		randomStakeString := fmt.Sprintf("%f", randomStake)
		scannedBalance = bufio.NewScanner(strings.NewReader(randomStakeString))
	}
	balance := 0.0
	var err error
	isMal := false
	for scannedBalance.Scan() {
		balance, err = strconv.ParseFloat(scannedBalance.Text(), 64)
		if err != nil {
			io.WriteString(conn, scannedBalance.Text()+" not a number")
			return
		}
		break
	}

	io.WriteString(conn, "Is this node malicious (y/n)\n")
	scannedMal := bufio.NewScanner(conn)
	if runType == "auto" {
		scannedMal = bufio.NewScanner(strings.NewReader(malString))
	}
	for scannedMal.Scan() {
		if scannedMal.Text() != "y" && scannedMal.Text() != "n" {
			io.WriteString(conn, scannedMal.Text()+" is not a valid response\n Please enter 'y' or 'n' ")
			return
		}
		if scannedMal.Text() == "y" {
			isMal = true
		}
		break
	}

	//Calculate address based on time
	t := time.Now()
	address := calculateHash(t.String())

	//Instantiate new validator
	unconfirmedTransactions := make(map[int]Transaction)
	confirmedTransactions := make(map[int]bool)
	curValidator := &Validator{
		conn:                    conn,
		incomingChannel:         make(chan interface{}),
		outgoingChannel:         make(chan interface{}),
		transactionChannel:      make(chan NewTransactionMessage),
		Address:                 address,
		Stake:                   balance,
		unconfirmedTransactions: unconfirmedTransactions,
		confirmedTransactions:   confirmedTransactions,
		IsMalicious:             isMal,
		validatorLock:           sync.Mutex{},
	}
	validators = append(validators, curValidator)

	if isMal {
		malValidators = append(malValidators, curValidator)
	}

	fmt.Printf("new validator count: %d\n", len(validators))

	//listen for transactions in transaction channel
	go func() {
		for {
			msg := <-curValidator.transactionChannel
			//Receiving unverified transactions
			io.WriteString(conn, "Received unverified transaction\n")
			isValid := isTransactionValid(msg.transaction, curValidator)
			if isValid {
				curValidator.unconfirmedTransactions[msg.transaction.ID] = msg.transaction
			}
		}
	}()

	//listen for messages in communication channel
	for {
		msg := <-curValidator.incomingChannel
		switch msg := msg.(type) {
		//Receiving block to validate
		case ValidateBlockMessage:
			io.WriteString(conn, "Received a Block to validate\n")
			isValid := isBlockValid(msg.newBlock, msg.oldBlock)
			validationStatusMessage := ValidationStatusMessage{
				isValid: isValid,
			}
			curValidator.outgoingChannel <- validationStatusMessage
		//Receiving verified transactions
		case VerifiedTransactionMessage:
			io.WriteString(conn, "Received verified transaction\n")
			//put verified transactions into confirmed slice for validator
			for _, transaction := range msg.transactions {
				curValidator.confirmedTransactions[transaction.ID] = true
			}

			//take transactions out of unconfirmed map
			for _, transaction := range msg.transactions {
				delete(curValidator.unconfirmedTransactions, transaction.ID)
			}
		default:
			io.WriteString(conn, "Received an unknown struct: %+v\n")
		}
	}

}

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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Validator struct {
	conn                       net.Conn
	incomingChannel            chan interface{}
	outgoingChannel            chan interface{}
	transactionChannel         chan NewTransactionMessage
	delegateVoteRequestChannel chan DelegateVoteRequestMessage
	delegateVoteChannel        chan DelegateVoteMessage
	Address                    string
	Stake                      float64
	unconfirmedTransactions    map[int]Transaction
	confirmedTransactions      map[int]bool
	IsMalicious                bool
	validatorLock              sync.Mutex
	transactionPoolLock        sync.Mutex
	committeeCount             int
	proposerCount              int
	blockSuccessCount          int
	reputation                 float64
	HeadBlockchain             *Block
	CurrBlockchain             *Block
}

// generateBlock creates a new block using previous block's hash
func generateBlock(proposer *Validator) (Block, error) {

	var newBlock Block

	//read transactions from local mempool if there are enough
	transactions := []Transaction{}
	if len(proposer.unconfirmedTransactions) > 0 {
		proposer.validatorLock.Lock()
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
		proposer.validatorLock.Unlock()
	} else {
		//else return an error
		err := errors.New("No transactions to validate")
		return newBlock, err
	}

	//set block information

	t := time.Now()
	oldBlock := proposer.CurrBlockchain
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Validator = proposer.Address
	newBlock.Transactions = transactions
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.IsMalicious = proposer.IsMalicious

	return newBlock, nil
}

func isBlockValid(newBlockPtr *Block) bool {
	newBlock := *newBlockPtr
	oldBlock := newBlock.PrevBlock
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
		conn:                       conn,
		incomingChannel:            make(chan interface{}),
		outgoingChannel:            make(chan interface{}),
		transactionChannel:         make(chan NewTransactionMessage),
		delegateVoteRequestChannel: make(chan DelegateVoteRequestMessage),
		delegateVoteChannel:        make(chan DelegateVoteMessage),
		Address:                    address,
		Stake:                      balance,
		unconfirmedTransactions:    unconfirmedTransactions,
		confirmedTransactions:      confirmedTransactions,
		IsMalicious:                isMal,
		validatorLock:              sync.Mutex{},
		transactionPoolLock:        sync.Mutex{},
		committeeCount:             0,
		proposerCount:              0,
		reputation:                 5.0,
		HeadBlockchain:             CertifiedBlockchain,
		CurrBlockchain:             CertifiedBlockchain}

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
			curValidator.transactionPoolLock.Lock()
			if isValid {
				curValidator.unconfirmedTransactions[msg.transaction.ID] = msg.transaction
			}
			curValidator.transactionPoolLock.Unlock()
		}
	}()

	//listen for delegate vote requests in delegate vote channel
	go func() {
		for {
			msg := <-curValidator.delegateVoteRequestChannel
			io.WriteString(conn, "Received delegate vote requests\n")
			validatorsCopy := make([]*Validator, len(validators))
			copy(validatorsCopy, validators)

			sort.Slice(validatorsCopy, func(i, j int) bool {
				return validatorsCopy[i].reputation > validatorsCopy[j].reputation
			})
			voteMsg := DelegateVoteMessage{
				delegateVotes: validatorsCopy[:msg.delegateSize],
			}

			curValidator.delegateVoteChannel <- voteMsg

		}
	}()

	//listen for messages in communication channel
	for {
		msg := <-curValidator.incomingChannel
		switch msg := msg.(type) {
		//Receiving block to validate
		case ValidateBlockMessage:
			io.WriteString(conn, "Received a Block to validate\n")
			isValid := isBlockValid(msg.newBlock)
			validationStatusMessage := ValidationStatusMessage{
				isValid: isValid,
			}
			curValidator.outgoingChannel <- validationStatusMessage
		//Receiving verified transactions
		case VerifiedBlockMessage:
			io.WriteString(conn, "Received verified transaction\n")
			//put verified transactions into confirmed slice for validator
			curValidator.transactionPoolLock.Lock()
			for _, transaction := range msg.transactions {
				curValidator.confirmedTransactions[transaction.ID] = true
			}

			//take transactions out of unconfirmed map
			for _, transaction := range msg.transactions {
				delete(curValidator.unconfirmedTransactions, transaction.ID)
			}
			curValidator.transactionPoolLock.Unlock()

			//add new block to blockchain and move currBlockchain pointer
			curValidator.CurrBlockchain.NextBlock = msg.newBlock
			msg.newBlock.PrevBlock = curValidator.CurrBlockchain
			curValidator.CurrBlockchain = msg.newBlock
		case ConsensusMessage:
			io.WriteString(conn, "Received new consensus state\n")

			//Update blockchain
			curValidator.CurrBlockchain = msg.CurrBlockchain
			curValidator.unconfirmedTransactions = msg.unconfirmedTransactions
			curValidator.confirmedTransactions = msg.confirmedTransactions

		default:
			io.WriteString(conn, "Received an unknown struct: %+v\n")
		}
	}

}

package pos

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gonum.org/v1/gonum/stat/sampleuv"
)

// Blockchain is a series of validated Blocks
var CertifiedBlockchain []Block

//Temporary blocks if we want to look into finality attacks
// var tempChain []Block

// Slice of validator pointers
var validators = make([]*Validator, 0)

// Slice of user pointers who can make transactions
var users = make(map[string]*User)

// Current block proposer
var proposer *Validator = nil

// Validators that will validate the proposed block
var validationCommittee = make([]*Validator, 0)

// Malicious validators
var malValidators = make([]*Validator, 0)

var validatorsSliceLock = &sync.Mutex{}

var committeeSize = 0

var runConsensusCounter = 0

func Run(runType string, numValidators int, numUsers int, numMal int, comSize int, attack string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	committeeSize = comSize

	// create genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{Index: 0, Timestamp: t.String(), Transactions: []Transaction{}, Hash: calculateBlockHash(genesisBlock), PrevHash: "", Validator: ""}
	CertifiedBlockchain = append(CertifiedBlockchain, genesisBlock)

	tcpPort := os.Getenv("PORT")

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server Listening on port :", tcpPort)
	defer server.Close()

	//Advances time slots, choosing new proposers that add blocks to the chain and new validation committees
	go func() {
		for {
			nextTimeSlot()
		}
	}()

	//auto create connections
	go func() {
		if runType == "manual" {
			return
		}
		for numValidators > 0 {
			conn, err := net.Dial("tcp", ":9000")
			if err != nil {
				fmt.Println("Error connecting:", err)
				return
			}
			malString := "n"
			if numMal > 0 {
				malString = "y"
				numMal--
			}
			go handleConnection(conn, runType, "v", malString)
			numValidators--
		}
		for numUsers > 0 {
			conn, err := net.Dial("tcp", ":9000")
			if err != nil {
				fmt.Println("Error connecting:", err)
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			go handleConnection(conn, runType, "u", "")
			numUsers--
		}
	}()

	//Accepts connections joining the network
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, runType, "", "")
	}

}

func chooseValidationCommittee(validators []*Validator, committeeSize int) []*Validator {
	//make a slice of stakes for weighted dsitribution
	validatorsSliceLock.Lock()
	stakeWeights := make([]float64, len(validators))
	for i, validator := range validators {
		stakeWeights[i] = validator.Stake
	}
	validatorsSliceLock.Unlock()

	validationCommittee := make([]*Validator, 0)
	weightedDist := sampleuv.NewWeighted(stakeWeights, nil)
	for i := 0; i < committeeSize; i++ {
		index, isOk := weightedDist.Take()
		if isOk {
			validationCommittee = append(validationCommittee, validators[index])
		} else {
			break
		}
	}
	return validationCommittee
}

func chooseBlockProposer() *Validator {
	if len(validationCommittee) == 0 {
		return nil
	}

	totalWeight := 0.0
	for _, validator := range validationCommittee {
		totalWeight += validator.Stake
	}

	randomNumber := 0.0
	rand.Seed(time.Now().UnixNano())
	randomNumber = rand.Float64() * totalWeight

	weightSum := 0.0
	for _, validator := range validationCommittee {
		weightSum += validator.Stake
		if weightSum >= randomNumber {
			return validator
		}
	}
	return nil
}

func printInfo() {
	//prints blockchain
	printString := ""
	for _, block := range CertifiedBlockchain {
		printString += "->["
		for _, transaction := range block.Transactions {
			printString += fmt.Sprintf("%d,", transaction.ID)
		}
		printString = printString[:len(printString)-1]
		printString += "]"
	}
	printString = printString[1:]
	println("BLOCKCHAIN")
	println(printString)

	//prints User balances
	// println("User balances")
	// for user := range users {
	// 	fmt.Printf("%s: %f\n", users[user].Name, users[user].Balance)
	// }

	//prints Validator balances
	// println("Validator balances")
	// for _, validator := range validators {
	// 	fmt.Printf("%s: %f, %d\n", validator.Address[:3], validator.Stake, validator.committeeCount)
	// }
}

func nextTimeSlot() {
	//wait 5 seconds every slot
	time.Sleep(5 * time.Second)
	fmt.Printf("\nTime slot %s\n\n", time.Now().Format("15:04:05"))
	runConsensusCounter += 1

	if runConsensusCounter >= 5 {
		longestChainConsensus()
		runConsensusCounter = 0
	}

	//randomly choose new committee of a third of all validators who will validate the new block
	validationCommittee = chooseValidationCommittee(validators, committeeSize)
	fmt.Println("New validation committee chosen")
	for _, commit := range validationCommittee {
		commit.committeeCount += 1
		fmt.Println(commit.Address[:3])
	}
	//Choose a new block proposer based on stake
	proposer = chooseBlockProposer()
	if proposer == nil {
		return
	}
	proposer.proposerCount += 1
	fmt.Printf("Proposer %s chosen as new block proposer\n", proposer.Address[:3])

	//block proposer chooses a new block

	// oldBlock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(proposer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Block %d chosen as new block\n", newBlock.Index)

	//validation committee validates blocks
	//broadcast block to all members of committee
	for _, validator := range validationCommittee {
		msg := ValidateBlockMessage{
			newBlock: newBlock,
		}
		validator.incomingChannel <- msg
	}

	// Process validation results
	validCount := 0
	invalidCount := 0
	validationResults := make(map[string]bool)
	for _, validator := range validationCommittee {
		msg := <-validator.outgoingChannel
		switch msg := msg.(type) { // Use type assertion to determine the type of the received message
		case ValidationStatusMessage:
			validationResults[validator.Address] = msg.isValid
			if msg.isValid == true {
				validCount++
			} else {
				invalidCount++
			}
		default:
			fmt.Printf("Received an unknown struct: %+v\n", msg)
			fmt.Printf("%T\n", msg)
		}
	}
	fmt.Printf("Voting results\nInvalid Count: %d\nValid Count: %d\nCommittee size: %d\n", invalidCount, validCount, len(validationCommittee))

	//add block if majority believe block is valid
	isValid := validCount >= len(validationCommittee)/2
	if isValid {
		// proposer.Blockchain = append(proposer.Blockchain, newBlock)
		println("Valid block added to blockchain")

		//broadcast the verified transactions to all blocks
		msg := VerifiedBlockMessage{
			transactions: newBlock.Transactions,
			newBlock:     newBlock,
		}
		for _, validator := range validators {
			validator.incomingChannel <- msg
		}

		//Update transactional amounts and reward proposer
		for _, transaction := range newBlock.Transactions {
			transaction.Sender.Balance -= (transaction.Amount + transaction.Reward)
			transaction.Receiver.Balance += transaction.Amount
			proposer.Stake += transaction.Reward

			senderString := fmt.Sprintf("New balance: %f\n", transaction.Sender.Balance)
			io.WriteString(transaction.Sender.conn, senderString)

			receiverString := fmt.Sprintf("New balance: %f\n", transaction.Receiver.Balance)
			io.WriteString(transaction.Receiver.conn, receiverString)
		}
	} else {
		println("Committee votes block invalid")
	}
	//punish validators who voted against the majority
	slashPercentage := 0.2
	for _, validator := range validationCommittee {
		if isValid {
			if validationResults[validator.Address] == false {
				validator.Stake *= slashPercentage
			}
		} else {
			if validationResults[validator.Address] == true {
				validator.Stake *= slashPercentage
			}
		}
	}

	printInfo()

}

func handleConnection(conn net.Conn, runType string, connectionType string, malString string) {
	defer conn.Close()

	//Determine user or validator connection
	io.WriteString(conn, "Is this node a user or validator (u/v)\n")
	scannedType := bufio.NewScanner(conn)
	if runType == "auto" {
		scannedType = bufio.NewScanner(strings.NewReader(connectionType))
	}
	for scannedType.Scan() {
		if scannedType.Text() == "u" {
			handleUserConnection(conn, runType)
		} else if scannedType.Text() == "v" {
			handleValidatorConnection(conn, runType, malString)
		} else {
			fmt.Printf("%s is not a valid response\n Please enter 'u' or 'v' ", scannedType.Text())
		}
		break
	}
}

func longestChainConsensus() {
	longestLength := -1
	var longestValidator *Validator = nil
	for _, validator := range validators {
		if len(validator.Blockchain) > longestLength {
			longestValidator = validator
			longestLength = len(validator.Blockchain)
		}
	}

	CertifiedBlockchain = make([]Block, len(longestValidator.Blockchain))
	copy(CertifiedBlockchain, longestValidator.Blockchain)

	for _, validator := range validators {
		//broadcast the verified transactions to all blocks
		if validator.Address == longestValidator.Address {
			continue
		}
		msg := ConsensusMessage{
			blockchain:              longestValidator.Blockchain,
			unconfirmedTransactions: longestValidator.unconfirmedTransactions,
			confirmedTransactions:   longestValidator.confirmedTransactions,
		}
		validator.incomingChannel <- msg
	}
}

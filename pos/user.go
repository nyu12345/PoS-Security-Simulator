package pos

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	mathrand "math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type User struct {
	conn        net.Conn
	commChannel chan interface{}
	Name        string
	Address     string
	Balance     float64
	PublicKey   *rsa.PublicKey
	privateKey  *rsa.PrivateKey
}

type Transaction struct {
	ID        int
	Sender    *User
	Receiver  *User
	Signature string
	Amount    float64
	Reward    float64
}

var transactionID = 0

var userID = 0

func generateTransaction(index int, sender *User, receiver *User, amount float64, reward float64) Transaction {
	transaction := Transaction{
		ID:       index,
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
		Reward:   reward,
	}
	signTransaction(&transaction, sender.privateKey)
	return transaction
}

func signTransaction(t *Transaction, privateKey *rsa.PrivateKey) error {
	// Concatenate the transaction data into a single string
	data := fmt.Sprintf("%d%p%p%f%f", t.ID, t.Sender, t.Receiver, t.Amount, t.Reward)

	// Hash the data using SHA256
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	// Sign the hashed data using the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return err
	}

	// Encode the signature as a base64 string
	t.Signature = base64.StdEncoding.EncodeToString(signature)

	return nil
}

func handleUserConnection(conn net.Conn, runType string) {
	defer conn.Close()

	//Enter initial stake and whether or not validator is malicious
	io.WriteString(conn, "Enter initial token balance:\n")
	scannedBalance := bufio.NewScanner(conn)
	if runType == "auto" {
		randomBalance := 0.0
		mathrand.Seed(time.Now().UnixNano())
		randomBalance = mathrand.Float64()*1000 + 10
		randomBalanceString := fmt.Sprintf("%f", randomBalance)
		scannedBalance = bufio.NewScanner(strings.NewReader(randomBalanceString))
	}
	balance := 0.0
	var err error
	for scannedBalance.Scan() {
		balance, err = strconv.ParseFloat(scannedBalance.Text(), 64)
		if err != nil {
			io.WriteString(conn, scannedBalance.Text()+" not a number")
			return
		}
		break
	}

	//Enter name
	io.WriteString(conn, "Enter user name:\n")
	scannedName := bufio.NewScanner(conn)
	name := ""
	if runType == "auto" {
		mutex.Lock()
		name = fmt.Sprintf("user%d", userID)
		userID++
		mutex.Unlock()
		scannedName = bufio.NewScanner(strings.NewReader(name))
	}
	for scannedName.Scan() {
		name = scannedName.Text()
		if _, ok := users[name]; ok {
			fmt.Printf("Name: %s already taken: \n", name)
		} else {
			break
		}
	}

	//Calculate address based on time
	t := time.Now()
	address := calculateHash(t.String())

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}

	publicKey := &privateKey.PublicKey

	//Instantiate new validator
	curUser := &User{
		conn:        conn,
		commChannel: make(chan interface{}),
		Name:        name,
		Address:     address,
		Balance:     float64(balance),
		privateKey:  privateKey,
		PublicKey:   publicKey,
	}

	users[name] = curUser

	fmt.Printf("new user count: %d\n", len(users))

	for {
		io.WriteString(conn, "Starting new transaction\n")
		io.WriteString(conn, "Enter receiver name:\n")
		scannedReceiver := bufio.NewScanner(conn)
		receiverName := ""

		if runType == "auto" {
			mutex.Lock()
			mathrand.Seed(time.Now().UnixNano())
			randomIndex := 0
			if len(users)-1 > 0 {
				randomIndex = mathrand.Intn(len(users) - 1)
			}
			counter := 0
			randomUser := ""
			for userName := range users {
				if counter == randomIndex {
					randomUser = userName
					break
				}
				counter++
			}
			scannedReceiver = bufio.NewScanner(strings.NewReader(randomUser))
			mutex.Unlock()
		}

		for scannedReceiver.Scan() {
			receiverName = scannedReceiver.Text()
			break
		}

		io.WriteString(conn, "Enter transaction amount:\n")
		scannedAmount := bufio.NewScanner(conn)
		if runType == "auto" {
			randomAmount := 0.0
			mathrand.Seed(time.Now().UnixNano())
			randomAmount = mathrand.Float64()*100 + 1
			randomAmountString := fmt.Sprintf("%f", randomAmount)
			scannedAmount = bufio.NewScanner(strings.NewReader(randomAmountString))
		}
		amount := 0.0
		for scannedAmount.Scan() {
			amount, err = strconv.ParseFloat(scannedAmount.Text(), 64)
			if err != nil {
				io.WriteString(conn, scannedAmount.Text()+" not a number")
				return
			}
			break
		}

		io.WriteString(conn, "Enter transaction reward:\n")
		scannedReward := bufio.NewScanner(conn)
		if runType == "auto" {
			randomReward := 0.0
			mathrand.Seed(time.Now().UnixNano())
			randomReward = mathrand.Float64()*5 + 0
			randomRewardString := fmt.Sprintf("%f", randomReward)
			scannedReward = bufio.NewScanner(strings.NewReader(randomRewardString))
		}
		reward := 0.0
		for scannedReward.Scan() {
			reward, err = strconv.ParseFloat(scannedReward.Text(), 64)
			if err != nil {
				io.WriteString(conn, scannedReward.Text()+" not a number")
				return
			}
			break
		}

		mutex.Lock()
		curTransactionID := transactionID
		transactionID++
		mutex.Unlock()

		curTransaction := &Transaction{
			ID:       curTransactionID,
			Sender:   users[curUser.Name],
			Receiver: users[receiverName],
			Amount:   amount,
			Reward:   reward,
		}

		//Broadcast current transaction to all validators
		mutex.Lock()
		validatorsCopy := validators
		mutex.Unlock()
		transactionString := fmt.Sprintf("Sent transaction %d\n", curTransaction.ID)
		io.WriteString(conn, transactionString)
		for _, validator := range validatorsCopy {
			msg := NewTransactionMessage{
				transaction: *curTransaction,
			}
			validator.commChannel <- msg
		}
		time.Sleep(5 * time.Second)
	}

}

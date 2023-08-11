# Testing Security for Proof of Stake Implementations

## Abstract

Here is a link to our paper where you can read about background, methodology, and results: [Testing Security for Proof of Stake Implementations](CS_512_Final_Paper.pdf)

Security is a fundamental issue in blockchains, as it ensures
the integrity of the blockchain, the accuracy of transactions, and the trust of users. This paper presents the design and implementation of our Proof of Stake (PoS) blockchain simulation to test different security strategies against various attacks. PoS is a consensus mechanism used in blockchain networks as an alternative to the traditional Proof of Work (PoW) mechanism. In PoS, validators are selected to validate blocks and secure the network based on the amount of cryptocurrency they hold and "stake" as collateral. A common malicious behavior is double spending, where the same funds are used for multiple different transactions. Attacks on the blockchain that give the dishonest user opportunities to conduct fraudulent activities like double spending compromise the accuracy of transactions. Our simulation supports a few different methods of enhancing blockchain security, including the standard PoS consensus,slashing of misbehaving validators, and reputation based PoS. There are a couple of attacks that can be simulated: network partition and balance attack. We report the measurements about how the different strategies for blockchain security perform against our simulated attacks.

## Installation
Run `go get ./...`

## Simulation

To run our simulation, run `go run main.go`. There are a variety of parameters for the simulation that can be changed manually in the `main()` function of `main.go`. These include

- runType
    - either "auto" or "manual"
    - "auto" will automatically generate validators, users, malicious nodes, and transactions to operate the blockchain
    - "manual" will let you do this manually, disregarding the following parameters
- numValidators
    - The number of validator nodes
- numUsers
    - The number of users making transactions
- numMal
    - The number of malicious nodes
- committeeSize
    - The size of the committee confirming a block into the blockchain
- delegateSize
    - The delegate committee size for reputation proof of stake blockchain
- blockchainType
    - "pos" - A generic proof of stake blockchain
    - "slashing" - the generic proof of stake blockchain with stake slashing punishments
    - "reputation" - A delegated proof of stake blockchain with elected delegates
- attack
    - Either "network_partition" or "balance" attack types

### auto

Simply run `go run main.go` and it will instantiate all validators and users in the blockchain while randomly generating transactions. The state of the blockchain will be printed out every few seconds showing new confirmed transactions in the blockchain and results of elections and block proposals

### manual

Run `go run main.go` which will start listening for incoming server connection requests on port 9000

Next open another terminal and open a connection using `nc localhost 9000`. You will then be asked to input 'u' for creating a user or 'v' for creating a validator. Each time you want to create a user of validator you must open another connection.

When creating validators you must enter token stake and malicious status

When creating users you must enter their name, and balance, then you will continously be prompted to create new transactions

As you make transactions between your different users in their different terminals the state of the blockchain will be output in the terminal of the original listening global server.

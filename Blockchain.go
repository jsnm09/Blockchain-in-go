package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Block struct {
	Index        int    `json:"index"`
	Timestamp    string `json:"timestamp"`
	Data         string `json:"data"`
	PreviousHash string `json:"previousHash"`
	Hash         string `json:"hash"`
	Nonce        int    `json:"nonce"`
}

type Message struct {
	Data string `json:"data"`
}

var Blockchain []Block
var mutex = &sync.Mutex{}

const difficulty = 2

// function to generate a SHA-256 hash for a block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%d",
		block.Index, block.Timestamp, block.Data, block.PreviousHash, block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// function to create new block based on previous block's hash
func generateBlock(oldBlock Block, data string) (Block, error) {
	newBlock := Block{
		Index:        oldBlock.Index + 1,
		Timestamp:    time.Now().Format(time.RFC3339),
		Data:         data,
		PreviousHash: oldBlock.Hash,
		Hash:         "",
		Nonce:        0,
	}

	// Proof of Work
	for {
		newBlock.Hash = calculateHash(newBlock)

		//check if hash satisfies difficulty
		if strings.HasPrefix(newBlock.Hash, strings.Repeat("0", difficulty)) {
			fmt.Printf("Block mined! Hash: %s\n", newBlock.Hash)
			break
		} else {
			newBlock.Nonce++
		}
	}

	return newBlock, nil
}

// bool function to check if new block is valid based on previous blocks
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	// Verify proof of work
	if !strings.HasPrefix(newBlock.Hash, strings.Repeat("0", difficulty)) {
		return false
	}
	return true
}

func runServer() {
	// Get the blockchain
	http.HandleFunc("/blockchain", getBlockchain)

	// Add new block
	http.HandleFunc("/mine", mineBlock)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Lock mutex to safely read blockchain
	mutex.Lock()
	defer mutex.Unlock()

	// Encode and return the blockchain
	json.NewEncoder(w).Encode(Blockchain)
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Parse message
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		http.Error(w, "Error parsing request JSON", http.StatusBadRequest)
		return
	}

	// Lock mutex before modifying blockchain
	mutex.Lock()
	defer mutex.Unlock()

	// Generate new block
	oldBlock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(oldBlock, message.Data)
	if err != nil {
		http.Error(w, "Error generating block", http.StatusInternalServerError)
		return
	}

	// Validate new block
	if isBlockValid(newBlock, oldBlock) {
		Blockchain = append(Blockchain, newBlock)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newBlock)
	} else {
		http.Error(w, "Invalid block", http.StatusBadRequest)
	}
}

func main() {
	// Initialize mutex
	mutex = &sync.Mutex{}

	// Create genesis block
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Format(time.RFC3339),
		Data:         "Genesis Block",
		PreviousHash: "",
		Hash:         "",
		Nonce:        0,
	}

	// Calculate hash for genesis block
	genesisBlock.Hash = calculateHash(genesisBlock)

	// Initialize blockchain with genesis block
	Blockchain = append(Blockchain, genesisBlock)

	// Start the HTTP server
	runServer()
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/triedb"
)

// CheckErr simply checks errors and panics, as a guilty pleasure.
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

// MakeDirIfNot handles dir creation operations
func MakeDirIfNot(dir string) bool {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, 0750)
		CheckErr(err)
		return true
	}
	return false
}

func NewTrieDB(N int) *trie.Trie {
	var trieInstance *trie.Trie
	dir := fmt.Sprintf("./data-%d", N)
	created := MakeDirIfNot(dir)
	db, err := rawdb.NewLevelDBDatabase(dir, 16, 16, "leveldb", false)
	if err != nil {
		log.Fatalf("Failed to create LevelDB database: %v", err)
	}
	diskdb := triedb.NewDatabase(db, triedb.HashDefaults)
	trieInstance = trie.NewEmpty(diskdb)
	if created {
		// trieInstance = trie.NewEmpty(diskdb)
		fmt.Println("Creating new trie")
	} else {

	}
	// rootHash := storeNewRootHash(N, trieInstance)
	// fmt.Printf("Trie hash: %x\n", rootHash)
	// // print the root of the trie
	// fmt.Printf("Trie root: %x\n", rootHash)
	return trieInstance
}

func main() {
	// Initialize the database
	db := rawdb.NewMemoryDatabase()

	// Create a new trie database
	triedb := triedb.NewDatabase(db, triedb.HashDefaults)

	// Initialize the state with an empty root
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(triedb), nil)
	if err != nil {
		log.Fatalf("Failed to create state: %v", err)
	}

	// Update the state with some data
	address := common.HexToAddress("0x123456")
	stateDB.AddBalance(address, 1000)

	// Commit the state to the trie
	root, err := stateDB.Commit(false)
	if err != nil {
		log.Fatalf("Failed to commit state: %v", err)
	}

	// Persist the trie to the database
	if err := triedb.Commit(root, false); err != nil {
		log.Fatalf("Failed to commit trie: %v", err)
	}

	fmt.Printf("State committed with root: %x\n", root)
}

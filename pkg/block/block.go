package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

/*
# Block struct

	Index		- The position of the data record in the blockchain
	Timestamp	- Automatically determined, it is the time when the data is written
	BPM			- The BPM (beats per minute) in the blockchain (Like the transaction information contained in the blockchain)
	Hash		- The SHA256 identifier representing this data record
	PrevHash	- The SHA256 identifier of the previous record in the chain
	Difficulty  - Mining difficulty
	Nonce       - Random number
*/
type Block struct {
	Index      int
	Timestamp  string
	BPM        int
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

var (
	Blockchain []Block
	Mutex      = &sync.Mutex{}
)

// init create genesis block
func init() {
	t := time.Now()
	genesisBlock := Block{Index: 0, Timestamp: t.String(), BPM: 0, Hash: "", PrevHash: ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)
}

// difficulty is a constant that defines the number of 0s we want leading the hash. The more zeros we have to get, the harder it is to find the correct hash.
func getDifficulty() int {
	difficulty, err := strconv.Atoi(os.Getenv("Difficulty"))
	if err != nil {
		difficulty = 0
	}
	return difficulty
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)
	newBlock.Difficulty = getDifficulty()

	// mining algorithm
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !IsHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(CalculateHash(newBlock), "do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(CalculateHash(newBlock), "wrok done!")
			newBlock.Hash = CalculateHash(newBlock)
			break
		}
	}

	return newBlock, nil
}

// make sure the chain we're checking is longer than the current blockchain
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

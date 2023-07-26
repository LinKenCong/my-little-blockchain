package main

import (
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	"github.com/LinKenCong/my-little-blockchain/pkg/block"
	"github.com/LinKenCong/my-little-blockchain/pkg/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := block.Block{}
		genesisBlock = block.Block{0, t.String(), 0, block.CalculateHash(genesisBlock), ""}
		spew.Dump(genesisBlock)
		block.Mutex.Lock()
		block.Blockchain = append(block.Blockchain, genesisBlock)
		block.Mutex.Unlock()
	}()

	log.Fatal(http.Run())
}

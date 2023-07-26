package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"

	"github.com/LinKenCong/my-little-blockchain/pkg/block"
)

// bcServer handles incoming concurrent Blocks
var bcServer chan []block.Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []block.Block)

	// create genesis block
	t := time.Now()
	genesisBlock := block.Block{Index: 0, Timestamp: t.String(), BPM: 0, Hash: "", PrevHash: ""}
	spew.Dump(genesisBlock)
	block.Blockchain = append(block.Blockchain, genesisBlock)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a new BPM:")

	scanner := bufio.NewScanner(conn)

	// take in BPM from stdin and add it to blockchain after conducting necessary validation
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}
			newBlock, err := block.GenerateBlock(block.Blockchain[len(block.Blockchain)-1], bpm)
			if err != nil {
				log.Println(err)
				continue
			}
			if block.IsBlockValid(newBlock, block.Blockchain[len(block.Blockchain)-1]) {
				newBlockchain := append(block.Blockchain, newBlock)
				block.ReplaceChain(newBlockchain)
			}

			bcServer <- block.Blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
	}()

	// simulate receiving broadcast
	go func() {
		for {
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(block.Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string("\n"))
			io.WriteString(conn, string(output))
			io.WriteString(conn, string("\nEnter a new BPM:"))
		}
	}()

	for range bcServer {
		spew.Dump(block.Blockchain)
	}
}

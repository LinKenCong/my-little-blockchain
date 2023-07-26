package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/LinKenCong/my-little-blockchain/pkg/block"
	"github.com/LinKenCong/my-little-blockchain/pkg/utils"
)

type Message struct {
	BPM int
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(block.Blockchain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var msg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		utils.RespondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	block.Mutex.Lock()
	prevBlock := block.Blockchain[len(block.Blockchain)-1]
	newBlock := block.GenerateBlock(prevBlock, msg.BPM)

	if block.IsBlockValid(newBlock, prevBlock) {
		block.Blockchain = append(block.Blockchain, newBlock)
	}
	block.Mutex.Unlock()

	utils.RespondWithJSON(w, r, http.StatusCreated, newBlock)
}

package http

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LinKenCong/my-little-blockchain/pkg/handles"
	"github.com/gorilla/mux"
)

func Run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handles.HandleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handles.HandleWriteBlock).Methods("POST")
	return muxRouter
}

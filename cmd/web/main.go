package main

import (
	"log"
	"net/http"
	"peer2peerchat/internals/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting a channel listener")
	go handlers.ListenForWsChannel()

	log.Println("Starting server on :7070")
	http.ListenAndServe(":7070", mux)
}

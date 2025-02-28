package main

import (
	"log"
	"net/http"

	"github.com/ipfs/kubo/client/rpc"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	node, err := rpc.NewLocalApi()
	if err != nil {
		log.Fatalf("Failed to create IPFS client: %v", err)
	}

	http.HandleFunc("/ask", ask(config, node))

	serverAddr := "localhost:9123"
	log.Printf("Starting server on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

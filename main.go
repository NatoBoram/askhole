package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	multi, err := multiaddr.NewMultiaddr(env.Multiaddr)
	if err != nil {
		log.Fatalf("Failed to parse Kubo address: %v", err)
	}

	kubo, err := rpc.NewApi(multi)
	if err != nil {
		log.Fatalf("Failed to create Kubo client: %v", err)
	}

	http.HandleFunc("/ask", ask(env, kubo))

	port := strconv.FormatUint(uint64(env.Port), 10)
	serverAddr := "localhost:" + port
	log.Printf("Starting server on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

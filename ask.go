package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/ipfs/boxo/path"
	"github.com/ipfs/kubo/client/rpc"
)

func askIpfs(w http.ResponseWriter, r *http.Request, kubo *rpc.HttpApi, hash string) {
	path, err := path.NewPath("/ipfs/" + hash)
	if err != nil {
		log.Printf("Failed to create path from hash %q: %v\n", hash, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	channel, err := kubo.Routing().FindProviders(ctx, path)
	if err != nil {
		log.Printf("Failed to find providers for %q: %v\n", hash, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	select {
	case provider, ok := <-channel:
		if !ok {
			log.Printf("No providers found for %q\n", hash)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Found provider %q for %q\n", provider.ID, hash)
		w.WriteHeader(http.StatusOK)
		return

	case <-ctx.Done():
		log.Printf("Context cancelled while finding providers for %q\n", hash)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func askIpns(w http.ResponseWriter, r *http.Request, kubo *rpc.HttpApi, name string) {
	path, err := path.NewPath("/ipns/" + name)
	if err != nil {
		log.Printf("Failed to create path from name %q: %v\n", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	channel, err := kubo.Routing().FindProviders(ctx, path)
	if err != nil {
		log.Printf("Failed to find providers for %q: %v\n", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	select {
	case provider, ok := <-channel:
		if !ok {
			log.Printf("No providers found for %q\n", name)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("Found provider %q for %q\n", provider.ID, name)
		w.WriteHeader(http.StatusOK)
		return
	case <-ctx.Done():
		log.Printf("Context cancelled while finding providers for %q\n", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ask(config *Env, kubo *rpc.HttpApi) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			log.Printf("Received a request with no domain\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !strings.HasSuffix(domain, config.Domain) {
			log.Printf("Domain %q does not match base domain %q\n", domain, config.Domain)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if domain == config.Domain {
			log.Printf("Domain matches base domain exactly, no subdomain present\n")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		remaining := strings.TrimSuffix(domain, "."+config.Domain)
		parts := strings.Split(remaining, ".")

		if len(parts) < 2 {
			log.Printf("Invalid subdomain format: %s\n", domain)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lastPart := parts[len(parts)-1]
		switch lastPart {
		case "ipfs":
			if len(parts) != 2 {
				log.Printf("Invalid IPFS subdomain format: %s\n", domain)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hash := parts[0]
			log.Printf("Finding providers for IPFS %q...\n", hash)
			askIpfs(w, r, kubo, hash)
			return

		case "ipns":
			if len(parts) < 2 {
				log.Printf("Invalid IPNS subdomain format: %s\n", domain)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			name := strings.Join(parts[:len(parts)-1], ".")
			log.Printf("Finding providers for IPNS %q...\n", name)
			askIpns(w, r, kubo, name)
			return

		default:
			log.Printf("Invalid protocol %q; expected ipfs or ipns\n", lastPart)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

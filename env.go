package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

type Env struct {
	Domain    string
	Multiaddr string
	Port      uint16
}

func getEnvironment() string {
	environment := os.Getenv("GO_ENV")
	if environment != "" {
		return environment
	}

	if testing.Testing() {
		os.Setenv("GO_ENV", "test")
		return "test"
	}

	os.Setenv("GO_ENV", "development")
	return "development"
}

func loadEnv() (*Env, error) {
	environment := getEnvironment()

	files := []string{
		".env." + environment + ".local",
		".env." + environment,
		".env.local",
		".env",
	}

	for _, file := range files {
		err := godotenv.Load(file)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load environment variables from %q: %w", file, err)
		}
	}

	domain := os.Getenv("KUBO_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("KUBO_DOMAIN is required")
	}

	multiaddr := os.Getenv("KUBO_MULTIADDR")
	if multiaddr == "" {
		return nil, fmt.Errorf("KUBO_MULTIADDR is required")
	}

	sPort := os.Getenv("PORT")
	if sPort == "" {
		return nil, fmt.Errorf("PORT is required")
	}

	uPort, err := strconv.ParseUint(sPort, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PORT: %w", err)
	}

	return &Env{
		Domain:    os.Getenv("KUBO_DOMAIN"),
		Multiaddr: os.Getenv("KUBO_MULTIADDR"),
		Port:      uint16(uPort),
	}, nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bur4kbey/go-ddns/internal/crypto"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development). 
	// In Docker, these will be set via environment variables and this will fail silently.
	_ = godotenv.Load()

	domain := flag.String("domain", "", "The base public domain (e.g., example.com)")
	key := flag.String("key", "", "The TXT record key/subdomain (e.g., home-node)")
	secret := flag.String("secret", "", "The decryption key")
	out := flag.String("out", "ip.txt", "The output file path")
	keepAlive := flag.Bool("keep-alive", false, "Keep running and check every 5 minutes")

	flag.Parse()

	if *domain == "" || *key == "" || *secret == "" {
		// Read from ENV as fallback for Docker usage
		if *domain == "" {
			*domain = os.Getenv("GO_DDNS_DOMAIN")
		}
		if *key == "" {
			*key = os.Getenv("GO_DDNS_KEY")
		}
		if *secret == "" {
			*secret = os.Getenv("GO_DDNS_SECRET")
		}
	}

	if *domain == "" || *key == "" || *secret == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	recordName := fmt.Sprintf("%s.%s", *key, *domain)
	if *key == "@" {
		recordName = *domain
	}

	var lastTXT string

	run(recordName, *secret, *out, &lastTXT)

	if *keepAlive {
		log.Printf("Entering keep-alive mode. Checking every 5 minutes.")
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			run(recordName, *secret, *out, &lastTXT)
		}
	}
}

func run(recordName, secret, out string, lastTXT *string) {
	// Look up TXT record
	txts, err := net.LookupTXT(recordName)
	if err != nil {
		log.Printf("Error looking up TXT record for %s: %v", recordName, err)
		return
	}

	if len(txts) == 0 {
		log.Printf("No TXT records found for %s", recordName)
		return
	}

	// We assume the first TXT record is the one we want. Cloudflare handles multiple,
	// but our server updates a specific one. Let's just pick the first.
	// Sometimes TXT records are split into multiple strings by DNS protocol, net.LookupTXT handles concatenation.
	encryptedIP := strings.Join(txts, "")

	if encryptedIP == *lastTXT {
		// No change
		log.Printf("TXT record hasn't changed. Skipping decryption.")
		return
	}

	ip, err := crypto.Decrypt(encryptedIP, secret)
	if err != nil {
		log.Printf("Error decrypting IP: %v", err)
		return
	}

	// Write to file
	err = os.WriteFile(out, []byte(ip), 0644)
	if err != nil {
		log.Printf("Error writing IP to file %s: %v", out, err)
		return
	}

	log.Printf("Successfully decrypted IP (%s) and wrote to %s", ip, out)
	*lastTXT = encryptedIP
}

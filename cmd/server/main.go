package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bur4kbey/go-ddns/internal/cloudflare"
	"github.com/bur4kbey/go-ddns/internal/crypto"
	"github.com/bur4kbey/go-ddns/internal/ipfetcher"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development). 
	// In Docker, these will be set via environment variables and this will fail silently.
	_ = godotenv.Load()

	secret := os.Getenv("GO_DDNS_SECRET")
	cfToken := os.Getenv("GO_DDNS_CLOUDFLARE_TOKEN")
	zoneID := os.Getenv("GO_DDNS_ZONE_ID")
	domain := os.Getenv("GO_DDNS_DOMAIN")
	key := os.Getenv("GO_DDNS_KEY")

	if secret == "" || cfToken == "" || zoneID == "" || domain == "" || key == "" {
		log.Fatal("Missing required environment variables: GO_DDNS_SECRET, GO_DDNS_CLOUDFLARE_TOKEN, GO_DDNS_ZONE_ID, GO_DDNS_DOMAIN, GO_DDNS_KEY")
	}

	recordName := fmt.Sprintf("%s.%s", key, domain)
	if key == "@" || key == "" {
		recordName = domain
	}

	cfClient, err := cloudflare.NewClient(cfToken)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudflare client: %v", err)
	}

	log.Printf("Starting GO-DDNS server for record: %s", recordName)

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	var lastIP string

	// Run immediately on start
	updateIP(cfClient, zoneID, recordName, secret, &lastIP)

	for range ticker.C {
		updateIP(cfClient, zoneID, recordName, secret, &lastIP)
	}
}

func updateIP(cfClient *cloudflare.Client, zoneID, recordName, secret string, lastIP *string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ip, err := ipfetcher.FetchPublicIP("")
	if err != nil {
		log.Printf("Error fetching public IP: %v", err)
		return
	}

	if ip == *lastIP {
		log.Printf("IP (%s) has not changed. Skipping update.", ip)
		return
	}

	encryptedIP, err := crypto.Encrypt(ip, secret)
	if err != nil {
		log.Printf("Error encrypting IP: %v", err)
		return
	}

	err = cfClient.UpdateTXTRecord(ctx, zoneID, recordName, encryptedIP)
	if err != nil {
		log.Printf("Error updating Cloudflare TXT record: %v", err)
		return
	}

	log.Printf("Successfully updated Cloudflare TXT record with encrypted IP. (Raw IP: %s)", ip)
	*lastIP = ip
}

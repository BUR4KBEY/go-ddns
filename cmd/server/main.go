package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bur4kbey/go-ddns/internal/cloudflare"
	"github.com/bur4kbey/go-ddns/internal/crypto"
	"github.com/bur4kbey/go-ddns/internal/env"
	"github.com/bur4kbey/go-ddns/internal/ipfetcher"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	secret  string
	cfToken string
	zoneID  string
	domain  string
	key     string
)

func main() {
	// Load .env file if it exists (for local development).
	// In Docker, these will be set via environment variables and this will fail silently.
	_ = godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "go-ddns server updates Cloudflare TXT records with your public IP",
		Run: func(cmd *cobra.Command, args []string) {
			config := env.Load()

			if secret == "" {
				secret = config.Secret
			}
			if cfToken == "" {
				cfToken = config.CloudflareToken
			}
			if zoneID == "" {
				zoneID = config.ZoneID
			}
			if domain == "" {
				domain = config.Domain
			}
			if key == "" {
				key = config.Key
			}

			if secret == "" || cfToken == "" || zoneID == "" || domain == "" || key == "" {
				fmt.Println("Usage:")
				_ = cmd.Help()
				os.Exit(1)
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
		},
	}

	rootCmd.Flags().StringVar(&secret, "secret", "", "The encryption secret")
	rootCmd.Flags().StringVar(&cfToken, "cf-token", "", "Cloudflare API token")
	rootCmd.Flags().StringVar(&zoneID, "zone-id", "", "Cloudflare Zone ID")
	rootCmd.Flags().StringVar(&domain, "domain", "", "The base public domain")
	rootCmd.Flags().StringVar(&key, "key", "", "The TXT record key/subdomain")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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

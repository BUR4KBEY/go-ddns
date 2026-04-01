package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bur4kbey/go-ddns/internal/crypto"
	"github.com/bur4kbey/go-ddns/internal/env"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	domain    string
	key       string
	secret    string
	out       string
	keepAlive bool
)

func main() {
	// Load .env file if it exists (for local development).
	// In Docker, these will be set via environment variables and this will fail silently.
	_ = godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "client",
		Short: "go-ddns client reads encrypted IP from Cloudflare TXT record",
		Run: func(cmd *cobra.Command, args []string) {
			if domain == "" || key == "" || secret == "" {
				// Read from ENV as fallback for Docker usage
				config := env.Load()
				if domain == "" {
					domain = config.Domain
				}
				if key == "" {
					key = config.Key
				}
				if secret == "" {
					secret = config.Secret
				}
			}

			if domain == "" || key == "" || secret == "" {
				fmt.Println("Usage:")
				_ = cmd.Help()
				os.Exit(1)
			}

			recordName := fmt.Sprintf("%s.%s", key, domain)
			if key == "@" {
				recordName = domain
			}

			var lastTXT string

			process(recordName, secret, out, &lastTXT)

			if keepAlive {
				log.Printf("Entering keep-alive mode. Checking every 5 minutes.")
				ticker := time.NewTicker(5 * time.Minute)
				defer ticker.Stop()

				for range ticker.C {
					process(recordName, secret, out, &lastTXT)
				}
			}
		},
	}

	rootCmd.Flags().StringVar(&domain, "domain", "", "[Required] The base public domain (e.g., domain.tld)")
	rootCmd.Flags().StringVar(&key, "key", "", "[Required] The TXT record key/subdomain (e.g., _go_ddns)")
	rootCmd.Flags().StringVar(&secret, "secret", "", "[Required] The decryption key")
	rootCmd.Flags().StringVarP(&out, "out", "o", "ip.txt", "The output file path")
	rootCmd.Flags().BoolVarP(&keepAlive, "keep-alive", "d", false, "Keep running and check every 5 minutes")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func process(recordName, secret, out string, lastTXT *string) {
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

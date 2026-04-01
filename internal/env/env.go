package env

import "os"

// Config holds all tracked environment variables used by the application
type Config struct {
	Domain          string
	Key             string
	Secret          string
	CloudflareToken string
	ZoneID          string
}

// Load retrieves all required environment variables and returns them in a Config object
func Load() *Config {
	return &Config{
		Domain:          os.Getenv("GO_DDNS_DOMAIN"),
		Key:             os.Getenv("GO_DDNS_KEY"),
		Secret:          os.Getenv("GO_DDNS_SECRET"),
		CloudflareToken: os.Getenv("GO_DDNS_CLOUDFLARE_TOKEN"),
		ZoneID:          os.Getenv("GO_DDNS_ZONE_ID"),
	}
}

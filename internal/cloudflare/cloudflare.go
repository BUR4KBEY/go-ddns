package cloudflare

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
)

type Client struct {
	api *cf.API
}

// NewClient initializes a new Cloudflare client using an API token.
func NewClient(token string) (*Client, error) {
	api, err := cf.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}
	return &Client{api: api}, nil
}

// UpdateTXTRecord updates or creates a TXT record with the given content.
// `name` should be the full domain name (e.g., "my-ip.example.com").
func (c *Client) UpdateTXTRecord(ctx context.Context, zoneID, name, content string) error {
	rc := cf.ZoneIdentifier(zoneID)
	
	// Cloudflare expects TXT record content to be wrapped in quotation marks.
	quotedContent := fmt.Sprintf(`"%s"`, content)

	records, _, err := c.api.ListDNSRecords(ctx, rc, cf.ListDNSRecordsParams{
		Type: "TXT",
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Name == name {
			// Record exists, check if it needs updating
			if record.Content == quotedContent {
				return nil // No change required
			}

			// Update the record
			_, err = c.api.UpdateDNSRecord(ctx, rc, cf.UpdateDNSRecordParams{
				ID:      record.ID,
				Type:    "TXT",
				Name:    name,
				Content: quotedContent,
				TTL:     300, // Optional: 300 seconds (5 minutes)
			})
			if err != nil {
				return fmt.Errorf("failed to update DNS record: %w", err)
			}
			return nil
		}
	}

	// Record doesn't exist, create it
	_, err = c.api.CreateDNSRecord(ctx, rc, cf.CreateDNSRecordParams{
		Type:    "TXT",
		Name:    name,
		Content: quotedContent,
		TTL:     300,
	})
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	return nil
}

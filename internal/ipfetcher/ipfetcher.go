package ipfetcher

import (
	"io"
	"net/http"
	"strings"
)

// FetchPublicIP retrieves the public IP address from a given URL (e.g., https://api.ipify.org).
func FetchPublicIP(url string) (string, error) {
	if url == "" {
		url = "https://api.ipify.org"
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ipBytes)), nil
}

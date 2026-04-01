# Go DDNS

Go DDNS is a lightweight, secure Dynamic DNS (DDNS) utility for Cloudflare users. It enables you to track your dynamic public IP address by publishing it as an encrypted TXT record on your Cloudflare DNS, making it accessible from anywhere without exposing your real IP to the public.

## How It Works

Go DDNS operates with two components:

1.  **The Reporter (Server):** Runs on the machine whose public IP you want to track. It fetches its own public IP, encrypts it using AES-256-GCM, and updates a Cloudflare TXT record.
2.  **The Retriever (Client):** Runs on any machine where you need the updated IP. It looks up the encrypted TXT record, decrypts it using your shared secret, and saves the IP to a local file.

By using encryption, your public IP remains hidden from anyone who doesn't have your shared secret, even if they see your TXT record.

## Features

- **End-to-End Encryption:** Your IP is encrypted locally before being sent to Cloudflare.
- **Efficient:** Only updates Cloudflare when your public IP actually changes.
- **Docker-Ready:** Minimalist Docker images for easy deployment on any server or NAS.

## Prerequisites

- A Cloudflare account with a managed domain.
- A Cloudflare API Token with `Zone:DNS:Edit` permissions.
- Your Cloudflare **Zone ID**.
- A **Secret Key** for encryption (a 32-byte hex string). You can generate one with:
  ```bash
  openssl rand -hex 32
  ```

## Configuration

The easiest way to configure Go DDNS is through environment variables or a `.env` file:

| Variable                   | Description                                              |
| :------------------------- | :------------------------------------------------------- |
| `GO_DDNS_SECRET`           | Your 32-byte hex encryption secret.                      |
| `GO_DDNS_CLOUDFLARE_TOKEN` | Your Cloudflare API Token (required for Reporter).       |
| `GO_DDNS_ZONE_ID`          | Your Cloudflare Zone ID (required for Reporter).         |
| `GO_DDNS_DOMAIN`           | Your base domain (e.g., `domain.tld`).                   |
| `GO_DDNS_KEY`              | The subdomain/key for the TXT record (e.g., `_go_ddns`). |

## Usage Examples

### Standalone (Using Just)

If you have [just](https://github.com/casey/just) installed, you can build and run easily:

**1. Start the Reporter (Updates Cloudflare):**

```bash
just build
./bin/server --domain "domain.tld" \
  --key "_go_ddns" \
  --secret "$(openssl rand -hex 32)" \
  --cf-token "your_cloudflare_token_here" \
  --zone-id "your_zone_id"
```

**2. Start the Retriever (Gets the IP):**

```bash
just build
./bin/client --domain "domain.tld" \
  --key "_go_ddns" \
  --secret "$(openssl rand -hex 32)" \
  --out "ip.txt" \
  --keep-alive
```

### Docker Compose

You can easily run both or either component using Docker Compose. This is ideal for home servers or remote gateways.

```yaml
services:
  # This machine will report its public IP to Cloudflare
  reporter:
    image: ghcr.io/bur4kbey/go-ddns-server:main
    restart: unless-stopped
    environment:
      GO_DDNS_SECRET= your_32_byte_hex_secret
      GO_DDNS_CLOUDFLARE_TOKEN: your_cloudflare_token
      GO_DDNS_ZONE_ID: your-zone-id
      GO_DDNS_DOMAIN: domain.tld
      GO_DDNS_KEY: _go_ddns

  # This machine will fetch the IP and save it to a local file
  retriever:
    image: ghcr.io/bur4kbey/go-ddns-client:main
    restart: unless-stopped
    environment:
      GO_DDNS_SECRET: your_32_byte_hex_secret
      GO_DDNS_DOMAIN: domain.tld
      GO_DDNS_KEY: _go_ddns
    command: ['--keep-alive', '--out', '/data/ip.txt']
    volumes:
      - ./ip.txt:/data/ip.txt
```

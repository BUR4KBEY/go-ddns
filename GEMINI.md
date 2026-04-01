# Go DDNS Project

A Go-based Dynamic DNS (DDNS) solution that updates Cloudflare DNS records. It consists of a server (monitoring and updating Cloudflare) and a client (reporting the public IP).

## Project Structure

- `cmd/`: Application entry points.
  - `server/`: Server component that interacts with Cloudflare API.
  - `client/`: Client component that fetches and reports public IP.
- `internal/`: Private library code.
  - `cloudflare/`: Cloudflare API client wrapper.
  - `crypto/`: Encryption utilities (AES-256-GCM) for secure IP reporting.
  - `env/`: Environment variable loading and configuration.
  - `ipfetcher/`: Public IP fetching logic.
- `build/package/`: Dockerfiles for client and server.
- `bin/`: Compiled binaries (ignored by git).

## Core Technologies

- **Language:** Go 1.26.1+
- **CLI Framework:** Cobra
- **DNS Provider:** Cloudflare (via `cloudflare-go`)
- **Config:** `.env` via `godotenv`
- **Task Runner:** `just`

## Common Commands (via Justfile)

@./Justfile

## Development Guidelines

- **Environment Variables:** Use `GO_DDNS_*` prefixed variables. See `.env.example`.
- **Encryption:** The client and server share a `GO_DDNS_SECRET` to encrypt/decrypt IP data using AES-256-GCM.
- **Adding Features:** Follow the existing pattern of putting core logic in `internal/` and CLI commands in `cmd/`.
- **Testing:** Add tests in the same directory as the code (e.g., `crypto_test.go`). Run `go test ./...`
- **GEMINI.md:** Update **Project Structure** section at `GEMINI.md` file whenever the project structure changes.

## Project-Specific Rules

- Always ensure that `internal/env/env.go` is updated when adding new environment variables.
- Maintain consistency in Cobra command flags across `cmd/client` and `cmd/server`.
- Use the `cloudflare` package for all DNS-related operations.

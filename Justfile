set dotenv-load

# Default recipe
default:
    @just --list

# Build the server binary
build-server:
    go build -o bin/server ./cmd/server

# Build the client binary
build-client:
    go build -o bin/client ./cmd/client

# Build both server and client binaries
build: build-server build-client

# Run the server (Requires .env file with GO_DDNS_* variables)
run-server: build-server
    ./bin/server

# Run the client (Requires .env file with GO_DDNS_* variables)
run-client: build-client
    ./bin/client --keep-alive

# Run both server and client simultaneously
run: build
    #!/usr/bin/env bash
    set -e
    echo "Starting GO-DDNS Server..."
    ./bin/server &
    SERVER_PID=$!
    
    echo "Starting GO-DDNS Client..."
    ./bin/client --keep-alive &
    CLIENT_PID=$!
    
    trap "kill $SERVER_PID $CLIENT_PID" EXIT INT TERM
    wait

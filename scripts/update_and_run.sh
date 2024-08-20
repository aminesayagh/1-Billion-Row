#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Print each command before executing it
set -x

# Check if Go is installed
if ! [ -x "$(command -v go)" ]; then
  echo 'Error: Go is not installed.' >&2
  exit 1
fi

# Determine OS and ARCH for Go installation
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map machine architecture to Go architecture names
case "$ARCH" in
  x86_64)
    GOARCH='amd64'
    ;;
  aarch64 | armv8*)
    GOARCH='arm64'
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Set the version to install
GOVERSION="1.22.5"

# Download the Go binary
GOTARBALL="go$GOVERSION.$OS-$GOARCH.tar.gz"
wget "https://go.dev/dl/$GOTARBALL"

# Remove any previous Go installation
sudo rm -rf /usr/local/go

# Install the new Go version
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Remove the downloaded tarball
rm "$GOTARBALL"

# Update the PATH environment variable for the current session
export PATH=/usr/local/go/bin:$PATH

# Add the PATH update to the shell profile for future sessions
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.profile

# use ../.env file to server the environment variables for the project
source ../.env

# Verify the installation
go version

# Update the Go modules
go mod tidy

# Install the Go packages
go install

# Run the main.go file
go run main.go
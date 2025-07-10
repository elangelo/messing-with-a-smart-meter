#!/bin/bash

# Smart Meter Reader Setup Script

echo "Setting up Smart Meter Reader..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Go version $GO_VERSION is installed, but version $REQUIRED_VERSION or later is required."
    exit 1
fi

echo "Go version $GO_VERSION detected ✓"

# Initialize Go module if not already done
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init smart-meter-reader
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Create .env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please edit .env file with your specific configuration!"
fi

# Build the application
echo "Building application..."
go build -o smart-meter-reader .

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "Setup complete! Next steps:"
    echo "1. Edit .env file with your configuration"
    echo "2. Connect your smart meter to the P1 port"
    echo "3. Run: ./smart-meter-reader"
    echo ""
    echo "For systemd service setup, see the systemd/ directory"
else
    echo "✗ Build failed. Please check the error messages above."
    exit 1
fi

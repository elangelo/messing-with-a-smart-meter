.PHONY: build run clean test deps setup deploy build-arm logs status restart stop start help help

# Variables - Customize these for your setup
PI_HOST = samuel@arthur
PI_SERVICE = smart-meter-reader
PI_BINARY_PATH = /opt/smart-meter-reader/smart-meter-reader

# Build the application
build:
	go build -o smart-meter-reader .

# Build for Raspberry Pi (ARM)
build-arm:
	GOOS=linux GOARCH=arm GOARM=7 go build -o smart-meter-reader .

# Deploy to Raspberry Pi (stop service, upload binary, restart service)
deploy: build-arm
	@echo "Stopping service on $(PI_HOST)..."
	ssh $(PI_HOST) "sudo systemctl stop $(PI_SERVICE)"
	@echo "Copying binary to $(PI_HOST)..."
	scp smart-meter-reader $(PI_HOST):/tmp/smart-meter-reader-new
	@echo "Installing new binary..."
	ssh $(PI_HOST) "sudo mv /tmp/smart-meter-reader-new $(PI_BINARY_PATH) && sudo chown root:root $(PI_BINARY_PATH) && sudo chmod +x $(PI_BINARY_PATH)"
	@echo "Starting service on $(PI_HOST)..."
	ssh $(PI_HOST) "sudo systemctl start $(PI_SERVICE)"
	@echo "Checking service status..."
	ssh $(PI_HOST) "sudo systemctl status $(PI_SERVICE) --no-pager -l"
	@echo "Deployment complete!"

# Quick deploy (same as deploy, for convenience)
deploy-quick: deploy
	@echo "Quick deployment completed"

# Check logs on Raspberry Pi
logs:
	ssh $(PI_HOST) "sudo journalctl -u $(PI_SERVICE) -f"

# Check status on Raspberry Pi
status:
	ssh $(PI_HOST) "sudo systemctl status $(PI_SERVICE) --no-pager -l"

# Restart service on Raspberry Pi
restart:
	ssh $(PI_HOST) "sudo systemctl restart $(PI_SERVICE)"
	ssh $(PI_HOST) "sudo systemctl status $(PI_SERVICE) --no-pager -l"

# Stop service on Raspberry Pi
stop:
	ssh $(PI_HOST) "sudo systemctl stop $(PI_SERVICE)"

# Start service on Raspberry Pi
start:
	ssh $(PI_HOST) "sudo systemctl start $(PI_SERVICE)"

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	rm -f smart-meter-reader

# Run tests
test:
	go test ./...

# Download dependencies
deps:
	go mod tidy
	go mod download

# Setup everything (deps + build)
setup: deps build

# Install as systemd service (requires sudo)
install-service:
	sudo cp systemd/smart-meter-reader.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable smart-meter-reader
	@echo "Service installed. Edit /etc/systemd/system/smart-meter-reader.service if needed"
	@echo "Start with: sudo systemctl start smart-meter-reader"

# Development helpers
dev:
	go run . -dev

# Check for updates
update:
	go get -u ./...
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Help target
help:
	@echo "Smart Meter Reader - Available Make Targets:"
	@echo ""
	@echo "Building:"
	@echo "  build       - Build for local architecture"
	@echo "  build-arm   - Build for Raspberry Pi (ARM)"
	@echo "  clean       - Remove build artifacts"
	@echo ""
	@echo "Development:"
	@echo "  run         - Run locally"
	@echo "  dev         - Run in development mode"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps        - Download dependencies"
	@echo "  update      - Update dependencies"
	@echo "  setup       - Setup everything (deps + build)"
	@echo ""
	@echo "Deployment (to your Pi):"
	@echo "  deploy      - Build + deploy to Pi + restart service"
	@echo "  deploy-quick- Same as deploy (alias)"
	@echo ""
	@echo "Remote Management:"
	@echo "  status      - Check service status on Pi"
	@echo "  logs        - Follow logs on Pi (Ctrl+C to exit)"
	@echo "  restart     - Restart service on Pi"
	@echo "  start       - Start service on Pi"
	@echo "  stop        - Stop service on Pi"
	@echo ""
	@echo "Local Installation:"
	@echo "  install-service - Install systemd service locally"

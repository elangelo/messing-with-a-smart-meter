package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize configuration
	config := NewConfig()
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize P1 reader
	p1Reader, err := NewP1Reader(config.SerialPort, config.SerialBaudRate)
	if err != nil {
		log.Fatalf("Failed to initialize P1 reader: %v", err)
	}
	defer p1Reader.Close()

	// Initialize InfluxDB client
	influxClient, err := NewInfluxClient(config)
	if err != nil {
		log.Fatalf("Failed to initialize InfluxDB client: %v", err)
	}
	defer influxClient.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	// Start reading and processing data
	log.Println("Starting smart meter reader...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down...")
			return
		default:
			// Read P1 telegram
			telegram, err := p1Reader.ReadTelegram()
			if err != nil {
				log.Printf("Error reading telegram: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Parse the telegram
			data, err := ParseP1Telegram(telegram)
			if err != nil {
				log.Printf("Error parsing telegram: %v", err)
				continue
			}

			// Send to InfluxDB
			if err := influxClient.WriteData(ctx, data); err != nil {
				log.Printf("Error writing to InfluxDB: %v", err)
			} else {
				log.Printf("Successfully sent data to InfluxDB: Power=%dW",
					data.CurrentPowerUsage)
			}

			// Wait a bit before next reading
			time.Sleep(10 * time.Second)
		}
	}
}

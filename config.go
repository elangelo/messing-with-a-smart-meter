package main

import (
	"errors"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	SerialPort       string
	SerialBaudRate   int
	InfluxDBURL      string
	InfluxDBDatabase string
	InfluxDBUsername string
	InfluxDBPassword string
	LogLevel         string
}

// NewConfig creates a new configuration from environment variables
func NewConfig() *Config {
	baudRate := 115200
	if br := os.Getenv("SERIAL_BAUD_RATE"); br != "" {
		if parsed, err := strconv.Atoi(br); err == nil {
			baudRate = parsed
		}
	}

	return &Config{
		SerialPort:       getEnv("SERIAL_PORT", "/dev/ttyUSB0"),
		SerialBaudRate:   baudRate,
		InfluxDBURL:      getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBDatabase: getEnv("INFLUXDB_DATABASE", "smart_meter"),
		InfluxDBUsername: os.Getenv("INFLUXDB_USERNAME"),
		InfluxDBPassword: os.Getenv("INFLUXDB_PASSWORD"),
		LogLevel:         getEnv("LOG_LEVEL", "INFO"),
	}
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	// For InfluxDB v1.x, database name is required but auth is optional
	if c.InfluxDBDatabase == "" {
		return errors.New("INFLUXDB_DATABASE is required")
	}
	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

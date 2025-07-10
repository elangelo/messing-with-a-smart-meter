package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// InfluxClient handles writing data to InfluxDB v1.x via HTTP API
type InfluxClient struct {
	baseURL  string
	database string
	username string
	password string
	client   *http.Client
}

// NewInfluxClient creates a new InfluxDB v1.x client
func NewInfluxClient(config *Config) (*InfluxClient, error) {
	client := &InfluxClient{
		baseURL:  config.InfluxDBURL,
		database: config.InfluxDBDatabase,
		username: config.InfluxDBUsername,
		password: config.InfluxDBPassword,
		client:   &http.Client{Timeout: 10 * time.Second},
	}

	// Test connection by pinging
	if err := client.ping(); err != nil {
		return nil, fmt.Errorf("failed to ping InfluxDB: %v", err)
	}

	// Create database if it doesn't exist
	if err := client.createDatabase(); err != nil {
		fmt.Printf("Warning: Could not create database (it might already exist): %v\n", err)
	}

	return client, nil
}

// ping tests the connection to InfluxDB
func (ic *InfluxClient) ping() error {
	resp, err := ic.client.Get(ic.baseURL + "/ping")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ping failed with status: %d", resp.StatusCode)
	}
	return nil
}

// createDatabase creates the database if it doesn't exist
func (ic *InfluxClient) createDatabase() error {
	query := fmt.Sprintf("CREATE DATABASE %s", ic.database)
	return ic.executeQuery(query)
}

// executeQuery executes a query against InfluxDB
func (ic *InfluxClient) executeQuery(query string) error {
	values := url.Values{}
	values.Set("q", query)
	if ic.username != "" {
		values.Set("u", ic.username)
		values.Set("p", ic.password)
	}

	resp, err := ic.client.PostForm(ic.baseURL+"/query", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("query failed with status: %d", resp.StatusCode)
	}
	return nil
}

// WriteData writes smart meter data to InfluxDB v1.x using Line Protocol
func (ic *InfluxClient) WriteData(ctx context.Context, data *SmartMeterData) error {
	// Build line protocol data
	timestamp := data.Timestamp.UnixNano()

	// Electricity measurement
	electricityLine := fmt.Sprintf("electricity,meter=smart_meter "+
		"consumed_tariff1=%f,consumed_tariff2=%f,produced_tariff1=%f,produced_tariff2=%f,"+
		"current_power_usage=%d,current_power_production=%d,total_consumed=%f,total_produced=%f %d",
		data.ElectricityConsumed1, data.ElectricityConsumed2,
		data.ElectricityProduced1, data.ElectricityProduced2,
		data.CurrentPowerUsage, data.CurrentPowerProduction,
		data.ElectricityConsumed1+data.ElectricityConsumed2,
		data.ElectricityProduced1+data.ElectricityProduced2,
		timestamp)

	// Write to InfluxDB
	return ic.writeLineProtocol(electricityLine)
}

// writeLineProtocol writes line protocol data to InfluxDB
func (ic *InfluxClient) writeLineProtocol(data string) error {
	writeURL := ic.baseURL + "/write"

	// Add query parameters
	values := url.Values{}
	values.Set("db", ic.database)
	values.Set("precision", "ns")
	if ic.username != "" {
		values.Set("u", ic.username)
		values.Set("p", ic.password)
	}

	writeURL += "?" + values.Encode()

	// Make the request
	resp, err := ic.client.Post(writeURL, "application/octet-stream", bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("write failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Close is a no-op for HTTP client
func (ic *InfluxClient) Close() {
	// HTTP client doesn't need explicit closing
}

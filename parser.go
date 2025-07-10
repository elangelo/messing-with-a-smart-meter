package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SmartMeterData represents the parsed data from a P1 telegram
type SmartMeterData struct {
	Timestamp              time.Time
	ElectricityConsumed1   float64 // kWh (tariff 1)
	ElectricityConsumed2   float64 // kWh (tariff 2)
	ElectricityProduced1   float64 // kWh (tariff 1)
	ElectricityProduced2   float64 // kWh (tariff 2)
	CurrentPowerUsage      int     // W
	CurrentPowerProduction int     // W
}

// ParseP1Telegram parses a DSMR P1 telegram and extracts meter data
func ParseP1Telegram(telegram string) (*SmartMeterData, error) {
	data := &SmartMeterData{
		Timestamp: time.Now(),
	}

	lines := strings.Split(telegram, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse electricity consumption (1.8.1 = consumed tariff 1)
		if matches := extractValue(line, `1-0:1\.8\.1\((\d+\.\d+)\*kWh\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.ElectricityConsumed1 = val
			}
		}

		// Parse electricity consumption (1.8.2 = consumed tariff 2)
		if matches := extractValue(line, `1-0:1\.8\.2\((\d+\.\d+)\*kWh\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.ElectricityConsumed2 = val
			}
		}

		// Parse electricity production (2.8.1 = produced tariff 1)
		if matches := extractValue(line, `1-0:2\.8\.1\((\d+\.\d+)\*kWh\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.ElectricityProduced1 = val
			}
		}

		// Parse electricity production (2.8.2 = produced tariff 2)
		if matches := extractValue(line, `1-0:2\.8\.2\((\d+\.\d+)\*kWh\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.ElectricityProduced2 = val
			}
		}

		// Parse current power usage (1.7.0)
		if matches := extractValue(line, `1-0:1\.7\.0\((\d+\.\d+)\*kW\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.CurrentPowerUsage = int(val * 1000) // Convert kW to W
			}
		}

		// Parse current power production (2.7.0)
		if matches := extractValue(line, `1-0:2\.7\.0\((\d+\.\d+)\*kW\)`); matches != nil {
			if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
				data.CurrentPowerProduction = int(val * 1000) // Convert kW to W
			}
		}
	}

	return data, nil
}

// extractValue extracts a value using regex pattern
func extractValue(line, pattern string) []string {
	re := regexp.MustCompile(pattern)
	return re.FindStringSubmatch(line)
}

package main

import (
	"testing"
	"time"
)

func TestParseP1Telegram(t *testing.T) {
	// Example P1 telegram (simplified)
	telegram := `/XMX5LGBBFG1012327662

1-3:0.2.8(42)
0-0:1.0.0(230315123456W)
1-0:1.8.1(000123.456*kWh)
1-0:1.8.2(000234.567*kWh)
1-0:2.8.1(000012.345*kWh)
1-0:2.8.2(000023.456*kWh)
1-0:1.7.0(00.324*kW)
1-0:2.7.0(00.000*kW)
0-1:24.2.1(230315120000W)(00012.345*m3)
!1234`

	data, err := ParseP1Telegram(telegram)
	if err != nil {
		t.Fatalf("Failed to parse telegram: %v", err)
	}

	// Check parsed values
	if data.ElectricityConsumed1 != 123.456 {
		t.Errorf("Expected ElectricityConsumed1 to be 123.456, got %f", data.ElectricityConsumed1)
	}

	if data.ElectricityConsumed2 != 234.567 {
		t.Errorf("Expected ElectricityConsumed2 to be 234.567, got %f", data.ElectricityConsumed2)
	}

	if data.CurrentPowerUsage != 324 {
		t.Errorf("Expected CurrentPowerUsage to be 324W, got %d", data.CurrentPowerUsage)
	}

	if data.GasConsumption != 12.345 {
		t.Errorf("Expected GasConsumption to be 12.345, got %f", data.GasConsumption)
	}
}

func TestParseP1Timestamp(t *testing.T) {
	timestamp := "230315123456W"
	parsed, err := parseP1Timestamp(timestamp)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	expected := time.Date(2023, 3, 15, 12, 34, 56, 0, time.Local)
	if !parsed.Equal(expected) {
		t.Errorf("Expected timestamp %v, got %v", expected, parsed)
	}
}

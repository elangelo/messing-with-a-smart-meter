package main

import (
	"testing"
)

func TestParseP1Telegram(t *testing.T) {
	// Example P1 telegram (simplified, electricity only)
	telegram := `/XMX5LGBBFG1012327662

1-3:0.2.8(42)
0-0:1.0.0(230315123456W)
1-0:1.8.1(000123.456*kWh)
1-0:1.8.2(000234.567*kWh)
1-0:2.8.1(000012.345*kWh)
1-0:2.8.2(000023.456*kWh)
1-0:1.7.0(00.324*kW)
1-0:2.7.0(01.234*kW)
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

	if data.ElectricityProduced1 != 12.345 {
		t.Errorf("Expected ElectricityProduced1 to be 12.345, got %f", data.ElectricityProduced1)
	}

	if data.ElectricityProduced2 != 23.456 {
		t.Errorf("Expected ElectricityProduced2 to be 23.456, got %f", data.ElectricityProduced2)
	}

	if data.CurrentPowerUsage != 324 {
		t.Errorf("Expected CurrentPowerUsage to be 324W, got %d", data.CurrentPowerUsage)
	}

	if data.CurrentPowerProduction != 1234 {
		t.Errorf("Expected CurrentPowerProduction to be 1234W, got %d", data.CurrentPowerProduction)
	}
}

func TestParseP1TelegramEdgeCases(t *testing.T) {
	t.Run("Empty telegram", func(t *testing.T) {
		data, err := ParseP1Telegram("")
		if err != nil {
			t.Fatalf("Expected no error for empty telegram, got: %v", err)
		}

		// Should return zero values for all fields
		if data.ElectricityConsumed1 != 0 {
			t.Errorf("Expected ElectricityConsumed1 to be 0, got %f", data.ElectricityConsumed1)
		}
	})

	t.Run("Minimal valid telegram", func(t *testing.T) {
		telegram := `/XMX5LGBBFG1012327662
1-0:1.7.0(00.500*kW)
!1234`

		data, err := ParseP1Telegram(telegram)
		if err != nil {
			t.Fatalf("Failed to parse minimal telegram: %v", err)
		}

		if data.CurrentPowerUsage != 500 {
			t.Errorf("Expected CurrentPowerUsage to be 500W, got %d", data.CurrentPowerUsage)
		}
	})

	t.Run("Production only", func(t *testing.T) {
		telegram := `/XMX5LGBBFG1012327662
1-0:2.7.0(02.500*kW)
!1234`

		data, err := ParseP1Telegram(telegram)
		if err != nil {
			t.Fatalf("Failed to parse production telegram: %v", err)
		}

		if data.CurrentPowerProduction != 2500 {
			t.Errorf("Expected CurrentPowerProduction to be 2500W, got %d", data.CurrentPowerProduction)
		}

		if data.CurrentPowerUsage != 0 {
			t.Errorf("Expected CurrentPowerUsage to be 0W when not specified, got %d", data.CurrentPowerUsage)
		}
	})
}

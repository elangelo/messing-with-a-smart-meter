package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// P1Reader handles reading data from the P1 port
type P1Reader struct {
	port io.ReadWriteCloser
}

// NewP1Reader creates a new P1 reader
func NewP1Reader(portName string, baudRate int) (*P1Reader, error) {
	config := &serial.Config{
		Name:        portName,
		Baud:        baudRate,
		ReadTimeout: time.Second * 5,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port %s: %v", portName, err)
	}

	return &P1Reader{
		port: port,
	}, nil
}

// ReadTelegram reads a complete P1 telegram from the serial port
func (p *P1Reader) ReadTelegram() (string, error) {
	scanner := bufio.NewScanner(p.port)
	var telegram strings.Builder
	inTelegram := false

	for scanner.Scan() {
		line := scanner.Text()

		// P1 telegram starts with "/"
		if strings.HasPrefix(line, "/") {
			inTelegram = true
			telegram.Reset()
		}

		if inTelegram {
			telegram.WriteString(line)
			telegram.WriteString("\n")

			// P1 telegram ends with "!" followed by CRC
			if strings.HasPrefix(line, "!") && len(line) == 5 {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from serial port: %v", err)
	}

	telegramStr := telegram.String()
	if telegramStr == "" {
		return "", fmt.Errorf("no complete telegram received")
	}

	return telegramStr, nil
}

// Close closes the serial port connection
func (p *P1Reader) Close() error {
	if p.port != nil {
		return p.port.Close()
	}
	return nil
}

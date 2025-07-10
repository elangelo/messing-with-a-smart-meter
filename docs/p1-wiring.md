# P1 Port Wiring Guide

## What you need
- Raspberry Pi (any model with USB ports)
- USB to Serial adapter (FTDI chip recommended)
- RJ12 cable or individual wires
- Breadboard/connector for wiring

## P1 Port Pinout
The P1 port on Dutch smart meters typically uses RJ12 (6-pin) connector:

```
Pin 1: +5V (power supply from meter, optional)
Pin 2: RTS (Request to Send) - connect to DTR of USB-Serial
Pin 3: GND (Ground) - connect to GND of USB-Serial
Pin 4: NC (Not Connected)
Pin 5: RXD (Receive Data) - connect to TX of USB-Serial
Pin 6: GND (Ground) - can also connect to GND
```

## Wiring Diagram
```
Smart Meter P1 Port    USB-Serial Adapter
─────────────────      ──────────────────
Pin 1 (+5V)       ──── (not connected)
Pin 2 (RTS)       ──── DTR
Pin 3 (GND)       ──── GND
Pin 4 (NC)        ──── (not connected)
Pin 5 (RXD)       ──── TX
Pin 6 (GND)       ──── GND (optional)
```

## USB-Serial Adapter Setup
1. Connect the USB-Serial adapter to your Raspberry Pi
2. Check if it's detected: `lsusb`
3. Find the device: `ls /dev/ttyUSB*` or `ls /dev/ttyACM*`
4. Add your user to the dialout group: `sudo usermod -a -G dialout $USER`
5. Logout and login again for group changes to take effect

## Testing Connection
You can test the connection manually:
```bash
# Install screen or minicom
sudo apt install screen

# Read data from the port
screen /dev/ttyUSB0 115200

# You should see data flowing every 10 seconds
# Press Ctrl+A then K to exit screen
```

## Common Issues
- **Permission denied**: Make sure user is in dialout group
- **No data**: Check wiring, especially RTS connection
- **Garbled data**: Wrong baud rate (should be 115200 for most meters)
- **Device not found**: USB-Serial adapter not detected by system

## Supported Meters
This application supports DSMR (Dutch Smart Meter Requirements) versions:
- DSMR 2.2
- DSMR 4.0
- DSMR 5.0

Most Dutch smart meters from Landis+Gyr, Kaifa, and Sagemcom are supported.

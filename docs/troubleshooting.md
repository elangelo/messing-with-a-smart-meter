# Troubleshooting Guide

## Common Issues and Solutions

### 1. "Permission denied" when accessing serial port

**Problem**: Application can't access `/dev/ttyUSB0`

**Solution**:
```bash
# Add your user to the dialout group
sudo usermod -a -G dialout $USER

# Logout and login again, or use:
newgrp dialout

# Check permissions
ls -la /dev/ttyUSB0
```

### 2. "No such file or directory" for serial port

**Problem**: Serial device not found

**Solutions**:
```bash
# Check what USB devices are connected
lsusb

# List all serial ports
ls /dev/tty*

# Check if USB-Serial adapter is detected
dmesg | grep -i usb
dmesg | grep -i tty
```

Common device paths:
- `/dev/ttyUSB0` - USB-to-Serial adapters
- `/dev/ttyACM0` - Arduino-compatible devices
- `/dev/ttyAMA0` - Raspberry Pi built-in UART

### 3. "Failed to connect to InfluxDB"

**Problem**: Can't reach InfluxDB server

**Solutions**:
```bash
# Check if InfluxDB is running
sudo systemctl status influxdb
# or for Docker:
docker ps | grep influxdb

# Test connection manually
curl http://localhost:8086/health

# Check if port is open
netstat -tlnp | grep 8086
```

### 4. "No complete telegram received"

**Problem**: Not receiving data from smart meter

**Solutions**:
1. **Check wiring**: Especially RTS (pin 2) connection
2. **Test manually**:
   ```bash
   # Install screen if not available
   sudo apt install screen
   
   # Monitor port directly
   screen /dev/ttyUSB0 115200
   # You should see data every 10 seconds
   # Press Ctrl+A then K to exit
   ```
3. **Check if meter is transmitting**: Some meters need to be enabled for P1 output
4. **Verify baud rate**: Most use 115200, but some older meters use 9600

### 5. "InfluxDB write error"

**Problem**: Data not being written to InfluxDB

**Solutions**:
1. **Verify token permissions**: Token needs write access to the bucket
2. **Check bucket exists**: Create bucket if missing
3. **Verify organization name**: Must match exactly

### 6. Application starts but no data

**Problem**: Application runs but doesn't process data

**Debug steps**:
```bash
# Run with verbose logging
LOG_LEVEL=DEBUG ./smart-meter-reader

# Check what's on the serial port
sudo cat /dev/ttyUSB0

# Test serial port settings
stty -F /dev/ttyUSB0 115200 raw -echo
```

### 7. Building/compilation errors

**Problem**: Go build fails

**Solutions**:
```bash
# Update Go to latest version (1.21+)
go version

# Clean module cache
go clean -modcache

# Retry dependency download
go mod download
go mod tidy

# Build with verbose output
go build -v .
```

## Testing Your Setup

### 1. Test serial connection
```bash
# Method 1: Using screen
screen /dev/ttyUSB0 115200

# Method 2: Using cat (might be messy)
timeout 30 cat /dev/ttyUSB0

# Method 3: Using minicom
sudo apt install minicom
minicom -D /dev/ttyUSB0 -b 115200
```

### 2. Test InfluxDB connection
```bash
# Check health
curl http://localhost:8086/health

# Test with your token
curl -H "Authorization: Token YOUR_TOKEN" \
     "http://localhost:8086/api/v2/buckets?org=YOUR_ORG"
```

### 3. Test parsing
```bash
# Run tests
go test .

# Run with test data
echo 'your-p1-telegram-here' | go run . -test
```

## Getting Help

1. **Check logs**: Application logs to stdout/stderr
2. **Run in foreground**: Don't use systemd while debugging
3. **Test components separately**: Serial port, InfluxDB, parsing
4. **Check smart meter manual**: Some meters have specific requirements

## Log Analysis

Common log messages and what they mean:

- `"Starting smart meter reader..."` - Application started successfully
- `"Successfully sent data to InfluxDB"` - Data flow working correctly
- `"Error reading telegram"` - Serial port or connection issue
- `"Error parsing telegram"` - Data format issue or corrupt transmission
- `"Error writing to InfluxDB"` - Database connection or permission issue

Enable debug logging for more details:
```bash
LOG_LEVEL=DEBUG ./smart-meter-reader
```

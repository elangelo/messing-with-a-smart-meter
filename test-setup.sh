#!/bin/bash

# Test script for smart meter reader

echo "🔧 Testing Smart Meter Reader Setup..."
echo ""

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please copy .env.example to .env and configure it."
    exit 1
fi

echo "✅ .env file found"

# Source the .env file
source .env

# Test InfluxDB connection
echo "🌐 Testing InfluxDB connection..."
if curl -s "${INFLUXDB_URL}/ping" > /dev/null; then
    echo "✅ InfluxDB is reachable at ${INFLUXDB_URL}"
else
    echo "❌ Cannot reach InfluxDB at ${INFLUXDB_URL}"
    exit 1
fi

# Check if database exists
echo "🗄️  Checking database..."
DB_CHECK=$(curl -s "${INFLUXDB_URL}/query" --data-urlencode "q=SHOW DATABASES" | grep -o "\"${INFLUXDB_DATABASE}\"")
if [ -n "$DB_CHECK" ]; then
    echo "✅ Database '${INFLUXDB_DATABASE}' exists"
else
    echo "⚠️  Database '${INFLUXDB_DATABASE}' not found, will be created automatically"
fi

# Check if serial port exists (if specified and not empty)
if [ -n "$SERIAL_PORT" ] && [ "$SERIAL_PORT" != "" ]; then
    echo "📡 Checking serial port..."
    if [ -e "$SERIAL_PORT" ]; then
        echo "✅ Serial port ${SERIAL_PORT} exists"
        # Check permissions
        if [ -r "$SERIAL_PORT" ] && [ -w "$SERIAL_PORT" ]; then
            echo "✅ Serial port permissions OK"
        else
            echo "⚠️  Serial port permission issue. Make sure user is in dialout group:"
            echo "   sudo usermod -a -G dialout \$USER"
            echo "   Then logout and login again"
        fi
    else
        echo "⚠️  Serial port ${SERIAL_PORT} not found"
        echo "   Common ports: /dev/ttyUSB0, /dev/ttyUSB1, /dev/ttyACM0"
        echo "   Check with: ls /dev/ttyUSB* /dev/ttyACM*"
    fi
else
    echo "⚠️  SERIAL_PORT not configured in .env"
fi

# Test building the application
echo "🔨 Testing build..."
if go build -o smart-meter-reader-test . > /dev/null 2>&1; then
    echo "✅ Application builds successfully"
    rm -f smart-meter-reader-test
else
    echo "❌ Build failed"
    echo "Run 'go build .' for details"
    exit 1
fi

echo ""
echo "🎉 Setup test completed!"
echo ""
echo "Next steps:"
echo "1. Connect your smart meter P1 port (see docs/p1-wiring.md)"
echo "2. Update SERIAL_PORT in .env if needed"
echo "3. Run: ./smart-meter-reader"
echo ""
echo "To test with fake data (no serial port needed):"
echo "Run the application - it will show connection errors but InfluxDB integration will work"

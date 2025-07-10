# Smart Meter Reader

A Go application to read smart meter data from the P1 port and send it to InfluxDB for monitoring with Grafana.

## Features

- 📊 Reads DSMR (Dutch Smart Meter Requirements) data from P1 port
- ⚡ Parses electricity consumption and production data (solar panels)
- 🗄️ Sends data to InfluxDB v1.x
- 🔧 Configurable via environment variables  
- 🛡️ Robust error handling and reconnection logic
- 🚀 Easy deployment to Raspberry Pi via Makefile
- 📈 Includes Grafana dashboard configuration

## Quick Start

1. **Clone and setup:**
   ```bash
   git clone <your-repo>
   cd smart-meter-reader
   make setup
   ```

2. **Configure:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Run locally:**
   ```bash
   make run
   ```

4. **Deploy to Raspberry Pi:**
   ```bash
   make deploy
   ```

## Configuration

Create a `.env` file based on `.env.example`:

```bash
# Serial port configuration
SERIAL_PORT=/dev/ttyUSB0
SERIAL_BAUD_RATE=115200

# InfluxDB v1.x configuration  
INFLUXDB_URL=http://your-influxdb-host:8086
INFLUXDB_DATABASE=smart_meter
INFLUXDB_USERNAME=
INFLUXDB_PASSWORD=
# Leave username/password empty if no auth configured

# Optional settings
LOG_LEVEL=INFO
```

## Available Make Commands

```bash
make help          # Show all available commands
make build         # Build for local architecture
make build-arm     # Build for Raspberry Pi (ARM)
make deploy        # Build + deploy to Pi + restart service
make status        # Check service status on Pi
make logs          # Follow logs on Pi (Ctrl+C to exit)
make restart       # Restart service on Pi
```

## P1 Port Connection

Make sure your Raspberry Pi is connected to the smart meter's P1 port. The P1 port typically uses a RJ12 connector and provides data at 115200 baud rate.

Common P1 port pinout:
- Pin 1: +5V (not always used)
- Pin 2: RTS (Request to Send)
- Pin 3: GND (Ground)
- Pin 4: NC (Not Connected)
- Pin 5: RXD (Receive Data)
- Pin 6: GND (Ground)

## Usage

### Local Development
```bash
make run
```

### Production Deployment
```bash
# Edit Makefile variables for your Pi:
# PI_HOST = your-user@your-pi-hostname

make deploy        # Automated deployment
```

### Manual Installation on Raspberry Pi
```bash
make build-arm
scp smart-meter-reader pi@your-pi:/tmp/
# Then install as systemd service (see docs/)
```

## Documentation

- 📖 [InfluxDB Setup](docs/influxdb-setup.md)
- 📊 [Grafana Dashboard Setup](docs/grafana-setup.md)
- ⚙️ [Systemd Service Setup](systemd/)

## Data Collected

The application collects the following electricity data:
- ⚡ Current power usage (W)
- 🌞 Current power production (W) - solar panels
- 📊 Total consumed energy (kWh) - tariff 1 & 2
- 📈 Total produced energy (kWh) - tariff 1 & 2

Data is stored every 10 seconds and can be visualized in Grafana.

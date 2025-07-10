# InfluxDB v1.x Setup Guide

## Your Current Setup
Your InfluxDB v1.8.10 is already running and accessible at `http://192.168.0.53:8086`. 

## Database Creation
The smart meter reader will automatically create the `smart_meter` database if it doesn't exist. You can also create it manually:

```bash
curl -i http://192.168.0.53:8086/query --data-urlencode "q=CREATE DATABASE smart_meter"
```

## Configuration
Since you're not using authentication, your `.env` file should be:

```env
# InfluxDB v1.x configuration
INFLUXDB_URL=http://192.168.0.53:8086
INFLUXDB_DATABASE=smart_meter
INFLUXDB_USERNAME=
INFLUXDB_PASSWORD=
# Leave username/password empty if no auth is configured
```

## Data Schema

The smart meter reader writes data to two measurements:

### Electricity Measurement
- **Measurement**: `electricity`
- **Tags**: `meter=smart_meter`
- **Fields**:
  - `consumed_tariff1` (kWh)
  - `consumed_tariff2` (kWh)
  - `produced_tariff1` (kWh)
  - `produced_tariff2` (kWh)
  - `current_power_usage` (W)
  - `current_power_production` (W)
  - `total_consumed` (kWh)
  - `total_produced` (kWh)

### Gas Measurement
- **Measurement**: `gas`
- **Tags**: `meter=smart_meter`
- **Fields**:
  - `consumption` (m³)

## Sample Queries (InfluxQL)

### Current Power Usage
```sql
SELECT last("current_power_usage") FROM "electricity" WHERE time > now() - 1h
```

### Daily Energy Consumption
```sql
SELECT max("total_consumed") - min("total_consumed") as "daily_consumption" 
FROM "electricity" 
WHERE time > now() - 1d 
GROUP BY time(1h)
```

### Gas Usage Over Time
```sql
SELECT "consumption" FROM "gas" WHERE time > now() - 7d
```

### Show all measurements
```bash
curl -G http://192.168.0.53:8086/query?db=smart_meter --data-urlencode "q=SHOW MEASUREMENTS"
```

### Show recent electricity data
```bash
curl -G http://192.168.0.53:8086/query?db=smart_meter --data-urlencode "q=SELECT * FROM electricity ORDER BY time DESC LIMIT 10"
```

## Grafana Integration

1. Install Grafana
2. Add InfluxDB as data source:
   - URL: `http://192.168.0.53:8086`
   - Database: `smart_meter`
   - No authentication needed (unless you enable it)

3. Create dashboards using InfluxQL queries

## Backup and Maintenance

### Backup
```bash
# Backup InfluxDB data
influx backup /path/to/backup --org your-organization
```

### Retention Policy
Consider setting up a retention policy to manage disk space:
```bash
influx bucket update \
  --name smart-meter \
  --retention 365d
```

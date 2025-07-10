# Smart Meter Grafana Dashboard Setup

## Step 1: Add InfluxDB Data Source

1. **Open Grafana** in your browser
2. **Go to Configuration → Data Sources** (gear icon in sidebar)
3. **Click "Add data source"**
4. **Select "InfluxDB"**
5. **Configure the data source:**
   - **Name**: `Smart Meter InfluxDB`
   - **URL**: `http://192.168.0.53:8086`
   - **Database**: `smart_meter`
   - **User**: (leave empty - no auth)
   - **Password**: (leave empty - no auth)
   - **HTTP Method**: `GET`
6. **Click "Save & Test"** - should show green "Data source is working"

## Step 2: Create the Dashboard

### Option A: Import Pre-built Dashboard (Recommended)
1. **Go to Dashboards → Import** (+ icon in sidebar → Import)
2. **Copy the JSON** from `smart-meter-dashboard.json` (created below)
3. **Paste it** into the import field
4. **Click "Load"**
5. **Select your InfluxDB data source**
6. **Click "Import"**

### Option B: Create Manually
Follow the panel configurations below to create each visualization.

## Panel Configurations

### 1. Current Power Usage (Gauge)
- **Panel Type**: Gauge
- **Query**: 
  ```sql
  SELECT last("current_power_usage") FROM "electricity" WHERE $timeFilter
  ```
- **Display**: 
  - Min: 0, Max: 5000
  - Unit: Watts (W)
  - Thresholds: Green (0-1000), Yellow (1000-3000), Red (3000+)

### 2. Current Power Production (Gauge)
- **Panel Type**: Gauge  
- **Query**:
  ```sql
  SELECT last("current_power_production") FROM "electricity" WHERE $timeFilter
  ```
- **Display**:
  - Min: 0, Max: 3000
  - Unit: Watts (W)
  - Thresholds: Green (1000+), Yellow (500-1000), Red (0-500)

### 3. Net Power Flow (Stat)
- **Panel Type**: Stat
- **Query A** (Usage):
  ```sql
  SELECT last("current_power_usage") FROM "electricity" WHERE $timeFilter
  ```
- **Query B** (Production):
  ```sql
  SELECT last("current_power_production") FROM "electricity" WHERE $timeFilter
  ```
- **Transform**: Add field from calculation: `Production - Usage`
- **Display**: 
  - Unit: Watts (W)
  - Positive = Exporting to grid
  - Negative = Importing from grid

### 4. Power Over Time (Time Series)
- **Panel Type**: Time series
- **Query A** (Usage):
  ```sql
  SELECT mean("current_power_usage") FROM "electricity" WHERE $timeFilter GROUP BY time($__interval) fill(null)
  ```
- **Query B** (Production):
  ```sql
  SELECT mean("current_power_production") FROM "electricity" WHERE $timeFilter GROUP BY time($__interval) fill(null)
  ```
- **Display**: 
  - Y-axis unit: Watts (W)
  - Different colors for usage vs production

### 5. Daily Energy Consumption (Bar Chart)
- **Panel Type**: Barchart
- **Query**:
  ```sql
  SELECT max("total_consumed") - min("total_consumed") as "consumption" FROM "electricity" WHERE $timeFilter GROUP BY time(1d) fill(null)
  ```
- **Display**: 
  - Unit: kWh
  - X-axis: Time (days)

### 6. Daily Energy Production (Bar Chart)
- **Panel Type**: Barchart
- **Query**:
  ```sql
  SELECT max("total_produced") - min("total_produced") as "production" FROM "electricity" WHERE $timeFilter GROUP BY time(1d) fill(null)
  ```
- **Display**: 
  - Unit: kWh
  - X-axis: Time (days)

### 7. Cumulative Consumption vs Production (Time Series)
- **Panel Type**: Time series
- **Query A** (Total Consumed):
  ```sql
  SELECT last("total_consumed") FROM "electricity" WHERE $timeFilter GROUP BY time($__interval) fill(null)
  ```
- **Query B** (Total Produced):
  ```sql
  SELECT last("total_produced") FROM "electricity" WHERE $timeFilter GROUP BY time($__interval) fill(null)
  ```
- **Display**: 
  - Unit: kWh
  - Show cumulative totals over time

### 8. Energy Balance Table
- **Panel Type**: Table
- **Query**:
  ```sql
  SELECT 
    last("total_consumed") as "Total Consumed (kWh)",
    last("total_produced") as "Total Produced (kWh)",
    last("total_produced") - last("total_consumed") as "Net Balance (kWh)"
  FROM "electricity" WHERE $timeFilter
  ```

## Dashboard Settings

### Time Range
- **Default**: Last 24 hours
- **Quick ranges**: 1h, 6h, 12h, 24h, 7d, 30d

### Refresh Rate
- **Auto-refresh**: 30s or 1m
- **Live updates** for real-time monitoring

### Variables (Optional)
Create dashboard variables for:
- **Time interval**: `$__interval`
- **Meter filter**: `meter` (if you have multiple meters)

## Tips for Better Dashboards

1. **Use appropriate time ranges** - shorter intervals for real-time monitoring
2. **Set up alerts** for unusual power consumption/production
3. **Add annotations** for events (maintenance, weather, etc.)
4. **Use template variables** for filtering
5. **Organize panels** logically (current status at top, historical data below)

## Sample Queries for Advanced Analysis

### Peak Power Hours
```sql
SELECT max("current_power_usage") FROM "electricity" WHERE $timeFilter GROUP BY time(1h)
```

### Solar Efficiency (Production vs Time of Day)
```sql
SELECT mean("current_power_production") FROM "electricity" WHERE $timeFilter GROUP BY time(1h)
```

### Energy Cost Calculation (if you have tariff data)
```sql
SELECT 
  sum("consumed_tariff1") * 0.25 + sum("consumed_tariff2") * 0.30 as "estimated_cost"
FROM "electricity" WHERE $timeFilter GROUP BY time(1d)
```

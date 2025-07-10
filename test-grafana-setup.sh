#!/bin/bash

# Quick test script for Grafana dashboard setup

echo "🎯 Smart Meter Grafana Dashboard Setup"
echo "======================================="
echo ""

# Test InfluxDB connection
echo "1. Testing InfluxDB connection..."
if curl -s "http://192.168.0.53:8086/ping" > /dev/null; then
    echo "   ✅ InfluxDB is accessible"
else
    echo "   ❌ Cannot reach InfluxDB"
    exit 1
fi

# Check if we have data
echo ""
echo "2. Checking for smart meter data..."
DATA_COUNT=$(curl -s "http://192.168.0.53:8086/query?db=smart_meter" --data-urlencode "q=SELECT COUNT(*) FROM electricity" | grep -o '"values":\[\["[^"]*",[0-9]*\]\]' | grep -o '[0-9]*' | tail -1)

if [ "$DATA_COUNT" -gt 0 ]; then
    echo "   ✅ Found $DATA_COUNT electricity measurements"
else
    echo "   ⚠️  No electricity data found yet"
fi

# Show sample recent data
echo ""
echo "3. Sample recent data:"
echo "   Querying last 5 readings..."
curl -s "http://192.168.0.53:8086/query?db=smart_meter" --data-urlencode "q=SELECT time, current_power_usage, current_power_production, total_consumed, total_produced FROM electricity ORDER BY time DESC LIMIT 5" | jq -r '.results[0].series[0].values[] | @tsv' | while IFS=$'\t' read -r time usage production consumed produced; do
    echo "   📊 $time: Usage=${usage}W, Production=${production}W, Consumed=${consumed}kWh, Produced=${produced}kWh"
done

echo ""
echo "4. Dashboard Setup Instructions:"
echo "   📋 1. Open Grafana in your browser"
echo "   📋 2. Go to Configuration → Data Sources"
echo "   📋 3. Add InfluxDB data source:"
echo "        - URL: http://192.168.0.53:8086"
echo "        - Database: smart_meter"
echo "        - No authentication needed"
echo "   📋 4. Import dashboard:"
echo "        - Go to Dashboards → Import"
echo "        - Copy content from docs/smart-meter-dashboard.json"
echo "        - Paste and import"
echo ""
echo "🎉 Ready for dashboard creation!"

# Show what panels will be available
echo ""
echo "5. Your dashboard will show:"
echo "   🔋 Current power usage (gauge)"
echo "   ☀️  Current solar production (gauge)"  
echo "   ⚖️  Net power flow (import/export)"
echo "   📈 Power trends over time"
echo "   📊 Cumulative energy consumption vs production"
echo "   📋 Energy summary table"
echo ""
echo "📖 See docs/grafana-setup.md for detailed instructions!"

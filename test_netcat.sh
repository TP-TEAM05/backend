#!/bin/bash

echo "╔════════════════════════════════════════════════════════╗"
echo "║        Netcat UDP Packet Monitor - Quick Test         ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""

# Test 1: Basic connectivity
echo "📡 Test 1: Basic UDP Connectivity"
echo "Starting listener on port 5051..."
echo ""

# Note: This is an interactive command - it will block until you press Ctrl+C
echo "Run this command in a terminal:"
echo ""
echo "  docker exec -it reco-backend-1 nc -ul 5051"
echo ""
echo "Then in another terminal, send a test packet:"
echo ""
echo "  echo '{\"type\":\"test\",\"message\":\"hello\"}' | docker exec -i reco-backend-1 nc -u localhost 5051"
echo ""
echo "You should see the JSON appear in the first terminal."
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Show current production monitoring (better alternative)
echo "💡 Better Alternative: Use Existing Logs"
echo ""
echo "Instead of netcat, use the production logging we already set up:"
echo ""
echo "View packets being SENT:"
echo "  docker logs reco-car-integration-1 -f | grep BACKEND-TX"
echo ""
echo "View packets being RECEIVED:"
echo "  docker logs reco-backend-1 -f | grep BACKEND-RX"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Show current status
echo "📊 Current Packet Flow Status:"
echo ""

TX_COUNT=$(docker logs reco-car-integration-1 2>&1 | grep -c BACKEND-TX)
RX_COUNT=$(docker logs reco-backend-1 2>&1 | grep -c BACKEND-RX)

echo "Packets SENT (TX):     $TX_COUNT"
echo "Packets RECEIVED (RX): $RX_COUNT"
echo ""

if [ $TX_COUNT -gt 0 ] && [ $RX_COUNT -gt 0 ]; then
    echo "✅ Packet flow is WORKING!"
else
    echo "⚠️  No packets detected in logs"
fi

echo ""
echo "Last 3 TX packets:"
docker logs reco-car-integration-1 2>&1 | grep BACKEND-TX | tail -3
echo ""
echo "Last 3 RX packets:"
docker logs reco-backend-1 2>&1 | grep BACKEND-RX | tail -3

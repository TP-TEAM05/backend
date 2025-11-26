#!/usr/bin/env python3
"""
UDP Packet Sniffer for Backend
Listens on port 5050 and logs all incoming UDP packets from car-integration
"""

import socket
import json
import sys
from datetime import datetime

def main():
    # Configuration
    HOST = '0.0.0.0'  # Listen on all interfaces
    PORT = 5051  # Use different port to not conflict with backend (we'll test this separately)

    if len(sys.argv) > 1:
        PORT = int(sys.argv[1])

    print(f"[UDP-SNIFFER] Starting UDP packet sniffer on {HOST}:{PORT}")
    print(f"[UDP-SNIFFER] Timestamp: {datetime.now().isoformat()}")
    print("-" * 80)

    # Create UDP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((HOST, PORT))

    print(f"[UDP-SNIFFER] Listening for packets...")
    print("-" * 80)

    packet_count = 0

    try:
        while True:
            # Receive data
            data, addr = sock.recvfrom(65536)  # Max UDP packet size
            packet_count += 1

            timestamp = datetime.now().isoformat()

            # Try to parse as JSON
            try:
                json_data = json.loads(data.decode('utf-8'))
                print(f"\n[PACKET #{packet_count}] {timestamp}")
                print(f"From: {addr[0]}:{addr[1]}")
                print(f"Size: {len(data)} bytes")
                print(f"Type: {json_data.get('type', 'unknown')}")
                print(f"Index: {json_data.get('index', 'N/A')}")

                # Print relevant fields based on type
                packet_type = json_data.get('type', '')

                if packet_type == 'update_vehicle_position':
                    vehicle = json_data.get('vehicle', {})
                    print(f"VIN: {vehicle.get('vin', 'N/A')}")
                    print(f"Position: ({vehicle.get('latitude', 'N/A')}, {vehicle.get('longitude', 'N/A')})")
                    print(f"Speed: {vehicle.get('speed', 'N/A')}")
                elif packet_type == 'update_vehicles':
                    vehicles = json_data.get('vehicles', [])
                    print(f"Vehicles count: {len(vehicles)}")
                    for v in vehicles:
                        print(f"  - VIN: {v.get('vin', 'N/A')}")
                elif packet_type == 'acknowledge':
                    print(f"Acknowledging Index: {json_data.get('acknowledgingIndex', 'N/A')}")

                # Optional: print full JSON for debugging
                if '--verbose' in sys.argv:
                    print("Full JSON:")
                    print(json.dumps(json_data, indent=2))

            except json.JSONDecodeError:
                print(f"\n[PACKET #{packet_count}] {timestamp}")
                print(f"From: {addr[0]}:{addr[1]}")
                print(f"Size: {len(data)} bytes")
                print(f"[WARNING] Could not parse as JSON")
                print(f"Raw data (first 200 bytes): {data[:200]}")
            except Exception as e:
                print(f"\n[ERROR] Exception processing packet: {e}")

            print("-" * 80)

    except KeyboardInterrupt:
        print(f"\n\n[UDP-SNIFFER] Stopped. Total packets received: {packet_count}")
        sock.close()

if __name__ == "__main__":
    main()

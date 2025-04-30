#!/usr/bin/env python3
"""
Find a specific packet in the JSON file and display its parsed data.
"""

import json
import sys

def find_packet(json_file, packet_id):
    """Find a packet with the specified ID in the JSON file."""
    with open(json_file, 'r') as f:
        data = json.load(f)
    
    for packet in data['packets']:
        if packet['packet_id'] == packet_id:
            print(f"Found packet {packet_id}: {packet['description']}")
            print(f"Direction: {packet['direction']}")
            print(f"Server type: {packet['server_type']}")
            print(f"Timestamp: {packet['timestamp']}")
            print(f"Size: {packet['size']} bytes")
            print("\nParsed data:")
            if packet['parsed_data']:
                for key, value in packet['parsed_data'].items():
                    print(f"  {key}: {value}")
            else:
                print("  No parsed data available")
            return True
    
    print(f"Packet {packet_id} not found")
    return False

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <json_file> <packet_id>")
        sys.exit(1)
    
    json_file = sys.argv[1]
    packet_id = sys.argv[2]
    
    if not find_packet(json_file, packet_id):
        sys.exit(1)
#!/usr/bin/env python3
"""
Send and Receive Packet Comparison Tool

This script compares packet definitions between the Send and Receive Go files
to identify common packet IDs and differences in their definitions.

Purpose:
    After converting packet definitions from Perl to Go format using convert_perl_to_go.py,
    this script analyzes the differences between the Send and Receive packet definitions.
    It identifies:
    - Packet IDs that are only in the Send file
    - Packet IDs that are only in the Receive file
    - Packet IDs that are common to both files
    - Differences in packet definitions for common packet IDs
    
    This is useful because the same packet ID can have different meanings in the Send and
    Receive directions, and this script helps identify those cases.

Usage:
    python3 compare_send_receive.py [--send SEND_FILE] [--receive RECEIVE_FILE] [--server-type SERVERTYPE]

Arguments:
    --send SEND_FILE         Path to the Send Go file
    --receive RECEIVE_FILE   Path to the Receive Go file
    --server-type SERVERTYPE Server type to compare (e.g., servertype0, servertype1)
                            This will look for:
                            - goKore/network/send/servers/SERVERTYPE.go
                            - goKore/network/receive/servers/SERVERTYPE.go

Output:
    Prints a comparison report showing:
    - Number of Send-only packet IDs
    - Number of Receive-only packet IDs
    - Number of common packet IDs
    - Detailed differences for common packet IDs
"""

import re
import os
import argparse

def extract_packet_ids(file_path):
    """Extract packet IDs from a Go file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Extract packet IDs
    pattern = r'packets\["([0-9A-F]{4})"\]\s*='
    packet_ids = set(re.findall(pattern, content))
    
    return packet_ids

def extract_packet_details(file_path):
    """Extract packet details from a Go file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Extract packet details
    pattern = r'packets\["([0-9A-F]{4})"\]\s*=\s*common\.PacketConstruction\{\s*ID:\s*"[0-9A-F]{4}",\s*Name:\s*"([^"]+)",\s*Format:\s*"([^"]*)",\s*FieldNames:\s*\[\]string\{([^}]*)\}'
    matches = re.findall(pattern, content, re.DOTALL)
    
    packets = {}
    for match in matches:
        packet_id = match[0]
        packet_name = match[1]
        format_string = match[2]
        field_names_str = match[3]
        
        # Extract field names
        field_names = []
        if field_names_str.strip():
            field_names_pattern = r'"([^"]+)"'
            field_names = re.findall(field_names_pattern, field_names_str)
        
        packets[packet_id] = {
            'name': packet_name,
            'format': format_string,
            'field_names': field_names
        }
    
    return packets

def compare_packets(send_file, receive_file):
    """Compare packet definitions between send and receive files."""
    # Extract packet IDs
    send_ids = extract_packet_ids(send_file)
    receive_ids = extract_packet_ids(receive_file)
    
    # Find common packet IDs
    common_ids = send_ids.intersection(receive_ids)
    
    # Extract packet details
    send_packets = extract_packet_details(send_file)
    receive_packets = extract_packet_details(receive_file)
    
    # Compare common packets
    differences = []
    for packet_id in sorted(common_ids):
        send_packet = send_packets[packet_id]
        receive_packet = receive_packets[packet_id]
        
        # Check for differences
        if (send_packet['name'] != receive_packet['name'] or
            send_packet['format'] != receive_packet['format'] or
            send_packet['field_names'] != receive_packet['field_names']):
            differences.append({
                'id': packet_id,
                'send': send_packet,
                'receive': receive_packet
            })
    
    return {
        'send_only': sorted(list(send_ids - receive_ids)),
        'receive_only': sorted(list(receive_ids - send_ids)),
        'common': sorted(list(common_ids)),
        'differences': differences
    }

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Compare Send and Receive packet definitions')
    parser.add_argument('--send', help='Path to the Send Go file')
    parser.add_argument('--receive', help='Path to the Receive Go file')
    parser.add_argument('--server-type', help='Server type to compare (e.g., servertype0, servertype1)')
    args = parser.parse_args()
    
    # Determine file paths based on arguments
    if args.server_type:
        server_type = args.server_type
        send_file = f"goKore/network/send/servers/{server_type}.go"
        receive_file = f"goKore/network/receive/servers/{server_type}.go"
    elif args.send and args.receive:
        send_file = args.send
        receive_file = args.receive
    else:
        # Default paths
        send_file = "goKore/network/send/servers/servertype0.go"
        receive_file = "goKore/network/receive/servers/servertype0.go"
    
    # Check if files exist
    if not os.path.exists(send_file):
        print(f"Error: Send file not found: {send_file}")
        return
    
    if not os.path.exists(receive_file):
        print(f"Error: Receive file not found: {receive_file}")
        return
    
    print(f"Comparing Send file: {send_file}")
    print(f"With Receive file: {receive_file}")
    
    # Compare packets
    result = compare_packets(send_file, receive_file)
    
    # Print results
    print("\nComparison Results:")
    print("-------------------")
    print(f"Send-only packet IDs: {len(result['send_only'])}")
    print(f"Receive-only packet IDs: {len(result['receive_only'])}")
    print(f"Common packet IDs: {len(result['common'])}")
    print(f"Packets with differences: {len(result['differences'])}")
    
    # Print packets with differences
    if result['differences']:
        print("\nPackets with differences:")
        for diff in result['differences']:
            packet_id = diff['id']
            send_packet = diff['send']
            receive_packet = diff['receive']
            
            print(f"\nPacket ID: {packet_id}")
            print(f"  Send name: {send_packet['name']}")
            print(f"  Receive name: {receive_packet['name']}")
            print(f"  Send format: {send_packet['format']}")
            print(f"  Receive format: {receive_packet['format']}")
            print(f"  Send field names: {send_packet['field_names']}")
            print(f"  Receive field names: {receive_packet['field_names']}")

if __name__ == "__main__":
    main()
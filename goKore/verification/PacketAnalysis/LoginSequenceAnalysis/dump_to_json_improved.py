#!/usr/bin/env python3
"""
DUMP to JSON Converter for Ragnarok Online Packet Analysis (Improved Version)

This script converts a Ragnarok Online packet dump file to a structured JSON format
containing all packet information. The output can be used as test data for integration
tests to validate that application behavior is deterministic and repeatable.
"""

import os
import re
import sys
import json
import argparse
import binascii
from dataclasses import dataclass, asdict, field
from typing import Dict, List, Optional, Tuple, Any

@dataclass
class PacketData:
    """Represents the raw and parsed data of a packet."""
    hex_bytes: List[str]  # List of hex bytes
    ascii_representation: str  # ASCII representation of the data
    binary: bytes = field(default=None, repr=False)  # Binary representation
    
    def __post_init__(self):
        # Convert hex bytes to binary
        if self.hex_bytes:
            self.binary = binascii.unhexlify(''.join(self.hex_bytes))
    
    def to_dict(self):
        """Convert to dictionary for JSON serialization."""
        return {
            "hex_bytes": self.hex_bytes,
            "ascii_representation": self.ascii_representation,
            "binary_base64": binascii.b2a_base64(self.binary).decode('ascii').strip() if self.binary else None
        }

@dataclass
class Packet:
    """Represents a single packet in the communication sequence."""
    direction: str  # "sent" or "received"
    packet_id: str  # Hex ID of the packet
    description: str  # Description of the packet
    size: int  # Size in bytes
    timestamp: str  # Timestamp when the packet was sent/received
    raw_data: List[PacketData]  # Raw packet data
    parsed_data: Optional[Dict[str, Any]] = None  # Any parsed information about the packet
    server_type: Optional[str] = None  # Account, Character, or Map server
    
    def to_dict(self):
        """Convert to dictionary for JSON serialization."""
        return {
            "direction": self.direction,
            "packet_id": self.packet_id,
            "description": self.description,
            "size": self.size,
            "timestamp": self.timestamp,
            "raw_data": [data.to_dict() for data in self.raw_data],
            "parsed_data": self.parsed_data,
            "server_type": self.server_type
        }

def parse_hex_data_line(line: str) -> Tuple[List[str], str]:
    """Parse a line of hex data from the dump file."""
    # Extract the hex bytes and ASCII representation
    match = re.search(r'\d+>  ((?:[0-9A-F]{2} ){8})   ([0-9A-F]{2} .+)', line)
    if match:
        hex_part1 = match.group(1).strip().split()
        hex_part2_and_ascii = match.group(2)
        
        # Split the second part into hex and ASCII
        hex_ascii_parts = hex_part2_and_ascii.split('  ')
        hex_part2 = hex_ascii_parts[0].strip().split()
        ascii_part = hex_ascii_parts[1] if len(hex_ascii_parts) > 1 else ""
        
        hex_bytes = hex_part1 + hex_part2
        return hex_bytes, ascii_part
    
    return [], ""

def parse_raw_data_block(raw_data_block: str) -> List[PacketData]:
    """Parse a block of raw data into structured format."""
    packet_data_list = []
    
    for line in raw_data_block.strip().split('\n'):
        hex_bytes, ascii_representation = parse_hex_data_line(line)
        if hex_bytes:
            packet_data = PacketData(
                hex_bytes=hex_bytes,
                ascii_representation=ascii_representation
            )
            packet_data_list.append(packet_data)
    
    return packet_data_list

def extract_packets(file_path: str) -> List[Packet]:
    """Extract all packets from a dump file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    packets = []
    current_server = "Unknown"
    
    # Track server connections
    if "Connecting to Account Server" in content:
        current_server = "Account Server"
    
    # Extract all packet exchanges
    lines = content.split('\n')
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Check for server transitions
        if "Connecting to Character Server" in line:
            current_server = "Character Server"
        elif "Connecting to Map Server" in line:
            current_server = "Map Server"
        
        # Extract sent packets
        sent_match = re.search(r'>> Sent packet: (\w+)\s+\[(.*?)\] \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})', line)
        if sent_match:
            packet_id, description, size, timestamp = sent_match.groups()
            
            # Find the raw data block that follows this packet header
            raw_data_block = ""
            j = i + 1
            while j < len(lines) and not lines[j].startswith('='):
                if re.match(r'\s+\d+>\s+', lines[j]):
                    raw_data_block += lines[j] + '\n'
                j += 1
            
            # Parse the raw data block
            raw_data = parse_raw_data_block(raw_data_block)
            
            packets.append(Packet(
                direction="sent",
                packet_id=packet_id,
                description=description,
                size=int(size),
                timestamp=timestamp,
                raw_data=raw_data,
                server_type=current_server
            ))
        
        # Extract received packets
        recv_match = re.search(r'<< Received packet:\s+(\w+) - (.*?) \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})', line)
        if recv_match:
            packet_id, description, size, timestamp = recv_match.groups()
            
            # Find the raw data block that follows this packet header
            raw_data_block = ""
            j = i + 1
            while j < len(lines) and not lines[j].startswith('='):
                if re.match(r'\s+\d+>\s+', lines[j]):
                    raw_data_block += lines[j] + '\n'
                j += 1
            
            # Parse the raw data block
            raw_data = parse_raw_data_block(raw_data_block)
            
            # Look for parsed data that might follow the raw data
            parsed_data = {}
            
            # Find the next non-raw-data line after the packet header
            j = i + 1
            while j < len(lines) and j < i + 50:  # Limit search to 50 lines
                if not re.match(r'\s+\d+>\s+', lines[j]) and lines[j].strip() and not lines[j].startswith('='):
                    # Check if this is the start of a parsed data section
                    if "----------" in lines[j]:
                        # Skip the header line
                        j += 1
                        # Collect all lines until the next separator or empty line
                        parsed_text = []
                        while j < len(lines) and not lines[j].startswith('=') and not lines[j].startswith('----------') and lines[j].strip():
                            parsed_text.append(lines[j])
                            j += 1
                        
                        # Process the parsed text
                        for line in parsed_text:
                            # Look for key-value pairs
                            kv_match = re.search(r'^([^:]+):\s*(.+)$', line)
                            if kv_match:
                                key = kv_match.group(1).strip()
                                value = kv_match.group(2).strip()
                                parsed_data[key] = value
                            else:
                                # If not a key-value pair, add as a note
                                note_key = f"note_{len(parsed_data)}"
                                parsed_data[note_key] = line.strip()
                    else:
                        # This is a single line of parsed data (like "Map Change: gef_fild07.gat (246, 191)")
                        parsed_data["info"] = lines[j].strip()
                        j += 1
                        # Check for more lines
                        while j < len(lines) and not lines[j].startswith('=') and lines[j].strip():
                            parsed_data[f"info_{j-i}"] = lines[j].strip()
                            j += 1
                    break
                j += 1
            
            packets.append(Packet(
                direction="received",
                packet_id=packet_id,
                description=description,
                size=int(size),
                timestamp=timestamp,
                raw_data=raw_data,
                parsed_data=parsed_data if parsed_data else None,
                server_type=current_server
            ))
        
        i += 1
    
    return packets

def main():
    """Main function to convert dump file to JSON."""
    parser = argparse.ArgumentParser(description="Convert a Ragnarok Online packet dump file to JSON format.")
    parser.add_argument("file", help="Path to the packet dump file")
    parser.add_argument("-o", "--output", help="Output JSON file path (default: <input_file>.json)")
    parser.add_argument("--pretty", action="store_true", help="Pretty-print the JSON output")
    args = parser.parse_args()
    
    if not os.path.exists(args.file):
        print(f"Error: File '{args.file}' not found.")
        return 1
    
    # Determine output file path
    output_file = args.output
    if not output_file:
        output_file = f"{args.file}.json"
    
    print(f"Converting {args.file} to JSON...")
    packets = extract_packets(args.file)
    
    if not packets:
        print("No packets found in the dump file.")
        return 1
    
    # Convert packets to dictionaries for JSON serialization
    packet_dicts = [packet.to_dict() for packet in packets]
    
    # Create the final JSON structure
    json_data = {
        "file_name": os.path.basename(args.file),
        "packet_count": len(packets),
        "packets": packet_dicts
    }
    
    # Write the JSON output
    with open(output_file, "w") as f:
        if args.pretty:
            json.dump(json_data, f, indent=2)
        else:
            json.dump(json_data, f)
    
    print(f"Conversion complete! JSON data written to {output_file}")
    print(f"Total packets extracted: {len(packets)}")
    
    # Print some statistics
    sent_packets = sum(1 for p in packets if p.direction == "sent")
    received_packets = sum(1 for p in packets if p.direction == "received")
    
    print(f"Sent packets: {sent_packets}")
    print(f"Received packets: {received_packets}")
    
    # Count packets by server type
    server_counts = {}
    for packet in packets:
        server_type = packet.server_type or "Unknown"
        server_counts[server_type] = server_counts.get(server_type, 0) + 1
    
    print("Packets by server type:")
    for server_type, count in server_counts.items():
        print(f"  {server_type}: {count}")
    
    # Count packets with parsed data
    parsed_data_count = sum(1 for p in packets if p.parsed_data)
    print(f"Packets with parsed data: {parsed_data_count}")
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
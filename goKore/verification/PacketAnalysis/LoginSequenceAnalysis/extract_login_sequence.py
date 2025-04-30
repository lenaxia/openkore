#!/usr/bin/env python3
"""
Extract Essential Login Sequence Packets

This script extracts only the essential packets needed for a successful login sequence
from a Ragnarok Online packet dump file. It helps to quickly analyze the critical
components of the login process without the noise of other packets.
"""

import os
import re
import sys
import argparse
from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple

@dataclass
class Packet:
    """Represents a single packet in the communication sequence."""
    direction: str  # "sent" or "received"
    packet_id: str  # Hex ID of the packet
    description: str  # Description of the packet
    size: int  # Size in bytes
    timestamp: str  # Timestamp when the packet was sent/received
    raw_data: str  # Raw hex data
    parsed_data: Optional[str] = None  # Any parsed information about the packet

# Define the essential packets for login sequence
ESSENTIAL_PACKETS = {
    # Account Server
    "sent:0064": "Account Server Login",
    "received:0AC4": "Account Info With Server Info",
    
    # Character Server
    "sent:0065": "Character Server Login",
    "received:082D": "Received characters from Game Login Server",
    "received:006B": "Received characters from Game Login Server",
    "received:08B9": "PinCode Request",
    "sent:0066": "Char Login",
    "received:0AC5": "Received character ID and Map IP",
    
    # Map Server
    "sent:0436": "Map Login",
    "received:0283": "Account ID",
    "received:02EB": "Enter Map",
    "sent:007D": "Map Loaded"
}

def extract_packets(file_path: str) -> List[Packet]:
    """Extract all packets from a dump file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    packets = []
    
    # Pattern for sent packets
    sent_pattern = r'>> Sent packet: (\w+)\s+\[(.*?)\] \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})'
    # Pattern for received packets
    recv_pattern = r'<< Received packet:\s+(\w+) - (.*?) \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})'
    
    # Extract raw data blocks
    data_blocks = re.findall(r'((?:  \d+>  (?:[0-9A-F]{2} ){8}   [0-9A-F]{2} .+\n)+)', content)
    
    # Extract sent packets
    for match in re.finditer(sent_pattern, content):
        packet_id, description, size, timestamp = match.groups()
        # Find the raw data block that follows this packet header
        pos = match.end()
        raw_data = ""
        for block in data_blocks:
            if content.find(block, pos, pos + 1000) != -1:
                raw_data = block
                break
        
        packets.append(Packet(
            direction="sent",
            packet_id=packet_id,
            description=description,
            size=int(size),
            timestamp=timestamp,
            raw_data=raw_data
        ))
    
    # Extract received packets
    for match in re.finditer(recv_pattern, content):
        packet_id, description, size, timestamp = match.groups()
        # Find the raw data block that follows this packet header
        pos = match.end()
        raw_data = ""
        for block in data_blocks:
            if content.find(block, pos, pos + 1000) != -1:
                raw_data = block
                break
        
        # Look for parsed data that might follow
        parsed_data = None
        next_section = content[pos:pos + 500]
        parsed_section = re.search(r'(?:----------.*?----------\n)(.*?)(?:\n-{10,}|\n={10,})', next_section, re.DOTALL)
        if parsed_section:
            parsed_data = parsed_section.group(1).strip()
        
        packets.append(Packet(
            direction="received",
            packet_id=packet_id,
            description=description,
            size=int(size),
            timestamp=timestamp,
            raw_data=raw_data,
            parsed_data=parsed_data
        ))
    
    return packets

def filter_essential_packets(packets: List[Packet]) -> List[Packet]:
    """Filter only the essential packets for login sequence."""
    essential = []
    for packet in packets:
        key = f"{packet.direction}:{packet.packet_id}"
        if key in ESSENTIAL_PACKETS:
            essential.append(packet)
    return essential

def format_packet_for_output(packet: Packet) -> str:
    """Format a packet for output."""
    direction = "→" if packet.direction == "sent" else "←"
    output = [f"{direction} {packet.packet_id}: {packet.description} [{packet.size} bytes] {packet.timestamp}"]
    
    if packet.raw_data:
        output.append("Raw data:")
        output.append(packet.raw_data.strip())
    
    if packet.parsed_data:
        output.append("Parsed data:")
        output.append(packet.parsed_data)
    
    return "\n".join(output)

def main():
    """Main function to extract essential login packets."""
    parser = argparse.ArgumentParser(description="Extract essential login sequence packets from a Ragnarok Online packet dump.")
    parser.add_argument("file", help="Path to the packet dump file")
    parser.add_argument("-o", "--output", help="Output file path (default: stdout)")
    args = parser.parse_args()
    
    if not os.path.exists(args.file):
        print(f"Error: File '{args.file}' not found.")
        return 1
    
    print(f"Extracting essential login packets from {args.file}...")
    packets = extract_packets(args.file)
    essential_packets = filter_essential_packets(packets)
    
    if not essential_packets:
        print("No essential login packets found in the dump file.")
        return 1
    
    output = ["# Essential Login Sequence Packets", ""]
    
    # Group packets by server
    account_server_packets = []
    char_server_packets = []
    map_server_packets = []
    
    for packet in essential_packets:
        if packet.packet_id in ["0064", "0AC4"]:
            account_server_packets.append(packet)
        elif packet.packet_id in ["0065", "082D", "006B", "08B9", "0066", "0AC5"]:
            char_server_packets.append(packet)
        else:
            map_server_packets.append(packet)
    
    # Add account server packets
    if account_server_packets:
        output.append("## Account Server")
        output.append("")
        for packet in account_server_packets:
            output.append(format_packet_for_output(packet))
            output.append("")
    
    # Add character server packets
    if char_server_packets:
        output.append("## Character Server")
        output.append("")
        for packet in char_server_packets:
            output.append(format_packet_for_output(packet))
            output.append("")
    
    # Add map server packets
    if map_server_packets:
        output.append("## Map Server")
        output.append("")
        for packet in map_server_packets:
            output.append(format_packet_for_output(packet))
            output.append("")
    
    # Write the output
    output_text = "\n".join(output)
    if args.output:
        with open(args.output, "w") as f:
            f.write(output_text)
        print(f"Essential login packets written to {args.output}")
    else:
        print(output_text)
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
#!/usr/bin/env python3
"""
Packet Analyzer for Ragnarok Online Login Sequence
This script analyzes packet dumps from successful login attempts to identify
the common sequence of packets required for a successful login.
"""

import os
import re
import sys
from collections import defaultdict
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

@dataclass
class LoginSequence:
    """Represents a complete login sequence with all packets."""
    packets: List[Packet]
    account_server_packets: List[Packet]
    char_server_packets: List[Packet]
    map_server_packets: List[Packet]

def parse_dump_file(file_path: str) -> LoginSequence:
    """Parse a dump file and extract all packets."""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split the content into sections based on server connections
    account_server_section = re.search(r'Connecting to Account Server.*?Closing connection to Account Server', 
                                      content, re.DOTALL)
    char_server_section = re.search(r'Connecting to Character Server.*?Closing connection to Character Server', 
                                   content, re.DOTALL)
    map_server_section = re.search(r'Connecting to Map Server.*', 
                                  content, re.DOTALL)
    
    # Extract all packets
    all_packets = []
    account_server_packets = []
    char_server_packets = []
    map_server_packets = []
    
    if account_server_section:
        account_server_packets = extract_packets(account_server_section.group(0))
        all_packets.extend(account_server_packets)
    
    if char_server_section:
        char_server_packets = extract_packets(char_server_section.group(0))
        all_packets.extend(char_server_packets)
    
    if map_server_section:
        map_server_packets = extract_packets(map_server_section.group(0))
        all_packets.extend(map_server_packets)
    
    return LoginSequence(
        packets=all_packets,
        account_server_packets=account_server_packets,
        char_server_packets=char_server_packets,
        map_server_packets=map_server_packets
    )

def extract_packets(content: str) -> List[Packet]:
    """Extract all packets from a section of the dump file."""
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

def analyze_login_sequences(sequences: List[LoginSequence]) -> Dict:
    """Analyze multiple login sequences to find common patterns."""
    # Count packet types in each server connection
    account_server_packets = defaultdict(int)
    char_server_packets = defaultdict(int)
    map_server_packets = defaultdict(int)
    
    # Track the order of packets
    account_server_flow = []
    char_server_flow = []
    map_server_flow = []
    
    for seq in sequences:
        # Account server packets
        flow = []
        for packet in seq.account_server_packets:
            key = f"{packet.direction}:{packet.packet_id}"
            account_server_packets[key] += 1
            flow.append(key)
        if flow and flow not in account_server_flow:
            account_server_flow.append(flow)
        
        # Character server packets
        flow = []
        for packet in seq.char_server_packets:
            key = f"{packet.direction}:{packet.packet_id}"
            char_server_packets[key] += 1
            flow.append(key)
        if flow and flow not in char_server_flow:
            char_server_flow.append(flow)
        
        # Map server packets
        flow = []
        for packet in seq.map_server_packets:
            key = f"{packet.direction}:{packet.packet_id}"
            map_server_packets[key] += 1
            flow.append(key)
        if flow and flow not in map_server_flow:
            map_server_flow.append(flow)
    
    return {
        "account_server": {
            "packets": dict(account_server_packets),
            "flows": account_server_flow
        },
        "char_server": {
            "packets": dict(char_server_packets),
            "flows": char_server_flow
        },
        "map_server": {
            "packets": dict(map_server_packets),
            "flows": map_server_flow
        }
    }

def find_common_sequence(sequences: List[LoginSequence]) -> Dict:
    """Find the common packet sequence across all login attempts."""
    # Extract just the packet IDs in order for each server connection
    account_flows = []
    char_flows = []
    map_flows = []
    
    for seq in sequences:
        account_flows.append([f"{p.direction}:{p.packet_id}" for p in seq.account_server_packets])
        char_flows.append([f"{p.direction}:{p.packet_id}" for p in seq.char_server_packets])
        map_flows.append([f"{p.direction}:{p.packet_id}" for p in seq.map_server_packets])
    
    # Find the longest common subsequence for each server
    account_common = find_longest_common_subsequence(account_flows)
    char_common = find_longest_common_subsequence(char_flows)
    map_common = find_longest_common_subsequence(map_flows)
    
    return {
        "account_server": account_common,
        "char_server": char_common,
        "map_server": map_common
    }

def find_longest_common_subsequence(sequences: List[List[str]]) -> List[str]:
    """Find the longest common subsequence among multiple sequences."""
    if not sequences:
        return []
    
    # Start with the first sequence as the common one
    common = sequences[0]
    
    # Compare with each other sequence
    for seq in sequences[1:]:
        common = lcs(common, seq)
    
    return common

def lcs(seq1: List[str], seq2: List[str]) -> List[str]:
    """Find the longest common subsequence between two sequences."""
    m, n = len(seq1), len(seq2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    
    # Fill the dp table
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if seq1[i - 1] == seq2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])
    
    # Reconstruct the LCS
    i, j = m, n
    lcs_result = []
    
    while i > 0 and j > 0:
        if seq1[i - 1] == seq2[j - 1]:
            lcs_result.append(seq1[i - 1])
            i -= 1
            j -= 1
        elif dp[i - 1][j] > dp[i][j - 1]:
            i -= 1
        else:
            j -= 1
    
    return list(reversed(lcs_result))

def generate_summary(sequences: List[LoginSequence], common_sequence: Dict) -> str:
    """Generate a summary of the login sequence."""
    # Get packet descriptions for the common sequence
    packet_descriptions = {}
    for seq in sequences:
        for packet in seq.packets:
            key = f"{packet.direction}:{packet.packet_id}"
            if key not in packet_descriptions:
                packet_descriptions[key] = packet.description
    
    summary = []
    summary.append("# Ragnarok Online Login Sequence Analysis")
    summary.append("\n## Account Server Login Sequence")
    for packet_key in common_sequence["account_server"]:
        direction = "→" if packet_key.startswith("sent") else "←"
        packet_id = packet_key.split(":")[1]
        desc = packet_descriptions.get(packet_key, "Unknown")
        summary.append(f"{direction} {packet_id}: {desc}")
    
    summary.append("\n## Character Server Login Sequence")
    for packet_key in common_sequence["char_server"]:
        direction = "→" if packet_key.startswith("sent") else "←"
        packet_id = packet_key.split(":")[1]
        desc = packet_descriptions.get(packet_key, "Unknown")
        summary.append(f"{direction} {packet_id}: {desc}")
    
    summary.append("\n## Map Server Login Sequence")
    for packet_key in common_sequence["map_server"]:
        direction = "→" if packet_key.startswith("sent") else "←"
        packet_id = packet_key.split(":")[1]
        desc = packet_descriptions.get(packet_key, "Unknown")
        summary.append(f"{direction} {packet_id}: {desc}")
    
    return "\n".join(summary)

def main():
    """Main function to analyze the dump files."""
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    dump_dir = os.path.join(base_dir, "dumps")
    dump_files = [os.path.join(dump_dir, f) for f in os.listdir(dump_dir) if f.startswith("DUMP")]
    
    if not dump_files:
        print("No dump files found!")
        return
    
    print(f"Found {len(dump_files)} dump files: {', '.join(os.path.basename(f) for f in dump_files)}")
    
    # Parse all dump files
    sequences = []
    for file_path in dump_files:
        print(f"Parsing {os.path.basename(file_path)}...")
        sequences.append(parse_dump_file(file_path))
    
    # Analyze the sequences
    analysis = analyze_login_sequences(sequences)
    common_sequence = find_common_sequence(sequences)
    
    # Generate summary
    summary = generate_summary(sequences, common_sequence)
    
    # Write the results
    output_dir = os.path.join(base_dir, "analysis")
    
    with open(os.path.join(output_dir, "login_sequence_summary.md"), "w") as f:
        f.write(summary)
    
    print("Analysis complete! Results written to analysis/login_sequence_summary.md")

if __name__ == "__main__":
    main()
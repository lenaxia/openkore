#!/usr/bin/env python3
"""
Packet Structure Analyzer for Ragnarok Online Login Sequence
This script analyzes the structure of key packets in the login sequence
by comparing multiple instances of the same packet across different login attempts.
"""

import os
import re
import sys
from collections import defaultdict
from dataclasses import dataclass
from typing import Dict, List, Optional, Set, Tuple

@dataclass
class PacketField:
    """Represents a field within a packet structure."""
    offset: int
    length: int
    values: List[bytes]
    description: str = ""
    is_variable: bool = False

@dataclass
class PacketStructure:
    """Represents the structure of a packet type."""
    packet_id: str
    direction: str
    description: str
    fields: List[PacketField]
    raw_examples: List[str]

def parse_hex_data(raw_data: str) -> bytes:
    """Parse hex data from the raw packet dump format."""
    hex_bytes = []
    for line in raw_data.strip().split('\n'):
        # Extract the hex bytes from the line
        match = re.search(r'\d+>  ((?:[0-9A-F]{2} ){8})   ([0-9A-F]{2} .+)', line)
        if match:
            hex_part1 = match.group(1).strip()
            hex_part2 = match.group(2).split('  ')[0].strip()
            hex_bytes.extend(hex_part1.split())
            hex_bytes.extend(hex_part2.split())
    
    return bytes([int(h, 16) for h in hex_bytes])

def extract_packet_instances(dump_files: List[str], packet_id: str, direction: str) -> Tuple[List[bytes], List[str], str]:
    """Extract all instances of a specific packet type from the dump files."""
    packet_instances = []
    raw_examples = []
    description = ""
    
    for file_path in dump_files:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Pattern for the packet header
        if direction == "sent":
            pattern = fr'>> Sent packet: {packet_id}\s+\[(.*?)\] \[(\d+) bytes\]'
        else:
            pattern = fr'<< Received packet:\s+{packet_id} - (.*?) \[(\d+) bytes\]'
        
        for match in re.finditer(pattern, content):
            desc = match.group(1)
            if not description:
                description = desc
            
            # Find the raw data block that follows this packet header
            pos = match.end()
            end_pos = content.find('=====', pos)
            if end_pos == -1:
                end_pos = content.find('\n\n', pos)
            
            if end_pos != -1:
                raw_data_section = content[pos:end_pos]
                # Extract the raw hex data
                data_block = re.search(r'((?:  \d+>  (?:[0-9A-F]{2} ){8}   [0-9A-F]{2} .+\n)+)', raw_data_section)
                if data_block:
                    raw_data = data_block.group(1)
                    raw_examples.append(raw_data)
                    packet_instances.append(parse_hex_data(raw_data))
    
    return packet_instances, raw_examples, description

def analyze_packet_structure(packet_instances: List[bytes], packet_id: str, direction: str, description: str, raw_examples: List[str]) -> PacketStructure:
    """Analyze the structure of a packet by comparing multiple instances."""
    if not packet_instances:
        return PacketStructure(packet_id=packet_id, direction=direction, description=description, fields=[], raw_examples=[])
    
    # Determine the minimum length of all instances
    min_length = min(len(p) for p in packet_instances)
    
    # Analyze each byte position
    fields = []
    current_field = None
    
    for i in range(min_length):
        values = [p[i] for p in packet_instances]
        is_constant = len(set(values)) == 1
        
        # If this is the first byte or the constancy changed, start a new field
        if current_field is None or current_field.is_variable != (not is_constant):
            if current_field:
                fields.append(current_field)
            
            current_field = PacketField(
                offset=i,
                length=1,
                values=[bytes([v]) for v in values],
                is_variable=not is_constant
            )
        else:
            # Extend the current field
            current_field.length += 1
            current_field.values = [v + bytes([values[j]]) for j, v in enumerate(current_field.values)]
    
    # Add the last field
    if current_field:
        fields.append(current_field)
    
    # Add descriptions to fields based on common patterns
    for field in fields:
        if field.offset == 0 and field.length == 2:
            field.description = "Packet ID"
        elif field.offset == 2 and field.length == 2:
            field.description = "Packet Length"
    
    return PacketStructure(
        packet_id=packet_id,
        direction=direction,
        description=description,
        fields=fields,
        raw_examples=raw_examples[:3]  # Limit to 3 examples to keep the output manageable
    )

def generate_packet_structure_report(structures: Dict[str, PacketStructure]) -> str:
    """Generate a detailed report of packet structures."""
    report = []
    report.append("# Ragnarok Online Packet Structure Analysis")
    report.append("\n## Key Login Sequence Packets\n")
    
    for key, structure in structures.items():
        direction = "→ Sent to server" if structure.direction == "sent" else "← Received from server"
        report.append(f"### {structure.packet_id} - {structure.description}")
        report.append(f"**Direction:** {direction}")
        report.append(f"**Packet ID:** 0x{structure.packet_id}")
        report.append("\n**Structure:**\n")
        report.append("| Offset | Length | Description | Constant | Values |")
        report.append("|--------|--------|-------------|----------|--------|")
        
        for field in structure.fields:
            constant = "Yes" if not field.is_variable else "No"
            if field.is_variable:
                values = "Variable"
            else:
                hex_value = field.values[0].hex().upper()
                hex_value = ' '.join(hex_value[i:i+2] for i in range(0, len(hex_value), 2))
                values = f"0x{hex_value}"
            
            report.append(f"| {field.offset} | {field.length} | {field.description} | {constant} | {values} |")
        
        report.append("\n**Raw Examples:**\n")
        for i, example in enumerate(structure.raw_examples):
            report.append(f"Example {i+1}:")
            report.append("```")
            report.append(example.strip())
            report.append("```\n")
    
    return "\n".join(report)

def main():
    """Main function to analyze packet structures."""
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    dump_dir = os.path.join(base_dir, "dumps")
    dump_files = [os.path.join(dump_dir, f) for f in os.listdir(dump_dir) if f.startswith("DUMP")]
    
    if not dump_files:
        print("No dump files found!")
        return
    
    print(f"Found {len(dump_files)} dump files: {', '.join(os.path.basename(f) for f in dump_files)}")
    
    # Key packets in the login sequence
    key_packets = [
        ("0064", "sent"),     # Account Server Login
        ("0AC4", "received"), # Account Info With Server Info
        ("0065", "sent"),     # Character Server Login
        ("082D", "received"), # Received characters from Game Login Server
        ("006B", "received"), # Received characters from Game Login Server
        ("08B9", "received"), # PinCode Request
        ("0066", "sent"),     # Char Login
        ("0AC5", "received"), # Received character ID and Map IP
        ("0436", "sent"),     # Map Login
        ("02EB", "received"), # Enter Map
        ("007D", "sent"),     # Map Loaded
    ]
    
    # Analyze each key packet
    structures = {}
    for packet_id, direction in key_packets:
        print(f"Analyzing {direction} packet {packet_id}...")
        instances, raw_examples, description = extract_packet_instances(dump_files, packet_id, direction)
        if instances:
            structure = analyze_packet_structure(instances, packet_id, direction, description, raw_examples)
            structures[f"{direction}:{packet_id}"] = structure
        else:
            print(f"  No instances of {packet_id} found!")
    
    # Generate the report
    report = generate_packet_structure_report(structures)
    
    # Write the results
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    output_dir = os.path.join(base_dir, "analysis")
    with open(os.path.join(output_dir, "packet_structures.md"), "w") as f:
        f.write(report)
    
    print("Analysis complete! Results written to analysis/packet_structures.md")

if __name__ == "__main__":
    main()
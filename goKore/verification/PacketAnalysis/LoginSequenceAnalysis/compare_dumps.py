#!/usr/bin/env python3
"""
Compare Packet Sequences Across Multiple Dumps

This script compares the packet sequences across multiple JSON dump files
to identify any differences or inconsistencies in the login process.
"""

import os
import sys
import json
import argparse
from typing import Dict, List, Set, Tuple

def load_json_dump(file_path: str) -> Dict:
    """Load a JSON dump file."""
    with open(file_path, 'r') as f:
        return json.load(f)

def extract_packet_sequence(dump_data: Dict) -> List[Tuple[str, str, str]]:
    """Extract the packet sequence from a dump file."""
    sequence = []
    for packet in dump_data['packets']:
        # Create a tuple of (direction, packet_id, server_type)
        sequence.append((
            packet['direction'],
            packet['packet_id'],
            packet['server_type'] or "Unknown"
        ))
    return sequence

def compare_sequences(sequences: Dict[str, List[Tuple[str, str, str]]]) -> Dict:
    """Compare packet sequences across multiple dumps."""
    results = {
        "common_packets": [],
        "unique_packets": {},
        "sequence_differences": {},
        "server_transitions": {}
    }
    
    # Find common packets across all sequences
    all_packets = set()
    for dump_name, sequence in sequences.items():
        for packet in sequence:
            all_packets.add(packet)
    
    # Check which packets are common to all dumps
    for packet in all_packets:
        is_common = True
        for dump_name, sequence in sequences.items():
            if packet not in sequence:
                is_common = False
                break
        
        if is_common:
            results["common_packets"].append(packet)
        else:
            # Track which dumps have this packet
            for dump_name, sequence in sequences.items():
                if packet in sequence:
                    if packet not in results["unique_packets"]:
                        results["unique_packets"][packet] = []
                    results["unique_packets"][packet].append(dump_name)
    
    # Check for sequence differences
    # For each common packet, check if it appears in the same order in all dumps
    common_packets = results["common_packets"]
    for i, packet in enumerate(common_packets):
        if i == 0:
            continue  # Skip the first packet
        
        prev_packet = common_packets[i-1]
        
        for dump_name, sequence in sequences.items():
            try:
                packet_idx = sequence.index(packet)
                prev_idx = sequence.index(prev_packet)
                
                # Check if the order is different
                if packet_idx <= prev_idx:
                    if packet not in results["sequence_differences"]:
                        results["sequence_differences"][packet] = {}
                    
                    results["sequence_differences"][packet][dump_name] = {
                        "expected_after": prev_packet,
                        "actual_position": packet_idx,
                        "expected_position": prev_idx + 1
                    }
            except ValueError:
                # This shouldn't happen for common packets
                pass
    
    # Track server transitions
    for dump_name, sequence in sequences.items():
        transitions = []
        current_server = None
        
        for direction, packet_id, server_type in sequence:
            if server_type != current_server:
                if current_server is not None:
                    transitions.append((current_server, server_type))
                current_server = server_type
        
        results["server_transitions"][dump_name] = transitions
    
    return results

def generate_report(results: Dict, sequences: Dict[str, List[Tuple[str, str, str]]]) -> str:
    """Generate a report of the comparison results."""
    report = []
    report.append("# Packet Sequence Comparison Report")
    report.append("\n## Overview")
    
    # Add basic statistics
    report.append(f"Number of dumps compared: {len(sequences)}")
    for dump_name, sequence in sequences.items():
        report.append(f"- {dump_name}: {len(sequence)} packets")
    
    # Common packets
    report.append("\n## Common Packets")
    report.append(f"Number of packets common to all dumps: {len(results['common_packets'])}")
    
    # Group common packets by server type
    server_packets = {}
    for direction, packet_id, server_type in results['common_packets']:
        if server_type not in server_packets:
            server_packets[server_type] = []
        server_packets[server_type].append((direction, packet_id))
    
    for server_type, packets in server_packets.items():
        report.append(f"\n### {server_type}")
        for direction, packet_id in packets:
            arrow = "→" if direction == "sent" else "←"
            report.append(f"- {arrow} {packet_id}")
    
    # Unique packets
    if results['unique_packets']:
        report.append("\n## Unique Packets")
        report.append("These packets appear in some dumps but not others:")
        
        for (direction, packet_id, server_type), dump_names in results['unique_packets'].items():
            arrow = "→" if direction == "sent" else "←"
            report.append(f"- {arrow} {packet_id} ({server_type}): Found in {', '.join(dump_names)}")
    
    # Sequence differences
    if results['sequence_differences']:
        report.append("\n## Sequence Differences")
        report.append("These packets appear in a different order in some dumps:")
        
        for packet, differences in results['sequence_differences'].items():
            direction, packet_id, server_type = packet
            arrow = "→" if direction == "sent" else "←"
            report.append(f"\n### {arrow} {packet_id} ({server_type})")
            
            for dump_name, diff in differences.items():
                expected_direction, expected_packet, expected_server = diff['expected_after']
                expected_arrow = "→" if expected_direction == "sent" else "←"
                
                report.append(f"- In {dump_name}:")
                report.append(f"  - Expected after: {expected_arrow} {expected_packet} ({expected_server})")
                report.append(f"  - Actual position: {diff['actual_position']}")
                report.append(f"  - Expected position: {diff['expected_position']}")
    
    # Server transitions
    report.append("\n## Server Transitions")
    
    all_transitions = set()
    for dump_name, transitions in results['server_transitions'].items():
        for transition in transitions:
            all_transitions.add(transition)
    
    if all_transitions:
        report.append("Server transitions in each dump:")
        
        for dump_name, transitions in results['server_transitions'].items():
            report.append(f"\n### {dump_name}")
            for from_server, to_server in transitions:
                report.append(f"- {from_server} → {to_server}")
    
    return "\n".join(report)

def main():
    """Main function to compare packet sequences across multiple dumps."""
    parser = argparse.ArgumentParser(description="Compare packet sequences across multiple JSON dump files.")
    parser.add_argument("files", nargs="+", help="JSON dump files to compare")
    parser.add_argument("-o", "--output", help="Output report file path")
    args = parser.parse_args()
    
    # Load all dump files
    dumps = {}
    for file_path in args.files:
        if not os.path.exists(file_path):
            print(f"Error: File '{file_path}' not found.")
            return 1
        
        dump_name = os.path.basename(file_path)
        dumps[dump_name] = load_json_dump(file_path)
    
    print(f"Comparing {len(dumps)} dump files...")
    
    # Extract packet sequences
    sequences = {}
    for dump_name, dump_data in dumps.items():
        sequences[dump_name] = extract_packet_sequence(dump_data)
    
    # Compare sequences
    results = compare_sequences(sequences)
    
    # Generate report
    report = generate_report(results, sequences)
    
    # Write or print the report
    if args.output:
        with open(args.output, "w") as f:
            f.write(report)
        print(f"Comparison report written to {args.output}")
    else:
        print(report)
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
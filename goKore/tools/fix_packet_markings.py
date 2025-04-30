#!/usr/bin/env python3
"""
Packet Markings Fixer

This script fixes packet markings in the packets file based on the implementation status in the Go file.

Purpose:
    The script reads a Go file containing packet implementations and a packets file with packet definitions.
    It updates the packet markings in the packets file to reflect whether each packet is implemented or not.
    
    Packet markings in the packets file look like:
    [X] '0064' => ['packet_name', ...] - Implemented packet
    [ ] '0065' => ['packet_name', ...] - Not implemented packet
    
    The script changes [ ] to [X] for implemented packets and [X] to [ ] for not implemented packets.

Usage:
    python3 fix_packet_markings.py [--packets PACKETS_FILE] [--go GO_FILE]

Arguments:
    --packets PACKETS_FILE  Path to the packets file (default: network/send/servers/servertype0.packets)
    --go GO_FILE           Path to the Go file (default: network/send/servers/servertype0.go)
"""

import re
import os
import argparse

def extract_packets_from_go_file(file_path):
    """Extract packet codes from the Go file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Regular expression to match packet codes
    pattern = r'packets\["([0-9A-F]{4})"\]'
    matches = re.findall(pattern, content)
    
    return set(matches)

def fix_packet_markings(packets_file_path, go_file_path):
    """Fix packet markings in the packets file."""
    # Extract packet codes from Go file
    implemented_packets = extract_packets_from_go_file(go_file_path)
    
    # Read the packets file
    with open(packets_file_path, 'r') as f:
        lines = f.readlines()
    
    # Regular expression to match packet definitions
    pattern = r"(\[)([X ])(]\s*)'([0-9A-F]{4})'"
    
    # Count variables
    fixed_to_implemented = 0
    fixed_to_not_implemented = 0
    already_correct = 0
    
    # Process each line
    new_lines = []
    for line in lines:
        match = re.search(pattern, line)
        if match:
            prefix = match.group(1)
            current_mark = match.group(2)
            suffix = match.group(3)
            packet_code = match.group(4)
            
            is_implemented = packet_code in implemented_packets
            
            if is_implemented and current_mark != 'X':
                # Change [ ] to [X]
                new_line = line.replace(f"[{current_mark}]", f"[X]", 1)
                fixed_to_implemented += 1
                print(f"Marking {packet_code} as implemented")
            elif not is_implemented and current_mark == 'X':
                # Change [X] to [ ]
                new_line = line.replace(f"[{current_mark}]", f"[ ]", 1)
                fixed_to_not_implemented += 1
            else:
                # Already correct
                new_line = line
                already_correct += 1
        else:
            # Not a packet definition line
            new_line = line
        
        new_lines.append(new_line)
    
    # Write the updated content to a new file
    new_file_path = packets_file_path + '.new'
    with open(new_file_path, 'w') as f:
        f.writelines(new_lines)
    
    # Replace the original file with the new file
    os.rename(new_file_path, packets_file_path)
    
    # Print statistics
    print(f"Fixed {fixed_to_implemented} packets to mark as implemented")
    print(f"Fixed {fixed_to_not_implemented} packets to mark as not implemented")
    print(f"Found {already_correct} packets already correctly marked")
    print(f"Total packets processed: {fixed_to_implemented + fixed_to_not_implemented + already_correct}")

def handle_special_cases(packets_file_path):
    """Handle special cases for packet markings."""
    with open(packets_file_path, 'r') as f:
        lines = f.readlines()
    
    modified = False
    for i, line in enumerate(lines):
        if "# 'C350'" in line and "[ ]" in line:
            lines[i] = line.replace("[ ]", "[X]")
            print(f"Fixed C350 packet marking (special case)")
            modified = True
            break
    
    if modified:
        with open(packets_file_path, 'w') as f:
            f.writelines(lines)

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Fix packet markings in the packets file')
    parser.add_argument('--packets', default="network/send/servers/servertype0.packets",
                        help='Path to the packets file')
    parser.add_argument('--go', default="network/send/servers/servertype0.go",
                        help='Path to the Go file')
    args = parser.parse_args()
    
    packets_file_path = args.packets
    go_file_path = args.go
    
    # Check if files exist
    if not os.path.exists(packets_file_path):
        print(f"Error: Packets file not found: {packets_file_path}")
        return
    
    if not os.path.exists(go_file_path):
        print(f"Error: Go file not found: {go_file_path}")
        return
    
    # Fix packet markings
    fix_packet_markings(packets_file_path, go_file_path)
    
    # Handle special cases
    handle_special_cases(packets_file_path)
    
    print(f"Packet markings fixed in {packets_file_path}")

if __name__ == "__main__":
    main()
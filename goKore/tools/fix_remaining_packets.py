#!/usr/bin/env python3
"""
Remaining Packets Fixer

This script fixes specific packet markings in the packets file based on predefined rules.

Purpose:
    The script reads a packets file and applies specific fixes to certain packet definitions
    based on their packet codes. This is useful for handling special cases that cannot be
    automatically determined by the fix_packet_markings.py script.
    
    For example, some packets might be marked as implemented ([X]) but should actually be
    marked as not implemented ([ ]) due to specific requirements or compatibility issues.

Usage:
    python3 fix_remaining_packets.py [--packets PACKETS_FILE] [--by-code]

Arguments:
    --packets PACKETS_FILE  Path to the packets file (default: network/send/servers/servertype0.packets)
    --by-code              Use packet codes instead of line numbers for fixes (more robust)

Note:
    Using --by-code is recommended as it's more robust against changes in the file structure.
    Line numbers can change if the file is modified, but packet codes remain the same.
"""

import re
import argparse
import os

def fix_packet_markings_by_line(packets_file_path):
    """Fix packet markings using line numbers."""
    # Read the packets file
    with open(packets_file_path, 'r') as f:
        lines = f.readlines()
    
    # Process each line
    fixed_lines = []
    fixed_count = 0
    
    # Fix specific lines by line number
    line_fixes = {
        267: lambda line: line.replace("[X]", "[ ]"),  # '01F8' => ['adopt_unknown']
        394: lambda line: line.replace("[X]", "[ ]"),  # '02E1' => ['actor_action', ...]
        443: lambda line: line.replace("[X]", "[ ]"),  # '07E6' => ['captcha_session_ID', ...]
        17: lambda line: line.replace("[X]", "[ ]"),   # '0078' => ['actor_exists', ...]
        19: lambda line: line.replace("[X]", "[ ]"),   # '0079' => ['actor_connected', ...]
        22: lambda line: line.replace("[X]", "[ ]"),   # '007B' => ['actor_moved', ...]
        23: lambda line: line.replace("[X]", "[ ]"),   # '007C' => ['actor_exists', ...]
    }
    
    for i, line in enumerate(lines):
        line_num = i + 1
        if line_num in line_fixes:
            new_line = line_fixes[line_num](line)
            if new_line != line:
                fixed_count += 1
                print(f"Fixed line {line_num}: {line.strip()}")
            fixed_lines.append(new_line)
        else:
            fixed_lines.append(line)
    
    # Write the updated content back to the file
    with open(packets_file_path, 'w') as f:
        f.writelines(fixed_lines)
    
    print(f"Fixed {fixed_count} packet markings by line number")
    return fixed_count

def fix_packet_markings_by_code(packets_file_path):
    """Fix packet markings using packet codes (more robust)."""
    # Read the packets file
    with open(packets_file_path, 'r') as f:
        lines = f.readlines()
    
    # Process each line
    fixed_lines = []
    fixed_count = 0
    
    # Fix specific packets by code
    packet_fixes = {
        '01F8': "[ ]",  # adopt_unknown
        '02E1': "[ ]",  # actor_action
        '07E6': "[ ]",  # captcha_session_ID
        '0078': "[ ]",  # actor_exists
        '0079': "[ ]",  # actor_connected
        '007B': "[ ]",  # actor_moved
        '007C': "[ ]",  # actor_exists
    }
    
    # Regular expression to match packet definitions
    pattern = r"(\[)([X ])(]\s*)'([0-9A-F]{4})'"
    
    for line in lines:
        match = re.search(pattern, line)
        if match:
            packet_code = match.group(4)
            if packet_code in packet_fixes:
                current_mark = match.group(2)
                target_mark = packet_fixes[packet_code].strip('[]')
                
                if current_mark != target_mark:
                    new_line = line.replace(f"[{current_mark}]", packet_fixes[packet_code], 1)
                    fixed_count += 1
                    print(f"Fixed packet {packet_code}: {line.strip()}")
                    fixed_lines.append(new_line)
                else:
                    fixed_lines.append(line)
            else:
                fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    # Write the updated content back to the file
    with open(packets_file_path, 'w') as f:
        f.writelines(fixed_lines)
    
    print(f"Fixed {fixed_count} packet markings by packet code")
    return fixed_count

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Fix remaining packet markings in the packets file')
    parser.add_argument('--packets', default="network/send/servers/servertype0.packets",
                        help='Path to the packets file')
    parser.add_argument('--by-code', action='store_true',
                        help='Use packet codes instead of line numbers for fixes (more robust)')
    args = parser.parse_args()
    
    packets_file_path = args.packets
    
    # Check if file exists
    if not os.path.exists(packets_file_path):
        print(f"Error: Packets file not found: {packets_file_path}")
        return
    
    # Fix packet markings
    if args.by_code:
        fix_packet_markings_by_code(packets_file_path)
    else:
        fix_packet_markings_by_line(packets_file_path)
    
    print(f"Packet markings fixed in {packets_file_path}")

if __name__ == "__main__":
    main()
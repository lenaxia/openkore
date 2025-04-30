#!/usr/bin/env python3
"""
Packet Format Converter

This script converts packet definitions in servertype0.go from one format to another.

Purpose:
    The script reads the original servertype0.go file which contains packet definitions
    in two different formats:
    
    Format 1 (older):
    "0064": {
        ID:         "0064",
        Name:       "master_login",
        Format:     "v a24 a24 C",
        FieldNames: []string{"version", "username", "password", "clienttype"},
    },
    
    Format 2 (newer):
    packets["006B"] = common.PacketConstruction{
        ID:         "006B",
        Name:       "received_characters_info",
        Format:     "v C3 x20 a*",
        FieldNames: []string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
    }
    
    The script converts all packet definitions to Format 2 for consistency.

Usage:
    python3 convert_packet_format.py

Output:
    Creates a new file servertype0.5.go with all packet definitions in Format 2.
"""

import re
import os

def main():
    # Path to the input and output files
    input_file = "../network/send/servers/servertype0.go"
    output_file = "../network/send/servers/servertype0.5.go"
    
    # Read the input file
    with open(input_file, 'r') as f:
        content = f.read()
    
    # Extract the package declaration, imports, and function signature
    header_match = re.search(r'(// Package.*?func ServerType0PacketConstructions\(\) map\[string\]common\.PacketConstruction \{)', content, re.DOTALL)
    if not header_match:
        print("Error: Could not find the function header")
        return
    
    header = header_match.group(1)
    
    # Extract all packet definitions in the first format
    # This regex looks for patterns like "0064": { ... },
    first_format_pattern = r'"([0-9A-F]{4})"\s*:\s*\{\s*ID:\s*"([0-9A-F]{4})",\s*Name:\s*"([^"]+)",\s*Format:\s*"([^"]*)",\s*FieldNames:\s*\[\]string\{([^}]*)\},?\s*\}'
    first_format_matches = re.finditer(first_format_pattern, content)
    
    # Extract all packet definitions in the second format
    # This regex looks for patterns like packets["006B"] = common.PacketConstruction{ ... }
    second_format_pattern = r'packets\["([0-9A-F]{4})"\]\s*=\s*common\.PacketConstruction\{\s*ID:\s*"([0-9A-F]{4})",\s*Name:\s*"([^"]+)",\s*Format:\s*"([^"]*)",\s*FieldNames:\s*\[\]string\{([^}]*)\},?\s*\}'
    second_format_matches = re.finditer(second_format_pattern, content)
    
    # Store all packet definitions
    packets = {}
    
    # Process first format matches
    for match in first_format_matches:
        packet_id = match.group(1)
        name = match.group(3)
        format_str = match.group(4)
        field_names = match.group(5)
        
        packets[packet_id] = {
            'id': packet_id,
            'name': name,
            'format': format_str,
            'field_names': field_names
        }
    
    # Process second format matches
    for match in second_format_matches:
        packet_id = match.group(1)
        name = match.group(3)
        format_str = match.group(4)
        field_names = match.group(5)
        
        packets[packet_id] = {
            'id': packet_id,
            'name': name,
            'format': format_str,
            'field_names': field_names
        }
    
    # Generate the output file content
    output_content = header + "\n"
    output_content += "\tpackets := map[string]common.PacketConstruction{}\n\n"
    
    # Add all packet definitions in the second format
    for packet_id in sorted(packets.keys()):
        packet = packets[packet_id]
        output_content += f'\tpackets["{packet_id}"] = common.PacketConstruction{{\n'
        output_content += f'\t\tID:         "{packet_id}",\n'
        output_content += f'\t\tName:       "{packet["name"]}",\n'
        output_content += f'\t\tFormat:     "{packet["format"]}",\n'
        output_content += f'\t\tFieldNames: []string{{{packet["field_names"]}}},\n'
        output_content += '\t}\n'
    
    # Add the function return statement
    output_content += "\n\treturn packets\n}"
    
    # Write the output file
    with open(output_file, 'w') as f:
        f.write(output_content)
    
    print(f"Conversion complete. Output written to {output_file}")

if __name__ == "__main__":
    main()
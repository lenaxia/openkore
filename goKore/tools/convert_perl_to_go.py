#!/usr/bin/env python3
"""
Perl to Go Packet Definition Converter

This script converts packet definitions from Perl format to Go format.

Purpose:
    The script extracts packet definitions from the original Perl files:
    - src/Network/Send/ServerType0.pm
    - src/Network/Receive/ServerType0.pm
    
    And converts them to Go format, creating separate output files:
    - goKore/network/send/servers/servertype0.send.go
    - goKore/network/send/servers/servertype0.receive.go
    
    The Perl format looks like:
    'XXXX' => ['packet_name', 'format_string', [qw(field1 field2 ...)]]
    
    The Go format looks like:
    packets["XXXX"] = common.PacketConstruction{
        ID:         "XXXX",
        Name:       "packet_name",
        Format:     "format_string",
        FieldNames: []string{"field1", "field2", ...},
    }

Usage:
    python3 convert_perl_to_go.py

Output:
    Creates two files with packet definitions in Go format:
    - servertype0.send.go - Send packets
    - servertype0.receive.go - Receive packets
"""

import re
import os
import sys

def extract_packets_from_send_file(file_path):
    """Extract packet definitions from the Send Perl file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Find the packets hash
    packets_match = re.search(r'my\s+%packets\s*=\s*\((.*?)\);', content, re.DOTALL)
    if not packets_match:
        print(f"Error: Could not find packets hash in {file_path}")
        return {}
    
    packets_content = packets_match.group(1)
    
    # Extract packet definitions
    # Format: 'XXXX' => ['packet_name', 'format_string', [qw(field1 field2 ...)]]
    pattern = r"'([0-9A-F]{4})'\s*=>\s*\['([^']+)'(?:,\s*'([^']*)')?(?:,\s*\[qw\((.*?)\)\])?\]"
    matches = re.findall(pattern, packets_content)
    
    packets = {}
    for match in matches:
        packet_id = match[0]
        packet_name = match[1]
        format_string = match[2] if match[2] else ""
        field_names = match[3].split() if match[3] else []
        
        packets[packet_id] = {
            'id': packet_id,
            'name': packet_name,
            'format': format_string,
            'field_names': field_names
        }
    
    return packets

def extract_packets_from_receive_file(file_path):
    """Extract packet definitions from the Receive Perl file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Find the packet_list hash
    packet_list_match = re.search(r'\$self->\{packet_list\}\s*=\s*\{(.*?)\};', content, re.DOTALL)
    if not packet_list_match:
        print(f"Error: Could not find packet_list hash in {file_path}")
        return {}
    
    packet_list_content = packet_list_match.group(1)
    
    # Extract packet definitions
    # Format: 'XXXX' => ['packet_name', 'format_string', [qw(field1 field2 ...)]]
    pattern = r"'([0-9A-F]{4})'\s*=>\s*\['([^']+)'(?:,\s*'([^']*)')?(?:,\s*\[qw\((.*?)\)\])?\]"
    matches = re.findall(pattern, packet_list_content)
    
    packets = {}
    for match in matches:
        packet_id = match[0]
        packet_name = match[1]
        format_string = match[2] if match[2] else ""
        field_names = match[3].split() if match[3] else []
        
        packets[packet_id] = {
            'id': packet_id,
            'name': packet_name,
            'format': format_string,
            'field_names': field_names
        }
    
    return packets

def convert_to_go_format(packets, output_file):
    """Convert packet definitions to Go format and write to file."""
    # Generate the Go file content
    content = """// Package servers provides server-specific packet constructions for different server types.
package servers

import (
	"github.com/lenaxia/goKore/network/common"
)

// ServerType0PacketConstructions provides packet constructions for ServerType0
func ServerType0PacketConstructions() map[string]common.PacketConstruction {
	packets := map[string]common.PacketConstruction{}

"""
    
    # Add all packet definitions in the Go format
    for packet_id in sorted(packets.keys()):
        packet = packets[packet_id]
        content += f'\tpackets["{packet_id}"] = common.PacketConstruction{{\n'
        content += f'\t\tID:         "{packet_id}",\n'
        content += f'\t\tName:       "{packet["name"]}",\n'
        content += f'\t\tFormat:     "{packet["format"]}",\n'
        
        # Handle field names
        if packet["field_names"]:
            field_names_str = '", "'.join(packet["field_names"])
            content += f'\t\tFieldNames: []string{{"{field_names_str}"}},\n'
        else:
            content += '\t\tFieldNames: []string{},\n'
        
        content += '\t}\n'
    
    # Add the function return statement
    content += "\n\treturn packets\n}"
    
    # Write the output file
    with open(output_file, 'w') as f:
        f.write(content)
    
    print(f"Conversion complete. Output written to {output_file}")

def main():
    # File paths
    send_file = "src/Network/Send/ServerType0.pm"
    receive_file = "src/Network/Receive/ServerType0.pm"
    send_output_file = "goKore/network/send/servers/servertype0.send.go"
    receive_output_file = "goKore/network/receive/servers/servertype0.receive.go"
    
    # Extract packets from both files
    print("Extracting packets from Send file...")
    send_packets = extract_packets_from_send_file(send_file)
    print(f"Found {len(send_packets)} packet definitions in the Send file.")
    
    print("Extracting packets from Receive file...")
    receive_packets = extract_packets_from_receive_file(receive_file)
    print(f"Found {len(receive_packets)} packet definitions in the Receive file.")
    
    # Convert to Go format and write to files
    print("Converting Send packets to Go format...")
    convert_to_go_format(send_packets, send_output_file)
    
    print("Converting Receive packets to Go format...")
    convert_to_go_format(receive_packets, receive_output_file)
    
    print("Done!")

if __name__ == "__main__":
    main()
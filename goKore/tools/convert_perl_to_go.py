#!/usr/bin/env python3
"""
Perl to Go Packet Definition Converter

This script converts packet definitions from Perl format to Go format.

Purpose:
    The script extracts packet definitions from the original Perl files:
    - src/Network/Send/ServerTypeX.pm
    - src/Network/Receive/ServerTypeX.pm
    
    And converts them to Go format, creating separate output files:
    - goKore/network/send/servers/servertypeX.go
    - goKore/network/receive/servers/servertypeX.go
    
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
    python3 convert_perl_to_go.py [--server-type SERVERTYPE]

Arguments:
    --server-type SERVERTYPE  Server type to convert (default: ServerType0)
                             This will look for src/Network/Send/SERVERTYPE.pm and
                             src/Network/Receive/SERVERTYPE.pm

Output:
    Creates two files with packet definitions in Go format:
    - goKore/network/send/servers/SERVERTYPE_LOWERCASE.go - Send packets
    - goKore/network/receive/servers/SERVERTYPE_LOWERCASE.go - Receive packets
"""

import re
import os
import sys
import argparse

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
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Convert Perl packet definitions to Go format')
    parser.add_argument('--server-type', default="ServerType0",
                        help='Server type to convert (e.g., ServerType0, ServerType1)')
    args = parser.parse_args()
    
    # Extract server type and create lowercase version for file names
    server_type = args.server_type
    server_type_lower = server_type.lower()
    
    # File paths
    send_file = f"src/Network/Send/{server_type}.pm"
    receive_file = f"src/Network/Receive/{server_type}.pm"
    send_output_file = f"goKore/network/send/servers/{server_type_lower}.go"
    receive_output_file = f"goKore/network/receive/servers/{server_type_lower}.go"
    
    # Check if input files exist
    if not os.path.exists(send_file):
        print(f"Error: Send file not found: {send_file}")
        return
    
    if not os.path.exists(receive_file):
        print(f"Error: Receive file not found: {receive_file}")
        return
    
    # Extract packets from both files
    print(f"Extracting packets from Send file: {send_file}")
    send_packets = extract_packets_from_send_file(send_file)
    print(f"Found {len(send_packets)} packet definitions in the Send file.")
    
    print(f"Extracting packets from Receive file: {receive_file}")
    receive_packets = extract_packets_from_receive_file(receive_file)
    print(f"Found {len(receive_packets)} packet definitions in the Receive file.")
    
    # Create output directories if they don't exist
    os.makedirs(os.path.dirname(send_output_file), exist_ok=True)
    os.makedirs(os.path.dirname(receive_output_file), exist_ok=True)
    
    # Convert to Go format and write to files
    print(f"Converting Send packets to Go format: {send_output_file}")
    convert_to_go_format(send_packets, send_output_file)
    
    print(f"Converting Receive packets to Go format: {receive_output_file}")
    convert_to_go_format(receive_packets, receive_output_file)
    
    print("Done!")

if __name__ == "__main__":
    main()
#!/usr/bin/env python3

import re
import os

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

if __name__ == "__main__":
    packets_file_path = "network/send/servers/servertype0.packets"
    go_file_path = "network/send/servers/servertype0.go"
    
    fix_packet_markings(packets_file_path, go_file_path)
    
    # Special case for C350 packet
    with open(packets_file_path, 'r') as f:
        lines = f.readlines()
    
    for i, line in enumerate(lines):
        if "# 'C350'" in line:
            lines[i] = line.replace("[ ]", "[X]")
            print(f"Fixed C350 packet marking")
            break
    
    with open(packets_file_path, 'w') as f:
        f.writelines(lines)
    
    print(f"Packet markings fixed in {packets_file_path}")
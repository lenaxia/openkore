#!/usr/bin/env python3

import re

def fix_packet_markings():
    packets_file_path = "../network/send/servers/servertype0.packets"
    
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
    
    print(f"Fixed {fixed_count} packet markings")

if __name__ == "__main__":
    fix_packet_markings()
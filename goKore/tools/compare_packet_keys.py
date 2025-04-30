#!/usr/bin/env python3
"""
Packet Keys Comparison Tool

This script compares packet keys between two Go files to ensure all packets were properly converted.

Purpose:
    After converting packet definitions from one format to another using convert_packet_format.py,
    this script verifies that all packet keys from the original file are present in the new file.
    It helps identify any missing or extra packet keys that might have occurred during conversion.

Usage:
    python3 compare_packet_keys.py [--original ORIGINAL_FILE] [--new NEW_FILE]

Arguments:
    --original ORIGINAL_FILE  Path to the original Go file (default: ../network/send/servers/servertype0.go)
    --new NEW_FILE            Path to the new Go file (default: ../network/send/servers/servertype0.5.go)

Output:
    Prints a comparison report showing:
    - Number of unique packet keys in the original file
    - Number of unique packet keys in the new file
    - Any missing packet keys (in original but not in new)
    - Any extra packet keys (in new but not in original)
"""

import re
import os
import argparse

def extract_keys_from_original_file(file_path):
    """Extract packet keys from the original servertype0.go file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Remove commented-out lines
    content = re.sub(r'//.*$', '', content, flags=re.MULTILINE)
    
    # Extract keys from first format: "0064": { ... }
    first_format_pattern = r'"([0-9A-F]{4})"\s*:'
    first_format_keys = set(re.findall(first_format_pattern, content))
    
    # Extract keys from second format: packets["006B"] = common.PacketConstruction{ ... }
    second_format_pattern = r'packets\["([0-9A-F]{4})"\]\s*='
    second_format_keys = set(re.findall(second_format_pattern, content))
    
    # Combine both sets of keys
    all_keys = first_format_keys.union(second_format_keys)
    
    return all_keys

def extract_keys_from_new_file(file_path):
    """Extract packet keys from the new servertype0.5.go file."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Remove commented-out lines
    content = re.sub(r'//.*$', '', content, flags=re.MULTILINE)
    
    # Extract keys from the new format: packets["0064"] = common.PacketConstruction{ ... }
    pattern = r'packets\["([0-9A-F]{4})"\]\s*='
    keys = set(re.findall(pattern, content))
    
    return keys

def compare_keys(original_keys, new_keys):
    """Compare keys from both files and report differences."""
    # Check for missing keys (in original but not in new)
    missing_keys = original_keys - new_keys
    
    # Check for extra keys (in new but not in original)
    extra_keys = new_keys - original_keys
    
    return missing_keys, extra_keys

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='Compare packet keys between two Go files')
    parser.add_argument('--original', default="../network/send/servers/servertype0.go",
                        help='Path to the original Go file')
    parser.add_argument('--new', default="../network/send/servers/servertype0.5.go",
                        help='Path to the new Go file')
    args = parser.parse_args()
    
    # File paths
    original_file = args.original
    new_file = args.new
    
    # Extract keys from both files
    print(f"Extracting keys from original file: {original_file}")
    original_keys = extract_keys_from_original_file(original_file)
    print(f"Found {len(original_keys)} unique packet keys in the original file.")
    
    print(f"Extracting keys from new file: {new_file}")
    new_keys = extract_keys_from_new_file(new_file)
    print(f"Found {len(new_keys)} unique packet keys in the new file.")
    
    # Compare keys
    missing_keys, extra_keys = compare_keys(original_keys, new_keys)
    
    # Report results
    print("\nComparison Results:")
    print("-------------------")
    
    if not missing_keys and not extra_keys:
        print("All keys from the original file are present in the new file.")
        print("No extra keys were found in the new file.")
        print("The conversion was successful!")
    else:
        if missing_keys:
            print(f"Missing keys (in original but not in new): {sorted(list(missing_keys))}")
        if extra_keys:
            print(f"Extra keys (in new but not in original): {sorted(list(extra_keys))}")
    
    # Print summary
    print("\nSummary:")
    print(f"Original file: {len(original_keys)} unique packet keys")
    print(f"New file: {len(new_keys)} unique packet keys")
    print(f"Missing keys: {len(missing_keys)}")
    print(f"Extra keys: {len(extra_keys)}")

if __name__ == "__main__":
    main()

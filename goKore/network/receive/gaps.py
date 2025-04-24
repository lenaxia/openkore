#!/usr/bin/env python3
import re
import sys

def parse_line_numbers(line):
    """Extract all line number ranges from a line"""
    matches = re.findall(r'lines (\d+(?:-\d+)?(?:,\d+(?:-\d+)?)*)', line)
    if not matches:
        return []
    
    ranges = []
    for group in matches[0].split(','):
        if '-' in group:
            start, end = map(int, group.split('-'))
            ranges.append((start, end))
        elif group:
            num = int(group)
            ranges.append((num, num))
    return line.strip(), ranges

def process_file(filename, start_line=1, end_line=12435, min_gap_size=1):
    """Process the file and identify gaps and overlaps"""
    with open(filename, 'r') as f:
        lines = f.readlines()
    
    # Extract all covered ranges with their content
    all_entries = []
    for line in lines:
        content, ranges = parse_line_numbers(line)
        for r in ranges:
            all_entries.append((r[0], r[1], content))
    
    # Sort entries by start line
    all_entries.sort(key=lambda x: x[0])
    
    # Find gaps with surrounding context
    gaps = []
    prev_end = start_line - 1
    prev_content = None
    
    for current_start, current_end, content in all_entries:
        if current_start > prev_end + 1:
            gap_size = current_start - prev_end - 1
            if gap_size >= min_gap_size:
                # Get the content after the gap
                next_content = content
                gaps.append({
                    'before': prev_content,
                    'after': next_content,
                    'gap_start': prev_end + 1,
                    'gap_end': current_start - 1,
                    'size': gap_size
                })
        prev_end = max(prev_end, current_end)
        prev_content = content
    
    # Check for gap after last range
    if prev_end < end_line:
        gap_size = end_line - prev_end
        if gap_size >= min_gap_size:
            gaps.append({
                'before': prev_content,
                'after': None,
                'gap_start': prev_end + 1,
                'gap_end': end_line,
                'size': gap_size
            })
    
    # Find overlaps with context
    overlaps = []
    for i in range(1, len(all_entries)):
        prev_start, prev_end, prev_content = all_entries[i-1]
        current_start, current_end, current_content = all_entries[i]
        
        if current_start <= prev_end:
            overlap_start = max(prev_start, current_start)
            overlap_end = min(prev_end, current_end)
            overlaps.append({
                'first': prev_content,
                'second': current_content,
                'overlap_start': overlap_start,
                'overlap_end': overlap_end,
                'size': overlap_end - overlap_start + 1
            })
    
    return gaps, overlaps

def print_results(gaps, overlaps):
    """Print the results with context"""
    print("\nGaps (missing coverage):")
    for gap in gaps:
        print(f"[ ] Gap: {gap['gap_start']}-{gap['gap_end']} ({gap['size']} lines)")
        print(f"    Before: {gap['before'] or 'N/A'}")
        print(f"    After:  {gap['after'] or 'N/A'}")
    
    print("\nOverlaps (duplicate coverage):")
    for overlap in overlaps:
        print(f"[ ] Overlap: {overlap['overlap_start']}-{overlap['overlap_end']} ({overlap['size']} lines)")
        print(f"    First coverage:  {overlap['first']}")
        print(f"    Second coverage: {overlap['second']}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <filename> [start_line] [end_line] [min_gap_size]")
        sys.exit(1)
    
    filename = sys.argv[1]
    start_line = int(sys.argv[2]) if len(sys.argv) > 2 else 1
    end_line = int(sys.argv[3]) if len(sys.argv) > 3 else 12435
    min_gap_size = int(sys.argv[4]) if len(sys.argv) > 4 else 1
    
    gaps, overlaps = process_file(filename, start_line, end_line, min_gap_size)
    print_results(gaps, overlaps)

#!/usr/bin/env python3
"""
Sequence Diagram Generator for Ragnarok Online Login Sequence
This script generates a PlantUML sequence diagram showing the packet exchange
between the client and the different servers during the login process.
"""

import os
import re
import sys
from typing import Dict, List, Tuple

def extract_login_sequence(dump_file: str) -> List[Tuple[str, str, str, str]]:
    """Extract the login sequence from a dump file."""
    with open(dump_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    sequence = []
    current_server = "Unknown"
    
    # Track server connections
    if "Connecting to Account Server" in content:
        current_server = "Account Server"
    
    # Extract all packet exchanges
    lines = content.split('\n')
    for i, line in enumerate(lines):
        # Check for server transitions
        if "Connecting to Character Server" in line:
            current_server = "Character Server"
        elif "Connecting to Map Server" in line:
            current_server = "Map Server"
        
        # Extract sent packets
        sent_match = re.search(r'>> Sent packet: (\w+)\s+\[(.*?)\] \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})', line)
        if sent_match:
            packet_id, description, size, timestamp = sent_match.groups()
            sequence.append(("Client", current_server, packet_id, description))
        
        # Extract received packets
        recv_match = re.search(r'<< Received packet:\s+(\w+) - (.*?) \[(\d+) bytes\]\s+(\d{4}\.\d{2}\.\d{2} \d{2}:\d{2}:\d{2})', line)
        if recv_match:
            packet_id, description, size, timestamp = recv_match.groups()
            sequence.append((current_server, "Client", packet_id, description))
    
    return sequence

def generate_sequence_diagram(sequences: List[List[Tuple[str, str, str, str]]]) -> str:
    """Generate a PlantUML sequence diagram from the login sequences."""
    # Combine all sequences to find all unique packets
    all_packets = set()
    for sequence in sequences:
        for _, _, packet_id, _ in sequence:
            all_packets.add(packet_id)
    
    # Create a mapping of packet IDs to colors for consistent coloring
    colors = ["#E6F7FF", "#E6FFE6", "#FFE6E6", "#F7E6FF", "#FFE6F7", "#E6FFF7", "#FFF7E6", "#F7FFE6"]
    packet_colors = {}
    for i, packet_id in enumerate(sorted(all_packets)):
        packet_colors[packet_id] = colors[i % len(colors)]
    
    # Start the PlantUML diagram
    diagram = []
    diagram.append("@startuml")
    diagram.append("skinparam sequenceMessageAlign center")
    diagram.append("skinparam sequenceArrowThickness 2")
    diagram.append("skinparam responseMessageBelowArrow true")
    diagram.append("skinparam maxMessageSize 200")
    diagram.append("skinparam wrapWidth 200")
    
    diagram.append("\ntitle Ragnarok Online Login Sequence\n")
    
    diagram.append("participant \"Client\" as Client")
    diagram.append("participant \"Account Server\" as AccountServer")
    diagram.append("participant \"Character Server\" as CharServer")
    diagram.append("participant \"Map Server\" as MapServer")
    
    diagram.append("\n== Account Server Login ==\n")
    
    # Use the first sequence as the reference
    reference_sequence = sequences[0]
    current_server = "Unknown"
    
    for from_entity, to_entity, packet_id, description in reference_sequence:
        # Check for server transitions
        if current_server != to_entity and to_entity != "Client" and current_server != from_entity and from_entity != "Client":
            current_server = to_entity if to_entity != "Client" else from_entity
            section_name = f"== {current_server} Login =="
            if section_name not in diagram:
                diagram.append(f"\n{section_name}\n")
        
        # Format the entities for PlantUML
        from_entity_id = from_entity.replace(" ", "")
        to_entity_id = to_entity.replace(" ", "")
        
        # Format the message
        message = f"{packet_id}: {description}"
        color = packet_colors.get(packet_id, "#FFFFFF")
        
        # Add the message to the diagram
        diagram.append(f"{from_entity_id} -[{color}]> {to_entity_id}: {message}")
    
    diagram.append("\n@enduml")
    return "\n".join(diagram)

def main():
    """Main function to generate the sequence diagram."""
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    dump_dir = os.path.join(base_dir, "dumps")
    dump_files = [os.path.join(dump_dir, f) for f in os.listdir(dump_dir) if f.startswith("DUMP")]
    
    if not dump_files:
        print("No dump files found!")
        return
    
    print(f"Found {len(dump_files)} dump files: {', '.join(os.path.basename(f) for f in dump_files)}")
    
    # Extract login sequences from all dump files
    sequences = []
    for file_path in dump_files:
        print(f"Extracting login sequence from {os.path.basename(file_path)}...")
        sequences.append(extract_login_sequence(file_path))
    
    # Generate the sequence diagram
    diagram = generate_sequence_diagram(sequences)
    
    # Write the results
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    output_dir = os.path.join(base_dir, "analysis")
    with open(os.path.join(output_dir, "login_sequence_diagram.puml"), "w") as f:
        f.write(diagram)
    
    print("Sequence diagram generated! Saved to analysis/login_sequence_diagram.puml")
    print("You can visualize this diagram using PlantUML or online tools like https://www.planttext.com/")

if __name__ == "__main__":
    main()
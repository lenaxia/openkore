# Ragnarok Online Login Sequence Analysis

This directory contains tools and analysis for understanding the packet exchange sequence required for a successful login to a Ragnarok Online server. The analysis is based on examining multiple successful login attempts captured in the `../dumps/` directory.

## Folder Structure

The project is organized into the following directories:

- `../dumps/` - Contains the raw packet dump files (DUMP1, DUMP2, etc.)
- `../extracteddata/` - Contains JSON files extracted from the dump files
- `../analysis/` - Contains markdown files and other analysis results
- `./` (this directory) - Contains the analysis scripts

## Overview

The login process for Ragnarok Online involves communication with three different servers:

1. **Account Server** - Handles authentication and provides server list
2. **Character Server** - Manages character selection and provides character data
3. **Map Server** - Handles the actual game world and gameplay

Each server requires a specific sequence of packets to be exchanged for successful login. This analysis aims to document and understand that sequence.

## Tools

### 1. Packet Analyzer (`packet_analyzer.py`)

This script analyzes the dump files to extract the common packet sequence across all login attempts. It identifies the essential packets required for a successful login.

**Usage:**
```bash
python packet_analyzer.py
```

**Output:**
- `../analysis/login_sequence_summary.md` - A markdown file summarizing the login sequence for each server

### 2. Packet Structure Analyzer (`packet_structure_analyzer.py`)

This script analyzes the structure of key packets in the login sequence by comparing multiple instances of the same packet across different login attempts. It helps understand what fields are constant and what fields vary between login attempts.

**Usage:**
```bash
python packet_structure_analyzer.py
```

**Output:**
- `../analysis/packet_structures.md` - A markdown file detailing the structure of key packets

### 3. Sequence Diagram Generator (`sequence_diagram_generator.py`)

This script generates a PlantUML sequence diagram showing the packet exchange between the client and the different servers during the login process. This provides a visual representation of the login flow.

**Usage:**
```bash
python sequence_diagram_generator.py
```

**Output:**
- `../analysis/login_sequence_diagram.puml` - A PlantUML file that can be rendered into a sequence diagram

To visualize the sequence diagram, you can use:
- PlantUML (https://plantuml.com/)
- Online PlantUML editor (https://www.planttext.com/)

### 4. Dump to JSON Converter (`dump_to_json.py` and `dump_to_json_improved.py`)

These scripts convert a packet dump file to a structured JSON format containing all packet information. The output can be used as test data for integration tests to validate that application behavior is deterministic and repeatable.

**Usage:**
```bash
python dump_to_json.py <dump_file> -o <output_file> [--pretty]
python dump_to_json_improved.py <dump_file> -o <output_file> [--pretty]
```

**Example:**
```bash
python dump_to_json_improved.py ../dumps/DUMP1 --pretty -o ../extracteddata/dump1_packets.json
```

**Output:**
- JSON file with all packet information

### 5. Essential Login Sequence Extractor (`extract_login_sequence.py`)

This script extracts only the essential packets needed for a successful login sequence from a packet dump file. It helps to quickly analyze the critical components of the login process without the noise of other packets.

**Usage:**
```bash
python extract_login_sequence.py <dump_file> -o <output_file>
```

**Example:**
```bash
python extract_login_sequence.py ../dumps/DUMP1 -o ../analysis/essential_login_packets.md
```

**Output:**
- Markdown file with essential login packets

### 6. Dump Comparison Tool (`compare_dumps.py`)

This script compares the packet sequences across multiple JSON dump files to identify any differences or inconsistencies in the login process.

**Usage:**
```bash
python compare_dumps.py <json_file1> <json_file2> [<json_file3> ...] -o <output_file>
```

**Example:**
```bash
python compare_dumps.py ../extracteddata/dump1_packets.json ../extracteddata/dump2_packets.json -o ../analysis/packet_comparison_report.md
```

**Output:**
- Markdown file with comparison results

### 7. Run All Analysis (`run_analysis.py`)

This script runs all the analysis tools and generates a comprehensive report.

**Usage:**
```bash
python run_analysis.py
```

**Output:**
- `../analysis/login_sequence_summary.md`
- `../analysis/packet_structures.md`
- `../analysis/login_sequence_diagram.puml`
- `../analysis/comprehensive_login_analysis.md` (combined report)

## Process All Dumps

To process all dumps and generate JSON files:

```bash
cd /home/mikekao/personal/openkore
for i in {1..7}; do 
  python3 goKore/verification/PacketAnalysis/LoginSequenceAnalysis/dump_to_json_improved.py \
  goKore/verification/PacketAnalysis/dumps/DUMP$i \
  --pretty -o goKore/verification/PacketAnalysis/extracteddata/dump${i}_packets.json
done
```

## Key Findings

The login sequence can be summarized as follows:

1. **Account Server Login**
   - Client sends account credentials (username/password)
   - Server responds with account info and server list
   - Client disconnects from Account Server

2. **Character Server Login**
   - Client connects to Character Server with session info from Account Server
   - Server sends character list
   - Client selects a character
   - Server sends character and map info
   - Client disconnects from Character Server

3. **Map Server Login**
   - Client connects to Map Server with session info from Character Server
   - Client sends map login packet
   - Server confirms login and sends initial map data
   - Client acknowledges map load
   - Login process complete

## Implementation Notes

When implementing a client that needs to authenticate with a Ragnarok Online server:

1. The session IDs received from the Account Server must be preserved and passed to the Character Server
2. The session IDs and character ID received from the Character Server must be preserved and passed to the Map Server
3. The client must properly acknowledge packet receipt where required (e.g., map loaded)
4. PIN code verification may be required depending on server configuration

## Further Analysis

For a deeper understanding of specific packets, refer to the generated `../analysis/packet_structures.md` file, which details the structure of key packets in the login sequence.
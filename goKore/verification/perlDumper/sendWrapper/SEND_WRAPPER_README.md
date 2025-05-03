# Send.pm Wrapper for Validation Data Generation

This tool allows you to call individual Send.pm subroutines to generate packet output for validation purposes. It's designed to help with generating sample data for reimplementation of the OpenKore network protocol.

## Overview

The `send_wrapper.pl` script creates a mock environment that:

1. Loads the original Send.pm module
2. Mocks all necessary dependencies
3. Creates a mock connection that captures packets instead of sending them
4. Allows calling any Send.pm method and captures the output
5. Outputs the generated packets in both human-readable and JSON formats

## Usage

```bash
# Make the script executable
chmod +x send_wrapper.pl

# List all available Send.pm methods
./send_wrapper.pl list

# Call a specific method with arguments
./send_wrapper.pl call sendMove 100 100

# Get help
./send_wrapper.pl help
```

## Output Format

When calling a method, the script will:

1. Print the generated packets in hex format to the console
2. Show the message ID and raw bytes for each packet
3. Save a JSON file with the complete packet information

The JSON output file is named `packet_output_<method>.json` and contains:
- The method name
- The arguments used
- An array of packets, each with:
  - Hex representation
  - Message ID
  - Array of byte values

## Example

```bash
# Call the sendMove method
./send_wrapper.pl call sendMove 100 100
```

Output:
```
Calling sendMove with args: 100, 100
Generated packets:
Packet 0: 85640000640000a9860100
  Message ID: 0085
  Raw bytes: 85 64 00 00 64 00 00 A9 86 01 00 
JSON output saved to packet_output_sendMove.json
```

## Extending

If you need to add support for additional methods or server types:

1. Modify the mock objects in the script to provide any required dependencies
2. Add any necessary packet definitions to the Send object
3. For complex methods that require specific state, you may need to set up additional mock data

## Troubleshooting

If you encounter errors:

1. Check that the paths to the OpenKore source files are correct
2. Ensure all required dependencies are mocked properly
3. For methods that require specific state or configuration, you may need to add additional mock data

## Integration with Go Reimplementation

The JSON output files can be used as test fixtures for your Go reimplementation. You can:

1. Create test cases that compare your Go-generated packets with the original Perl-generated packets
2. Use the byte arrays for direct comparison
3. Build a validation suite that ensures compatibility between implementations
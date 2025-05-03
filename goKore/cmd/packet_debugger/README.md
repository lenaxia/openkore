# goKore Network Packet Debugger

This tool allows you to step through packet exchanges in the goKore network stack. It helps validate the core functionality of send and receive stacks by visually comparing packet construction and decoding.

## Features

- Load and parse packet dumps from JSON files
- Step through packets one by one
- Display sent and received packets in a 4-column layout:
  - Column 1: Client representation of received packets (decoded values vs expected)
  - Column 2: Send packets from client to server
  - Column 3: Receive packets from server to client
  - Column 4: Server representation of sent packets (encoded vs expected)
- Compare encoded/decoded packets with the expected data from dumps

## Usage

1. Build the tool:
   ```
   cd /path/to/goKore/network/tools
   go build packet_debugger.go
   ```

2. Run the tool:
   ```
   ./packet_debugger
   ```

3. Select a packet dump file from the list by entering its number.

4. Step through packets by pressing Enter. Type 'q' to quit.

## Packet Dump Format

The tool expects packet dumps in JSON format with the following structure:

```json
{
  "file_name": "DUMP_NAME",
  "packet_count": 123,
  "packets": [
    {
      "direction": "sent",
      "packet_id": "0064",
      "description": "Account Server Login",
      "size": 55,
      "timestamp": "2025.04.30 16:24:58",
      "raw_data": [
        {
          "hex_bytes": ["64", "00", "1C", ...],
          "ascii_representation": "",
          "binary_base64": "..."
        }
      ],
      "parsed_data": null,
      "server_type": "Account Server"
    },
    // More packets...
  ]
}
```

## How It Works

1. The tool loads packet dumps from the specified directory.
2. It instantiates both Send and Receive factories to process packets.
3. For sent packets:
   - Column 2 shows the raw packet data from the dump
   - Column 4 shows how the Send factory would encode the packet
4. For received packets:
   - Column 3 shows the raw packet data from the dump
   - Column 1 shows how the Receive factory would decode the packet
5. This allows you to visually verify that the packet encoding and decoding is working correctly.

## Extending the Tool

To fully implement packet encoding and decoding:

1. Enhance the `encodeSendPacket` function to:
   - Extract field values from the packet
   - Use the Send instance to encode the packet
   - Compare the encoded packet with the raw data from the dump

2. Enhance the `decodeReceivePacket` function to:
   - Use the Receive instance to decode the packet
   - Extract and display the field values
   - Compare with expected values if available

## Troubleshooting

If you encounter issues:

- Make sure the packet dump files are in the correct format
- Check that the server type matches what's registered in the factories
- Verify that the packet IDs in the dumps match those in the server type definitions
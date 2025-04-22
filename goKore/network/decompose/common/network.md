**Server Connection Handlers:**

- parse() - Packet parser wrapper (lines 610-626)
  - Wraps SUPER::parse method
  - Adds debug logging when configured (debugPacket_received=3)
  - Logs packet variables when:
    - debugPacket_received=3
    - Packet switch is in debugPacket_include list
  - Returns parsed packet arguments

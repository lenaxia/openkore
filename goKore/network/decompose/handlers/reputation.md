**Reputation System Handlers:**

- repute_info - Reputation information handler (lines 12387-12405)
  - Processes reputation information notifications (0B8D)
  - Clears existing reputation_list global array
  - Defines unpacking structure:
    * 16-byte entries with V4 unpacking
    * Keys: type, type2, points, points2
  - Parses reputation data from binary format:
    * Processes data in 16-byte chunks
    * Unpacks each chunk into a reputation entry
    * Adds each entry to reputation_list array
  - No direct user feedback or messages
  - Packet: PACKET_ZC_REPUTE_INFO
# Boss Monster Related Handlers

**Method Implementations:**
- boss_map_info - Boss monster information handler (lines 9539-9555)
  - Processes boss monster information notifications
  - Gets boss name from packet (converted from bytes)
  - Handles multiple flag values:
    * 0: No boss found - "You cannot find any trace of a Boss Monster in this area"
    * 1: Boss location - "MVP Boss X is now on location: (x, y)"
    * 2: Boss detected - "MVP Boss X has been detected on this map!"
    * 3: Boss respawn - "MVP Boss X is dead, but will spawn again in X hour(s) and X minutes(s)"
    * Other: Unknown result (outputs debug and warning messages)
  - Uses "info" message category for all messages
  - Packet: 0293
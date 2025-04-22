**Status Change Handlers:**

- unit_levelup() - Handles level up and effect notifications (lines 2819-2844)
  - Processes different effect types:
    * Base level up (triggers base_level hook)
    * Job level up (triggers job_level hook)
    * Refine success/failure
    * Pharmacy success/failure
    * Game over
  - Displays appropriate messages for each effect type
  - Uses Actor::get() to retrieve actor information

- stats_added() - Handles status change acknowledgments (lines 1676-1738)
  - Processes ZC_STATUS_CHANGE_ACK packets
  - Handles success/failure cases (207 = failure)
  - Updates character stats including:
    - Basic stats (STR, AGI, VIT, INT, DEX, LUK)
    - Special stats (POW, STA, WIS, SPL, CON, CRT)
  - Calls packet_charStats hook after processing
  - Debug logs each stat update

- stats_info() - Character status updates (lines 1743-1790)
  - Processes ZC_STATUS packets
  - Requires in-game state (changeToInGameState)
  - Updates comprehensive character stats including:
    - Base stats (STR, AGI, VIT, INT, DEX, LUK) with points
    - Attack values (physical and magic)
    - Defense values (physical and magic)
    - Combat stats (hit, flee, critical)
    - Free stat points
  - Detailed debug output showing all updated stats
- stat_info2() - Extended status information (lines 1794-1825)
  - Processes ZC_COUPLESTATUS packets
  - Requires in-game state (changeToInGameState)
  - Updates base and bonus values for:
    - STR, AGI, VIT, INT, DEX, LUK
  - Triggers inventory callback (onStatInfo2) for certain server types
  - Debug logs each updated stat with base+bonus values
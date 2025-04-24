**Stat Handlers:**

- stats_added - Status change request result handler (lines 1670-1738)
  - Handles the result of a status change request (ZC_STATUS_CHANGE_ACK)
  - Processes success/failure of stat point allocation
  - Updates character stats based on the type:
    - Base stats (STR, AGI, VIT, INT, DEX, LUK)
    - Special stats (POW, STA, WIS, SPL, CON, CRT)
  - Triggers packet_charStats plugin hook

- stats_info - Character status information handler (lines 1740-1790)
  - Handles complete character status information (ZC_STATUS)
  - Updates character's base stats and required points
  - Updates derived stats (attack, defense, hit, flee, etc.)
  - Provides detailed debug output of all stats

- stat_info2 - Character parameter change handler (lines 1792-1825)
  - Handles character parameter changes (ZC_COUPLESTATUS)
  - Updates both base and bonus values for stats
  - Triggers inventory update when stats change

- stat_info - Character parameter change handler (lines 1623-1668)
  - Handles various stat updates:
    - Base stats (STR, AGI, VIT, INT, DEX, LUK)
    - Combat stats (attack, defense, hit, flee, etc)
    - Experience and level changes
    - Weight and zeny
    - Special stats for homunculus/mercenary
  - Uses stat_info_handlers hash for specific stat processing
  - Handles different packet types:
    - 00B0: Character parameter change
    - 00B1: Long parameter change  
    - 00BE: Status change
    - 0141: Couple status
    - 01AB: Other player parameter change
    - 02A2: Mercenary parameter change
    - 07DB: Homunculus parameter change
    - 0ACB: Long parameter change v2
  - Maintains default walk speed if not set
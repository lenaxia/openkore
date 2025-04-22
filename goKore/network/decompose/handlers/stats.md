**Stat Handlers:**

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
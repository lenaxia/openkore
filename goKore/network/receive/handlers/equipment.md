# Equipment Related Handlers

**Method Implementations:**
- unequip_item_switch - Equipment switch unequip handler (lines 11643-11670)
  - Processes equipment switch unequip notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from inventory by ID
  - Clears eqswitch flag from item
  - Handles special equipment types:
    * For arrows (type 10 or 32768):
      - Clears character's eqswitch arrow field
    * For other equipment:
      - Iterates through equipSlot_rlut
      - Skips arrow slots (10 and 32768)
      - Clears equipment from character's eqswitch hash
      - Triggers unequipped_item_sw hook with slot and item
  - Displays message with item name, binID, and equipment type
  - Uses "inventory" message category
  - Prefixes message with "[Equip Switch]"
  - Contains comment about packet format (0A9A)
- unequip_item - Item unequip handler (lines 11612-11639)
  - Processes item unequip notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from inventory by ID
  - Clears equipped flag from item
  - Handles special equipment types:
    * For arrows (type 10 or 32768):
      - Clears character's arrow equipment
      - Clears character's arrow field
    * For other equipment:
      - Iterates through equipSlot_rlut
      - Skips arrow slots (10 and 32768)
      - Clears equipment from character's equipment hash
      - Triggers unequipped_item hook with slot and item
  - Displays message with item name, binID, and equipment type
  - Uses "inventory" message category
  - Contains comments about packet formats:
    * 00AC: Basic packet (ZC_REQ_TAKEOFF_EQUIP_ACK)
    * 08D1: Extended packet (ZC_REQ_TAKEOFF_EQUIP_ACK2)
    * 099A: Version 5 packet (ZC_ACK_TAKEOFF_EQUIP_V5)
  - Contains comment about result values (inversed for v2 v5):
    * 0 = failure
    * 1 = success
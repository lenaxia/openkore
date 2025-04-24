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
- arrow_none - Arrow status and weight limit handler (lines 10160-10182)
  - Processes notifications about arrow status and weight limit issues
  - Handles different type values:
    * Type 0: No arrows equipped
      - Deletes character's arrow field
      - Handles dcOnEmptyArrow configuration:
        * If enabled: Logs error, sends quit message, and exits
        * If disabled: Displays "Please equip arrow first" error
    * Type 1: Weight limit exceeded - "You can't Attack or use Skills"
    * Type 2: Weight limit exceeded - "You can't use Skills"
    * Type 3: Arrow equipped - "Arrow equipped"
  - Simple implementation focused on notification and auto-disconnect
  - No plugin hooks triggered
  - Packet: 013B
  - Format: 'v' (type)
- equip_item - Item equip handler (lines 10627-10651)
  - Processes item equip notifications
  - Gets item from character's inventory by ID
  - Handles success/failure based on packet type:
    * 00AA with success=0: Displays failure message
    * 0999 with success=1: Displays failure message
    * Other cases: Processes successful equip
  - For successful equips:
    * Sets item's equipped flag to the type value
    * Updates character's equipment references:
      - Special handling for arrows (type 10 or 32768)
      - For other equipment: Maps type to equipment slots
      - Skips arrow slots (workaround for Arrow bug)
    * Triggers equipped_item hook with slot and item
    * Displays success message with item details
  - Decrements AI waitForEquip counter if set
  - Packets:
    * 00AA: Format 'a2 v C' or 'a2 v2 C' (ID type [viewid] success)
    * 0999: Format 'a2 V v C' (ID type viewID success)
    * 08D0: Format 'a2 v2 C' (ID type viewid success)
- equip_item_switch - Equipment switch equip handler (lines 10653-10680)
  - Processes equipment switch equip notifications
  - Gets item from character's inventory by ID
  - Handles success/failure based on success value:
    * 1 or 2: Displays failure message
    * Other values: Processes successful equip
  - For successful equips:
    * Sets item's eqswitch flag to the type value
    * Updates character's equipment references:
      - Special handling for arrows (type 10 or 32768)
      - For other equipment: Maps type to equipment slots
      - Skips arrow slots (workaround for Arrow bug)
    * Triggers equipped_item_sw hook with slot and item
    * Displays success message with item name and type
  - Decrements AI waitForEquip counter if set
  - Prefixes all messages with "[Equip Switch]"
  - Packets:
    * 0A98: Format 'a2 V v' or 'a2 V2' (ID type success)
- equip_switch_log - Equipment switch list handler (lines 10694-10703)
  - Processes full list of items in the equip switch window
  - Parses variable-length log data in 6-byte chunks
  - For each chunk:
    * Unpacks item index and position
    * Gets item from character's inventory by ID
    * Updates character's eqswitch hash with item reference
  - No messages or notifications
  - No plugin hooks triggered
  - Simple implementation focused on data structure updates
  - Packet: 0A9B
  - Format: 'v a*' (len log)
- equip_switch_run_res - Equipment switch execution result handler (lines 10683-10692)
  - Processes acknowledgement packet for full equip switch operation
  - Handles success/failure based on success value:
    * 1: Displays failure message
    * 0: Displays success message
  - Simple implementation focused on notification
  - Uses "info" message category for all messages
  - Prefixes all messages with "[Equip Switch]"
  - No plugin hooks triggered
  - Packet: 0A9D
  - Format: 'v' (success)
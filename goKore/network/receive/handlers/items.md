# Items Related Handlers

- use_item - Item usage handler (lines 11673-11681)
  - Processes item usage notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from inventory by ID
  - Displays message with item name, binID, and amount
  - Calls inventoryItemRemoved to update inventory
  - Uses "useItem" message category
  - Contains TODO comment: "only used to report failure? $args->{success}"
  - Simple implementation focused on notification and inventory update

- item_preview - Item preview handler (lines 12182-12193)
  - Processes item preview notifications
  - Gets item reference from inventory using ID
  - Updates item properties:
    * broken status (if defined)
    * upgrade level
    * card information
    * option information
  - Updates item name using itemName function
  - Simple implementation focused on item data updates
- item_appeared - Ground item appearance handler (lines 7011-7052)
  - Processes item appearance notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from itemsList by ID
  - If item doesn't exist:
    * Creates new Actor::Item
    * Sets appear_time, amount, nameID, identified
    * Gets item name using itemName function
    * Sets ID and marks for addition to itemsList
  - Updates item position data
  - Adds item to itemsList if new
  - Implements auto-pickup logic:
    * Checks if AI is in AUTO state
    * Verifies item is in pickupitems list (priority 2)
    * Checks itemsTakeAuto or itemsGatherAuto configuration
    * Respects itemsGatherAuto_notInTown setting
    * Verifies character weight is below itemsMaxWeight
    * Ensures item is within 5 cells distance
    * Sends take packet if all conditions met
  - Displays item appearance message with details
  - Triggers item_appeared hook with item and type
  - Supports multiple packet formats:
    * 009E: Basic format (ZC_ITEM_FALL_ENTRY)
    * 084B: Extended format with type (ZC_ITEM_FALL_ENTRY4)
    * 0ADD: Extended format with drop effects (ZC_ITEM_FALL_ENTRY5)
- item_disappeared - Ground item disappearance handler (lines 7086-7119)
  - Processes item disappearance notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from itemsList by ID
  - Implements attackLooters functionality:
    * Checks if attackLooters is enabled and item is in pickupitems list
    * Iterates through monster list to find potential looters
    * Applies mon_control filters (attack_auto, attack_lvl, etc.)
    * Checks if monster is within attackLooters_dist of item
    * Initiates attack on the monster if conditions met
    * Displays "Attack Looter" message
  - Outputs debug message about item disappearance
  - Creates deep copy of item in items_old hash
  - Sets disappeared flag and gone_time in the copy
  - Removes item from itemsList
  - Packet: 00A1 (ZC_ITEM_DISAPPEAR)
  - Format: 'L' (ID)
- item_exists - Ground item existence handler (lines 7054-7084)
  - Processes item existence notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from itemsList by ID
  - If item doesn't exist:
    * Creates new Actor::Item
    * Sets appear_time, amount, nameID, identified
    * Gets item name using itemName function
    * Sets ID and marks for addition to itemsList
  - Updates item position data
  - Adds item to itemsList if new
  - Displays item existence message with details
  - Triggers item_exists hook with:
    * item: Item reference
    * type: Item type
    * show_effect: Show effect flag
    * effect_type: Effect type
  - Similar to item_appeared but without auto-pickup logic
- special_item_obtain - Special item acquisition handler (lines 9903-9953)
  - Processes special item acquisition notifications
  - Gets item name using itemNameSimple function
  - Gets holder name and strips language code
  - Handles different acquisition types:
    * TYPE_BOXITEM (box/container items):
      - Determines unpacking format based on etc value
      - Extracts box_nameID from etc field
      - Gets box item name
      - Formats message using msgTable[1629] or fallback format
      - Logs to "GM" chat log if logSystemChat is enabled
      - Uses "schat" message category
    * TYPE_MONSTER_ITEM (monster drops):
      - Extracts monster name from etc field
      - Strips language code from monster name
      - Formats message with holder, item, and monster names
      - Logs to "GM" chat log if logSystemChat is enabled
      - Uses "schat" message category
    * Unknown types:
      - Displays warning with type number
      - Uses "schat" warning category
  - Triggers packet_special_item_obtain hook with:
    * ObtainType, ItemName, ItemID, Holder
    * SourceItemID, SourceName, Msg
  - Comprehensive implementation for special item notifications

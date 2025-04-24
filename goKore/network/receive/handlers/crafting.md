# Crafting Related Handlers

**Method Implementations:**

- arrowcraft_list - Arrow crafting and poison creation list handler (lines 10183-10215)
  - Processes list of items that can be used for arrow crafting or poison creation
  - Parses raw message data to extract item IDs
  - Resets character's selected_craft flag
  - Clears existing arrowCraftID array
  - For each item ID in the packet:
    * Looks up item in character's inventory by nameID
    * Adds item's binID to arrowCraftID array
  - Special handling for GC_POISONINGWEAPON skill (ID 2027):
    * Displays "Received Possible Poison List" message
    * Handles autoPoison configuration:
      - Looks up configured poison item in inventory
      - Automatically sends selection if found
      - Sets selected_craft flag to 1
      - Displays error if configured poison not available
  - For other skills:
    * Displays "Received Possible Item List" message
    * Suggests using 'arrowcraft' or 'poison' commands
  - No plugin hooks triggered
  - Packet: 01AD
- makable_item_list - Craftable item list handler (lines 4999-5019)
  - Processes list of items that can be crafted
  - Resets makableList array
  - Gets unpack format from makable_item_list_pack or defaults to 'v4'
  - Calculates packed length for iteration
  - Creates formatted header for item list display
  - For each item in the list:
    * Unpacks nameID and material IDs (1-3)
    * Adds nameID to makableList array
    * Formats item entry with index, name, and ID
  - Displays formatted item list with header and footer
  - Informs user about 'create' command availability
  - Triggers makable_item_list hook with item_list
  - Uses "list" message category for item list
  - Uses "info" message category for command hint
  - Packet: 018d
  - Format: 'W { W { W }*3 }*' (len item_list)
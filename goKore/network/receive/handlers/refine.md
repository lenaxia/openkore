# Item Refining and Upgrading Handlers

**Method Implementations:**
- upgrade_message - Weapon upgrade result handler (lines 9179-9192)
  - Processes weapon upgrade result notifications
  - Gets item name using itemNameSimple function
  - Handles multiple result types:
    * 0: Success - "Weapon upgraded"
    * 1: Failure - "Weapon not upgraded"
    * 2: Skill level too low - "Cannot upgrade until you level up the upgrade weapon skill"
    * 3: Missing materials - "You lack item to upgrade the weapon"
  - Contains commented alternative message for type 1
  - Uses "info" message category for types 0, 1, 3
  - Uses "error" message category for type 2
  - Packet: 0223
- refine_result - Refine/craft result handler (lines 9161-9176)
  - Processes refining and crafting result notifications
  - Handles multiple result types:
    * 0: Successful weapon refine
    * 1: Failed weapon refine
    * 2: Successful potion creation
    * 3: Failed potion creation
    * 6: Successful cooking
    * Other: Unknown result (shows code)
  - Displays appropriate message for each result type
  - Shows item nameID in all messages
  - Simple implementation focused on result notification
- upgrade_list - Refine item list handler (lines 9114-9134)
  - Processes list of refinable items
  - Clears previous refine list
  - Creates formatted header with centered title
  - Iterates through raw message data in 13-byte chunks
  - For each item:
    * Extracts item index and nameID
    * Retrieves item from character inventory
    * Stores item ID in refineList array
    * Displays item index, name, and binID
  - Adds footer with separator line
  - Shows instruction for using 'refine' command
  - Uses "list" message category for item list
  - Uses "info" message category for instruction
  - Packet: 0221
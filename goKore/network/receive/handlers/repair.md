# Item Repair Related Handlers

**Method Implementations:**
- repair_result - Repair result handler (lines 11327-11339)
  - Processes repair result notifications
  - Adjusts index by subtracting 2
  - Gets item from character inventory by ID
  - Handles two flag values:
    * 1: Failure - "Repair of X failed"
    * 0: Success - "Successfully repaired 'X'"
  - Clears repairList after processing
  - Contains comments about packet format (01FE)
  - Contains comments about parameters:
    * index: ignored (inventory index)
    * result: 0 = success, 1 = failure
- repair_list - Repair item list handler (lines 11265-11318)
  - Processes repair item list notifications
  - Clears existing repairList
  - Creates formatted header with "Repair List" title
  - Processes each item in the list:
    * Unpacks index, nameID, upgrade, and cards
    * Gets item from character inventory
    * Stores item name
  - Handles two scenarios:
    * Self repair (myself=1):
      - Checks if item IDs match inventory
      - If mismatch detected, sets myself=0 and rebuilds list
      - Uses inventory binID as index
    * Other player repair (myself=0):
      - Rebuilds entire array
      - Uses item index directly
      - Gets item names using itemNameSimple and itemName
  - Formats each entry with:
    * Index number
    * Short name (30 chars)
    * Full name
  - Adds footer separator line
  - Displays complete formatted list
  - Uses "list" message category
  - Contains comment about packet format (01FC)
  - Contains comment about "dirty hack" for item ID mismatch detection
**Roulette System Handlers:**

- roulette_window - Roulette window opening handler (lines 12056-12086)
  - Processes roulette window opening notifications (0A1A)
  - Stores roulette data in global roulette hash
  - Handles error conditions:
    * Result 1: Generic error - "Something went wrong"
    * Result 2: Insufficient points - "No enough Point (coin) to roll"
  - Displays formatted roulette information:
    * Header with serial number
    * Result, row, column, and bonus item
    * Coin counts (gold, silver, bronze)
  - Warns user when reaching final stage (stage 6)
  - Uses "info" message category for normal messages
  - Uses "warning" message category for alerts

- roulette_info - Roulette rewards information handler (lines 12088-12107)
  - Processes roulette rewards information (0A1C)
  - Handles both old and new packet formats:
    * Old: 8-byte entries with v4 unpacking
    * New: 12-byte entries with v2 V2 unpacking (>= 20180516)
  - Parses item information from binary data:
    * Extracts level, column, item_id, amount
    * Gets item name using itemNameSimple
  - Stores parsed data in roulette.items hash
  - Outputs debug information for each item
  - Organized by level and column for UI display

- roulette_recv_item - Roulette item reward handler (lines 12109-12115)
  - Processes roulette item reward notifications (0A22)
  - Displays message with reward type and bonus item name
  - Uses "info" message category
  - Simple implementation focused on notification

- roulette_window_update - Roulette window update handler (lines 12117-12146)
  - Processes roulette window update notifications (0A20)
  - Updates roulette data in global roulette hash
  - Handles error conditions:
    * Result 1: Generic error - "Something went wrong"
    * Result 2: Insufficient points - "No enough Point (coin) to roll"
  - Displays formatted roulette information:
    * Header with serial number
    * Result, row, column, and bonus item
    * Coin counts (gold, silver, bronze)
    * Highlighted result item name
  - Warns user when reaching final stage (stage 6)
  - Uses "info" message category for normal messages
  - Uses "warning" message category for alerts
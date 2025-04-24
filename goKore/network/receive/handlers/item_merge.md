# Item Merging Related Handlers

**Method Implementations:**
- parse_merge_item_result - Item merge result parser (lines 10095-10098)
  - Parses raw item merge result data
  - Unpacks itemIndex and adjusts by -2
  - Simple implementation focused on data parsing
- merge_item_result - Item merge result handler (lines 10072-10093)
  - Processes item merge result notifications
  - Handles multiple result codes:
    * 0: Success
      - Gets item from character inventory by ID
      - Displays success message
      - If item found:
        * Updates item amount
        * Shows message with old and new amounts
      - If item not found:
        * Shows error about item being moved
    * 1: Cannot merge - "Items cannot be merged"
    * 2: Stack limit - "The amount of merged item will be exceed stack limit"
    * Other: Unknown error with code
  - Outputs debug message with result details
  - Uses "info" message category for success
  - Uses "error" message category for failures
  - Contains comment from author [Cydh]
  - Packet: 096F
- parse_merge_item_open - Item merge list parser (lines 10063-10066)
  - Parses raw item merge list data
  - Unpacks itemList into array of item IDs
  - Creates list of item hashes with ID field
  - Contains comment that received index from server is +2
  - Simple implementation focused on data parsing
- merge_item_open - Item merge list handler (lines 10046-10061)
  - Processes item merge list notifications
  - Clears mergeItemList hash
  - Outputs debug message with number of mergeable items
  - Groups items by nameID for easier merging:
    * For each item in the list:
      - Gets item from character inventory by ID
      - Creates entry in mergeItemList if not exists
      - Sets name from item
      - Initializes empty list array
      - Adds item to list with ID and item info
      - Outputs debug message with item details
  - Displays message with number of mergeable items
  - Uses "info" message category
  - Contains comment from author [Cydh]
  - Packet: 096D
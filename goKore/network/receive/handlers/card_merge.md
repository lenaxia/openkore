# Card Merging Related Handlers

**Method Implementations:**
- card_merge_status - Card merge result handler (lines 10543-10584)
  - Processes card merge result notifications
  - Contains comment: "something about successful compound?"
  - Gets item_index, card_index, and fail flag from packet
  - Handles two fail values:
    * 1: Failure - "Card merging failed"
    * 0: Success
      - Gets item and card from character inventory
      - Displays success message with card and item names
      - Removes one card from inventory
      - Updates item cards data:
        * Uses bytes and disables utf8 encoding
        * Iterates through card slots (4 slots)
        * Keeps existing cards
        * Adds new card to first empty slot
        * Fills remaining slots with zeros
      - Updates item name with new card data
  - Clears cardMergeItemsID array and cardMergeIndex variable
  - Uses "success" message category for successful merges
  - Contains FIXME comment: "this is unoptimized"
- card_merge_list - Card merge list handler (lines 10525-10541)
  - Processes card merge list notifications
  - Contains comment: "You just requested a list of possible items to merge a card into"
  - Contains comment: "The RO client does this when you double click a card"
  - Gets raw message and unpacks length
  - Processes each item in the list:
    * Unpacks item index
    * Gets item from character inventory by ID
    * Adds item binID to cardMergeItemsID array
  - Runs 'card mergelist' command after processing
  - Simple implementation focused on building merge list
# Search Store Related Handlers

**Method Implementations:**
- search_store_pos - Search store position handler (lines 9280-9284)
  - Processes search store position notifications
  - Displays message with store coordinates (x, y)
  - Simple implementation focused on notification
- search_store_result - Search store results handler (lines 9241-9278)
  - Processes search store results
  - Determines packet structure based on length:
    * 114 bytes per entry: Standard format
    * 131 bytes per entry: Extended format with additional field
  - Clears previous results if first page
  - Updates has_next flag for pagination
  - Processes each store entry:
    * Unpacks store data: storeID, accountID, shopName, nameID, etc.
    * Extracts card information
    * Creates universalCatalogInfo hash with all fields
    * Adds to current page array
    * Triggers search_store hook for each entry
  - Returns early if no entries found
  - Adds current page to universalCatalog list
  - Calls Misc::searchStoreInfo with page index
  - Complex implementation for search results processing
- search_store_fail - Search store failure handler (lines 9221-9237)
  - Processes search store failure notifications
  - Displays generic error message with reason code
  - Handles specific reason codes:
    * 0: Shows message from msgTable[1804]
    * 1: Shows message from msgTable[1785]
    * 2: Shows message from msgTable[1799]
    * 3: Shows message from msgTable[1801]
    * 4: Shows message from msgTable[1798]
    * Other: Shows "Unknown reason"
  - Uses "error" message category
  - Simple implementation focused on error notification
- search_store_open - Search store open handler (lines 9209-9219)
  - Processes search store open notifications
  - Outputs debug message with catalog type:
    * 0: Universal Catalog Silver
    * 1: Universal Catalog Gold
  - Displays message with remaining search count
  - Updates universalCatalog data:
    * Sets open flag to 1
    * Sets type to catalog type
  - Uses debug level 2 with "search_store" category
  - Simple implementation focused on notification
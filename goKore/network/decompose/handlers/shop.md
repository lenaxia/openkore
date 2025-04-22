# Shop Related Handlers

**Method Implementations:**
- cash_buy_fail - Cash shop purchase failure handler (lines 10614-10617)
  - Processes cash shop purchase failure notifications
  - Outputs debug message with:
    * Cash points
    * Kafra points
    * Failure code
  - Contains TODO comment indicating incomplete implementation
  - Packet: 0289
  - Simple implementation focused on debugging
- cash_dealer - Cash shop dealer handler (lines 10008-10040)
  - Processes cash shop dealer notifications
  - Clears talk hash
  - Sets AI talk state to 'cash'
  - Updates AI talk timestamp
  - Clears cashList
  - Parses item list using complex unpacking:
    * Extracts price, price_discount, type, nameid for each item
  - Contains comment about keeping cash_points and kafra_points locally
  - Displays formatted header with cash points and kafra points
  - Processes each item:
    * Creates new Actor::Item
    * Sets price, type, nameID
    * Sets ID to cashList size
    * Gets item name using itemName function
    * Adds item to cashList
    * Outputs debug message with item name and price
    * Displays formatted item details with ID, name, type, discounted price
  - Displays footer separator
  - Complex implementation for cash shop initialization
- vender_buy_fail - Vendor purchase failure handler (lines 9986-10002)
  - Processes vendor purchase failure notifications
  - Handles multiple failure codes:
    * 1: Insufficient zeny
    * 2: Overweight
    * 4: Requested more than available stock
    * 6: Vendor refreshed shop before purchase
    * 8: Vendor would exceed max zeny
    * Other: Unknown error with code
  - Displays detailed error messages with:
    * Amount attempted to purchase
    * Item ID
    * Failure reason
    * Error code
  - Uses "error" message category for all messages
- sell_result - Item selling result handler (lines 9519-9527)
  - Processes item selling result notifications
  - Handles two result states:
    * fail=1: Failure - "Sell failed"
    * fail=0: Success - "Sold X items" and "Sell completed"
  - Clears sellList array after processing
  - Checks if AI is in sellAuto mode
  - Uses "error" message category for failures
  - Uses "success" message category for successful sales
  - Packet: 00CB
- npc_market_purchase_result - Market purchase result handler (lines 7763-7781)
  - Processes NPC market purchase results (PACKET_ZC_NPC_MARKET_OPEN)
  - Outputs debug message with result code
  - Handles multiple result codes:
    * MARKET_BUY_RESULT_ERROR (-1): General error
    * MARKET_BUY_RESULT_SUCCESS (0): Success
    * MARKET_BUY_RESULT_NO_ZENY (1): Insufficient zeny
    * MARKET_BUY_RESULT_OVER_WEIGHT (2): Overweight
    * MARKET_BUY_RESULT_OUT_OF_SPACE (3): No inventory space
    * MARKET_BUY_RESULT_AMOUNT_TOO_BIG (4): Amount exceeds available stock
    * Other: Unknown error with code
  - Displays appropriate success/error messages
  - Uses "info" message category for all messages
  - Packet: 09D7
- npc_market_info - NPC market item list handler (lines 7710-7751)
  - Processes NPC market shop item listings (PACKET_ZC_NPC_MARKET_OPEN)
  - Uses server-specific pack format (npc_market_info_pack) or default 'v C V2 v'
  - Clears storeList and talk hash
  - Processes each item in the itemList:
    * Creates new Actor::Item
    * Unpacks item data (nameID, type, price, amount, weight)
    * Skips items with zero amount (client behavior)
    * Sets item ID to storeList size (workaround for duplicate items)
    * Sets item name using itemName function
    * Adds item to storeList
  - Returns early if no items in store
  - Automatically runs 'store' command if not in buyAuto mode
  - Sets in_market flag to 1
  - Updates AI variables:
    * Sets npc_talk state to 'store'
    * Records timestamp
  - Outputs debug messages for each item added
  - Packets: 09D5 (two versions with different nameID field sizes)
- buy_result - Shop purchase result handler (lines 7681-7704)
  - Processes NPC shop purchase results (ZC_PC_PURCHASE_RESULT)
  - Handles multiple result codes:
    * 0: Success - "Buy completed"
    * 1: Insufficient zeny
    * 2: Insufficient weight capacity
    * 3: Too many different inventory items
    * 4: Item does not exist in store
    * 5: Item cannot be exchanged
    * 6: Invalid store
    * Other: Generic failure with code
  - Sets recv_buy_packet flag if in buyAuto mode
  - Triggers buy_result hook with fail code
  - Uses appropriate message categories:
    * "success" for successful purchases
    * Default (error) for failures
  - Packet: 00CA
- npc_sell_list - NPC sellable items handler (lines 7637-7663)
  - Processes list of items that can be sold to NPC (ZC_PC_SELL_ITEMLIST)
  - Checks if message contains item data
  - Processes each item in itemsdata:
    * Unpacks index, price, and overcharge price
    * Gets item reference from inventory by ID
    * Flags item as sellable
    * Displays item amount, name, and price
  - Flags all other inventory items as unsellable:
    * Skips equipped items
    * Skips already flagged sellable items
    * Sets unsellable flag
  - Clears talk hash
  - Displays "Ready to start selling items" message
  - Updates AI variables:
    * Sets npc_talk state to 'sell'
    * Records timestamp
  - Packet: 00C7
- npc_store_info - NPC shop item list handler (lines 7588-7633)
  - Processes NPC shop item listings (ZC_PC_PURCHASE_ITEMLIST)
  - Supports multiple packet versions:
    * 0B77: Uses V3 C v V pack format with nameID, price, type, sprite_id, location
    * Others: Uses server-specific pack format (npc_store_info_pack) or default 'V V C v'
  - Clears storeList and talk hash
  - Processes each item in the raw message:
    * Creates new Actor::Item
    * Unpacks item data using appropriate format
    * Sets item ID to storeList size (workaround for duplicate items)
    * Sets item name using itemName function
    * Adds item to storeList
  - Updates AI variables:
    * Sets npc_talk state to 'store'
    * Records timestamp
  - Automatically runs 'store' command if not in buyAuto mode
  - Outputs debug messages for each item added
  - Packets: 00C6, 0B77
- npc_store_begin - NPC shop dialog handler (lines 7573-7581)
  - Processes NPC shop dialog notifications (ZC_SELECT_DEALTYPE)
  - Clears talk hash to reset dialog state
  - Sets talk{ID} to NPC ID
  - Updates AI variables:
    * Sets npc_talk state to 'buy_or_sell'
    * Records timestamp
  - Sets storeList NPC name using getNPCName
  - Falls back to "Unknown" if name not found
  - Simple implementation focused on shop dialog initialization
  - Packet: 00C4
- cash_shop_buy_result - Cash shop purchase result handler (lines 4493-4526)
  - Processes cash shop purchase results
  - Handles multiple result codes:
    * 0: Success
    * 1: Wrong Tab
    * 2: Shortage cash
    * 3: Unknown item
    * 4: Inventory weight
    * 5: Inventory item count
    * 9: Rune overcount
    * 10: Eachitem overcount
    * 11: Unknown
    * 12: Busy
  - Displays appropriate success/error messages
  - Updates cash points on successful purchase
  - Provides detailed debug logging
- cash_shop_open_result - Cash shop opening handler (lines 4483-4491)
  - Processes cash shop opening result
  - Displays current cash points and kafra points
  - Stores points in cashShop{points} hash:
    * cash => cash_points
    * kafra => kafra_points
  - Packet: 0845 (cash_window_shop_open)
- cash_shop_list - Cash shop item list handler (lines 4443-4481)
  - Processes cash shop item listings by tab
  - Supports multiple tab types:
    * New (0)
    * Popular (1)
    * Limited (2)
    * Rental (3)
    * Perpetuity (4)
    * Buff (5)
    * Recovery (6)
    * Etc (7)
  - Unpacks item ID and price information
  - Stores items in cashShop{list} array
  - Formats and displays item list with prices
- vender_items_list - Vendor shop items handler (lines 3942-3991)
  - Processes vendor shop item listings
  - Stores vendor ID and character ID
  - Clears and populates venderItemList
  - Formats and displays vendor inventory
  - Handles item expiration dates
  - Triggers packet_vender_store and packet_vender_store2 hooks
  - Debug logs items in vendor store
- vending_start - Shop opening handler (lines 3908-3940)
  - Processes shop opening notifications
  - Initializes articles array for shop items
  - Unpacks item data from server packet
  - Formats and displays shop inventory
  - Tracks shop earnings
  - Debug logs items added to shop
  - TODO: Read shop title from packet instead of using $shop{title}
- shop_sold_long - Extended shop sale handler (lines 3861-3905)
  - Enhanced version of shop_sold with additional data
  - Processes item sales with timestamp and buyer ID
  - Includes formatted date in sale messages
  - Provides buyer character ID in hook data
  - Handles sold-out notifications
  - Closes shop when all items are sold
  - Packet: (Not specified in comments)
- shop_sold - Shop item sale handler (lines 3820-3859)
  - Processes individual item sales in player's shop
  - Updates article quantities and earnings
  - Logs sales to shop log if configured
  - Triggers vending_item_sold hook
  - Handles sold-out notifications
  - Closes shop when all items are sold
  - Packet: (Not specified in comments)
- shop_skill - Shop skill usage handler (lines 3810-3816)
  - Processes shop skill usage notifications
  - Displays number of items that can be sold
  - Packet: 012D
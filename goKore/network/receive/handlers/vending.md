# Shop Related Handlers

**Method Implementations:**
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


- open_store_status - Store setup result handler (lines 11839-11851)
  - Processes store setup result notifications
  - Handles two flag values:
    * 0: Success - "Store set up successfully"
      - Triggers open_store_success hook
      - Uses "success" message category
    * Other: Failure - "Failed setting up shop with error code X"
      - Triggers open_store_fail hook with flag
      - Uses error function
  - Simple implementation focused on notification
- vender_lost - Vender removal handler (lines 11704-11710)
  - Processes vender removal notifications
  - Gets vender ID from packet
  - Removes ID from venderListsID array using binRemove
  - Deletes vender information from venderLists hash
  - Simple implementation focused on cleanup
- vender_found - Vender discovery handler (lines 11689-11702)
  - Processes vender discovery notifications
  - Gets vender ID from packet
  - Checks if vender is already in list
  - For new venders:
    * Adds ID to venderListsID array using binAdd
    * Triggers packet_vender hook with ID and title
  - Updates vender information in venderLists hash:
    * Sets title (converted from bytes)
    * Sets id
  - Contains comment: "You see a vender! Add them to the visible venders list."


- buying_buy_fail - Buying failure handler (lines 9886-9897)
  - Processes buying failure notifications
  - Handles multiple result codes:
    * 3: Insufficient zeny - "Failed to buying (insufficient zeny)"
    * 4: Buying complete - "Buying up complete"
    * Other: Unknown error with code
  - For result 4:
    * Sets buyershopstarted to 0
    * Triggers buyer_shop_closed hook
    * Uses regular message category
  - For other results:
    * Uses "error" message category
- buyer_lost - Buyer removal handler (lines 9878-9884)
  - Processes buyer removal notifications
  - Gets buyer ID from packet
  - Removes ID from buyerListsID array
  - Deletes buyer from buyerLists hash
  - Similar to buying_store_lost but for buyers
  - Simple implementation focused on buyer cleanup
- buyer_found - Buyer discovery handler (lines 9866-9876)
  - Processes buyer discovery notifications
  - Gets buyer ID from packet
  - Checks if buyer exists in buyerLists hash
  - If new buyer:
    * Adds ID to buyerListsID array
    * Triggers packet_buyer hook with ID
  - Updates buyer information:
    * title: Buyer title (converted from bytes)
    * id: Buyer ID
  - Similar to buying_store_found but with different hook
- buying_store_update - Buying store update handler (lines 9853-9864)
  - Processes buying store update notifications
  - Checks if selfBuyerItemList exists
  - Iterates through all items in selfBuyerItemList
  - For matching item (by nameID):
    * Displays message with count and item name
    * Updates item amount in selfBuyerItemList
  - Simple implementation focused on inventory update and notification
- buying_store_fail - Buying store failure handler (lines 9840-9851)
  - Processes buying store failure notifications
  - Handles multiple result codes:
    * 5: Deal failure - "The deal has failed"
    * 6: Insufficient items - "X item could not be sold because you do not have the wanted amount of items"
    * 7: Insufficient zeny - "Failed to deal because you have not enough Zeny"
    * Other: Unknown result with code
  - Uses itemNameSimple function to get item name for result 6
  - Contains comments referencing msgstring IDs (58, 1748, 1746)
  - Uses "error" message category for all messages
- buying_store_item_delete - Buying store item sale handler (lines 9829-9838)
  - Processes buying store item sale notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from character inventory by ID
  - Calculates total zeny earned (amount * zeny)
  - If item exists:
    * Calls inventoryItemRemoved with binID and amount
  - Displays message with item name, amount, and total zeny
  - Contains comment referencing msgstring 1747
  - Simple implementation focused on inventory update and notification
- buying_store_items_list - Buying store items list handler (lines 9758-9827)
  - Processes buying store items list notifications
  - Clears global variables:
    * buyerPriceLimit: Price limit
    * buyerID: Buyer ID
    * buyingStoreID: Store ID
  - Clears buyerItemList
  - Sets global variables from packet
  - Gets player actor using Actor::get
  - Uses server-specific pack format or default 'V v C v'
  - Calculates item entry length
  - Creates formatted header with buyer name
  - Processes each item in the list:
    * Creates new Actor::Item
    * Unpacks item data: price, amount, type, nameID
    * Gets item name using itemName function
    * Sets item ID to position in list
    * Adds item to buyerItemList
    * Outputs debug message with item name and price
    * Triggers packet_buying_store hook with item data
    * Formats and adds item to display message
  - Adds price limit to display message
  - Displays complete message with "list" category
  - Handles expiration date if present:
    * Converts to Unix timestamp
    * Displays formatted date
  - Triggers packet_buying_store2 hook with:
    * buyerID: Buyer ID
    * buyingStoreID: Store ID
    * itemList: Complete item list
    * expireDate: Expiration date
  - Complex implementation for buying store display
- buying_store_lost - Buying store removal handler (lines 9750-9756)
  - Processes buying store removal notifications
  - Gets store ID from packet
  - Removes ID from buyerListsID array
  - Deletes store from buyerLists hash
  - Simple implementation focused on store cleanup
- buying_store_found - Buying store discovery handler (lines 9738-9748)
  - Processes buying store discovery notifications
  - Gets store ID from packet
  - Checks if store exists in buyerLists hash
  - If new store:
    * Adds ID to buyerListsID array
    * Triggers packet_buying hook with ID
  - Updates store information:
    * title: Store title (converted from bytes)
    * id: Store ID
  - Simple implementation focused on store tracking
- open_buying_store_item_list - Buying store item list handler (lines 9697-9736)
  - Processes buying store item list notifications
  - Gets raw message and size
  - Sets header length (12 bytes)
  - Uses server-specific pack format or default 'V v C v'
  - Calculates item entry length
  - Clears selfBuyerItemList array
  - Displays "Buying Shop opened!" message with "BuyShop" category
  - Contains commented code about @articles
  - Processes each item in the list:
    * Unpacks item data: price, amount, type, nameID
    * Gets item name using itemName function
    * Adds item to selfBuyerItemList array
    * Triggers packet_open_buying_store hook with item data
  - Runs 'bs' command after processing all items
  - Complex implementation for buying store initialization
- buyer_items - Buyer items list handler (lines 9674-9695)
  - Processes buyer items list notifications
  - Gets vendor information:
    * Binary ID from packet
    * Player actor using Actor::get
    * Player name from actor
  - Calculates header length (12 bytes)
  - Gets total value from packet
  - Processes each item in the list:
    * Unpacks item data: price, amount, nameID
    * Creates item hash with extracted data
  - Contains TODO comment indicating incomplete implementation
  - No output or storage of processed items
- open_buying_store - Buying store open handler (lines 9667-9671)
  - Processes buying store open notifications
  - Gets maximum number of items from packet
  - Displays message with item limit
  - Simple implementation focused on notification
- open_buying_store_fail - Buying store open failure handler (lines 9194-9207)
  - Processes buying store open failure notifications
  - Handles multiple failure reasons:
    * 1: Generic failure - "Failed to open Purchasing Store"
    * 2: Weight limit - "The total weight of the item exceeds your weight limit"
    * 8: Invalid shop info - "Shop information is incorrect and cannot be opened"
    * Other: Generic message - "Failed opening your buying store"
  - Sets buyershopstarted flag to 0
  - Uses "info" message category
  - Packet: 0x812

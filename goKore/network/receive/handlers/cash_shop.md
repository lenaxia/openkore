# Cash Shop Related Handlers

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
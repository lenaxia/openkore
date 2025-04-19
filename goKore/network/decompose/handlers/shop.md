**Shop Handlers:**

**Market Handlers:**
- npc_market_info (lines 7710-7751)
  - Handles NPC market shop item list (PACKET_ZC_NPC_MARKET_OPEN)
  - Features:
    - Processes different packet formats
    - Creates item entries with price/type/amount/weight
    - Handles duplicate items
    - Updates store state
    - Triggers 'store' command if not in buyAuto

- npc_market_purchase_result (lines 7764-7826)
  - Handles market purchase results (PACKET_ZC_NPC_MARKET_OPEN)
  - Features:
    - Processes all result codes:
      - Success
      - No zeny
      - Overweight
      - No inventory space
      - Amount too big
      - Unknown error
    - Updates store list after purchase
    - Maintains market state

**Trade Deal Handlers:**
- deal_add_other (lines 5826-5844)
  - Handles items/zeny added to trade by other player
  - Processes:
    - Items: nameID, amount, identified, broken, upgrade, cards, options
    - Zeny: amount added to deal_other_zeny
  - Displays formatted messages for added items/zeny

- deal_begin (lines 5847-5880)
  - Handles trade initiation responses:
    - 0: Target too far
    - 2: Target already trading
    - 3: Success case
    - 5: Target opening storage
  - Sets up currentDeal with name/ID
  - Calls 'engaged_deal' hook on success
  - Uses "deal" message type

- deal_cancelled (lines 5882-5888)
  - Cleans up deal state variables
  - Calls 'cancelled_deal' hook
  - Displays cancellation message

- deal_complete (lines 5890-5896)
  - Cleans up deal state variables
  - Calls 'complete_deal' hook
  - Displays completion message

- deal_finalize (lines 5898-5911)
  - Handles trade finalization:
    - type 1: Other player finalized
    - else: You finalized (deducts zeny)
  - Calls 'finalized_deal' hook

- deal_request (lines 5913-5927)
  - Handles incoming trade requests
  - Sets timeout for auto-cancel
  - Displays request message with level
  - Calls 'incoming_deal' hook

**Cash Shop Handlers:**
- cash_shop_list (lines 4443-4481)
  - Lists items in cash shop tabs:
    - 0: New
    - 1: Popular
    - 2: Limited
    - 3: Rental
    - 4: Perpetuity
    - 5: Buff
    - 6: Recovery
    - 7: Etc
  - Displays formatted item list with prices
  - Stores items in cashShop{list} structure

- cash_shop_open_result (lines 4483-4491)
  - Shows available cash points
  - Processes ZC_CASH_SHOP_OPEN packet
  - Stores points in cashShop{points} structure
  - Format: "Cash Points: XC - Kafra Points: XC"

- cash_shop_buy_result (lines 4493-4526)
  - Handles purchase results:
    - 0: Success
    - 1: Wrong Tab
    - 2: Insufficient Cash
    - 3: Unknown Item
    - 4: Overweight
    - 5: Full Inventory
    - 9: Rune Overcount
    - 10: Item Overcount
    - 11: Unknown Error
    - 12: Busy
  - Updates cash points on success
  - Shows appropriate error messages

**Shop/Vending System Handlers:**

- shop_skill (lines 3810-3816)
  - Handles shop skill activation
  - Shows message with number of items that can be sold
  - Processes ZC_SHOP_SKILL packet

- shop_sold (lines 3820-3859)
  - Processes item sales (short format)
  - Updates sold item quantities and earnings
  - Logs sales transactions
  - Triggers hooks:
    - 'vending_item_sold'
    - 'vending_item_sold_out'
  - Handles shop closure when items sold out

- shop_sold_long (lines 3861-3905)
  - Processes item sales (long format with buyer info)
  - Includes additional details:
    - Buyer character ID
    - Transaction timestamp
    - Exact zeny earned
  - Same hooks as shop_sold

- vending_start (lines 3908-3940)
  - Handles shop opening
  - Processes item list and prices
  - Displays formatted shop interface
  - Initializes shop state variables

- vender_items_list (lines 3942-3991)
  - Processes vending shop item lists
  - Handles both player and NPC vendors
  - Displays formatted item list with:
    - Prices
    - Quantities
    - Item types
  - Supports expiration dates for timed shops
  - Triggers 'packet_vender_store' and 'packet_vender_store2' hooks
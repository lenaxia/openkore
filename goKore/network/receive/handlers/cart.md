# Cart Related Handlers

**Method Implementations:**
- cart_add_failed - Cart addition failure handler (lines 5196-5208)
  - Processes cart addition failure notifications
  - Handles different failure reasons:
    * 0: Overweight (cart weight limit exceeded)
    * 1: Too many items (cart item count limit exceeded)
    * Other: Unknown code
  - Displays appropriate error message
- cart_info - Cart status information handler (lines 5190-5194)
  - Processes cart status information
  - Updates character's cart object with:
    * Current item count
    * Maximum item count
    * Current weight
    * Maximum weight
  - Logs debug message when received
  - Packet: 0121 (ZC_NOTIFY_CARTITEM_COUNTINFO)
- cart_item_removed - Cart item removal handler (lines 5176-5186)
  - Processes item removal from cart
  - Extracts item ID and amount from args
  - Retrieves item from cart by ID
  - Calls Misc::cartItemRemoved helper for actual removal
  - Simple implementation that delegates to utility function
  - Similar to storage_item_removed but for cart
- cart_item_added - Cart item addition handler (lines 5146-5174)
  - Processes individual item additions to cart
  - Handles both new items and existing items:
    * New items: Creates Actor::Item and adds to cart
    * Existing items: Updates amount
  - Sets item properties (nameID, type, identified, etc.)
  - Displays addition message with item details
  - Updates itemChange tracking
  - Stores item reference in args for hooks
  - Similar to storage_item_added but for cart
- cart_items_nonstackable - Non-stackable cart items handler (lines 5133-5144)
  - Processes non-stackable items in cart
  - Uses _items_list helper with:
    * Actor::Item class
    * packet_cart hook
    * parse_items_nonstackable parser
    * Cart-specific getter and adder
  - Similar to storage_items_nonstackable but for cart
- cart_items_stackable - Stackable cart items handler (lines 5120-5131)
  - Processes stackable items in cart
  - Uses _items_list helper with:
    * Actor::Item class
    * packet_cart hook
    * parse_items_stackable parser
    * Cart-specific getter and adder
  - Similar to storage_items_stackable but for cart
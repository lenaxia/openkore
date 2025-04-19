**Cart Item Handlers:**
- cart_items_stackable (lines 5121-5131)
  - Processes stackable cart items
  - Uses _items_list helper
  - Handles ZC_CART_ITEMLIST packets

- cart_items_nonstackable (lines 5133-5144)
  - Processes non-stackable cart items
  - Uses _items_list helper
  - Handles ZC_CART_EQUIPMENTLIST packets

- cart_item_added (lines 5147-5174)
  - Handles new cart items
  - Creates new item or updates count
  - Logs item additions
  - Processes ZC_ADD_ITEM_TO_CART packets

- cart_item_removed (lines 5176-5186)
  - Processes cart item removal
  - Updates item counts
  - Handles ZC_DELETE_ITEM_FROM_CART packets

**Cart Management Handlers:**
- cart_info (lines 5191-5194)
  - Updates cart weight/count info
  - Processes ZC_NOTIFY_CARTITEM_COUNTINFO

- cart_add_failed (lines 5197-5208)
  - Handles cart add failures
  - Processes error codes:
    - 0: Overweight
    - 1: Too many items
  - Shows appropriate error messages

**Cart System Handlers:**

- cart_off (lines 3804-3807)
  - Handles cart release/closure
  - Updates cart state via $char->cart->close
  - Shows success message "Cart released"
  - Processes ZC_CART_OFF packet
# Item Related Handlers

**Method Implementations:**
- use_item - Item usage handler (lines 11673-11681)
  - Processes item usage notifications
  - Requires in-game state (changeToInGameState)
  - Gets item from inventory by ID
  - Displays message with item name, binID, and amount
  - Calls inventoryItemRemoved to update inventory
  - Uses "useItem" message category
  - Contains TODO comment: "only used to report failure? $args->{success}"
  - Simple implementation focused on notification and inventory update

- item_preview - Item preview handler (lines 12182-12193)
  - Processes item preview notifications
  - Gets item reference from inventory using ID
  - Updates item properties:
    * broken status (if defined)
    * upgrade level
    * card information
    * option information
  - Updates item name using itemName function
  - Simple implementation focused on item data updates


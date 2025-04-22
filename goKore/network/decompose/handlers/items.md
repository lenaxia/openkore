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
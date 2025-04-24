# Item List Related Handlers

**Method Implementations:**
- item_list_end - Item list finalization handler (lines 5341-5354)
  - Processes item list finalization for different container types
  - Logs debug message with container type
  - Handles different container types:
    * INVTYPE_INVENTORY: Character inventory
    * INVTYPE_CART: Character cart
    * INVTYPE_STORAGE/INVTYPE_GUILD_STORAGE: Storage/guild storage
  - Calls appropriate onitemListEnd method for each container
  - Clears current_item_list global variable
  - Warns on unsupported container types
- item_list_nonstackable - Non-stackable item list handler (lines 5291-5339)
  - Processes non-stackable items for different container types
  - Requires in-game state (changeToInGameState)
  - Creates arguments for _items_list helper:
    * Actor::Item class
    * Container-specific hook, getter, and adder
    * parse_items_nonstackable parser
  - Special handling for equipped items:
    * Updates character equipment references
    * Skips special slots (arrow bug workaround)
  - Supports multiple container types:
    * INVTYPE_INVENTORY: Character inventory
    * INVTYPE_CART: Character cart
    * INVTYPE_STORAGE/INVTYPE_GUILD_STORAGE: Storage/guild storage
  - Warns on unsupported container types
- item_list_stackable - Stackable item list handler (lines 5249-5289)
  - Processes stackable items for different container types
  - Requires in-game state (changeToInGameState)
  - Creates arguments for _items_list helper:
    * Actor::Item class
    * Container-specific hook, getter, and adder
    * parse_items_stackable parser
  - Special handling for equipped arrows
  - Supports multiple container types:
    * INVTYPE_INVENTORY: Character inventory
    * INVTYPE_CART: Character cart
    * INVTYPE_STORAGE/INVTYPE_GUILD_STORAGE: Storage/guild storage
  - Warns on unsupported container types
- item_list_start - Item list initialization handler (lines 5232-5247)
  - Processes item list initialization for different container types
  - Sets current_item_list global variable
  - Logs debug message with container type and name
  - Handles different container types:
    * INVTYPE_INVENTORY: Character inventory
    * INVTYPE_CART: Character cart
    * INVTYPE_STORAGE/INVTYPE_GUILD_STORAGE: Storage/guild storage
  - Calls appropriate onitemListStart method for each container
  - Warns on unsupported container types
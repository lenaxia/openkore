# Inventory Related Handlers

**Method Implementations:**
- identify - Item identification result handler (lines 6917-6930)
  - Processes item identification result notifications
  - Handles success/failure based on flag value:
    * 0: Successful identification
      - Gets item from character's inventory by ID
      - Sets identified flag to 1
      - Updates item's type_equip based on nameID
      - Displays success message with item name and binID
    * Other: Failed identification
      - Displays "Item Appraisal has failed" message
  - Clears identifyID array after processing
  - Uses "info" message category for success messages
  - Packet: 0179 (ZC_ACK_ITEMIDENTIFY)
  - Format: 'W B' (ID flag)
- identify_list - Identifiable items list handler (lines 6898-6915)
  - Processes list of items that can be identified
  - Clears existing identifyID array
  - Parses raw message data to extract item IDs:
    * Processes message in 2-byte chunks
    * Gets item from character's inventory by ID
    * Adds item's binID to identifyID array
  - Counts number of identifiable items
  - Displays message with item count
  - Suggests using 'identify' command
  - Uses "info" message category
  - Packet: 0177 (ZC_ITEMIDENTIFY_LIST)
  - Format: 'W a*' (len itemList)
- item_upgrade - Item upgrade result handler (lines 7120-7131)
  - Processes item upgrade result notifications
  - Extracts type, ID, and upgrade level from args
  - Gets item from character's inventory by ID
  - Updates item's upgrade level
  - Displays message with item name and new upgrade level
  - Uses "parseMsg/upgrade" message category
  - Updates item name to reflect new upgrade level
  - Simple implementation focused on inventory update
- item_used - Item usage result handler (lines 6956-7002)
  - Processes item usage result notifications
  - Extracts ID, itemID, actorID, remaining, and success from args
  - Prepares hook arguments for packet_useitem hook
  - Handles different user scenarios:
    * Self usage (ID matches accountID):
      - Gets item from character's inventory by ID
      - For successful usage:
        * Calculates amount used
        * Displays success message with item details
        * Calls inventoryItemRemoved to update inventory
        * Adds item details to hook arguments
      - For failed usage:
        * Displays failure message with item name
      - Handles unknown items with generic messages
    * Other actor usage:
      - Gets actor reference using Actor::get
      - Gets item name using itemNameSimple
      - Displays usage message with actor name and item
  - Triggers packet_useitem hook with all arguments
  - Uses "useItem" message category (priority 1 for self, 2 for others)
- inventory_items_stackable - Stackable inventory items handler (lines 5210-5230)
  - Processes stackable items in inventory
  - Requires in-game state (changeToInGameState)
  - Uses _items_list helper with:
    * Actor::Item class
    * packet_inventory hook
    * parse_items_stackable parser
    * Inventory-specific getter and adder
  - Special handling for equipped arrows
  - Similar to storage_items_stackable but for inventory
- cart_off - Cart release handler (lines 3804-3807)
  - Closes the character's cart
  - Displays cart release message
  - Packet: 012B
- rental_expired - Rental item expiration handler (lines 3789-3801)
  - Processes rental item expiration notifications
  - Removes expired items from inventory
  - Displays expiration message to user
  - Triggers rental_expired hook with item details
  - Packet: 0299
- arrow_equipped - Arrow/bullet equipment handler (lines 3653-3670)
  - Updates character's equipped arrow/bullet
  - Sets the arrow field in character data
  - Updates equipment inventory reference
  - Handles waitForEquip AI variable
  - Triggers equipped_item hook
  - Packet: 013C

**Method Implementations:**
- inventory_item_favorite - Inventory item favorite status handler (lines 9955-9963)
  - Processes inventory item favorite status notifications
  - Gets item from character inventory by ID
  - Handles two flag values:
    * 1: Item removed from favorite tab - "Inventory Item removed from favorite tab"
    * 0: Item moved to favorite tab - "Inventory Item move to favorite tab"
  - Displays message with item name
  - Uses "storage" message category
  - Simple implementation focused on notification
- inventory_item_added - Item addition handler (lines 3679-3753)
  - Handles adding new items to inventory
  - Manages both new items and stackable additions
  - Processes item properties (ID, name, type, amount, etc.)
  - Handles failure cases (inventory full, frozen, etc.)
  - Triggers item_gathered hook
  - Implements auto-drop logic
  - Packets: 00A0, 029A, 02D4, 0990, 0A0C, 0A37

- inventory_item_removed - Item removal handler (lines 3756-3778)
  - Handles item removal from inventory
  - Processes different removal reasons:
    - Skill casting (1)
    - Refinement failure (2)
    - Chemical reaction (3)
    - Storage transfer (4)
    - Cart transfer (5)
    - Selling (6)
    - Four Spirit Analysis skill (7)
    - Unknown reasons (default)
  - Packets: 00AF, 07FA

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

# Storage Related Handlers

**Method Implementations:**
- storage_password_result - Storage password result handler (lines 11056-11089)
  - Processes storage password result notifications
  - Handles multiple type values:
    * 4 (STORE_PASSWORD_CHANGE_OK): "Successfully changed storage password"
    * 5 (STORE_PASSWORD_CHANGE_NG): "Error: Incorrect storage password"
    * 6 (STORE_PASSWORD_CHECK_OK): "Successfully entered storage password"
    * 7 (STORE_PASSWORD_CHECK_NG): "Error: Incorrect storage password"
      - Also disables storageAuto configuration
      - Removes storageAuto from AI queue
  - Contains TODO comment with constant definitions:
    * STORE_PASSWORD_EMPTY = 0x0
    * STORE_PASSWORD_EXIST = 0x1
    * STORE_PASSWORD_CHANGE = 0x2
    * STORE_PASSWORD_CHECK = 0x3
    * STORE_PASSWORD_PANALTY = 0x8
  - Contains comment about unknown purpose of val parameter
  - Uses "success" message category for success
  - Uses error function for errors
- storage_password_request - Storage password request handler (lines 10971-11053)
  - Processes storage password requests
  - Handles multiple flag values:
    * 0: New password request
      - Handles both character and storage passwords (switch 023E vs others)
      - For character password: Prompts for new character password
      - For storage password:
        * If storageAuto_password is empty: Prompts for new storage password
        * Updates configuration with new password
      - Gets encryption key from masterServer
      - Creates Crypton object with key
      - Formats password with length prefix
      - Encrypts password
      - Sends storage password packet with type 2 and 3
      - Displays success message
    * 1: Password verification request
      - Handles both character and storage passwords
      - If password is empty: Prompts for password input
      - Updates configuration with input
      - Gets encryption key from masterServer
      - Creates Crypton object with key
      - Formats password with length prefix
      - Encrypts password
      - Sends storage password packet with type 3
    * 8: Too many wrong attempts
      - Displays error message
      - Disables storageAuto configuration
      - Removes storageAuto from AI queue
    * Other: Debug message about unknown flag
  - Contains error handling for missing storageEncryptKey
  - Uses "success" message category for success
  - Uses error function for errors
  - Uses debug function for unknown flags
- storage_item_removed - Storage item removal handler (lines 5108-5118)
  - Processes item removal from storage
  - Extracts item ID and amount from args
  - Retrieves item from storage by ID
  - Calls Misc::storageItemRemoved helper for actual removal
  - Simple implementation that delegates to utility function
- storage_item_added - Storage item addition handler (lines 5077-5106)
  - Processes individual item additions to storage
  - Handles both new items and existing items:
    * New items: Creates Actor::Item and adds to storage
    * Existing items: Updates amount
  - Sets item properties (nameID, type, identified, etc.)
  - Displays addition message with item details
  - Updates itemChange tracking
  - Stores item reference in args for hooks
- storage_items_nonstackable - Non-stackable storage items handler (lines 5062-5075)
  - Processes non-stackable items in storage
  - Uses _items_list helper with:
    * Actor::Item class
    * packet_storage hook
    * parse_items_nonstackable parser
    * Storage-specific getter and adder
  - Sets storageTitle if provided
  - Similar to storage_items_stackable but without clearing storage
- storage_items_stackable - Stackable storage items handler (lines 5040-5060)
  - Processes stackable items in storage
  - Clears existing storage data
  - Uses _items_list helper with:
    * Actor::Item class
    * packet_storage hook
    * parse_items_stackable parser
    * Storage-specific getter and adder
  - Handles high bit in amount field (amount & ~0x80000000)
  - Sets storageTitle if provided
- storage_closed - Storage closing handler (lines 5026-5027)
  - Processes storage closing notification
  - Calls character storage close method
  - Simple one-line implementation
- storage_opened - Storage opening handler (lines 5021-5024)
  - Processes storage opening notification
  - Calls character storage open method
  - Passes arguments to storage object

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

- inventory_expansion_result - Inventory expansion result handler (lines 12155-12180)
  - Processes inventory expansion result notifications (0B18)
  - Handles multiple result codes:
    * EXPAND_INVENTORY_RESULT_SUCCESS (0x0): Success message
    * EXPAND_INVENTORY_RESULT_FAILED (0x1): Generic failure message
    * EXPAND_INVENTORY_RESULT_OTHER_WORK (0x2): Window closure required
    * EXPAND_INVENTORY_RESULT_MISSING_ITEM (0x3): Missing required item
    * EXPAND_INVENTORY_RESULT_MAX_SIZE (0x4): Maximum limit reached
    * Other: Unknown result message with code
  - Uses "info" message category for all messages
  - References msgstringtable for messages
  - Comprehensive error handling

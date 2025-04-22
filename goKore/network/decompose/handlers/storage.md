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
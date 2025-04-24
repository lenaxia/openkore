# Vending Related Handlers

**Method Implementations:**
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
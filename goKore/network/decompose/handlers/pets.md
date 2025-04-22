# Pet Related Handlers

**Method Implementations:**
- pet_info2 - Extended pet information handler (lines 9046-9098)
  - Processes extended pet information notifications
  - Extracts type, ID, and value from packet
  - Handles multiple information types:
    * 0: No pet owned - Unsets pet ID
    * 1: Friendship update - Sets pet friendly value
    * 2: Hunger update - Sets pet hungry value
    * 3: Accessory update - Sets pet accessory value
    * 4: Performance info - No action (commented debug)
    * 5: Pet ownership - Sets pet ID
  - Contains extensive commented code about pet spawning
  - Contains comments about related Freya functions
  - Outputs debug messages for most update types
  - Simple implementation focused on specific data updates
- pet_info - Pet information handler (lines 9034-9044)
  - Processes pet information notifications
  - Updates pet data structure with:
    * name: Pet name (converted from bytes)
    * renameflag: Flag indicating if pet can be renamed
    * level: Pet level
    * hungry: Pet hunger level
    * friendly: Pet intimacy/friendship level
    * accessory: Pet accessory item ID
    * type: Pet type (if provided)
  - Outputs detailed debug message with all pet status information
  - Uses "pet" debug category
  - Simple implementation focused on data updates
- pet_food - Pet feeding result handler (lines 9025-9032)
  - Processes pet feeding result notifications
  - Handles two result states:
    * success=1: Displays success message with food item name
    * success=0: Displays error message about missing food
  - Uses itemNameSimple function to get food name
  - Uses "pet" message category for success
  - Uses "error" message category for failure
  - Simple implementation focused on result notification
- pet_evolution_result - Pet evolution result handler (lines 9008-9023)
  - Processes pet evolution result notifications
  - Handles multiple result codes:
    * 0x0: General error - "Pet evolution error"
    * 0x3: Accessory error - "Unequip pet accessories first to start evolution"
    * 0x4: Material error - "Insufficient materials for evolution"
    * 0x5: Intimacy error - "Loyal Intimacy is required to evolve"
    * 0x6: Success - "Pet evolution success"
  - Contains commented codes:
    * PET_EVOL_NO_CALLPET = 0x1
    * PET_EVOL_NO_PETEGG = 0x2
  - Uses "error" message category for failures
  - Uses "success" message category for successful evolution
- pet_emotion - Pet emotion display handler (lines 8999-9006)
  - Processes pet emotion notifications
  - Extracts ID and emotion type from packet
  - Looks up emotion display string from emotions_lut table
  - Falls back to "/e$type" format if not found
  - Checks if pet exists in pets hash
  - Displays formatted message with pet name and emotion
  - Uses "emotion" message category
  - Simple implementation focused on display
- pet_capture_result - Pet capture result handler (lines 8990-8997)
  - Processes pet capture result notifications
  - Handles two result states:
    * success=1: Displays "Pet capture success" message
    * success=0: Displays "Pet capture failed" message
  - Uses "info" message category
  - Simple implementation focused on result notification
- pet_capture_process - Pet capture attempt handler (lines 8985-8988)
  - Processes pet capture attempt notifications
  - Displays message about slot machine capture attempt
  - Uses "info" message category
  - Simple implementation focused on notification
- egg_list - Pet egg hatching list handler (lines 5948-5959)
  - Processes list of hatchable pet eggs in inventory
  - Creates formatted header with centered title
  - Iterates through raw message data in 2-byte chunks
  - For each egg:
    * Extracts item index from message
    * Retrieves item from character inventory
    * Displays item binID and name
  - Adds footer with separator line
  - Shows instruction for using 'pet hatch' command
  - Outputs to 'list' message channel
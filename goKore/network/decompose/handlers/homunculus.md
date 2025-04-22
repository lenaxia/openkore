**Homunculus Handlers:**

- enforce_homun_state() - Validates homunculus state (lines 2927-2936)
  - Checks if character exists
  - Creates temporary homunculus if missing
  - Used by homunculus_property() before updates
  - Debug logs when creating new homunculus

- homunculus_state_handler() - Manages homunculus state changes (lines 2938-3002)
  - Handles state bit flags:
    * 0x1 - Named status
    * 0x2 - Vaporized status
    * 0x4 - Dead status
  - Tracks state transitions:
    * New state initialization
    * State change detection
  - Generates appropriate messages for:
    * Renaming/un-renaming
    * Vaporizing/recalling
    * Death/resurrection
  - Updates homunculus_info structure

- homunculus_info() - Handles homunculus state notifications (lines 3022-3055)
  - Processes ZC_CHANGESTATE_MER (0230) packets
  - State types:
    * HO_PRE_INIT (0x0) - Initialization/resurrection
    * HO_RELATIONSHIP_CHANGED (0x1) - Intimacy changes
    * HO_FULLNESS_CHANGED (0x2) - Hunger changes
    * HO_ACCESSORY_CHANGED (0x3) - Accessory changes
    * HO_HEADTYPE_CHANGED (0x4) - Head type changes
  - Special cases:
    * Handles resurrection after teleport
    * Maintains homunculus object reference
    * Adds to SlaveManager if needed

- homunculus_property() - Processes homunculus property packets (lines 2897-2925)
  - Supports multiple packet versions:
    * ZC_PROPERTY_HOMUN (022e)
    * ZC_PROPERTY_HOMUN_2 (09f7) 
    * ZC_PROPERTY_HOMUN3 types 1/2 (0b2f/0b76)
    * ZC_PROPERTY_HOMUN4 (0ba4)
  - Validates homunculus state via enforce_homun_state()
  - Updates homunculus properties:
    * Name (bytesToString conversion)
    * Stats via slave_calcproperty_handler()
    * State via homunculus_state_handler()
  - Handles attack distance configuration:
    * Auto-adjusts homunculus_attackDistance if too large
    * Sets homunculus_attackMaxDistance from server data
  - TODO: Could be refactored with mercenary/char property handlers
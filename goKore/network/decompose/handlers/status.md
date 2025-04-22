# Character Status Related Handlers

**Method Implementations:**

- actor_status_active - Status effect handler (lines 4220-4239)
  - Processes status effect activation notifications
  - Requires in-game state (changeToInGameState)
  - Gets status type, ID, and duration from packet
  - Translates numeric status type to named status using statusHandle
  - Falls back to "UNKNOWN_STATUS_$type" for undefined statuses
  - Special handling for cart activation (type 673)
  - Sets status on actor with appropriate duration
  - Special handling for Rolling Cutter counters (type 0x153):
    * Updates character's spirit counter value
    * Displays appropriate message based on target (self vs other)
    * Uses "parseMsg_statuslook" message category
- resurrection - Resurrection handler (lines 11341-11375)
  - Processes resurrection notifications
  - Gets target ID and type from packet
  - Gets player reference from playersList
  - Handles different target scenarios:
    * Self resurrection (targetID = accountID):
      - Displays "You have been resurrected" message
      - Clears character dead and dead_time flags
      - Sets character resurrected flag
    * Other player resurrection:
      - Clears player dead flag
      - Resets player deltaHp to 0
    * Slave resurrection (isMySlaveID):
      - Enforces homunculus state
      - Gets slave reference
      - For homunculus slaves:
        * Displays "Slave Resurrected" message
        * Clears homunculus_info dead flag
        * Adds homunculus to SlaveManager if needed
      - Displays resurrection message with actor name
  - Uses "info" message category
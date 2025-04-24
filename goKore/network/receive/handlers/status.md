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
- actor_trapped - Actor trap status handler (lines 8160-8166)
  - Processes notifications when an actor (player, monster, NPC) is trapped
  - Gets actor reference using Actor::get with the provided ID
  - Outputs debug message with actor name and trapped status
  - Simple implementation focused on status notification
  - Contains historical comment about ID validity on different servers
  - Does not trigger any plugin hooks
  - Does not modify actor state beyond logging
  - Packet: 01AC (ZC_NOTIFY_ACTORSET)
  - Format: 'a4' (ID)
- character_status - Character status update handler (lines 5550-5567)
  - Processes character status updates
  - Gets actor reference using Actor::get with the provided ID
  - Handles multiple packet formats:
    * 028A: Updates level (lv) and opt3 fields
    * 0229/0119: Updates opt1 and opt2 fields
  - Always updates the option field
  - Calls setStatus to apply status changes to the actor
  - No debug messages or notifications
  - No plugin hooks triggered
  - Packets:
    * 0119: Format 'a4 v3 C' (ID opt1 opt2 option stance)
    * 0229: Format 'a4 v2 V C' (ID opt1 opt2 option stance)
    * 028A: Format 'a4 V3' (ID option lv opt3)
- devotion - Devotion skill status handler (lines 5928-5946)
  - Processes devotion skill notifications (Crusader/Paladin skill)
  - Gets source actor using Actor::get with sourceID
  - Clears existing devotion data for the source actor
  - Processes up to 5 target IDs from the packet:
    * Extracts each 4-byte ID from targetIDs field
    * Stops processing when encountering a zero ID
    * Updates devotionList with source-target relationships
    * Gets target actor references
  - Builds message string using skillUseNoDamage_string
  - Updates devotion range in devotionList
  - Displays message with "devotion" category
  - Contains FIXME comment about needing better display
  - No plugin hooks triggered
  - Packet: 01CF
  - Format: 'a4 a20 v' (sourceID targetIDs range)
- hp_sp_changed - HP/SP update handler (lines 7156-7170)
  - Processes HP and SP change notifications
  - Requires in-game state (changeToInGameState)
  - Handles different type values:
    * Type 5: HP change
      - Adds amount to character's current HP
      - Caps HP at character's maximum HP
    * Type 7: SP change
      - Adds amount to character's current SP
      - Caps SP at character's maximum SP
  - Simple implementation focused on stat updates
  - No messages or notifications
  - No plugin hooks triggered
- overweight_percent - Weight percentage handler (lines 9422-9426)
  - Processes overweight percentage notifications
  - Minimal implementation that only outputs debug message
  - Contains TODO comment indicating incomplete implementation
  - Outputs debug message with percent value
  - No state changes or actions taken
  - No plugin hooks triggered
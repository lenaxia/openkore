**Movement and Actor Display Handlers:**

- move_interrupt - Movement interruption notification handler (lines 11961-11966)
  - Processes movement interruption notifications (0AB8)
  - Outputs debug message when movement is interrupted
  - Simple implementation with minimal functionality
  - Handles interruptions caused by:
    * Casting skills
    * Fleeing from mobs
    * Other movement-stopping events

- no_teleport - Teleport failure handler (lines 11186-11198)
  - Processes teleport failure notifications
  - Handles multiple fail values:
    * 0: Area unavailable - "Unavailable Area To Teleport"
      - Also clears teleport from AI queue
    * 1: Memo unavailable - "Unavailable Area To Memo"
    * Other: Generic failure - "Unavailable Area To Teleport (fail code X)"
  - Uses error function for all messages
  - Simple implementation focused on error handling

- actor_movement_interrupted - Movement interruption handler (lines 8136-8158)
  - Processes movement interruption notifications (ZC_STOPMOVE)
  - Requires in-game state (changeToInGameState)
  - Creates coordinates hash from packet x,y values
  - Gets actor reference using Actor::get
  - Updates actor position data:
    * Sets pos and pos_to to new coordinates
    * Records movement time (time_move)
    * Resets time_move_calc to 0
  - Handles special actor types:
    * For players/character: Resets sitting state to 0
    * For character: Clears AI move queue and outputs debug message
    * For homunculus: Clears AI move queue
  - Contains detailed comments about packet types:
    * 0088: Basic packet version
    * 08CD: Newer packet version
  - Uses "parseMsg_move" debug category

- high_jump - Instant movement handler (lines 7134-7154)
  - Processes instant movement/teleport notifications
  - Requires in-game state (changeToInGameState)
  - Gets actor reference using Actor::get
  - Handles unknown actors:
    * Creates new Actor::Unknown
    * Sets appear_time and nameID
  - Detects failed movements:
    * Checks if destination matches current pos_to
    * Displays failure message
    * Returns early
  - Updates actor position:
    * Sets pos and pos_to to new coordinates
  - Displays movement message with actor name and coordinates
  - Updates movement timing variables:
    * Sets time_move to current time
    * Resets time_move_calc to 0
  - Uses "skill" message category

- character_moves - Character movement handler (lines 5496-5531)
  - Processes character movement notifications (ZC_NOTIFY_PLAYERMOVE)
  - Requires in-game state (changeToInGameState)
  - Updates character position data:
    * Sets pos_to coordinates from packet
    * Calculates movement distance
    * Records movement start time
  - Calculates movement parameters:
    * Uses character walk speed
    * Finds path solution
    * Computes movement time
  - Updates character direction:
    * Calculates vector between positions
    * Converts vector to degree
    * Sets body and head direction
  - Contains AI integration:
    * Handles mapRoute escape logic
    * Triggers escape AI if needed
    * Sends escape message if configured
  - Packet: 0087

- actor_action() - Handles various actor actions (lines 2588-2710)
  - Processes different action types:
    * ACTION_ITEMPICKUP: Item pickup handling
      - Tracks item pickup events
      - Updates item takenBy field
    * ACTION_SIT/ACTION_STAND: Sitting/standing state changes
      - Updates character state
      - Triggers sitAuto AI if needed
      - Handles both player and other actors
    * ACTION_ATTACK_*: Various attack types
      - Calculates and displays damage (including critical/miss)
      - Updates damage tables
      - Handles dual wield damage
      - Manages attack messages and status displays
      - Differentiates between player/other actors
      - Tracks damage taken
  - Calls Misc::checkValidity() for state validation
  - Uses attack_string() for message formatting
  - Calls Plugins::callHook('packet_attack')

- actor_display_compatibility() - Compatibility layer for actor display events (lines 1833-1839)
  - Provides backward compatibility for actor display hooks
  - Calls packet_pre/actor_display hook before processing
  - Delegates to actor_display() unless return flag set
  - Calls packet/actor_display hook after processing
  - Used by:
    * actor_exists
    * actor_connected
    * actor_moved
    * actor_spawned
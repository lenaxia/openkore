**Movement and Actor Display Handlers:**

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
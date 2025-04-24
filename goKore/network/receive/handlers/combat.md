# Combat Related Handlers

**Method Implementations:**
- monster_ranged_attack - Monster ranged attack handler (lines 11140-11163)
  - Processes monster ranged attack notifications
  - Gets monster ID and attack range
  - Creates coordinate hashes for source and target positions
  - Updates monster's movetoattack_pos and movetoattack_time if monster exists
  - Updates character's movetoattack_pos and movetoattack_time
  - Outputs debug message with coordinates and range
  - Triggers monster_ranged_attack hook with monster ID
  - Uses "parseMsg_move" debug category
  - Handles failed attack attempts due to range issues
- combo_delay - Combo delay handler (lines 10587-10597)
  - Processes combo delay notifications
  - Updates character's combo_packet with delay value
  - Contains comment questioning the formula: "How was the above formula derived?"
  - Contains comment suggesting moving manipulation to functions.pl
  - Gets actor using Actor::get
  - Determines appropriate verb form (have/has)
  - Outputs debug message with actor name, verb, and delay
  - Uses "parseMsg_comboDelay" debug category
- attack_range - Attack range update handler (lines 10219-10237)
  - Processes attack range update notifications
  - Gets range type from packet
  - Outputs debug message with range value
  - Returns early if not in game state
  - Updates character's attack_range
  - If attackDistanceAuto config is enabled:
    * Decreases attackDistance if it's greater than the new range
    * Displays success message with new attackDistance
    * Updates attackMaxDistance to match the new range
    * Displays success message with new attackMaxDistance
  - Uses "debug" message category for range notification
  - Uses "success" message category for config changes
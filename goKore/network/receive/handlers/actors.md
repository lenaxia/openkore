**Method Implementations:**

- actor_info() - Handles actor information updates (lines 2712-2804)
  - Processes different actor types:
    * Players:
      - Updates name, party/guild info, titles
      - Calls updatePlayerNameCache()
      - Triggers charNameUpdate hook
    * Monsters:
      - Updates monster name and info
      - Updates monster LUT if new
      - Triggers mobNameUpdate hook
    * NPCs:
      - Updates NPC name and info
      - Updates NPC LUT if new
      - Triggers npcNameUpdate hook
    * Pets:
      - Updates pet name and info
      - Triggers petNameUpdate hook
    * Slaves:
      - Updates slave name and info
      - Triggers slaveNameUpdate hook
    * Elementals:
      - Updates elemental name and info
      - Triggers elementalNameUpdate hook
  - Uses bytesToString() for string conversion
  - Calls Plugins::callHook() for various update types

- actor_look_at - Actor direction handler (lines 8121-8129)
  - Processes actor direction change notifications
  - Requires in-game state (changeToInGameState)
  - Gets actor reference using Actor::get
  - Updates actor's look properties:
    * head direction
    * body direction
  - Outputs debug message with actor name and direction values
  - Uses "parseMsg" debug category
  - Simple implementation focused on direction tracking
- offline_clone_lost - Offline player clone removal handler (lines 7994-8015)
  - Processes offline player clone removal notifications
  - Removes player from multiple lists:
    * From playersList:
      - Gets player reference by ID
      - Sets gone_time to current time
      - Creates deep copy in players_old hash
      - Triggers player_disappeared hook
      - Removes from playersList
    * From vender list:
      - Removes ID from venderListsID array
      - Deletes entry from venderLists hash
    * From buyer list:
      - Removes ID from buyerListsID array
      - Deletes entry from buyerLists hash
  - Comprehensive cleanup of player references
- offline_clone_found - Offline player clone handler (lines 7954-7992)
  - Processes offline player clone notifications
  - Gets actor reference from playersList by ID
  - Creates new Actor::Player if not found:
    * Sets object_type to 0x0 (player)
    * Sets clone flag to 1
    * Sets ID and nameID
    * Sets name (converted from bytes)
    * Sets appear_time to current time
    * Sets jobID and type
    * Sets position coordinates
    * Sets movement timing variables
    * Sets walk_speed to 1 (hack)
    * Sets level to 1
    * Sets appearance properties:
      - Robe, clothes_color
      - Headgear (low, mid, top)
      - Weapon, shield
      - Sex, hair_color
  - Adds actor to playersList
  - Triggers hooks:
    * add_player_list
    * player (backwards compatibility)
    * player_exist
- sprite_change - Actor appearance change handler (lines 4528-4574)
  - Processes appearance changes for players and characters
  - Handles multiple sprite types:
    * Type 0: Job change
    * Type 2: Weapon and shield
    * Type 3: Lower headgear
    * Type 4: Upper headgear
    * Type 5: Middle headgear
    * Type 6: Hair color
    * Type 9: Shoes
    * Type 12: Robe
    * Type 7/13: Body palette/color
  - Updates player appearance properties
  - Displays appropriate messages for each change
  - Triggers sprite_job_change hook
- monster_hp_info_tiny - Monster HP bar handler (lines 4086-4094)
  - Processes simplified monster HP information
  - Updates monster's HP percentage (in 5% increments)
  - Logs approximate HP percentage remaining
  - Packet: 0A36 (ZC_HP_INFO_TINY)
- monster_hp_info - Monster HP display handler (lines 4073-4082)
  - Processes monster HP information packets
  - Updates monster's current and maximum HP values
  - Calculates and logs HP percentage
  - Packet: 0977 (ZC_HP_INFO)
- monster_typechange - Monster type change handler (lines 4050-4069)
  - Processes monster type/sprite changes
  - Updates monster name based on monster lookup table
  - Resets damage tracking statistics
  - Displays notification message about the change
  - Packet: 01B0 (ZC_NPCSPRITE_CHANGE)
- received_characters_blockSize() - Character data block size handler (lines 695-708)
  - Determines the block size for character data packets
  - Can be overridden by server-specific configurations (charBlockSize)
  - Defaults to 155 bytes (standard for kRO and most official/emulator servers)
  - Used when parsing character selection/creation packets
  - Last updated: 2020-11-13

- received_characters_unpackString() - Character data format handler (lines 711-719)
  - Defines unpack formats for character data packets
  - Supports different versions:
    - 175 bytes (PACKETVER >= 20201007): handles uint64 HP/SP fields
    - 155 bytes (PACKETVER >= 20170830): handles uint64 exp fields
  - Unpacks key character attributes:
    - Basic info (charID, name, job, stats)
    - Appearance (hair style, colors)
    - Status (HP, SP, exp)
    - Equipment and position

**Actor Display Handlers:**

- actor_display - Main actor display handler (lines 1833-2397)
  - Handles all actor types:
    - Players
    - Monsters
    - NPCs
    - Portals
    - Pets
    - Homunculus/Mercenary
    - Elementals
  - Processes different actor states:
    - Exists (standing)
    - Connected (new)
    - Moved
    - Spawned
  - Manages actor properties:
    - Position and movement
    - Visual appearance
    - Status effects
    - Equipment
    - Guild information
  - Maintains actor lists:
    - playersList
    - monstersList
    - npcsList
    - portalsList
    - petsList
    - slavesList
    - elementalsList
  - Handles special cases:
    - Off-map coordinates
    - Client sight distance
    - Guild flags
    - Visual effects

**Death and Disappearance Handlers:**
- actor_died_or_disappeared() - Handles all entity death/disappearance events (lines 2407-2586)
  - Processes different types of disappearances:
    - 0 = out of sight
    - 1 = died  
    - 2 = logged out
    - 3 = teleport
    - 4 = trickdead
  - Handles special cases for:
    - Player character death (updates death count, triggers hooks)
    - Monster deaths/disappearances (handles loot taking)
    - Player disconnections/teleports
    - Portal/NPC/Pet/Slave/Elemental disappearances
    - Unknown entities (logs debug info)
  - Maintains old actor lists for disappeared entities
  - Calls appropriate plugin hooks for each entity type



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

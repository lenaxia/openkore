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

**Social Interaction Handlers:**

- marriage_partner_name() - Marriage partner name handler (lines 4112-4116)
  - Processes marriage partner name packets
  - Converts byte data to string using bytesToString()
  - Displays partner name in message log
  - Triggered before "I miss you" skill is cast

- married - Marriage notification handler (lines 7004-7009)
  - Processes marriage event notifications
  - Gets actor reference using Actor::get
  - Displays "[Actor] got married!" message
  - Simple implementation focused on notification


- ignore_player_result - Individual ignore result handler (lines 6946-6955)
  - Processes results of ignoring individual players
  - Handles two operation types:
    * 0 = Enable ignore: Displays "Player ignored" message
    * 1 = Disable ignore:
      - Checks for error code 0 (success)
      - Displays "Player unignored" message
  - Contains TODO comment about storing list of ignored players
  - Simple implementation focused on user feedback

- ignore_all_result - Global ignore result handler (lines 6933-6943)
  - Processes results of ignore all players commands
  - Handles two operation types:
    * 0 = Enable ignore all:
      - Sets ignored_all flag to 1
      - Displays "All Players ignored" message
    * 1 = Disable ignore all:
      - Checks for error code 0 (success)
      - Displays "All players unignored" message
  - Contains TODO comment about storing this state
  - Simple implementation focused on user feedback

- friend_response - Friend request response handler (lines 6235-6258)
  - Processes friend request response notifications (ZC_ADD_FRIENDS_LIST)
  - Extracts response information:
    * Response type
    * Friend name (converted from bytes to string)
  - Handles different response types:
    * 0 = Success:
      - Adds new friend to friendsID array
      - Extracts accountID and charID from raw message
      - Sets name and online status (1)
      - Displays "You have become friends with" message
    * 1 = Rejection: "[Player] does not want to be friends"
    * 2 = Your list full: "Your Friend List is full"
    * 3 = Their list full: "[Player]'s Friend List is full"
    * Other: "[Player] rejected to be your friend"
  - Result codes (from packet comments):
    * 0 = "You have become friends with (%s)."
    * 1 = "(%s) does not want to be friends with you."
    * 2 = "Your Friend List is full."
    * 3 = "(%s)'s Friend List is full."

- friend_removed - Friend removal handler (lines 6212-6226)
  - Processes friend removal notifications (PACKET_ZC_DELETE_FRIENDS)
  - Extracts friend information:
    * Account ID
    * Character ID
  - Searches through friendsID array for matching friend
  - When found:
    * Displays "[Friend] is no longer your friend" message
    * Removes friend ID from friendsID array
    * Deletes friend entry from friends hash
  - Simple implementation focused on friend removal

- friend_request - Friend request handler (lines 6194-6208)
  - Processes incoming friend requests (ZC_REQ_ADD_FRIENDS)
  - Stores request information in incomingFriend hash:
    * Account ID
    * Character ID
    * Name (converted from bytes to string)
  - Displays notification messages:
    * "[Player] wants to be your friend"
    * Instructions for accepting/rejecting request
  - Triggers friend_request hook with:
    * Account ID
    * Character ID
    * Name
  - Simple implementation focused on request notification

- friend_logon - Friend online status handler (lines 6171-6190)
  - Processes friend online/offline notifications (ZC_FRIENDS_STATE)
  - Extracts friend information:
    * Account ID
    * Character ID
    * Online status (isNotOnline flag)
  - Searches through friendsID array for matching friend
  - Updates online status in friends hash (1-isNotOnline)
  - Displays appropriate message:
    * "[Friend] has disconnected" when isNotOnline=1
    * "[Friend] has connected" when isNotOnline=0
  - State values (from packet comments):
    * 0 = online
    * 1 = offline

- friend_list - Friend list handler (lines 6143-6163)
  - Processes friend list packets (ZC_FRIENDS_LIST)
  - Resets existing friend data structures:
    * Clears friendsID array
    * Clears friends hash
  - Parses raw message data in 32-byte chunks
  - For each friend entry:
    * Adds ID to friendsID array
    * Extracts accountID, charID, and name
    * Converts name bytes to string
    * Initializes online status to 0 (offline)
    * Increments ID counter
  - Simple implementation focused on data structure setup

- emoticon - Emoticon display handler (lines 5961-6024)
  - Processes emoticon/emotion display packets
  - Looks up emotion display text from emotions_lut
  - Handles different actor types:
    * Current character:
      - Displays simple message with name and emotion
    * Other players:
      - Calculates distance to player
      - Displays message with distance, name, and emotion
      - Handles follow emotion response if configured
      - Mirrors specific emotions (30↔31) when following
    * Monsters/slaves:
      - Calculates distance
      - Displays message with actor type, name, and emotion
    * Other/unknown actors:
      - Uses generic display with nameIdx
  - Logs emotions to chat log if configured
  - Triggers packet_emotion hook with emotion and ID

- show_eq - Equipment display handler (lines 3216-3281)
 - Handles multiple packet versions for equipment display:
   - 02D7: Default packet version
   - 0906: Unimplemented on eAthena
   - 0859: Added in 20101124
   - 0997: Added in 20120925
   - 0A2D: Added in 20150226
   - 0B03: Added in 20150226 (alternative)
 - Parses equipment info with different formats per version
 - Supports robe equipment (PACKETVER >= 20100629)
 - Formats and displays equipment info with:
   - Centered title with character name
   - List of equipment by slot
   - Proper item naming and identification
 - Uses internationalized strings (T())
 - Outputs to 'list' message channel

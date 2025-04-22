**Method Implementations:**
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
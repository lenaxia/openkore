# Party Related Handlers

**Method Implementations:**
- partylv_info - Party member level info handler (lines 9428-9435)
  - Processes party member level information updates
  - Gets member ID from packet
  - Checks if member exists in party users hash
  - Updates party member data if found:
    * job: Job ID
    * lv: Level
  - Simple implementation focused on data updates
- party_dead - Party member death notification handler (lines 8472-8488)
  - Processes party member death/revival notifications
  - Gets member name:
    * From party users hash if available
    * Falls back to ID if not found
  - Handles two states based on isDead flag:
    * 1 (dead): Sets dead flag and records death time
    * 0 (alive): Clears dead and dead_time flags
  - Displays appropriate messages:
    * "Party member X is dead"
    * "Party member X is alive"
  - Uses "info" message category
  - Packet: 0AB2 (includes GID and dead flag)
- party_users_info - Party member information handler (lines 8415-8468)
  - Processes party member information packets
  - Requires in-game state (changeToInGameState)
  - Supports multiple packet versions:
    * 0A44: PACKETVER >= 20151007 (50 bytes per member)
    * 0AE5: PACKETVER >= 20171207 (54 bytes per member)
    * 00FB: Default/old version (46 bytes per member)
  - Defines different data structures per version:
    * Basic: ID, name, map, admin, online
    * Extended: Adds jobID, lv
    * 20171207+: Adds GID field
  - Sets party name from packet
  - Processes each member's information:
    * Skips last 6 bytes (item rules in 0a43)
    * Adds ID to partyUsersID array if not present
    * Creates new Actor::Party object
    * Unpacks member data using version-specific format
    * Converts name from bytes to string
    * Inverts admin and online flags (0=true, 1=false)
    * Clears dead status flags
  - Outputs debug message with member name and map
  - Triggers party_users_info_ready hook
  - Contains comment about server behavior with saveMap
- party_show_picker - Party item pickup notification handler (lines 8399-8413)
  - Processes party member item pickup notifications
  - Returns early if sourceID matches accountID (own character)
  - Gets member name:
    * From party users hash if available
    * Falls back to sourceID if not found
  - Creates item object with properties:
    * nameID: Item ID
    * identified: Identification status
    * upgrade: Upgrade/refine level
    * cards: Card information
    * broken: Broken status
  - Displays message with member name and item name
  - Uses "info" message category
  - Contains comment about server sending this packet for own character (rRo)
- party_organize_result - Party creation result handler (lines 8383-8397)
  - Processes party creation result notifications
  - Handles success case (fail=0):
    * Sets admin flag for character in party users hash
  - Handles multiple failure cases:
    * 1: Party name already exists
    * 2: Already in a party
    * 3: Not allowed in current map
    * Other: Unknown error (shows code)
  - Displays appropriate warning messages for each failure case
  - Uses "warning" message category for failures
- party_location - Party member location update handler (lines 8371-8382)
  - Processes party member location updates
  - Gets member ID from packet
  - Updates party member location if member exists:
    * Sets x and y coordinates
    * Sets online status to 1
  - Outputs debug message with member name and coordinates
  - Uses "parseMsg" debug category
  - Simple implementation focused on position tracking
- party_leave - Party member leave handler (lines 8348-8369)
  - Processes party member leave notifications
  - Gets member ID from packet
  - Gets actor reference from party users hash
  - Removes member from party data structures:
    * Deletes from char->{party}{users} hash
    * Removes from partyUsersID array
  - Special handling for own character:
    * Sets actor reference to $char
    * Deletes entire party structure
    * Clears partyUsersID array
    * Sets party joined flag to 0
  - Handles different leave reasons:
    * GROUPMEMBER_DELETE_LEAVE: Normal leave
    * GROUPMEMBER_DELETE_EXPEL: Kicked from party
    * Other: Unknown reason (shows code)
  - Displays appropriate message for each leave reason
- party_invite_result - Party invitation response handler (lines 8324-8346)
  - Processes party invitation response notifications
  - Converts invitee name from bytes to string
  - Handles multiple response types:
    * ANSWER_ALREADY_OTHERGROUPM: Target already in another party
    * ANSWER_JOIN_REFUSE: Target denied request
    * ANSWER_JOIN_ACCEPT: Target accepted request
    * ANSWER_MEMBER_OVERSIZE: Party is full
    * ANSWER_DUPLICATE: Same account already in party
    * ANSWER_JOINMSG_REFUSE: Join message refused
    * ANSWER_UNKNOWN_ERROR: Unknown error
    * ANSWER_UNKNOWN_CHARACTER: Character offline or doesn't exist
    * ANSWER_INVALID_MAPPROPERTY: Invalid map property
  - Displays appropriate messages for each response type
  - Uses "warning" category for failures, "info" for other responses
- party_invite - Party invitation handler (lines 8311-8322)
  - Processes incoming party invitations
  - Converts party name from bytes to string
  - Displays invitation message with party name
  - Stores invitation information:
    * Party ID
    * ACK packet type (02C7 or 00FF based on request packet)
  - Sets up auto-deny timeout
  - Triggers party_invite hook with:
    * partyID: Party identifier
    * partyName: Party name
  - Supports multiple packet versions:
    * 02C6: Newer packet version
    * Default: Older packet version
- party_hp_info - Party member HP update handler (lines 8301-8309)
  - Processes party member HP updates
  - Extracts member ID from packet
  - Updates party member HP information if member exists:
    * Current HP value
    * Maximum HP value
  - Simple implementation focused on data updates
  - No message display or hooks
- party_leader - Party leadership change handler (lines 8288-8299)
  - Processes party leadership change notifications
  - Iterates through partyUsersID array
  - Updates admin status for affected members:
    * Removes admin flag from old leader
    * Sets admin flag for new leader
  - Displays message with new leader's name
  - Uses "party" message category with priority 1
- party_exp - Party experience settings handler (lines 8258-8286)
  - Processes party experience sharing settings
  - Updates character's party share setting
  - Handles EXP distribution types:
    * 0: Individual Take - Each member gets own exp
    * 1: Even Share - Experience divided evenly
  - Optionally processes item settings if present:
    * itemPickup: Item pickup policy
    * itemDivision: Item division policy
  - Displays appropriate messages for each setting:
    * EXP distribution method
    * Item pickup method
    * Item division method
  - Shows error message for invalid settings
  - Uses "party" message category with priority 1
- party_chat - Party chat message handler (lines 8233-8256)
  - Processes party chat messages
  - Converts message bytes to string
  - Parses message format:
    * Extracts username and message content
    * Removes trailing space from username
  - Processes message content:
    * Strips language code
    * Solves message using solveMessage function
  - Displays formatted message with [Party] prefix
  - Handles logging:
    * Writes to party chat log if logPartyChat is enabled
    * Adds to ChatQueue with 'p' type
    * Outputs debug message with "partychat" category
  - Triggers packet_partyMsg hook with:
    * MsgUser: Sender username
    * Msg: Parsed message content
    * RawMsg: Original message
  - Uses "partychat" message category
- party_allow_invite - Party invitation permission handler (lines 8223-8231)
  - Processes party invitation permission notifications
  - Handles two states based on type flag:
    * type=1: Displays "Not allowed other player invite to Party"
    * type=0: Displays "Allowed other player invite to Party"
  - Simple implementation focused on notification
  - Contains TODO comment about storing this state
  - Uses "party" message category with priority 1
- party_join - Party join handler (lines 8168-8220)
  - Processes party join notifications
  - Requires in-game state (changeToInGameState)
  - Supports multiple packet versions:
    * 0104: Default old packet
    * 01E9: PACKETVER >= 2015
    * 0A43: PACKETVER >= 2016
    * 0AE4: PACKETVER >= 2017
  - Handles different data fields per version:
    * Basic: ID, role, x, y, type, name, user, map
    * Extended: Adds lv, item_pickup, item_share
    * 2016+: Adds jobID
    * 2017+: Adds charID
  - Processes party member joining:
    * Adds ID to partyUsersID array if not present
    * Handles self-join vs. other member join
    * Displays appropriate join messages
  - Updates party member data:
    * Creates new Actor::Party if needed
    * Sets admin status, online state, position
    * Updates name, ID, level, job information
  - Updates party settings:
    * Party name
    * Item pickup policy
    * Item division policy
  - Triggers hooks:
    * packet_partyJoin (for self-join)
  - Uses message category 1 (undef)
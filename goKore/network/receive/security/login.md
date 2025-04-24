**Login Handlers:**

**secure_login_key** - Secure login key handler (lines 11377-11381)
  - Processes secure login key notifications
  - Stores secure key in secureLoginKey variable
  - Outputs debug message with hexadecimal representation of key
  - Uses "connection" debug category
  - Simple implementation focused on key storage
  
- **received_login_token** - Login token handler (lines 8022-8029)
  - Processes login token notifications
  - Skips processing for XKore mode 1 (version == 1)
  - Gets master server configuration
  - Sends token to server with:
    - Username and password
    - Master version and client version
    - Login token and length
    - OTP IP and port
  - Comment notes that rathena uses 0064 packet instead of 0825
  - Simple implementation focused on token forwarding
  
- **character_ban_list** - Character ban list handler (lines 7945-7948)
  - Processes character ban list notifications
  - Contains only comment about packet structure:
    - Header + Len + CharList[character_name(size:24)]
  - Empty implementation (stub function)
  - No actual processing or message display
  
- **login_error_game_login_server** - Character server login error handler (lines 5454-5456)
  - Handles login errors specific to character server
  - Displays error message for invalid character selection
  - Resets network state to 1 (disconnected)
  - Simple implementation with minimal error handling
  
- **login_error** - Login error handler (lines 5356-5452)
  - Processes login error responses from server
  - Handles multiple error types:
    - REFUSE_INVALID_ID/REFUSE_INVALID_ID2: Account doesn't exist
    - REFUSE_INVALID_PASSWD/REFUSE_INVALID_PASSWD2: Password error
    - REFUSE_BAN_BY_GM/REFUSE_NOT_CONFIRMED: Account blocked
    - REFUSE_INVALID_VERSION: Version mismatch
    - REFUSE_BLOCK_TEMPORARY: Temporary connection block
    - REFUSE_USER_PHONE_BLOCK: Phone lock
    - REFUSE_EMAIL_NOT_CONFIRMED: Email not confirmed
    - REFUSE_BLOCKED_ID: User blocked
    - REFUSE_BLOCKED_COUNTRY: Country blocked
    - REFUSE_BILLING: Billing issues
    - REFUSE_CHANGE_PASSWD_FORCE2: Password change required
    - REFUSE_ACCOUNT_NOT_PREMIUM: Premium server access denied
    - REFUSE_NOT_ALLOWED_IP_ON_TESTING: Connection delayed
    - REFUSE_TOKEN_EXPIRED: Token expired
  - Manages reconnection attempts:
    - Prompts for username/password re-entry
    - Resets connection timeouts
    - Handles offline mode transitions
  - Triggers hooks:
    - invalid_password
    - dial
  
- **received_characters_slots_info** - Character slot information handler (lines 811-836)
  - Manages character slot configuration during login
  - Handles different slot types:
    - Normal slots
    - Premium slots 
    - Billing slots
  - Tracks slot counts and validity
  - Updates connection state to CONNECTED_TO_LOGIN_SERVER
  - Stores character server information
  - Calls received_characters if character info is present
  - Parameters:
    - total_slot: Total available character slots
    - premium_start_slot: First premium slot index
    - premium_end_slot: Last premium slot index  
    - normal_slot: Normal slot count
    - premium_slot: Premium slot count
    - billing_slot: Billing slot count
    - producible_slot: Slots available for character creation
    - valid_slot: Valid/usable slots
    - options: Additional connection options
    - charInfo: Character information if present
  
- **received_characters** - Character data handler (lines 841-915)
  - Processes character information packets during login
  - Manages character objects and slots
  - Handles character reuse/creation logic:
    - Reuses existing character if IDs match
    - Creates new Actor::You if needed
  - Updates character attributes:
    - Name, ID, levels, headgear
    - Last map, gender, delete date
  - Maintains character slot array ($chars)
  - Handles login state transitions
  - Parameters:
    - charInfo: Packed character data
    - Uses masterServer->{charBlockSize} for parsing
  - Related methods:
    - received_characters_blockSize()
    - received_characters_unpackString()
  - Post-processing:
    - Sends ban check if needed
    - Handles PIN code requests
    - Manages login pause timeouts
  
- **sync_received_characters** - Character synchronization handler (lines 920-930)
  - Handles character page synchronization during login
  - Manages sync count and state in $charSvrSet
  - Tracks number of pages in character selection screen
  - Parameters:
    - sync_Count: Total pages to sync
  - Behavior:
    - Resets sync_received_characters counter when new sync_Count received
    - Sends sync packets to server if client is alive
    - Only works with DirectConnection network type
  - Related packets:
    - PACKET_HC_CHARLIST_NOTIFY (0x09A0)
  
- **reconstruct_received_characters** - Character data reconstruction handler (lines 935-940)
  - Reconstructs packed character data for network transmission
  - Uses received_characters_unpackString helper method
  - Parameters:
    - chars: Array of character data structures
    - charInfo: Output packed character data
  - Packing format:
    - Uses masterServer->{charBlockSize} for packing
    - Processes each character in chars array
    - Maps character attributes using unpackString keys/types
  - Related methods:
    - received_characters_unpackString()
    - reconstruct_received_characters_info()
  
- **reconstruct_received_characters_info** - Character info reconstruction handler (lines 942-947)
  - Reconstructs packed character info for network transmission
  - Similar to reconstruct_received_characters but for info packets
  - Parameters:
   - chars: Array of character data structures
   - charInfo: Output packed character info
  - Packing format:
   - Uses masterServer->{charBlockSize} for packing
   - Processes each character in chars array
   - Maps character attributes using unpackString keys/types
  - Related packets:
   - PACKET_HC_ACCEPT_MAKECHAR_NEO_UNION (0x006E)
  
- **character_creation_successful** - Character creation handler (lines 952-965)
  - Handles successful character creation notification
  - Creates new Actor::You instance for the character
  - Parameters:
   - charInfo: Packed character data
  - Initialization:
   - Unpacks character attributes using received_characters_unpackString
   - Sets character ID to accountID
   - Initializes headgear (top/bottom)
   - Tracks initial job/base levels for exp calculation
  - Related methods:
   - received_characters_unpackString()
  - Related packets:
   - PACKET_HC_ACCEPT_MAKECHAR_NEO_UNION (0x006E)
  
- **character_creation_failed** - Character creation failure handler (lines 999-1021)
  - Handles failed character creation attempts
  - Processes different failure reasons:
   - 0x00: Character name already exists
   - 0xFF: General creation denied
   - 0x01: Underage restriction
   - 0x02: Invalid symbols in name
   - 0x03: Slot eligibility issue
  - Displays appropriate error messages to user
  - Manages state transition:
   - Returns to character selection screen
   - Updates network state if needed
   - Maintains first login map flag
   - Preserves starting zeny amount
  - Parameters:
   - flag: Failure reason code
  - Related methods:
   - charSelectScreen()
  - Related packets:
   - PACKET_HC_REFUSE_MAKECHAR (0x006F)
  
- **received_characters_info** - Character info timeout handler (lines 1027-1044)
  - Manages character info timeout during login
  - Sets up timeout hooks for character selection screen:
   - 6 second timeout for server connection
   - Cleans up timeout hooks when character selected
  - Calls received_characters_slots_info with provided args
  - Tracks character login time
  - Parameters:
   - args: Character info arguments passed to received_characters_slots_info
  - Related methods:
   - received_characters_slots_info()
  - Related packets:
   - PACKET_HC_ACCEPT_ENTER_NEO_UNION (0x006B)
   - PACKET_HC_ACCEPT2 (0x082D)

- received_characters_blockSize - Character data block size handler (lines 695-708)
 - Determines the block size for character data packets
 - Can be overridden by server-specific configurations (charBlockSize)
 - Defaults to 155 bytes (standard for kRO and most official/emulator servers)
 - Used when parsing character selection/creation packets
 - Last updated: 2020-11-13

- received_characters_unpackString - Character data format handler (lines 711-719)
 - Defines unpack formats for character data packets
 - Supports different versions:
   - 175 bytes (PACKETVER >= 20201007): handles uint64 HP/SP fields
   - 155 bytes (PACKETVER >= 20170830): handles uint64 exp fields
 - Unpacks key character attributes:
   - Basic info (charID, name, job, stats)
   - Appearance (hair style, colors)
   - Status (HP, SP, exp)
   - Equipment and position

- switch_character - Character switching handler (lines 11593-11603)
 - Processes character switching notifications
 - Sets network state to CONNECTED_TO_MASTER_SERVER
 - Disconnects from server
 - Clears character list (@chars)
 - Outputs debug message with result value
 - Contains comments about:
   * X-Kore character switching
   * Packet format (00B3)
   * Type values (1 = disconnect, char-select)
   * TODO for adding real client messages and logic
   * ClientLogic: LoginStartMode = 5; ShowLoginScreen
 - Contains FIXME comment about better support for multiple received_characters packets

- received_character_ID_and_Map - Character and map info handler (lines 8055-8113)
 - Processes character ID and map server information
 - Sets network state to 4 and clears connection retry counter
 - Stores character ID from packet
 - Handles XKore version 1 master server setup
 - Processes map information:
   * Extracts map name (removes .gat extension)
   * Cleans up instance ID from map name
   * Creates new Field object if needed
   * Handles field loading errors
 - Extracts map server connection details:
   * From mapUrl/mapPort if provided in URL format
   * Otherwise from mapIP/mapPort with makeIP conversion
   * Uses masterServer IP for private servers
 - Implements XKore mode 1 workaround:
   * Finds matching character by ID in chars array
   * Updates char configuration
   * Sets $char reference
 - Displays detailed game connection information:
   * Character ID (hex and decimal)
   * Map name
   * Map IP and port
 - Checks if map is allowed
 - Disconnects from character server
 - Initializes stat variables
 - Uses "connection" message category

- character_name - Character name handler (lines 5533-5549)
 - Processes character name information packets
 - Converts name bytes to string
 - Updates guild member name if applicable:
   * Searches guild member list by character ID
   * Updates name field when found
 - Logs debug message with received name
 - Simple implementation focused on name updates

- character_deletion_failed - Character deletion failure handler (lines 5483-5492)
 - Processes failed character deletion notifications
 - Displays error message about incorrect email address
 - Clears deletion index reference
 - Manages state transition:
   * Returns to character selection screen
   * Updates network state if needed
   * Sets first login map flag
   * Preserves starting zeny amount
   * Sets welcome message flag
 - Similar state handling to character_deletion_successful

- character_deletion_successful - Character deletion success handler (lines 5463-5481)
 - Processes successful character deletion notifications
 - Handles two scenarios:
   * Known deletion index: Shows character name and index
   * Unknown deletion index: Shows generic message
 - Cleans up character data:
   * Removes character from chars array
   * Clears deletion index reference
   * Removes empty character slots
 - Manages state transition:
   * Returns to character selection screen
   * Updates network state if needed
   * Sets first login map flag
   * Preserves starting zeny amount
   * Sets welcome message flag

- char_delete2_result - Character deletion result handler (lines 3562-3583)
 - Processes character deletion result packets
 - Handles successful deletion case:
   - Sets deletion date for character
   - Displays time remaining until deletion
 - Handles various error cases:
   - Character already planned for deletion (error 0)
   - Database error (error 3)
   - Guild membership prevents deletion (error 4)
   - Party membership prevents deletion (error 5)
   - Unknown errors (default case)
 - Returns to character select screen after processing
 - Uses internationalized strings (T()) for messages

- char_delete2_accept_result - Character deletion acceptance handler (lines 3586-3633)
 - Processes character deletion acceptance packets
 - Handles successful deletion:
   - Removes character data from memory
   - Cleans up empty character slots
   - Returns to character select screen
 - Handles various error cases:
   - Birthday required (error 0)
   - System settings prevent deletion (error 2)
   - Database error (error 3)
   - Temporary deletion block (error 4)
   - Birthday mismatch (error 5)
   - Incorrect email (error 7)
   - Unknown errors (default case)
 - Always returns to character select screen
 - Uses internationalized strings (T()) for messages

- char_delete2_cancel_result - Character deletion cancellation handler (lines 3636-3650)
 - Processes character deletion cancellation packets
 - Handles successful cancellation:
   - Clears deletion date
   - Shows confirmation message
 - Handles error cases:
   - Database error (error 2)
   - Unknown errors (default case)
 - Returns to character select screen
 - Uses internationalized strings (T()) for messages

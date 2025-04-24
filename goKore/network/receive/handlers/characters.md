# Character Related Handlers

**Method Implementations:**
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
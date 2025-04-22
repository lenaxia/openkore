**Login Handlers:**

**received_characters_slots_info** - Character slot information handler (lines 811-836)
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

**received_characters** - Character data handler (lines 841-915)
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

**sync_received_characters** - Character synchronization handler (lines 920-930)
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

**reconstruct_received_characters** - Character data reconstruction handler (lines 935-940)
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

**reconstruct_received_characters_info** - Character info reconstruction handler (lines 942-947)
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

**character_creation_successful** - Character creation handler (lines 952-965)
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

**character_creation_failed** - Character creation failure handler (lines 999-1021)
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

**received_characters_info** - Character info timeout handler (lines 1027-1044)
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
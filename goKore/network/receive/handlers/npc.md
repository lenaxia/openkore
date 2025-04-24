**Method Implementations:**

- dynamicnpc_create_result - Dynamic NPC creation result handler (lines 12414-12432)
  - Processes dynamic NPC creation result notifications (0A17)
  - Translates result codes to readable status messages:
    * DYNAMICNPC_RESULT_SUCCESS: "Success"
    * DYNAMICNPC_RESULT_UNKNOWN: "Unknown"
    * DYNAMICNPC_RESULT_UNKNOWNNPC: "Unknown NPC"
    * DYNAMICNPC_RESULT_DUPLICATE: "Duplicate"
    * DYNAMICNPC_RESULT_OUTOFTIME: "Out of time"
  - Displays status message
  - Simple implementation focused on notification
  - Packet: PACKET_ZC_DYNAMICNPC_CREATE_RESULT
- npc_clear_dialog - NPC dialog clear handler (lines 7665-7669)
  - Processes NPC dialog clear notifications
  - Extracts NPC ID from packet
  - Outputs debug message about dialog closure
  - Simple implementation focused on notification
  - Uses "parseMsg" debug category
- npc_talk_text - NPC text input handler (lines 7561-7569)
  - Processes NPC text input dialog notifications (ZC_OPEN_EDITDLGSTR)
  - Extracts NPC ID from packet
  - Gets NPC name using getNPCName
  - Updates AI variables:
    * Sets npc_talk state to 'text'
    * Records timestamp
  - Simple implementation focused on state tracking
  - Packet: 01D4
- npc_talk_responses - NPC dialog menu handler (lines 7501-7557)
  - Processes NPC dialog menu notifications (ZC_MENU_LIST)
  - Extracts NPC ID and nameID from raw message
  - Auto-creates Task::TalkNPC if not active (similar to npc_talk)
  - Stores dialog information:
    * Sets talk{ID} and talk{nameID}
    * Unpacks and converts response string
  - Processes response options:
    * Splits responses by colon separator
    * Removes RO color codes
    * Handles special item references (^nItemID^)
    * Filters empty responses
    * Adds "Cancel Chat" option at the end
  - Updates AI variables:
    * Sets npc_talk state to 'select'
    * Records timestamp
  - Automatically runs 'talk resp' command
  - Gets NPC name using getNPCName
  - Triggers npc_talk_responses hook with:
    * ID: NPC ID
    * name: NPC name
    * responses: Response options array
  - Packet: 00B7
- npc_talk_number - NPC number input handler (lines 7489-7497)
  - Processes NPC number input dialog notifications (ZC_OPEN_EDITDLG)
  - Extracts NPC ID from packet
  - Gets NPC name using getNPCName
  - Updates AI variables:
    * Sets npc_talk state to 'number'
    * Records timestamp
  - Simple implementation focused on state tracking
  - Packet: 0142
- npc_talk_continue - NPC dialog continue handler (lines 7478-7485)
  - Processes NPC dialog continue notifications (ZC_WAIT_DIALOG)
  - Extracts NPC ID from raw message
  - Gets NPC name using getNPCName
  - Updates AI variables:
    * Sets npc_talk state to 'next'
    * Records timestamp
  - Simple implementation focused on state tracking
  - Packet: 00B5
- npc_talk_close - NPC dialog close handler (lines 7455-7474)
  - Processes NPC dialog close notifications (ZC_CLOSE_DIALOG)
  - Validates dialog context:
    * Checks if npc_talk ID matches current ID
    * Ignores unexpected close events
    * Skips processing for buy_or_sell state
  - Gets NPC name using getNPCName
  - Updates AI variables:
    * Sets npc_talk state to 'close'
    * Records timestamp
  - Clears talk hash
  - Triggers npc_talk_done hook with NPC ID
  - Packet: 00B6
- npc_talk - NPC dialog handler (lines 7410-7451)
  - Processes NPC dialog messages (ZC_SAY_DIALOG)
  - Auto-creates Task::TalkNPC if not active:
    * Checks if NPC or route+TalkNPC task is running
    * Creates and activates new TalkNPC task
    * Triggers npc_autotalk hook
  - Stores dialog information:
    * Sets talk{ID} and talk{nameID}
    * Converts message bytes to string
    * Removes RO color codes (^[a-fA-F0-9]{6})
    * Appends to existing conversation
  - Updates AI variables:
    * Sets npc_talk state to 'initiated'
    * Records timestamp
  - Gets NPC name using getNPCName
  - Triggers npc_talk hook with:
    * ID: NPC ID
    * nameID: Numeric NPC ID
    * name: NPC name
    * msg: Dialog message
  - Displays formatted message with NPC name
  - Uses "npc" message category
  - Packet: 00B4
- npc_chat - NPC chat display handler (lines 4957-4997)
  - Processes NPC chat messages
  - Handles message formatting:
    * Extracts name from "Name : Message" format
    * Calculates distance to NPC
    * Formats position information
  - Logs chat to npc.txt if logChat is enabled
  - Displays message with distance information
  - Triggers npc_chat hook with:
    * actor: Actor reference
    * ID: NPC ID
    * message: Formatted message
- parse_npc_chat - NPC chat parser (lines 4951-4955)
  - Pre-processes NPC chat packets
  - Retrieves actor reference using Actor::get
  - Sets actor field in args hash for use by npc_chat handler
  - Packet: 02C1
- parse_npc_image - NPC image parser (lines 3102-3106)
  - Converts NPC image bytes to string
  - Handles raw image data from packets

- reconstruct_npc_image - NPC image reconstructor (lines 3108-3112)
  - Converts NPC image string back to bytes
  - Prepares image data for sending

- npc_image - NPC image handler (lines 3117-3133)
  - Handles ZC_SHOW_IMAGE and ZC_SHOW_IMAGE2 packets
  - Displays/hides NPC illustrations
  - Supports different image types (type=2 for show, 255 for hide)
  - Manages talk{image} state
  - Logs image operations with debug messages

- show_script - NPC script message handler (lines 3375-3387)
  - Handles script/show messages from NPCs
  - Processes message content and NPC ID
  - Uses bytesToString for message conversion
  - Looks up NPC by ID in npcsList
  - Outputs debug message with NPC name and message
  - Calls plugin hook 'show_script' with:
    - ID: NPC ID
    - message: Decoded message content
  - Handles packet type:
    - 08B3: NPC script/show packet

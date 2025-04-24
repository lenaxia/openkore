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

- npc_market_purchase_result - Market purchase result handler (lines 7763-7781)
  - Processes NPC market purchase results (PACKET_ZC_NPC_MARKET_OPEN)
  - Outputs debug message with result code
  - Handles multiple result codes:
    * MARKET_BUY_RESULT_ERROR (-1): General error
    * MARKET_BUY_RESULT_SUCCESS (0): Success
    * MARKET_BUY_RESULT_NO_ZENY (1): Insufficient zeny
    * MARKET_BUY_RESULT_OVER_WEIGHT (2): Overweight
    * MARKET_BUY_RESULT_OUT_OF_SPACE (3): No inventory space
    * MARKET_BUY_RESULT_AMOUNT_TOO_BIG (4): Amount exceeds available stock
    * Other: Unknown error with code
  - Displays appropriate success/error messages
  - Uses "info" message category for all messages
  - Packet: 09D7

- npc_market_info - NPC market item list handler (lines 7710-7751)
  - Processes NPC market shop item listings (PACKET_ZC_NPC_MARKET_OPEN)
  - Uses server-specific pack format (npc_market_info_pack) or default 'v C V2 v'
  - Clears storeList and talk hash
  - Processes each item in the itemList:
    * Creates new Actor::Item
    * Unpacks item data (nameID, type, price, amount, weight)
    * Skips items with zero amount (client behavior)
    * Sets item ID to storeList size (workaround for duplicate items)
    * Sets item name using itemName function
    * Adds item to storeList
  - Returns early if no items in store
  - Automatically runs 'store' command if not in buyAuto mode
  - Sets in_market flag to 1
  - Updates AI variables:
    * Sets npc_talk state to 'store'
    * Records timestamp
  - Outputs debug messages for each item added
  - Packets: 09D5 (two versions with different nameID field sizes)

- npc_sell_list - NPC sellable items handler (lines 7637-7663)
  - Processes list of items that can be sold to NPC (ZC_PC_SELL_ITEMLIST)
  - Checks if message contains item data
  - Processes each item in itemsdata:
    * Unpacks index, price, and overcharge price
    * Gets item reference from inventory by ID
    * Flags item as sellable
    * Displays item amount, name, and price
  - Flags all other inventory items as unsellable:
    * Skips equipped items
    * Skips already flagged sellable items
    * Sets unsellable flag
  - Clears talk hash
  - Displays "Ready to start selling items" message
  - Updates AI variables:
    * Sets npc_talk state to 'sell'
    * Records timestamp
  - Packet: 00C7

- npc_store_info - NPC shop item list handler (lines 7588-7633)
  - Processes NPC shop item listings (ZC_PC_PURCHASE_ITEMLIST)
  - Supports multiple packet versions:
    * 0B77: Uses V3 C v V pack format with nameID, price, type, sprite_id, location
    * Others: Uses server-specific pack format (npc_store_info_pack) or default 'V V C v'
  - Clears storeList and talk hash
  - Processes each item in the raw message:
    * Creates new Actor::Item
    * Unpacks item data using appropriate format
    * Sets item ID to storeList size (workaround for duplicate items)
    * Sets item name using itemName function
    * Adds item to storeList
  - Updates AI variables:
    * Sets npc_talk state to 'store'
    * Records timestamp
  - Automatically runs 'store' command if not in buyAuto mode
  - Outputs debug messages for each item added
  - Packets: 00C6, 0B77

- npc_store_begin - NPC shop dialog handler (lines 7573-7581)
  - Processes NPC shop dialog notifications (ZC_SELECT_DEALTYPE)
  - Clears talk hash to reset dialog state
  - Sets talk{ID} to NPC ID
  - Updates AI variables:
    * Sets npc_talk state to 'buy_or_sell'
    * Records timestamp
  - Sets storeList NPC name using getNPCName
  - Falls back to "Unknown" if name not found
  - Simple implementation focused on shop dialog initialization
  - Packet: 00C4

- sell_result - Item selling result handler (lines 9519-9527)
  - Processes item selling result notifications
  - Handles two result states:
    * fail=1: Failure - "Sell failed"
    * fail=0: Success - "Sold X items" and "Sell completed"
  - Clears sellList array after processing
  - Checks if AI is in sellAuto mode
  - Uses "error" message category for failures
  - Uses "success" message category for successful sales
  - Packet: 00CB
- buy_result - Shop purchase result handler (lines 7681-7704)
  - Processes NPC shop purchase results (ZC_PC_PURCHASE_RESULT)
  - Handles multiple result codes:
    * 0: Success - "Buy completed"
    * 1: Insufficient zeny
    * 2: Insufficient weight capacity
    * 3: Too many different inventory items
    * 4: Item does not exist in store
    * 5: Item cannot be exchanged
    * 6: Invalid store
    * Other: Generic failure with code
  - Sets recv_buy_packet flag if in buyAuto mode
  - Triggers buy_result hook with fail code
  - Uses appropriate message categories:
    * "success" for successful purchases
    * Default (error) for failures
  - Packet: 00CA

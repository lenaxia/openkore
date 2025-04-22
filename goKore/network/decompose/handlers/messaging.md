**Method Implementations:**
- self_chat - Self chat message handler (lines 11383-11406)
  - Processes self chat notifications
  - Converts message bytes to string
  - Parses message to extract user and content:
    * Format: "User : Message"
    * Note: May be undefined on eAthena servers (used for non-chat messages)
  - For defined chat messages:
    * Strips language code from message
    * Solves message using solveMessage function
    * Reconstructs full message with user and parsed content
  - Logs to chat log if logChat is enabled
  - Displays message with "selfchat" category
  - Triggers packet_selfChat hook with user and message
  - Contains comments about message format and eAthena behavior
- private_message - Private message handler (lines 11200-11239)
  - Processes private message notifications
  - Requires in-game state (changeToInGameState)
  - Gets sender name and message from packet
  - Converts bytes to string
  - Strips language code from message
  - Solves message using solveMessage function
  - Manages private message users list:
    * Adds new users to privMsgUsers array
    * Triggers parseMsg/addPrivMsgUser hook with user details
  - Logs message to chat log if logPrivateChat is enabled
  - Displays formatted message with sender and content
  - Adds message to ChatQueue with 'pm' type
  - Triggers packet_privMsg hook with message details
  - Handles auto-disconnect on PM:
    * If dcOnPM is enabled and AI state is AUTO:
      - Displays auto-disconnect message
      - Logs to chat log
      - Sends quit packet
      - Quits client
  - Contains comment about parameter types
- manner_message - Manner point notification handler (lines 10457-10470)
  - Processes manner point notifications
  - Handles multiple flag values:
    * 0: Point aligned - "A manner point has been successfully aligned"
    * 3: GM chat block - "Chat Block has been applied by GM due to your ill-mannerous action"
    * 4: Anti-spam block - "Automated Chat Block has been applied due to Anti-Spam System"
    * 5: Good point - "You got a good point"
    * Other: Unknown result with flag value
  - Uses "info" message category for known results
  - Uses "warning" message category for unknown results
- talkie_box - Talkie box message handler (lines 10452-10455)
  - Processes talkie box message notifications
  - Gets actor name using Actor::get and nameString
  - Displays message with actor name and message content
  - Uses "info" message category
  - Simple implementation focused on notification
- private_message_sent - Private message sent handler (lines 9965-9984)
  - Processes private message sent notifications
  - Handles multiple type values:
    * 0: Message sent successfully
      - Displays message with recipient and content
      - Logs to PM chat log if logPrivateChat is enabled
      - Triggers packet_sentPM hook with recipient and message
    * 1: Recipient offline - "X is not online"
    * 2: Message ignored - "Player X ignored your message"
    * Other: Messages disabled - "Player X doesn't want to receive messages"
  - Uses "pm/sent" message category for successful messages
  - Uses "warning" message category for failures
  - Removes first entry from lastpm array after processing
- message_string - Message string handler (lines 9302-9324)
  - Processes message string notifications
  - Increments index from packet for msgTable lookup
  - Converts param from bytes to string if present
  - Handles multiple packet types:
    * 07E2: Message with numeric value (ZC_MSG_VALUE)
    * 09CD: Message with color (ZC_MSG_COLOR)
    * 0A6F: Message with string value (ZC_FORMATSTRING_MSG)
  - For valid message IDs:
    * For 07E2/0A6F with param: Formats message with parameter
    * For others: Shows plain message
  - For unknown message IDs:
    * Displays warning about missing msgstringtable.txt entry
    * Shows index and param values
    * Suggests updating msgstringtable.txt from data.grf
  - Special handling for mercenary dismissal (index 1267-1270)
  - Triggers packet_message_string hook with:
    * index: Message index
    * val: Parameter value
  - Uses "warning" message category
- rodex_delete - Rodex mail deletion handler (lines 8819-8831)
  - Processes Rodex mail deletion notifications
  - Returns early if mail ID doesn't exist in rodexList
  - Displays deletion confirmation message with mail ID
  - Triggers rodex_mail_deleted hook with:
    * mailID: Mail identifier
  - Removes mail entry from rodexList mails hash
  - Simple implementation focused on cleanup
- rodex_get_item - Rodex item retrieval handler (lines 8805-8817)
  - Processes Rodex mail item retrieval results
  - Handles failure case:
    * Displays error message
    * Returns early
  - For successful retrieval:
    * Displays success message
    * Clears items array in mail data
    * Updates attach flag:
      - If attach was 'i' (item only), sets to undef (no attachment)
      - Otherwise sets to 'z' (zeny only)
  - Simple implementation focused on result notification
- rodex_get_zeny - Rodex zeny retrieval handler (lines 8791-8803)
  - Processes Rodex mail zeny retrieval results
  - Handles failure case:
    * Displays error message
    * Returns early
  - For successful retrieval:
    * Displays success message
    * Resets zeny1 value to 0 in mail data
    * Updates attach flag:
      - If attach was 'z' (zeny only), sets to 0 (no attachment)
      - Otherwise sets to 'i' (item only)
  - Simple implementation focused on result notification
- rodex_write_result - Rodex mail sending result handler (lines 8779-8789)
  - Processes Rodex mail sending results
  - Handles failure case:
    * Displays error message
    * Returns early
  - For successful sending:
    * Displays success message
    * Clears rodexWrite structure
  - Simple implementation focused on result notification
- rodex_check_player - Rodex player verification handler (lines 8748-8772)
  - Processes Rodex player name verification results
  - Converts player name from bytes to string
  - Handles player not found case:
    * Displays error message
    * Deletes target from rodexWrite
    * Returns early
  - Supports multiple packet versions:
    * 0A14: Basic format with char_id, class, base_level
    * 0A51: Extended format that includes name field
  - For 0A51, updates target name in rodexWrite
  - Displays formatted player information:
    * Name and base level
    * Character ID and job class
  - Uses centered header with separator lines
  - Uses swrite for formatted output
- rodex_open_write - Rodex mail composition handler (lines 8734-8746)
  - Processes Rodex mail composition window opening
  - Initializes rodexWrite structure:
    * Creates new InventoryList for items
    * Sets default title to "TITLE"
  - Handles target name if provided:
    * Converts from bytes to string
    * Stores in rodexWrite target structure
    * Sends checkname request to server
  - Outputs debug message with target and title
  - Simple implementation focused on initialization
- rodex_add_item - Rodex item addition handler (lines 8692-8732)
  - Processes Rodex mail item addition results
  - Handles multiple failure cases:
    * 1: Weight error - "Item attachment has been failed"
    * 2: General failure - "Item attachment has been failed"
    * 3: Maximum attachments exceeded
    * 4: Banned item type
    * Other: Unknown error with code
  - Returns early on any failure
  - For successful addition:
    * Gets item reference from rodexWrite items list
    * If item exists, increases amount
    * If item doesn't exist:
      - Creates new Actor::Item
      - Sets all item properties (ID, nameID, type, etc.)
      - Gets item name using itemName function
      - Adds to rodexWrite items list
    * Displays message with item details:
      - Name, binID, amount, item type
  - Uses "drop" message category
- rodex_remove_item - Rodex item removal handler (lines 8672-8690)
  - Processes Rodex mail item removal results
  - Handles failure case:
    * Displays error message
    * Returns early
  - For successful removal:
    * Gets item reference from rodexWrite items list
    * Displays message with item details:
      - Name, binID, amount, item type
    * Updates item amount in rodexWrite
    * Removes item from list if amount reaches zero
  - Uses "drop" message category
- unread_rodex - Unread Rodex mail notification handler (lines 8666-8670)
  - Processes unread Rodex mail notifications
  - Displays simple notification message
  - Triggers rodex_unread_mail hook
  - Simple implementation focused on notification
- rodex_read_mail - Rodex mail content handler (lines 8580-8664)
  - Processes Rodex mail content packets
  - Extracts mail body:
    * Gets raw message and size
    * Defines header pack format and calculates length
    * Extracts text portion based on text_len
    * Converts from bytes to string
    * Removes trailing newlines
    * Solves message using solveMSG function
  - Stores mail metadata:
    * Zeny amounts (zeny1, zeny2)
    * Mail type with human-readable mapping
  - Processes attached items:
    * Uses server-specific item pack format
    * Calculates item length
    * Creates empty items array
    * Unpacks each item's properties
    * Gets item name using itemName function
    * Formats display string with amount
    * Adds to mail items array
  - Displays formatted mail content:
    * Mail header with ID and sender
    * Mail type and title
    * Message body
    * Item count and zeny amount
    * List of attached items
  - Updates mail status in rodexList:
    * Copies body, items, zeny values
    * Sets isRead flag to 1
    * Sets current_read to mailID1
  - Triggers rodex_mail hook with:
    * mailID: Mail identifier
    * from: Sender name
    * title: Mail title
    * content: Mail body
    * zeny: Zeny amount
    * itemCount: Number of items
    * items: Array of item objects
  - Uses "list" message category
- rodex_mail_list - Rodex mail list handler (lines 8490-8578)
  - Processes Rodex mail list packets
  - Supports multiple packet versions:
    * 0B5F: Newer version with 45-byte entries
    * 0AC2: Version with 41-byte entries
    * 09F0/0A7D: Default version with 44-byte entries
  - Defines different data structures per version:
    * Basic: mailID1, mailID2, isRead, attach, sender, regDateTime, expireSecconds, Titlelength
    * Others: Add openType field
  - Manages mail list state:
    * For 0A7D/0AC2/0B5F: Resets list and current page
    * For others: Increments current page
    * Sets last_page flag if isEnd=1
    * Tracks mails_per_page from amount field
  - Processes each mail entry:
    * Unpacks mail data using version-specific format
    * Extracts and decodes title from message
    * Converts sender name from bytes to string
    * Calculates expire days from seconds
    * Maps attachment type to readable format
  - Displays formatted mail list with:
    * Index, ID, sender, attachment type, read status, expiration, title
  - Triggers rodex_mail_list hook with:
    * mails: Mail list hash
    * current_page: Current page number
    * last_mailID: Last mail ID on current page
    * isEnd: End of list flag
  - Uses "list" message category
- whisper_list - Whisper list handler (lines 5571-5577)
  - Processes whisper list packets
  - Unpacks list of whisper names from raw message
  - Each name is 24 bytes in length
  - Calculates number of names based on message size
  - Outputs debug message with list of names
  - Simple implementation with minimal processing
- local_broadcast - Local broadcast message handler (lines 3137-3148)
  - Processes ZC_BROADCAST2 packets
  - Handles formatted messages with color codes
  - Converts message bytes to string
  - Formats color as 6-digit hex code
  - Supports logging via chatLog when configured
  - Triggers packet_localBroadcast plugin hook
  - Displays messages in schat channel

- system_chat - System chat message handler (lines 3482-3495)
  - Processes various system chat message formats
  - Handles different message prefixes:
    - "ssss": War of Emperium messages (yellow color)
    - "micc": Player broadcast messages (with color codes)
    - "blue": System messages (blue color)
  - Extracts color codes from micc format messages
  - Uses internationalized strings (T()) for prefixes
  - Supports Chinese character detection in names
  - Handles null-padded message formats
  - Processes message bytes to string conversion
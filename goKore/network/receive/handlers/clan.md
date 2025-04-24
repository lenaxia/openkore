**Clan System Handlers:**

- clan_leave - Clan leave handler (lines 8967-8974)
  - Processes clan leave notifications
  - Checks if character is in a clan (clan_name exists)
  - Displays message with clan name when leaving
  - Clears clan data structure (undef %clan)
  - Simple implementation focused on cleanup

- clan_chat - Clan chat message handler (lines 8944-8965)
  - Processes clan chat messages
  - Requires in-game state (changeToInGameState)
  - Converts message components from bytes:
    * charname: Sender name
    * message: Chat message content
  - Solves message using solveMessage function
  - Handles logging:
    * Writes to clan chat log if logClanChat is enabled
    * Adds to ChatQueue with 'clan' type
    * Outputs debug message with "clanchat" category
  - Displays formatted message with [Clan] prefix
  - Triggers packet_clanMsg hook with:
    * MsgUser: Sender username
    * Msg: Parsed message content
    * RawMsg: Original message
  - Uses "clanchat" message category

- clan_info - Clan information handler (lines 8913-8942)
  - Processes clan information notifications
  - Updates clan data structure with:
    * clan_ID: Clan identifier
    * clan_name: Clan name (converted from bytes)
    * clan_master: Clan leader name (converted from bytes)
    * clan_map: Clan map name (converted from bytes)
    * alliance_count: Number of allied clans
    * antagonist_count: Number of enemy clans
  - Processes allied clan names:
    * Extracts from ally_antagonist_names field
    * Unpacks each 24-byte name segment
    * Converts from bytes to string
    * Appends to ally_names string with comma separator
  - Processes enemy clan names:
    * Extracts from ally_antagonist_names field (continuing from allies)
    * Unpacks each 24-byte name segment
    * Converts from bytes to string
    * Appends to antagonist_names string with comma separator
  - Uses foreach loop to process multiple fields

- clan_user - Clan user count handler (lines 8904-8911)
  - Processes clan user count notifications
  - Updates clan data structure with:
    * onlineuser: Number of online clan members
    * totalmembers: Total number of clan members
  - Uses foreach loop to process multiple fields
  - Simple implementation focused on data updates
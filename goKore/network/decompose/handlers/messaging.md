**Chat System Handlers:**
- whisper_list (lines 5571-5577)
  - Processes whisper ignore list
  - Unpacks list of character names (24B each)
  - Debug logs the list

- chat_created (lines 5585-5595)
  - Handles successful chat room creation
  - Sets currentChatRoom to accountID
  - Stores room info in chatRooms hash
  - Calls 'chat_created' hook

- chat_info (lines 5604-5629)
  - Displays chat room info
  - Processes room title, owner, limit, type
  - Calls 'packet_chatinfo' hook with details
  - Supports private/public/arena/PK room types

- chat_users (lines 5637-5666)
  - Handles entering a chat room
  - Processes user list with roles (owner/normal)
  - Updates currentChatRoomUsers array
  - Calls 'chat_joined' hook

- chat_join_result (lines 5680-5701)
  - Handles chat join failures:
    - 0: Room full
    - 1: Wrong password
    - 2: Kicked
    - 4: Not enough zeny
    - 5: Level too low
    - 6: Level too high
    - 7: Wrong job class

- chat_modified (lines 5711-5739)
  - Handles chat room property changes
  - Updates title, limit, type, user count
  - Calls 'chat_modified' hook with old/new data

- chat_newowner (lines 5747-5771)
  - Processes owner change notifications
  - Updates chatRooms ownerID
  - Sets user role (0=owner, 1=normal)

- chat_user_join (lines 5776-5785)
  - Handles new user joining chat
  - Updates user list and count
  - Displays join message

- chat_user_leave (lines 5793-5809)
  - Handles user leaving chat
  - Removes from user list
  - Special handling for self-leaving
  - Calls 'chat_leave' hook when leaving

- chat_removed (lines 5814-5823)
  - Handles chat room deletion
  - Cleans up chatRooms data
  - Calls 'chat_removed' hook
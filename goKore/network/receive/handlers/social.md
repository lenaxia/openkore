**Social Interaction Handlers:**

- married - Marriage notification handler (lines 7004-7009)
  - Processes marriage event notifications
  - Gets actor reference using Actor::get
  - Displays "[Actor] got married!" message
  - Simple implementation focused on notification

- ignore_player_result - Individual ignore result handler (lines 6946-6955)
  - Processes results of ignoring individual players
  - Handles two operation types:
    * 0 = Enable ignore: Displays "Player ignored" message
    * 1 = Disable ignore:
      - Checks for error code 0 (success)
      - Displays "Player unignored" message
  - Contains TODO comment about storing list of ignored players
  - Simple implementation focused on user feedback

- ignore_all_result - Global ignore result handler (lines 6933-6943)
  - Processes results of ignore all players commands
  - Handles two operation types:
    * 0 = Enable ignore all:
      - Sets ignored_all flag to 1
      - Displays "All Players ignored" message
    * 1 = Disable ignore all:
      - Checks for error code 0 (success)
      - Displays "All players unignored" message
  - Contains TODO comment about storing this state
  - Simple implementation focused on user feedback

- friend_response - Friend request response handler (lines 6235-6258)
  - Processes friend request response notifications (ZC_ADD_FRIENDS_LIST)
  - Extracts response information:
    * Response type
    * Friend name (converted from bytes to string)
  - Handles different response types:
    * 0 = Success:
      - Adds new friend to friendsID array
      - Extracts accountID and charID from raw message
      - Sets name and online status (1)
      - Displays "You have become friends with" message
    * 1 = Rejection: "[Player] does not want to be friends"
    * 2 = Your list full: "Your Friend List is full"
    * 3 = Their list full: "[Player]'s Friend List is full"
    * Other: "[Player] rejected to be your friend"
  - Result codes (from packet comments):
    * 0 = "You have become friends with (%s)."
    * 1 = "(%s) does not want to be friends with you."
    * 2 = "Your Friend List is full."
    * 3 = "(%s)'s Friend List is full."

- friend_removed - Friend removal handler (lines 6212-6226)
  - Processes friend removal notifications (PACKET_ZC_DELETE_FRIENDS)
  - Extracts friend information:
    * Account ID
    * Character ID
  - Searches through friendsID array for matching friend
  - When found:
    * Displays "[Friend] is no longer your friend" message
    * Removes friend ID from friendsID array
    * Deletes friend entry from friends hash
  - Simple implementation focused on friend removal

- friend_request - Friend request handler (lines 6194-6208)
  - Processes incoming friend requests (ZC_REQ_ADD_FRIENDS)
  - Stores request information in incomingFriend hash:
    * Account ID
    * Character ID
    * Name (converted from bytes to string)
  - Displays notification messages:
    * "[Player] wants to be your friend"
    * Instructions for accepting/rejecting request
  - Triggers friend_request hook with:
    * Account ID
    * Character ID
    * Name
  - Simple implementation focused on request notification

- friend_logon - Friend online status handler (lines 6171-6190)
  - Processes friend online/offline notifications (ZC_FRIENDS_STATE)
  - Extracts friend information:
    * Account ID
    * Character ID
    * Online status (isNotOnline flag)
  - Searches through friendsID array for matching friend
  - Updates online status in friends hash (1-isNotOnline)
  - Displays appropriate message:
    * "[Friend] has disconnected" when isNotOnline=1
    * "[Friend] has connected" when isNotOnline=0
  - State values (from packet comments):
    * 0 = online
    * 1 = offline

- friend_list - Friend list handler (lines 6143-6163)
  - Processes friend list packets (ZC_FRIENDS_LIST)
  - Resets existing friend data structures:
    * Clears friendsID array
    * Clears friends hash
  - Parses raw message data in 32-byte chunks
  - For each friend entry:
    * Adds ID to friendsID array
    * Extracts accountID, charID, and name
    * Converts name bytes to string
    * Initializes online status to 0 (offline)
    * Increments ID counter
  - Simple implementation focused on data structure setup

- emoticon - Emoticon display handler (lines 5961-6024)
  - Processes emoticon/emotion display packets
  - Looks up emotion display text from emotions_lut
  - Handles different actor types:
    * Current character:
      - Displays simple message with name and emotion
    * Other players:
      - Calculates distance to player
      - Displays message with distance, name, and emotion
      - Handles follow emotion response if configured
      - Mirrors specific emotions (30↔31) when following
    * Monsters/slaves:
      - Calculates distance
      - Displays message with actor type, name, and emotion
    * Other/unknown actors:
      - Uses generic display with nameIdx
  - Logs emotions to chat log if configured
  - Triggers packet_emotion hook with emotion and ID

- chat_removed - Chat room removal handler (lines 5813-5823)
  - Processes chat room removal packets (ZC_DESTROY_ROOM)
  - Updates chat room data structures:
    * Removes room ID from chatRoomsID list
    * Deletes room from chatRooms hash
  - Stores removed chat data for hook
  - Triggers chat_removed hook with:
    * Room ID
    * Room data (before removal)
  - Simple implementation for room removal events

- chat_user_leave - Chat room user leave handler (lines 5792-5809)
  - Processes user leave notifications (ZC_MEMBER_EXIT)
  - Converts user name bytes to string
  - Updates chat room data structures:
    * Removes user from room's users hash
    * Removes user from currentChatRoomUsers list
    * Updates room's user count
  - Handles two scenarios:
    * Current character leaving:
      - Removes room from chatRoomsID list
      - Deletes room from chatRooms hash
      - Clears currentChatRoomUsers array
      - Resets currentChatRoom
      - Displays "You left" message
      - Triggers chat_leave hook
    * Other user leaving:
      - Displays "[user] has left" message
  - Leave flags (from packet comments):
    * 0 = left
    * 1 = kicked

- chat_user_join - Chat room user join handler (lines 5775-5785)
  - Processes user join notifications (ZC_MEMBER_NEWENTRY)
  - Converts user name bytes to string
  - Only processes if in a chat room (currentChatRoom not empty)
  - Updates chat room data structures:
    * Adds user to currentChatRoomUsers list
    * Sets user's role in room to normal (1)
    * Updates room's user count
  - Displays join message with user name
  - Simple implementation for user join events

- chat_newowner - Chat room owner change handler (lines 5746-5771)
  - Processes chat room owner change packets (ZC_ROLE_CHANGE)
  - Converts user name bytes to string
  - Handles owner role assignment (type 0):
    * If current character is new owner:
      - Sets room ownerID to accountID
    * If another player is new owner:
      - Searches playersList for matching name
      - Sets room ownerID to player ID when found
    * Updates user's role in room to owner (2)
  - Handles normal role assignment (type 1):
    * Updates user's role in room to normal (1)
  - User roles (from packet comments):
    * 0 = owner (menu)
    * 1 = normal

- chat_modified - Chat room modification handler (lines 5710-5739)
  - Processes chat room property changes (ZC_CHANGE_CHATROOM)
  - Converts title bytes to string
  - Extracts room properties:
    * Owner ID, chat ID, user limit
    * Public/private status
    * Current user count
  - Determines room ID based on ownership
  - Creates temporary chat data structure
  - Triggers chat_modified hook with:
    * Room ID
    * Old room data
    * New room data
  - Updates chatRooms global with new properties
  - Displays modification message to user
  - Chat room types (from packet comments):
    * 0 = private (password protected)
    * 1 = public
    * 2 = arena (NPC waiting room)
    * 3 = PK zone (non-clickable)

- chat_join_result - Chat room join result handler (lines 5679-5701)
  - Processes chat room join result packets (ZC_REFUSE_ENTER_ROOM)
  - Displays appropriate message based on result type:
    * 0: Room is full
    * 1: Incorrect password
    * 2: You're kicked (note: has duplicate cases)
    * Other types: Unknown reason
  - Note: Has implementation issues with duplicate type=2 cases for:
    * Joined Chat Room
    * Not enough zeny
    * Low level
    * High level
    * Unsuitable job class
  - Simple implementation focused on user feedback

- chat_users - Chat room users handler (lines 5636-5666)
  - Processes chat room user list packets (ZC_ENTER_ROOM)
  - Extracts chat room ID from raw message
  - Sets currentChatRoom to this ID
  - Creates or retrieves chat room data structure
  - Resets user count for the room
  - Processes each user entry (28 bytes each):
    * Extracts user type and name
    * Converts name bytes to string
    * Adds new users to currentChatRoomUsers list
    * Sets user role (2 for owner, 1 for normal)
    * Increments room user count
  - Displays room join message with room title
  - Triggers chat_joined hook with room data
  - User roles (from packet comments):
    * 0 = owner (menu)
    * 1 = normal

- chat_info - Chat room information handler (lines 5604-5629)
  - Processes chat room information packets (ZC_ROOM_NEWENTRY)
  - Converts title bytes to string
  - Manages chat room data structures:
    * Creates new chat room entry if needed
    * Updates existing chat room information
    * Adds room ID to chatRoomsID list if new
  - Updates room properties:
    * Title, owner ID, user limit
    * Public/private status
    * Current user count
  - Triggers packet_chatinfo hook with room data
  - Chat room types (from packet comments):
    * 0 = private (password protected)
    * 1 = public
    * 2 = arena (NPC waiting room)
    * 3 = PK zone (non-clickable)

- chat_created - Chat room creation handler (lines 5585-5595)
  - Processes chat room creation result packets (ZC_ACK_CREATE_CHATROOM)
  - Sets up chat room data structures:
    * Sets currentChatRoom to accountID
    * Copies createdChatRoom data to chatRooms
    * Adds room ID to chatRoomsID list
    * Adds character name to currentChatRoomUsers
  - Displays success message to user
  - Triggers chat_created hook with room data
  - Flag values (not used in handler):
    * 0 = Room successfully created
    * 1 = Room limit exceeded
    * 2 = Same room already exists

- marriage_partner_name() - Marriage partner name handler (lines 4112-4116)
  - Processes marriage partner name packets
  - Converts byte data to string using bytesToString()
  - Displays partner name in message log
  - Triggered before "I miss you" skill is cast


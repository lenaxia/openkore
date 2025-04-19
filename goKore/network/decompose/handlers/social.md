**Social Interaction Handlers:**

- adopt_reply() - Adoption system responses (lines 9557-9566)
  - Handles adoption request responses:
    * Already have child (type 0)
    * Level too low (type 1)
    * Cannot adopt married person (type 2)
  - Displays appropriate rejection messages

- ignore_all_result (lines 6933-6943)
  - Manages global ignore state
  - Handles ZC_ACK_ALL_IGNORE packets
  - Features:
    - Toggles ignored_all flag
    - Handles enable/disable cases
    - Provides user feedback
    - Supports error reporting

- ignore_player_result (lines 6946-6955)
  - Processes individual ignore operations
  - Handles ZC_ACK_IGNORE packets
  - Features:
    - Supports ignore/unignore
    - Handles error cases
    - Provides user feedback
    - Maintains player ignore state

- married (lines 7004-7009)
  - Handles marriage announcements
  - Processes ZC_NOTIFY_MARRIAGE packets
  - Features:
    - Displays marriage notifications
    - Supports actor-based announcements
    - Uses proper string formatting

**Friend List Handlers:**
- friend_list (lines 6140-6155)
  - Manages friend relationships
  - Processes:
    - Account IDs
    - Character IDs
    - Names (when available)
  - Maintains:
    - @friendsID array

- friend_status (lines 6157-6170)
  - Tracks friend online/offline status changes
  - Features:
    - Updates friend list display
    - Triggers 'friend_status' hook

- friend_request (lines 6172-6185)
  - Processes incoming friend requests
  - Handles:
    - Requestor account ID

- friend_response (lines 6187-6202)
  - Processes friend request responses
  - Handles:
    - Response status (accept=1/reject=0)

- friend_remove (lines 6204-6218)
  - Processes friend removal operations
  - Handles:
    - Removed account ID
    - Removal confirmation
  - Features:
    - Updates friend list
    - Triggers 'friend_removed' hook
    - Validates removal integrity
    - Target account ID
    - Response timestamp
  - Features:
    - Updates friend list on accept
    - Triggers 'friend_response' hook
    - Validates response integrity
    - Character name
    - Request timestamp
  - Features:
    - Triggers 'friend_request' hook
    - Provides accept/decline options
    - Validates requestor info
    - Handles both online (1) and offline (0) states
    - Matches status to friend IDs
    - %friends hash structure
  - Handles both name and nameless packet formats
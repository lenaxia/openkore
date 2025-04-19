**Mail Notification Handlers:**
- unread_rodex() - Notifies of new mail (lines 8666-8670)
  - Displays unread mail notification
  - Triggers rodex_unread_mail hook

**Mail Item Management:**
- rodex_remove_item() - Handles item removal (lines 8672-8690)
  - Validates removal success
  - Updates item amounts
  - Removes item if amount reaches 0
  - Logs removal details

- rodex_add_item() - Handles item attachment (lines 8692-8732)
  - Processes various failure cases:
    * Weight error
    * Attachment limit
    * Banned items
  - Updates existing items or adds new ones
  - Logs addition details

**Mail Composition Handlers:**
- rodex_open_write() - Initializes mail compose (lines 8734-8746)
  - Creates new mail structure
  - Sets default title
  - Initiates name check if recipient provided

- rodex_check_player() - Validates recipient (lines 8748-8777)
  - Handles different packet versions
  - Displays recipient details:
    * Name, level, class
    * Character ID
  - Stores validated recipient info

- rodex_write_result() - Processes send result (lines 8779-8789)
  - Handles send failures
  - Cleans up after successful send

**Mail Retrieval Handlers:**
- rodex_get_zeny() - Handles zeny collection (lines 8791-8803)
  - Processes collection failures
  - Updates mail attachment status

- rodex_get_item() - Handles item collection (lines 8805-8817)
  - Processes collection failures
  - Updates mail attachment status

**Mail Deletion Handler:**
- rodex_delete() - Handles mail deletion (lines 8819-8831)
  - Confirms deletion
  - Triggers rodex_mail_deleted hook
  - Removes mail from list

**Mail Listing Handler:**
- rodex_mail_list() - Processes mail list packets (lines 8490-8578)
  - Handles multiple packet versions (0B5F, 0AC2, 09F0, 0A7D)
  - Parses mail metadata:
    * Mail IDs, read status
    * Sender, expiration
    * Attachment types
    * Title length/content
  - Manages pagination state
  - Formats output display
  - Triggers mail_list hook

**Mail Reading Handler:**
- rodex_read_mail() - Processes mail content (lines 8580-8598)
  - Extracts mail body text
  - Processes zeny amounts
  - Handles mail type
  - Parses header information
  - Cleans up message formatting

**Mail System Handlers:**
- mail_list (lines 6220-6235)
  - Manages mail inbox listing
  - Processes:
    - Mail IDs
    - Sender names
    - Timestamps
    - Read/unread status
  - Maintains:
    - @mailList array
    - %mail hash structure

- mail_read (lines 6237-6245)
  - Handles mail content retrieval

- mail_send_result (lines 6247-6260)
  - Processes mail sending outcomes
  - Handles:
    - Success/failure status codes
    - Error messages

- mail_delete_result (lines 6262-6275)
  - Processes mail deletion outcomes
  - Handles:
    - Deletion status codes
    - Mail ID references
  - Features:
    - Updates mail list on success
    - Triggers 'mail_deleted' hook
    - Provides user feedback
    - Recipient validation
  - Features:
    - Triggers 'mail_sent' hook
    - Provides user feedback
    - Logs send attempts
  - Features:
    - Processes mail body text
    - Updates read status
    - Handles attachments
    - Triggers 'mail_read' hook
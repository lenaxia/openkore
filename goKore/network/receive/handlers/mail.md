# Mail System Related Handlers

**Method Implementations:**
- mail_new - New mail notification handler (lines 10871-10874)
  - Processes new mail notifications
  - Displays message with:
    * Sender name (converted from bytes)
    * Mail title (converted from bytes)
  - Uses "info" message category
  - Simple implementation focused on notification
- mail_send - Mail sending result handler (lines 10864-10869)
  - Processes mail sending result notifications
  - Uses ternary operator for conditional message:
    * If fail: "Failed to send mail, the recipient does not exist"
    * If success: "Mail sent succesfully"
  - Uses "info" error category for failure
  - Uses "info" message category for success
  - Simple implementation focused on notification
- mail_setattachment - Mail attachment setting handler (lines 10835-10862)
  - Processes mail attachment setting result notifications
  - Handles failure case:
    * Clears mailAttachAmount if defined
    * Displays failure message with item or zeny info
  - Handles success cases:
    * For items:
      - Gets item from inventory by ID
      - Displays success message with item info
      - If mailAttachAmount is defined:
        * Calculates change amount (min of item amount and mailAttachAmount)
        * Removes item from inventory
        * Triggers packet_item_removed hook
        * Clears mailAttachAmount
    * For zeny:
      - Displays success message
      - If mailAttachAmount is defined:
        * Calculates change amount (min of zeny and mailAttachAmount)
        * Deducts zeny from character
        * Displays message about zeny loss
  - Uses "info" message category for all messages
  - Complex implementation handling both item and zeny attachments
- mail_getattachment - Mail attachment retrieval handler (lines 10824-10833)
  - Processes mail attachment retrieval result notifications
  - Handles multiple fail values:
    * 0: Success - "Successfully added attachment to inventory"
    * 2: Weight failure - "Failed to get the attachment to inventory due to your weight"
    * Other: General failure - "Failed to get the attachment to inventory"
  - Uses "info" message category for success
  - Uses "info" error category for failures
  - Simple implementation focused on notification
- mail_refreshinbox - Mail inbox refresh handler (lines 10780-10822)
  - Processes mail inbox refresh notifications
  - Stores old mail count and clears mailList
  - Handles empty inbox case:
    * Displays "There is no mail in your inbox" message
    * Returns early
  - Skips processing if count hasn't changed
  - Displays mail count message
  - Creates formatted inbox display:
    * Centered "Inbox" header with separator lines
    * Column headers: #, R, Title, Sender, Date
    * Separator line
  - Processes each mail entry in 73-byte chunks:
    * Unpacks mailID, title, read status, sender, timestamp
    * Converts title and sender from bytes to string
    * Formats each entry with:
      - Index number
      - Read status
      - Title (truncated to 34 chars)
      - Sender (truncated to 24 chars)
      - Formatted date (without year)
  - Adds footer separator line
  - Displays complete formatted list
  - Uses "info" message category for notifications
  - Uses "list" message category for the mail list
  - Contains comments about truncation decisions
- mail_read - Mail content display handler (lines 10756-10778)
  - Processes mail content notifications
  - Creates item object with:
    * nameID, upgrade, cards, broken from packet
    * name from itemName function
  - Builds formatted message with:
    * Centered "Mail" header with separator lines
    * Title and sender (converted from bytes)
    * Message content (converted from bytes)
    * Separator line
    * Item details (name and amount if present)
    * Zeny amount (formatted with formatNumber)
    * Footer separator line
  - Displays complete formatted message
  - Uses "info" message category
  - Complex implementation focused on formatted display
- mail_return - Mail return handler (lines 10749-10754)
  - Processes mail return result notifications
  - Uses ternary operator for conditional message:
    * If fail: "The mail with ID: X does not exist"
    * If success: "The mail with ID: X is returned to the sender"
  - Uses "info" error category for failure
  - Uses "info" message category for success
  - Simple implementation focused on notification
- mail_window - Mail window status handler (lines 10739-10747)
  - Processes mail window status notifications
  - Handles two flag values:
    * 1: Closed - "Mail window is now closed"
    * 0: Opened - "Mail window is now opened"
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
- mail_delete - Mail deletion handler (lines 10729-10737)
  - Processes mail deletion result notifications
  - Handles two fail values:
    * 1: Failure - "Failed to delete mail with ID: X"
    * 0: Success - "Succeeded to delete mail with ID: X"
  - Displays message with mail ID
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
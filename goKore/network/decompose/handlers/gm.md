# Game Master Related Handlers

**Method Implementations:**
- GM_silence - GM silence status handler (lines 9568-9575)
  - Processes GM silence status notifications
  - Gets GM name from packet (converted from bytes)
  - Handles two flag values:
    * 1: Muted - "You have been: muted by X"
    * 0: Unmuted - "You have been: unmuted by X"
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
- GM_req_acc_name - GM account name request handler (lines 9533-9536)
  - Processes GM account name request responses
  - Displays message with target ID and account name
  - Uses "info" message category
  - Simple implementation focused on notification
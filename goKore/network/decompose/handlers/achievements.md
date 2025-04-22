# Achievement Related Handlers

**Method Implementations:**
- achievement_list - Achievement list handler (lines 9452-9484)
  - Processes achievement list notifications
  - Clears existing achievementList
  - Gets raw message and size
  - Defines header length (22 bytes)
  - Sets up achievement pack format and calculates length
  - Processes each achievement entry:
    * Unpacks achievement data with multiple fields:
      - achievementID: Achievement identifier
      - completed: Completion status
      - objective1-10: Progress on objectives
      - completed_at: Completion timestamp
      - reward: Reward information
    * Updates achievementList with achievement
    * Displays message with achievement title and ID
  - Uses achievements hash to look up titles
  - Falls back to empty string if achievement not found
  - Uses "info" message category
- achievement_update - Achievement update handler (lines 9442-9450)
  - Processes achievement update notifications
  - Creates achievement hash with multiple fields:
    * achievementID: Achievement identifier
    * completed: Completion status
    * objective1-10: Progress on objectives
    * completed_at: Completion timestamp
    * reward: Reward information
  - Updates achievementList with new/updated achievement
  - Displays message with achievement title and ID
  - Uses achievements hash to look up title
  - Falls back to empty string if achievement not found
  - Uses "info" message category
- achievement_reward_ack - Achievement reward acknowledgment handler (lines 9437-9440)
  - Processes achievement reward notifications
  - Displays message with achievement title and ID
  - Uses achievements hash to look up title
  - Falls back to empty string if achievement not found
  - Uses "info" message category
  - Simple implementation focused on notification
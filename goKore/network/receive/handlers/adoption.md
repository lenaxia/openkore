# Adoption Related Handlers

**Method Implementations:**
- adopt_request - Adoption request handler (lines 10117-10120)
  - Processes adoption request notifications
  - Displays message with requester's name
  - Uses "info" message category
  - Contains TODO comment about sourceID and targetID
  - Simple implementation focused on notification
- adopt_reply - Adoption request reply handler (lines 9557-9566)
  - Processes adoption request reply notifications
  - Handles multiple type values:
    * 0: Multiple children - "You cannot adopt more than 1 child"
    * 1: Level requirement - "You must be at least character level 70 in order to adopt someone"
    * 2: Married target - "You cannot adopt a married person"
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
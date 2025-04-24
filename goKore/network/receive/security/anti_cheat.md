# Anti-Cheat Related Handlers

**Method Implementations:**
- hack_shield_alarm - Hack Shield detection handler (lines 10447-10450)
  - Processes Hack Shield detection notifications
  - Displays error message about forced disconnect
  - Mentions Poseidon (server component)
  - Runs 'relog 100000000' command (very long delay)
  - Uses "connection" error category
  - Simple implementation focused on notification and automatic relog
- gameguard_request - GameGuard query handler (lines 6318-6326)
  - Processes GameGuard verification requests
  - Skips processing for:
    * XKore mode 1 with gameGuard != "2"
    * Servers with gameGuard = 0
  - Forwards raw packet data to Poseidon client:
    * Uses Poseidon::Client::getInstance()->query()
    * Passes exact raw message bytes
  - Outputs debug message when querying Poseidon
  - Simple implementation focused on forwarding requests
- gameguard_grant - GameGuard login grant handler (lines 6300-6316)
  - Processes GameGuard login permission responses
  - Handles different server types:
    * 0: Login denied (incorrect/delayed response)
      - Displays error about Poseidon server version
    * 1: Account server login granted
      - Displays success message
    * Other: Char/map server login granted
      - Displays success message
      - Calls change_to_constate25() for gameGuard="2"
  - Updates network state from 1.2 to 1.3 if needed
  - Uses "poseidon" message category
- EAC_key - Easy Anti-Cheat detection handler (lines 6293-6298)
  - Detects Easy Anti-Cheat protection on server
  - Skips processing if ignoreAntiCheatWarning is enabled in masterServer
  - Logs detection message to chat log with "k" category
  - Displays error message about lack of EAC support
  - Forces client to quit
  - Simple implementation focused on graceful exit
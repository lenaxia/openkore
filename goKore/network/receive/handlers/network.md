# Network Related Handlers

**Method Implementations:**

- errors - Server error and disconnection handler (lines 6026-6138)
  - Processes server error notifications and forced disconnects
  - Triggers disconnected hook if in-game
  - Handles auto-disconnect based on configuration:
    * dcOnDisconnect: Auto-disconnects on most errors
    * dcOnServerShutDown: Auto-disconnects on server shutdown
    * dcOnServerClose: Auto-disconnects when server closes
    * dcOnDualLogin: Auto-disconnects or reconnects on dual login
  - Resets network state and connection variables
  - Sets reconnect timeout
  - Handles numerous error types with specific messages:
    * 0: Server shutting down
    * 1: Server closed
    * 2: Dual login prohibited
    * 3: Out of sync with server
    * 4: Server jammed due to over-population
    * 5: Underaged player restriction
    * 6: Payment required
    * 8: Server still recognizes last connection
    * 9: IP capacity of Internet Cafe is full
    * 10: Out of available paid time
    * 15: Forced disconnect by GM
    * 101: Account suspended for possible 3rd party program use
    * 102: Too many connections from same IP
    * Others: Unknown error
  - Contains extensive comments documenting all error codes
  - Packet: 0081
  - Format: 'B' (type)
- isvr_disconnect - isvr disconnection handler (lines 11739-11742)
  - Processes isvr disconnection notifications
  - Minimal implementation that only outputs debug message
  - No state changes or actions taken
  - No plugin hooks triggered
- sync_request - Synchronization request handler (lines 11407-11425)
  - Processes server synchronization request notifications
  - Contains comments about uncertainty of packet purpose
  - Only processes for serverType 1
  - Validates that ID matches accountID:
    * If matched:
      - Updates ai_sync timeout timestamp
      - Sends sync packet if client not marked as alive
      - Outputs debug message about sync request
      - Uses "connection" debug category
    * If mismatched:
      - Displays warning about wrong ID
      - Uses "warning" message category
  - Packet: 0187
  - Format: 'a4' (ID)
- sync_request_ex - Extended synchronization request handler (lines 4412-4441)
  - Processes extended server synchronization request notifications
  - Skips processing for XKore modes 1 and 3 (lets client handle it)
  - Contains commented debug logging code
  - Gets packet ID from switch value
  - Looks up corresponding sync_ex_reply ID from table
  - Cleans leading zeros from both packet ID and sync ID
  - Converts sync ID to hexadecimal number
  - Sends reply using sendReplySyncRequestEx with the sync ID
  - Contains detailed comments about the process
  - Implementation by Fr3DBr (credited in comments)
  - Simple implementation focused on anti-bot protection response
  - Debug message: "Received the package 'isvr_disconnect'"
- received_sync - Server synchronization handler (lines 8114-8119)
  - Processes server synchronization notifications
  - Requires in-game state (changeToInGameState)
  - Outputs debug message about sync reception
  - Updates play timeout timestamp
  - Simple implementation focused on keepalive functionality
  - Uses "parseMsg" debug category (priority 2)
  - No plugin hooks triggered
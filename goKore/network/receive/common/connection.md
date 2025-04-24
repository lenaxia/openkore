**Connection Handlers:**

- map_load_error - Map server connection error handler (lines 1297-1328)
  - Error code handling:
    - 0: Wrong client type
    - 1: Wrong ID
    - 2: Timeout
    - 3: Already logged in
    - 4: Waiting state
    - Other: Unknown error
  - Disconnection behavior:
    - Calls disconnected hook
    - Optionally auto-disconnects based on config
  - State management:
    - Resets network state to 1
    - Clears connection tries counter
    - Sets reconnect timeout

- connection_refused - Connection error handler (lines 1227-1231)
  - Displays error message when server denies connection
  - Parameters:
    - args->{error}: Numeric error code from server
  - Behavior:
    - Shows formatted error message to user
    - Uses translation function (TF) for localization


# Connection Related Handlers

- change_to_constate25 - Connection state change handler (lines 10111-10114)
  - Changes network state to 2.5
  - Unsets accountID
  - Simple implementation focused on state change
- quit_response - Disconnect request response handler (lines 9492-9499)
  - Processes disconnect request response notifications
  - Handles two result states:
    * fail=1: Cannot disconnect - "Please wait 10 seconds before trying to log out"
    * fail=0: Disconnect successful - "Logged out from the server succesfully"
  - Uses "error" message category for failures
  - Uses "success" message category for successful logout
  - Contains comments about result codes:
    * 0 = disconnect (quit) - DISCONNECTABLE_STATE
    * 1 = cannot disconnect (wait 10 seconds) - NOTDISCONNECTABLE_STATE
    * ? = ignored
  - Packet: 018B

- changeToInGameState() - Network state transition handler (lines 673-693)
  - Handles transitioning network state to IN_GAME or IN_GAME_BUT_UNINITIALIZED
  - Different behavior based on network version (version 1 vs others)
  - For version 1:
    - Sets IN_GAME state if accountID exists and character is initialized
    - Sets IN_GAME_BUT_UNINITIALIZED state otherwise
    - Sends welcome message if verbose config is enabled
  - For other versions, simply returns success (1)

- ping - Server ping handler (lines 12195-12199)
  - Processes ping notifications (0B1D)
  - Skips processing for XKore modes 1 and 3
  - Responds by sending ping back to server
  - Simple implementation for maintaining connection

- load_confirm - Client input permission handler (lines 12148-12153)
  - Processes client input permission notifications (0B01)
  - Outputs debug message about keyboard usage permission
  - Simple implementation with minimal functionality
  - Primarily relevant for ragexe client

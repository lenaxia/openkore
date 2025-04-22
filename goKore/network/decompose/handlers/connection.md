# Connection Related Handlers

**Method Implementations:**
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
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
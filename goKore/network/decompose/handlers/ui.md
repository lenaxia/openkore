**UI State Handlers:**

- map_loaded - Post-map load initialization (lines 1237-1287)
  - Network state transitions:
    - Sets network state to IN_GAME
    - Handles version-specific state changes
  - Character initialization:
    - Sets up character position and look direction
    - Initializes movement tracking
    - Sets initial status for private servers
  - Post-load actions:
    - Sends various initialization packets (map loaded, sync, guild info, cash items)
    - Calls in-game hook for plugins
    - Resets ignoreAll state
  - Server-specific behaviors:
    - Different handling for version 1 vs other versions
    - Special cases for bRO, idRO_Renewal, twRO
    - Private server status handling
  - User feedback:
    - Displays connection and coordinate messages
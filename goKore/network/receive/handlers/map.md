# Map Related Handlers

**Method Implementations:**
- map_change_cell - Map cell change handler (lines 10426-10429)
  - Processes map cell change notifications
  - Outputs debug message with:
    * Cell coordinates (x, y)
    * New cell type
    * Map name
  - Uses "info" debug category
  - Contains TODO comment: "add actual functionality, maybe alter field?"
  - Packet: 0192
- map_changed - Server change map handler (lines 7248-7361)
  - Processes map server change notifications (ZC_NPCACK_SERVERMOVE)
  - Sets network state to 4
  - Similar to map_change but for different map server
  - Handles map loading:
    * Extracts map name without instance ID
    * Checks if map is allowed
    * Creates new Field object if needed
  - Updates character position:
    * Sets pos and pos_to to new coordinates
    * Resets movement timers and solution path
  - Resets connection state variables
  - Initializes map change variables
  - Sets mapChanged flag for all AI sequences
  - Updates portalTrace timestamp
  - Processes server connection info:
    * Extracts map IP and port from URL or IP field
    * Displays map connection information
  - Disconnects from current map server (except version 1)
  - Resets item and skill timers:
    * Clears useSelf_item timers
    * Clears useSelf_skill timers
    * Clears doCommand timers
  - Resets character state:
    * Clears statuses, spirits, permitSkill, encoreSkill
    * Clears guild information
    * Closes and clears cart if active
  - Triggers map_changed hook with old map info
  - Updates AI timeout
  - Packets: 0092 (normal), 0AC7 (with DNS host)
- map_change - Map change/teleport handler (lines 7180-7243)
  - Processes map change notifications (ZC_NPCACK_MAPMOVE)
  - Requires in-game state (changeToInGameState)
  - Stops continuous skills if active
  - Handles map loading:
    * Extracts map name without instance ID
    * Checks if map is allowed
    * Creates new Field object if needed
  - Manages AI state:
    * Clears AI queue if configured
    * Initializes map change variables
    * Sets mapChanged flag for all AI sequences
    * Updates portalTrace timestamp
  - Updates character position:
    * Sets pos and pos_to to new coordinates
    * Resets movement timers and solution path
  - Displays map change message with coordinates
  - Handles network protocol differences:
    * Version 1: Suspends AI client
    * Other versions:
      - Sends mapLoaded notification
      - Sends blockingPlayerCancel if needed
      - Updates AI timeout
  - Triggers map_changed hook with old map info
  - Packets: 0091 (normal), 0A4B (airship)
- map_property3() - Map property and PvP mode handler (lines 4242-4265)
  - Processes map type information
  - Updates character status based on map type
  - Handles map property info table flags
  - Sets PvP mode based on map type:
    * Type 6: PvP mode (1)
    * Type 8: GvG mode (2)
    * Type 19: Battleground mode (3)
  - Triggers pvp_mode hook with mode information
  - Packet: 099B

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

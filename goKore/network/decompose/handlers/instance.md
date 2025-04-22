# Instance Dungeon Related Handlers

**Method Implementations:**
- instance_window_leave - Instance leave notification handler (lines 10507-10523)
  - Processes instance leave notifications
  - Handles multiple flag values:
    * 0: Update - "Received Memory Dungeon reservation update"
    * 1: Expired - "The Memorial Dungeon expired it has been destroyed"
    * 2: Entry time expired - "The Memorial Dungeon's entry time limit expired it has been destroyed"
    * 3: Removed - "The Memorial Dungeon has been removed"
    * 4: System error - "The instance windows has been removed, possibly due to party/guild leave"
    * Other: Unknown result with flag value
  - Contains detailed comments about flag meanings:
    * TYPE_NOTIFY = 0x0 - Will pop up Memory Dungeon Window
    * TYPE_DESTROY_LIVE_TIMEOUT = 0x1
    * TYPE_DESTROY_ENTER_TIMEOUT = 0x2
    * TYPE_DESTROY_USER_REQUEST = 0x3
    * TYPE_CREATE_FAIL = 0x4
  - Uses "debug" message category for flag 0
  - Uses "info" message category for flags 1-4
  - Uses "warning" message category for unknown flags
  - Contains TODO comment: "test if correct message displays, no type == 0 ?"
  - Packet: 02CE
- instance_window_join - Instance join notification handler (lines 10491-10496)
  - Processes instance join notifications
  - Outputs debug message with packet information
  - Triggers instance_ready hook
  - Contains TODO comment indicating incomplete implementation
  - Packet: 02CD
  - Simple implementation focused on debugging and hook triggering
- instance_window_queue - Instance queue notification handler (lines 10484-10487)
  - Processes instance queue notifications
  - Outputs debug message with packet information
  - Contains TODO comment indicating incomplete implementation
  - Contains comment: "To announce Instancing queue creation if no maps available"
  - Packet: 02CC
  - Simple implementation focused on debugging
- instance_window_start - Instance window start handler (lines 10476-10479)
  - Processes instance window start notifications
  - Outputs debug message with packet information
  - Contains TODO comment indicating incomplete implementation
  - Contains comment: "Required to start the instancing information window on Client"
  - Contains comment: "This window re-appear each 'refresh' of client automatically until 02CD is send to client"
  - Packet: 02CB
  - Simple implementation focused on debugging
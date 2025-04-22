**UI State Handlers:**

- attendance_ui - Attendance system UI handler (lines 11915-11959)
  - Processes attendance system UI notifications (0AE2 type=0x7)
  - Validates attendance period dates against current date
  - Parses attendance data:
    * Extracts already_requested flag
    * Calculates attendance_count
  - Displays formatted attendance information:
    * Period start/end dates
    * Current attendance day
    * List of daily rewards with status
  - Handles automatic attendance requests:
    * Checks attendanceAuto configuration
    * Runs 'attendance request' command when eligible
  - Error handling:
    * Warns if attendance_rewards.txt is outdated
    * Errors if attendance_rewards.txt doesn't exist

- action_ui - UI action handler (lines 11909-11913)
  - Processes UI action notifications
  - Outputs debug message about UI type
  - Contains comment about packet format (0AF0)
  - Contains comment about action types:
    * 0x0: close current UI
  - Simple implementation with minimal functionality

- open_ui - UI opening handler (lines 11878-11903)
  - Processes UI opening notifications
  - Outputs debug message about UI type
  - Handles multiple UI types:
    * BANK_UI (0x0): Bank UI
    * STYLIST_UI (0x1): Stylist UI
    * CAPTCHA_UI (0x2): Captcha UI
    * MACRO_UI (0x3): Macro Recorder UI
    * UI_UNUSED (0x4): Unused UI
    * TIPBOX_UI (0x5): Tip Box UI
    * RENEWQUEST_UI (0x6): Quest UI
    * ATTENDANCE_UI (0x7): Attendance UI
      - Calls attendance_ui handler
  - Displays appropriate message for each UI type
  - Displays error for unknown UI types
  - Contains TODO comments about implementing bank and stylist systems
  - Contains comment about packet format (0AE2)

- stylist_res - Stylist UI result handler (lines 11853-11861)
  - Processes stylist UI result notifications
  - Handles two res values:
    * true: Success - "[Stylist UI] Success"
      - Uses "info" message category
    * false: Failure - "[Stylist UI] Fail"
      - Uses error function
  - Simple implementation focused on notification

- progress_bar_unit - Progress bar unit handler (lines 11241-11244)
  - Processes progress bar unit notifications
  - Outputs debug message with GID and time values
  - Simple implementation focused on debugging

- font - Font usage handler (lines 10707-10710)
  - Processes font usage notifications
  - Outputs debug message with account ID and font ID
  - Uses "info" debug category
  - Contains TODO comment indicating incomplete implementation
  - Packet: 02EF
  - Simple implementation focused on debugging

- hotkeys - Hotkey list handler (lines 8032-8053)
  - Processes hotkey list notifications
  - Clears existing hotkeyList
  - Creates formatted display header with:
    * Column headers: #, Name, Type, Lv
    * Separator lines
  - Processes hotkey data in 7-byte chunks:
    * Unpacks type, ID, and level for each hotkey
    * Determines name based on type:
      - For skills: Gets name using Skill->new(idn)
      - For items: Uses itemNameSimple function
    * Formats each entry with:
      - Index number
      - Name (skill or item)
      - Type ("skill" or "item")
      - Level
  - Outputs formatted list with debug category "list"
  - Contains TODO comment about implementing rotate functionality
  - Comment notes this info is sent to xkore 2 clients

- refine_status - Refine result handler (lines 7938-7943)
  - Processes refine result notifications
  - Determines message based on status flag:
    * Success (status=1): Message index 3272
    * Failure (status=0): Message index 3273
  - Formats message with:
    * Item name (converted from bytes)
    * Refine level
    * Item name (simplified)
  - Displays warning message with result
  - Uses "info" message category
  - Simple implementation focused on result notification

- refineui_info - Refine UI item info handler (lines 7895-7936)
  - Processes refine UI item information (0AA2)
  - Handles two scenarios:
    * Valid item (len > 7):
      - Stores item index and blessing requirement
      - Gets item reference from inventory
      - Gets blessing item reference from inventory
      - Displays formatted information:
        * Target equipment details (index, name)
        * Blessing requirements (needed, owned)
        * Material options with:
          - Material ID
          - Success chance percentage
          - Zeny cost
          - Current inventory count
      - Shows continuation instructions with command syntax
    * Invalid item (len <= 7):
      - Displays error message about incompatible equipment
      - Suggests 'i' command to check equipment
  - Uses complex unpacking for materials list
  - Uses "info" message category for success, "error" for failures
  - Author attribution: [Cydh] (from packet comment)

- refineui_opened - Refine UI opening handler (lines 7885-7889)
  - Processes refine UI opening notifications (0AA0)
  - Displays instructional message to user:
    * Explains how to check equipment ('i' command)
    * Provides command syntax for continuing (refineui select [ItemIdx])
  - Sets refineUI->{open} flag to 1
  - Uses "info" message category
  - Simple implementation focused on user guidance
  - Author attribution: [Cydh]

- progress_bar_stop - Progress bar completion handler (lines 4590-4593)
  - Processes progress bar completion notifications
  - Displays completion message to user
  - Simple handler for server-initiated progress bar termination

- progress_bar - Progress bar display handler (lines 4576-4588)
  - Processes progress bar display requests
  - Shows loading message with time information
  - Sets character progress_bar flag
  - Creates task chain to:
    * Wait for specified time
    * Send progress completion message
    * Display completion message
    * Reset progress_bar flag
  - Uses TaskManager for timing and callbacks

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
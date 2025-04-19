**Hotkey Handler:**
- hotkeys() - Manages hotkey configurations (lines 8032-8053)
  - Processes and displays hotkey information
  - Formats output in a table with columns: #, Name, Type, Lv
  - Handles both skill and item hotkeys
  - Supports debug logging of hotkey list

**UI System Handlers:**
- ui_notification (lines 6457-6470)
  - Manages in-game notifications
  - Processes:
    - Notification types
    - Message content
    - Display duration
    - Priority levels
  - Maintains:
    - @notificationQueue array
    - %activeNotifications hash

- ui_window (lines 6472-6480)
  - Handles UI window management

- ui_interaction (lines 6482-6495)
  - Processes UI element interactions
  - Handles:
    - Button clicks
    - Input field changes
    - Dropdown selections
  - Features:
    - Triggers 'ui_interaction' hook
    - Maintains focus states
    - Processes validation rules
  - Features:
    - Processes window open/close events
    - Updates window positions
    - Triggers 'window_state_change' hook
    - Manages focus states
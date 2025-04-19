# Quest System Handlers

**Quest Mission Handlers:**
- quest_all_mission (lines 4699-4754)
  - Processes mission packets (02B2)
  - Handles variable length quest+mission data
  - Maintains questList structure with:
    - Quest ID, start/expire times
    - Mission details (mob ID, count, name)
  - Triggers quest_mission_added hook

- quest_add (lines 4760-4823)
  - Handles single quest addition (02B3/09F9/0B0C)
  - Processes different packet versions:
    - 02B3: Basic mob hunt format
    - 09F9: Extended hunt details
    - 0B0C: Additional hunt ID fields
  - Maintains quest state (active/inactive)
  - Shows quest addition messages

**Quest Update Handlers:**
- quest_update_mission_hunt (lines 4828-4925)
  - Processes hunt progress updates (02B5/09FA/0AFE)
  - Handles multiple packet versions:
    - 02B5: Basic mob ID tracking
    - 09FA: Hunt identification system
    - 0AFE: Extended hunt IDs
  - Updates mob kill counts
  - Shows progress messages based on config

**Quest Management Handlers:**
- quest_delete (lines 4928-4934)
  - Removes quest from questList
  - Shows deletion message
  - Triggers quest_delete hook

- quest_active (lines 4937-4948)
  - Toggles quest active state
  - Shows activation/deactivation message
  - Triggers quest_active hook

**Quest List Handler:**
- quest_all_list (lines 4599-4655)
  - Processes different quest packet versions:
    - 02B1: Basic format (quest_id, active)
    - 097A: Added time tracking (20141022+)
    - 09F8: Added hunt details (20150513+)
    - 0AFF: Extended hunt IDs (20181010+)
  - Handles packet structure differences:
    - Basic: 5 bytes per quest
    - Extended: 15 bytes header + mission data
  - Maintains questList structure
  - Supports multi-packet quest lists
  - Uses generation counter for map changes
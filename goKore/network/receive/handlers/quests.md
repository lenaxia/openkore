# Quest Related Handlers

**Method Implementations:**
- quest_active - Quest activation handler (lines 4937-4948)
  - Processes quest activation/deactivation notifications
  - Displays appropriate message based on active state:
    * Active: "Quest X is now active"
    * Inactive: "Quest X is now inactive"
  - Updates quest active status in questList
  - Includes quest title lookup for better readability
  - Triggers quest_active hook
  - Packet: 02B7
- quest_delete - Quest deletion handler (lines 4928-4934)
  - Processes quest deletion notifications
  - Displays deletion message with quest title lookup
  - Removes quest from questList global structure
  - Triggers quest_delete hook
  - Packet: 02B4
- quest_update_mission_hunt - Quest progress update handler (lines 4827-4925)
  - Processes quest mission progress updates
  - Supports multiple packet versions:
    * 02B5/08FE: Basic format with mob_id
    * 09FA: Extended format with hunt_id (SERVERTYPE >= 20150513)
    * 0AFE: Advanced format with continuous hunt IDs (SERVERTYPE >= 20181010)
  - Handles complex mission ID resolution:
    * Matches hunt_id to hunt_id
    * Matches mob_id to mob_id
    * Cross-matches hunt_id to mob_id and vice versa
  - Updates mission count and goal values
  - Displays progress based on questDisplayStyle config:
    * Style 1: Simple "[mob_name] [count/goal]"
    * Style 2+: Detailed with quest title and ID
  - Triggers quest_mission_updated and quest_update_mission_hunt_end hooks
  - Packets: 02B5, 08FE, 09FA, 0AFE
- quest_add - Quest addition handler (lines 4759-4823)
  - Processes new quest addition packets
  - Supports multiple packet versions:
    * 02B3: Default format with basic mission info
    * 09F9: Extended format with hunt details (SERVERTYPE >= 20150513)
    * 0B0C: Advanced format with continuous hunt IDs (SERVERTYPE >= 20181010)
  - Updates questList with new quest information:
    * Quest ID, active status, start/expire times
    * Mission amount and details
  - Handles up to 3 missions per quest
  - Displays quest addition message with title lookup
  - Converts mob names from bytes to string
  - Triggers quest_mission_added and quest_added hooks
- quest_all_mission - Quest mission details handler (lines 4698-4754)
  - Processes detailed mission information for quests
  - Uses fixed format for quest and mission data:
    * Quest data: quest_id, time_start, time_expire, mission_amount
    * Mission data: mob_id, mob_count, mob_name
  - Updates existing quest entries in questList
  - Handles up to 3 missions per quest
  - Converts mob names from bytes to string
  - Provides detailed debug logging
  - Triggers quest_mission_added and quest_all_mission_end hooks
  - Packet: 02B2
- quest_all_list - Quest list handler (lines 4599-4694)
  - Processes complete quest list packets
  - Supports multiple packet versions:
    * 02B1: Basic quest list (quest ID and active status)
    * 097A: Extended quest list with time and mission info (SERVERTYPE >= 20141022)
    * 09F8: Enhanced quest list with hunt details (SERVERTYPE >= 20150513)
    * 0AFF: Advanced quest list with continuous hunt IDs (SERVERTYPE >= 20181010)
  - Handles multi-packet quest lists across map changes
  - Maintains quest generation tracking
  - Unpacks quest and mission data
  - Updates questList global structure
  - Processes mission details:
    * Mob ID, name, count, and goal
    * Hunt ID and type (for newer versions)
    * Level requirements (for newer versions)
  - Triggers quest_mission_added and quest_all_list_end hooks
  - Provides detailed debug logging
# Taekwon Class Related Handlers

**Method Implementations:**
- taekwon_rank - Taekwon mission rank update handler (lines 10966-10969)
  - Processes Taekwon mission rank updates
  - Displays message with rank value
  - Uses "info" message category
  - Contains comment: "Updates the fame rank points for the Taekwon ranking"
  - Contains comment about packet format (0224)
  - Simple implementation focused on notification
- taekwon_packets - Taekwon class special packet handler (lines 10944-10962)
  - Processes Taekwon class special notifications
  - Determines string value based on value parameter:
    * 1: "Sun"
    * 2: "Moon"
    * 3: "Stars"
    * Other: "Unknown (X)"
  - Handles multiple flag values:
    * 0: Map registration - "You have now marked: X as Place of the Y"
    * 1: Map information - "X is marked as Place of the Y"
    * 10: Hate mob registration - "You have now marked X as Target of the Y"
    * 11: Hate mob information - "X is marked as Target of the Y"
    * 20: TK_MISSION target - "[TaeKwon Mission] Target Monster : X (Y%)"
    * 30: Reset - "Your Hate and Feel targets have been resetted"
    * Other: Unknown result with flag value
  - Uses bytesToString to convert names
  - Uses "info" message category for most messages
  - Uses no category for flag 11 messages
  - Uses "warning" message category for unknown flags
  - Contains TODO comment: "test if we must use ID to know if the packets are meant for us"
  - Comment notes that ID is monsterID
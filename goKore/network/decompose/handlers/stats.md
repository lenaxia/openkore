**Character Stat Handlers:**

- repute_info() - Character reputation info handler (lines 12387-12435)
  - Handles character reputation updates
  * Processes ZC_REPUTE_INFO packets
  * Updates fame and reputation values
  * Displays reputation change notifications

- hp_sp_changed (lines 7157-7170)
  - Manages HP/SP change notifications
  - Processes ZC_HP_SP_CHANGED packets
  - Features:
    - Handles both HP (type 5) and SP (type 7) changes
    - Maintains stat boundaries (max HP/SP)
    - Supports direct stat modifications
    - Works with in-game state validation

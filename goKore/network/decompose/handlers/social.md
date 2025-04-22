**Social Interaction Handlers:**

- actor_info() - Handles actor information updates (lines 2712-2804)
  - Processes different actor types:
    * Players:
      - Updates name, party/guild info, titles
      - Calls updatePlayerNameCache()
      - Triggers charNameUpdate hook
    * Monsters:
      - Updates monster name and info
      - Updates monster LUT if new
      - Triggers mobNameUpdate hook
    * NPCs:
      - Updates NPC name and info
      - Updates NPC LUT if new
      - Triggers npcNameUpdate hook
    * Pets:
      - Updates pet name and info
      - Triggers petNameUpdate hook
    * Slaves:
      - Updates slave name and info
      - Triggers slaveNameUpdate hook
    * Elementals:
      - Updates elemental name and info
      - Triggers elementalNameUpdate hook
  - Uses bytesToString() for string conversion
  - Calls Plugins::callHook() for various update types
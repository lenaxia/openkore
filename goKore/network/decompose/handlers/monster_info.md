# Monster Information Related Handlers

**Method Implementations:**
- sense_result - Monster sense skill result handler (lines 11427-11446)
  - Processes monster sense skill result notifications
  - Defines lookup tables for race and size:
    * race_lut: Formless, Undead, Beast, Plant, etc.
    * size_lut: Small, Medium, Large
  - Displays formatted message with monster information:
    * Monster name and level
    * Size and race
    * Defense and magic defense
    * Element type and HP
    * Damage modifiers for elements:
      - Ice, Earth, Fire, Wind
      - Poison, Holy, Dark, Spirit
      - Undead
  - Uses monsters_lut to get monster name from nameID
  - Uses elements_lut to get element name
  - Uses "list" message category
  - Contains comment about parameter order
# Cooking Related Handlers

**Method Implementations:**
- cooking_list - Cooking recipe list handler (lines 9137-9159)
  - Processes list of available cooking recipes
  - Clears previous cooking list and type
  - Sets current cooking type from packet
  - Creates formatted header with centered title
  - Iterates through raw message data in 2-byte chunks
  - For each recipe:
    * Extracts item nameID
    * Stores nameID in cookingList array
    * Displays item index and name
  - Adds footer with separator line
  - Shows instruction for using 'cook' command
  - Uses "list" message category for recipe list
  - Uses "info" message category for instruction
  - Triggers cooking_list hook with:
    * cooking_list: Array of recipe item IDs
  - Packet: 025A
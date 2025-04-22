# Inventory Related Handlers

**Method Implementations:**
- inventory_item_added - Item addition handler (lines 3679-3753)
  - Handles adding new items to inventory
  - Manages both new items and stackable additions
  - Processes item properties (ID, name, type, amount, etc.)
  - Handles failure cases (inventory full, frozen, etc.)
  - Triggers item_gathered hook
  - Implements auto-drop logic
  - Packets: 00A0, 029A, 02D4, 0990, 0A0C, 0A37

- inventory_items_stackable - Stackable items handler (lines 5210-5230)
  - Handles adding stackable items to inventory
  - Uses _items_list helper method
  - Processes item additions through inventory->add()
  - Special handling for arrow equipment
  - Triggers packet_inventory hook
  - Packets: 00A3, 01EE, 02E8, 0900, 0991

- inventory_item_removed - Item removal handler (lines 3756-3778)
  - Handles item removal from inventory
  - Processes different removal reasons:
    - Skill casting (1)
    - Refinement failure (2)
    - Chemical reaction (3)
    - Storage transfer (4)
    - Cart transfer (5)
    - Selling (6)
    - Four Spirit Analysis skill (7)
    - Unknown reasons (default)
  - Packets: 00AF, 07FA

- inventory_items_nonstackable - Non-stackable items handler (lines 1101-1126)
  - Handles adding non-stackable items to inventory
  - Manages equipment updates
  - Uses _items_list helper method
  - Processes item additions through inventory->add()
  - Handles equipment slot updates
  - Triggers packet_inventory hook
  - Packets: 00A4, 0295, 02D0, 0992, 0A0D

**Method Implementations:**
- inventory_items_nonstackable - Non-stackable items handler (lines 1101-1126)
  - Handles adding non-stackable items to inventory
  - Manages equipment updates
  - Uses _items_list helper method
  - Processes item additions through inventory->add()
  - Handles equipment slot updates
  - Triggers packet_inventory hook
  - Packets: 00A4, 0295, 02D0, 0992, 0A0D
# Item Related Handlers

**Method Implementations:**
- arrow_equipped - Arrow/Bullet equipment handler (lines 3653-3670)
  - Handles equipping arrows/bullets to character
  - Updates character's arrow equipment slot
  - Triggers equipment hooks
  - Displays message when equipped
  - Packet: 013C

- inventory_item_added - Item pickup handler (lines 3679-3753)
  - Handles adding items to inventory from various sources
  - Supports multiple packet versions (00A0, 029A, 02D4, etc)
  - Manages both new items and stackable items
  - Handles pickup failure cases (inventory full, frozen, etc)
  - Triggers item_gathered hook
  - Supports auto-drop functionality

- inventory_item_removed - Item removal handler (lines 3756-3786)
  - Handles item removal from inventory
  - Supports multiple removal reasons (used, broken, sold, etc)
  - Triggers packet_item_removed hook
  - Provides debug messages for different removal cases
  - Packets: 00AF, 07FA

- rental_expired - Rental item expiration handler (lines 3789-3801)
  - Handles rental item expiration notifications
  - Removes expired item from inventory
  - Triggers rental_expired hook
  - Displays expiration message
  - Packet: 0299
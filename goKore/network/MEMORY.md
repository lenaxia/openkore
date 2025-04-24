## Task
We have now decomposed Receive.pm into a bunch of separate files in preparation for a re-implemnetation in golang. But we need to ensure that the structure that we're going to reimplement it into follows golang best practices. This implementation is intended to provide a Receive golang module that will be used in a larger go program, so please keep that in mind.

We will do this in three steps. Make sure to regularly update MEMORY.md. Make sure to mark step 1 and step 2 boxes as complete `[X]` when the file is done being reviewed. Review the sections below in case files have been reviewed but not checked off.

Step 1: Review each md file one by one and documemt the purpose and content in the `Receive Descriptions` section below. Do these in batches, read 5-7 files then write descriptions to MEMORY.md, don't do too many or else context will get too long.

Step 2: Review each md file again, but this time verify each of the objects (methods, constants, types, etc) in each markdown file. Answer each of the two following questions: 1/ Does this object belong in this file or is there a better place for it? 2/ If there is a better place for it, where should it be? Document this migration in the `Receive Migrations` section below. DO NOT MAKE ANY CHANGES TO MD FILES YET. You can do these in batches of 2-3 files depending on the length of the file then write updates to MEMORY.md, but keep it shorter for large files to prevent long context. If the file is longer than 500 lines, iteratively review sections of the file until you have parsed the entire contents. Only review complete sections. If a section seems to be partial or incomplete, ignore it and ensure that the entire block is reviewed when fetching the next batch of lines.

Step 3: Based on the file descriptions recommend any changes to the folder structure that better follows goland best practices. This may mean consolidating files.

## File List

### common/
 1   2
[X] [X] account.md
[X] [X] connection.md
[X] [X] constants.md
[X] [X] network.md

### security/
 1   2
[X] [X] anti_cheat.md
[X] [X] encryption.md
[X] [X] login.md
[X] [X] pin.md

### handlers/
 1   2
[X] [X] account.md
[X] [X] account_rates.md
[X] [X] achievements.md
[X] [X] actors.md
[X] [X] adoption.md
[X] [X] auction.md
[X] [X] banking.md
[X] [X] battleground.md
[X] [X] booking.md
[X] [X] books.md
[X] [X] boss.md
[X] [X] buying_store.md
[X] [X] captcha.md
[X] [X] card_merge.md
[X] [X] cart.md
[X] [X] character.md
[X] [X] characters.md
[X] [X] clan.md
[X] [X] combat.md
[X] [X] config.md
[X] [X] connection.md
[X] [X] cooking.md
[X] [X] deaths.md
[X] [X] effects.md
[X] [X] elemental.md
[X] [X] emblems.md
[X] [X] equipment.md
[X] [X] experience.md
[X] [X] gm.md
[X] [X] guild.md
[X] [X] homunculus.md
[X] [X] instance.md
[X] [X] inventory.md
[X] [X] item_merge.md
[X] [X] items_list.md
[X] [X] items.md
[X] [X] mail.md
[X] [X] map.md
[X] [X] marriage.md
[X] [X] memo.md
[X] [X] mercenary.md
[X] [X] messaging.md
[X] [X] minimap.md
[X] [X] misc.md
[X] [X] monster_info.md
[X] [X] movement.md
[X] [X] mvp.md
[X] [X] navigation.md
[X] [X] npc.md
[X] [X] party.md
[X] [X] pc_cafe.md
[X] [X] pets.md
[X] [X] portals.md
[X] [X] pvp.md
[X] [X] quests.md
[X] [X] ranking.md
[X] [X] refine.md
[X] [X] rental.md
[X] [X] repair.md
[X] [X] reputation.md
[X] [X] roulette.md
[X] [X] script.md
[X] [X] search_store.md
[X] [X] shop.md
[X] [X] skills.md
[X] [X] social.md
[X] [X] stats.md
[X] [X] status.md
[X] [X] storage.md
[X] [X] system.md
[X] [X] system_misc.md
[X] [X] taekwon.md
[X] [X] trade.md
[X] [X] transportation.md
[X] [X] ui.md
[X] [X] vending.md


# Receive Descriptions

## Common Directory

- common/account.md: Contains handlers for account server information processing. Includes methods for parsing and reconstructing server info packets for different server types, and handling account information from login servers.

- common/connection.md: Contains handlers for connection-related events, including map server connection errors and connection refusal. Manages error codes, disconnection behavior, and state management.

- common/constants.md: Contains constant definitions used across the network module. Currently only includes homunculus state constants used with the homunculus_info handler. Contains references to other blocks of constances in src/Network/Receive.pm.

- common/network.md: Contains core network functionality, specifically the parse() method which wraps the parent class's parse method and adds debug logging capabilities for packet inspection.

## Security Directory

- security/anti_cheat.md: Contains handlers for various anti-cheat systems including Hack Shield, GameGuard, and Easy Anti-Cheat. Manages detection responses, verification requests, and graceful exits when unsupported protection is detected.

- security/encryption.md: Contains handlers for message ID encryption initialization. Sets up complex encryption using nibble extraction, XOR operations, and bit shifting to establish secure communication.

- security/login.md: Contains comprehensive login-related handlers including secure login key processing, character data management, login error handling, character creation, and character slot management.

- security/pin.md: Contains handlers for PIN code validation, requests, and management. Handles various PIN states, validates PIN format, and manages PIN storage and reuse.

## Handlers Directory

- handlers/account.md: Contains handlers for account-related information, including displaying account ID in debug logs and formatting account payment information (subscription time remaining in days/hours/minutes).

- handlers/account_rates.md: Contains handlers for account rate information, including detailed rate information for experience, death penalties, and drop rates, as well as premium account rate bonuses.

- handlers/achievements.md: Contains handlers for the achievement system, including listing achievements, updating achievement progress, and acknowledging achievement rewards.

- handlers/actors.md: Contains comprehensive handlers for actor management, including appearance changes, direction changes, HP information, and character data formatting. Handles various actor types (players, monsters, NPCs) and their display properties.

- handlers/adoption.md: Contains handlers for the adoption system, processing adoption requests and replies with appropriate validation messages.

- handlers/auction.md: Contains handlers for the auction system, including auction results, item addition, window status, and auction ending with appropriate success/failure messages.

- handlers/banking.md: Contains handlers for the banking system, including account balance checking, deposits, and withdrawals with appropriate success/failure messages and balance updates.

- handlers/battleground.md: Contains stub handlers for battleground-related functionality, including emblem and message notifications. Currently focused on debugging with TODO comments.

- handlers/booking.md: Contains handlers for the party booking system, including creation, deletion, updating, and searching for booking entries with appropriate success/failure messages.

- handlers/books.md: Contains a simple handler for book reading notifications, currently only logging debug information with a TODO comment about adding table file support.

- handlers/boss.md: Contains handlers for boss monster information, including location tracking, detection notifications, and respawn timers with appropriate messages for each state.

- handlers/buying_store.md: Contains comprehensive handlers for the buying store system, including store creation, item listing, store discovery, item purchases, and various failure conditions with appropriate messages.

- handlers/captcha.md: Contains extensive handlers for CAPTCHA and macro detection systems, including image processing, answer validation, upload/download functionality, and macro reporter features.

- handlers/card_merge.md: Contains handlers for the card merging system, including listing mergeable items and processing merge results with appropriate success/failure messages and inventory updates.

- handlers/cart.md: Contains handlers for cart-related functionality, including item addition/removal, weight/capacity tracking, and handling of stackable and non-stackable items with appropriate messages.

- handlers/character.md: Contains handlers for character deletion operations, including deletion results, deletion acceptance, and deletion cancellation with appropriate success/failure messages and state management.

- handlers/characters.md: Contains handlers for character-related operations, including character switching, character ID and map information, character name updates, character deletion, and equipment display.

- handlers/clan.md: Contains handlers for the clan system, including clan information updates, clan chat messages, clan member counts, and clan leaving operations with appropriate data structure management.

- handlers/combat.md: Contains handlers for combat-related functionality, including monster ranged attacks, combo delays, and attack range updates with appropriate debug messages and configuration adjustments.

- handlers/config.md: Contains handlers for player configuration settings, including equipment visibility, summoning permissions, and pet/homunculus autofeeding with appropriate feedback messages.

- handlers/connection.md: Contains handlers for connection state changes and disconnect request responses, including state management and appropriate success/failure messages.

- handlers/cooking.md: Contains handlers for the cooking system, including recipe list processing, displaying available recipes, and providing instructions for using the cooking command.

- handlers/deaths.md: Contains handlers for entity death and disappearance events, handling different types of disappearances (out of sight, death, logout, teleport) for various entity types.

- handlers/effects.md: Contains handlers for various visual and sound effects, including level up effects, skill effects, area spells, hat effects, sound effects, and minimap indicators.

- handlers/elemental.md: Contains handlers for elemental information updates, creating or updating elemental actors with data from server packets.

- handlers/emblems.md: Contains a stub handler for character emblem updates, currently only logging debug information with a TODO comment indicating incomplete implementation.

- handlers/equipment.md: Contains handlers for equipment-related operations, including item equipping and unequipping for both regular equipment and equipment switch functionality.

- handlers/experience.md: Contains handlers for experience point notifications, handling both base and job experience from battle and quest sources with appropriate formatting.

- handlers/gm.md: Contains handlers for Game Master interactions, including silence status changes and account name request responses with appropriate notification messages.

- handlers/guild.md: Contains comprehensive handlers for guild-related functionality, including member management, position settings, alliances, notices, storage logs, and emblem handling.

- handlers/homunculus.md: Contains handlers for homunculus-related functionality, including property calculation, feeding, state management, and property updates with appropriate state transitions and messages.

- handlers/instance.md: Contains handlers for instance dungeon functionality, including window start, queue notifications, join notifications, and leave notifications with appropriate status messages.

- handlers/inventory.md: Contains handlers for inventory-related operations, including item addition/removal, stackable and non-stackable item management, cart operations, rental item expiration, and favorite item status.

- handlers/item_merge.md: Contains handlers for the item merging system, including merge list processing and merge result handling with appropriate success/failure messages and inventory updates.
- handlers/items.md: Contains handlers for item-related operations, specifically item usage notifications with appropriate inventory updates and messages.

- handlers/items_list.md: Contains handlers for managing item lists across different container types (inventory, cart, storage), including initialization, processing of stackable and non-stackable items, and finalization.

- handlers/mail.md: Contains comprehensive handlers for the mail system, including sending, receiving, reading, deleting mail, and managing attachments with appropriate success/failure messages and formatted displays.

- handlers/map.md: Contains handlers for map-related operations, including map changes, server changes, cell type changes, and property updates with appropriate state management and hooks.

- handlers/marriage.md: Contains handlers for marriage-related notifications, specifically divorce notifications with appropriate messages.
- handlers/memo.md: Contains handlers for memo (save point) operations, processing success/failure results with appropriate messages and hooks.

- handlers/mercenary.md: Contains handlers for mercenary-related functionality, including initialization and removal with appropriate property management and configuration adjustments.

- handlers/messaging.md: Contains comprehensive handlers for various messaging systems, including private messages, system messages, chat, broadcasts, whispers, and the Rodex mail system with appropriate formatting and hooks.

- handlers/minimap.md: Contains handlers for minimap-related functionality, including indicator processing with actor references and color management, with a placeholder for reconstruction functionality.

- handlers/misc.md: Contains miscellaneous handlers that don't fit elsewhere, including online user count notifications, remaining time information, and placeholder functions.
- handlers/monster_info.md: Contains handlers for monster information, specifically the monster sense skill result handler that displays formatted monster statistics including level, size, race, defenses, and elemental properties.

- handlers/movement.md: Contains handlers for character and actor movement, including movement interruption, teleport failures, high jumps, and character movement with appropriate position tracking and state management.

- handlers/mvp.md: Contains handlers for MVP (Most Valuable Player) related notifications, including player MVP status, other player MVP status, and MVP item rewards with appropriate messages.

- handlers/navigation.md: Contains handlers for the navigation system, processing server navigation requests for monster hunting and location navigation with appropriate hooks and messages.

- handlers/npc.md: Contains comprehensive handlers for NPC interactions, including dialog management, text/number input, menu responses, chat messages, and image display with appropriate state tracking and hooks.
- handlers/party.md: Contains extensive handlers for party-related functionality, including member management, invitation handling, chat messages, experience settings, HP updates, and location tracking with appropriate data structures and messages.

- handlers/pc_cafe.md: Contains handlers for PC cafe point system, currently only implementing debug logging for point notifications with a TODO comment indicating incomplete implementation.

- handlers/pets.md: Contains handlers for pet-related functionality, including pet information updates, feeding, evolution, emotion display, capture results, and egg hatching with appropriate data structure management.

- handlers/portals.md: Contains handlers for warp portal functionality, specifically the warp portal list handler that processes memo points, updates save map configuration, and displays formatted portal lists.

- handlers/pvp.md: Contains handlers for PvP (Player vs Player) functionality, specifically the rank update handler that tracks and displays the player's current PvP ranking.
- handlers/quests.md: Contains comprehensive handlers for the quest system, including quest list management, mission tracking, progress updates, and quest activation/deletion with appropriate hooks and messages.

- handlers/ranking.md: Contains handlers for various ranking systems, including top 10 rankings for blacksmiths, alchemists, taekwon fighters, and PK players, as well as ranking points updates with appropriate formatting and display.

- handlers/refine.md: Contains handlers for item refining and upgrading, including refinable item listing, refine/craft results, and weapon upgrade results with appropriate success/failure messages.

- handlers/rental.md: Contains handlers for item rental functionality, specifically the rental time notification handler that displays remaining rental time for items.

- handlers/repair.md: Contains handlers for item repair functionality, including repair item listing and repair results with appropriate success/failure messages and inventory management.
- handlers/reputation.md: Contains handlers for the reputation system, specifically the reputation information handler that processes and stores reputation data in the global reputation list.

- handlers/roulette.md: Contains handlers for the roulette game system, including window opening/updating, rewards information, and item reward notifications with appropriate data structure management and messages.

- handlers/script.md: Contains handlers for NPC script messages, processing message content and NPC identification with appropriate hooks and debug messages.

- handlers/search_store.md: Contains handlers for the item search store system, including search results processing, position tracking, failure handling, and store opening with appropriate data structure management.

- handlers/shop.md: Contains comprehensive handlers for various shop systems, including NPC shops, player vending, cash shops, and market systems with appropriate item listing, purchase/sale handling, and result notifications.
- handlers/skills.md: Contains extensive handlers for skill-related functionality, including skill casting, cooldowns, failures, autospell management, and gospel buff effects with appropriate state tracking and hooks.

- handlers/social.md: Contains comprehensive handlers for social interactions, including chat rooms, friends list, emotions, marriage notifications, and ignore player functionality with appropriate data structure management.

- handlers/stats.md: Contains handlers for character statistics, including stat point allocation, status information updates, and parameter changes for both base and derived stats with appropriate hooks.

- handlers/status.md: Contains handlers for character status effects, including status activation and resurrection with appropriate state management and messages.

- handlers/storage.md: Contains handlers for storage-related functionality, including item addition/removal, password management, and storage opening/closing with appropriate data structure management.
- handlers/system.md: Contains handlers for system-level functionality, specifically the network state transition handler that manages game state changes based on network version.

- handlers/system_misc.md: Contains miscellaneous system handlers, including client input permission, inventory expansion, item preview, server ping, and Star Gladiator map confirmation with appropriate messages.

- handlers/taekwon.md: Contains handlers for Taekwon class-specific functionality, including mission rank updates and special packet handling for celestial/hate target registration with appropriate messages.

- handlers/trade.md: Contains comprehensive handlers for the player trading system, including trade requests, item addition, finalization, completion, and cancellation with appropriate data structure management.

- handlers/transportation.md: Contains handlers for transportation systems, specifically the private airship usage handler with appropriate success/failure messages.
- handlers/ui.md: Contains extensive handlers for user interface elements, including attendance system, stylist, refine UI, progress bars, hotkeys, and map loading with appropriate state management and messages.

- handlers/vending.md: Contains handlers for the vending system, including store setup results, vender discovery, and vender removal with appropriate data structure management and hooks.











# Step 3: Recommended Go Structure

Based on the review of the files, here's a recommended structure that better follows Go best practices:

## Principles Applied
1. Organize packages around functionality, not types
2. Use short, concise, and descriptive package names
3. Ensure each package has a single, well-defined purpose
4. Avoid deeply nested package hierarchies
5. Follow Go's standard library conventions

## Recommended Structure

```
network/
├── receive/
│   ├── core/           # Core network functionality (replaces common/)
│   │   ├── parse.go    # Network parsing functionality
│   │   └── account.go  # Account-related core functionality
│   ├── security/       # Security-related functionality
│   │   ├── login.go    # Login and authentication
│   │   ├── pin.go      # PIN code handling
│   │   └── anticheat.go # Anti-cheat systems
│   ├── game/           # Game feature handlers (replaces handlers/)
│   │   ├── actor/      # Actor-related functionality
│   │   │   ├── player.go
│   │   │   ├── monster.go
│   │   │   └── npc.go
│   │   ├── item/       # Item-related functionality
│   │   │   ├── inventory.go
│   │   │   ├── storage.go
│   │   │   └── equipment.go
│   │   ├── social/     # Social features
│   │   │   ├── chat.go
│   │   │   ├── guild.go
│   │   │   ├── party.go
│   │   │   └── friends.go
│   │   ├── economy/    # Economic features
│   │   │   ├── shop.go
│   │   │   ├── vending.go
│   │   │   ├── auction.go
│   │   │   └── banking.go
│   │   └── world/      # World interaction
│   │       ├── map.go
│   │       ├── movement.go
│   │       └── npc.go
│   └── types/          # Shared types and constants
│       ├── constants.go
│       └── packets.go
└── send/               # For sending packets (future implementation)
```

## Key Improvements

1. **Logical Grouping**: The structure groups related functionality together, making it easier to navigate and understand the codebase.

2. **Reduced File Count**: Instead of having 70+ separate files, related handlers are consolidated into logical files based on functionality.

3. **Better Dependency Management**: The structure makes dependencies clearer and reduces circular dependencies.

4. **Go Idioms**: Follows Go's convention of organizing by functionality rather than by type.

5. **Scalability**: The structure can easily accommodate new features without becoming unwieldy.

6. **Testability**: The organization makes it easier to write focused tests for each component.

## Migration Strategy

1. Start by creating the new directory structure
2. Move the homunculus constants from common/constants.md to game/actor/homunculus.go
3. Consolidate the account-related functionality from various files into core/account.go
4. Move character-related functionality from security/login.md to game/actor/player.go
5. Group related handlers into their respective files based on functionality
6. Ensure each file has a clear, single responsibility

This structure will provide a solid foundation for implementing the Receive module in Go while following best practices.

# Receive Migrations
Format Example:
```
[ ] source_file.md:some_method_name:line_number - destination_file.md
[ ] reputation.md:repair_weapon:56 - items.md
```

## Common Directory

[ ] common/constants.md:HO_PRE_INIT:5 - handlers/homunculus.md
[ ] common/constants.md:HO_RELATIONSHIP_CHANGED:6 - handlers/homunculus.md
[ ] common/constants.md:HO_FULLNESS_CHANGED:7 - handlers/homunculus.md
[ ] common/constants.md:HO_ACCESSORY_CHANGED:8 - handlers/homunculus.md
[ ] common/constants.md:HO_HEADTYPE_CHANGED:9 - handlers/homunculus.md

## Security Directory

[ ] security/pin.md:queryLoginPinCode:26 - common/account.md
[ ] security/pin.md:queryAndSaveLoginPinCode:35 - common/account.md

## Handlers Directory

[ ] handlers/account.md:account_id:3 - common/account.md

[ ] handlers/actors.md:received_characters_blockSize:83 - security/login.md
[ ] handlers/actors.md:received_characters_unpackString:90 - security/login.md

[ ] handlers/character.md:char_delete2_result:2 - security/login.md
[ ] handlers/character.md:char_delete2_accept_result:16 - security/login.md
[ ] handlers/character.md:char_delete2_cancel_result:33 - security/login.md

[ ] handlers/characters.md:switch_character:4 - security/login.md
[ ] handlers/characters.md:received_character_ID_and_Map:17 - security/login.md
[ ] handlers/characters.md:character_name:43 - security/login.md
[ ] handlers/characters.md:character_deletion_failed:51 - security/login.md
[ ] handlers/characters.md:character_deletion_successful:62 - security/login.md
[ ] handlers/characters.md:show_eq:77 - security/login.md

[ ] handlers/connection.md:change_to_constate25:4 - common/connection.md

[ ] handlers/deaths.md:actor_died_or_disappeared:2 - handlers/actors.md

[ ] handlers/effects.md:minimap_indicator:90 - handlers/minimap.md

[ ] handlers/emblems.md:char_emblem_update:4 - handlers/guild.md

[ ] handlers/inventory.md:cart_off:14 - handlers/cart.md

[ ] handlers/items_list.md:item_list_start:42 - handlers/inventory.md (inventory part)
[ ] handlers/items_list.md:item_list_start:42 - handlers/cart.md (cart part)
[ ] handlers/items_list.md:item_list_start:42 - handlers/storage.md (storage part)

[ ] handlers/items_list.md:item_list_stackable:29 - handlers/inventory.md (inventory part)
[ ] handlers/items_list.md:item_list_stackable:29 - handlers/cart.md (cart part)
[ ] handlers/items_list.md:item_list_stackable:29 - handlers/storage.md (storage part)

[ ] handlers/items_list.md:item_list_nonstackable:14 - handlers/inventory.md (inventory part)
[ ] handlers/items_list.md:item_list_nonstackable:14 - handlers/cart.md (cart part)
[ ] handlers/items_list.md:item_list_nonstackable:14 - handlers/storage.md (storage part)

[ ] handlers/items_list.md:item_list_end:4 - handlers/inventory.md (inventory part)
[ ] handlers/items_list.md:item_list_end:4 - handlers/cart.md (cart part)
[ ] handlers/items_list.md:item_list_end:4 - handlers/storage.md (storage part)

[ ] handlers/movement.md:actor_action:80 - handlers/actors.md
[ ] handlers/movement.md:actor_display_compatibility:100 - handlers/actors.md

[ ] handlers/script.md:show_script:2 - handlers/npc.md

[ ] handlers/social.md:actor_info:278 - handlers/actors.md

[ ] handlers/system.md:changeToInGameState:2 - common/connection.md

[ ] handlers/ui.md:map_loaded:150 - handlers/map.md
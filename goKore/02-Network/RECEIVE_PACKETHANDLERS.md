# Packet Handlers

This file details the handling logic for each packet type in Receive.pm.

## Handler Structure

Each handler is described with:
- Packet ID and Name
- Purpose
- Detailed logic steps
- Game state modifications
- Error handling
- Plugin hook points

## Handlers

### 1 Login/Authentication Handlers

#### 0x0AC4 - Login-server Connect Handler
- Purpose: Handle the connection to the login server
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0B02 - Login PIN Code Request Handler
- Purpose: Handle the request for login PIN code
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 2 Character/Player Handlers

#### 0x0119 - Character Status Handler
- Purpose: Handle updates to character status
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x02B9 - Character Status Handler (Renewal)
- Purpose: Handle updates to character status in renewal servers
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0ACB - Stat Info Handler (Extended)
- Purpose: Handle extended stat information updates
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 3 Inventory/Item Handlers

#### 0x0122 - Inventory Items List Handler
- Purpose: Handle the list of items in the inventory
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0295 - Item Use Response Handler
- Purpose: Handle the response after using an item
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 4 Skill Handlers

#### 0x07E1 - Skill List Handler
- Purpose: Handle the list of skills
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0ADE - Extended Skills Handler
- Purpose: Handle extended skill information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 5 Chat Handlers

#### 0x0383 - Guild Chat Handler
- Purpose: Handle guild chat messages
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 6 NPC Interaction Handlers

#### 0x0A37 - NPC Dialog Handler (Extended)
- Purpose: Handle extended NPC dialog interactions
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 7 Party Handlers

#### 0x0194 - Party Invite Handler
- Purpose: Handle party invitations
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0195 - Party Invite Result Handler
- Purpose: Handle the result of party invitations
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 8 Guild Handlers

#### 0x01EC - Guild Member Position Info Handler
- Purpose: Handle information about guild member positions
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 9 Storage Handlers

#### 0x0995 - Storage Password Request Handler
- Purpose: Handle requests for storage password
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 11 Homunculus Handlers

#### 0x0B2F - Homunculus Info Handler (Extended)
- Purpose: Handle extended homunculus information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 12 Pet Handlers

#### 0x01AE - Pet Info Handler
- Purpose: Handle pet information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 14 Quest Handlers

#### 0x09F9 - Quest List Handler (Extended)
- Purpose: Handle extended quest list information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 15 Rodex Mail System Handlers

#### 0x09F1 - Check Player Handler
- Purpose: Handle player checks for the Rodex mail system
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 16 Captcha/Macro Detection System Handlers

#### 0x0A53 - Captcha Upload Request Handler
- Purpose: Handle requests to upload captcha images
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A55 - Captcha Upload Status Handler
- Purpose: Handle the status of captcha image uploads
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A57 - Macro Reporter Status Handler
- Purpose: Handle the status of macro reporters
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A58 - Macro Detector Request Handler
- Purpose: Handle requests for macro detection
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A59 - Macro Detector Download Handler
- Purpose: Handle downloads related to macro detection
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A5B - Macro Detector Show Handler
- Purpose: Handle the display of macro detection information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0A5D - Macro Detector Status Handler
- Purpose: Handle the status of macro detection
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 17 Attendance System Handlers

#### 0x0AE2 - Open UI Handler
- Purpose: Handle opening of UI elements
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0AF0 - UI Action Response Handler
- Purpose: Handle responses to UI actions
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 18 Bank System Handlers

#### 0x09AC - Bank Account Status Handler
- Purpose: Handle bank account status information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 19 Roulette System Handlers

#### 0x0A24 - Coin Conversion Handler
- Purpose: Handle coin conversions in the roulette system
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 21 Clan System Handlers

#### 0x098A - Clan User Count Handler
- Purpose: Handle clan user count information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0991 - Clan Skill Handler
- Purpose: Handle clan skill information
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 26 Other Miscellaneous Handlers

#### 0x01C8 - Use Card Handler
- Purpose: Handle the use of cards
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0221 - Upgrade Skill Level Handler
- Purpose: Handle skill level upgrades
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x07D8 - Bargain Sale Handler
- Purpose: Handle bargain sale events
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x07D9 - Bargain Sale Open Handler
- Purpose: Handle opening of bargain sales
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x07DA - Bargain Sale Add Item Handler
- Purpose: Handle adding items to bargain sales
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0803 - Booking Register Handler
- Purpose: Handle registration for bookings
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0805 - Booking Search Handler
- Purpose: Handle searches for bookings
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0807 - Booking Delete Handler
- Purpose: Handle deletion of bookings
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x0809 - Booking Insert Handler
- Purpose: Handle insertion of new bookings
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x080A - Booking Update Handler
- Purpose: Handle updates to existing bookings
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x080B - Booking Delete Handler
- Purpose: Handle deletion of bookings (alternative)
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x09A0 - Equip Switch Log Handler
- Purpose: Handle logs for equipment switches
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

#### 0x09A2 - Equip Switch Handler
- Purpose: Handle equipment switches
- Detailed logic steps: [To be implemented]
- Game state modifications: [To be implemented]
- Error handling: [To be implemented]
- Plugin hook points: [To be implemented]

### 0x0097 - Private Message Handler
- Purpose: Process incoming private messages
- Logic:
  1. Extract message length, sender name, flag, and message content
  2. Decode message content from server's character set
  3. Check if sender is in the friends list
  4. Update last message received time
  5. Display message in client UI
- Game State Modifications:
  - Update `$lastMsgFromUser`
  - Add sender to `@privMsgUsers` if not present
- Error Handling:
  - Check for empty sender name or message
  - Validate message length
- Plugin Hook: `packet_privMsg`

### 0x00A0 - Inventory Item Added Handler
- Purpose: Handle addition of items to inventory
- Logic:
  1. Parse item details (ID, amount, properties)
  2. Check if item already exists in inventory
  3. If exists, update quantity
  4. If new, create new inventory entry
  5. Apply any item modifiers or cards
  6. Check for auto-drop conditions
- Game State Modifications:
  - Update `$char->{inventory}`
  - Modify `$char->{weight}` and `$char->{weight_max}`
- Error Handling:
  - Validate item ID and amount
  - Check for inventory capacity
- Plugin Hook: `item_added`

### 0x0078 - Character Move Handler
- Purpose: Update character position
- Logic:
  1. Extract new coordinates and direction
  2. Calculate movement vector
  3. Update character position
  4. Check for collisions or warps
- Game State Modifications:
  - Update `$char->{pos}` and `$char->{pos_to}`
  - Update `$char->{look}{body}` for direction
- Error Handling:
  - Validate coordinates within map bounds
- Plugin Hook: `character_move`

### 0x00B0 - Stats Info Handler
- Purpose: Update character statistics
- Logic:
  1. Identify stat type
  2. Update corresponding stat value
  3. Recalculate derived stats if necessary
- Game State Modifications:
  - Update relevant stat in `$char->{stats}`
- Error Handling:
  - Validate stat type and value range
- Plugin Hook: `stat_info`

### 0x00B4 - NPC Talk Handler
- Purpose: Process NPC dialog
- Logic:
  1. Extract NPC ID and message
  2. Parse message for dialog options
  3. Update NPC conversation state
  4. Display message and options to user
- Game State Modifications:
  - Update `$talk{ID}` and `$talk{msg}`
- Error Handling:
  - Check for valid NPC ID
- Plugin Hook: `npc_talk`

### 0x00BE - Update Status Handler
- Purpose: Handle status effect changes
- Logic:
  1. Identify status effect type
  2. Apply or remove status effect
  3. Update character appearance if necessary
- Game State Modifications:
  - Update `$char->{statuses}`
  - Modify `$char->{look}` if applicable
- Error Handling:
  - Validate status effect type
- Plugin Hook: `status_change`

### 0x0080 - Remove Entity Handler
- Purpose: Remove an entity from the game world
- Logic:
  1. Identify entity type (player, monster, NPC, etc.)
  2. Remove entity from appropriate list
  3. Update relevant game state
- Game State Modifications:
  - Remove from `$playersList`, `$monstersList`, or `$npcsList`
- Error Handling:
  - Check for entity existence before removal
- Plugin Hook: `entity_removed`

## Server-Specific Logic

### 0x0078 - Character Move (bRO variation)
- Additional step: Apply bRO-specific coordinate transformation
- Logic:
  1. Decrypt coordinates using bRO algorithm
  2. Apply standard move logic

### 0x00A0 - Inventory Item Added (Newer servers)
- Additional step: Process extended card information
- Logic:
  1. Extract additional card slots
  2. Apply card effects to item properties

Note: This list covers the main packet handlers. Receive.pm contains handlers for many more packet types, each with its own specific logic and state modifications.

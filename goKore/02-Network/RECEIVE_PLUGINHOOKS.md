# Plugin Hooks

This file lists all plugin hook points in Receive.pm, their timing in the packet handling process, and the data they provide.

## Hook Structure

Each hook is described with:
- Hook name
- Timing in packet handling process
- Data passed to the hook
- Purpose

## Hooks List

### packet_pre/packet_post
- Timing: Before and after processing any packet
- Data:
  - `args->{switch}`: Packet ID
  - `args->{RAW_MSG}`: Raw packet data
- Purpose: Allow plugins to intercept or modify any packet before/after processing

### packet_privMsg
- Timing: After processing a private message packet
- Data:
  - `args->{privMsgUser}`: Sender's name
  - `args->{privMsg}`: Message content
  - `args->{isAdmin}`: Boolean indicating if sender is admin
- Purpose: Allow plugins to intercept or modify private messages

### packet_pubMsg
- Timing: After processing a public chat message
- Data:
  - `args->{pubMsgUser}`: Speaker's name
  - `args->{pubMsg}`: Message content
  - `args->{pubID}`: Speaker's ID
- Purpose: Allow plugins to intercept or modify public chat messages

### packet_sentPM
- Timing: After sending a private message
- Data:
  - `args->{to}`: Recipient's name
  - `args->{msg}`: Message content
- Purpose: Allow plugins to log or modify sent private messages

### item_added
- Timing: After adding an item to the inventory
- Data:
  - `args->{item}`: Item object with all properties
  - `args->{amount}`: Quantity added
- Purpose: Allow plugins to react to inventory changes

### packet_charNameUpdate
- Timing: After receiving updated character name information
- Data:
  - `args->{player}`: Player object
- Purpose: Allow plugins to track or modify character name changes

### packet_mapChange
- Timing: Before changing the current map
- Data:
  - `args->{map}`: New map name
- Purpose: Allow plugins to perform actions before map change

### packet_sendMapLoaded
- Timing: After sending the map loaded packet to the server
- Data:
  - `args->{map}`: Current map name
- Purpose: Allow plugins to perform actions after confirming map load

### packet_skill_use
- Timing: After processing a skill use packet
- Data:
  - `args->{sourceID}`: ID of skill user
  - `args->{targetID}`: ID of skill target
  - `args->{skillID}`: ID of used skill
  - `args->{amount}`: Skill level or amount
  - `args->{x}`: X coordinate (for ground-targeted skills)
  - `args->{y}`: Y coordinate (for ground-targeted skills)
- Purpose: Allow plugins to react to or modify skill usage

### packet_attack
- Timing: After processing an attack packet
- Data:
  - `args->{sourceID}`: Attacker's ID
  - `args->{targetID}`: Target's ID
  - `args->{damage}`: Damage dealt
  - `args->{type}`: Attack type
- Purpose: Allow plugins to react to or modify attacks

### packet_npc_talk
- Timing: After receiving NPC dialog
- Data:
  - `args->{ID}`: NPC's ID
  - `args->{msg}`: Dialog message
- Purpose: Allow plugins to intercept or modify NPC conversations

### packet_storage_open
- Timing: After opening storage
- Data:
  - `args->{items}`: List of storage items
- Purpose: Allow plugins to react to storage access

### packet_buy_result
- Timing: After attempting to buy an item
- Data:
  - `args->{result}`: Success or failure
  - `args->{itemID}`: ID of item attempted to buy
- Purpose: Allow plugins to react to shop transactions

### packet_sell_result
- Timing: After attempting to sell an item
- Data:
  - `args->{result}`: Success or failure
  - `args->{itemID}`: ID of item attempted to sell
- Purpose: Allow plugins to react to shop transactions

### packet_vender_store
- Timing: After receiving vender shop information
- Data:
  - `args->{venderID}`: Vender's ID
  - `args->{items}`: List of items in the shop
- Purpose: Allow plugins to process or modify vender information

### packet_actor_status_active
- Timing: After receiving status effect information
- Data:
  - `args->{ID}`: Actor's ID
  - `args->{type}`: Status effect type
  - `args->{flag}`: Status effect flag
  - `args->{tick}`: Duration of the effect
- Purpose: Allow plugins to react to status effect changes

### packet_actor_movement
- Timing: After processing an actor's movement
- Data:
  - `args->{ID}`: Actor's ID
  - `args->{coords}`: New coordinates
- Purpose: Allow plugins to track or modify actor movements

### packet_actor_exists
- Timing: After an actor appears in view
- Data:
  - `args->{actor}`: Actor object
- Purpose: Allow plugins to react to new actors in the vicinity

### packet_actor_died
- Timing: After an actor dies
- Data:
  - `args->{actor}`: Actor object
- Purpose: Allow plugins to react to actor deaths

### packet_actor_info
- Timing: After receiving detailed actor information
- Data:
  - `args->{ID}`: Actor's ID
  - `args->{name}`: Actor's name
  - `args->{party_name}`: Party name (if applicable)
  - `args->{guild_name}`: Guild name (if applicable)
  - `args->{guild_title}`: Guild title (if applicable)
- Purpose: Allow plugins to process or modify actor information

### packet_inventory_item_removed
- Timing: After an item is removed from inventory
- Data:
  - `args->{index}`: Item's inventory index
  - `args->{amount}`: Amount removed
- Purpose: Allow plugins to react to inventory changes

### packet_map_property
- Timing: After receiving map property information
- Data:
  - `args->{type}`: Map property type
- Purpose: Allow plugins to react to changes in map properties

### packet_quest_update
- Timing: After receiving quest update information
- Data:
  - `args->{questID}`: Quest ID
  - `args->{active}`: Quest active status
  - `args->{time_start}`: Quest start time
  - `args->{time_expire}`: Quest expiration time
- Purpose: Allow plugins to track or modify quest information

### packet_homunculus_info
- Timing: After receiving homunculus information
- Data:
  - `args->{homunculus}`: Homunculus object with all properties
- Purpose: Allow plugins to process or modify homunculus information

### packet_storage_closed
- Timing: After closing storage
- Data: None
- Purpose: Allow plugins to react to storage closure

### packet_rodex_mail_list
- Timing: After receiving Rodex mail list
- Data:
  - `args->{mails}`: List of mail objects
- Purpose: Allow plugins to process or modify Rodex mail information

### packet_achievement_update
- Timing: After receiving achievement update
- Data:
  - `args->{achievementID}`: Achievement ID
  - `args->{completed}`: Completion status
- Purpose: Allow plugins to track or modify achievement progress

### packet_guild_chat
- Timing: After receiving a guild chat message
- Data:
  - `args->{chatMsgUser}`: Sender's name
  - `args->{chatMsg}`: Message content
- Purpose: Allow plugins to process or modify guild chat messages

### packet_market_info
- Timing: After receiving market information
- Data:
  - `args->{items}`: List of market items
- Purpose: Allow plugins to process or modify market information

### packet_skill_cast
- Timing: After receiving skill cast information
- Data:
  - `args->{sourceID}`: Caster's ID
  - `args->{targetID}`: Target's ID
  - `args->{x}`: X coordinate
  - `args->{y}`: Y coordinate
  - `args->{skillID}`: Skill ID
  - `args->{amount}`: Skill level or amount
  - `args->{startTime}`: Cast start time
  - `args->{castTime}`: Total cast time
- Purpose: Allow plugins to react to or modify skill casting

## Server-Specific Hooks

### bRO_packet_mapChange
- Server: bRO (Brazilian Ragnarok Online)
- Timing: Before standard packet_mapChange
- Data: Same as packet_mapChange
- Purpose: Allow bRO-specific actions before map change

### idRO_packet_item_list
- Server: idRO (Indonesian Ragnarok Online)
- Timing: After receiving item list
- Data:
  - `args->{type}`: List type (inventory, storage, cart)
  - `args->{items}`: List of items
- Purpose: Allow idRO-specific item list processing

Note: This expanded list covers most of the hooks used in Receive.pm. There are still more specialized hooks for various game events and server-specific implementations.

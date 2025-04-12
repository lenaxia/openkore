# Game State Variables

This file lists all game state variables that are modified by packet handlers in Receive.pm.

## Character Information

- `$char->{name}`: Character name
- `$char->{ID}`: Character ID
- `$char->{accountID}`: Account ID
- `$char->{lv}`: Base level
- `$char->{lv_job}`: Job level
- `$char->{exp}`: Base experience
- `$char->{exp_job}`: Job experience
- `$char->{jobID}`: Job ID
- `$char->{sex}`: Character gender

## Stats

- `$char->{str}`: Strength
- `$char->{agi}`: Agility
- `$char->{vit}`: Vitality
- `$char->{int}`: Intelligence
- `$char->{dex}`: Dexterity
- `$char->{luk}`: Luck
- `$char->{points_free}`: Stat points available
- `$char->{points_skill}`: Skill points available

## Health and Status

- `$char->{hp}`: Current HP
- `$char->{hp_max}`: Maximum HP
- `$char->{sp}`: Current SP
- `$char->{sp_max}`: Maximum SP
- `$char->{weight}`: Current weight
- `$char->{weight_max}`: Maximum weight capacity
- `$char->{walk_speed}`: Walking speed
- `$char->{attack_speed}`: Attack speed
- `$char->{statuses}`: Hash of active status effects

## Inventory and Equipment

- `$char->{inventory}`: Hash of inventory items
- `$char->{equipment}`: Hash of equipped items
- `$char->{slots}{inventory}`: Number of inventory slots
- `$char->{slots}{storage}`: Number of storage slots

## Skills

- `$char->{skills}`: Hash of character skills
- `$char->{permitSkill}`: Currently permitted skill
- `$char->{encoreSkill}`: Skill available for encore

## Map and Position

- `$char->{pos}`: Current position (x, y coordinates)
- `$char->{pos_to}`: Destination position for movement
- `$char->{look}{body}`: Body direction
- `$char->{look}{head}`: Head direction
- `$field->{name}`: Current map name
- `$field->{baseName}`: Base name of the current map
- `$field->{width}`: Map width
- `$field->{height}`: Map height

## Other Actors

- `$playersList`: List of visible players
- `$monstersList`: List of visible monsters
- `$npcsList`: List of visible NPCs
- `$petsList`: List of visible pets
- `$slavesList`: List of slaves (homunculus, mercenary)

## Game Environment

- `$itemsList`: List of items on the ground
- `$portalsList`: List of visible portals
- `$spellsID`: List of active spells/effects

## Chat and Social

- `@chatRooms`: List of available chat rooms
- `%guild`: Guild information
- `@privMsgUsers`: List of users who have sent private messages
- `$currentChatRoom`: Current chat room information
- `%friends`: List of friends
- `%incomingFriend`: Incoming friend requests

## Economy

- `$char->{zeny}`: Current zeny (money) amount
- `@venderListsID`: List of vender IDs
- `%venderLists`: Hash of vender information
- `$buyerList`: List of buyer shops

## Quest and Mission

- `%questList`: Hash of active quests
- `$questID`: Current quest ID

## Connection State

- `$conState`: Current connection state
- `$accountID`: Account ID
- `$charID`: Character ID
- `$sessionID`: Session ID
- `$sessionID2`: Secondary session ID
- `$accountSex`: Account gender

## Misc

- `$char->{party}`: Party information
- `$char->{muted}`: Mute status
- `$char->{exp_last}`: Last recorded experience
- `$char->{time_move}`: Last movement time
- `$char->{time_move_calc}`: Calculated movement time
- `$char->{last_skill_used}`: Last used skill
- `$char->{homunculus}`: Homunculus information
- `$char->{mercenary}`: Mercenary information
- `$char->{cart}`: Cart information

Note: This list covers the main game state variables. Receive.pm interacts with and modifies many more variables throughout the OpenKore system.

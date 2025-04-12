# Packet Definitions

This file contains a comprehensive list of all packet types and their structures handled in Receive.pm.

## Login/Authentication Packets

### Login-server Connect (0x0AC4)
- Purpose: Handles initial connection to login server
- Structure: `v V2 C2 Z24 V2 v`
- Fields:
  - `version`: Client version
  - `ID`: Account ID
  - `sessionID`: Session ID
  - `sex`: Account gender
  - `serverName`: Server name
  - `serverUsers`: Number of users
  - `serverIP`: Server IP
  - `serverPort`: Server port

### Login PIN Code Request (0x0B02)
- Purpose: Requests PIN code for account login
- Structure: `v V C`
- Fields:
  - `accountID`: Account ID
  - `flag`: Request type
  - `seed`: PIN encryption seed

## Character/Player Packets

### Character Status (0x0119)
- Purpose: Updates character status information
- Structure: `v V4 v V5 v5`
- Fields:
  - `type`: Status type
  - `value`: Status value
  - `bonus`: Status bonus
  - `upgrade`: Status upgrade points
  - `val`: Additional value

### Character Status Renewal (0x02B9)
- Purpose: Updates character status for renewal servers
- Structure: `v V4 v V5 v5 V2`
- Fields:
  - (Same as 0x0119 plus:)
  - `type2`: Additional status type
  - `value2`: Additional status value

### Extended Stat Info (0x0ACB)
- Purpose: Provides extended character statistics
- Structure: `v V2 v V5 v10`
- Fields:
  - `type`: Stat type
  - `base`: Base value
  - `bonus`: Bonus value
  - `equip`: Equipment bonus
  - `enchant`: Enchant bonus

## Packet Structure Format

Each packet is defined with:
- Packet ID (hex)
- Name and description
- Structure (using format specifiers)
- Field names and descriptions
- Server variations if applicable

Format specifiers:
- v: 16-bit unsigned integer
- V: 32-bit unsigned integer
- C: 8-bit unsigned integer
- c: 8-bit signed integer
- x: padding byte
- a: fixed-length string
- Z: null-terminated string
- b: bit string
- f: 32-bit float
- d: 64-bit float

## Core Game Packets

### Skill-related Packets

#### 0x0114 - Skill Use Failed
- Structure: `v V v2 C`
- Fields:
  - `skillID`: Skill ID
  - `btype`: Basis type
  - `unknown`: Unknown (always 0)
  - `cause`: Failure cause
  - `itemID`: Item ID (if applicable)

#### 0x01B9 - Cast Cancelled
- Structure: `V`
- Fields:
  - `ID`: Actor ID

#### 0x07E6 - Skill Cast (with coordinates)
- Structure: `v V4 v3 C`
- Fields:
  - `skillID`: Skill ID
  - `sourceID`: Caster ID
  - `targetID`: Target ID
  - `startTime`: Cast start time
  - `endTime`: Cast end time
  - `xPos`: X coordinate
  - `yPos`: Y coordinate
  - `range`: Skill range
  - `fail`: Failure flag

### Party-related Packets

#### 0x0101 - Party Share EXP Setting
- Structure: `V2`
- Fields:
  - `type`: Share type (0: Individual, 1: Even Share)
  - `flag`: Additional flags

#### 0x0105 - Party Settings
- Structure: `V3`
- Fields:
  - `exp`: Experience sharing (0: Individual, 1: Even Share)
  - `item`: Item sharing (0: Individual, 1: Even Share)
  - `itemPickup`: Item pickup (0: Individual, 1: Even Share)

#### 0x0109 - Party Chat
- Structure: `v2 Z*`
- Fields:
  - `len`: Packet length
  - `ID`: Speaker ID
  - `message`: Chat message

### Storage-related Packets

#### 0x00F2 - Storage Closed
- Structure: `v`
- Fields:
  - `flag`: Closure flag

#### 0x00F4 - Storage Item Added
- Structure: `v3 C3 a8 V C2 a4 v`
- Fields:
  - `index`: Item index
  - `amount`: Item amount
  - `nameID`: Item ID
  - `identified`: Identification flag
  - `broken`: Broken status
  - `upgrade`: Upgrade level
  - `cards`: Card slots
  - `type`: Item type
  - `result`: Operation result

#### 0x00F6 - Storage Item Removed
- Structure: `v2`
- Fields:
  - `index`: Item index
  - `amount`: Amount removed

### Expanded Inventory System Packets

#### 0x0B0B - Inventory Expansion Cost Info
- Structure: `v V2`
- Fields:
  - `result`: Result code
  - `itemID`: Required item ID
  - `zeny`: Required zeny amount

#### 0x0B0D - Inventory Expansion Result
- Structure: `v`
- Fields:
  - `result`: Expansion result (0: Success, others: Failure reasons)

### Detailed Actor Info Packets

#### 0x0229/0x022A - Character Status/Equipment
- Structure: `v V5 C v11 V3 v17 Z24 C8 v4 C`
- Fields:
  - `ID`: Character ID
  - `GID`: Character GID
  - `stance`: Character stance
  - `sex`: Gender
  - `hairColor`: Hair color
  - `headBottom`: Bottom headgear
  - `shield`: Shield
  - `headTop`: Top headgear
  - `headMid`: Middle headgear
  - `hairStyle`: Hair style
  - Various equipment and appearance fields
  - `name`: Character name
  - Various stat fields

#### 0x02DB/0x02DC - Character Look Change
- Structure: `V C2 v`
- Fields:
  - `GID`: Character GID
  - `type`: Change type
  - `value`: New value
  - `index`: Item index (if applicable)

#### 0x0B19 - Character Look Change with Rotation
- Structure: `V C2 v2`
- Fields:
  - `GID`: Character GID
  - `type`: Change type
  - `value`: New value
  - `index`: Item index (if applicable)
  - `rotation`: Body rotation

### Detailed Trade Packets

#### 0x00E6 - Trade Request
- Structure: `V Z24`
- Fields:
  - `ID`: Requester ID
  - `name`: Requester name

#### 0x00E7 - Trade Response
- Structure: `C`
- Fields:
  - `type`: Response type (3: Accept, 4: Cancel)

#### 0x00E8 - Trade Item Added
- Structure: `v2 C3 a8 V`
- Fields:
  - `index`: Item index
  - `amount`: Item amount
  - `type`: Item type
  - `identified`: Identification flag
  - `broken`: Broken status
  - `cards`: Card slots
  - `price`: Item price (if in vending)

#### 0x00EA - Trade Complete
- Structure: `C`
- Fields:
  - `result`: Trade result (0: Success, 1: Failure)

### Detailed Pet System Packets

#### 0x01A2 - Pet Catch Process
- Structure: `C`
- Fields:
  - `result`: Catch result (0: Failure, 1: Success)

#### 0x01A3 - Pet Catch Result
- Structure: `C`
- Fields:
  - `result`: Final result (0: Failure, 1: Success)

#### 0x01A4 - Pet Feed
- Structure: `V C`
- Fields:
  - `ID`: Pet ID
  - `success`: Feeding success (0: Failure, 1: Success)

### Detailed Homunculus Packets

#### 0x022E - Homunculus Info
- Structure: `V Z24 V5 v11 V2 v12`
- Fields:
  - `ID`: Homunculus ID
  - `name`: Homunculus name
  - Various stat and info fields

#### 0x0230 - Homunculus Food
- Structure: `V C v`
- Fields:
  - `ID`: Homunculus ID
  - `success`: Feeding success
  - `foodID`: Food item ID

#### 0x0231 - Homunculus Name
- Structure: `V Z24`
- Fields:
  - `ID`: Homunculus ID
  - `name`: New name

### Detailed Mercenary Packets

#### 0x029B - Mercenary Info
- Structure: `V Z24 V5 v11 V2 v12`
- Fields:
  - `ID`: Mercenary ID
  - `name`: Mercenary name
  - Various stat and info fields

#### 0x02A2 - Mercenary Skill Update
- Structure: `v V v3`
- Fields:
  - `skillID`: Skill ID
  - `type`: Update type
  - `level`: New skill level
  - `sp`: SP cost
  - `range`: Skill range

### Detailed Quest Packets

#### 0x02B1 - Quest List
- Structure: `v3 a*`
- Fields:
  - `len`: Packet length
  - `questID`: Quest ID
  - `active`: Quest active state
  - `questInfo`: Quest details

#### 0x02B2 - Quest Add
- Structure: `V C V2 v a*`
- Fields:
  - `questID`: Quest ID
  - `active`: Quest active state
  - `startTime`: Start time
  - `endTime`: End time
  - `mobs`: Number of mob objectives
  - `mobInfo`: Mob objective details

#### 0x02B3 - Quest Update
- Structure: `V C V2 v a*`
- Fields:
  - `questID`: Quest ID
  - `active`: Quest active state
  - `startTime`: Start time
  - `endTime`: End time
  - `mobs`: Number of mob objectives
  - `mobInfo`: Updated mob objective details

### Detailed Mail Packets

#### 0x0242 - Mail Read
- Structure: `v V2 v Z40 Z24 V v2 V2 v a*`
- Fields:
  - `len`: Packet length
  - `mailID`: Mail ID
  - `title`: Mail title
  - `body`: Mail body
  - `sender`: Sender name
  - `zeny`: Attached zeny
  - `itemCount`: Number of attached items
  - Various item details

#### 0x0245 - Mail Delete
- Structure: `V C`
- Fields:
  - `mailID`: Mail ID
  - `result`: Deletion result (0: Success, 1: Failure)

#### 0x0246 - Mail Return
- Structure: `V C`
- Fields:
  - `mailID`: Mail ID
  - `result`: Return result (0: Success, 1: Failure)

### Detailed Achievement Packets

#### 0x0839 - Achievement Update
- Structure: `V V10 V C V`
- Fields:
  - `achievementID`: Achievement ID
  - `objective`: Objective progress (10 fields)
  - `completed`: Completion time
  - `reward`: Reward claim status
  - `date`: Date of completion

#### 0x083A - Achievement Reward
- Structure: `V C`
- Fields:
  - `achievementID`: Achievement ID
  - `result`: Reward claim result

#### 0x083B - Achievement List
- Structure: `v V a*`
- Fields:
  - `len`: Packet length
  - `achievementCount`: Number of achievements
  - `achievementList`: List of achievement data

### Detailed Clan Packets

#### 0x098F - Clan Leave
- Structure: `V`
- Fields:
  - `result`: Leave result (0: Success, 1: Failure)

#### 0x0990 - Clan Message
- Structure: `v2 Z24 Z*`
- Fields:
  - `len`: Packet length
  - `messageLen`: Message length
  - `sender`: Sender name
  - `message`: Clan message

### Detailed Macro Detection Packets

#### 0x0A4E - Macro Detector Result
- Structure: `C V2`
- Fields:
  - `result`: Detection result
  - `captchaKey`: Captcha key
  - `imageSize`: Image size

#### 0x0A4F - Macro Reporter Result
- Structure: `C`
- Fields:
  - `result`: Reporting result

### Login/Authentication

#### 0x0069 - Account Info
- Structure: `v V a4 a4 a4 a4 C`
- Fields:
  - `len`: Packet length
  - `accountID`: Account ID
  - `sessionID`: Session ID
  - `sessionID2`: Secondary session ID
  - `lastLoginIP`: Last login IP
  - `lastLoginTime`: Last login timestamp
  - `sex`: Account gender

#### 0x006A - Login Error
- Structure: `v C`
- Fields:
  - `type`: Error type
  - `error`: Error code

### Character/Player

#### 0x0078 - Character Move
- Structure: `v3 x9`
- Fields:
  - `coordX`: Target X coordinate
  - `coordY`: Target Y coordinate
  - `direction`: Movement direction
- Variations:
  - bRO: `v3 x5`
  - idRO: `v3 x7`
  - Zero: `v3 x11`

#### 0x007C - Spawn Entity
- Structure: `v V C x3 v3 x2`
- Fields:
  - `ID`: Entity ID
  - `type`: Entity type
  - `state`: Entity state
  - `x,y,dir`: Position and direction
- Types:
  - 0: PC
  - 1: NPC
  - 2: Item
  - 3: Monster
  - 4: Pet
  - 5: Homunculus
  - 6: Mercenary

#### 0x0080 - Remove Entity
- Structure: `v C`
- Fields:
  - `ID`: Entity ID
  - `type`: Removal type
- Types:
  - 0: Out of sight
  - 1: Died
  - 2: Logged out
  - 3: Teleported
  - 4: Trickdead

### Inventory/Items

#### 0x00A0 - Inventory Item Added
- Structure: `v3 C3 a8 V C2 a4 v`
- Fields:
  - `index`: Item index
  - `amount`: Item amount
  - `nameID`: Item ID
  - `identified`: Identification flag
  - `broken`: Broken status
  - `upgrade`: Upgrade level
  - `cards`: Card slots
  - `type_equip`: Equipment type
  - `type`: Item type
  - `fail`: Failure flag
  - `expire`: Expiration time
- Variations:
  - Renewal: `v3 C3 a8 V C2 a4 v4 V`
  - Zero: `v3 C3 a8 V C2 a4 v4 V2`

#### 0x00AF - Inventory Item Removed
- Structure: `v2`
- Fields:
  - `index`: Item index
  - `amount`: Amount removed

### NPC Interaction

#### 0x00B4 - NPC Talk
- Structure: `v2 Z*`
- Fields:
  - `ID`: NPC ID
  - `message`: Dialog text

#### 0x00B5 - NPC Talk Continue
- Structure: `v`
- Fields:
  - `ID`: NPC ID

#### 0x00B6 - NPC Talk Close
- Structure: `v`
- Fields:
  - `ID`: NPC ID

#### 0x00B7 - NPC Talk Responses
- Structure: `v2 Z*`
- Fields:
  - `ID`: NPC ID
  - `message`: Response options

### Stats/Status

#### 0x00B0 - Stats Info
- Structure: `v V`
- Fields:
  - `type`: Stat type
  - `val`: New value
- Types:
  - 0: Speed
  - 1: BaseExp
  - 2: JobExp
  - 3: Karma
  - 4: Manner
  - 5: HP
  - 6: MaxHP
  - 7: SP
  - 8: MaxSP

#### 0x00BD - Full Stats Info
- Structure: `v16`
- Fields:
  - `points`: Status points
  - Basic stats: `str,agi,vit,int,dex,luk`
  - Bonus stats: `str_bonus,agi_bonus,vit_bonus,int_bonus,dex_bonus,luk_bonus`

#### 0x00BE - Status Change
- Structure: `v C`
- Fields:
  - `statusID`: Status effect ID
  - `flag`: 1=on, 0=off

### Combat

#### 0x008A - Attack Info
- Structure: `v V2 v3 V2`
- Fields:
  - `attackMT`: Attack motion
  - `sourceID`: Attacker ID
  - `targetID`: Target ID
  - `startTime`: Start time
  - `sourceSpeed`: Attacker speed
  - `targetSpeed`: Target speed
  - `damage`: Damage amount
  - `level`: Skill level

#### 0x0110 - Skill Use
- Structure: `v V4 v2 C2 v`
- Fields:
  - `skillID`: Skill ID
  - `sourceID`: Caster ID
  - `targetID`: Target ID
  - `startTime`: Cast start time
  - `endTime`: Cast end time
  - `x,y`: Target coordinates
  - `amount`: Skill level/damage
  - `fail`: Failure flag
  - `type`: Skill type

### Chat/Communication

#### 0x008E - Public Chat
- Structure: `v2 Z*`
- Fields:
  - `len`: Message length
  - `ID`: Speaker ID
  - `message`: Chat text

#### 0x0097 - Private Message
- Structure: `v Z24 V Z*`
- Fields:
  - `len`: Message length
  - `privMsgUser`: Sender name
  - `flag`: Message flags
  - `privMsg`: Message text

#### 0x009A - System Message
- Structure: `v2 Z*`
- Fields:
  - `len`: Message length
  - `type`: Message type
  - `message`: System message

### Party/Guild

#### 0x00FB - Party Info
- Structure: `v V Z24 Z24 V2`
- Fields:
  - `len`: Packet length
  - `partyID`: Party ID
  - `partyName`: Party name
  - `playerName`: Player name
  - `role`: Party role
  - `state`: Party state

#### 0x0162 - Guild Info
- Structure: `v V Z24 Z24 V4`
- Fields:
  - `len`: Packet length
  - `guildID`: Guild ID
  - `guildName`: Guild name
  - `masterName`: Guild master name
  - `castleID`: Castle ID
  - `members`: Member count
  - `maxMembers`: Maximum members
  - `average`: Average level

### Map/Environment

#### 0x0091 - Map Change
- Structure: `Z16 v2`
- Fields:
  - `map`: Map name
  - `x,y`: Spawn coordinates

#### 0x0092 - Map Changed
- Structure: `Z16 v2 V v2`
- Fields:
  - `map`: Map name
  - `x,y`: Character position
  - `IP`: Map server IP
  - `port`: Map server port

#### 0x009D - Weather Effects
- Structure: `v V`
- Fields:
  - `type`: Weather type
  - `flag`: Effect flags

### Trade/Vending

#### 0x00E9 - Deal Request
- Structure: `V Z24`
- Fields:
  - `ID`: Requester ID
  - `name`: Requester name

#### 0x00EF - Deal Added
- Structure: `v2 C3 a8 V`
- Fields:
  - `index`: Item index
  - `amount`: Item amount
  - `type`: Item type
  - `identified`: ID flag
  - `broken`: Break status
  - `cards`: Card data
  - `price`: Item price

## Server-Specific Packets

### bRO (Brazilian)
- 0x0A36: Monster HP Info (Tiny)
  - Structure: `v V C`
  - Fields: `ID`, `HP`, `HP_PERC`

### iRO (International)
- 0x0B1B: Attendance System
  - Structure: `v V2 v V2 v`
  - Fields: `PacketType`, `PacketLength`, `TimeStamp`, `Reward_Count`, `Reward_1`, `Reward_2`, `Reward_3`

- 0x0B20: Macro Detector
  - Structure: `v V2 v V`
  - Fields: `PacketType`, `PacketLength`, `AID`, `MacroDetectorNum`, `ImageSize`

### kRO (Korean)
- 0x0B18: Inventory Expansion
  - Structure: `v V2`
  - Fields: `PacketType`, `PacketLength`, `Result`

- 0x0B1D: Ping
  - Structure: `v V`
  - Fields: `PacketType`, `AID`

## Modern Packet Extensions

### 0x0AAA - Extended Inventory
- Structure: `v V2 C v2 V2 v a24 C v2 a4 a4 V`
- Fields:
  - `index`: Item index in inventory
  - `amount`: Item amount
  - `nameID`: Item ID
  - `identified`: Identification flag
  - `damaged`: Damage flag
  - `refine`: Refine level
  - `card1`, `card2`, `card3`, `card4`: Card slots
  - `equip`: Equipment location
  - `type`: Item type
  - `result`: Result of operation
  - `expire_time`: Expiration time
  - `bind_on_equip_type`: Bind on equip type
  - `unique_id`: Unique item ID
  - `favorite`: Favorite flag
  - `look`: View ID

### 0x0B77 - NPC Store Item List
- Structure: `v { V3 C v V }*`
- Fields:
  - `nameID`: Item ID
  - `price`: Regular price
  - `discountPrice`: Discounted price
  - `itemType`: Type of item
  - `viewSprite`: View sprite ID
  - `location`: Item location/tab

### 0x0B79 - Item List
- Structure: `v a*`
- Fields:
  - `packetLen`: Packet length
  - `itemList`: Raw item list data (format varies)

### 0x0B7B - Character Equipment Window
- Structure: `v Z24 v8 C { a57 }*`
- Fields:
  - `name`: Character name
  - `class`: Character class
  - `hairstyle`: Hair style
  - `bottom-viewid`: Bottom view ID
  - `mid-viewid`: Middle view ID
  - `up-viewid`: Upper view ID
  - `robe`: Robe view ID
  - `haircolor`: Hair color
  - `cloth-dye`: Cloth dye
  - `gender`: Character gender
  - `equip`: Equipment data (57 bytes per item)

### 0x00A0 - Inventory Item Added (Updated)
- Structure: `v3 C3 a8 V C2 a4 v V2` (for newer servers)
- Fields:
  - (All previous fields)
  - `option_count`: Number of item options
  - `option_data`: Item option data

### 0x0078 - Character Move
- Structure: `v3 x9` (Default)
- Variations:
  - bRO: `v3 x5`
  - idRO: `v3 x7`
  - Zero: `v3 x11`
- Fields:
  - `coordX`: Target X coordinate
  - `coordY`: Target Y coordinate
  - `direction`: Movement direction

### 0x00BE - Status Change
- Structure: `v C`
- Fields:
  - `statusID`: Status effect ID
  - `flag`: 1=on, 0=off
- Status Effect Types:
  - 0: Standard statuses (e.g., Poison, Silence)
  - 1: Buffs (e.g., Blessing, Increase AGI)
  - 2: Debuffs (e.g., Curse, Slow)
  - 3: Special states (e.g., Hiding, Cloaking)
  - (Add more specific status types here)

## Server-Specific Variations

### bRO (Brazilian Ragnarok Online)
- 0x0078: Uses different padding (x5 instead of x9)
- (Add other bRO-specific packet variations here)

### Zero Server
- 0x0078: Uses extended padding (x11)
- (Add other Zero server-specific packet variations here)

### Private Servers
- Note: Private servers may have custom packet structures
- Common variations:
  - Extended item properties
  - Custom status effects
  - Modified character stats
- Always verify packet structures for specific private servers

### 0x0ACB - Extended Stats
- Structure: `v V v11 V12`
- Fields:
  - `points`: Status points
  - `str`, `agi`, `vit`, `int`, `dex`, `luk`: Base stats
  - `str_bonus`, `agi_bonus`, `vit_bonus`, `int_bonus`, `dex_bonus`, `luk_bonus`: Bonus stats
  - `attack`, `attack_bonus`: Attack power
  - `attack_magic_min`, `attack_magic_max`: Magic attack range
  - `defense`, `defense_bonus`: Defense
  - `defense_magic`, `defense_magic_bonus`: Magic defense
  - `hit`, `flee`, `flee_bonus`, `critical`: Battle stats

### 0x0ADE - Extended Skills
- Structure: `v V2 v3 V v6 a24 C`
- Fields:
  - `skill_id`: Skill ID
  - `type`: Skill type
  - `level`: Skill level
  - `sp_cost`: SP cost
  - `attack_range`: Attack range
  - `skill_name`: Skill name
  - `upgradable`: Upgradable flag

### 0x0B1C - Macro Detector Image
- Structure: `v V2 v a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `AID`: Account ID
  - `MacroDetectorNum`: Macro detector number
  - `ImageData`: Compressed image data

### 0x0B2F - Homunculus Info
- Structure: `v V Z24 v14 V4 v V2`
- Fields:
  - `name`: Homunculus name
  - `level`, `hunger`, `intimacy`: Basic stats
  - `atk`, `matk`, `hit`, `crit`, `def`, `mdef`, `flee`, `aspd`: Battle stats
  - `hp`, `max_hp`, `sp`, `max_sp`: HP and SP
  - `exp`, `max_exp`: Experience
  - `skill_points`: Available skill points

Note: This list covers the core packet types and some modern extensions. The complete set of packets handled by Receive.pm includes many more specialized packets for various game features and server implementations.

## Captcha System Packets

### 0x07E6 - Captcha Session ID
- Structure: `v V`
- Fields:
  - `PacketType`: Packet identifier
  - `SessionID`: Captcha session ID

### 0x07E8 - Captcha Image
- Structure: `v V2 a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `ImageSize`: Size of the captcha image
  - `ImageData`: Captcha image data

### 0x07E9 - Captcha Answer
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Answer result (0 = correct, 1 = incorrect)

## Macro Detection Packets

### 0x0A53 - Captcha Upload Request
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Status`: Upload request status

### 0x0A55 - Captcha Upload Status
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Status`: Upload status result

### 0x0A57 - Macro Reporter Status
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Status`: Reporter status

### 0x0A58 - Macro Detector Request
- Structure: `v V V`
- Fields:
  - `PacketType`: Packet identifier
  - `ImageSize`: Size of the detector image
  - `CaptchaKey`: Unique key for the captcha
- Status Codes:
  - 0: Detection not started
  - 1: Detection in progress
  - 2: Detection failed
  - 3: Detection successful

### 0x0A59 - Macro Detector Download
- Structure: `v a*`
- Fields:
  - `PacketType`: Packet identifier
  - `ImageData`: Detector image data
- Image Format: BMP (Windows Bitmap)
- Compression: RLE (Run-Length Encoding)

### 0x0A5B - Macro Detector Show
- Structure: `v v V`
- Fields:
  - `PacketType`: Packet identifier
  - `RemainingChances`: Number of chances left
  - `RemainingTime`: Time remaining in milliseconds

### 0x0A5D - Macro Detector Result
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Detection result
- Result Types:
  - 0: Successful detection (human player)
  - 1: Detected macro use
  - 2: False positive
  - 3: Inconclusive result
- Fields:
  - `PacketType`: Packet identifier
  - `RemainingChances`: Number of chances left
  - `RemainingTime`: Time remaining in milliseconds

### 0x0A5D - Macro Detector Status
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Status`: Detector status

## Rodex Mail System

### 0x09E8 - Rodex Mail List
- Structure: `v V2 C V a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `MailBoxType`: Type of mailbox
  - `MailAmount`: Number of mails
  - `IsEnd`: Flag indicating if this is the last packet
  - `MailList`: List of mail data
- Error Codes:
  - 1: Exceeded attachment limit
  - 2: Invalid recipient
  - 3: Mailbox full
- Status Types:
  - 0: Unread
  - 1: Read
  - 2: Replied

### 0x09EB - Rodex Read Mail
- Structure: `v V2 V2 Z24 Z40 V x4 C8 v a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `MailID`: Unique mail ID
  - `Title`: Mail title
  - `SenderName`: Name of the sender
  - `Message`: Mail message
  - `ZenyAmount`: Amount of zeny attached
  - `ItemAmount`: Number of items attached
  - `ItemList`: List of attached items
- Item Attachment Structure:
  - `ItemID`: ID of the item
  - `Amount`: Quantity of the item
  - `Identified`: Identification status
  - `Broken`: Broken status
  - `Card1`, `Card2`, `Card3`, `Card4`: Card slots
  - `Options`: Item options data

### 0x09EB - Rodex Read Mail
- Structure: `v V2 V2 Z24 Z40 V x4 C8 v a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `MailID`: Unique mail ID
  - `Title`: Mail title
  - `SenderName`: Name of the sender
  - `Message`: Mail message
  - `ZenyAmount`: Amount of zeny attached
  - `ItemAmount`: Number of items attached
  - `ItemList`: List of attached items

### 0x09ED - Rodex Delete Mail
- Structure: `v V C`
- Fields:
  - `PacketType`: Packet identifier
  - `MailID`: ID of the mail to delete
  - `Result`: Deletion result

### 0x09EF - Rodex Write Mail
- Structure: `v V C`
- Fields:
  - `PacketType`: Packet identifier
  - `ReceiverID`: ID of the mail receiver
  - `Result`: Write result

### 0x09F1 - Rodex Check Player
- Structure: `v V Z24 v V`
- Fields:
  - `PacketType`: Packet identifier
  - `CharID`: Character ID
  - `Name`: Character name
  - `Level`: Character level
  - `CharacterExists`: Flag indicating if the character exists

### 0x09F3 - Rodex Send Mail Result
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Send result

## Attendance System

### 0x0AE2 - Open UI (Attendance)
- Structure: `v C V`
- Fields:
  - `PacketType`: Packet identifier
  - `UIType`: Type of UI to open (7 for Attendance)
  - `Data`: Additional data for the UI
- Attendance Reward Types:
  - 1: Item
  - 2: Zeny
  - 3: Experience Points
- Attendance Streak Bonuses:
  - 7 days: Extra item
  - 14 days: Choice box
  - 30 days: Special costume
- Error Codes:
  - 0: Success
  - 1: Already claimed today
  - 2: Attendance period ended
  - 3: Network error

### 0x0AF0 - UI Action Response
- Structure: `v V V`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Attendance claim result
  - `RewardData`: Reward details
  - `UIType`: Type of UI
  - `Data`: Response data

## Bank System

### 0x09A6 - Banking Check
- Structure: `v V v`
- Fields:
  - `PacketType`: Packet identifier
  - `Balance`: Current bank balance
  - `Reason`: Reason code
- Error Codes:
  - 0: Success
  - 1: Exceeded daily limit
  - 2: Insufficient funds
  - 3: Account restricted

### 0x09A8 - Banking Deposit Result
- Structure: `v C V V`
- Fields:
  - `PacketType`: Packet identifier
  - `Reason`: Result reason
  - `Money`: Amount deposited
  - `Balance`: New balance
- Transaction Limits:
  - Daily deposit cap: 1,000,000,000 zeny
  - Minimum deposit: 1,000 zeny

### 0x09AA - Banking Withdraw Result
- Structure: `v C V V`
- Fields:
  - `PacketType`: Packet identifier
  - `Reason`: Result reason
  - `Money`: Amount withdrawn
  - `Balance`: New balance
- Transaction Limits:
  - Daily withdraw cap: 500,000,000 zeny
  - Minimum withdraw: 1,000 zeny

### 0x09AC - Bank Account Status
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `AccountStatus`: Status of the bank account
- Account Status Types:
  - 0: Normal
  - 1: VIP
  - 2: Restricted
  - 3: Frozen

## Roulette System

### 0x0A1A - Roulette Window
- Structure: `v C V C C v V3`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Open result
  - `Serial`: Roulette serial
  - `Stage`: Current stage
  - `Price`: Price index
  - `AdditionalItem`: Additional item ID
  - `GoldPoint`: Gold coin amount
  - `SilverPoint`: Silver coin amount
  - `BronzePoint`: Bronze coin amount
- Reward Types:
  - 0: Common
  - 1: Rare
  - 2: Jackpot
- Result Codes:
  - 0: Free spin available
  - 1: Paid spin
  - 2: Out of coins

### 0x0A1C - Roulette Info
- Structure: `v V a*`
- Fields:
  - `PacketType`: Packet identifier
  - `Serial`: Roulette serial
  - `RouletteInfo`: Roulette item information

### 0x0A20 - Roulette Window Update
- Structure: `v C v2 v V3`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Update result
  - `Stage`: Current stage
  - `Price`: Price index
  - `AdditionalItem`: Additional item ID
  - `GoldPoint`: Gold coin amount
  - `SilverPoint`: Silver coin amount
  - `BronzePoint`: Bronze coin amount

### 0x0A22 - Roulette Reward
- Structure: `v C v`
- Fields:
  - `PacketType`: Packet identifier
  - `Type`: Reward type
  - `ItemID`: Rewarded item ID

### 0x0A24 - Coin Conversion
- Structure: `v C V3`
- Fields:
  - `PacketType`: Packet identifier
  - `ConversionType`: Type of conversion
  - `BronzeAmount`: Bronze coins
  - `SilverAmount`: Silver coins
  - `GoldAmount`: Gold coins
- Conversion Rates:
  - 10 Bronze = 1 Silver
  - 10 Silver = 1 Gold

## Achievement System

### 0x0A23 - Achievement List
- Structure: `v V2 a*`
- Fields:
  - `PacketType`: Packet identifier
  - `PacketLength`: Length of the packet
  - `AchievementCount`: Number of achievements
  - `AchievementList`: List of achievement data

### 0x0A24 - Achievement Update
- Structure: `v V10 V C V`
- Fields:
  - `PacketType`: Packet identifier
  - `AchievementID`: Achievement ID
  - `ObjectiveCount`: Number of objectives
  - `Objectives`: Objective progress data
  - `CompleteTime`: Completion timestamp
  - `Reward`: Reward data

### 0x0A26 - Achievement Reward
- Structure: `v V C`
- Fields:
  - `PacketType`: Packet identifier
  - `AchievementID`: Achievement ID
  - `Result`: Reward claim result

## Clan System

### 0x098A - Clan User Count
- Structure: `v v v`
- Fields:
  - `PacketType`: Packet identifier
  - `OnlineCount`: Number of online clan members
  - `TotalCount`: Total number of clan members

### 0x098D - Clan Information
- Structure: `v V Z24 Z24 Z16 v2`
- Fields:
  - `PacketType`: Packet identifier
  - `ClanID`: Clan ID
  - `ClanName`: Clan name
  - `MasterName`: Clan master name
  - `MapName`: Clan map name
  - `AllyCount`: Number of alliances
  - `AntagonistCount`: Number of antagonists
- Clan Status Codes:
  - 0: Active
  - 1: Disbanded
  - 2: In Formation
  - 3: Suspended

### 0x098F - Clan Alliance/Antagonist Data
- Structure: `v V v Z24 V`
- Fields:
  - `PacketType`: Packet identifier
  - `OtherClanID`: ID of allied/antagonist clan
  - `RelationType`: 0 for ally, 1 for antagonist
  - `ClanName`: Name of allied/antagonist clan
  - `RelationDuration`: Duration of the relationship in seconds

### 0x0991 - Clan Skill Information
- Structure: `v V v V`
- Fields:
  - `PacketType`: Packet identifier
  - `SkillID`: ID of the clan skill
  - `SkillLevel`: Current level of the skill
  - `ExpireTime`: Expiration time of the skill (0 if permanent)

### 0x098E - Clan Chat
- Structure: `v v Z24 Z*`
- Fields:
  - `PacketType`: Packet identifier
  - `MessageLength`: Length of the message
  - `SenderName`: Name of the message sender
  - `Message`: Chat message content

### 0x0990 - Clan Leave
- Structure: `v`
- Fields:
  - `PacketType`: Packet identifier

## Navigation System

### 0x08E2 - Navigate To
- Structure: `v C3 Z16 v2 v`
- Fields:
  - `PacketType`: Packet identifier
  - `Type`: Navigation type
  - `Flag`: Navigation flag
  - `Hide`: Hide flag
  - `MapName`: Target map name
  - `X`: X coordinate
  - `Y`: Y coordinate
  - `MonsterID`: Target monster ID (if applicable)
- Navigation Types:
  - 0: NPC
  - 1: Monster
  - 2: Coordinates
  - 3: Quest Objective
- Flags:
  - 0x01: Auto-routing
  - 0x02: Use waypoint system
  - 0x04: Avoid monsters
- Coordinate Handling:
  - For absolute positioning: X and Y are direct map coordinates
  - For relative positioning: X and Y are offsets from current position

## Expanded Inventory System

### 0x0B18 - Inventory Expansion Result
- Structure: `v v`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Expansion result

## Item Preview System

### 0x0999 - Item Preview
- Structure: `v v C C C a8 a*`
- Fields:
  - `PacketType`: Packet identifier
  - `Index`: Item index
  - `Identified`: Identification flag
  - `Broken`: Broken status
  - `Upgrade`: Upgrade level
  - `Cards`: Card data
  - `Options`: Item options data

## Reputation System

### 0x0B8D - Repute Info
- Structure: `v a*`
- Fields:
  - `PacketType`: Packet identifier
  - `ReputeInfo`: Reputation data for various types
- Reputation Types:
  - 1: PvP
  - 2: Quest
  - 3: Trade
  - 4: Crafting
  - 5: Hunting
- Point Calculation Methods:
  - PvP: (Kills * 10) - (Deaths * 5)
  - Quest: Completed quest points
  - Trade: (Successful trades * 2) + (Trade volume / 10000)
  - Crafting: (Successful crafts * 3) + (Craft quality bonus)
  - Hunting: (Monsters killed * 1) + (Boss monsters * 50)
- Reputation Ranks:
  - 0-999: Novice
  - 1000-4999: Apprentice
  - 5000-9999: Veteran
  - 10000+: Master
- Rank Benefits:
  - Novice: Basic access to reputation NPCs
  - Apprentice: 5% discount on reputation shops
  - Veteran: Access to special quests
  - Master: Unique title and cosmetic effects

## Gold PC Cafe System

### 0x0A15 - Gold PC Cafe Point
- Structure: `v C V V`
- Fields:
  - `PacketType`: Packet identifier
  - `IsActive`: Activation status
  - `Point`: Current points
  - `PlayedTime`: Time played in seconds

## Dynamic NPC System

### 0x0A17 - Dynamic NPC Create Result
- Structure: `v C`
- Fields:
  - `PacketType`: Packet identifier
  - `Result`: Creation result

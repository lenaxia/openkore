**Method Implementations:**

- guild_members_list - Guild member list handler (lines 6435-6478)
  - Handles ZC_MEMBERMGR_INFO (0154), 0AA5, and 0B7D packets
  - Processes different packet formats with varying member data structures:
    - 0B7D: 58 bytes, includes name (ID, charID, hair_style, hair_color, sex, jobID, lv, contribution, online, position, lastLoginTime, name)
    - 0AA5: 34 bytes, no name (ID, charID, hair_style, hair_color, sex, jobID, lv, contribution, online, position, lastLoginTime)
    - 0154: 104 bytes, includes memo (ID, charID, hair_style, hair_color, sex, jobID, lv, contribution, online, position, memo, name)
  - Converts name from bytes to string when present
  - Requests character name for 0AA5 packets
  - Stores member data in $guild{member} array

- guild_invite_result - Guild invite response handler (lines 6480-6503)
  - Handles ZC_ACK_REQ_JOIN_GUILD (0169) packet
  - Processes different response types:
    - 0 = Already in guild
    - 1 = Offer rejected
    - 2 = Offer accepted
    - 3 = Guild full
  - Displays appropriate translated messages for each response type
  - Handles unknown response types with fallback message

- guild_location - Guild member position handler (lines 6505-6522)
  - Handles ZC_NOTIFY_POSITION_TO_GUILDM (01EB) packet
  - Updates position for online guild members
  - Uses account ID to match members (supports multiple chars per account)
  - Skips updates for invalid coordinates (0,0)
  - Updates both current and destination positions

- guild_leave - Guild member leave handler (lines 6527-6548)
  - Handles ZC_ACK_LEAVE_GUILD (015A, 0A83) packets
  - Processes both direct name and charID based member identification
  - Removes leaving member from $guild{member} array
  - Displays formatted leave message with reason
  - Handles byte string conversion for name and reason fields

- guild_expulsion - Guild member expulsion handler (lines 6547-6569)
  - Handles ZC_ACK_BAN_GUILD (015C, 0839, 0A82) packets
  - Similar to guild_leave but for forced removals
  - Processes both direct name and charID based member identification
  - Removes expelled member from $guild{member} array
  - Displays formatted expulsion message with reason
  - Handles byte string conversion for name and reason fields

- guild_member_online_status - Guild member login/logout handler (lines 6571-6591)
  - Handles ZC_UPDATE_CHARSTAT (016D) and ZC_UPDATE_CHARSTAT2 (01F2) packets
  - Tracks online status (0=offline, 1=online) for guild members
  - Updates member status in $guild{member} array
  - Displays login/logout messages in guild chat
  - TODO: Could also update sex, hair_style, hair_color from ZC_UPDATE_CHARSTAT2

- guild_update_member_position - Guild member position update handler (lines 6593-6615)
  - Handles ZC_ACK_REQ_CHANGE_MEMBERS (0156) packets
  - Processes position changes for multiple guild members
  - Uses 12-byte member position info structure (ID, charID, position)
  - Updates member positions in $guild{member} array
  - Displays formatted position change messages
  - References $guild{positions} array for title lookups

- guild_members_title_list - Guild position titles handler (lines 6617-6630)
  - Handles ZC_POSITION_ID_NAME_INFO (0166) packets
  - Processes guild position title information:
    - Each entry contains position ID (4 bytes) and title name (24 bytes)
    - Stores titles in $guild{positions} array
    - Converts byte strings to Perl strings for title names
  - Packet format: <packet len>.W { <position id>.L <position name>.24B }*
  - Used by guild_update_member_position for title lookups

- guild_name - Guild basic information handler (lines 6632-6666)
  - Handles ZC_UPDATE_GDID (016C) packets
  - Processes guild information:
    - guildID (4 bytes)
    - emblemID (4 bytes)
    - mode flags (4 bytes):
      - 0x01 = allow invite
      - 0x10 = allow expel
    - ismaster flag (1 byte)
    - guildName (24 bytes)
  - Stores information in $char->{guild} hash
  - Server-specific behavior:
    - twRO: Requests basic info (0), hostile alliances (3), members list (1)
    - jRO: Requests members list (1)
    - Others: Requests master check, expulsion list (4), basic info (0), members list (1), emblem
  - Packet format: <guild id>.L <emblem id>.L <mode>.L <ismaster>.B <inter sid>.L <guild name>.24B

- guild_request - Guild invite request handler (lines 6668-6680)
  - Handles ZC_REQ_JOIN_GUILD (016A) packets
  - Processes guild invite requests:
    - guildID (4 bytes)
    - guildName (24 bytes)
  - Stores request information in %incomingGuild hash:
    - ID = guildID
    - Type = 1 (guild request)
  - Sets timeout for auto-deny if configured ($timeout{'ai_guildAutoDeny'})
  - Displays formatted invite message
  - Packet format: <guild id>.L <guild name>.24B

- guild_master_member - Guildmaster status handler (lines 6682-6703)
  - Handles guildmaster status notifications (ZC_ACK_GUILD_MENUINTERFACE, 014E)
  - Processes different status types:
    - 0xd7 = Guildmaster status
    - 0x57 = Not guildmaster status
  - Displays appropriate messages based on status
  - Handles unknown status types with warning
  - Packet format: <menu flag>.L
  - Menu flags:
    - 0x00 = Basic Info (always on)
    - 0x01 = Member manager
    - 0x02 = Positions
    - 0x04 = Skills
    - 0x10 = Expulsion list
    - 0x40 = Unknown (GMENUFLAG_ALLGUILDLIST)
    - 0x80 = Notice

- guild_alliance - Guild alliance result handler (lines 6705-6729)
 - Handles ZC_ACK_REQ_ALLY_GUILD (0173) packets
 - Processes different alliance result codes:
   - 0 = Already allied
   - 1 = You rejected the offer
   - 2 = You accepted the offer
   - 3 = They have too many alliances
   - 4 = You have too many alliances
   - 5 = Alliances are disabled
 - Displays appropriate translated messages for each result code
 - Handles unknown result codes with warning
 - Packet format: <answer>.B

- guild_member_setting_list - Guild position settings handler (lines 6731-6751)
 - Handles ZC_POSITION_INFO (0160) packets
 - Processes guild position settings:
   - Each position uses 16 bytes (position id, mode flags, ranking, pay rate)
   - Mode flags:
     - 0x01 = allow invite
     - 0x10 = allow expel
     - 0x100 = guild storage access
   - Stores settings in $guild{positions} array
   - Updates invite/punish/gstorage permissions based on mode flags
   - Stores feeEXP (pay rate) for each position
 - Packet format: <packet len>.W { <position id>.L <mode>.L <ranking>.L <pay rate>.L }*
 - TODO: Properly handle ranking field

- guild_skills_list - Guild skills information handler (lines 6753-6773)
 - Handles ZC_GUILD_SKILLINFO (0162) packets
 - Processes guild skill information:
   - Each skill uses 37 bytes (skillID, targetType, level, sp, range, skillName, up)
   - Skill data structure:
     - skillID (2 bytes)
     - targetType (4 bytes)
     - level (2 bytes)
     - sp cost (2 bytes)

- guild_expulsion_list - Guild ban list handler (lines 6775-6807)
 - Handles ZC_BAN_LIST (0163) and 0B7C packets
 - Processes two different packet formats:
   - 0B7C format (68 bytes per entry):
     - charID (4 bytes)
     - reason (40 bytes)
     - name (24 bytes)
   - 0163 format (88 bytes per entry):
     - name (24 bytes)
     - account name (24 bytes)
     - reason (40 bytes)
 - Clears existing expulsion list before processing
 - Converts name and reason fields from bytes to strings
 - Stores expulsion data in $guild{expulsion} array
 - Packet formats:
   - 0163: <packet len>.W { <char name>.24B <account name>.24B <reason>.40B }*
   - 0B7C: <packet len>.W { <charID>.L <reason>.40B <name>.24B }*
     - attack range (2 bytes)
     - skill name (24 bytes)

- guild_member_map_change - Guild member map change handler (lines 6809-6823)
 - Handles ZC_NOTIFY_MAPPROPERTY (01EC) packets
 - Processes guild member map changes:
   - Packet format: <account id>.L <char id>.L <status>.L <map name>.16B
   - Updates member's map information in guild data
   - Clears position (pos) and destination position (pos_to) data
   - Converts map name from bytes to string
   - Logs map change with debug message
   - Matches member by charID
   - Packet triggers when guild member changes maps
     - upgradable flag (1 byte)
   - Stores skill information in $guild{skills} hash
   - Converts skill name from bytes to string
   - Only sets level if not previously set
 - Packet format: <packet len>.W <skill points>.W { <skill id>.W <type>.L <level>.W <sp cost>.W <atk range>.W <skill name>.24B <upgradable>.B }*
 - TODO: Merge with skills_list handler

- guild_member_add - New guild member handler (lines 6824-6840)
 - Handles ZC_ADD_MEMBER_TO_GUILD (0182) and 0B7E packets
 - Processes new guild member information:
   - Packet format: <account>.L <char id>.L <hair style>.W <hair color>.W <gender>.W <class>.W <level>.W <contrib exp>.L <state>.L <position>.L <memo>.50B <name>.24B
   - Adds member to $guild{member} array
   - Converts member name from bytes to string
   - Displays join message in guild chat
   - Handles both old (0182) and new (0B7E) packet formats
   - Stores comprehensive member data:
     - Account ID, character ID
     - Appearance details (hair, gender, class)
     - Level and contribution
     - Position and memo
     - Character name

- guild_notice - Guild notice handler (lines 6842-6857)
 - Handles ZC_GUILD_NOTICE (016F) packets
 - Processes guild notice messages:
   - Packet format: <subject>.60B <notice>.120B
   - Strips language codes from subject and notice
   - Formats message with header/footer if content exists
   - Skips display if both subject and notice are empty
   - Output format:
     ---Guild Notice---
     [Subject]
     
     [Notice]
     ------------------

- guild_storage_log - Guild storage log handler (lines 9577-9613)
 - Handles guild storage access logs
 - Processes different result codes:
   - 0/1: Successful get/put operations
   - 2: Empty storage
   - 3: Storage not in use
 - Action types:
   - 0: Get item
   - 1: Put item
 - Log format details:
   - Item structure: ID, nameID, amount, action, upgrade, uniqueID, identified, type_equip, cards, charName, time, attribute
   - Displays formatted log with:
     - Character name
     - Item name
     - Amount
     - Action type
     - Timestamp
   - Includes header/footer separators

- guild_emblem - Guild emblem handler (lines 10368-10373)
  - Handles ZC_GUILD_EMBLEM_IMG (0152) packets
  - Currently just debug stub (TODO implementation)
  
- guild_emblem_update - Guild emblem update handler (lines 10375-10380)
  - Handles ZC_UPDATE_GDID (01B4) packets
  - Currently just debug stub (TODO implementation)

- char_emblem_update - Character emblem update handler (lines 10382-10387)
  - Handles ZC_UPDATE_CHARSTAT2 (0B47) packets
  - Currently just debug stub (TODO implementation)

- guild_position_changed - Guild position change handler (lines 10389-10394)
  - Handles ZC_POSITION_CHANGE (0174) packets
  - Currently just debug stub (TODO implementation)

- guild_position - Guild position handler (lines 10396-10401)
  - Handles ZC_GUILD_POSITION (0AFD) packets
  - Currently just debug stub (TODO implementation)

- guild_unally - Guild alliance removal handler (lines 10403-10408)
  - Handles ZC_DELETE_RELATED_GUILD (0184) packets
  - Currently just debug stub (TODO implementation)

- guild_opposition_result - Guild opposition result handler (lines 10410-10415)
  - Handles ZC_GUILD_OPPOSITION_RESULT (0181) packets
  - Currently just debug stub (TODO implementation)

- guild_alliance_added - Guild alliance addition handler (lines 10417-10422)
  - Handles ZC_ADD_RELATED_GUILD (0185) packets
  - Currently just debug stub (TODO implementation)
  - Note: This packet doesn't exist in eA

- guild_info() - Guild basic information handler (lines 6419-6433)
  - Handles ZC_GUILD_INFO, ZC_GUILD_INFO2, ZC_GUILD_INFO3 packets
  - Stores guild information in %guild hash
  - Converts byte strings to Perl strings for name/master fields
  - Increments member count automatically

- guild_members_list() - Guild member list handler (lines 6442-6478)
  - Handles ZC_MEMBERMGR_INFO packets
  - Supports multiple packet versions (0154, 0AA5, 0B7D)
  - Parses member data with different formats based on packet type
  - Stores member information in $guild{member} array
  - Handles character name resolution for some packet types

- guild_create_result() - Handles guild creation results (lines 6396-6417)
  - Processes ZC_RESULT_MAKE_GUILD packet
  - Handles different creation result types:
    - 0 = Success
    - 1 = Already in guild
    - 2 = Guild name exists
    - 3 = Missing Emperium item
  - Parameters:
    - $args->{type} - Result code (0-3)
  - Comments:
    - "0167 <result>.B" (packet format)
    - Detailed result code descriptions in comments

- guild_ally_request() - Handles guild alliance requests (lines 6359-6371)
  - Processes ZC_REQ_ALLY_GUILD packet
  - Stores inviter's account ID and guild name
  - Sets timeout for auto-deny if configured
  - Parameters:
    - $args->{ID} - Inviter's account ID
    - $args->{guildName} - Guild name (24B string)
  - Comments:
    - "0171 <inviter account id>.L <guild name>.24B" (packet format)
    - "Freya calls it an account ID" (note about ID type)

- guild_broken() - Handles guild break notifications (lines 6378-6394)
  - Processes ZC_ACK_DISORGANIZE_GUILD_RESULT packet
  - Handles different break result flags:
    - 0 = success (clears guild data)
    - 1 = invalid key
    - 2 = members still in guild
  - Parameters:
    - $args->{flag} - Result code (0-2)
  - Comments:
    - "015E <reason>.L" (packet format)
    - Detailed flag descriptions in comments

- guild_allies_enemy_list() - Guild alliance/enemy list handler (lines 6330-6357)
  - Handles ZC_MYGUILD_BASIC_INFO packet (014C)
  - Packet format: <packet len>.W { <relation>.L <guild id>.L <guild name>.24B }*
  - Clears and rebuilds guild ally/enemy lists
  - Processes each guild entry (32 bytes each):
    - type=0 Ally
    - type=1 Enemy
  - Stores guild IDs and names in %guild hash
  - Debugs ally/enemy relationships

# Emblem Related Handlers

**Method Implementations:**
- char_emblem_update - Character emblem update handler (lines 10384-10387)
  - Processes character emblem update notifications
  - Outputs debug message with packet information
  - Contains TODO comment indicating incomplete implementation
  - Packet: 0B47
  - Simple implementation focused on debugging

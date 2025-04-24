# Skill Related Handlers

**Method Implementations:**

- skill_post_delay - Skill cooldown icon handler (lines 3389-3398)
  - Processes skill cooldown notifications (ZC_SKILL_POSTDELAY)
  - Creates Skill object from skill ID
  - Gets skill name from the Skill object
  - Uses statusName for 'EFST_DELAY' or falls back to 'Delay'
  - Sets character status with format: "[SkillName] Delay"
  - Applies the cooldown time from packet

- skill_post_delaylist - Skill cooldown list handler (lines 3400-3430)
  - Processes multiple skill cooldowns at once
  - Supports two packet formats:
    * 043E: Basic format with skill ID and remaining time
    * 0985: Extended format with total time included
  - Dynamically selects unpacking format based on packet type
  - Iterates through skill list data in fixed-length chunks
  - Creates Skill objects to get skill names
  - Sets character status for each skill cooldown

- gospel_buff_aligned - Gospel buff message handler (lines 3432-3474)
  - Processes gospel skill effect messages (ZC_SKILLMSG)
  - Maps numeric status codes to descriptive messages
  - Handles various gospel effects:
    * Status removal (21)
    * Status immunity (22)
    * Max HP/SP increases (23-24)
    * Stat increases (25)
    * Holy weapon/armor enchantment (28-29)
    * Defense increase (30)
    * Attack increase (31)
    * Accuracy/Flee increase (32)
  - Displays appropriate info messages for each effect
  - Contains commented code for unknown effects

- parse_sage_autospell - Sage autospell parser (lines 3150-3158)
  - Processes autospell skill list data
  - Handles both regular autospell and shadow spell variants
  - Unpacks skill IDs from packet data:
    * Uses autoshadowspell_list (v* format) if available
    * Falls back to autospell_list (V* format)
  - Creates Skill objects for each ID
  - Sorts skills numerically for consistent display

- reconstruct_sage_autospell - Sage autospell packet builder (lines 3160-3166)
  - Converts skill objects back to packet format
  - Extracts skill IDs from provided skill objects
  - Packs data in both formats:
    * autoshadowspell_list: 2-byte format (v*)
    * autospell_list: 4-byte format (V*)
  - Used for packet reconstruction

- sage_autospell - Autospell skill list handler (lines 3168-3206)
  - Handles Sage's Hindsight and Shadow Chaser's Auto Shadow Spell
  - Requires in-game state (changeToInGameState)
  - Displays formatted list of available autospell skills
  - Implements autoSpell configuration:
    * Parses comma-separated skill list from config
    * Tests each configured skill against available skills
    * Respects autoSpell_safe setting for skill validation
    * Sends appropriate packet based on context:
      - sendSkillSelect if 'why' parameter is provided
      - sendAutoSpell otherwise
  - Provides helpful error messages and hints
- skill_use_failed - Skill failure handler (lines 11744-11837)
  - Processes skill use failure notifications
  - Defines extensive error type mappings:
    * basefailtype: Basic failure types (emotions, sit, chat, etc.)
    * failtype: Detailed failure reasons (SP, HP, requirements, etc.)
  - Determines error message based on:
    * For skillID 1 and cause 0: Uses basefailtype with btype
    * For known causes: Uses failtype with cause
    * For cause 71: Appends item ID to message
    * For unknown errors: Uses "Unknown error"
  - Clears character's casting state
  - Sets up hook arguments with:
    * skillID, btype, itemId, flag, cause
    * failMessage and warn flag
  - Triggers packet_skillfail hook
  - Displays warning message with:
    * Skill name, error message, cause number
    * Uses "skill" warning category
  - Handles special homunculus skill cases:
    * Skill 247 (Resurrect Homunculus): Updates dead flag
    * Skill 243 (Call Homunculus): Updates vaporized flag
  - Contains detailed debug messages for homunculus skills
- skill_add - Skill addition handler (lines 11712-11738)
  - Processes skill addition notifications
  - Requires in-game state (changeToInGameState)
  - Gets skill handle from name or creates from skillID
  - Updates character's skills hash with:
    * ID, SP, range, upgradable status
    * Target type, level
    * Sets new flag to 1
  - Adds handle to skillsID array if not already present
  - Contains comment about fixing bug with "Night" status received twice
  - Adds skill to Skill::DynamicInfo with:
    * skillID, handle, level, SP, target
    * Sets owner to Skill::OWNER_CHAR
  - Triggers packet_charSkills hook with:
    * ID, handle, level, upgradable, level2
  - Complex implementation for skill management
- cast_cancelled - Skill cast cancellation handler (lines 11568-11584)
  - Processes skill cast cancellation notifications
  - Gets actor reference using ID
  - Sets cast_cancelled timestamp
  - Gets skill name from casting data (or "Unknown")
  - Determines message domain (selfSkill vs skill)
  - Displays formatted message with source and skill name
  - Triggers packet_castCancelled hook with sourceID
  - Deletes casting state from source actor
  - Contains comment about packet format (01B9)
  - Contains comment about packet purpose (ZC_DISPEL)
- skill_cast - Skill casting handler (lines 11450-11564)
  - Processes skill casting notifications
  - Requires in-game state (changeToInGameState)
  - Gets source and target actors
  - Creates Skill object from skillID
  - Sets up casting state in source actor:
    * Stores skill, target, coordinates, start time, cast time
    * Uses Scalar::Util::weaken to prevent memory leaks
  - Handles target determination:
    * For ground-targeted skills (x,y coordinates):
      - Calculates distance to skill area
      - Sets target string to location
      - Unsets targetID
    * For actor-targeted skills:
      - Uses target's nameString
  - Performs trigger actions for self-cast:
    * Sets time_cast and time_cast_wait
    * Clears cast_cancelled flag
  - Calls countCastOn helper function
  - Determines message domain (selfSkill vs skill)
  - Displays formatted message using skillCast_string
  - Triggers is_casting hook with detailed parameters
  - Handles monster skill cancellation:
    * Gets monster control settings
    * For AI::AUTO with skillcancel_auto:
      - Switches target to casting monster
      - Stops current attack
      - Dequeues AI actions
      - Starts new attack on caster
    * For area skills:
      - Calculates position behind target monster
      - Uses vector math to determine movement
      - Routes character to safe position
      - Displays avoidance message
  - Uses Misc::checkValidity at multiple points
  - Contains TODO comment about 'dispose' support

- starplace - Star Gladiator map confirmation handler (lines 12201-12206)
  - Processes Star Gladiator's Feeling map confirmation (0253)
  - Displays which value from server
  - Simple implementation with minimal functionality
  - Packet name: ZC_STARPLACE

- area_spell - Area effect spell handler (lines 4267-4313)
  - Processes area effect spell notifications (including traps)
  - Supports multiple packet formats:
    * 011F: Basic format with ID, sourceID, x, y, type, isVisible
    * 01C9: Extended format with scribble message support
    * 08C7: Newer format with range parameter
  - Manages spell tracking:
    * Adds spell to spellsID array
    * Updates spell properties in spells hash
    * Tracks position, type, visibility
  - Special handling for warp portals (type 0x81)
  - Displays debug message with spell name and location
  - Special handling for scribble messages (01C9 packet)
  - Triggers packet_areaSpell hook with all parameters
  - Packets: 011F, 01C9, 08C7

- area_spell_multiple2 - Multiple area spells handler (lines 4315-4361)
  - Processes multiple area effect spells in a single packet
  - Parses variable-length data with fixed-size entries (18 bytes each)
  - Unpacks each entry with format: 'a4 a4 v2 V C2'
  - Extracts: ID, sourceID, x, y, type, range, isVisible
  - Similar spell tracking to area_spell:
    * Adds spells to spellsID array
    * Updates spell properties in spells hash
    * Tracks position, type, visibility, range
  - Special handling for warp portals (type 0x81)
  - Displays debug message for each spell
  - Triggers packet_areaSpell hook with parameters from last spell
  - Packet: 099F

- area_spell_multiple3 - Multiple area spells with level handler (lines 4363-4411)
  - Processes multiple area effect spells with level information
  - Parses variable-length data with fixed-size entries (19 bytes each)
  - Unpacks each entry with format: 'a4 a4 v2 V C3'
  - Extracts: ID, sourceID, x, y, type, range, isVisible, lvl
  - Similar spell tracking to area_spell_multiple2:
    * Adds spells to spellsID array
    * Updates spell properties in spells hash
    * Tracks position, type, visibility, range, level
  - Special handling for warp portals (type 0x81)
  - Displays debug message for each spell including level
  - Triggers packet_areaSpell hook with parameters from last spell
  - Packet: 09CA
- skill_delete - Skill removal handler (lines 9614-9624)
  - Processes skill removal notifications
  - Creates Skill object from skillID
  - Performs validation checks:
    * Returns early if skill object creation fails
    * Returns early if skill doesn't exist in character's skills
  - Displays "Lost skill" message with skill name
  - Removes skill from character's skills hash
  - Removes skill handle from skillsID array
  - Uses "skill" message category
  - Simple implementation focused on skill removal
- skill_exchange_item - Material exchange skill handler (lines 7863-7880)
  - Processes material exchange skill notifications
  - Handles different skill types:
    * Type 0: Change Material
      - Displays "Change Material is ready" message
      - Suggests using 'cm' command
    * Other types: Four Spirit Analysis
      - Displays "Four Spirit Analysis is ready" message
      - Suggests using 'analysis' command
  - Contains detailed comments about type values:
    * 0: Change Material
    * 1: Elemental Analysis Lv 1
    * 2: Elemental Analysis Lv 2
  - Sets skillExchangeItem global variable to type+1
  - Uses "info" message category
  - Simple implementation focused on notification
- skill_msg - Skill message handler (lines 9285-9294)
  - Processes skill-related messages
  - Looks up message ID in msgTable array
  - For known message IDs:
    * Creates Skill object from skill ID
    * Gets skill name from the Skill object
    * Displays formatted message with skill name and message text
    * Uses "info" message category
  - For unknown message IDs:
    * Displays warning about missing message
    * Suggests updating msgstringtable.txt from data.grf
    * Includes message ID and skill ID in warning
    * Uses "warning" message category
  - Simple implementation focused on message display
- skills_list - Skills list handler (lines 9326-9393)
  - Processes skill list notifications for different actors
  - Requires in-game state (changeToInGameState)
  - Supports multiple packet formats:
    * 0B32: Compact format with 15-byte entries
    * Others: Extended format with 37-byte entries including skill handle
  - Determines owner type and appropriate hook based on packet switch:
    * 010F/0B32: Character skills (OWNER_CHAR)
    * 0235: Homunculus skills (OWNER_HOMUN)
    * 029D: Mercenary skills (OWNER_MERC)
  - Selects appropriate skillsID reference based on owner type
  - Clears existing skills from the actor's skills hash
  - Parses variable-length data with fixed-size entries
  - For each skill entry:
    * Unpacks skill data using appropriate format
    * Gets skill handle from ID
    * Updates actor's skills hash with all skill properties
    * Adds handle to skillsID array if not already present
    * Adds skill to Skill::DynamicInfo
    * Triggers appropriate hook with skill details
  - Contains TODO comments about moving skillsID to Actor
  - Complex implementation for comprehensive skill management
  - Packets: 010F, 0235, 029D, 0B32
- skill_update - Skill update handler (lines 9395-9420)
  - Processes skill update notifications
  - Extracts skill details from args:
    * skillID, lv, sp, range, up
  - Creates Skill object from skillID
  - Gets handle and name from the Skill object
  - Updates character's skills hash with:
    * Level, SP cost, range, upgradable status
  - Adds skill to Skill::DynamicInfo with:
    * ID, handle, level, SP, range, target type
    * Sets owner to Skill::OWNER_CHAR
  - Triggers packet_charSkills hook with:
    * ID, handle, level, upgradable, level2
  - Outputs debug message with skill name and level
  - Contains TODO comment about using type parameter
  - Simple implementation focused on skill data update

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
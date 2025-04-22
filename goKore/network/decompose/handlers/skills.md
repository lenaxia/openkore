**Method Implementations:**
- parse_sage_autospell - Sage auto-spell parser (lines 3150-3158)
  - Processes auto-spell lists for Sage class
  - Handles both autoshadowspell_list (v* format) and autospell_list (V* format)
  - Creates Skill objects from unpacked skill IDs
  - Sorts skill IDs numerically
  - Filters out empty/undefined skill IDs

- reconstruct_sage_autospell - Sage auto-spell reconstructor (lines 3160-3168)
  - Rebuilds auto-spell lists for network transmission
  - Converts Skill objects back to packed format
  - Maintains compatibility with both spell list types

- sage_autospell - Sage auto-spell display handler (lines 3174-3185)
  - Formats and displays auto-spell list to user
  - Requires in-game state
  - Creates centered header with title
  - Lists skills with ID and name in formatted columns
  - Uses internationalized strings (T())
  - Outputs to 'list' message channel

- skill_post_delay - Single skill cooldown handler (lines 3391-3398)
  - Handles ZC_SKILL_POSTDELAY packet (043D)
  - Creates skill name from ID
  - Sets status effect showing cooldown
  - Uses EFST_DELAY status type
  - Applies cooldown timer to character

- skill_post_delaylist - Multiple skill cooldown handler (lines 3403-3430)
  - Handles both ZC_SKILL_POSTDELAYLIST (043E) and ZC_SKILL_POSTDELAYLIST2 (0985) packets
  - Processes list of skill cooldowns
  - Supports two different packet formats:
    - 043E: Basic format with skill ID and remaining time
    - 0985: Extended format with total time and remaining time
  - Creates skill names from IDs
  - Sets status effects showing cooldowns
  - Uses EFST_DELAY status type
  - Applies cooldown timers to character

- gospel_buff_aligned - Gospel skill buff messages (lines 3447-3474)
  - Handles ZC_SKILLMSG packet (0215)
  - Processes various Gospel skill effect messages
  - Status IDs and corresponding messages:
    - 21: All abnormal status effects removed
    - 22: Immunity to abnormal status effects
    - 23: Max HP increased
    - 24: Max SP increased
    - 25: All stats increased
    - 28: Weapon blessed with Holy power
    - 29: Armor blessed with Holy power
    - 30: Defense increased
    - 31: Attack strength increased
    - 32: Accuracy and Flee Rate increased
  - Uses internationalized strings (T())
  - Outputs to 'info' message channel
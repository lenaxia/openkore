**Method Implementations:**

- unit_levelup - Level up and special effects handler (lines 2809-2844)
  - Processes various visual effects notifications
  - Handles multiple effect types:
    * 0/7/9: Base level up effects (regular, super novice, taekwon)
    * 1/8: Job level up effects (regular, super novice)
    * 2: Refining failure effect
    * 3: Refining success effect
    * 4: Game over effect
    * 5: Pharmacy (potion creation) success effect
    * 6: Pharmacy (potion creation) failure effect
  - Displays appropriate messages based on effect type
  - Triggers plugin hooks for base_level and job_level events
  - Uses "refine" message category for refining-related effects
- blade_stop - Blade Stop skill effect handler (lines 10433-10440)
  - Processes Blade Stop skill effect notifications
  - Handles two active values:
    * 0: Deactivated - "Blade Stop by X on Y is deactivated"
    * 1: Active - "Blade Stop by X on Y is active"
  - Gets source and target actor names using Actor::get and nameString
  - Uses "info" message category for all messages
  - Contains TODO comment: "the actual status is sent to us in opt3"
  - Packet: 01D1
- area_spell_disappears - Area spell disappearance handler (lines 10151-10159)
  - Processes area spell disappearance notifications
  - Gets spell ID from packet
  - Gets spell from spells hash
  - Outputs debug message with:
    * Spell name (from getSpellName)
    * Spell binID
    * Source actor name (from getActorName)
    * Position coordinates
  - Deletes spell from spells hash
  - Removes ID from spellsID array
  - Uses "skill" debug category with level 2
- hat_effect - Hat effect display handler (lines 7375-7406)
  - Processes hat effect display notifications
  - Gets actor reference using Actor::get
  - Builds effect name string from effect list:
    * Looks up effect handle in hatEffectHandle hash
    * Uses hatEffectName if available, otherwise uses handle
    * Falls back to "Unknown #[ID]" for undefined effects
    * Joins multiple effects with commas
  - Handles two flag states:
    * 1 = Effect active: "[Actor] use/uses effect: [effects]"
    * Other = Effect removed: "[Actor] are/is no longer: [effects]"
  - Uses actor's verb() method for grammatically correct messages
  - Outputs to 'effect' message category
  - Contains TODO comment about storing effects in actor
- parse_hat_effect - Hat effect parser (lines 7367-7371)
  - Processes hat effect data from raw packet
  - Unpacks effect information into HatEFID structure
  - Uses complex unpacking pattern:
    * Extracts effect data as series of 2-byte values
    * Maps each value to a hash with HatEFID key
  - Outputs debug message with flag and effect IDs
  - Used by hat_effect handler
  - Packet: 0A3B
- sound_effect - Sound effect handler (lines 6882-6896)
  - Processes sound effect notifications (ZC_SOUND)
  - Gets actor reference if ID exists
  - Handles different sound types with appropriate messages:
    * 0 = play once: "[Actor] play/plays: [sound]"
    * 1 = play repeat: "[Actor] are/is now playing: [sound]"
    * 2 = stop: "[Actor] stopped playing: [sound]"
  - Uses actor's verb() method for grammatically correct messages
  - Falls back to "Now playing: [sound]" when no actor
  - Outputs to 'effect' message category
  - Comment notes continuous sounds could be actor statuses
  - Packet details:
    * File name relative to data\wav
    * act: 0=play once, 1=repeat, 2=stop
    * term: unknown purpose for act=1
- misc_effect - Miscellaneous effect display handler (lines 6861-6869)
  - Processes generic visual effect notifications
  - Gets actor reference using Actor::get
  - Uses actor's verb() method for grammatically correct messages
  - Displays effect name from effectName lookup table
  - Falls back to "Unknown #[effect]" for undefined effects
  - Outputs to 'effect' message category
  - Simple implementation focused on effect notification
- revolving_entity - Character special effects handler (lines 3999-4044)
  - Handles monk spirits, gunslinger coins, ninja amulets, and soul energy
  - Processes packets: 01D0 (monk spirits), 01E1 (gunslinger coins), 08CF (ninja amulets), 0B73 (soul energy)
  - Updates character or actor spirit count and type
  - Handles elemental properties when applicable
  - Displays appropriate messages based on entity type and owner
  - Supports both self and other actors
- minimap_indicator - Minimap indicator handler (lines 3071-3099)
  - Handles showing/clearing minimap indicators
  - Takes parameters: show (bool), actor (Actor), x/y coordinates, RGB color values
  - Supports special effects like quest markers (effect=1)
  - Logs indicator changes with color information
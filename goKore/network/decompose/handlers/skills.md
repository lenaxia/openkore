**Skill Messaging:**
- skill_msg() - Skill-related messages (lines 9286-9294)
  - Displays messages from msgstringtable.txt
  - Formats with skill name
  - Handles unknown message IDs

- message_string() - General game messages (lines 9302-9324)
  - Processes various message types:
    * Simple messages (ZC_MSG_COLOR)
    * Messages with values (ZC_MSG_VALUE)
    * Formatted strings (ZC_FORMATSTRING_MSG)
  - Handles mercenary-related messages
  - Triggers packet_message_string hook

**Skill Management:**
- skills_list() - Updates skill lists (lines 9327-9393)
  - Handles different owner types:
    * Character (010F, 0B32)
    * Homunculus (0235)
    * Mercenary (029D)
  - Processes skill details:
    * ID, level, SP cost
    * Range, target type
    * Upgradable status
  - Maintains dynamic skill info
  - Triggers appropriate hooks

- skill_update() - Updates single skill (lines 9396-9417)
  - Processes skill level/SP changes
  - Updates character skill data
  - Maintains dynamic skill info
  - Triggers packet_charSkills hook

**Targeted Skill Handlers:**

**Skill Exchange Handlers:**
- skill_exchange_item (lines 7864-7880)
  - Handles skill exchange item preparation
  - Features:
    - Supports two types:
      - Change Material (type 0)
      - Four Spirit Analysis (type 1)
    - Provides appropriate user commands:
      - 'cm' for Change Material
      - 'analysis' for Four Spirit Analysis
    - Maintains skill exchange state
- devotion (lines 5929-5946)
  - Manages Devotion skill (target protection)
  - Tracks:
    - Source actor (devotion caster)
    - Up to 5 target actors
    - Protection range
  - Builds formatted message showing protection links
  - Stores data in devotionList structure

**Area Effect Handlers:**
- area_spell_multiple3 (lines 4364-4411)
  - Processes multiple area effect spells/traps
  - Handles warp portal detection (type 0x81)
  - Tracks spell properties:
    - Position (x,y)
    - Type and range
    - Visibility and level
  - Maintains spells{} structure with:
    - Source ID
    - Bin ID
    - Position data
  - Triggers packet_areaSpell hook
  - Debug logs spell details when enabled

**Skill Cooldown Handlers:**
- skill_post_delay (lines 3390-3398)
  - Single skill cooldown display
  - Shows delay status with skill name
  - Uses EFST_DELAY status type

- skill_post_delaylist (lines 3400-3430)
  - Multiple skill cooldowns display
  - Handles two packet versions (043E and 0985)
  - Processes packed skill list data
  - Supports both total_time and remain_time fields

**Gospel Buff Handler:**
- gospel_buff_aligned (lines 3447-3474)
  - Handles Gospel skill buff messages
  - Processes 11 different buff types:
    - Status effect removal (0x15)
    - Status immunity (0x16)
    - Max HP/SP increases (0x17-0x18)
    - Stat increases (0x19)
    - Holy element enchantments (0x1c-0x1d)
    - DEF/ATK/HIT/Flee boosts (0x1e-0x20)
  - Shows appropriate duration messages

- skill_post_delay (lines 3390-3398)
  - Single skill cooldown display
  - Shows delay status with skill name
  - Uses EFST_DELAY status type

- skill_post_delaylist (lines 3400-3430)
  - Multiple skill cooldowns display
  - Handles two packet versions (043E and 0985)
  - Processes packed skill list data
  - Supports both total_time and remain_time fields

#### Sage Auto Spell (lines 3174-3206)
```perl
sub sage_autospell {
	my ($self, $args) = @_;
	
	# Displays list of available auto spells
	my $msg = center(' ' . T('Auto Spell') . ' ', 40, '-') . "\n"
	. T("   # Skill\n")
	. (join '', map { swrite '@>>> @<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<', [$_->getIDN, $_] } @{$args->{skills}})
	. ('-'x40) . "\n";
	
	message $msg, 'list';

	# Handle auto spell configuration
	if ($config{autoSpell}) {
		my @autoSpells = split /\s*,\s*/, $config{autoSpell};
		for my $autoSpell (@autoSpells) {
			my $skill = new Skill(auto => $autoSpell);
			if (!$config{autoSpell_safe} || List::Util::first { $_->getIDN == $skill->getIDN } @{$args->{skills}}) {
				$messageSender->sendAutoSpell($skill->getIDN);
				return;
			}
		}
		error TF("Configured autoSpell (%s) not available.\n", $config{autoSpell});
	}
}
```

#### Skill Cooldown (lines 3391-3398, 3403-3430)
```perl
sub skill_post_delay {
	my ($self, $args) = @_;
	my $skillName = (new Skill(idn => $args->{ID}))->getName;
	my $status = defined $statusName{'EFST_DELAY'} ? $statusName{'EFST_DELAY'} : 'Delay';
	$char->setStatus($skillName." ".$status, 1, $args->{time});
}

sub skill_post_delaylist {
	my ($self, $args) = @_;
	my $skill_post_delay_info;
	if ($args->{switch} eq "0985") {
		$skill_post_delay_info = {
			len => 10,
			types => 'v V2',
			keys => [qw(ID total_time remain_time)],
		};
	} else {
		$skill_post_delay_info = {
			len => 6,
			types => 'v V',
			keys => [qw(ID remain_time)],
		};
	}

	for (my $i = 0; $i < length($args->{skill_list}); $i += $skill_post_delay_info->{len}) {
		my $skill;
		@{$skill}{@{$skill_post_delay_info->{keys}}} = unpack($skill_post_delay_info->{types}, substr($args->{skill_list}, $i, $skill_post_delay_info->{len}));
		$skill->{name} = (new Skill(idn => $skill->{ID}))->getName;
		my $status = defined $statusName{'EFST_DELAY'} ? $statusName{'EFST_DELAY'} : 'Delay';
		$char->setStatus($skill->{name}." ".$status, 1, $skill->{remain_time});
	}
}
```

### Key Features:
- Manages Sage's Auto Spell and Shadow Chaser's Auto Shadow Spell
- Displays formatted list of available spells
- Supports configurable auto spell selection (autoSpell)
- Implements safety check (autoSpell_safe)
- Integrates with skill system and configuration
- Handles both direct spell selection and why-based selection
- Provides user feedback about spell availability
- Handles skill cooldown timers:
  - Single skill (skill_post_delay)
  - Multiple skills (skill_post_delaylist)
  - Supports different packet versions
  - Maintains status effects for cooldown display
  - Properly unpacks skill data structures
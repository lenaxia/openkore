**Monster System Handlers:**

- boss_map_info() - MVP boss information (lines 9539-9555)
  - Handles MVP boss detection packets:
    * No boss detected (flag 0)
    * Boss location (flag 1)
    * Boss detected on map (flag 2)
    * Boss respawn timer (flag 3)
  - Displays appropriate messages
  - Handles unknown flag cases

**Actor Handlers:**
- actor_look_at() - Updates actor's look direction (lines 8121-8129)
  - Updates head/body look direction for specified actor
  - Logs look change in parseMsg debug channel
  - Requires actor ID, head and body parameters

- actor_movement_interrupted() - Handles interrupted movement (lines 8136-8158)
  - Resets actor position to specified coordinates
  - Clears movement calculations
  - Handles special cases for player/homunculus
  - Logs debug message for player character

- actor_trapped() - Handles trapped actor state (lines 8160-8166)
  - Logs debug message when actor is trapped
  - Works with valid actor IDs (contrary to old comment)

**Monster Status Handlers:**

**Player Clone Handlers:**
- offline_clone_found (lines 7954-7992)
  - Handles detection of offline player clones
  - Features:
    - Creates new player actor if not existing
    - Sets clone flag
    - Populates all player attributes:
      - Name, ID, position
      - Job, equipment, appearance
      - Movement parameters
    - Adds to players list
    - Triggers player-related hooks:
      - add_player_list
      - player (legacy)
      - player_exist
- monster_hp_info (lines 4073-4082)
  - Shows detailed monster HP information
  - Processes ZC_HP_INFO packets
  - Updates monster HP and max HP
  - Displays debug message with HP percentage

- monster_hp_info_tiny (lines 4086-4094)
  - Shows approximate monster HP percentage
  - Updates monster's hp_percent field
  - Displays debug message with estimated HP

**Entity State Handlers:**
- revolving_entity (lines 3999-4044)
  - Tracks special entity states:
    - Monk spirits (01D0)
    - Gunslinger coins (01E1)
    - Warlock amulets (08CF)
    - Soul energy (0B73)
  - Updates entity counts and types
  - Shows state change messages
  - Handles different packet types:
    - ZC_SPIRITS
    - ZC_SPIRITS2
    - ZC_SPIRITS3

**Monster Sprite Handler:**
- monster_typechange (lines 4051-4055)
  - Processes NPC sprite changes
  - Handles ZC_NPCSPRITE_CHANGE packets
  - Updates monster appearance
  - Currently unused type parameter

**Actor Action Handler:**
- actor_action (lines 3722-3755)
  - Processes actor action animations
  - Handles ZC_ACTION_FAILURE packets
  - Key features:
    - Shows attack/miss animations
    - Updates actor state (attacking/moving)
    - Supports different action types:
      - ATTACK (0)
      - SIT (1)
      - STAND (2)
      - DAMAGE (3)
    - Triggers action hooks:
      - 'packet_actorAction'
      - 'packet_actorAction_$type'
    - Maintains action cooldown timers

**NPC Script Display:**
- show_script (lines 3374-3387)
  - Handles NPC script messages
  - Processes NPC ID and message content
  - Triggers 'show_script' plugin hook
  - Supports UTF-8 message decoding

# Actor Management Handlers

## Core Actor Display Methods

### actor_display_compatibility() - lines 1833-1839
```perl
sub actor_display_compatibility {
    my ($self, $args) = @_;
    # compatibility; TODO do it in PacketParser->parse?
    Plugins::callHook('packet_pre/actor_display', $args);
    &actor_display unless $args->{return};
    Plugins::callHook('packet/actor_display', $args);
}
```

### actor_display() - lines 1842-2100
```perl
sub actor_display {
    my ($self, $args) = @_;
    return unless changeToInGameState();
    my ($actor, $mustAdd);

    #### Initialize ####
    my $nameID = unpack("V", $args->{ID});
    my $name = bytesToString($args->{name});
    $name =~ s/^\s+|\s+$//g;

    if ($args->{switch} eq "0086") {
        # Message 0086 contains less information about the actor than other similar
        # messages. So we use the existing actor information.
        my $coordsArg = $args->{coords};
        my $tickArg = $args->{tick};
        $args = Actor::get($args->{ID})->deepCopy();
        # Here we overwrite the $args data with the 0086 packet data.
        $args->{switch} = "0086";
        $args->{coords} = $coordsArg;
        $args->{tick} = $tickArg;
    }
```

## Actor Type Classification

```perl
    #### Step 0: determine object type ####
    my $object_class;
    if (defined $args->{object_type}) {
        if ($args->{type} == 45) { # portals use the same object_type as NPCs
            $object_class = 'Actor::Portal';
        } else {
            $object_class = {
                PC_TYPE, 'Actor::Player',
                NPC_MOB_TYPE, 'Actor::Monster',
                NPC_EVT_TYPE, 'Actor::NPC', # both NPCs and portals
                NPC_PET_TYPE, 'Actor::Pet',
                NPC_HO_TYPE, 'Actor::Slave::Homunculus',
                NPC_MERSOL_TYPE, 'Actor::Slave::Mercenary',
                NPC_TYPE2, 'Actor::NPC',
            }->{$args->{object_type}};
        }
    }

    unless (defined $object_class) {
        if ($jobs_lut{$args->{type}}) {
            if ($args->{type} <= 6000) {
                $object_class = 'Actor::Player';
            } elsif (($args->{type} >= 6001 && $args->{type} <= 6016) || ($args->{type} >= 6048 && $args->{type} <= 6052)) {
                $object_class = 'Actor::Slave::Homunculus';
            } elsif ($$args->{type} >= 6017 && $$args->{type} <= 6046) {
                $object_class = 'Actor::Slave::Mercenary';
            } else {
                $object_class = 'Actor::Slave::Unknown';
            }
        } elsif ($args->{type} == 45) {
            $object_class = 'Actor::Portal';
        } elsif ($args->{type} >= 1000) {
            if ($args->{hair_style} == 0x64) {
                $object_class = 'Actor::Pet';
            } else {
                $object_class = 'Actor::Monster';
            }
        } else {   # ($args->{type} < 1000 && $args->{type} != 45 && !$jobs_lut{$args->{type}})
            $object_class = 'Actor::NPC';
        }
    }
```

## Actor Creation and Management

```perl
    #### Step 1: create the actor object ####

    if ($object_class eq 'Actor::Player') {
        # Actor is a player
        $actor = $playersList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::Player();
            $actor->{appear_time} = time;
            $actor->{name} = $name if defined $name;
            $mustAdd = 1;
        }
        $actor->{nameID} = $nameID;
    } elsif ($object_class eq 'Actor::Slave') {
        require ErrorHandler;
        die "Unset Actor::Slave type, this shouldn't happen\n";
    } elsif ($object_class eq 'Actor::Slave::Homunculus' || $object_class eq 'Actor::Slave::Mercenary' || $object_class eq 'Actor::Slave::Unknown') {
        # Actor is a homunculus or a mercenary
        $actor = $slavesList->getByID($args->{ID});
        if (!defined $actor) {
            if ($char->{slaves} && $char->{slaves}{$args->{ID}}) {
                $actor = $char->{slaves}{$args->{ID}};
            } elsif ($char->{homunculus} && $char->{homunculus}{ID} && $char->{homunculus}{ID} eq $args->{ID}) {
                $actor = $char->{homunculus};
            } elsif ($char->{mercenary} && ($char->{mercenary}{ID} && $char->{mercenary}{ID} eq $args->{ID})) {
                $actor = $char->{mercenary};
            } else {
                $actor = $object_class->new();
            }

            $actor->{appear_time} = time;
            $actor->{name_given} = $name if defined $name;
            $actor->{jobID} = $args->{type} if exists $args->{type};
            $mustAdd = 1;
        }
        $actor->{nameID} = $nameID;
    } elsif ($object_class eq 'Actor::Portal') {
        # Actor is a portal
        $actor = $portalsList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::Portal();
            $actor->{appear_time} = time;
            my $exists = portalExists($field->baseName, \%coordsTo);
            $actor->{source}{map} = $field->baseName;
            if ($exists ne "") {
                $actor->setName("$portals_lut{$exists}{source}{map} -> " . getPortalDestName($exists));
            }
            $mustAdd = 1;
        }
        $actor->{nameID} = $nameID;
    } elsif ($object_class eq 'Actor::Pet') {
        # Actor is a pet
        $actor = $petsList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::Pet();
            $actor->{appear_time} = time;
            $actor->{name} = $name;
            $actor->{name_given} = defined $name ? $name : T("Unknown");
            $mustAdd = 1;

            # Previously identified monsters could suddenly be identified as pets.
            if ($monstersList->getByID($args->{ID})) {
                $monstersList->removeByID($args->{ID});
            }
            $actor->{nameID} = $args->{type};
        }
    } elsif ($object_class eq 'Actor::Monster') {
        $actor = $monstersList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::Monster();
            $actor->{appear_time} = time;
            if ($monsters_lut{$args->{type}}) {
                $actor->setName($monsters_lut{$args->{type}});
            }
            $actor->{name} = $name if defined $name;
            $actor->{name_given} = "Unknown";
            $actor->{binType} = $args->{type};
            $mustAdd = 1;
            $actor->{nameID} = $args->{type};
        }
    } elsif ($object_class eq 'Actor::NPC') {
        # Actor is an NPC
        $actor = $npcsList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::NPC();
            $actor->{appear_time} = time;
            $actor->{name} = $name if defined $name;
            $mustAdd = 1;
        }
        $actor->{nameID} = $nameID;
    } elsif ($object_class eq 'Actor::Elemental') {
        # Actor is a Elemental
        $actor = $elementalsList->getByID($args->{ID});
        if (!defined $actor) {
            $actor = new Actor::Elemental();
            $actor->{appear_time} = time;
            $mustAdd = 1;
        }
        $actor->{name} = $jobs_lut{$args->{type}};
    }
```

## Position and Visibility Handling

```perl
    #### Step 2: update actor information ####
    $actor->{ID} = $args->{ID};
    $actor->{charID} = $args->{charID} if $args->{charID} && $args->{charID} ne "\0\0\0\0";
    $actor->{jobID} = $args->{type};
    $actor->{type} = $args->{type};
    $actor->{lv} = $args->{lv};
    $actor->{walk_speed} = $args->{walk_speed} / 1000 if (exists $args->{walk_speed} && $args->{switch} ne "0086");
    $actor->{len} = $args->{len} if $args->{len};
    $actor->{object_type} = $args->{object_type} if (defined $args->{object_type});

    # Remove actors that are located outside the map
    if (defined $field && ($field->isOffMap($coordsFrom{x}, $coordsFrom{y}) || $field->isOffMap($coordsTo{x}, $coordsTo{y}))) {
        warning TF("Ignoring actor with off map coordinates: (%d, %d)->(%d, %d), field max: (%d, %d)\n",$coordsFrom{x},$coordsFrom{y},$coordsTo{x},$coordsTo{y},$field->width(),$field->height());
        $actor->{avoid} = 1;
        return;
    }

    if ( ($coordsFrom{x} == 0 && $coordsFrom{y} == 0) || ($coordsTo{x} == 0 && $coordsTo{y} == 0) ||
         (blockDistance(\%coordsFrom, \%coordsTo) > $config{clientSight}) ) {
            warning TF("Ignoring bugged actor moved packet (%s) (%d, %d)->(%d, %d)\n", $args->{switch}, $coordsFrom{x}, $coordsFrom{y}, $coordsTo{x}, $coordsTo{y});
            return;
    }

    $actor->{pos} = {%coordsFrom};
    $actor->{pos_to} = {%coordsTo};
    $actor->{time_move} = time;
    $actor->{time_move_calc} = calcTime(\%coordsFrom, \%coordsTo, $actor->{walk_speed});

    # Ignore actors with a distance greater than clientSight
    if ($config{clientSight}) {
        my $realMyPos = calcPosition($char);
        my $realActorPos = calcPosition($actor);
        my $realActorDist = blockDistance($realMyPos, $realActorPos);

        if ($realActorDist >= $config{clientSight}) {
            my ($actor_type) = $object_class =~ /\:\:(\w+)$/;
            warning TF("Avoiding out of sight %s: '%s' at (%d, %d) (distance: %d >= max %d) - check clientSight in config.txt\n", $actor_type, $actor->{name}, $actor->{pos_to}{x}, $actor->{pos_to}{y}, $realActorDist, $config{clientSight});
            $actor->{avoid} = 1;
        } else {
            $actor->{avoid} = 0;
        }
    }

    #### Actor Appearance and Equipment ####
    if (UNIVERSAL::isa($actor, "Actor::Player")) {
        $actor->{emblemID} = $args->{emblemID} if (exists $args->{emblemID});
        $actor->{guildID} = $args->{guildID} if (exists $args->{guildID});

        if (exists $args->{lowhead}) {
            $actor->{headgear}{low} = $args->{lowhead};
            $actor->{headgear}{mid} = $args->{midhead};
            $actor->{headgear}{top} = $args->{tophead};
            $actor->{weapon} = $args->{weapon};
            $actor->{shield} = $args->{shield};
        }

        $actor->{sex} = $args->{sex};

        if ($args->{act} == 1) {
            $actor->{dead} = 1;
        } elsif ($args->{act} == 2) {
            $actor->{sitting} = 1;
        }

        $actor->{hair_color} = $args->{hair_color} if (exists $args->{hair_color});
    } elsif (UNIVERSAL::isa($actor, "Actor::NPC") && $args->{type} == 722) { # guild flag has emblem
        $actor->{emblemID} = $args->{emblemID};
        $actor->{guildID} = $args->{guildID};
    }

    $actor->{hair_style} = $args->{hair_style} if (exists $args->{hair_style});
    $actor->{look}{body} = $args->{body_dir} if (exists $args->{body_dir});
    $actor->{look}{head} = $args->{head_dir} if (exists $args->{head_dir});

    #### Visual Effects ####
    $actor->{opt3} = $args->{opt3} if (exists $args->{opt3}); # stackable

    # Known visual effects:
    # 0x0001 = Yellow tint (eg, a quicken skill)
    # 0x0002 = Red tint (eg, power-thrust)
    # 0x0004 = Gray tint (eg, energy coat)
    # 0x0008 = Slow lightning (eg, mental strength)
    # 0x0010 = Fast lightning (eg, MVP fury)
    # 0x0020 = Black non-moving statue (eg, stone curse)
    # 0x0040 = Translucent weapon
    # 0x0080 = Translucent red sprite (eg, marionette control?)
    # 0x0100 = Spaztastic weapon image (eg, mystical amplification)
    # 0x0200 = Gigantic glowy sphere-thing
    # 0x0400 = Translucent pink sprite (eg, marionette control?)
    # 0x0800 = Glowy sprite outline (eg, assumptio)
    # 0x1000 = Bright red sprite, slowly moving red lightning (eg, MVP fury?)
    # 0x2000 = Vortex-type effect

    $actor->{opt1} = $args->{opt1}; # nonstackable
    $actor->{opt2} = $args->{opt2}; # stackable
    $actor->{option} = $args->{option}; # stackable

    if (setStatus($actor, $args->{opt1}, $args->{opt2}, $args->{option})) {
        $mustAdd = 0;
    }

    #### Step 3: Add actor to actor list ####
    if ($mustAdd) {
        if (UNIVERSAL::isa($actor, "Actor::Player")) {
            $playersList->add($actor);
            Plugins::callHook('add_player_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::Monster")) {
            $monstersList->add($actor);
            Plugins::callHook('add_monster_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::Pet")) {
            $petsList->add($actor);
            Plugins::callHook('add_pet_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::Portal")) {
            $portalsList->add($actor);
            Plugins::callHook('add_portal_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::NPC")) {
            my $ID = $args->{ID};
            my $location = $field->baseName . " $actor->{pos}{x} $actor->{pos}{y}";
            if ($npcs_lut{$location}) {
                $actor->setName($npcs_lut{$location});
            }
            $npcsList->add($actor);
            Plugins::callHook('add_npc_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::Slave")) {
            $slavesList->add($actor);
            Plugins::callHook('add_slave_list', $actor);
        } elsif (UNIVERSAL::isa($actor, "Actor::Elemental")) {
            $elementalsList->add($actor);
            Plugins::callHook('add_elemental_list', $actor);
        }
    }

    #### Packet Specific Handling ####
    if ($args->{switch} eq "0078" || $args->{switch} eq "01D8" || $args->{switch} eq "022A" ||
        $args->{switch} eq "02EE" || $args->{switch} eq "07F9" || $args->{switch} eq "0915" ||
        $args->{switch} eq "09DD" || $args->{switch} eq "09FF" ||
        $packetParser->{packet_list}->{$args->{switch}}[0] eq "actor_exists") {
        # Actor Exists (standing)
        if ($actor->isa('Actor::Player')) {
            debug "Player Exists: " . $actor->name . " ($actor->{binID}) Level $actor->{lv} $sex_lut{$actor->{sex}} $jobs_lut{$actor->{jobID}} ($coordsFrom{x}, $coordsFrom{y})\n";
            Plugins::callHook('player_exist', {player => $actor});
        } elsif ($actor->isa('Actor::NPC')) {
            message TF("NPC Exists: %s (%d, %d) (ID %d) - (%d)\n", $actor->name, $actor->{pos_to}{x}, $actor->{pos_to}{y}, $actor->{nameID}, $actor->{binID});
            Plugins::callHook('npc_exist', {npc => $actor});
        } elsif ($actor->isa('Actor::Monster')) {
            debug sprintf("Monster Exists: %s (%d)\n", $actor->name, $actor->{binID});
            Plugins::callHook('monster_exist', {monster => $actor});
        } elsif ($actor->isa('Actor::Pet')) {
            debug sprintf("Pet Exists: %s (%d)\n", $actor->name, $actor->{binID});
            Plugins::callHook('pet_exist', {pet => $actor});
        } elsif ($actor->isa('Actor::Slave')) {
            debug sprintf("Slave Exists: %s (%d)\n", $actor->name, $actor->{binID});
            Plugins::callHook('slave_exist', {slave => $actor});
        } elsif ($actor->isa('Actor::Elemental')) {
            debug sprintf("Elemental Exists: %s (%d)\n", $actor->name, $actor->{binID});
            Plugins::callHook('elemental_exist', {elemental => $actor});
        }
    } elsif ($args->{switch} eq "0079" || $args->{switch} eq "01DB" || $args->{switch} eq "022B" ||
             $args->{switch} eq "02ED" || $args->{switch} eq "01D9" || $args->{switch} eq "07F8" ||
             $args->{switch} eq "0858" || $args->{switch} eq "090F" || $args->{switch} eq "09DC" ||
             $args->{switch} eq "09FE" || $packetParser->{packet_list}->{$args->{switch}}[0] eq "actor_connected") {
        # Actor Connected (new)
        if ($actor->isa('Actor::Player')) {
            debug "Player Connected: ".$actor->name." ($actor->{binID}) Level $args->{lv} $sex_lut{$actor->{sex}} $jobs_lut{$actor->{jobID}} ($coordsTo{x}, $coordsTo{y})\n";
            Plugins::callHook('player_connected', {player => $actor});
        }
    } elsif ($args->{switch} eq "007B" || $args->{switch} eq "0086" || $args->{switch} eq "01DA" ||
             $args->{switch} eq "022C" || $args->{switch} eq "02EC" || $args->{switch} eq "07F7" ||
             $args->{switch} eq "0856" || $args->{switch} eq "0914" || $args->{switch} eq "09DB" ||
             $args->{switch} eq "09FD" || $packetParser->{packet_list}->{$args->{switch}}[0] eq "actor_moved") {
        # Actor Moved
        if ($actor->isa('Actor::Player')) {
            debug "Player Moved: " . $actor->name . " ($actor->{binID}) Level $actor->{lv} $sex_lut{$actor->{sex}} $jobs_lut{$actor->{jobID}} - ($coordsFrom{x}, $coordsFrom{y}) -> ($coordsTo{x}, $coordsTo{y})\n";
            Plugins::callHook('player_moved', $actor);
        } elsif ($actor->isa('Actor::Monster')) {
            debug "Monster Moved: " . $actor->nameIdx . " - ($coordsFrom{x}, $coordsFrom{y}) -> ($coordsTo{x}, $coordsTo{y})\n";
            Plugins::callHook('monster_moved', $actor);
        } elsif ($actor->isa('Actor::Pet')) {
            debug "Pet Moved: " . $actor->nameIdx . " - ($coordsFrom{x}, $coordsFrom{y}) -> ($coordsTo{x}, $coordsTo{y})\n";
            Plugins::callHook('pet_moved', $actor);
        } elsif ($actor->isa('Actor::Slave')) {
            debug "Slave Moved: " . $actor->nameIdx . " - ($coordsFrom{x}, $coordsFrom{y}) -> ($coordsTo{x}, $coordsTo{y})\n";
            Plugins::callHook('slave_moved', $actor);
        }
    } elsif ($args->{switch} eq "007C") {
        # Actor Spawned
        if ($actor->isa('Actor::Player')) {
            debug "Player Spawned: " . $actor->nameIdx . " $sex_lut{$actor->{sex}} $jobs_lut{$actor->{jobID}}\n";
            Plugins::callHook('player_spawned', {player => $actor});
        } elsif ($actor->isa('Actor::Monster')) {
            debug "Monster Spawned: " . $actor->nameIdx . "\n";
            Plugins::callHook('monster_spawned', {monster => $actor});
        } elsif ($actor->isa('Actor::Pet')) {
            debug "Pet Spawned: " . $actor->nameIdx . "\n";
            Plugins::callHook('pet_spawned', {pet => $actor});
        } elsif ($actor->isa('Actor::Slave')) {
            debug "Slave Spawned: " . $actor->nameIdx . " $jobs_lut{$actor->{jobID}}\n";
            Plugins::callHook('slave_spawned', {slave => $actor});
        }
    }

    if ($char->{elemental}{ID} eq $actor->{ID}) {
        $char->{elemental} = $actor;
    }
}

#### Actor Death/Disappearance Handling (lines 2401-2586)
```perl
sub actor_died_or_disappeared {
	my ($self,$args) = @_;
	return unless changeToInGameState();
	my $ID = $args->{ID};
	avoidList_ID($ID);

	# type:
	#     0 = out of sight
	#     1 = died
	#     2 = logged out
	#     3 = teleport
	#     4 = trickdead

	if ($ID eq $accountID) {
		# Player death handling
		message T("You have died\n") if (!$char->{dead});
		$char->{deathCount}++;
		$char->{dead} = 1;
		$char->{dead_time} = time;
	} elsif (defined $monstersList->getByID($ID)) {
		# Monster death/disappearance
		my $monster = $monstersList->getByID($ID);
		if ($args->{type} == 1) {
			debug "Monster Died: " . $monster->name . " ($monster->{binID})\n", "parseMsg_damage";
			$monster->{dead} = 1;
		} elsif ($args->{type} == 0) {
			debug "Monster Disappeared: " . $monster->name . " ($monster->{binID})\n", "parseMsg_presence";
			$monster->{disappeared} = 1;
		}
		$monstersList->remove($monster);
	} elsif (defined $playersList->getByID($ID)) {
		# Player disappearance
		my $player = $playersList->getByID($ID);
		if ($args->{type} == 1) {
			message TF("Player Died: %s (%d) %s %s\n", $player->name, $player->{binID}, $sex_lut{$player->{sex}}, $jobs_lut{$player->{jobID}});
			$player->{dead} = 1;
		} else {
			debug "Player Disappeared: " . $player->name . " ($player->{binID})\n", "parseMsg_presence";
			$player->{disappeared} = 1;
			$playersList->remove($player);
		}
	} elsif (defined $portalsList->getByID($ID)) {
		# Portal disappearance
		my $portal = $portalsList->getByID($ID);
		debug "Portal Disappeared: " . $portal->name . " ($portal->{binID})\n", "parseMsg";
		$portalsList->remove($portal);
	}
	# Similar handlers exist for NPCs, pets, slaves, elementals
}
```

#### Actor Actions (lines 2588-2700)
```perl
sub actor_action {
	my ($self,$args) = @_;
	return unless changeToInGameState();

	if ($args->{type} == ACTION_ITEMPICKUP) {
		# Item pickup action
		my $source = Actor::get($args->{sourceID});
		debug "$source picks up item\n", 'parseMsg_presence';
	} elsif ($args->{type} == ACTION_SIT) {
		# Sitting action
		if ($args->{sourceID} eq $accountID) {
			message T("You are sitting.\n");
			$char->{sitting} = 1;
		} else {
			message TF("%s is sitting.\n", getActorName($args->{sourceID}));
		}
	} elsif ($args->{type} == ACTION_STAND) {
		# Standing action
		if ($args->{sourceID} eq $accountID) {
			message T("You are standing.\n");
			$char->{sitting} = 0;
		} else {
			message TF("%s is standing.\n", getActorName($args->{sourceID}));
		}
	} else {
		# Attack action
		my $totalDamage = $args->{damage} + $args->{dual_wield_damage};
		my $dmgdisplay = $totalDamage == 0 ? "Miss!" : $args->{damage};
		
		updateDamageTables($args->{sourceID}, $args->{targetID}, $totalDamage);
		
		my $source = Actor::get($args->{sourceID});
		my $target = Actor::get($args->{targetID});
		my $msg = attack_string($source, $target, $dmgdisplay);
		
		if ($args->{sourceID} eq $accountID) {
			message("$msg", $totalDamage > 0 ? "attackMon" : "attackMonMiss");
		} elsif ($args->{targetID} eq $accountID) {
			message("$msg", $args->{damage} > 0 ? "attacked" : "attackedMiss");
		}
	}
}

#### Actor Information Updates (lines 2712-2804)
```perl
sub actor_info {
	my ($self, $args) = @_;
	my $name = bytesToString($args->{name});
	$name =~ s/^\s+|\s+$//g;
	
	# Handle player info updates
	my $player = $playersList->getByID($args->{ID});
	if ($player) {
		$player->setName($name);
		$player->{party}{name} = bytesToString($args->{partyName}) if defined $args->{partyName};
		$player->{guild}{name} = bytesToString($args->{guildName}) if defined $args->{guildName};
		$player->{guild}{title} = bytesToString($args->{guildTitle}) if defined $args->{guildTitle};
		message "Player Info: " . $player->nameIdx . "\n", "parseMsg_presence", 2;
	}

	# Handle monster info updates
	my $monster = $monstersList->getByID($args->{ID});
	if ($monster) {
		$monster->{name_given} = $name;
		if ($monsters_lut{$monster->{nameID}} eq "") {
			$monster->setName($name);
			$monsters_lut{$monster->{nameID}} = $name;
		}
	}

	# Similar handlers exist for NPCs, pets, slaves, elementals
}
```

#### Level Up Effects (lines 2819-2844)
```perl
sub unit_levelup {
	my ($self, $args) = @_;
	my $ID = $args->{ID};
	my $type = $args->{type};
	my $actor = Actor::get($ID);

	if ($type == LEVELUP_EFFECT || $type == LEVELUP_EFFECT2 || $type == LEVELUP_EFFECT3) {
		message TF("%s gained a level!\n", $actor);
	} elsif ($type == JOBLEVELUP_EFFECT || $type == JOBLEVELUP_EFFECT2) {
		message TF("%s gained a job level!\n", $actor);
	} elsif ($type == REFINING_FAIL_EFFECT) {
		message TF("%s failed to refine a weapon!\n", $actor), "refine";
	} elsif ($type == REFINING_SUCCESS_EFFECT) {
		message TF("%s successfully refined a weapon!\n", $actor), "refine";
	} elsif ($type == MAKEITEM_AM_SUCCESS_EFFECT) {
		message TF("%s successfully created a potion!\n", $actor), "refine";
	} elsif ($type == MAKEITEM_AM_FAIL_EFFECT) {
		message TF("%s failed to create a potion!\n", $actor), "refine";
	} elsif ($type == GAME_OVER_EFFECT) {
		message TF("%s received GAME OVER!\n", $actor);
	}
}
```

#### Homunculus Properties (lines 2897-3000)
```perl
sub homunculus_property {
	my ($self, $args) = @_;
	return 0 unless enforce_homun_state();

	my $slave = $char->{homunculus};
	$slave->{name} = bytesToString($args->{name});
	
	# Update homunculus stats from packet
	foreach (@{$args->{KEYS}}) {
		$slave->{$_} = $args->{$_};
	}

	# Handle homunculus states:
	# 0 - alive, unnamed and not vaporized
	# 1 - named
	# 2 - rest
	# 4 - dead
	if ($args->{state} & 1) {
		$char->{homunculus_info}{renameflag} = 1;
	}
	if ($args->{state} & 2) {
		$char->{homunculus_info}{vaporized} = 1;
	}
	if ($args->{state} & 4) {
		$char->{homunculus_info}{dead} = 0;
	}
}
```

### Key Features:
- Handles all actor types (players, monsters, NPCs, pets, portals, etc.)
- Manages actor creation and lifecycle
- Processes position and movement data
- Implements visibility/distance checks
- Supports plugin hooks for actor events
- Handles coordinate transformations
- Manages actor lists (playersList, monstersList, etc.)
- Provides debug logging and warnings
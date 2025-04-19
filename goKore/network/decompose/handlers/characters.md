# Character Movement Handlers:
- character_moves (lines 5494-5531)
  - Processes player movement notifications
  - Key features:
    - Updates character position (pos_to)
    - Calculates movement distance
    - Sets movement timing (time_move)
    - Computes movement solution path
    - Calculates movement duration based on walk speed
    - Adjusts character facing direction
  - Special cases:
    - Handles route escape when no portals available
    - Triggers escape shout if configured
    - Queues escape AI action

# Character Information Handlers

**Sprite Change Handler:**
- sprite_change (lines 4528-4574)
  - Handles visual character changes:
    - Type 0: Job change
    - Type 2: Weapon/Shield change
    - Type 3: Lower headgear
    - Type 4: Upper headgear
    - Type 5: Middle headgear
    - Type 6: Hair color
    - Type 9: Shoes
    - Type 12: Robe
    - Types 7/13: Body palette
  - Shows appropriate change messages
  - Triggers sprite_job_change hook

**Progress Bar Handlers:**
- progress_bar (lines 4576-4588)
  - Starts progress bar with specified duration
  - Sends completion packet when done
  - Sets char->{progress_bar} flag
  - Shows start/complete messages

- progress_bar_stop (lines 4590-4593)
  - Immediately stops progress bar
  - Clears progress_bar flag
  - Shows completion message

**Marriage System Handler:**
- marriage_partner_name (lines 4112-4116)
  - Displays marriage partner's name
  - Shows before "I miss you" skill cast
  - Processes ZC_MARRIAGE_PARTNER_NAME packets
  - Output format: "Marriage partner name: $name"

**Status Change Handler:**
- character_status_change (lines 3656-3689)
  - Processes character status effect changes
  - Handles ZC_STATUS_CHANGE packets
  - Key features:
    - Updates status effect timers (start/end times)
    - Shows activation/removal messages:
      - "Status Effect Added: $effectName"
      - "Status Effect Removed: $effectName"
    - Maintains status effect list in $char->{statuses}
    - Supports both positive and negative effects
    - Triggers plugin hooks:
      - 'packet_statusChange' for all changes
      - 'packet_statusChange_$effectID' for specific effects

## Character Block Processing

### received_characters_blockSize() - lines 700-708
```perl
sub received_characters_blockSize {
    if ($masterServer && $masterServer->{charBlockSize}) {
        return $masterServer->{charBlockSize};
    } else {
        # last change: 2020-11-13
        # default in kRO, most of official servers and emulators (rAthena, Hercules)
        return 155;
    }
}
```

### received_characters_unpackString() - lines 711-809
```perl
sub received_characters_unpackString {
    my $char_info;
    for ($masterServer && $masterServer->{charBlockSize}) {
        if ($_ == 175) {  # PACKETVER >= 20201007 [hp, hp_max, sp and sp_max are now uint64]
            $char_info = {
                types => 'a4 V2 V V2 V6 v V2 V2 V2 V2 v2 V v9 Z24 C8 v Z16 V4 C',
                keys => [qw(charID exp exp_2 zeny exp_job exp_job_2 lv_job body_state health_state effect_state stance manner status_point hp hp_2 hp_max hp_max_2 sp sp_2 sp_max sp_max_2 walkspeed jobID hair_style weapon lv skill_point head_bottom shield head_top head_mid hair_pallete clothes_color name str agi vit int dex luk slot hair_color is_renamed last_map delete_date robe slot_addon rename_addon sex)],
            };
        } elsif ($_ == 155) {  # PACKETVER >= 20170830 [base and job exp are now uint64]
            $char_info = {
                types => 'a4 V2 V V2 V6 v V2 v4 V v9 Z24 C8 v Z16 V4 C',
                keys => [qw(charID exp exp_2 zeny exp_job exp_job_2 lv_job body_state health_state effect_state stance manner status_point hp hp_max sp sp_max walkspeed jobID hair_style weapon lv skill_point head_bottom shield head_top head_mid hair_pallete clothes_color name str agi vit int dex luk slot hair_color is_renamed last_map delete_date robe slot_addon rename_addon sex)],
            };
        } elsif ($_ == 147) { # PACKETVER >= 20141022 [iRO Doram Update, walk_speed is now long]
            $char_info = {
                types => 'a4 V9 v V2 v4 V v9 Z24 C8 v Z16 V4 C',
                keys => [qw(charID exp zeny exp_job lv_job body_state health_state effect_state stance manner status_point hp hp_max sp sp_max walkspeed jobID hair_style weapon lv skill_point head_bottom shield head_top head_mid hair_pallete clothes_color name str agi vit int dex luk slot hair_color is_renamed last_map delete_date robe slot_addon rename_addon sex)],
            };
        } else {
            die "Unknown charBlockSize: $_";
        }
        return $char_info;
    }
    die "masterserver or charBlockSize is undefined";
}
```

### Key Features:
- Handles multiple character block formats (3 versions shown)
- Provides unpack patterns for different protocol versions
- Includes comprehensive field mappings
- Has fallback to default block size

## Character Slot Management

### received_characters_slots_info() - lines 811-836
```perl
sub received_characters_slots_info {
    return if ($net->getState() == Network::IN_GAME);
    my ($self, $args) = @_;
    $net->setState(Network::CONNECTED_TO_LOGIN_SERVER);
    $charSvrSet{total_slot} = $args->{total_slot} if (exists $args->{total_slot});
    $charSvrSet{premium_start_slot} = $args->{premium_start_slot} if (exists $args->{premium_start_slot});
    $charSvrSet{premium_end_slot} = $args->{premium_end_slot} if (exists $args->{premium_end_slot});

    $charSvrSet{normal_slot} = $args->{normal_slot} if (exists $args->{normal_slot});
    $charSvrSet{premium_slot} = $args->{premium_slot} if (exists $args->{premium_slot});
    $charSvrSet{billing_slot} = $args->{billing_slot} if (exists $args->{billing_slot});

    $charSvrSet{producible_slot} = $args->{producible_slot} if (exists $args->{producible_slot});
    $charSvrSet{valid_slot} = $args->{valid_slot} if (exists $args->{valid_slot});

    undef $conState_tries;

    Plugins::callHook('parseMsg/recvChars', $args->{options});
    if ($args->{options} && exists $args->{options}{charServer}) {
        $charServer = $args->{options}{charServer};
    } else {
        $charServer = $net->serverPeerHost . ":" . $net->serverPeerPort;
    }

    $self->received_characters($args) if ($args->{charInfo});
}
```

### Key Features:
- Manages character slot information
- Handles premium/normal slot differentiation
- Processes server-provided slot data
- Maintains connection state
- Triggers character data processing

## Character Creation Status

### received_char_create_status() - lines 1001-1021
```perl
sub received_char_create_status {
    my ($self, $args) = @_;
    if ($args->{flag} == 0x00) {
        message T("Charname already exists.\n"), "info";
    } elsif ($args->{flag} == 0xFF) {
        message T("Char creation denied.\n"), "info";
    } elsif ($args->{flag} == 0x01) {
        message T("You are underaged.\n"), "info";
    } elsif ($args->{flag} == 0x02) {
        message T("Symbols in Character Names are forbidden.\n"), "info";
    } elsif ($args->{flag} == 0x03) {
        message T("You are not elegible to open the Character Slot.\n"), "info";
    } else {
        message T("Character creation failed. " .
            "If you didn't make any mistake, then the name you chose already exists.\n"), "info";
    }
    if (charSelectScreen() == 1) {
        $net->setState(3);
        $firstLoginMap = 1;
        $startingzeny = $chars[$config{'char'}]{'zeny'} unless defined $startingzeny;
        $sentWelcomeMessage = 1;
    }
}
```

### Key Features:
- Handles character creation status codes
- Provides user feedback for various failure cases
- Manages state transitions
- Handles initial character selection

## Character Stat Management

### stat_info_handlers - lines 1330-1500
```perl
our %stat_info_handlers = (
    VAR_SPEED, sub { $_[0]{walk_speed} = $_[1] / 1000 },
    VAR_EXP, sub {
        my ($actor, $value) = @_;
        $actor->{exp_last} = $actor->{exp};
        $actor->{exp} = $value;
        # Handles base exp calculations and level changes
    },
    VAR_JOBEXP, sub {
        my ($actor, $value) = @_;
        $actor->{exp_job_last} = $actor->{exp_job};
        $actor->{exp_job} = $value;
        # Handles job exp calculations and level changes
    },
    VAR_HP, sub { $_[0]{hp} = $_[1] },
    VAR_MAXHP, sub { $_[0]{hp_max} = $_[1] },
    VAR_SP, sub { $_[0]{sp} = $_[1] },
    VAR_MAXSP, sub { $_[0]{sp_max} = $_[1] },
    VAR_POINT, sub { $_[0]{points_free} = $_[1] },
    VAR_CLEVEL, sub {
        my ($actor, $value) = @_;
        $actor->{lv} = $value;
        message sprintf($actor->verb(T("%s are now level %d\n"), T("%s is now level %d\n")), $actor, $value), "success", $actor->isa('Actor::You') ? 1 : 2;
    },
    VAR_SPPOINT, sub { $_[0]{points_skill} = $_[1] },
    VAR_MONEY, sub {
        my ($actor, $value) = @_;
        my $change = $value - $actor->{zeny};
        $actor->{zeny} = $value;
        # Handles zeny change notifications
    }
);
```

### Key Features:
- Centralized stat management system
- Handles all core character attributes
- Manages experience calculations
- Tracks HP/SP changes
- Processes level up notifications
- Handles zeny (currency) changes

## Advanced Stat Management

### stat_info() - lines 1608-1668
```perl
sub stat_info {
    my ($self, $args) = @_;
    return unless changeToInGameState;

    my $actor = {
        '00B0' => $char,
        '00B1' => $char,
        '00BE' => $char,
        '0141' => $char,
        '01AB' => exists $args->{ID} && Actor::get($args->{ID}),
        '07DB' => $char->{homunculus},
        '0ACB' => $char,
    }->{$args->{switch}};

    if ($args->{switch} eq "02A2") {
        if (!$char->{mercenary}) {
            $char->{mercenary} = new Actor::Slave::Mercenary;
        }
        $actor = $char->{mercenary};
    }

    unless ($actor) {
        warning sprintf "Actor is unknown or not ready for stat information (switch %s, type %d, val %d)\n", @{$args}{qw(switch type val)};
        return;
    }

    if (exists $stat_info_handlers{$args->{type}}) {
        debug "Stat: $args->{type} => $args->{val}\n", "parseMsg", $_[0]->isa('Actor::You') ? 1 : 2;
        $stat_info_handlers{$args->{type}}($actor, $args->{val});
    } else {
        warning sprintf "Unknown stat (%d => %d) received for %s\n", @{$args}{qw(type val)}, $actor;
    }

    if (!$char->{walk_speed}) {
        $char->{walk_speed} = 0.15; # Default speed
    }
}
```

### stats_added() - lines 1670-1738
```perl
sub stats_added {
    my ($self, $args) = @_;

    if ($args->{val} == 207) {
        error T("Not enough stat points to add\n");
    } else {
        if ($args->{type} == VAR_STR) {
            $char->{str} = $args->{val};
            debug "Strength: $args->{val}\n", "parseMsg";
        } elsif ($args->{type} == VAR_AGI) {
            $char->{agi} = $args->{val};
            debug "Agility: $args->{val}\n", "parseMsg";
        } elsif ($args->{type} == VAR_VIT) {
            $char->{vit} = $args->{val};
            debug "Vitality: $args->{val}\n", "parseMsg";
        } elsif ($args->{type} == VAR_INT) {
            $char->{int} = $args->{val};
            debug "Intelligence: $args->{val}\n", "parseMsg";
        } elsif ($args->{type} == VAR_DEX) {
            $char->{dex} = $args->{val};
            debug "Dexterity: $args->{val}\n", "parseMsg";
        } elsif ($args->{type} == VAR_LUK) {
            $char->{luk} = $args->{val};
            debug "Luck: $args->{val}\n", "parseMsg";
        }
    }
    Plugins::callHook('packet_charStats', {
        type => $args->{type},
        val => $args->{val}
    });
}
```

### stats_info() - lines 1740-1790
```perl
sub stats_info {
    my ($self, $args) = @_;
    return unless changeToInGameState();
    $char->{points_free} = $args->{points_free};
    $char->{str} = $args->{str};
    $char->{points_str} = $args->{points_str};
    $char->{agi} = $args->{agi};
    $char->{points_agi} = $args->{points_agi};
    $char->{vit} = $args->{vit};
    $char->{points_vit} = $args->{points_vit};
    $char->{int} = $args->{int};
    $char->{points_int} = $args->{points_int};
    $char->{dex} = $args->{dex};
    $char->{points_dex} = $args->{points_dex};
    $char->{luk} = $args->{luk};
    $char->{points_luk} = $args->{points_luk};
    $char->{attack} = $args->{attack};
    $char->{attack_bonus} = $args->{attack_bonus};
    $char->{attack_magic_min} = $args->{attack_magic_min};
    $char->{attack_magic_max} = $args->{attack_magic_max};
    $char->{def} = $args->{def};
    $char->{def_bonus} = $args->{def_bonus};
    $char->{def_magic} = $args->{def_magic};
    $char->{def_magic_bonus} = $args->{def_magic_bonus};
    $char->{hit} = $args->{hit};
    $char->{flee} = $args->{flee};
    $char->{flee_bonus} = $args->{flee_bonus};
    $char->{critical} = $args->{critical};
    debug "Strength: $char->{str} #$char->{points_str}\n"
        ."Agility: $char->{agi} #$char->{points_agi}\n"
        ."Vitality: $char->{vit} #$char->{points_vit}\n"
        ."Intelligence: $char->{int} #$char->{points_int}\n"
        ."Dexterity: $char->{dex} #$char->{points_dex}\n"
        ."Luck: $char->{luk} #$char->{points_luk}\n"
        ."Status Points: $char->{points_free}\n", "parseMsg";
}
```

### Key Features:
- Handles stat point allocation results
- Updates core character attributes
- Provides detailed debug logging
- Manages combat attributes (attack, defense, etc.)
- Supports plugin hooks for stat changes
- Handles stat point bonuses
- Processes character parameter changes
- Manages status effect updates

## Character Data Processing

### received_characters() - lines 838-900
```perl
sub received_characters {
    my ($self, $args) = @_;
    my $blockSize = $self->received_characters_blockSize();
    my $char_info = $self->received_characters_unpackString;

    # rAthena and Hercules send all pages
    # Official Server send only pages with characters + 1 empty (tested bRO, iRO) Jul-2020
    if (length($args->{charInfo} == 0) {
        $charSvrSet{sync_received_characters} = $charSvrSet{sync_Count} if (exists $charSvrSet{sync_received_characters});
    } else {
        $charSvrSet{sync_received_characters}++ if (exists $charSvrSet{sync_received_characters});
    }

    $net->setState(Network::CONNECTED_TO_LOGIN_SERVER) if $net->getState() != Network::CONNECTED_TO_LOGIN_SERVER;

    return unless exists $args->{charInfo};

    for (my $i = 0; $i < length($args->{charInfo}); $i += $masterServer->{charBlockSize}) {
        my $temporary_character;
        @{$temporary_character}{@{$char_info->{keys}}} = unpack($char_info->{types}, substr($args->{charInfo}, $i, $masterServer->{charBlockSize}));

        my $character;

        # Re-use existing $char object instead of re-creating it
        if ($char && $char->{ID} eq $accountID && $char->{charID} eq $temporary_character->{charID}) {
            $character = $char;
        } elsif (exists $chars[$temporary_character->{slot}] && $chars[$temporary_character->{slot}]->{charID} eq $temporary_character->{charID}) {
            $character = $chars[$temporary_character->{slot}];
        } else {
            $character = new Actor::You;
        }

        @{$character}{@{$char_info->{keys}}} = unpack($char_info->{types}, substr($args->{charInfo}, $i, $masterServer->{charBlockSize}));
        $character->{ID} = $accountID;
        $character->{name} = bytesToString($character->{name});
        $character->{lastJobLvl} = $character->{lv_job};
        $character->{lastBaseLvl} = $character->{lv};
        $character->{headgear}{low} = $character->{head_bottom};
        $character->{headgear}{top} = $character->{head_top};
        $character->{headgear}{mid} = $character->{head_mid};
        $character->{nameID} = unpack("V", $character->{ID});
        $character->{last_map} =~ s/\.gat.*//g if ($character->{last_map});

        if ((!exists($character->{sex})) || ($character->{sex} ne "0" && $character->{sex} ne "1")) { $character->{sex} = $accountSex2; }

        $chars[$character->{slot}] = $character;
        setCharDeleteDate($character->{slot}, $character->{delete_date}) if $character->{delete_date};
    }

    message T("Received characters from Character Server\n"), "connection";
    $messageSender->sendBanCheck($accountID) if (!$net->clientAlive && $masterServer->{serverType} == 2);
}
```

### Key Features:
- Processes character data from server
- Maintains character state
- Handles character reuse and creation
- Converts data formats (bytes to string)
- Manages character deletion dates
- Provides user feedback
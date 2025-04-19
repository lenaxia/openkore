**Network Handlers:**

- ping() - Ping response handler (lines 12195-12199)
  - Responds to server ping requests
  * Only active in DirectConnection mode
  * Sends ping response back to server
- remain_time_info() - Logs remaining time information (lines 8017-8020)
  - Displays result, expiration date and remaining time
  - Outputs to console debug channel

- received_login_token() - Handles login token (lines 8022-8029)
  - Processes login token for XKore mode
  - Sends token to server with version info
  - Handles rathena server differences

- received_character_ID_and_Map() - Processes character/map info (lines 8056-8113)
  - Sets character ID and map connection info
  - Handles field initialization
  - Processes map IP/port from different formats
  - XKore 1 compatibility handling
  - Displays comprehensive game info
  - Performs allowed map check

- received_sync() - Handles sync packet (lines 8115-8119)
  - Updates play timeout
  - Minimal debug logging

# Network Utilities

## Sync Request Handler:
- sync_request_ex (lines 4413-4441)
  - Handles extended sync requests
  - Skips processing in XKore modes 1/3
  - Converts packet IDs:
    - Looks up sync ID from mapping table
    - Strips leading zeros
    - Converts to hex
  - Sends sync reply via messageSender
  - Used for client-server synchronization

## Packet Parser (parse) - lines 610-626
```perl
sub parse {
    my $self = shift;
    my $args = $self->SUPER::parse(@_);

    if ($args && $config{debugPacket_received} == 3 &&
            existsInList($config{'debugPacket_include'}, $args->{switch})) {
        my $packet = $self->{packet_list}{$args->{switch}};
        my ($name, $packString, $varNames) = @{$packet};

        my @vars = ();
        for my $varName (@{$varNames}) {
            message "$varName = $args->{$varName}\n";
        }
    }

    return $args;
}
```

### Key Features:
- Wraps SUPER::parse for basic packet parsing
- Adds debug logging when configured (debugPacket_received=3)
- Logs packet variables when packet is in debugPacket_include list
- Returns parsed packet arguments

## State Management (changeToInGameState) - lines 673-693
```perl
sub changeToInGameState {
    if ($net->version() == 1) {
        if ($accountID && UNIVERSAL::isa($char, 'Actor::You')) {
            if ($net->getState() != Network::IN_GAME) {
                $net->setState(Network::IN_GAME);
            }
            return 1;
        } else {
            if ($net->getState() != Network::IN_GAME_BUT_UNINITIALIZED) {
                $net->setState(Network::IN_GAME_BUT_UNINITIALIZED);
                if ($config{verbose} && $messageSender && !$sentWelcomeMessage) {
                    $messageSender->injectAdminMessage("Please relogin to enable X-${Settings::NAME}.");
                    $sentWelcomeMessage = 1;
                }
            }
            return 0;
        }
    } else {
        return 1;
    }
}
```

### Key Features:
- Manages network state transitions
- Handles both initialized and uninitialized game states
- Version-specific behavior
- Welcome message injection for uninitialized state

## Server Info Parsing (parse_account_server_info) - lines 1048-1107
```perl
sub parse_account_server_info {
    my ($self, $args) = @_;
    my $server_info;

    if ($args->{switch} eq '0B60') { # tRO 2020, twRO 2021
        $server_info = {
            len => 164,
            types => 'a4 v Z20 v3 a128 V',
            keys => [qw(ip port name state users property ip_port unknown)],
        };
    } elsif ($args->{switch} eq '0AC4' || $args->{switch} eq '0B07') { # kRO Zero 2017, kRO ST 201703+, vRO 2021
        $server_info = {
            len => 160,
            types => 'a4 v Z20 v3 a128',
            keys => [qw(ip port name users state property ip_port)],
        };
    } elsif ($args->{switch} eq '0AC9') { # cRO 2017
        $server_info = {
            len => 154,
            types => 'a20 V v a126',
            keys => [qw(name users unknown ip_port)],
        };
    } elsif ($args->{switch} eq '0276' && ($masterServer->{serverType} eq "tRO" or $masterServer->{serverType} eq "aRO")) {
        $server_info = {
            len => 36,
            types => 'a4 v Z20 v5',
            keys => [qw(ip port name state users property sid unknown)],
        };
    } else { # 0069 [default] and 0276 [pRO]
        $server_info = {
            len => 32,
            types => 'a4 v Z20 v3',
            keys => [qw(ip port name users display unknown)],
        };
    }

    @{$args->{servers}} = map {
        my %server;
        @server{@{$server_info->{keys}}} = unpack($server_info->{types}, $_);
        if ($masterServer && $masterServer->{private}) {
            $server{ip} = $masterServer->{ip};
        } elsif (exists $server{ip_port} && $server{ip_port} =~ /.*\:\d+/) {
            @server{qw(ip port)} = split (/\:/, $server{ip_port});
            $server{ip} =~ s/^\s+|\s+$//g;
            $server{port} =~ tr/0-9//cd;
        } else {
            $server{ip} = inet_ntoa($server{ip});
        }
        $server{name} = bytesToString($server{name});
        \%server
    } unpack '(a'.$server_info->{len}.')*', $args->{serverInfo};

    if (length $args->{lastLoginIP} == 4 && $args->{lastLoginIP} ne "\0"x4) {
        $args->{lastLoginIP} = inet_ntoa($args->{lastLoginIP});
    } else {
        delete $args->{lastLoginIP};
    }
}
```

### Key Features:
- Handles multiple server info formats
- Supports various regional server types
- Converts network byte order
- Manages IP/port parsing
- Handles private server configurations

## Map Loading (map_loaded) - lines 1237-1287
```perl
sub map_loaded {
    my ($self, $args) = @_;
    $net->setState(Network::IN_GAME);
    undef $conState_tries;
    $char = $chars[$config{char}];
    return unless changeToInGameState();
    main::initMapChangeVars();

    if ($net->version == 1) {
        $net->setState(4);
        message(T("Waiting for map to load...\n"), "connection");
        ai_clientSuspend(0, $timeout{'ai_clientSuspend'}{'timeout'});
    } else {
        $messageSender->sendReqRemainTime() if (grep { $masterServer->{serverType} eq $_ } qw(Zero Sakray));
        $messageSender->sendMapLoaded();
        $messageSender->sendSync(1);
        $messageSender->sendGuildRequestInfo(0) if ($masterServer->{serverType} ne 'twRO');
        $messageSender->sendRequestCashItemsList() if (grep { $masterServer->{serverType} eq $_ } qw(bRO idRO_Renewal twRO));
        $messageSender->sendCashShopOpen() if ($config{whenInGame_requestCashPoints});
        $messageSender->sendBlockingPlayerCancel() if $masterServer->{blockingPlayerCancel} || $self->{blockingPlayerCancel};
    }

    message(T("You are now in the game\n"), "connection");
    Plugins::callHook('in_game');
    $timeout{'ai'}{'time'} = time;
    our $quest_generation++;

    $char->{pos} = {};
    makeCoordsDir($char->{pos}, $args->{coords}, \$char->{look}{body});
    $char->{pos_to} = {%{$char->{pos}}};
    message(TF("Your Coordinates: %s, %s\n", $char->{pos}{x}, $char->{pos}{y}), undef, 1);
    $char->{time_move} = 0;
    $char->{time_move_calc} = 0;
    $char->{solution} = [];
    push(@{$char->{solution}}, { x => $char->{pos}{x}, y => $char->{pos}{y} });

    if ($masterServer->{private}){ setStatus($char, $char->{opt1}, $char->{opt2}, $char->{option}); }
    $ignored_all = 0;
}
```

### Key Features:
- Handles game state transition to IN_GAME
- Manages version-specific initialization
- Sends required post-login packets
- Tracks character position and movement
- Handles private server configurations
- Manages initial status effects
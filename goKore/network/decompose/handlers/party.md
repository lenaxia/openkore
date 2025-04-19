**Party Health & Status Handlers:**
- party_hp_info() - Updates party member HP (lines 8301-8309)
  - Updates HP and max HP for specified party member
  - Only updates if member exists in party

**Party Invitation Handlers:**
- party_invite() - Handles incoming party invites (lines 8311-8322)
  - Displays invite message
  - Stores invite details for response
  - Sets auto-deny timeout
  - Triggers party_invite hook

- party_invite_result() - Processes invite responses (lines 8324-8346)
  - Handles all response types:
    * Already in party
    * Request denied
    * Request accepted
    * Party full
    * Duplicate account
    * Unknown error
    * Character offline
    * Invalid map property

**Party Member Management:**
- party_leave() - Handles member leaving (lines 8348-8369)
  - Removes member from party list
  - Handles self-leave case
  - Differentiates between voluntary leave and kick
  - Cleans up party data structures

- party_location() - Updates member positions (lines 8371-8382)
  - Updates x,y coordinates for online members
  - Sets online status flag
  - Debug logs position changes

- party_organize_result() - Handles party creation (lines 8383-8397)
  - Sets admin flag for leader
  - Handles failure cases:
    * Name exists
    * Already in party
    * Map restriction
    * Unknown errors

**Party Item Handling:**
- party_show_picker() - Tracks item pickups (lines 8399-8413)
  - Displays who picked up what item
  - Formats item details
  - Excludes self-pickups

**Party Information Handlers:**
- party_users_info() - Updates party roster (lines 8415-8468)
  - Handles multiple packet versions
  - Processes member details:
    * ID, name, map
    * Admin status
    * Online status
    * Job, level (newer packets)
  - Maintains party user list
  - Triggers info_ready hook

- party_dead() - Tracks member deaths (lines 8470-8488)
  - Updates death status
  - Sets death timestamp
  - Handles revival notifications

**Party Management Handlers:**
- party_join() - Handles party join events (lines 8168-8220)
  - Processes different packet versions (0104, 01E9, 0A43, 0AE4)
  - Updates party member information (ID, role, position, etc.)
  - Handles both self-join and others joining cases
  - Maintains party user list and properties

- party_allow_invite() - Controls party invite permissions (lines 8222-8231)
  - Toggles whether others can invite to party
  - Shows appropriate status messages

- party_chat() - Processes party chat messages (lines 8233-8256)
  - Parses and logs party chat messages
  - Handles language codes and message solving
  - Supports chat logging if configured

- party_exp() - Manages party experience settings (lines 8258-8286)
  - Handles EXP sharing mode (Individual/Even Share)
  - Manages item pickup and division settings
  - Shows appropriate status messages for each setting

- party_leader() - Handles party leadership changes (lines 8288-8298)
  - Updates admin status for old/new leaders
  - Announces new leader to party

**Party Join Request Handler:**
- party_join_request_by_name (lines 3591-3620)
  - Processes party join requests by player name
  - Handles ZC_ADD_MEMBER_TO_GROUP2 packet
  - Key features:
    - Validates party size limits (max 12 members)
    - Shows request notification: "Party Request from $name"
    - Supports auto-accept via config{partyAutoAccept}
    - Maintains party request timeout (30 seconds)
    - Stores request info in $incomingParty{ID} and $incomingParty{name}
    - Triggers 'packet_partyJoinRequest' plugin hook

#### Party Creation (lines 3656-3678)
```perl
sub party_created {
	my ($self, $args) = @_;
	$char->{party}{ID} = $args->{partyID};
	message T("Party created\n"), "party";
	Plugins::callHook('packet_partyCreated');
}
```

#### Party Member Updates (lines 3680-3721)
```perl
sub party_member_update {
	my ($self, $args) = @_;
	return unless $char->{party}{ID};
	
	my $ID = $args->{ID};
	my $member = $party{members}{$ID};
	
	if ($member) {
		$member->{name} = bytesToString($args->{name});
		$member->{jobID} = $args->{jobID};
		$member->{lv} = $args->{lv};
		$member->{online} = $args->{online};
	} else {
		$party{members}{$ID} = {
			name => bytesToString($args->{name}),
			jobID => $args->{jobID},
			lv => $args->{lv},
			online => $args->{online},
			ID => $ID
		};
	}
	
	message TF("%s (%s) - %s\n", $party{members}{$ID}{name}, $jobs_lut{$args->{jobID}}, $args->{online} ? T("Online") : T("Offline")), "party";
}
```

#### Party Invitation (lines 3723-3745)
```perl
sub party_invite {
	my ($self, $args) = @_;
	my $name = bytesToString($args->{name});
	message TF("%s has invited you to join a party\n", $name), "party";
	
	if ($config{partyAuto} == 1) {
		$messageSender->sendPartyJoin($args->{partyID});
		message T("Auto-accepted party invitation\n"), "party";
	} elsif ($config{partyAuto} == 2) {
		$messageSender->sendPartyJoin($args->{partyID});
		message T("Auto-rejected party invitation\n"), "party";
	}
	
	Plugins::callHook('packet_partyInvite', {name => $name, partyID => $args->{partyID}});
}
```

### Key Features:
- Handles party creation and ID assignment
- Processes party member updates (online status, level, job)
- Maintains party member roster in memory
- Supports automatic party invitation responses (accept/reject)
- Provides detailed member information display
- Uses internationalization (T()/TF()) for localized messages
- Integrates with plugin system via hooks
- Supports configurable party auto-response settings
- Maintains data consistency for party operations
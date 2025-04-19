**Warp Portal Handler:**
- warp_portal_list (lines 3521-3559)
  - Processes warp portal/memo point listing
  - Handles up to 4 memo locations (memo1-memo4)
  - Key features:
    - Auto-detects saveMap from packet type
    - Strips .gat extension from map names
    - Formats display with map names and coordinates
    - Integrates with teleport skill system
    - Supports both ZC_WARPLIST and ZC_MEMORIALDUNGEON_SUBSCRIPTION_INFO packets
  - Output format:
    ```
    [Memo #] MapName (X,Y)
    ```
  - Stores locations in $char->{memo}{$num}
  - Triggers 'packet/memo' plugin hook
```perl
sub warp_portal_list {
	my ($self, $args) = @_;

	# strip gat extension
	($args->{memo1}) = $args->{memo1} =~ /^(.*)\.gat/;
	($args->{memo2}) = $args->{memo2} =~ /^(.*)\.gat/;
	($args->{memo3}) = $args->{memo3} =~ /^(.*)\.gat/;
	($args->{memo4}) = $args->{memo4} =~ /^(.*)\.gat/;
	
	# Auto-detect saveMap
	if ($args->{type} == 26) {
		configModify('saveMap', $args->{memo2}) if ($args->{memo2} && $config{'saveMap'} ne $args->{memo2});
	} elsif ($args->{type} == 27) {
		configModify('saveMap', $args->{memo1}) if ($args->{memo1} && $config{'saveMap'} ne $args->{memo1});
		configModify( "memo$_", $args->{"memo$_"} ) foreach grep { $args->{"memo$_"} ne $config{"memo$_"} } 1 .. 4;
	}

	$char->{warp}{type} = $args->{type};
	undef @{$char->{warp}{memo}};
	push @{$char->{warp}{memo}}, $args->{memo1} if $args->{memo1} ne "";
	push @{$char->{warp}{memo}}, $args->{memo2} if $args->{memo2} ne "";
	push @{$char->{warp}{memo}}, $args->{memo3} if $args->{memo3} ne "";
	push @{$char->{warp}{memo}}, $args->{memo4} if $args->{memo4} ne "";

	my $msg = center(T(" Warp Portal "), 50, '-') ."\n".
		T("#  Place                           Map\n");
	for (my $i = 0; $i < @{$char->{warp}{memo}}; $i++) {
		$msg .= swrite(
			"@< @<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< @<<<<<<<<<<<<<<<",
			[$i, $maps_lut{$char->{warp}{memo}[$i].'.rsw'}, $char->{warp}{memo}[$i]]);
	}
	$msg .= ('-'x50) . "\n";
	message $msg, "list";

	if ($args->{type} == 26 && AI::inQueue('teleport')) {
		$messageSender->sendWarpTele(26, AI::args->{lv} == 2 ? "$config{saveMap}.gat" : "Random");
		AI::dequeue;
	}
}
```

### Key Features:
- Manages warp portals and memo points
- Handles two types of portal lists (type 26 and 27)
- Automatically updates saveMap configuration
- Processes up to 4 memo points per portal
- Strips .gat extension from map names
- Displays formatted portal list to user
- Integrates with AI system for teleport handling
- Maintains portal state in character data
- Supports internationalization through T() function
- Provides clean user interface with aligned columns
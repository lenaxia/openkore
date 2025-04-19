#### Equipment Display (lines 3216-3281)
```perl
sub show_eq {
	my ($self, $args) = @_;
	
	# Handles multiple packet versions for equipment display:
	# - 02D7 (default)
	# - 0859 (20101124+)
	# - 0997 (20120925+)
	# - 0A2D (20150226+)
	# - 0B03 (20150226+)
	
	my $item_info;
	if ($args->{switch} eq '02D7') {
		$item_info = { len => 26, ... };
	} elsif ($args->{switch} eq '0859') {
		$item_info = { len => 28, ... };
	} elsif ($args->{switch} eq '0997') {
		$item_info = { len => 31, ... };
	} elsif ($args->{switch} eq '0A2D') {
		$item_info = { len => 57, ... };
	} elsif ($args->{switch} eq '0B03') {
		$item_info = { len => 67, ... };
	}

	# Format and display equipment info
	my $name = bytesToString($args->{name});
	my $msg = center(" $name " . T("Equip Info") . " ", 50, '-') . "\n";
	for (my $i = 0; $i < length($args->{equips_info}); $i += $item_info->{len}) {
		my $item = unpack_item($args->{equips_info}, $i, $item_info);
		$msg .= sprintf("%-20s: %s\n", $equipTypes_lut{$item->{equipped}}, itemName($item));
	}
	message($msg, "list");
}
```

### Key Features:
- Supports multiple packet versions with different item formats
- Handles equipment display for other players
- Formats detailed equipment information
- Maintains backward compatibility
- Processes item properties including:
  - Equipment slots
  - Item names and IDs
  - Upgrade levels
  - Cards and options
  - Binding status
- Displays in consistent tabular format
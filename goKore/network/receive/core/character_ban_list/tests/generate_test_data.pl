#!/usr/bin/perl
use strict;
use warnings;

# This script generates test data for the character_ban_list packet
# The packet format is: 'v a*' (unsigned short + remaining bytes)
# Field names: [len, charList]

# Create a mock packet
sub create_mock_packet {
    my ($num_entries) = @_;
    
    my $packet = '';
    
    # Add the number of entries (2 bytes, unsigned short)
    $packet .= pack('v', $num_entries);
    
    # Add character names (each 24 bytes)
    for my $i (1..$num_entries) {
        my $char_name = "TestChar$i";
        my $padded_name = $char_name . ("\0" x (24 - length($char_name)));
        $packet .= $padded_name;
    }
    
    return $packet;
}

# Output test data for Go tests
sub output_test_data {
    my ($num_entries) = @_;
    
    my $mock_packet = create_mock_packet($num_entries);
    
    print "Test data for $num_entries entries:\n";
    print "Hex: " . unpack('H*', $mock_packet) . "\n";
    
    my @bytes = unpack('C*', $mock_packet);
    print "Bytes: [" . join(', ', @bytes) . "]\n";
    
    if ($num_entries > 0) {
        print "Character names:\n";
        for my $i (0..$num_entries-1) {
            my $name = substr($mock_packet, 2 + $i * 24, 24);
            $name =~ s/\0+$//;  # Remove trailing null bytes
            print "  $i: '$name'\n";
        }
    }
    
    print "\n";
}

# Output test data for Go tests
print "=== Test Data for character_ban_list packet ===\n\n";
output_test_data(0);  # Empty list
output_test_data(1);  # One entry
output_test_data(2);  # Two entries
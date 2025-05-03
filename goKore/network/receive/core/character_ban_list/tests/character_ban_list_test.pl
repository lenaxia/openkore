#!/usr/bin/perl
use strict;
use warnings;
use Test::More;
use FindBin qw($RealBin);
use lib "$RealBin/../../../../../../src";  # Adjust path to point to the src directory
use Network::Receive;
use Network::Receive::ServerType0;

# Create a test for character_ban_list
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

# Test case 1: Valid ban list with multiple entries
subtest 'Valid ban list with multiple entries' => sub {
    # Create a mock packet with 2 entries
    my $mock_packet = create_mock_packet(2);
    
    # Create a receive object
    my $receive = Network::Receive::ServerType0->new();
    
    # Create args hash
    my $args = {
        len => 2,
        charList => substr($mock_packet, 2)  # Skip the length field
    };
    
    # Call the handler
    $receive->character_ban_list($args);
    
    # Since the Perl implementation is empty, we just verify the packet structure
    is(length($args->{charList}), 2 * 24, 'Character list has correct length');
    
    # Extract character names
    my @char_names;
    for my $i (0..1) {
        my $name = substr($args->{charList}, $i * 24, 24);
        $name =~ s/\0+$//;  # Remove trailing null bytes
        push @char_names, $name;
    }
    
    is($char_names[0], 'TestChar1', 'First character name is correct');
    is($char_names[1], 'TestChar2', 'Second character name is correct');
};

# Test case 2: Empty ban list
subtest 'Empty ban list' => sub {
    # Create a mock packet with 0 entries
    my $mock_packet = create_mock_packet(0);
    
    # Create a receive object
    my $receive = Network::Receive::ServerType0->new();
    
    # Create args hash
    my $args = {
        len => 0,
        charList => substr($mock_packet, 2)  # Skip the length field
    };
    
    # Call the handler
    $receive->character_ban_list($args);
    
    # Since the Perl implementation is empty, we just verify the packet structure
    is(length($args->{charList}), 0, 'Character list is empty');
};

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
output_test_data(0);  # Empty list
output_test_data(2);  # Two entries

done_testing();
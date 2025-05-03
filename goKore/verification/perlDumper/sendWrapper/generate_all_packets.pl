#!/usr/bin/env perl
# generate_all_packets.pl - Automatically generate test data for all Send.pm subroutines

use strict;
use warnings;
use FindBin qw($RealBin);
use lib "$RealBin/../../../../src";
use JSON::PP;
use Data::Dumper;
use File::Path qw(make_path);

# Load required modules
require "$RealBin/../../../../src/Network/MessageTokenizer.pm";
require "$RealBin/../../../../src/Network/PacketParser.pm";
require "$RealBin/../../../../src/Network/Send.pm";

# Mock necessary dependencies
package Globals;
our $accountID = pack("V", 0x12345678);
our $char = { 
    look => {},
    skills => { 1 => { ID => 1, lv => 10 } },
    inventory => [],
    cart => [],
    storage => [],
};
our $masterServer = { 
    serverType => 0,
    ip => '127.0.0.1',
    port => 6900,
    master_version => 1,
    version => 1,
    sendCryptKeys => '',
};
our $messageSender;
our $net;
our $packetParser;
our $bytesSent = 0;
our %packetDescriptions = (
    Send => {}
);
our $enc_val1 = 0;
our $enc_val2 = 0;
our $syncSync = pack("V", 12345);
our %timeout = ();
our %talk = ();
our $skillExchangeItem = 0;
our $rodexList = {};
our $rodexWrite = {};
our %universalCatalog = ();
our %rpackets = ();
our $mergeItemList = [];
our $repairList = [];
our %cashShop = ();

package Log;
sub debug { }
sub message { }
sub warning { }
sub error { }

package Misc;
sub visualDump { }
sub dumpData { }
sub stripLanguageCode { }

package Plugins;
sub callHook { }
sub hasHook { return 0; }

package Utils;
sub getTickCount { return time; }
sub getHex { return $_[0]; }
sub existsInList { return 0; }
sub getCoordString { return pack("C3", $_[0], $_[1], 0); }
sub makeCoordsDir { }

package I18N;
sub bytesToString { return $_[0]; }
sub stringToBytes { return $_[0]; }

package Settings;
our %sys = (locale => 'en');

# Create a mock connection that captures packets instead of sending them
package MockConnection;
sub new {
    my $class = shift;
    return bless { packets => [] }, $class;
}

sub getState { return 5; } # IN_GAME
sub serverSend {
    my ($self, $msg) = @_;
    push @{$self->{packets}}, $msg;
    return 1;
}
sub serverAlive { return 1; }
sub clientSend { 
    my ($self, $msg) = @_;
    push @{$self->{packets}}, $msg;
    return 1;
}
sub version { return 0; }

# Main program
package main;

# Create output directory
my $output_dir = "$RealBin/testdata";
make_path($output_dir) unless -d $output_dir;

print "Creating mock connection...\n";
my $connection = MockConnection->new();

print "Creating Send object...\n";
my $send = Network::Send->create($connection, 0); # ServerType 0
$Globals::messageSender = $send;
$Globals::net = $connection;

# Get all send* methods from Send.pm
print "Finding all send* methods...\n";
my @methods = grep { /^send[A-Z]/ } keys %{Network::Send::};
print "Found " . scalar(@methods) . " methods\n";

# Define argument generators for each method
my %arg_generators = (
    # Login-related methods
    sendMasterLogin => sub { return ("username", "password", 1, 1); },
    sendGameLogin => sub { return (123456, "sessionID", "sessionID2", 1); },
    sendCharLogin => sub { return (0); },
    sendMapLogin => sub { return (123456, 789012, "sessionID", 1); },
    
    # Movement and actions
    sendMove => sub { return (100, 100); },
    sendLook => sub { return (0, 0); },
    sendAction => sub { return (12345, 0); },
    
    # Chat
    sendChat => sub { return ("Hello, world!"); },
    sendPrivateMsg => sub { return ("Player", "Hello there!"); },
    sendPartyChat => sub { return ("Party message"); },
    sendGuildChat => sub { return ("Guild message"); },
    
    # Items
    sendItemUse => sub { return (1234, 5678); },
    sendDrop => sub { return (1234, 10); },
    sendTake => sub { return (1234); },
    sendStorageAdd => sub { return (1234, 10); },
    sendStorageGet => sub { return (1234, 10); },
    sendStoragePassword => sub { return ("password", 3); },
    
    # Skills
    sendSkillUse => sub { return (1, 10, 12345); },
    sendSkillUseLoc => sub { return (1, 10, 100, 100); },
    
    # Other common methods
    sendSync => sub { return (); },
    sendRestart => sub { return (0); },
    
    # NPC interaction
    sendTalk => sub { return (12345); },
    sendTalkCancel => sub { return (12345); },
    sendTalkContinue => sub { return (12345); },
    sendTalkResponse => sub { return (12345, 1); },
    sendTalkNumber => sub { return (12345, 100); },
    sendTalkText => sub { return (12345, "Hello NPC"); },
    
    # Party
    sendPartyLeader => sub { return (12345); },
    sendPartyOption => sub { return (1, 1, 1); },
    
    # Default generator for methods without specific generators
    DEFAULT => sub { return (); }
);

# Process each method
foreach my $method (sort @methods) {
    print "Processing $method...\n";
    
    # Skip methods that are not callable or are just utility functions
    next if $method =~ /^sendRaw$/;
    next if $method =~ /^sendToServer$/;
    
    # Clear previous packets
    $connection->{packets} = [];
    
    # Generate arguments
    my @args;
    if (exists $arg_generators{$method}) {
        @args = $arg_generators{$method}->();
    } else {
        print "  No specific argument generator for $method, using default\n";
        @args = $arg_generators{DEFAULT}->();
    }
    
    # Call the method
    print "  Calling $method with args: " . join(", ", map { defined $_ ? $_ : 'undef' } @args) . "\n";
    eval {
        $send->$method(@args);
    };
    
    if ($@) {
        print "  Error calling $method: $@\n";
        next;
    }
    
    # Get the captured packets
    my $packets = $connection->{packets};
    
    if (!@$packets) {
        print "  No packets generated for $method\n";
        next;
    }
    
    # Create output data
    my $output = {
        method => $method,
        args => \@args,
        packets => []
    };
    
    # Process each packet
    for my $i (0..$#{$packets}) {
        my $packet = $packets->[$i];
        my $hex = unpack("H*", $packet);
        
        # Try to identify the packet type
        my $messageID = uc(unpack("H2", substr($packet, 1, 1)) . unpack("H2", substr($packet, 0, 1)));
        
        # Convert to byte array
        my @bytes = map { ord($_) } split //, $packet;
        
        push @{$output->{packets}}, {
            hex => $hex,
            messageID => $messageID,
            bytes => \@bytes
        };
        
        print "  Packet $i: $messageID ($hex)\n";
    }
    
    # Save JSON output to a file
    my $json_file = "$output_dir/${method}.json";
    open my $fh, ">", $json_file or die "Cannot open $json_file: $!";
    print $fh encode_json($output);
    close $fh;
    print "  Output saved to $json_file\n";
}

print "\nAll done! Generated test data for " . scalar(@methods) . " methods\n";
print "Test data is available in $output_dir\n";
#!/usr/bin/env perl
# generate_from_analysis.pl - Generate test data using method analysis

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

# Load analysis results
my $analysis_file = "$RealBin/analysis/send_methods_analysis.json";
open my $fh, '<', $analysis_file or die "Cannot open $analysis_file: $!";
my $analysis_json = do { local $/; <$fh> };
close $fh;

my $analysis = decode_json($analysis_json);

print "Creating mock connection...\n";
my $connection = MockConnection->new();

print "Creating Send object...\n";
my $send = Network::Send->create($connection, 0); # ServerType 0
$Globals::messageSender = $send;
$Globals::net = $connection;

# Define default values for common argument types
my %default_values = (
    # Numeric values
    'ID' => 12345,
    'accountID' => 123456,
    'charID' => 789012,
    'sessionID' => pack("V", 54321),
    'sessionID2' => pack("V", 98765),
    'targetID' => 12345,
    'amount' => 10,
    'index' => 1,
    'slot' => 0,
    'type' => 1,
    'flag' => 0,
    'value' => 100,
    'x' => 100,
    'y' => 100,
    'lv' => 10,
    'skillID' => 1,
    
    # String values
    'username' => 'username',
    'password' => 'password',
    'message' => 'Hello, world!',
    'privMsg' => 'Hello there!',
    'privMsgUser' => 'Player',
    'storeName' => 'My Shop',
    'text' => 'Sample text',
    'code' => '123456',
    
    # Hash values
    'items' => [{ ID => 1, itemID => 1, amount => 1 }],
    'coords' => pack("C3", 100, 100, 0),
    
    # Special values
    'master_version' => 1,
    'version' => 1,
    'sex' => 1,
    'accountSex' => 1,
    'time' => time(),
    'tick' => time(),
);

# Generate arguments for each method
foreach my $method (sort keys %{$analysis->{parameters}}) {
    print "Processing $method...\n";
    
    # Skip methods that are not callable or are just utility functions
    next if $method =~ /^(sendRaw|sendToServer)$/;
    
    # Clear previous packets
    $connection->{packets} = [];
    
    my $params = $analysis->{parameters}{$method};
    my @args;
    
    # Skip $self parameter
    shift @$params;
    
    # Generate arguments based on parameter types
    foreach my $param (@$params) {
        if ($param =~ /\$args/) {
            # Handle $args parameter
            my $args_hash = {};
            
            if (exists $analysis->{args_keys}{$method}) {
                # Use analyzed keys
                foreach my $key (@{$analysis->{args_keys}{$method}}) {
                    $args_hash->{$key} = exists $default_values{$key} ? $default_values{$key} : 1;
                }
            } else {
                # No keys found, use empty hash
                $args_hash = {};
            }
            
            push @args, $args_hash;
        } elsif ($param =~ /\$(\w+)/) {
            # Handle regular parameter
            my $param_name = $1;
            push @args, exists $default_values{$param_name} ? $default_values{$param_name} : 1;
        } else {
            # Unknown parameter type
            push @args, 1;
        }
    }
    
    # Call the method
    print "  Calling $method with args: " . join(", ", map { defined $_ ? (ref $_ ? ref $_ : $_) : 'undef' } @args) . "\n";
    eval {
        $send->$method(@args);
    };
    
    if ($@) {
        print "  Error calling $method: $@\n";
        next;
    }
    
    # Get the captured packets
    my $packets = $connection->{packets};
    
    # Create output data
    my $output = {
        method => $method,
        args => \@args,
        packets => []
    };
    
    # Process each packet
    if (@$packets) {
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
    } else {
        print "  No packets generated for $method\n";
        
        # For methods that don't generate packets, save the method info anyway
        my $json_file = "$output_dir/${method}.json";
        open my $fh, ">", $json_file or die "Cannot open $json_file: $!";
        print $fh encode_json($output);
        close $fh;
        print "  Method info saved to $json_file\n";
    }
}

print "\nAll done! Generated test data based on method analysis\n";
print "Test data is available in $output_dir\n";

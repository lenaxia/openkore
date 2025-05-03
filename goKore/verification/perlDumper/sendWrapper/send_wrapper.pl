#!/usr/bin/env perl
# send_wrapper.pl - Wrapper for Send.pm to generate validation data

use strict;
use warnings;
use FindBin qw($RealBin);
use lib "$RealBin/../../../../src";  # Updated path to OpenKore src
use lib "$RealBin/../../../../";     # Updated path to OpenKore root
use JSON::PP;
use Data::Dumper;

# Load required modules
require "$RealBin/../../../../src/Network/MessageTokenizer.pm";  # Updated path
require "$RealBin/../../../../src/Network/PacketParser.pm";      # Updated path
require "$RealBin/../../../../src/Network/Send.pm";              # Updated path

# Mock necessary dependencies
package Globals;
our $accountID = pack("V", 0x12345678);
our $char = { look => {} };
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

# Parse command line arguments
my $command = $ARGV[0] || "help";

if ($command eq "help") {
    print "Usage: $0 [command] [args...]\n";
    print "Commands:\n";
    print "  list                  - List available Send.pm methods\n";
    print "  call [method] [args]  - Call a specific Send.pm method with arguments\n";
    print "  help                  - Show this help message\n";
    exit;
}

# Create mock connection
my $connection = MockConnection->new();

# Create Send object
my $send = Network::Send->create($connection, 0); # ServerType 0
$Globals::messageSender = $send;
$Globals::net = $connection;

if ($command eq "list") {
    # List available methods
    my @methods = grep { /^send[A-Z]/ } keys %{Network::Send::};
    print "Available Send.pm methods:\n";
    foreach my $m (sort @methods) {
        print "  $m\n";
    }
    exit;
} elsif ($command eq "call") {
    my $method = $ARGV[1] || die "No method specified. Use 'list' to see available methods.\n";
    my @args = @ARGV[2..$#ARGV];
    
    # Call the method
    if ($send->can($method)) {
        print "Calling $method with args: " . join(", ", @args) . "\n";
        $send->$method(@args);
        
        # Get the captured packets
        my $packets = $connection->{packets};
        
        # Output the packets in both hex and JSON format
        my $output = {
            method => $method,
            args => \@args,
            packets => []
        };
        
        print "Generated packets:\n";
        for my $i (0..$#{$packets}) {
            my $packet = $packets->[$i];
            my $hex = unpack("H*", $packet);
            print "Packet $i: $hex\n";
            
            # Try to identify the packet type
            my $messageID = uc(unpack("H2", substr($packet, 1, 1)) . unpack("H2", substr($packet, 0, 1)));
            print "  Message ID: $messageID\n";
            
            # Print raw bytes for debugging
            print "  Raw bytes: ";
            my @bytes = map { ord($_) } split //, $packet;
            for my $byte (@bytes) {
                printf("%02X ", $byte);
            }
            print "\n";
            
            push @{$output->{packets}}, {
                hex => $hex,
                messageID => $messageID,
                bytes => \@bytes
            };
        }
        
        # Save JSON output to a file
        my $json_file = "packet_output_${method}.json";
        open my $fh, ">", $json_file or die "Cannot open $json_file: $!";
        print $fh encode_json($output);
        close $fh;
        print "JSON output saved to $json_file\n";
    } else {
        die "Unknown method: $method. Use 'list' to see available methods.\n";
    }
} else {
    die "Unknown command: $command. Use 'help' for usage information.\n";
}
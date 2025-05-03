#!/usr/bin/env perl
# generate_all_packets_comprehensive.pl - Automatically generate test data for all Send.pm subroutines

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

# Get all methods from Send.pm
print "Finding all methods in Send.pm...\n";

# Get methods that start with "send"
my @send_methods = grep { /^send[A-Z]/ } keys %{Network::Send::};
print "Found " . scalar(@send_methods) . " send methods\n";

# Get methods that start with "reconstruct"
my @reconstruct_methods = grep { /^reconstruct_/ } keys %{Network::Send::};
print "Found " . scalar(@reconstruct_methods) . " reconstruct methods\n";

# Get methods that start with "parse"
my @parse_methods = grep { /^parse_/ } keys %{Network::Send::};
print "Found " . scalar(@parse_methods) . " parse methods\n";

# Get other methods
my @other_methods = grep { 
    !/^send[A-Z]/ && 
    !/^reconstruct_/ && 
    !/^parse_/ && 
    !/^(import|AUTOLOAD|DESTROY|BEGIN|END|UNITCHECK|CHECK|INIT)$/ &&
    /^[a-z]/  # Only include methods that start with lowercase letters
} keys %{Network::Send::};
print "Found " . scalar(@other_methods) . " other methods\n";

# Combine all methods
my @all_methods = (@send_methods, @reconstruct_methods, @parse_methods, @other_methods);
print "Total methods to process: " . scalar(@all_methods) . "\n";

# Define argument generators for each method type
my %arg_generators = (
    # Send methods
    sendMasterLogin => sub { return ("username", "password", 1, 1); },
    sendGameLogin => sub { return (123456, "sessionID", "sessionID2", 1); },
    sendCharLogin => sub { return (0); },
    sendMapLogin => sub { return (123456, 789012, "sessionID", 1); },
    sendMove => sub { return (100, 100); },
    sendLook => sub { return (0, 0); },
    sendAction => sub { return (12345, 0); },
    sendChat => sub { return ("Hello, world!"); },
    sendPrivateMsg => sub { return ("Player", "Hello there!"); },
    sendPartyChat => sub { return ("Party message"); },
    sendGuildChat => sub { return ("Guild message"); },
    sendItemUse => sub { return (1234, 5678); },
    sendDrop => sub { return (1234, 10); },
    sendTake => sub { return (1234); },
    sendStorageAdd => sub { return (1234, 10); },
    sendStorageGet => sub { return (1234, 10); },
    sendStoragePassword => sub { return ("password", 3); },
    sendSkillUse => sub { return (1, 10, 12345); },
    sendSkillUseLoc => sub { return (1, 10, 100, 100); },
    sendSync => sub { return (); },
    sendRestart => sub { return (0); },
    sendTalk => sub { return (12345); },
    sendTalkCancel => sub { return (12345); },
    sendTalkContinue => sub { return (12345); },
    sendTalkResponse => sub { return (12345, 1); },
    sendTalkNumber => sub { return (12345, 100); },
    sendTalkText => sub { return (12345, "Hello NPC"); },
    sendPartyLeader => sub { return (12345); },
    sendPartyOption => sub { return (1, 1, 1); },
    
    # Reconstruct methods
    reconstruct_master_login => sub { 
        return ({ 
            username => "username", 
            password => "password", 
            version => 1, 
            master_version => 1 
        }); 
    },
    reconstruct_game_login => sub { 
        return ({ 
            accountID => 123456, 
            sessionID => "sessionID", 
            sessionID2 => "sessionID2", 
            accountSex => 1 
        }); 
    },
    reconstruct_sync => sub { return ({}); },
    reconstruct_character_move => sub { 
        return ({ 
            x => 100, 
            y => 100 
        }); 
    },
    reconstruct_public_chat => sub { 
        return ({ 
            message => "Hello, world!" 
        }); 
    },
    reconstruct_private_message => sub { 
        return ({ 
            privMsg => "Hello there!", 
            privMsgUser => "Player" 
        }); 
    },
    reconstruct_storage_password => sub { 
        return ({ 
            type => 3, 
            pass => "password" 
        }); 
    },
    reconstruct_party_chat => sub { 
        return ({ 
            message => "Party message" 
        }); 
    },
    reconstruct_buy_bulk_vender => sub { 
        return ({ 
            venderID => 12345, 
            venderCID => 67890, 
            items => [{ amount => 1, itemIndex => 1 }] 
        }); 
    },
    reconstruct_buy_bulk_buyer => sub { 
        return ({ 
            buyerID => 12345, 
            buyingStoreID => 67890, 
            items => [{ ID => 1, itemID => 1, amount => 1 }] 
        }); 
    },
    reconstruct_buy_bulk_openShop => sub { 
        return ({ 
            limitZeny => 1000000, 
            result => 1, 
            storeName => "My Shop", 
            items => [{ nameID => 1, amount => 1, price => 1000 }] 
        }); 
    },
    reconstruct_guild_chat => sub { 
        return ({ 
            message => "Guild message" 
        }); 
    },
    reconstruct_char_delete2_accept => sub { 
        return ({ 
            charID => 12345, 
            code => "123456" 
        }); 
    },
    reconstruct_client_hash => sub { 
        return ({ 
            type => 1 
        }); 
    },
    reconstruct_actor_move => sub { 
        return ({ 
            coords => pack("C3", 100, 100, 0) 
        }); 
    },
    reconstruct_cash_shop_buy => sub { 
        return ({ 
            kafra_points => 1000, 
            items => [{ amount => 1, item_id => 1 }] 
        }); 
    },
    reconstruct_item_list_window_selected => sub { 
        return ({ 
            ID => 1, 
            amount => 1 
        }); 
    },
    
    # Parse methods - these typically don't generate packets but modify args
    parse_master_login => sub { 
        return ({ 
            password_md5_hex => "0123456789abcdef0123456789abcdef" 
        }); 
    },
    parse_character_move => sub { 
        return ({ 
            coords => pack("C3", 100, 100, 0) 
        }); 
    },
    parse_public_chat => sub { 
        return ({}); 
    },
    parse_private_message => sub { 
        return ({ 
            privMsg => "Hello there!", 
            privMsgUser => "Player" 
        }); 
    },
    parse_party_chat => sub { 
        return ({}); 
    },
    parse_buy_bulk_vender => sub { 
        return ({ 
            itemInfo => pack("v2", 1, 1) . pack("v2", 2, 2) 
        }); 
    },
    parse_guild_chat => sub { 
        return ({}); 
    },
    parse_actor_move => sub { 
        return ({ 
            coords => pack("C3", 100, 100, 0) 
        }); 
    },
    parse_pet_evolution => sub { 
        return ({}); 
    },
    
    # Other methods
    encryptMessageID => sub { 
        my $msg = pack("v", 0x0064); # Example packet ID
        return (\$msg); 
    },
    cryptKeys => sub { 
        return (0x12345678, 0x12345678, 0x12345678); 
    },
    injectMessage => sub { 
        return ("Hello, world!"); 
    },
    injectAdminMessage => sub { 
        return ("Admin message"); 
    },
    pinEncode => sub { 
        return (12345, 1234); 
    },
    secureLoginHash => sub { 
        return ("password", "salt", 1); 
    },
    encrypt_password => sub { 
        return ("password", "salt"); 
    },
    rodex_delete_mail => sub { 
        return (12345); 
    },
    rodex_request_zeny => sub { 
        return (12345); 
    },
    rodex_request_items => sub { 
        return (12345); 
    },
    rodex_cancel_write_mail => sub { 
        return (); 
    },
    rodex_add_item => sub { 
        return (1, 1); 
    },
    rodex_remove_item => sub { 
        return (1); 
    },
    rodex_open_write_mail => sub { 
        return (0); 
    },
    rodex_checkname => sub { 
        return ("Player"); 
    },
    rodex_send_mail => sub { 
        return (1, "Player", "Title", "Body", [1, 2, 3], 1000); 
    },
    rodex_refresh_maillist => sub { 
        return (0, 0, 0); 
    },
    rodex_read_mail => sub { 
        return (12345); 
    },
    rodex_next_maillist => sub { 
        return (0, 0, 0); 
    },
    rodex_open_mailbox => sub { 
        return (0, 0); 
    },
    rodex_close_mailbox => sub { 
        return (); 
    },
    
    # Default generator for methods without specific generators
    DEFAULT => sub { return (); }
);

# Process each method
foreach my $method (sort @all_methods) {
    print "Processing $method...\n";
    
    # Skip methods that are not callable or are just utility functions
    next if $method =~ /^(sendRaw|sendToServer)$/;
    
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

print "\nAll done! Generated test data for " . scalar(@all_methods) . " methods\n";
print "Test data is available in $output_dir\n";
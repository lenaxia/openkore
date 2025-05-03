#!/usr/bin/env perl
# analyze_send_methods.pl - Analyze Send.pm methods to determine required arguments

use strict;
use warnings;
use FindBin qw($RealBin);
use lib "$RealBin/../../../../src";
use JSON::PP;
use Data::Dumper;
use File::Path qw(make_path);

# Path to Send.pm
my $send_pm_path = "$RealBin/../../../../src/Network/Send.pm";

# Output directory for analysis results
my $output_dir = "$RealBin/analysis";
make_path($output_dir) unless -d $output_dir;

# Read the Send.pm file
open my $fh, '<', $send_pm_path or die "Cannot open $send_pm_path: $!";
my $send_pm_content = do { local $/; <$fh> };
close $fh;

print "Analyzing Send.pm methods...\n";

# Extract method definitions
my %methods;
my $current_method = '';
my $in_method = 0;
my $brace_count = 0;
my @method_lines;

foreach my $line (split /\n/, $send_pm_content) {
    # Check for method definition
    if ($line =~ /^\s*sub\s+(\w+)\s*\{/) {
        $current_method = $1;
        $in_method = 1;
        $brace_count = 1;
        @method_lines = ($line);
        next;
    }
    
    if ($in_method) {
        push @method_lines, $line;
        
        # Count braces to track method boundaries
        $brace_count += () = $line =~ /\{/g;
        $brace_count -= () = $line =~ /\}/g;
        
        if ($brace_count == 0) {
            $in_method = 0;
            $methods{$current_method} = join("\n", @method_lines);
        }
    }
}

print "Found " . scalar(keys %methods) . " methods\n";

# Analyze method parameters
my %method_params;
my %method_args_keys;
my %method_param_types;
my %complex_methods;

foreach my $method (sort keys %methods) {
    print "Analyzing method: $method\n";
    
    my $code = $methods{$method};
    
    # Extract parameter list
    if ($code =~ /^\s*sub\s+\w+\s*\{\s*\n\s*my\s*\(([^)]+)\)\s*=\s*\@_;/m) {
        my $param_list = $1;
        my @params = split /\s*,\s*/, $param_list;
        $method_params{$method} = \@params;
        
        print "  Parameters: " . join(", ", @params) . "\n";
        
        # Analyze parameter types
        my @param_types;
        my $has_complex_params = 0;
        
        # Skip $self
        for (my $i = 1; $i < scalar(@params); $i++) {
            my $param = $params[$i];
            my $param_name = $param;
            $param_name =~ s/^\s+|\s+$//g; # Trim whitespace
            
            # Determine parameter type based on name
            my $param_type = "unknown";
            
            # Check for reference parameters
            if ($param_name =~ /\$r_/) {
                $param_type = "reference";
                $has_complex_params = 1;
            }
            # Check for common parameter types by name
            elsif ($param_name =~ /\$(name|title|message|msg|text|input|password|email|notice|map|body|user|playerName|receiver|sender|code|salt)/) {
                $param_type = "string";
            }
            elsif ($param_name =~ /\$(ID|id|accountID|charID|playerID|targetID|venderID|buyerID|mailID|guildID|monID|itemID|skillID)/) {
                $param_type = "numeric_id";
            }
            elsif ($param_name =~ /\$(amount|count|num|flag|type|level|lvl|lv|position|alignment|point|version|sex|page|seed|limit|public)/) {
                $param_type = "numeric";
            }
            elsif ($param_name =~ /\$(x|y)/) {
                $param_type = "coordinate";
            }
            elsif ($param_name =~ /\$items/) {
                $param_type = "items_array";
                $has_complex_params = 1;
            }
            elsif ($param_name =~ /\$args/) {
                $param_type = "args_hash";
                $has_complex_params = 1;
            }
            elsif ($param_name =~ /\@items/) {
                $param_type = "items_array";
                $has_complex_params = 1;
            }
            else {
                # Default to unknown
                $param_type = "unknown";
                $has_complex_params = 1;
            }
            
            push @param_types, $param_type;
        }
        
        $method_param_types{$method} = \@param_types;
        
        # Mark methods with complex parameters
        if ($has_complex_params) {
            $complex_methods{$method} = 1;
        }
        
        # Check if method takes $args parameter
        if (@params >= 2 && $params[1] =~ /\$args/) {
            print "  Method takes \$args parameter, analyzing usage...\n";
            
            # Extract keys accessed in $args hash
            my @keys;
            while ($code =~ /\$args\s*->\s*\{['"]?([^'"}]+)['"]?\}/g) {
                push @keys, $1;
            }
            while ($code =~ /\$args\s*\{\s*['"]?([^'"}\s]+)['"]?\s*\}/g) {
                push @keys, $1;
            }
            
            # Look for reconstruct calls
            if ($code =~ /\$self->reconstruct\(\{([^}]+)\}\)/) {
                my $reconstruct_args = $1;
                while ($reconstruct_args =~ /([^=>\s,]+)\s*=>\s*\$args\s*->\s*\{['"]?([^'"}]+)['"]?\}/g) {
                    push @keys, $2;
                }
                while ($reconstruct_args =~ /([^=>\s,]+)\s*=>\s*\$args\s*\{\s*['"]?([^'"}\s]+)['"]?\s*\}/g) {
                    push @keys, $2;
                }
            }
            
            # Remove duplicates
            my %seen;
            @keys = grep { !$seen{$_}++ } @keys;
            
            $method_args_keys{$method} = \@keys if @keys;
            
            print "  Args keys: " . join(", ", @keys) . "\n" if @keys;
        }
    }
    
    # Check for reconstruct method calls
    if ($code =~ /\$self->reconstruct\(\{([^}]+)\}\)/) {
        my $reconstruct_args = $1;
        print "  Calls reconstruct with: $reconstruct_args\n";
        
        # Extract switch value
        if ($reconstruct_args =~ /switch\s*=>\s*['"]([^'"]+)['"]/) {
            print "  Switch: $1\n";
        }
    }
    
    # Check for complex code patterns
    if ($code =~ /\$\$/ || $code =~ /\@\{/ || $code =~ /\%\{/) {
        $complex_methods{$method} = 1;
    }
}

# Save analysis results
my $analysis = {
    methods => \%methods,
    parameters => \%method_params,
    args_keys => \%method_args_keys,
    param_types => \%method_param_types,
    complex_methods => \%complex_methods
};

my $json_file = "$output_dir/send_methods_analysis.json";
open my $out_fh, '>', $json_file or die "Cannot open $json_file: $!";
print $out_fh encode_json($analysis);
close $out_fh;

# Save complex methods to a separate file for LLM processing
my $complex_file = "$output_dir/complex_methods.json";
open my $complex_fh, '>', $complex_file or die "Cannot open $complex_file: $!";

my $complex_data = {
    methods => {}
};

foreach my $method (sort keys %complex_methods) {
    $complex_data->{methods}->{$method} = {
        code => $methods{$method},
        parameters => $method_params{$method},
        param_types => $method_param_types{$method},
    };
    
    if (exists $method_args_keys{$method}) {
        $complex_data->{methods}->{$method}->{args_keys} = $method_args_keys{$method};
    }
}

print $complex_fh encode_json($complex_data);
close $complex_fh;

print "\nAnalysis complete. Results saved to $json_file\n";
print "Complex methods saved to $complex_file for LLM processing\n";

# Now create a script that uses this analysis to generate test data
my $generator_script = "$RealBin/generate_from_analysis.pl";
open my $gen_fh, '>', $generator_script or die "Cannot open $generator_script: $!";

print $gen_fh <<'EOT';
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
EOT

close $gen_fh;
chmod 0755, $generator_script;

print "\nCreated generator script: $generator_script\n";
print "To generate test data based on the analysis, run:\n";
print "  $generator_script\n";
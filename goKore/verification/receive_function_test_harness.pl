#!/usr/bin/perl
use strict;
use warnings;
use FindBin qw($RealBin);
use lib "$RealBin/../../..";
use JSON;
use Data::Dumper;

# Add OpenKore source directory to @INC
use lib "$RealBin/../../../src";
use Network::Receive;

# Function to test Network::Receive functions
sub test_receive_functions {
    my ($data) = @_;
    
    my $function_name = $data->{function_name};
    my $args = $data->{args} || {};
    
    print "Testing Network::Receive function: $function_name\n";
    
    # Create a Network::Receive object
    my $server = $args->{server} || 'ServerType0';
    my $receive = Network::Receive->create($server);
    
    # Check if the function exists
    if (!$receive->can($function_name)) {
        print "Error: Function $function_name does not exist in Network::Receive\n";
        return "ERROR: Function not found";
    }
    
    # Prepare arguments
    my @func_args = ();
    if (exists $args->{arguments} && ref($args->{arguments}) eq 'ARRAY') {
        @func_args = @{$args->{arguments}};
    }
    
    # Call the function
    my $result;
    eval {
        if (@func_args) {
            $result = $receive->$function_name(@func_args);
        } else {
            $result = $receive->$function_name();
        }
    };
    
    if ($@) {
        print "Error executing function: $@\n";
        return "ERROR: $@";
    }
    
    # Format the result
    my $formatted_result;
    if (defined $result) {
        if (ref($result) eq 'HASH') {
            $formatted_result = to_json($result, { pretty => 1 });
        } elsif (ref($result) eq 'ARRAY') {
            $formatted_result = to_json($result, { pretty => 1 });
        } else {
            $formatted_result = $result;
        }
    } else {
        $formatted_result = "undef";
    }
    
    print "Function executed successfully\n";
    print "Result: $formatted_result\n";
    
    return $formatted_result;
}

1;
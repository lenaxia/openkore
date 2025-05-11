#!/usr/bin/perl
use strict;
use warnings;
use PPI;
use JSON;
use Data::Dumper;

# Initialize JSON and logging
my $json = JSON->new->pretty;
my $output = {
    constants => {},
    exports => [],
    subs => [],
    gaps => [],
    other => [],
    export_tags => {},
    hashes => {}
};

# Process files
for my $file ('src/Network/Receive.pm', 'src/Network/Send.pm') {
    print "Processing $file...\n";
    
    my @constructs;
    my $grep_cmd = 'grep -E \'^[{}]|^);|^sub |^our |^use |^package |^#\' ' . $file . ' -n';

    open my $fh, '-|', $grep_cmd or die $!;
    
    while (<$fh>) {
        chomp;
        my ($line_num, $type, $content) = /^(\d+):\s*(sub|our|use|package|#|\}|use constant)?\s*(.*)/;
        next unless defined $type;
        if ($type eq 'sub') {
            # print out a print stderr if the subroutine name is parse_minimap_indicator

            $content =~ s/^\s*sub\s+//;  # Clean up the content
            next unless $content =~ /\S/;  # Only skip if completely empty
        }
        
        # Handle subroutines
        if ($type eq 'sub') {
            my $sub_name = $content =~ /^([^\s({]+)/ ? $1 : 'anonymous';
            my $sub_start = $line_num;
            my $sub_end;
            my $source_code = '';
            my $comments = '';
            
            # Get source code lines first
            open my $src_fh, '<', $file or die $!;
            my @lines = <$src_fh>;
            close $src_fh;
            
            # Find closing brace for this sub in the source file
            my $brace_count = 1;
            $sub_end = $sub_start;
            while ($sub_end <= $#lines) {
                $sub_end++;
                if ($lines[$sub_end-1] =~ /\{/) {
                    $brace_count++;
                }
                if ($lines[$sub_end-1] =~ /\}/) {
                    $brace_count--;
                    last if $brace_count == 0;
                }
            }
            
            # Extract just this subroutine's code
            $source_code = join('', @lines[$sub_start-1..$sub_end-1]);
            
            # Find comments before sub
            my $comment_line = $sub_start - 1;
            while ($comment_line > 0) {
                last unless $lines[$comment_line-1] =~ /^\s*(#|\=)/;
                $comments = $lines[$comment_line-1] . $comments;
                $comment_line--;
            }
            
            push @{$output->{subs}}, {
                name => $sub_name,
                file => $file,
                reference => "Network::Receive::$sub_name",
                start => $sub_start,
                end => $sub_end,
                comments => $comments,
                source => $source_code
            };
            next;
        }
        
        # Special handling for use constant blocks
        if ($type eq 'use' && $content =~ /^constant/) {
            my $start_line = $line_num;
            my $end_line;
            
            # Find the closing brace/paren
            while (<$fh>) {
                if (/^(\d+):(}|\);)/) {
                    $end_line = $1;
                    last;
                }
            }
            
            if ($end_line) {
                open my $src_fh, '<', $file or die "Can't open $file: $!";
                my @lines = <$src_fh>;
                close $src_fh;
                
                my $constant_block = join('', @lines[$start_line-1..$end_line-1]);
                push @constructs, {
                    line => $start_line,
                    type => 'constant_block',
                    content => $constant_block,
                    end => $end_line
                };
                
                # Parse different constant formats
                if ($constant_block =~ /use constant \{([^}]+)\}/s) {
                    # Hash style: use constant { NAME => value, ... }
                    my $constants_str = $1;
                    my @lines = split(/\n/, $constants_str);
                    my $current_line = $start_line + 1;
                    
                    foreach my $line (@lines) {
                        next if $line =~ /^\s*$/;
                        if ($line =~ /(\w+)\s*=>\s*([^,]+)/) {
                            $output->{constants}{$1} = {
                                name => $1,
                                value => $2,
                                file => $file,
                                line => $current_line,
                                start => $start_line,
                                end => $end_line,
                                export_tag => find_export_tag($1, $file)
                            };
                        }
                        $current_line++;
                    }
                } elsif ($constant_block =~ /use constant (\w+)\s*=>\s*\(([^)]+)\)/s) {
                    # Hashref style: use constant NAME => (...)
                    my $const_name = $1;
                    my $const_value = $2;
                    $output->{constants}{$const_name} = {
                        name => $const_name,
                        value => $const_value,
                        file => $file,
                        line => $start_line,
                        start => $start_line,
                        end => $end_line,
                        export_tag => find_export_tag($const_name, $file),
                        is_hashref => 1
                    };
                }
                next;
            }
        }
        
        push @constructs, {
            line => $line_num,
            type => $type,
            content => $content
        };
    }
    close $fh;
    
    # Rest of the processing logic...
}

# Generate JSON output
my $json_output = eval { $json->encode($output) };
if ($@) {
    warn "JSON generation error: $@";
    $json_output = "{}";
}

# Write JSON to file
open my $out_fh, '>', 'packet_analysis.json' or die "Can't write output file: $!";
print $out_fh $json_output;
close $out_fh;

# Helper to find which export tag a constant belongs to
sub find_export_tag {
    my ($const_name, $file) = @_;
    
    open my $fh, '<', $file or return undef;
    my @lines = <$fh>;
    close $fh;
    
    # First check if it's directly in @EXPORT
    foreach my $line (@lines) {
        if ($line =~ /\@EXPORT\s*=\s*\(/ .. $line =~ /\)/) {
            return 'direct' if $line =~ /\b$const_name\b/;
        }
    }
    
    # Then check export tags
    foreach my $line (@lines) {
        if ($line =~ /\%EXPORT_TAGS\s*=\s*\(/ .. $line =~ /\)/) {
            if ($line =~ /(\w+)\s*=>\s*\[.*\b$const_name\b/) {
                return $1;
            }
        }
    }
    
    return undef;
}

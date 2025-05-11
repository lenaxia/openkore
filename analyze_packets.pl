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
        
        # Special handling for hash declarations
        if ($type eq 'our' && $content =~ /^\%(\w+)\s*=\s*\(/) {
            my $hash_name = $1;
            my $hash_data = parse_hash($fh, $file, $line_num, $hash_name);
            if ($hash_data) {
                $output->{hashes}{$hash_name} = $hash_data;
                next;
            }
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
                   
                   # Parse hashref entries
                   my %hash_values;
                   while ($const_value =~ /(\w+)\s*=>\s*([^,]+)/g) {
                       my ($key, $val) = ($1, $2);
                       $val =~ s/^\s+|\s+$//g;
                       $hash_values{$key} = $val;
                   }
                   
                   $output->{constants}{$const_name} = {
                       name => $const_name,
                       value => \%hash_values,
                       file => $file,
                       line => $start_line,
                       start => $start_line,
                       end => $end_line,
                       export_tag => find_export_tag($const_name, $file),
                       is_hashref => 1,
                       hash_values => \%hash_values
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
# Extract handler signature from sub declaration
sub extract_handler_signature {
    my ($line) = @_;
    if ($line =~ /sub\s*\{([^}]*)\}/) {
        my $params = $1;
        $params =~ s/^\s+|\s+$//g;
        return $params || '()';
    }
    return '()';
}

sub find_export_tag {
    my ($const_name, $file) = @_;
    
    open my $fh, '<', $file or return undef;
    my @lines = <$fh>;
    close $fh;
    
    my %export_tags;
    my $in_export_tags = 0;
    my $current_tag;
    my $in_tag_array = 0;
    
    # First check if it's directly in @EXPORT
    my $in_export = 0;
    foreach my $line (@lines) {
        if ($line =~ /\@EXPORT\s*=\s*\(/) {
            $in_export = 1;
            next;
        }
        if ($in_export && $line =~ /\)/) {
            last;
        }
        if ($in_export && $line =~ /\b$const_name\b/) {
            print STDERR "Found $const_name in direct EXPORT\n";
            return 'direct';
        }
        
        # Parse EXPORT_TAGS structure
        if ($line =~ /\%EXPORT_TAGS\s*=\s*\(/) {
            print STDERR "Entered EXPORT_TAGS section\n";
            $in_export_tags = 1;
            next;
        }
        
        if ($in_export_tags) {
            # Start of new tag group
            if ($line =~ /^\s*(\w+)\s*=>\s*\[/) {
                $current_tag = $1;
                $in_tag_array = 1;
                $export_tags{$current_tag} = [];
                print STDERR "Found tag group: $current_tag\n";
                next;
            }
            
            # Inside a tag array - collect all constants
            if ($in_tag_array) {
                # Extract all constants from this line
                while ($line =~ /(\b\w+\b)/g) {
                    my $const = $1;
                    next if $const eq 'qw'; # Skip Perl quote operator
                    push @{$export_tags{$current_tag}}, $const;
                    print STDERR "Added $const to $current_tag\n";
                    
                    # Check if this is our target constant
                    if ($const eq $const_name) {
                        print STDERR "Matched $const_name to tag $current_tag\n";
                        return $current_tag;
                    }
                }
                
                # End of current tag array
                if ($line =~ /\]/) {
                    print STDERR "End of tag array for $current_tag\n";
                    $in_tag_array = 0;
                    $current_tag = undef;
                }
            }
            
            # End of EXPORT_TAGS
            if ($line =~ /\)/) {
                print STDERR "Exited EXPORT_TAGS section\n";
                $in_export_tags = 0;
                last;
            }
        }
    }
    
    # Check all collected tags if we didn't find a direct match
    foreach my $tag (keys %export_tags) {
        if (grep { $_ eq $const_name } @{$export_tags{$tag}}) {
            print STDERR "Found $const_name in tag $tag (delayed match)\n";
            return $tag;
        }
    }
    
    print STDERR "No export tag found for $const_name\n";
    return undef;
}

# Parse a hash declaration and its contents
sub parse_hash {
    my ($fh, $file, $line_num, $hash_name) = @_;
    my $start_line = $line_num;
    my $end_line;
    my @entries;
    my $in_pod = 0;
    my $pod_content = '';
    
    # Find closing parenthesis
    while (<$fh>) {
        if (/^(\d+):\)/) {
            $end_line = $1;
            last;
        }
    }
    
    if ($end_line) {
        open my $src_fh, '<', $file or die "Can't open $file: $!";
        my @lines = <$src_fh>;
        close $src_fh;
        
        my $current_line = $start_line + 1;
        my $current_entry;
        my $current_tag;
        my $in_tag_array = 0;
        
        for my $i ($start_line..$end_line-1) {
            my $line = $lines[$i];
            
            # Handle POD documentation
            if ($line =~ /^=pod/) {
                $in_pod = 1;
                $pod_content = '';
                next;
            }
            if ($line =~ /^=cut/) {
                $in_pod = 0;
                if ($current_entry) {
                    $current_entry->{pod} = $pod_content;
                }
                next;
            }
            if ($in_pod) {
                $pod_content .= $line;
                next;
            }
            
            # Special handling for EXPORT_TAGS format
            if ($hash_name eq 'EXPORT_TAGS') {
                # Start of new tag group
                if ($line =~ /^\s*(\w+)\s*=>\s*\[/) {
                    $current_tag = $1;
                    $in_tag_array = 1;
                    $current_entry = {
                        key => $current_tag,
                        line => $current_line,
                        type => 'tag_group',
                        values => [],
                        content => $line
                    };
                    push @entries, $current_entry;
                    next;
                }
                
                # Inside tag array - collect constants
                if ($in_tag_array) {
                    # Extract all constants from this line, handling qw() lists
                    if ($line =~ /qw\(([^)]+)\)/) {
                        my @consts = split(/\s+/, $1);
                        push @{$current_entry->{values}}, @consts;
                    } else {
                        while ($line =~ /(\b\w+\b)/g) {
                            my $const = $1;
                            next if $const eq 'qw'; # Skip Perl quote operator
                            push @{$current_entry->{values}}, $const;
                        }
                    }
                    
                    # End of current tag array
                    if ($line =~ /\]/) {
                        $in_tag_array = 0;
                        $current_entry->{content} .= $line;
                        $current_entry->{end} = $current_line;
                        $current_entry = undef;
                    } elsif ($current_entry) {
                        $current_entry->{content} .= $line;
                    }
                    next;
                }
            }
            
            # Default hash entry parsing
            if ($line =~ /^\s*(\w+)\s*,\s*sub\s*\{/) {
                my $handler_name = $1;
                $current_entry = {
                    key => $handler_name,
                    line => $current_line,
                    type => 'handler',
                    content => $line,
                    handler_name => $handler_name,
                    handler_type => 'sub',
                    signature => extract_handler_signature($line)
                };
                push @entries, $current_entry;
            } elsif ($line =~ /^\s*(\w+)\s*,\s*([^,]+)/) {
                $current_entry = {
                    key => $1,
                    line => $current_line,
                    type => 'value',
                    content => $2
                };
                push @entries, $current_entry;
            } elsif ($current_entry && $line =~ /\}/) {
                $current_entry->{content} .= $line;
                $current_entry->{end} = $current_line;
                $current_entry = undef;
            } elsif ($current_entry) {
                $current_entry->{content} .= $line;
            }
            
            $current_line++;
        }
        
        return {
            name => $hash_name,
            start => $start_line,
            end => $end_line,
            entries => \@entries,
            file => $file
        };
    }
    return undef;
}

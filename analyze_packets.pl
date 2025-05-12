#!/usr/bin/perl
use strict;
use warnings;
use PPI;
use JSON;
use Data::Dumper;
use File::Basename;

# Configuration
my $output_dir = './output';  # Configurable output directory

# Check command line arguments
unless (@ARGV) {
    die "Usage: $0 <input_file.pm>\n";
}
my $input_file = $ARGV[0];

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

# Process input file
my $file = $input_file;
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
    
    
    # Generate JSON output
my $json_output = eval { $json->encode($output) };
if ($@) {
    warn "JSON generation error: $@";
    $json_output = "{}";
}

    # Create output directory if needed
    unless (-d $output_dir) {
        mkdir $output_dir or die "Cannot create output directory $output_dir: $!";
    }

    # Generate output filename from input filename
my $base_name = basename($input_file, '.pm');
my $output_file = "$output_dir/${base_name}_analysis.json";

# Write JSON to file
open my $out_fh, '>', $output_file or die "Can't write output file $output_file: $!";
print $out_fh $json_output;
close $out_fh;

print "Analysis complete. Output written to $output_file\n";

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
    
    # First check if we have eval-parsed EXPORT_TAGS data
    if (exists $output->{hashes}{EXPORT_TAGS} &&
        exists $output->{hashes}{EXPORT_TAGS}{eval_data}) {
        my $tags = $output->{hashes}{EXPORT_TAGS}{eval_data};
        
        # Debug output showing full EXPORT_TAGS structure
        print STDERR "DEBUG: Full EXPORT_TAGS structure:\n";
        print STDERR Dumper($output->{hashes}{EXPORT_TAGS}) . "\n";
        
        print STDERR "DEBUG: Checking export tags for $const_name. Available tags: " .
            join(', ', keys %$tags) . "\n";
        
        foreach my $tag (keys %$tags) {
            # Handle both array ref and scalar cases
            my @tag_values;
            if (ref($tags->{$tag}) eq 'ARRAY') {
                @tag_values = @{$tags->{$tag}};
            } elsif (ref($tags->{$tag})) {
                warn "Unexpected ref type for tag $tag: " . ref($tags->{$tag});
                next;
            } else {
                @tag_values = split(/\s+/, $tags->{$tag});
            }
            
            # More detailed debug output
            print STDERR "DEBUG: Tag '$tag' has values:\n";
            print STDERR "  - $_\n" for @tag_values;
            
            # Exact match against full constant name
            if (grep { $_ eq $const_name } @tag_values) {
                print STDERR "DEBUG: Found exact match for $const_name in tag $tag\n";
                return $tag;
            }
            # Also check for case variations if exact match fails
            if (grep { lc($_) eq lc($const_name) } @tag_values) {
                print STDERR "DEBUG: Found case-insensitive match for $const_name in tag $tag\n";
                return $tag;
            }
        }
    }
    
    # Fallback to more robust @EXPORT check
    open my $fh, '<', $file or do {
        print STDERR "DEBUG: Couldn't open $file for export tag check\n";
        return undef;
    };
    
    my $file_content = do { local $/; <$fh> };
    close $fh;
    
    # Check for \@EXPORT = qw(...) style
    if ($file_content =~ /\@EXPORT\s*=\s*qw\(([^)]*)\)/s) {
        my $export_list = $1;
        if ($export_list =~ /\b$const_name\b/) {
            print STDERR "DEBUG: Found $const_name in \@EXPORT qw() list\n";
            return 'direct';
        }
    }
    
    # Check for \@EXPORT = (...) style
    if ($file_content =~ /\@EXPORT\s*=\s*\(([^)]*)\)/s) {
        my $export_list = $1;
        if ($export_list =~ /\b$const_name\b/) {
            print STDERR "DEBUG: Found $const_name in \@EXPORT list\n";
            return 'direct';
        }
    }
    
    print STDERR "DEBUG: No export tag found for $const_name after full search. Tags available: " .
        (exists $output->{hashes}{EXPORT_TAGS} ?
         join(', ', keys %{$output->{hashes}{EXPORT_TAGS}{eval_data}}) : 'none') . "\n";
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
                # Extract the full hash definition
                my @hash_lines = @lines[$start_line-1..$end_line-1];
                my %export_tags;
                
                # Join all lines and parse tag groups that may span multiple lines
                my $full_text = join('', @hash_lines);
                while ($full_text =~ /(\w+)\s*=>\s*\[qw\(([^)]+)\)\]/gs) {
                    my ($tag, $constants) = ($1, $2);
                    # Clean up constants by removing newlines and extra whitespace
                    $constants =~ s/\s+/ /g;
                    $constants =~ s/^\s+|\s+$//g;
                    my @values = split(/\s+/, $constants);
                    
                    push @entries, {
                        key => $tag,
                        line => $start_line,
                        type => 'tag_group',
                        values => \@values,
                        content => $line,
                        end => $end_line
                    };
                    
                    $export_tags{$tag} = \@values;
                }
                
                return {
                    eval_data => \%export_tags,  # Store parsed data
                    name => $hash_name,
                    start => $start_line,
                    end => $end_line,
                    entries => \@entries,
                    file => $file
                };
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

# Experience and Level Handlers

## Experience Processing

### exp() - lines 901-950
```perl
sub exp {
    my ($self, $args) = @_;
    my $type = $args->{type} || 0;
    my $amount = $args->{amount};
    my $expType = $args->{expType} || 0;
    my $totalExp = $args->{totalExp};

    if ($type == 0) {  # ZC_NOTIFY_EXP
        $amount = unpack("V", $amount) if (ref($amount) eq 'SCALAR');
        $totalExp = unpack("V", $totalExp) if (ref($totalExp) eq 'SCALAR');
    } elsif ($type == 1) {  # ZC_NOTIFY_EXP2
        $amount = unpack("V", $amount) if (ref($amount) eq 'SCALAR');
        $totalExp = unpack("V", $totalExp) if (ref($totalExp) eq 'SCALAR');
    }

    my $expPercent;
    if ($totalExp > 0) {
        $expPercent = ($amount / $totalExp) * 100;
    } else {
        $expPercent = 0;
    }

    if ($expType == 0) {  # Normal exp
        message TF("Gained %s experience (%.2f%%)\n", formatNumber($amount), $expPercent), "exp";
    } elsif ($expType == 1) {  # Quest exp
        message TF("Gained %s quest experience (%.2f%%)\n", formatNumber($amount), $expPercent), "exp";
    } elsif ($expType == 2) {  # Party exp
        message TF("Gained %s party experience (%.2f%%)\n", formatNumber($amount), $expPercent), "exp";
    }

    if ($char) {
        $char->{exp} += $amount;
        $char->{exp_max} = $totalExp;
        $char->{exp_percent} = $expPercent;
    }

    Plugins::callHook('packet_exp', {
        amount => $amount,
        expType => $expType,
        totalExp => $totalExp,
        expPercent => $expPercent
    });
}
```

### levelUp() - lines 952-1000
```perl
sub levelUp {
    my ($self, $args) = @_;
    my $type = $args->{type} || 0;  # 0=base, 1=job
    my $level = $args->{level};

    if ($type == 0) {  # Base level
        message TF("Reached Base Level %s!\n", $level), "info";
        if ($char) {
            $char->{lv} = $level;
            $char->{points_free} += $config{points_per_lvl} || 1;
        }
    } elsif ($type == 1) {  # Job level
        message TF("Reached Job Level %s!\n", $level), "info";
        if ($char) {
            $char->{lv_job} = $level;
            $char->{skill_points} += $config{skill_points_per_lvl} || 1;
        }
    }

    if ($char && $type == 0) {
        my $delta = $level - $char->{lastBaseLvl};
        if ($delta > 0) {
            $char->{hp_max} += $delta * $config{hp_per_lvl} if $config{hp_per_lvl};
            $char->{sp_max} += $delta * $config{sp_per_lvl} if $config{sp_per_lvl};
            $char->{lastBaseLvl} = $level;
        }
    }

    Plugins::callHook('packet_levelUp', {
        type => $type,
        level => $level
    });
}
```

### Key Features:
- Handles different experience types (normal, quest, party)
- Calculates experience percentages
- Updates character stats
- Provides detailed user feedback
- Supports plugin hooks
- Manages level-up bonuses

## Core Experience Calculations (stat_info_handlers) - lines 1330-1420
```perl
VAR_EXP, sub {
    my ($actor, $value) = @_;
    $actor->{exp_last} = $actor->{exp};
    $actor->{exp} = $value;

    if ($actor->{lastBaseLvl} eq $actor->{lv}) {
        $monsterBaseExp = $actor->{exp} - $actor->{exp_last};
    } else {
        $monsterBaseExp = $actor->{exp_max_last2} - $actor->{exp_last} + $actor->{exp};
        $actor->{lastBaseLvl} = $actor->{lv};
        $actor->{exp_max_last2} = $actor->{exp_max};
    }

    if ($monsterBaseExp > 0) {
        $totalBaseExp += $monsterBaseExp;
    }
},

VAR_JOBEXP, sub {
    my ($actor, $value) = @_;
    $actor->{exp_job_last} = $actor->{exp_job};
    $actor->{exp_job} = $value;

    if ($actor->{lastJobLvl} eq $actor->{lv_job}) {
        $monsterJobExp = $actor->{exp_job} - $actor->{exp_job_last};
    } else {
        $monsterJobExp = $actor->{exp_job_max_last2} - $actor->{exp_job_last} + $actor->{exp_job};
        $actor->{lastJobLvl} = $actor->{lv_job};
        $actor->{exp_job_max_last2} = $actor->{exp_job_max};
    }

    if ($monsterJobExp > 0) {
        $totalJobExp += $monsterJobExp;
    }

    my $basePercent = $char->{exp_max} ?
        ($monsterBaseExp / $char->{exp_max} * 100) :
        0;
    my $jobPercent = $char->{exp_job_max} ?
        ($monsterJobExp / $char->{exp_job_max} * 100) :
        0;
    message TF("%s have gained %d/%d (%.2f%%/%.2f%%) Exp\n", $char, $monsterBaseExp, $monsterJobExp, $basePercent, $jobPercent), "exp";
    Plugins::callHook('exp_gained');
}
```

### Key Features:
- Handles both base and job experience calculations
- Tracks experience across level changes
- Maintains historical exp data
- Calculates experience percentages
- Provides detailed exp gain messages
- Supports plugin hooks for exp events
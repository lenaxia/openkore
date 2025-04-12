# Status Effect System

## Validation Rules from Actor.pm
```go
// Thread-safe status registry with atomic counters
type StatusRegistry struct {
    mu            sync.RWMutex
    immunities    map[EntityID]*lockFreeSet[StatusID]
    activeEffects map[EntityID]*lockFreeMap[StatusID, *StatusEffect]
}

// NewStatusRegistry creates registry with Perl-compat defaults
func NewStatusRegistry() *StatusRegistry {
    return &StatusRegistry{
        immunities:    make(map[EntityID]*lockFreeSet[StatusID]),
        activeEffects: make(map[EntityID]*lockFreeMap[StatusID, *StatusEffect]),
    }
}

// CanApply validates status application with Perl-compatible rules
func (sr *StatusRegistry) CanApply(target Entity, effect StatusEffect) error {
    sr.mu.RLock()
    defer sr.mu.RUnlock()

    // Check global immunities
    if immSet := sr.immunities[target.ID()]; immSet != nil && immSet.Contains(effect.ID) {
        return fmt.Errorf("%w: global immunity", ErrStatusImmune)
    }

    // Check existing effects
    effectMap := sr.activeEffects[target.ID()]
    if effectMap != nil {
        // Check refresh rules from Actor::setStatus
        if current, exists := effectMap.Get(effect.ID); exists {
            if !current.Flags.AllowRefresh {
                return fmt.Errorf("%w: status exists", ErrStatusExists)
            }
            if current.Source != effect.Source && !current.Flags.AllowStack {
                return fmt.Errorf("%w: source mismatch", ErrStatusConflict)
            }
        }

        // Check bidirectional conflicts
        effectMap.Range(func(id StatusID, e *StatusEffect) bool {
            if slices.Contains(e.Overrides, effect.ID) {
                return fmt.Errorf("%w: %s overrides %s", 
                    ErrStatusConflict, e.ID, effect.ID)
            }
            if slices.Contains(effect.BlockedBy, id) {
                return fmt.Errorf("%w: %s blocked by %s", 
                    ErrStatusConflict, effect.ID, id)
            }
            return true
        })
    }

    // Check target-specific immunities
    if target.Immunities().Contains(effect.ID) {
        return fmt.Errorf("%w: target immunity", ErrStatusImmune)
    }

    return nil
}

// AddImmunity atomically updates immunity state
func (sr *StatusRegistry) AddImmunity(target Entity, status StatusID) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    
    if sr.immunities[target.ID()] == nil {
        sr.immunities[target.ID()] = newLockFreeSet[StatusID]()
    }
    sr.immunities[target.ID()].Insert(status)
}

// RemoveImmunity atomically updates immunity state  
func (sr *StatusRegistry) RemoveImmunity(target Entity, status StatusID) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    
    if set := sr.immunities[target.ID()]; set != nil {
        set.Remove(status)
    }
}

// Atomic immunity updates
func (sv *StatusValidator) AddImmunity(target Entity, status StatusID) {
    sv.mutex.Lock()
    defer sv.mutex.Unlock()
    
    if _, exists := sv.immunities[target.ID()]; !exists {
        sv.immunities[target.ID()] = make(map[StatusID]bool)
    }
    sv.immunities[target.ID()][status] = true
}

func (sv *StatusValidator) RemoveImmunity(target Entity, status StatusID) {
    sv.mutex.Lock()
    defer sv.mutex.Unlock()
    
    if immunities, exists := sv.immunities[target.ID()]; exists {
        delete(immunities, status)
    }
}
    
    // Check existing status conflicts
    current := target.GetStatus(status.ID)
    if current != nil {
        // From OpenKore's status refresh rules
        if !current.Flags.Refreshable {
            return ErrStatusExists
        }
        if status.Source != current.Source && !status.Flags.Stackable {
            return ErrStatusConflict
        }
    }
    
    // Check mutual blocking statuses
    for _, existing := range target.Statuses() {
        // Original Perl logic: check both directions
        if contains(existing.Overrides, status.ID) {
            return fmt.Errorf("%w: %s blocks %s", 
                ErrStatusConflict, existing.ID, status.ID)
        }
        if contains(status.BlockedBy, existing.ID) {
            return fmt.Errorf("%w: %s blocked by existing %s", 
                ErrStatusConflict, status.ID, existing.ID)
        }
    }
    
    // Check target immunities
    if target.Immunities().Has(status.ID) {
        return fmt.Errorf("%w: %s", ErrStatusImmune, status.ID)
    }
    
    return nil
}
```

## Core Types
```go
type StatusEffect struct {
    ID          StatusID
    StartedAt   time.Time
    Duration    time.Duration  
    Source      EntityID
    Stacks      int
    Flags       StatusFlags
    Overrides   []StatusID // Statuses this overrides
    BlockedBy   []StatusID // Statuses that block this
}

// StatusApplication rules
type StatusRules struct {
    MaxStacks      int
    RefreshPolicy  RefreshPolicy // Extend/Reset/Ignore
    StackPolicy    StackPolicy   // Independent/Combine/Replace
    CancelOnDamage bool          // If damage removes status
}

// StatusResolver handles interactions
type StatusResolver struct {
    rules map[StatusID]StatusRules
}

func (sr *StatusResolver) CanApply(
    target Entity, 
    status StatusEffect,
) bool {
    // Check immunity
    if target.Immunities().Has(status.ID) {
        return false
    }
    
    // Check existing status conflicts
    for _, existing := range target.Statuses() {
        if contains(existing.Overrides, status.ID) {
            return false
        }
        if contains(status.BlockedBy, existing.ID) {
            return false
        }
    }
    
    return true
}

func (sr *StatusResolver) CalculateStacks(
    target Entity,
    status StatusEffect,
) int {
    rules := sr.rules[status.ID]
    current := target.GetStatusStacks(status.ID)
    
    switch rules.StackPolicy {
    case StackPolicyIndependent:
        return current + 1
    case StackPolicyCombine:
        return min(current+1, rules.MaxStacks)
    case StackPolicyReplace:
        return 1
    default:
        return 1
    }
}

type StatusFlags uint32
const (
    FlagBuff StatusFlags = 1 << iota
    FlagDebuff
    FlagDispellable
    FlagRefreshable
)

type StatusHandler interface {
    OnApply(target Entity)
    OnTick(target Entity)
    OnRemove(target Entity)
    OnExpire(target Entity)
}
```

## Common Status Effects
| Status ID      | Type    | Flags               | Description            |
|----------------|---------|---------------------|------------------------|
| STATUS_POISON  | Debuff  | Dispellable         | HP drain over time     |
| STATUS_HASTE   | Buff    | Refreshable         | Increased movement speed |
| STATUS_SILENCE | Debuff  | Dispellable         | Blocks skill usage     |

## Status Application
```go
type StatusApplier struct {
    validator    *StatusValidator
    statusStore  *StatusStore
    taskManager  TaskManager
}

func (sa *StatusApplier) ApplyStatus(target Entity, status StatusEffect) error {
    if err := sa.validator.CanApply(target, status); err != nil {
        return err
    }

    sa.statusStore.Lock(target.ID())
    defer sa.statusStore.Unlock(target.ID())
    
    existing := sa.statusStore.Get(target.ID(), status.ID)
    if existing != nil {
        if existing.Flags&FlagRefreshable != 0 {
            // Refresh duration and return existing
            sa.statusStore.RefreshDuration(target.ID(), status.ID, status.Duration)
            return nil
        }
        return ErrStatusExists
    }

    // Add new status with expiration task
    sa.statusStore.Add(target.ID(), status)
    
    // Schedule expiration if duration > 0
    if status.Duration > 0 {
        sa.taskManager.Add(task.NewTimedTask(
            status.Duration,
            func() {
                sa.statusStore.Remove(target.ID(), status.ID)
                target.OnStatusExpired(status.ID)
            },
        ))
    }
    
    return nil
}
```

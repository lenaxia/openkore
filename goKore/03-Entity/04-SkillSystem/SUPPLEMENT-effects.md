# Skill Effect System

## Core Interfaces
```go
type Effect interface {
    Apply(caster, target Entity) EffectResult
    Duration() time.Duration
    Stackable() bool
}

type EffectResult struct {
    Success      bool
    Damage       int         // For damage effects  
    Heal         int         // For heal effects
    Status       StatusID    // For status effects
    Message      string      // Feedback message
    ShouldRetry  bool        // Whether effect can be retried
    StateVersion uint64      // Version for reconciliation
    IsStateful   bool        // Whether effect modifies entity state
}

type QuarantineEntry struct {
    EffectID     EffectID
    StartedAt    time.Time
    ExpiresAt    time.Time 
    RetryCount   int
    LastError    string
    LastAttempt  time.Time
}

type RetryPolicy struct {
    effect_retry_window  int // Base delay in ms
    quarantine_ttl       time.Duration
    max_retries          int
    jitter_factor        float64 // 0.0-1.0
}

// Effect Types
type DamageEffect struct {
    Element     ElementType
    Power       int
    Variance    float64 // 0.0-1.0
    IsMagical   bool
}

type HealEffect struct {
    Base        int
    MatkFactor  float64
}

type StatusEffect struct {
    Status      StatusID
    Duration    time.Duration
    Chance      float64 // 0.0-1.0
}

// Composite Effects
type CompositeEffect struct {
    Effects []Effect
}

func (ce *CompositeEffect) Apply(caster, target Entity) EffectResult {
    var combined EffectResult
    for _, effect := range ce.Effects {
        res := effect.Apply(caster, target)
        combined.Damage += res.Damage
        combined.Heal += res.Heal
        // ... combine other fields
    }
    return combined
}
```

## Effect Application
```go
type EffectApplier struct {
    rand          rand.Source
    metrics       ReconciliationMetrics
    quarantined   map[EntityID]QuarantineEntry
    retryPolicy   RetryPolicy
    effectSys     *EffectSystem
    mux           sync.RWMutex
}

func (ea *EffectApplier) Apply(effect Effect, caster, target Entity) EffectResult {
    start := time.Now()
    
    // Check quarantine status with retry policy
    ea.mux.RLock()
    if qEntry, ok := ea.quarantined[target.ID]; ok {
        baseDelay := time.Duration(ea.retryPolicy.effect_retry_window) * time.Millisecond
        jitter := time.Duration(float64(baseDelay) * ea.retryPolicy.jitter_factor * rand.Float64())
        delay := baseDelay << qEntry.retryCount + jitter
        
        if time.Since(qEntry.lastAttempt) < delay {
            ea.metrics.TrackQuarantinedEffect(effect.GetID(), qEntry.retryCount)
            ea.mux.RUnlock()
            return EffectResult{Success: false, Message: "effect in quarantine cooling-off period"}
        }
        qEntry.retryCount++
        qEntry.lastAttempt = time.Now()
        ea.quarantined[target.ID] = qEntry
    }
    ea.mux.RUnlock()

    // Check resistance
    if target.Resists(effect) {
        ea.metrics.TrackResistedEffect(effect.GetID())
        return EffectResult{
            Success: false,
            Message: "Target resisted effect",
        }
    }
    
    // Apply effect
    result := effect.Apply(caster, target)
    
    // Handle status effects
    if result.Status != 0 {
        target.Statuses().Add(StatusEffect{
            ID: result.Status,
            Duration: effect.Duration(),
            Source: caster.ID(),
        })
    }
    
    duration := time.Since(start)
    
    // Track effect outcome metrics with quarantine integration
    outcome := "success"
    if !result.Success {
        outcome = "failed"
        if shouldQuarantine(result) {
            ea.mux.Lock()
            ea.quarantined[target.ID] = QuarantineEntry{
                EffectID:    effect.GetID(),
                StartedAt:   time.Now(),
                ExpiresAt:   time.Now().Add(ea.retryPolicy.quarantine_ttl),
                RetryCount:  0,
                LastError:   result.Message,
            }
            ea.mux.Unlock()
            ea.metrics.TrackQuarantineStart(target.ID, effect.GetID())
        }
    } else if result.Damage > 0 {
        outcome = "damage"
    } else if result.Heal > 0 {
        outcome = "heal"
    } else if result.Status != 0 {
        outcome = "status"
    }
    
    ea.metrics.RecordEffectApplication(
        caster.Type().String(),
        target.Type().String(),
        outcome,
        duration,
        result.Message,
    )
    
    // Link to reconciliation system
    if result.Success && effect.IsStateful() {
        ea.effectSys.NotifyReconciliation(
            caster.ID,
            target.ID,
            effect.GetID(),
            result,
        )
    }
    
    // Metrics documentation linked to SUPPLEMENT-metrics.md
    // Tracked dimensions: damage/heal amounts, status applications, resistances
    // See "Effect Outcome Tracking" in SUPPLEMENT-metrics.md
    
    return result
}
```

## Common Effects
| Effect Type | Parameters | Description |
|-------------|------------|-------------|
| Damage      | Power, Element | Deals damage |
| Heal        | Base, MatkFactor | Restores HP |
| Status      | StatusID, Duration | Applies buff/debuff |
| Teleport    | Map, X, Y | Moves target |
| Summon      | MobID, Count | Spawns entities |

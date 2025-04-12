# Deferred Locking Supplement

## Migration Details
| C++ Pattern          | Go Approach          | Safety Guarantees     |
|----------------------|----------------------|-----------------------|
| MutexLocker          | defer Lock()         | Scope-based safety    |
| Explicit lock/unlock | Automatic via defer  | Compiler enforced     |

## TrackedMutex Design
```go
// SafeMutex implements Systems domain locking contract
type SafeMutex struct {
    mu      sync.Mutex
    systems SystemsProvider
}

// NewTrackedMutex creates a context-aware mutex with instrumentation
func NewTrackedMutex(name string, ctx context.Context, stats *MutexStats) *TrackedMutex {
    return &TrackedMutex{
        name:  name,
        ctx:   ctx,
        stats: stats,
    }
}

// MutexStats tracks contention metrics
type MutexStats struct {
    WaitCount    atomic.Int64
    WaitDuration atomic.Int64
    HoldDuration atomic.Int64 
}

// LockWithContext acquires mutex with context support
func (sm *SafeMutex) LockWithContext(ctx context.Context) error {
    locked := make(chan struct{})
    go func() {
        sm.mu.Lock()
        close(locked)
    }()
    
    select {
    case <-locked:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Unlock releases mutex with safety checks and metrics
func (tm *TrackedMutex) Unlock() {
    if tm.state.Load() != 1 {
        panic(fmt.Sprintf("%s: unlock of unlocked mutex", tm.name))
    }
    
    currentHolder := tm.holder.Load()
    if currentHolder != 0 && currentHolder != goroutineID() {
        panic(fmt.Sprintf("%s: unlock from non-owner goroutine", tm.name))
    }

    holdTime := time.Now().UnixNano() - tm.waitStart.Load()
    tm.stats.HoldDuration.Add(holdTime)
    
    tm.holder.Store(0)
    tm.state.Store(0)
    tm.base.Unlock()
}

// TryLockWithTimeout attempts lock with deadline
func (tm *TrackedMutex) TryLockWithTimeout(timeout time.Duration) bool {
    ctx, cancel := context.WithTimeout(tm.ctx, timeout)
    defer cancel()
    return tm.LockWithContext(ctx) == nil
}

// DetectDeadlocks checks for potential deadlocks
func (tm *TrackedMutex) DetectDeadlocks(threshold time.Duration) bool {
    return tm.waitStart != 0 && 
        (time.Now().UnixNano()-tm.waitStart) > int64(threshold)
}
```

## Migration Checklist
- [X] Define DeferLocker interface
- [X] Finalize TrackedMutex design
- [O] Update threadpool usage patterns
- [X] Document context propagation
- [O] Verify defer ordering semantics
- [X] Implement context-aware lock timeouts

## Performance Considerations
1. Defer overhead vs explicit unlocks
2. Channel vs condition variable timeouts
3. Stack trace capture costs
4. Goroutine ID tracking

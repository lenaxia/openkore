# Concurrency Patterns

## Entity State Management
```go
// Atomic state updates matching Actor.pm's deltaHP tracking
type EntityState struct {
    hp    atomic.Int32
    maxHp atomic.Int32 
    sp    atomic.Int32
    status atomic.Value // Stores []StatusEffect
}

// Thread-safe position updates
type EntityPosition struct {
    mu    sync.RWMutex
    pos   Position
    posTo Position
}

func (ep *EntityPosition) Update(from, to Position) {
    ep.mu.Lock()
    defer ep.mu.Unlock()
    ep.pos = from
    ep.posTo = to
}
```

## Repository Lock Hierarchy
```go
// Matches OpenKore's ActorList.pm synchronization
type EntityRepository struct {
    mu        sync.RWMutex // Main registry lock (matches ActorList.pm {IDmap} sync)
    entities  map[EntityID]*EntityEntry
    spatialMu sync.RWMutex // Spatial index lock
    
    // Nested locks must be acquired in order:
    // 1. Repository mu
    // 2. Spatial mu
    // 3. Individual entity locks
}

type EntityEntry struct {
    mu       sync.RWMutex
    entity   *Entity
    position *EntityPosition
}
```

## Worker Pool Pattern
```go
// Replaces Perl's callback lists with Go channels
const entityWorkers = 8

type EntityUpdate struct {
    ID     EntityID
    Update func(*Entity)
}

func StartEntityWorker(updates <-chan EntityUpdate) {
    go func() {
        for update := range updates {
            repo.Get(update.ID, func(e *Entity) {
                e.Lock()
                defer e.Unlock()
                update.Update(e)
            })
        }
    }()
}
```

## Thread Safety Matrix
| Operation          | Locks Required                  | Go Equivalent               |
|--------------------|---------------------------------|-----------------------------|
| Entity Update      | Entity RLock                   | sync.RWMutex.RLock()        |
| Inventory Modify   | Entity Lock + Inventory Lock   | sync.Mutex nested locks     |
| Status Effect Add  | Entity Lock + Status RWLock    | atomic.Value for status     |
| Position Update    | Entity Lock + SpatialMap Lock  | Position struct mutex       |
| Name Change        | Entity Lock + Registry Lock    | Repository mutex hierarchy  |

## Channel Patterns
```go
// Buffered channels match OpenKore's event queue sizing
type EntityEvents struct {
    NameChange     chan string       // Size 10
    PositionUpdate chan Position     // Size 100  
    StatusChange   chan []StatusEffect // Size 20
}

func NewEntityEvents() *EntityEvents {
    return &EntityEvents{
        NameChange:     make(chan string, 10),
        PositionUpdate: make(chan Position, 100),
        StatusChange:   make(chan []StatusEffect, 20),
    }
}
```

## Context Timeouts
```go
// Matches OpenKore's 50ms timeout for entity queries
func QueryEntities(ctx context.Context, pred func(*Entity) bool) ([]*Entity, error) {
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
    defer cancel()
    
    // Query logic with ctx.Done() checks
}
```

// Matches ActorList.pm's IDmap synchronization
func (r *EntityRegistry) GetByID(id EntityID) (*Entity, bool) {
    r.RLock()
    defer r.RUnlock()
    e, exists := r.entities[id]
    return e, exists
}
```

## Lock Hierarchy (From OpenKore's Core Patterns)
1. Entity registry lock 
2. Individual entity lock
3. Component locks (inventory/status/etc)
```go
func TransferItem(src, dest *Entity, item Item) error {
    src.Lock()
    defer src.Unlock()
    
    dest.Lock() 
    defer dest.Unlock()
    
    // Matches OpenKore's inventory item movement logic
    if err := src.Inventory.Remove(item); err != nil {
        return err
    }
    return dest.Inventory.Add(item)
}
```

## Atomic Updates (Like OpenKore's deltaHp Tracking)
```go
type EntityState struct {
    hp    atomic.Int32
    maxHp atomic.Int32
    sp    atomic.Int32
    status atomic.Value // stores StatusEffect slice
}

func (es *EntityState) ApplyDamage(amount int) int {
    current := es.hp.Load()
    newVal := current - int32(amount)
    if newVal < 0 {
        newVal = 0
    }
    es.hp.Store(newVal)
    return int(current - newVal) // Return actual damage dealt
}
```

## Channel Patterns (Replacing Perl's CallbackList)
```go
// Equivalent to Actor.pm's onNameChange CallbackList
type EntityEvents struct {
    NameChange chan string
    PositionUpdate chan Position
    StatusChange chan []StatusEffect
}

// Buffered channels match OpenKore's event queue sizing
func NewEntityEvents() *EntityEvents {
    return &EntityEvents{
        NameChange:     make(chan string, 10),
        PositionUpdate: make(chan Position, 100),
        StatusChange:   make(chan []StatusEffect, 20),
    }
}
```

## Deadlock Prevention Rules
1. Entity registry must be locked before individual entities
2. Status effects locked before inventory
3. Network events have priority over AI decisions
4. Use context timeouts for database operations:
```go
func QueryEntities(ctx context.Context, pred func(*Entity) bool) ([]*Entity, error) {
    ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond) // Match OpenKore's 50ms timeout
    defer cancel()
    
    // ... query logic ...
}
```

## Thread Safety Matrix
| Operation          | Locks Required                  | Perl Equivalent               |
|--------------------|---------------------------------|-------------------------------|
| Entity Update      | Entity RLock                   | Actor.pm field access         |
| Inventory Modify   | Entity Lock + Inventory Lock   | Actor::Item transaction       |
| Status Effect Add  | Entity Lock + Status RWLock    | Actor::setStatus              |
| Position Update    | Entity Lock + SpatialMap Lock  | Field.pm position validation  |
| Name Change        | Entity Lock + Registry Lock    | ActorList::getByName update   |

# Domain Interfaces Specification

## Core Design Principles
1. **Immutable State** - All interface methods that return state must provide immutable snapshots
2. **Thread Safety** - Implementations must be safe for concurrent access
3. **Versioned Contracts** - Interfaces include version identifiers for backward compatibility
4. **Error Handling** - All methods return standardized error types

## 1. Entity Domain Interfaces

### EntitySync
```go
// Tracks state changes for network synchronization
type EntitySync interface {
    // Version identifies the interface implementation
    Version() int

    // Get full current state (for initial sync)
    GetState() EntityState
    
    // Get only changed fields since last sync
    GetDelta() EntityDelta
    
    // Apply remote state updates
    ApplyDelta(EntityDelta) error
}

// EntityState contains all persistent entity data
type EntityState struct {
    ID        EntityID
    Position  Position
    Stats     BaseStats
    Statuses  []StatusEffect
    Equipment map[EquipSlot]ItemID
}

// EntityDelta contains changed fields only  
type EntityDelta struct {
    Modified map[string]interface{}
    Timestamp time.Time
}
```

### Combatant
```go
// Handles combat-related actions and state
type Combatant interface {
    // Core combat actions
    Attack(Target) (DamageReport, error)
    UseSkill(SkillContext) (SkillResult, error)
    
    // Status monitoring
    GetCooldowns() map[SkillID]time.Time
    GetAggroTable() map[EntityID]int
    
    // Event hooks
    OnDamage(func(DamageEvent))
    OnDeath(func(DeathEvent))
}
```

## 2. Network Domain Interfaces

### ProtocolCodec
```go
// Handles packet encoding/decoding
type ProtocolCodec interface {
    Version() int
    
    // Packet serialization
    Encode(NetworkMessage) ([]byte, error)
    Decode([]byte) (NetworkMessage, error)
    
    // Protocol features
    SupportsCompression() bool
    SupportsEncryption() bool
}

// NetworkMessage standard packet format
type NetworkMessage struct {
    Header  PacketHeader
    Payload []byte
    CRC     uint32
}
```

## 3. AI Domain Interfaces

### BehaviorController
```go
// Manages AI decision making
type BehaviorController interface {
    // Task management
    QueueTask(Task) (TaskID, error)
    CancelTask(TaskID) error
    
    // State evaluation
    EvaluateThreat(EntityID) ThreatLevel
    FindPath(Position) ([]Position, error)
    
    // Configuration
    SetPriority(PriorityProfile) error
}

// Task execution contract
type Task interface {
    ID() TaskID
    Priority() int
    Execute() TaskStatus
    Mutexes() []ResourceLock
}
```

## 4. Skill System Interfaces

### SkillResolver
```go
// Handles skill effect calculation
type SkillResolver interface {
    Version() int
    
    // Core resolution
    Resolve(SkillContext) (SkillResult, error)
    
    // Validation
    ValidateTarget(SkillContext) error
    CheckRange(SkillContext) bool
    
    // Cooldown tracking
    GetGlobalCooldown() time.Duration
}

type SkillContext struct {
    Caster    EntityState
    Skill     SkillDefinition
    Target    Target
    Timestamp time.Time
}
```

## Versioning Policy
1. Major version changes require interface redesign
2. Minor version changes allow additive changes
3. All implementations must handle at least N-1 version

## Error Handling Standards
```go
// Standard error codes
const (
    ErrInvalidState = 1
    ErrOutOfRange   = 2
    ErrCooldown     = 3
    ErrNoPath       = 4
)

// Error wrapper type
type DomainError struct {
    Code    int
    Message string
    Source  string
}
```

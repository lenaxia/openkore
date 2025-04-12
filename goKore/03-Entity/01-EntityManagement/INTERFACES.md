# Entity Management Interfaces

## Core Interfaces
```go
// Entity defines mandatory capabilities for all game entities
type Entity interface {
    // Identity (matches Actor.pm fields)
    ID() EntityID          // Actor->{ID} 
    LegacyID() uint32      // Actor->{nameID} conversion
    Type() EntityType      // Actor->{actorType}
    Name() string          // Actor->{name} with fallback
    
    // Spatial (mirrors Actor->position())
    Position() Position
    PreviousPosition() Position  // Actor->{pos}
    ValidatePosition() error     // Field.pm boundaries
    
    // State Management
    Stats() EntityStats          // HP/SP/BaseStats
    Inventory() InventoryHolder  // Actor::Item logic
    Statuses() []StatusEffect    // Actor->{statuses}
    
    // Relationships 
    Owner() Entity               // Pets/homunculus
    Party() PartyMember          // Party relationships
    
    // Perl Compatibility
    DistanceTo(other Entity) float64  // Actor->distance()
    Verb(you, other string) string    // Actor->verb()
    DeepCopy() Entity                 // Actor->deepCopy()
}

// EntityRepository handles persistence (ActorList.pm analog)
type EntityRepository interface {
    // Core Operations
    Get(id EntityID) (Entity, error)       // ActorList::getByID
    GetByLegacyID(id uint32) (Entity, error) // nameID lookup
    Add(e Entity) error                    // ActorList::add
    Remove(id EntityID) error              // ActorList::remove
    
    // Query Patterns
    FindByName(name string) []Entity       // Actor::getByName
    FindByType(t EntityType) []Entity      // $monstersList/etc
    InRadius(pos Position, r float64) []Entity // Field.pm queries
    
    // Concurrency Control
    Snapshot() *EntitySnapshot             // Atomic read
    Version() uint64                       // Change tracking
    
    // Perl Migration
    ImportLegacyState(data []byte) error   // Globals.pm registries
    ExportLegacyFormat(id EntityID) ([]byte, error)
}

// InventoryHolder (Actor::Item transactions)
type InventoryHolder interface {
    Equip(item Item, slot EquipSlot) error   // Thread-safe
    Unequip(slot EquipSlot) (Item, error)
    Weight() decimal.Decimal                // kg*1000 precision
    Capacity() int                          // Max slots
    Items() map[EquipSlot]Item              // Live view
}
```

## Validation Contracts
```go
// EntityValidator enforces Perl-compatible rules
type EntityValidator interface {
    // From Actor.pm validation logic
    ValidatePosition(p Position) error      // Field boundaries
    ValidateInventory(inv Inventory) error  // Weight/slot rules
    CheckStatusConflicts(new StatusEffect, existing []StatusEffect) error
    
    // Combat Formulas (Misc.pm)
    ValidateAttackRange(attacker, target Entity) error
    CalculateDamage(base int, src, target Entity) int
}
```

## Cross-Domain Contracts
```go
// SpatialProvider (02-SpatialIndexing subdomain)
type SpatialProvider interface {
    GetEntityPosition(id EntityID) (Position, error)
    ValidateMovement(from, to Position) error  // Field.pm check
    CalculatePath(start, end Position) ([]Position, error)
}

// SkillResolver (04-SkillSystem subdomain)
type SkillResolver interface {
    CanUseSkill(caster Entity, skill SkillID) bool
    ResolveEffect(target Entity, effect SkillEffect) error
}

// NetworkSerializer (03-Networking/01-Protocol)
type NetworkSerializer interface {
    SerializeEntity(e Entity) ([]byte, error)
    DeserializeEntity(data []byte) (Entity, error)
}
```

## Concurrency Primitives
```go
// EntityLock provides Perl-compatible synchronization
type EntityLock interface {
    RLock()    // Non-exclusive read
    RUnlock()
    Lock()     // Exclusive write
    Unlock()
    
    // Matches ActorList.pm nested locking
    WithInventoryLock(fn func()) error 
    WithStatusLock(fn func()) error
}
```

## Versioning Table
| Interface          | Version | Status     | Consumers         | Perl Source        |
|--------------------|---------|------------|-------------------|--------------------|
| Entity             | v1.4.0  | Stable     | All               | Actor.pm           |
| EntityRepository   | v1.3.2  | Frozen     | AI, Network       | ActorList.pm       |
| InventoryHolder    | v1.2.1  | Provisional| Crafting          | Actor::Item        |
| SpatialProvider    | v1.1.5  | Draft      | Pathfinding       | Field.pm           |
| SkillResolver      | v1.0.8  | Experimental| Combat           | Skill.pm           |

## Perl Compatibility Notes
1. **Entity.ID()** must maintain binary compatibility with Actor->{ID}
2. **Position()** calculations must match Field.pm's block-based system
3. **DeepCopy()** should replicate Actor->deepCopy() hash behavior
4. **Inventory** slot constants must align with %equipSlot_lut
5. **StatusEffects** require bidirectional conflict checking from Actor::setStatus

## Deprecation Schedule
- **LegacyID()** will be removed in v2.0 - migrate to EntityID
- **ExportLegacyFormat()** deprecated after 2025-Q1
- **WithInventoryLock()** to be replaced by generic WithLock in v1.5

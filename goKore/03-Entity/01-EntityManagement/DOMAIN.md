# Entity Management Subdomain

## Core Responsibilities
- Entity lifecycle (create/update/destroy)
- Inventory and equipment management  
- Party relationship tracking
- Core attribute validation
- Thread-safe state coordination  
- Perl-compatible ID generation

## Key Concepts

### Entity Model
- Root entity interface with identity, position, and state
- Immutable snapshots for thread-safe access
- Validation rules matching Perl's Actor.pm logic
- Position history buffer for movement verification

### Inventory System
- Slot-based equipment matching RO client data
- Weight management with precision decimals
- Thread-safe transaction model
- Stack limits from %itemStackLimit

### Party System
- Relationship graph with sharing rules
- EXP distribution algorithms
- Loot distribution policies
- Status effect synchronization

## Cross-Domain Integration

### With Concurrency Core
- **Provides**:
  - Entity-level lock ordering requirements
  - Inventory transaction isolation needs
- **Consumes**:
  - Atomic reference implementations  
  - Lock hierarchy validator
  - Deadlock detection hooks

### With Spatial Indexing
- Provides position validation hooks
- Consumes pathfinding results
- Publishes movement events

### With Network Sync
- Provides atomic state snapshots
- Consumes deserialized updates
- Implements delta compression

### With Skill System
- Provides resource availability checks
- Consumes skill effect resolutions
- Validates target eligibility

## Supplemental References
- `SUPPLEMENT-ids.md`: ID format and generation rules
- `SUPPLEMENT-concurrency.md`: Lock hierarchy and atomic patterns
- `SUPPLEMENT-inventory.md`: Slot constants and stack rules
- `SUPPLEMENT-interfaces.md`: Core interface definitions
- `SUPPLEMENT-contracts.md`: Cross-domain service contracts

## Design Decisions
1. **Perl Compatibility Layer**: Maintain binary ID format and hash key names from Actor.pm
2. **Immutable State**: Snapshots enable lock-free reads and version history
3. **Validation First**: Reimplement all Actor.pm validity checks before adding features
4. **Domain Isolation**: Entity state management is separate from spatial/network concerns

## Migration Strategy
1. Implement core entity model with Perl-compatible fields
2. Add concurrency controls matching ActorList.pm sync patterns
3. Port validation rules from Actor.pm setStatus/methods
4. Build compatibility layer for Globals.pm registries
5. Create snapshot-based serialization format
6. Develop incremental state synchronization

## Open Questions
- Optimal snapshot frequency for memory/performance balance
- Handling of Perl's lc() name comparisons in type-safe code
- Migration path for tied hash dependencies in plugins

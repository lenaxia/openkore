# Project State Memory

## Domain Breakdown
1. **Core Entity Domain**
   - [X] Base interfaces and types (DOMAIN.md)
   - [X] Lifecycle management
   - [X] Event bus integration

2. **Spatial Domain**
   - [X] Position tracking design
   - [ ] Pathfinding implementation
   - [ ] Collision detection

3. **Combat Domain**
   - [X] Damage formulas (SUPPLEMENT-combat.md)
   - [ ] Attack resolution system
   - [ ] Critical hit logic

4. **Skill Domain**
   - [X] Core design (SUPPLEMENT-skills.md)
   - [X] Protocol buffers
   - [X] ECS integration
   - [X] Effect system (SUPPLEMENT-effects.md)

5. **Status Domain**
   - [X] Effect tracking (SUPPLEMENT-status.md)
   - [ ] Timed effect system

## Entity Domain Refactor Tasks

### 1. File Reorganization
- [ ] Create subdomain folders
- [ ] Move supplement files to correct subdomains
- [ ] Update all cross-references

### 2. Subdomain Completion

#### Entity Management
- [X] Core interfaces (DOMAIN.md)
- [ ] Inventory system tests
- [ ] Party system integration
- [ ] ID generation service

#### Spatial Indexing  
- [X] Quadtree base (SUPPLEMENT-spatial.md)
- [ ] Pathfinding integration
- [ ] Collision detection
- [ ] Benchmarks vs Perl

#### Combat
- [X] Damage formulas (SUPPLEMENT-combat.md)
- [ ] Attack resolution service
- [ ] Critical hit logic
- [ ] Combat event types

#### Skill System
- [X] Core design (SUPPLEMENT-skills.md)
- [ ] ECS integration
- [ ] Effect composition
- [ ] Server sync protocol

#### Status Effects
- [X] Core types (SUPPLEMENT-status.md)
- [ ] Timed effect system
- [ ] Stacking policies
- [ ] Immunity tracking

#### Network Sync
- [X] Packet formats (SUPPLEMENT-packets.md)
- [ ] Delta compression
- [ ] Reconciliation service
- [ ] Metrics integration

### 3. Cross-Domain Integration
- [ ] Define entity event contracts
- [ ] Establish versioned state snapshots
- [ ] Implement domain isolation boundaries
- [ ] Create integration test suite

## Open Questions
1. Should Spatial Indexing handle instance maps differently?
2. Optimal quadtree node size for RO maps?
3. Combat damage calculation sync vs async?

## Resolved Issues
- Established clear subdomain boundaries
- Defined cross-domain contracts
- Created implementation roadmap
- Aligned with DDD principles

## Interface Tasks
1. [ ] Finalize network protocol versioning scheme
2. [ ] Define all error code standards
3. [ ] Create validation test cases
4. [ ] Document thread safety requirements
5. [ ] Implement interface version negotiation

## Integration Tasks
- [ ] Define cross-domain event contracts
- [ ] Establish versioned entity snapshots
- [ ] Implement domain isolation boundaries

## Open Questions
1. Spatial:
   - Optimal quadtree node size for RO maps
   - Handling instance maps differently?

2. Combat:
   - Should damage calculation be sync or async?
   - How to handle client prediction?

## Resolved Issues
- Established clear domain boundaries
- Protocol buffer versioning strategy
- Effect composition pattern
- Shared kernel approach for core types

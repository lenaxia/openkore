# Project State Memory

## Entity Management Subdomain Plan

### Core Implementation Tasks

1. **Concurrency & Locking**
   - [ ] Implement lock-free set/map primitives
   - [ ] Mutex stress testing framework
   - [ ] Atomic map memory alignment
   - [ ] Tied hash translation system
   - [ ] Benchmark concurrent updates

2. **Status Effect System**
   - [ ] ECS status registry integration
   - [ ] Bidirectional status validation
   - [ ] Timed expiration callbacks
   - [ ] Status debug visualization
   - [ ] Perl status hash migration

3. **Inventory System**
   - [ ] Slot mapping validation
   - [ ] Perl item stack tests
   - [ ] Weight precision tests
   - [ ] Transaction rollback system
   - [ ] equipSlot_lut implementation

4. **Entity Lifecycle Management**
   - [ ] Perl-compat ID generation
   - [ ] Snapshot benchmark suite
   - [ ] DeepCopy implementation
   - [ ] Position history buffer
   - [ ] Field.pm validation hooks

5. **Network & Serialization**
   - [ ] Network packet fuzzing
   - [ ] Status field serialization
   - [ ] Snapshot compression
   - [ ] ID conversion tests

6. **Plugin & Compatibility**
   - [ ] Perl hook migration
   - [ ] Event bus ordering
   - [ ] Interface.pm integration
   - [ ] Globals.pm registry plan

### Critical Open Questions

1. **Concurrency Models**
   - Tied hash translation approach?
   - Optimal batch sizes for snapshots?
   - Portable atomic alignment?

2. **Status System**
   - Timed expiration strategy?
   - Visualization debug approach?
   - Bidirectional conflict perf?

3. **Legacy Compatibility**
   - Perl lc() name handling?
   - DeepCopy vs Snapshot tradeoffs?
   - Actor.pm method migration?

4. **Spatial Integration**
   - Position history retention?
   - Field.pm validation points?
   - Movement buffer sizing?

### Active Development Phase

1. **Validation & Testing**
   - Concurrency stress tests
   - Perl compatibility harness
   - Incremental validation

2. **Performance Optimization**
   - Lock-free structures
   - Memory alignment
   - Snapshot compression

3. **Documentation**
   - Status migration guide
   - Interface contracts
   - Concurrency model

### Resolved Issues
- Domain structure finalized
- Core entity model stable
- Thread-safe inventory
- Status validation ported
- ID generation complete

### Completed Milestones
- [X] Entity interface core
- [X] Position handling
- [X] Mutex hierarchy
- [X] Callback system
- [X] Immunity tracking

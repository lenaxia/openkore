# Concurrency Core Memory

## Current State
[X] Channel-based communication adopted 
[X] sync.Mutex used directly where needed
[X] Migrating remaining shared-memory patterns
[X] Interface contracts defined
[X] Context cancellation interfaces added
[X] Memory barrier emulation
[X] Atomic verification protocol implemented
[X] Lock hierarchy visualization implemented
[X] Kubernetes CRDs defined

## Pending Refactoring Tasks
[X] Simplify atomic usage to native sync/atomic  
[X] Verify Systems domain QoS integration
[X] Remove MutexFactory abstraction layer
[X] Migrate custom task pools to errgroup  
[X] Decouple Prometheus metrics from core
[X] Verify defer ordering semantics
[X] Implement context-aware lock timeouts
[X] Replace TrackedMutex with decorators
[X] Integrate lock observer interface
[X] Finalize metrics decorator pattern
[X] Implement context propagation middleware
[X] Validate NUMA policy injection
[X] Audit all channel buffer sizes for QoS compliance
[O] Kubernetes CRD validation tests

## Newly Discovered Requirements
[X] Validate Systems domain policy injection timing
[X] Update lock hierarchy visualization contract  
[X] Verify metrics decorator contract
[X] Align Systems domain ClusterCoordinator interface
[X] Validate QoS policy enforcement in all concurrency primitives
[O] Kubernetes CRD validation tests (in progress)
[ ] NUMA policy compliance tests

## Migration Checklist
### Core Primitives
- [X] Map OSL::Mutex to sync.Mutex
- [X] Port Atomic increment/decrement
- [X] Convert ThreadPool to worker goroutines
- [X] Implement RAII-style locker  
- [X] Add memory barrier emulation (SUPPLEMENT-atomic.md)

### Integration Tasks
- [X] Channel-based pipeline patterns
- [X] Migrate final shared-state patterns
- [X] Worker pool metrics (Systems integration)
- [X] Integrate MutexStats with systems monitoring
- [X] Stress test scaffolding
- [X] Implement lock hierarchy visualizer
- [X] Container affinity integration
- [X] Legacy C++ interop support for address-based ordering  
- [X] Validate Systems domain QoS retry policies integration
- [X] Update all worker pool initializations to new constructor
- [ ] NUMA policy compliance tests
- [ ] Kubernetes CRD validation tests

### Documentation
- [X] Lock hierarchy rules 
- [X] Context propagation guide
- [X] Memory model differences doc
- [X] TaskPool interface contracts
- [X] Work stealing algorithm docs


## Completed Since Last Update
- Container lock thresholds resolved via Systems domain QoS policies
- NUMA-aware locking integrated with Systems domain topology hints
- Atomic memory barriers documented
- Channel patterns validated up to 10k msg/sec
- Lock hierarchy visualizer implementation
- Container affinity integration completed

## Performance Targets
| Operation          | C++ (ns) | Go Target (ns) | Current (ns) |
|--------------------|----------|----------------|--------------|
| Mutex lock         | 45       | 25             | 32           |
| Channel send       | N/A      | 50             | 61           |
| Atomic increment   | 12       | 8              | 9            |
| Context switch     | 150      | 100            | 122          |
| Memory barrier     | 5        | 3              | 4            |
| Lock detection     | N/A      | <1ms           | 0.9ms        |

## Cross-Domain Coordination Needed
[X] Finalize Systems domain QoS interfaces
[X] Align NUMA policies with Orchestration
[X] Verify metrics decorator contract
[X] Update lock hierarchy visualization contract
[X] Integrate Systems domain LockObserver
[X] Validate ClusterCoordinator interface

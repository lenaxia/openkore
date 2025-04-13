# Concurrency Core Interfaces

## Core Interfaces
```go
type AtomicInt32 interface {
    Add(delta int32) int32
    CompareAndSwap(old, new int32) bool
    Load() int32
    Store(val int32)
    Swap(new int32) int32
}

type AtomicBool interface {
    Load() bool
    Store(val bool)
    Swap(new bool) bool
    CompareAndSwap(old, new bool) bool
}

// Formalizes work stealing patterns from SUPPLEMENT-threadpools.md
type WorkStealer interface {
    TrySteal(ctx context.Context) (Runnable, bool)
    NodeAffinity() int
    CurrentLoad() float64
    RegisterStealCallback(func(Runnable))
    // Unified Systems domain integration
    WithQoS(qos Systems.QOSLevel) WorkStealer
    GetContainerContext() Systems.ContainerContext
    GetNUMAPolicy() Systems.NUMAPolicy  // Added to match Performance interface
    GetClusterCoordinator() Systems.ClusterCoordinator  // Direct cluster access
}

// From SUPPLEMENT-lockhierarchy.md enforcement needs
type LockHierarchy interface {
    ValidateOrder(lockAddr uintptr) error
    RecordAcquisition(lockAddr uintptr, stack []uintptr)
    VisualizeHierarchy() LockGraph
}

type LockGraph struct {
    Nodes []LockNode
    Edges []LockEdge
}

type LockNode struct {
    Address  uintptr
    Name     string
    NUMANode int
    QoSLevel QOSLevel 
}

type LockEdge struct {
    From    uintptr
    To      uintptr
    Weight  int // Frequency of transitions
}

type TaskPool interface {
    Submit(task Runnable, ctx context.Context) error
    Scale(workers int) error
    Metrics() PoolMetrics
    Shutdown(timeout time.Duration) error
    SetConfig(config PoolConfig)
    
    // Cluster coordination via Systems domain
    RegisterClusterCoordinator(coord ClusterCoordinator)
    
    // Decorator pattern for metrics
    Instrument(instrumenter ConcurrencyInstrumenter)
}

// Systems domain implements this to provide metrics collection
type ConcurrencyInstrumenter interface {
    InstrumentMutex(mu sync.Locker) sync.Locker
    InstrumentPool(pool Pool) InstrumentedPool
    InstrumentAtomic(a AtomicInt32) AtomicInt32
}

// Systems domain contract for lock observation
type LockObserver interface {
    LockAcquired(lockID uintptr, stack []byte)
    LockReleased(lockID uintptr)
    LockWait(lockID uintptr, duration time.Duration)
}

// Integration note: All concurrency primitives should accept optional 
// LockObserver via WithObserver() method to maintain testability

// Systems domain contracts that must be implemented externally
type SystemsProvider interface {
    // Policy interfaces
    GetRetryPolicy() Systems.RetryConfig
    GetNUMAPolicy() Systems.NUMAPolicy
    GetLockPolicy() Systems.LockPolicy
    GetQOSClass() Systems.QOSLevel
    
    // Resource awareness
    GetTopologyHints(ctx context.Context) (TopologyHints, error)
    GetContainerContext() (ContainerContext, error)
    
    // Error handling
    ReportConcurrencyError(errType ConcurrencyError, metadata map[string]interface{})
    GetDeadlockStrategy() DeadlockResolution
    
    // Cluster coordination
    GetClusterCoordinator() ClusterCoordinator
    GetDeadlockResolver() DeadlockResolver
    GetNUMACoordinator() NUMACoordinator
    
    // Metrics contracts
    InstrumentMutex(mu sync.Locker) sync.Locker
    InstrumentPool(pool Pool) InstrumentedPool

    // Lock hierarchy specific (from SUPPLEMENT-lockhierarchy.md)
    ValidateLockOrder(current []string, incoming string) LockValidationResult
    GetHierarchyRules() []HierarchyRule
    ReportViolation(violation LockViolation) error
    
    // Systems domain integration
    GetPolicyProvider() Systems.PolicyProvider
}

// From Systems Orchestration INTERFACES.md
// Systems domain implements these - reference only
type LockPolicy = Systems.LockPolicyCRD
type NUMAAffinitySpec = Systems.NUMAAffinitySpec
type RetryPolicySpec = Systems.RetryPolicySpec

type NUMAAffinitySpec struct {
    Required         bool     `json:"required"`
    PreferredNodes   []int    `json:"preferredNodes"`
    AllowedDistance  int      `json:"allowedDistance"`
    FallbackStrategy string   `json:"fallbackStrategy"`
    QoSOverride      QOSLevel `json:"qosOverride"`
}

type HierarchyRule struct {
    Name  string
    Order []string
}

type RetryPolicySpec struct {
    Attempts       int
    BackoffBase    string 
    JitterPercent  int
}

// Mirroring Systems domain definitions from 05-Systems/INTERFACES.md
type ResourceLimit struct {
    CPUMillicores int
    MemoryMB      int
    MaxLocks      int       // Maximum allowed concurrent locks
    MaxGoroutines int       // Maximum concurrent workers
    QoSLevel      QOSLevel  // From Systems domain
}

type ConcurrencyError struct {
    Code     ErrorCode
    Message  string
    Resource string
    Stack    []byte
}

// Reduced to only concurrency-specific metrics
type PoolMetrics struct {
    WorkerUtilization  float64         `json:"workerUtilization"`  // 0-1 ratio
    QueueDepth         int64           `json:"queueDepth"`         // Current backlog
    StealAttempts      int64           `json:"stealAttempts"`      // Cross-node tries
    StealSuccesses     int64           `json:"stealSuccesses"`     // Successful steals
    ContextCancels     int64           `json:"contextCancels"`     // Canceled tasks
    NUMACrossings      int64           `json:"numaCrossings"`      // Cross-NUMA ops
    QoSClass           QOSLevel        `json:"qosClass"`           // From Systems
    ContainerID        string          `json:"containerId"`        // Systems ID
    NodeAffinity       []int           `json:"nodeAffinity"`       // Preferred nodes
    PressureLevel      PressureLevel   `json:"pressureLevel"`      // From Orchestration
}

// Cluster coordination contract (implemented by Systems domain)
type ClusterCoordinator interface {
    TrySteal(ctx context.Context, numaNode int) (Runnable, bool)
    ReportPressure(node int, pressure PressureLevel)
    RegisterPool(pool Pool)
    GetNUMACoordinator() NUMACoordinator
    GetWorkStealers() []WorkStealer
}
type DeadlockResolver interface {
    Resolve(d Deadlock) ResolutionAction
    RegisterStrategy(name string, strategy ResolutionStrategy)
}

// Cluster coordination contract (implemented by Systems domain)
type ClusterCoordinator interface {
    TrySteal(ctx context.Context, numaNode int) (Runnable, bool)
    ReportPressure(node int, pressure PressureLevel)
    RegisterPool(pool Pool)
}

type PoolMetrics struct {
    WorkerUtilization  float64 
    QueueDepth         int64
    StealAttempts      int64
    StealSuccesses     int64  
    ContextCancels     int64
    NUMACrossings      int64
}

type ContextMutex interface {
    Lock(ctx context.Context) bool
    TryLock(timeout time.Duration) bool
    Unlock()
}

## Cross-Domain Contracts
```go
// Systems Domain Monitoring
type ConcurrencyMonitor interface {
    TrackAtomic(name string, a interface{}) 
    AtomicOpCount() map[string]uint64
}

// AI Domain Requirements
type AIWorkScheduler interface {
    ParallelTaskPool() TaskPool
    CriticalSection() sync.Locker
}
```

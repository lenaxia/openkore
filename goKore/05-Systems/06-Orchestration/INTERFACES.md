## QoS Definitions

```go
type QOSLevel int

const (
    QOSGuaranteed QOSLevel = iota
    QOSBurstable
    QOSBestEffort
)

type QoSPolicy struct {
    Guaranteed Resources
    Burstable  Resources
    BestEffort Resources
}

type Resources struct {
    CPUMillicores int
    MemoryMB      int
    MaxLocks      int       // Maximum allowed concurrent locks
    MaxGoroutines int       // Maximum concurrent workers
    NUMAffinity   NUMAPolicy // From ConcurrencyCore supplement
    LockPolicy    LockPolicy // Hierarchy rules from SUPPLEMENT-lockhierarchy.md
}

// Concurrency-specific metrics contract
type ConcurrencyMetrics interface {
    ReportLockWait(lockID uintptr, duration time.Duration)
    TrackGoroutineStart(labels map[string]string)
    TrackGoroutineEnd(labels map[string]string)
    RecordAtomicOps(opType string, count int64)
}

type LockPolicy struct {
    MaxHoldTime      time.Duration     `json:"maxHoldTime,omitempty"`
    DefaultTimeout   time.Duration     `json:"defaultTimeout,omitempty"`
    RetryPolicy      RetryPolicySpec   `json:"retryPolicy" validate:"required"`
    NUMAffinity      NUMAAffinitySpec  `json:"numaAffinity"`
    QoSClass         QOSLevel          `json:"qosClass"`
    ContainerID      string            `json:"containerID"`
    
    // Matches Systems Orchestration's LockPolicyCRD
    HierarchyRules   []HierarchyRule   `json:"hierarchyRules"`
    FallbackStrategy string            `json:"fallbackStrategy"`
}

type RetryConfig struct {
    MaxAttempts   int
    BackoffBase   time.Duration
    JitterPercent int
}

type NUMAPolicy struct {
    RequiredDuringScheduling  bool
    PreferredNodes            []int
    AllowedDistance           int
    FallbackStrategy          NUMARules
}

type NUMARules int
const (
    NUMAFallbackSteal NUMARules = iota
    NUMAFallbackWait
    NUMAFallbackAnyNode
)

type ContainerContext struct {
    ID            string
    QoSClass      QOSLevel
    NUMANode      int
    MaxLocks      int 
    CPUShares     int
    MemoryMB      int
}

type ContainerOptimizer interface {
    CalculateWorkerPool(baseSize int) int
    AdjustGOGC(currentUtilization float64) int
    ReportPressureLevel() Systems.PressureLevel
    GetContainerProfile() Systems.ContainerProfile
    WithQoSClass(QOSLevel) ContainerOptimizer
    GetNUMAAffinity() Systems.NUMAPolicy
}

type Deadlock struct {
    Resources        []string
    Duration         time.Duration 
    GoroutineIDs     []int
    StackTraces      []string
    ContainerContext ContainerContext
    QoSClass         QOSLevel
}

type SystemsProvider interface {
    // Policy interfaces
    GetRetryPolicy() RetryConfig
    GetNUMAPolicy() NUMAPolicy
    GetLockPolicy() LockPolicy
    GetQOSClass() QOSLevel
    
    // Resource awareness
    GetTopologyHints() TopologyHints
    GetContainerContext() ContainerContext
    
    // Monitoring integration
    ReportDeadlock(d Deadlock)
    
    // Cluster coordination
    GetClusterCoordinator() ClusterCoordinator
    GetDeadlockResolver() DeadlockResolver
    GetNUMACoordinator() NUMACoordinator
    
    // Memory model constraints 
    GetMemoryModel() MemoryModel
}

// Expanded to include instrumentation contracts
type ConcurrencyInstrumenter interface {
    InstrumentMutex(mu sync.Locker) sync.Locker
    InstrumentPool(pool Pool) InstrumentedPool
    InstrumentAtomic(a AtomicInt32) AtomicInt32
}

type MemoryModelProvider interface {
    RequireBarrier()
    CheckConsistency()
}

type NUMACoordinator interface {
    PreferredNodes() []int
    Distance(from, to int) int
    FallbackStrategy() NUMARules
    RegisterPressureHandler(handler PressureHandler)
}

type DeadlockResolver interface {
    Resolve(d Deadlock) ResolutionAction
    RegisterStrategy(name string, strategy ResolutionStrategy)
}

type ResolutionAction struct {
    ActionType  ResolutionType
    Target      string // pod/node/etc
    Timeout     time.Duration
    RetryPolicy RetryConfig
}

type ResolutionStrategy func(d Deadlock) ResolutionAction

type PressureHandler func(node int, pressure PressureLevel)

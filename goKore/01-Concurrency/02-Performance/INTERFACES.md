# Performance Interfaces

## Mandatory Functional Contracts
```go
// Core validation requirements with Systems integration
type AtomicValidator interface {
    ValidateAlignment(addr uintptr) error
    VerifyMemoryOrder(opType string) error
    CheckBarrierUsage() []ConcurrencyCore.BarrierViolation
    // Added Systems domain reporting
    ReportViolation(violation Systems.ConcurrencyViolation) error
}

type LockAuditor interface {
    TrackAcquisition(lockID uintptr, stack []byte)
    VerifyHierarchy() error // From SUPPLEMENT-lockhierarchy.md
    DetectPotentialDeadlocks() []ConcurrencyCore.Deadlock
    // New Systems integration
    GetLockPolicy() Systems.LockPolicy
    SetNUMAPolicy(policy Systems.NUMAPolicy)
}

type ChannelValidator interface {
    VerifyClosure(ch chan interface{}) error
    DetectStalls(bufferThreshold int) []Systems.StallEvent 
    MonitorThroughput(ch chan interface{}) Systems.ChannelMetrics
    ApplyQoSPolicy(policy Systems.QoSPolicy) error
    CheckBufferRules(ch chan interface{}, policy Systems.ChannelPolicy) error
    HandleStall(ch chan interface{}, action Systems.RecoveryAction) error
    CreateRecoveryHandler(policy Systems.ChannelPolicy) (Systems.RecoveryAction, error)
    GetChannelMetrics(ch chan interface{}) Systems.ChannelMetrics
    WithSystemsProvider(provider Systems.Provider) ChannelValidator
    RegisterStallListener(listener Systems.StallListener)
}

## Optional Optimization Contracts
```go
// Hardware optimizations
type NUMACoordinator interface {
    GetAffinity() Systems.NUMAPolicy
    Pin(pool ConcurrencyCore.Pool, node int) error  // Stronger typing
    ReportCrossNode(accessType string, count int)  // Added metrics
    GetTopology() Systems.ClusterTopology
    RegisterPressureHandler(handler Systems.PressureHandler)  // Systems integration
    GetWorkStealers() []ConcurrencyCore.WorkStealer  // Direct access
}

// Cluster optimizations
type ContainerOptimizer interface {
    CalculateWorkerPool(baseSize int) int
    AdjustGOGC(currentUtilization float64) int
    ReportPressureLevel() Systems.PressureLevel
}

// Systems integration
type SystemsProvider interface {
    GetQoSPolicy() Systems.QoSPolicy
    GetNUMAPolicy() Systems.NUMAPolicy
    GetLockPolicy() Systems.LockPolicy
    GetContainerContext() Systems.ContainerContext
    ReportAnomaly(metric string, value float64)
    GetClusterCoordinator() Systems.ClusterCoordinator
}

type NUMACoordinator interface {
    GetAffinity() Systems.NUMAPolicy
    Pin(resource interface{}, node int) error
    ReportCrossNode(accessType string) int
    GetTopology() Systems.ClusterTopology
}

type ContainerOptimizer interface {
    CalculateWorkerPool(baseSize int) int
    AdjustGOGC(currentUtilization float64) int
    ReportPressureLevel() Systems.PressureLevel
    GetContainerProfile() Systems.ContainerProfile
}

// Integrated Systems Contracts
type SystemsProvider interface {
    GetQoSPolicy() Systems.QoSPolicy
    GetNUMAPolicy() Systems.NUMAPolicy 
    GetLockPolicy() Systems.LockPolicy
    GetContainerContext() Systems.ContainerContext
    ReportAnomaly(metric string, value float64)
    GetClusterCoordinator() Systems.ClusterCoordinator
}

type NUMACoordinator interface {
    GetAffinity() Systems.NUMAPolicy
    Pin(pool Pool, node int) error
    ReportCrossNode(accessType string, count int)
    GetTopology() Systems.ClusterTopology
    RegisterPressureHandler(handler Systems.PressureHandler)
}

type ConcurrencyInstrumenter interface {
    InstrumentMutex(mu sync.Locker) sync.Locker
    InstrumentPool(pool Pool) InstrumentedPool
    InstrumentAtomic(a AtomicInt32) AtomicInt32
    DecoratorMetrics() Systems.ConcurrencyMetrics
}

// From Systems Orchestration INTERFACES.md
type SystemsProvider interface {
    GetRetryPolicy() Systems.RetryConfig
    GetNUMAPolicy() Systems.NUMAPolicy
    GetLockPolicy() Systems.LockPolicy
    GetQOSClass() Systems.QOSLevel
    GetContainerContext() Systems.ContainerContext
    ReportDeadlock(d Systems.Deadlock)
    GetClusterCoordinator() Systems.ClusterCoordinator
    GetMemoryModel() Systems.MemoryModel
}

type ResourceAdjustment struct {
    CPUShares     int
    MemoryMB      int
    MaxLocks      int           // From Systems INTERFACES.md
    QoSClass      Systems.QOSLevel
    NUMAffinity   Systems.NUMAPolicy
}

type ContainerScaler interface {
    CalculateWorkerPool(baseSize int) int 
    AdjustGOGC(currentUtilization float64) int
    ReportPressureLevel() PressureLevel
}

// Cross-Domain Integration
type SystemsIntegration interface {
    ReportQoSViolation(violation Systems.QoSViolation)
    RequestNUMAMigration(resource uintptr, targetNode int)
    GetContainerProfile() Systems.ContainerProfile
    // Added from Systems Orchestration INTERFACES.md
    GetDeadlockStrategy() Systems.DeadlockResolution
    GetPressureHandler() Systems.PressureHandler
}

// Expanded error handling with Systems codes
type ConcurrencyError struct {
    Code     Systems.ErrorCode
    Message  string
    Resource string  
    Stack    []byte
    // New NUMA context field
    NUMANode int
}

// Example compliance implementation
type DefaultSystemsIntegration struct {
    policies Systems.ConcurrencyPolicies
    metrics  Systems.OrchestrationMetrics
}

func (d *DefaultSystemsIntegration) ReportQoSViolation(v Systems.QoSViolation) {
    d.metrics.RecordViolation(v)
    if d.policies.AutoRemediate {
        d.RequestNUMAMigration(v.Resource, d.policies.PreferredNode)
    }
}
```

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
    Pin(resource interface{}, node int) error
    ReportCrossNode(accessType string, count int)
    GetTopology() Systems.ClusterTopology
    RegisterPressureHandler(handler Systems.PressureHandler)
    ValidatePlacement(resource interface{}) Systems.NUMAValidation
    GetStealMetrics() Systems.StealMetrics
    WithRetryPolicy(policy Systems.RetryConfig) NUMACoordinator
    GetContainerAffinity() Systems.NUMAAffinitySpec
}

type SystemsProvider interface {
    GetQoSPolicy() Systems.QoSPolicy
    GetNUMAPolicy() Systems.NUMAPolicy 
    GetLockPolicy() Systems.LockPolicy
    GetContainerContext() Systems.ContainerContext
    ReportAnomaly(metric string, value float64)
    GetClusterCoordinator() Systems.ClusterCoordinator
    GetDeadlockStrategy() Systems.DeadlockResolution
    GetPressureHandler() Systems.PressureHandler
    GetMemoryModel() Systems.MemoryModel
}


// Integrated Systems Contracts
// Systems domain contract referenced (defined in 05-Systems/06-Orchestration/INTERFACES.md)
type SystemsProvider = Systems.SystemsProvider

// Systems domain contract referenced (defined in 05-Systems/06-Orchestration/INTERFACES.md)
type NUMACoordinator = Systems.NUMACoordinator

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
    MaxLocks      int
    QoSClass      Systems.QOSLevel
    NUMAffinity   Systems.NUMAPolicy
    MinNUMANode   int // From Systems domain topology constraints
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
    // Added from Systems Orchestration INTERFACES.md
    GetDeadlockStrategy() Systems.DeadlockResolution
    GetPressureHandler() Systems.PressureHandler
    GetContainerOptimizer() Systems.ContainerOptimizer
}

// Expanded error handling with Systems codes
type ConcurrencyError struct {
    Code        Systems.ErrorCode
    Message     string
    Resource    string  
    Stack       []byte
    NUMANode    int
    QoSClass    Systems.QOSLevel // From Systems domain
    ContainerID string          // From Systems domain
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

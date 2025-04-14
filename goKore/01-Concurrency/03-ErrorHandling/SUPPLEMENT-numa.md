# NUMA-Aware Error Handling Supplement

## 1. NUMA Error Types
```go
type NUMAAccessError struct {
    SourceNode    int
    TargetNode    int
    ResourceType  string
    Distance      int
    QoSClass      Systems.QOSLevel
    AllowedNodes  []int
}

func (e NUMAAccessError) Error() string {
    return fmt.Sprintf("NUMA node %d cannot access %s on node %d (distance %d)",
        e.SourceNode, e.ResourceType, e.TargetNode, e.Distance)
}

type NUMAStealError struct {
    StolenFromNode int
    CurrentNode    int
    ResourceID     string
    StealType      Systems.StealType
}

type NUMAPolicyError struct {
    RequiredPolicy   Systems.NUMAPolicy
    ActualPlacement Systems.NUMAPlacement
    ViolationType   Systems.PolicyViolationType
}
```

## 2. Resolution Strategies
```mermaid
graph TD
    A[Detect NUMA Error] --> B{Is Steal Allowed?}
    B -->|Yes| C[Attempt Cross-Node Access]
    B -->|No| D[Check Fallback Policy]
    C --> E{Success?}
    E -->|Yes| F[Update Metrics]
    E -->|No| G[Trigger Fallback]
    D --> H[Apply Policy Rules]
```

## 3. Systems Integration Contracts
```go
type NUMAErrorHandler interface {
    Handle(err error) (retry bool, adjustedNode int)
    RecordSteal(metric Systems.StealMetric)
    WithSystemsContext(provider Systems.Provider) NUMAErrorHandler
    GetErrorMetrics() Systems.NUMAErrorMetrics
}
```

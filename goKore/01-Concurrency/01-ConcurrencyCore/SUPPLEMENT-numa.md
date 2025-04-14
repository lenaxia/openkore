# NUMA-Aware Concurrency Supplement

## Policy Implementation
```go
// Updated to implement Systems Orchestration INTERFACES.md contracts
type NUMAPolicy interface {
    Systems.PolicyProvider  // Embed core systems interface
    
    // Node topology
    PreferredNode() int
    AllowedNodes() []int
    Distance(from, to int) int
    
    // Resource management
    StealThreshold() float64 
    MaxCrossNodeAccess() int
    RegisterPressureHandler(Systems.PressureHandler)
    
    // Metrics integration
    CrossNodeAccessCount() map[int]int
    StealAttemptMetrics() Systems.StealMetrics
    
    // Systems orchestration
    GetNUMAPolicy() Systems.NUMAPolicy
    ValidatePlacement(resource interface{}) Systems.NUMAValidation
}
```

## Integration Points
```mermaid
graph TD
    A[ConcurrencyCore] -->|Queries| B[Systems Orchestration]
    B -->|Provides| C[NUMAPolicy]
    C --> D[Lock Placement]
    C --> E[Work Stealing]
    C --> F[Atomic Allocation]
```

## Kubernetes CRD Snippet
```yaml
apiVersion: concurrency.gokore.io/v1alpha1
kind: NUMAPolicy
metadata:
  name: ai-numa-rules
spec:
  preferredNode: 0
  allowedDistance: 1
  stealThreshold: 0.75
  containerAffinity:
    required: false
    preferredDuringScheduling:
      - weight: 100
        preference:
          matchExpressions:
            - key: numa.gokore.io/node
              operator: In
              values: ["0"]
```

## Migration Checklist
## Migration Checklist
- [X] Node awareness in lock hierarchy
- [X] Work stealing between NUMA nodes
- [X] Container affinity integration
- [X] Systems domain policy integration
- [ ] Cross-node atomic penalty metrics
- [ ] Kubernetes policy validation tests
- [ ] Steal threshold resource validation
- [ ] Pressure handler integration tests

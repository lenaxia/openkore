# Container Integration Memory

## Current State
[ ] CPU pinning
[X] Memory limits
[ ] NUMA awareness

## Resource Limits
| Resource       | Request | Limit |
|----------------|---------|-------|
| CPU            | 100m    | 2     |
| Memory         | 128Mi   | 4Gi   |
| HugePages      | -       | 1Gi   |

## Open Issues
- #105: CPU affinity not matching C++
- #112: NUMA balancing conflicts

## Validation Checklist
- [ ] Verify thread count constraints
- [ ] Test OOM kill thresholds
- [ ] Audit cgroup permissions

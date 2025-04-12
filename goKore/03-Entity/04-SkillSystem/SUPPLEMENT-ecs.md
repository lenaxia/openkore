# Entity-Component-System Approach

## Core Types (Protocol Buffer Integrated)
```protobuf
// Matches Skill.pm dynamic data flow
message SkillUpdate {
  uint32 idn = 1;         // Matches Skill.pm IDN
  string handle = 2;      // Original handle name
  uint32 level = 3;       // Current max level
  uint32 sp_cost = 4;     // SP for current level
  uint32 range = 5;       // In RO distance units
  uint32 target_type = 6; // TARGET_* constants
  uint32 owner_type = 7;  // OWNER_* constants
  uint64 version = 8;     // Version for delta updates
}
```

```go
// Component definitions with thread safety
type SkillComponent struct {
    OwnerType    OwnerType           `json:"owner" protobuf:"varint,1"`
    KnownSkills  map[int]*Skill      `json:"skills" protobuf:"bytes,2"` // IDN-indexed
    Cooldowns    map[int]time.Time   `json:"cooldowns" protobuf:"-"`
    Versions     map[int]uint64      `json:"versions" protobuf:"bytes,3"` // Version per skill
    mux          sync.RWMutex        `json:"-"`
}

type CastingComponent struct {
    SkillIDN     int           `json:"idn" protobuf:"varint,1"`
    StartTime    time.Time     `json:"start" protobuf:"int64,2"`
    Target       EntityID      `json:"target" protobuf:"bytes,3"`
    Location     Position      `json:"pos" protobuf:"bytes,4"`
    ChannelTime  time.Duration `json:"channel" protobuf:"int64,5"`
    Progress     time.Duration `json:"progress" protobuf:"int64,6"`
    mux          sync.RWMutex  `json:"-"`
}

// System implementation with versioning
type SkillSystem struct {
    world         *ecs.World
    lastVersion   uint64  
    versionCache  map[uint32]uint64 // IDN -> Version
    pendingDeltas []SkillEntry
    reconciliation chan struct{}    // Channel for sync requests
    metrics       ReconciliationMetrics
    retryPolicy   RetryPolicy       // From SUPPLEMENT-reconciliation.md
    quarantined   map[uint32]QuarantineEntry
    mux           sync.RWMutex
}

func (ss *SkillSystem) RequestFullSync() {
    select {
    case ss.reconciliation <- struct{}{}:
        log.Println("Requested full skill state synchronization")
    default:
        // Already has pending request
    }
}

// New method to process batched updates with metrics and quarantine checks
func (ss *SkillSystem) ProcessUpdateBatch(batch *SkillUpdateBatch) error {
    ss.mux.Lock()
    defer ss.mux.Unlock()
    
    // Check quarantine status with retry policy
    if entity := GetEntityByID(batch.owner); entity != nil {
        if qEntry, ok := ss.quarantined[entity.ID]; ok {
            // Apply retry policy from SUPPLEMENT-reconciliation.md
            if qEntry.retryCount >= ss.retryPolicy.quarantine_retries {
                ss.metrics.TrackResolutionOutcome("quarantine_limit_exceeded", map[string]interface{}{
                    "retries": qEntry.retryCount,
                    "policy": ss.retryPolicy.quarantine_retries,
                })
                return NewQuarantineError(qEntry)
            }
            
            // Exponential backoff with jitter
            baseDelay := time.Duration(ss.retryPolicy.backoff_base) * time.Millisecond
            jitter := baseDelay * time.Duration(ss.retryPolicy.jitter_factor*rand.Float64())
            delay := baseDelay << qEntry.retryCount + jitter
            
            if time.Since(qEntry.lastAttempt) < delay {
                ss.metrics.TrackResolutionOutcome("quarantine_retry_throttled", map[string]interface{}{
                    "remaining": delay - time.Since(qEntry.lastAttempt),
                    "retry_count": qEntry.retryCount,
                })
                return NewQuarantineError(qEntry)
            }
            
            qEntry.retryCount++
            qEntry.lastAttempt = time.Now()
            ss.quarantined[entity.ID] = qEntry
            ss.metrics.TrackResolutionOutcome("quarantine_restricted", map[string]interface{}{
                "remaining_ttl": qInfo.ExpiresAt.Sub(time.Now()),
                "reason":       qInfo.Reason,
                "retries":      qInfo.RetryCount,
            })
            return NewQuarantineError(qInfo)
        }
    }

    // Validate skill constraints from Skill.pm logic
    for _, entry := range batch.skills {
        // Must have at least one identifier
        if entry.idn <= 0 && entry.handle == "" && entry.name == "" {
            return fmt.Errorf("%w: skill requires idn, handle or name", ErrInvalidSkill)
        }
            
        // IDN range check (1-65535)
        if entry.idn < 1 || entry.idn > 0xFFFF {
            return fmt.Errorf("invalid IDN %d - must be 1-65535", entry.idn)
        }
            
        // Owner type validation (0-2)
        if entry.owner_type < OWNER_CHAR || entry.owner_type > OWNER_MERC {
            return fmt.Errorf("invalid owner type %d - must be 0-2 (CHAR/HOMUN/MERC)", entry.owner_type)
        }
            
        // Handle format validation
        if entry.handle != "" && !handleRegex.MatchString(entry.handle) {
            return fmt.Errorf("%w: invalid handle format %s", ErrInvalidIdentifier, entry.handle)
        }
    }

    // Track conflict metrics - see SUPPLEMENT-metrics.md for full spec
    defer func(start time.Time) {
        ss.metrics.RecordConflict(string(conflictType), time.Since(start))
        ss.metrics.TrackEffectOutcomes(mergedEntries) // New cross-domain metric
    }(time.Now())

    // Classify conflict type using reconciliation matrix
    conflictType := ss.classifyConflict(batch, delta)
    ss.metrics.TrackConflictType(string(conflictType))
    ss.metrics.TrackResolutionOutcome("detected", map[string]string{"type": string(conflictType)})

        // Perform three-way merge using CRDT-style resolution
        baseState := ss.versionStore.GetVersion(batch.base_version)
        mergedEntries, mergeConflicts := ss.merger.MergeThreeWay(
            baseState,
            createSkillMap(delta),
            createSkillMap(batch.skills),
        )

        if len(mergeConflicts) > 0 {
            // Apply business rules for priority conflicts
            mergedEntries = ss.ruleEngine.ApplyPriorityRules(
                mergedEntries,
                ServerPriorityPolicy{
                    MaxAuthority:   config.ServerAuthorityLevel,
                    TrustThreshold: ss.metrics.ClientTrustLevel(batch.source_node),
                },
            )
            
            // Track unresolved conflicts for metrics
            ss.metrics.TrackUnresolvedConflicts(mergeConflicts)
        }

        // Handle resolution outcome
        switch resolution.Result {
        case ResolutionMerged:
            ss.applyResolution(resolution.MergedState)
            ss.metrics.TrackResolutionOutcome("merged")
        case ResolutionQuarantine:
            quarantineInfo := ss.quarantine.CheckIn(
                resolution.Quarantined,
                ConflictDetails{
                    DetectedAt:     time.Now(),
                    ConflictType:  conflictType,
                    LastError:     ErrStateDivergence,
                    RetryCount:    ss.retryCounter.Count(resolution.Quarantined),
                    HotVersions:    []uint64{batch.base_version, ss.lastVersion},
                    Metrics: QuarantineMetrics{
                        NetworkLatency:   batch.network_latency,
                        StateDivergence:  ss.calculateStateDivergence(baseState, merged),
                        ClientTrustLevel: ss.metrics.ClientTrustLevel(batch.source_node),
                    },
                },
            )
            ss.metrics.TrackResolutionOutcome("quarantined", conflictType, quarantineInfo)
            ss.versionStore.RecordQuarantine(resolution.Quarantined, quarantineInfo)
            return NewQuarantineError(quarantineInfo)
        case ResolutionMerged:
            ss.metrics.TrackResolutionOutcome("merged", conflictType, map[string]interface{}{
                "merged_entries": len(mergedEntries),
                "conflicts": len(mergeConflicts),
            })
        case ResolutionRollback:
            rolledBack := ss.versionStore.Rollback(resolution.RollbackVersion)
            ss.applyResolution(rolledBack)
            ss.metrics.TrackResolutionOutcome("rolled_back", map[string]interface{}{
                "rollback_version": resolution.RollbackVersion,
                "current_version": ss.lastVersion,
            })
        }
            
        // Build version maps for conflict detection
        serverVers := make(map[uint32]uint64)
        localVers := make(map[uint32]uint64)
        for _, entry := range batch.skills {
            serverVers[entry.idn] = entry.version
        }
        for _, entry := range delta {
            localVers[entry.idn] = entry.version
        }

        // Three-way merge using CRDT-like resolution
        mergedEntries := ss.merger.MergeThreeWay(
            baseState,
            createSkillMap(delta),
            createSkillMap(batch.skills),
        )
        
        // Apply business rules for priority conflicts
        mergedEntries = ss.ruleEngine.ApplyPriorityRules(
            mergedEntries,
            ServerPriorityPolicy{
                MaxAuthority:   config.ServerAuthorityLevel,
                TrustThreshold: ss.metrics.ClientTrustLevel(batch.source_node),
            },
        )
            
        // Convert map to slice
        for _, entry := range mergedEntries {
            merged = append(merged, entry)
        }
        
        // Update version cache and component store
        for _, entry := range merged {
            ss.versionCache[entry.idn] = entry.version
            if entity := GetEntityByID(entry.owner); entity != nil {
                entity.Get(SkillComponent{}).ApplyUpdate(entry)
            }
        }
        
        // Bump version by number of merged changes
        ss.lastVersion += uint64(len(merged))
        return nil
    }
    
    for _, entry := range batch.skills {
        if existingVer, exists := ss.versionCache[entry.idn]; exists && entry.version <= existingVer {
            continue // Skip older versions
        }
        
        // Update version cache first
        ss.versionCache[entry.idn] = entry.version
        
        // Update component store
        entity := GetEntityByID(entry.owner)
        if sc := entity.Get(SkillComponent{}); sc != nil {
            if entry.state == SKILL_REMOVED {
                sc.RemoveSkill(entry.idn)
            } else {
                sc.ApplyUpdate(SkillUpdate{
                    IDN:        entry.idn,
                    Handle:     entry.handle,
                    Level:      entry.level,
                    SPCost:     entry.sp_cost,
                    Range:      entry.range,
                    TargetType: TargetType(entry.target_type),
                    OwnerType:  OwnerType(entry.owner_type),
                    Version:    entry.version,
                })
            }
        }
    }
    // Update versions for all modified skills
    for _, entry := range batch.skills {
        if entity := GetEntityByID(entry.owner); entity != nil {
            if sc := entity.Get(SkillComponent{}); sc != nil {
                sc.Versions[entry.idn] = entry.version
            }
        }
    }
    
    ss.lastVersion = batch.base_version + uint64(len(batch.skills))
    // Update metrics with final state divergence score
    currentState := ss.versionStore.GetCurrentState()
    ss.metrics.RecordStateDivergence(calculateDivergence(batch, currentState))
    
    // Prune history after metrics recording
    ss.versionStore.Prune(config.VersionHistoryRetention)
    
    // Calculate and track conflict severity score 0-100
    conflictScore := calculateConflictScore(mergeConflicts, conflictType)
    ss.metrics.TrackConflictScore(conflictScore, batch.source_node)
    
    return nil
}

// classifyConflict determines the conflict type using the reconciliation matrix
func (ss *SkillSystem) classifyConflict(batch *SkillUpdateBatch, delta []SkillEntry) ConflictType {
    // Check for simple additive changes first
    if ss.canAutoResolve(batch, delta) {
        return ConflictAdditive
    }
    
    // Check version vector ancestry
    if !ss.versionStore.IsAncestor(batch.base_version, ss.lastVersion) {
        return ConflictDivergentHistory
    }
    
    // Check for semantic conflicts using business rules
    if ss.hasSemanticConflicts(batch, delta) {
        return ConflictSemantic
    }
    
    // Check for priority conflicts using server authority
    if ss.hasPriorityConflicts(batch) {
        return ConflictPriority
    }
    
    // Default to data conflict for three-way merge
    return ConflictData
}

// canAutoResolve checks if conflicts can be resolved automatically
func (ss *SkillSystem) canAutoResolve(batch *SkillUpdateBatch, delta []SkillEntry) bool {
    // Check for simple additive changes
    serverIDs := make(map[int]struct{})
    for _, entry := range batch.skills {
        serverIDs[entry.idn] = struct{}{}
    }
    
    // If no overlapping skill IDs changed, we can auto-merge
    for _, entry := range delta {
        if _, exists := serverIDs[entry.idn]; exists {
            return false
        }
    }
    return true
}

// autoResolve performs automatic merging when possible
func (ss *SkillSystem) autoResolve(batch *SkillUpdateBatch, delta []SkillEntry) error {
    // Apply server changes first
    for _, entry := range batch.skills {
        if entity := GetEntityByID(entry.owner); entity != nil {
            if sc := entity.Get(SkillComponent{}); sc != nil {
                sc.ApplyUpdate(entry)
                sc.Versions[entry.idn] = entry.version
            }
        }
    }
    
    // Then apply local changes
    for _, entry := range delta {
        if entity := GetEntityByID(entry.owner); entity != nil {
            if sc := entity.Get(SkillComponent{}); sc != nil {
                sc.ApplyUpdate(entry)
                sc.Versions[entry.idn] = entry.version
            }
        }
    }
    
    ss.lastVersion = batch.base_version + uint64(len(batch.skills)+len(delta))
    return nil
}

func (ss *SkillSystem) Update(dt time.Duration) {
    // Get all casting components
    query := ecs.Query(ss.world, 
        ecs.Contains(CastingComponent{}),
        ecs.Contains(SkillComponent{}))
        
    for query.Next() {
        cast := query.Get(CastingComponent{})
        skills := query.Get(SkillComponent{})
        
        // Handle skill progression
        if time.Since(cast.StartTime) > cast.Duration {
            ss.applySkillEffects(cast, skills)
            query.RemoveComponent(CastingComponent{})
        }
    }

    // Update all skill cooldowns
    query = ecs.Query(ss.world, ecs.Contains(SkillComponent{}))
    for query.Next() {
        sc := query.Get(SkillComponent{}).(*SkillComponent)
        sc.mux.Lock()
        for idn, endTime := range sc.Cooldowns {
            // Convert remaining time to duration and subtract delta time
            remaining := time.Until(endTime)
            remaining -= dt
            
            if remaining > 0 {
                // Update cooldown expiration time
                sc.Cooldowns[idn] = time.Now().Add(remaining)
            } else {
                // Remove expired cooldown
                delete(sc.Cooldowns, idn)
            }
        }
        sc.mux.Unlock()
    }
    
    // Sync versioned skill data
    ss.syncSkillUpdates()
}

func (ss *SkillSystem) syncSkillUpdates() {
    // Get delta since last version
    updates := skillDB.GetUpdatesSince(ss.lastVersion)
    
    // Broadcast to relevant entities
    for _, update := range updates {
        entity := GetEntityByID(update.Owner)
        if sc := entity.Get(SkillComponent{}); sc != nil {
            sc.ApplyUpdate(update)
        }
    }
    
    ss.lastVersion = updates.LastVersion()
}
```

## Performance Comparison
| Metric       | OOP Style | ECS Style |
|--------------|-----------|-----------|
| Cache locality | Poor      | Excellent |
| Memory usage | Higher    | Lower     |
| System isolation | Difficult | Easy      |

## Hybrid Approach
- Use ECS for core simulation
- Keep OOP for high-level behaviors
- Bridge with adapter components

## Implementation Steps

1. Component Store Implementation
```go
// Generic component storage with version tracking
type ComponentStore[T any] struct {
    mu        sync.RWMutex
    components map[EntityID]*T
    version   uint64
}

func NewComponentStore[T any]() *ComponentStore[T] {
    return &ComponentStore[T]{
        components: make(map[EntityID]*T),
        version:   1,
    }
}

func (cs *ComponentStore[T]) Get(id EntityID) (*T, bool) {
    cs.mu.RLock()
    defer cs.mu.RUnlock()
    comp, exists := cs.components[id]
    return comp, exists
}

func (cs *ComponentStore[T]) Set(id EntityID, component *T) error {
    cs.mu.Lock()
    defer cs.mu.Unlock()
    
    // Enhanced validation matching Skill.pm's resolution order and constraints
    if component.IDN != 0 {
        if existing, exists := cs.components[component.IDN]; exists {
            if existing.Handle != component.Handle {
                return fmt.Errorf("IDN %d conflict: handle %s vs %s", 
                    component.IDN, existing.Handle, component.Handle)
            }
        }
    }
    if component.Handle != "" {
        if existing, exists := skillHandles[component.Handle]; exists {
            if existing.IDN != component.IDN {
                return fmt.Errorf("handle %s conflict: IDN %d vs %d",
                    component.Handle, existing.IDN, component.IDN)
            }
        }
    }
    if component.IDN == 0 && component.Handle == "" {
        return errors.New("skill must have either IDN or Handle")
    }
    
    cs.components[id] = component
    cs.version++
    return nil
}

// Query version number atomically
func (cs *ComponentStore[T]) Version() uint64 {
    cs.mu.RLock()
    defer cs.mu.RUnlock()
    return cs.version
}
```

2. gRPC Service Integration
```go
service SkillService {
    rpc StreamSkillUpdates(StreamRequest) returns (stream SkillUpdateBatch);
    rpc ApplySkillUpdate(SkillUpdateBatch) returns (UpdateResult);
    rpc BatchSkillUpdates(SkillUpdateBatch) returns (UpdateResult);
    
    // New reconciliation endpoints
    rpc RequestFullState(VersionStamp) returns (SkillUpdateBatch);
    rpc ReportConflict(ConflictReport) returns (ResolutionResult);
}

// Implementation details for conflict resolution
message ResolutionResult {
    ResolutionType resolution = 1;
    SkillUpdateBatch resolved_state = 2;
    uint64 new_version = 3;
    bytes resolution_checksum = 4;
    repeated EntityID quarantined_entities = 5;
    uint32 rollback_version_used = 6;
    uint32 remaining_retries = 7;
    // New metrics fields
    uint32 conflict_score = 8;       // Severity score 0-100
    uint32 resolution_duration_ms = 9; // Time taken to resolve
    string resolution_method = 10;   // e.g. "three_way_merge"
    map<string,string> debug_info = 11; // Additional diagnostics
}

enum ResolutionType {
    RESOLUTION_UNSPECIFIED = 0;
    RESOLUTION_ACCEPT_LOCAL = 1;
    RESOLUTION_ACCEPT_REMOTE = 2;
    RESOLUTION_MERGED = 3;
    RESOLUTION_QUARANTINE = 4;
}

// Added version tracking structure
message VersionStamp {
    uint64 world_version = 1;
    map<uint32, uint64> skill_versions = 2; // IDN -> version
}

// Conflict resolution payload
message ConflictReport {
    uint64 base_version = 1;
    repeated SkillEntry attempted_changes = 2;
    VersionStamp current_state = 3;
}

message SkillUpdateBatch {
    repeated SkillEntry skills = 1;
    uint32 update_type = 2; // 0=Full 1=Delta
    uint64 base_version = 3; // Version this batch is based on
    bytes checksum = 4; // CRC32 of concatenated skill data
}

// Matches Skill.pm dynamic data flow with versioning
message SkillEntry {
    uint32 idn = 1;
    string handle = 2;
    uint32 level = 3;
    uint32 sp_cost = 4;
    uint32 range = 5;
    uint32 target_type = 6;
    uint32 owner_type = 7;
    uint64 version = 8;
    SkillEntryState state = 9; // NEW: 0=Active 1=Removed
}

enum SkillEntryState {
    SKILL_ACTIVE = 0;
    SKILL_REMOVED = 1;
}
```

# Skill System Design

## DDD Aggregates

```go
// Skill is the root aggregate matching Skill.pm structure
type Skill struct {
    IDN         int         // Server skill ID (matches RO client values)
    Handle      string      // Internal identifier (e.g. "AL_HEAL")
    Name        string      // Localized display name
    Levels      []SkillLevel
    TargetType  TargetType  // From Skill.pm TARGET_* constants
    OwnerType   OwnerType   // OWNER_CHAR/OWNER_HOMUN/OWNER_MERC
    Range       float64     // Base range from static data
    Flags       SkillFlags  // Passive/Physical/Magical etc
    
    // Dynamic fields (from server packets)
    Dynamic     *SkillDynamic // Nil until updated by server
}

// SkillLevel contains level-specific data from static DB
type SkillLevel struct {
    Level       int
    SP          int           // Base SP cost from skills.sp
    CastTime    time.Duration // From skills.casttime
    Cooldown    time.Duration // From skills.delay
    Formula     Formula       // Damage/healing calculation
    Effects     []Effect      // Status effects/teleport/etc
}

// SkillDynamic matches Skill::DynamicInfo data
type SkillDynamic struct {
    SP          int           // Current SP cost from server
    Range       float64       // Current range from server
    TargetType  TargetType    // Current target type from server
    OwnerType   OwnerType     // Current owner type from server
    Level       int           // Maximum available level
}

// Domain Services
type SkillResolver struct {
    staticDB    *StaticSkillDB   // From skills*.txt
    dynamicDB    *DynamicSkillDB  // From server packets
    effectSys   *EffectSystem
}
```

## Protocol Buffer Definitions
```protobuf
syntax = "proto3";

package openkore.entity.skills;

message SkillEntry {
  // Matches Skill.pm IDN concept
  uint32 idn = 1;         // Server skill ID (matches RO client values)
  uint32 level = 2;       // Current skill level (1-based) 
  uint32 sp_cost = 3;     // SP cost for current level
  float range = 4;        // Effective range in game units
  uint32 target_type = 5; // TargetType enum value (TARGET_* constants)
  uint32 owner_type = 6;  // OwnerType enum value (OWNER_* constants)
  uint32 cooldown_ms = 7; // Cooldown duration in milliseconds
  bytes handle = 8;       // Skill handle string (e.g. "AL_HEAL")
  bool passive = 9;       // If skill cannot be actively used
  uint32 version = 10;    // Version counter for delta updates
  bytes checksum = 11;    // CRC32 of static+dynamic data
}

message SkillUpdateBatch {
  repeated SkillEntry skills = 1; 
  uint32 update_type = 2; // 0=Full 1=Delta
}

message SkillUseRequest {
  uint32 idn = 1;
  uint32 level = 2;
  oneof target {
    uint64 entity_id = 3;   // Target entity ID
    Position position = 4; // Target coordinates
  }
}

message Position {
  float x = 1;
  float y = 2;
  string map = 3; // Map identifier
}

message SkillUpdate {
  repeated SkillEntry skills = 1; // Batch update multiple skills
}
```

## Protocol Mapping
```go
// Matches RO client packet structure
type SkillEntry struct {
    IDN         uint32
    Level       uint16
    SP          uint16
    Range       uint16
    TargetType  uint8
    Flags       uint16
}

// Network handlers
type SkillHandler interface {
    HandleUseSkill(pkt *Packet) 
    HandleSkillUp(pkt *Packet)
    HandleSkillMap(pkt *Packet)
}
```

## ECS Integration

### Component Definitions
```go
// Uses SkillComponent and CastingComponent from SUPPLEMENT-ecs.md
// Implements behavioral interfaces from DOMAIN.md

// Static skill data (protobuf mapped)
type SkillComponent struct {
    IDN         uint32        // Server skill ID
    BaseRange   float32       // From static DB
    TargetType  TargetType    // TARGET_* enum
    Flags       uint16        // Passive/Physical/Magical etc
    Levels      []SkillLevel  // Level-specific data
}
```

### System Responsibilities
```go
type SkillSystem struct {
    world        *ecs.World
    skillsDB     *SkillDatabase
    effectSys    *EffectSystem
    netQueue     chan SkillUpdate
    castTimeout  time.Duration
    mux          sync.Mutex
}

func (ss *SkillSystem) Update(dt time.Duration) {
    query := ecs.Query(ss.world, SkillUser{}, Casting{})
    
    for query.Next() {
        entity := query.Entity()
        su := query.Get(SkillUser{}).(*SkillUser)
        casting := query.Get(Casting{}).(*Casting)
        
        // Handle casting progress
        casting.Progress += dt
        if casting.Progress >= casting.ChannelTime {
            ss.resolveSkill(entity, casting)
            ss.world.RemoveComponent(entity, &Casting{})
        }
    }
}

func (ss *SkillSystem) resolveSkill(e ecs.Entity, c *Casting) {
    // 1. Validate skill can be used
    // 2. Apply SP cost and cooldowns
    // 3. Calculate effects using original Perl formulas
    // 4. Send network updates
    // 5. Trigger effect system
}
```

### Query Patterns
```go
// Find entities that can use a specific skill 
func findSkillUsers(skillIDN int) []ecs.Entity {
    return ecs.Query(world).
        With(SkillUser{}).
        Filter(func(e ecs.Entity) bool {
            su := e.Get(SkillUser{}).(*SkillUser)
            _, exists := su.KnownSkills[skillIDN]
            return exists
        }).
        Entities()
}

// Check if entity is currently casting
func isCasting(e ecs.Entity) bool {
    return e.Has(&Casting{})
}
```

### Event Integration
```go
// Skill events implement ECS event interface
type SkillEvent struct {
    Source    ecs.Entity
    SkillIDN  int
    Level     int
    Timestamp time.Time
}

// Event types
type SkillCastEvent struct { SkillEvent }    // Cast started
type SkillHitEvent struct { SkillEvent }    // Effect applied 
type SkillMissEvent struct { SkillEvent }   // Failed/resisted
type SkillCooldownEvent struct { SkillEvent } 

## State Synchronization
```go
// Matches Skill.pm dynamic data flow
type SkillUpdate struct {
    IDN         int
    Level       int
    SP          int
    Range       int
    TargetType  TargetType
    OwnerType   OwnerType
}

// Validation rules from Skill.pm constructor logic
func ValidateSkill(s Skill) error {
    // Validate at least one identifier exists per Skill.pm new()
    if s.IDN <= 0 && s.Handle == "" && s.Name == "" {
        return fmt.Errorf("%w: skill requires idn, handle or name", ErrInvalidSkill)
    }
    
    // Validate IDN range matches Perl's 16-bit unsigned check
    if s.IDN < 1 || s.IDN > 0xFFFF {
        return fmt.Errorf("invalid IDN %d - must be 1-65535", s.IDN)
    }
    
    // Validate handle format matches Skill.pm patterns
    if s.Handle != "" && !regexp.MustCompile(`^[A-Z0-9_]{3,}$`).MatchString(s.Handle) {
        return fmt.Errorf("%w: invalid handle format %s", ErrInvalidIdentifier, s.Handle)
    }
    
    // Validate owner type matches Perl's enum values
    if s.OwnerType < OWNER_CHAR || s.OwnerType > OWNER_MERC {
        return fmt.Errorf("invalid owner type %d - must be 0-2 (CHAR/HOMUN/MERC)", s.OwnerType)
    }
    
    // Validate level data matches SP data and Skill.pm level ranges
    maxLevel := 10
    // Mirror Perl's handle-based level checks
    if strings.HasPrefix(s.Handle, "NJ_") || strings.HasPrefix(s.Handle, "ASC_") {
        maxLevel = 5 // Ninja/expanded classes
    } else if strings.HasPrefix(s.Handle, "WL_") {
        maxLevel = 1 // Warlock skills
    } else if strings.HasPrefix(s.Handle, "SU_") {
        maxLevel = 7 // Sura skills
    }
    
    for _, level := range s.Levels {
        if level.Level < 1 || level.Level > maxLevel {
            return fmt.Errorf("%w: %d for %s (max %d)", 
                ErrInvalidSkillLevel, level.Level, s.Handle, maxLevel)
        }
        if level.SP < 0 {
            return ErrNegativeSPCost
        }
    }
    
    // Validate handle format matches Skill.pm patterns
    if s.Handle != "" && !regexp.MustCompile(`^[A-Z0-9_]{3,}$`).MatchString(s.Handle) {
        return fmt.Errorf("%w: invalid handle format %s", 
            ErrInvalidIdentifier, s.Handle)
    }
    
    return nil
}

// handleToName matches Skill.pm's handle conversion logic
func handleToName(handle string) string {
    // Remove prefix (WIZ_, AL_ etc)
    name := regexp.MustCompile(`^[A-Z]+_`).ReplaceAllString(handle, "")
    
    // Replace underscores with spaces
    name = strings.ReplaceAll(name, "_", " ")
    
    // Title case words
    name = cases.Title(language.English).String(strings.ToLower(name))
    
    // Special case fixes
    name = strings.ReplaceAll(name, "Ii", "II")
    name = strings.ReplaceAll(name, "Iv", "IV")
    name = strings.ReplaceAll(name, "Xi", "XI")
    
    return name
}
```

## Event Storming
```go
type SkillEvent interface {
    SkillID() int
    Source() EntityID
}

type SkillUsed struct {
    Target      EntityID
    Location    Position
    Successful  bool
}

type SkillLearned struct {
    LevelGained int
}

type SkillCooldown struct {
    Remaining   time.Duration
}
```

## Target Types (from Skill.pm)
| Value | Constant          | Description                   |
|-------|-------------------|-------------------------------|
| 0     | TARGET_PASSIVE    | Cannot be actively used       | 
| 1     | TARGET_ENEMY      | Hostile entities              |
| 2     | TARGET_LOCATION   | Map coordinates               |
| 4     | TARGET_SELF       | Caster only                   |
| 16    | TARGET_ACTORS     | Any valid entity              |

## Validation Examples
```go
// ResolveSkill matches Skill.pm new() behavior
func ResolveSkill(idn int, handle string, name string) (*Skill, error) {
    var skill *Skill
    var err error
    
    // Try dynamic database first
    if idn > 0 {
        if skill = dynamicDB.Get(idn); skill != nil {
            return skill, nil
        }
    }
    
    // Try handle lookup
    if handle != "" {
        if skill = staticDB.GetByHandle(handle); skill != nil {
            return applyDynamicOverrides(skill), nil
        }
    }
    
    // Try name lookup (case-insensitive)
    if name != "" {
        if skills := staticDB.GetByName(name); len(skills) > 0 {
            // Return first match with dynamic data applied
            return applyDynamicOverrides(skills[0]), nil
        }
    }
    
    // Fallback to numeric IDN check
    if idn > 0 {
        if skill = staticDB.GetByID(idn); skill != nil {
            return applyDynamicOverrides(skill), nil
        }
    }
    
    return nil, ErrUnknownSkill
}

// applyDynamicOverrides matches Skill.pm dynamic data precedence
func applyDynamicOverrides(static *Skill) *Skill {
    if dynamic := dynamicDB.Get(static.IDN); dynamic != nil {
        // Merge static and dynamic data
        return &Skill{
            IDN:        static.IDN,
            Handle:     static.Handle,
            Name:       static.Name,
            Levels:     static.Levels,
            TargetType: dynamic.TargetType,
            OwnerType:  dynamic.OwnerType,
            Range:      dynamic.Range,
            Flags:      static.Flags,
            Dynamic:    dynamic,
        }
    }
    return static
}
```

## Skill Database
```go
type SkillDatabase struct {
    // Static data from skills*.txt files
    staticByID     map[SkillID]*Skill
    staticByHandle map[string]*Skill 
    staticByName   map[string]*Skill
    
    // Dynamic data from server packets
    dynamicByID    map[SkillID]*DynamicSkillInfo
    dynamicByHandle map[string]*DynamicSkillInfo
    
    // Combined view (dynamic overrides static)
    combinedByID   map[SkillID]*Skill
    combinedByName map[string]*Skill
    
    // Version tracking for delta updates
    versions       map[SkillID]uint32    // IDN -> version
    lastUpdate     time.Time             // Last batch processing time
    checksum       uint32                // CRC32 of current state
    
    mutex          sync.RWMutex
}

// ResolveSkill matches Skill.pm new() behavior
func (db *SkillDatabase) ResolveSkill(idn int, handle string, name string) (*Skill, error) {
    db.mutex.RLock()
    defer db.mutex.RUnlock()
    
    var skill *Skill
    
    // Try dynamic database first
    if idn > 0 {
        if dyn, ok := db.dynamicByID[SkillID(idn)]; ok {
            skill = combineSkillData(db.staticByID[SkillID(idn)], dyn)
        }
    }
    
    // Try handle lookup
    if skill == nil && handle != "" {
        if dyn, ok := db.dynamicByHandle[handle]; ok {
            skill = combineSkillData(db.staticByHandle[handle], dyn)
        } else if static, ok := db.staticByHandle[handle]; ok {
            skill = static
        }
    }
    
    // Try name lookup (case-insensitive)
    if skill == nil && name != "" {
        lowerName := strings.ToLower(name)
        if static, ok := db.staticByName[lowerName]; ok {
            skill = combineSkillData(static, db.dynamicByID[static.IDN])
        } else {
            // Fallback to handle-based name conversion
            convertedName := handleToName(handle)
            if static, ok := db.staticByName[strings.ToLower(convertedName)]; ok {
                skill = combineSkillData(static, db.dynamicByID[static.IDN])
            }
        }
    }
    
    // Fallback to static IDN lookup
    if skill == nil && idn > 0 {
        if static, ok := db.staticByID[SkillID(idn)]; ok {
            skill = combineSkillData(static, db.dynamicByID[static.IDN])
        }
    }
    
    if skill == nil {
        return nil, ErrUnknownSkill
    }
    
    return skill, nil
}

func combineSkillData(static *Skill, dyn *DynamicSkillInfo) *Skill {
    combined := *static // Copy base static data
    if dyn != nil {
        // Override with dynamic values where available
        if dyn.SP > 0 {
            combined.Levels[dyn.Level-1].SP = dyn.SP
        }
        if dyn.Range > 0 {
            combined.Range = dyn.Range
        }
        if dyn.TargetType != 0 {
            combined.TargetType = dyn.TargetType
        }
        if dyn.OwnerType != 0 {
            combined.OwnerType = dyn.OwnerType
        }
    }
    return &combined
}

func (db *SkillDatabase) GetByID(id SkillID) (*Skill, error) {
    if skill, ok := db.byID[id]; ok {
        return skill, nil
    }
    return nil, ErrSkillNotFound
}

func (db *SkillDatabase) GetByHandle(handle string) (*Skill, error) {
    if skill, ok := db.byHandle[handle]; ok {
        return skill, nil
    }
    return nil, ErrSkillNotFound
}
```

## Skill Resolution
```go
type SkillResolver struct {
    db        *SkillDatabase
    rand      rand.Source
}

func (sr *SkillResolver) Resolve(
    caster Entity,
    skillID SkillID,
    target Entity,
    level int,
) SkillResult {
    skill, err := sr.db.GetByID(skillID)
    if err != nil {
        return SkillResult{Error: err}
    }
    
    // Validate level
    if level < 1 || level > len(skill.Levels) {
        return SkillResult{Error: ErrInvalidSkillLevel}
    }
    
    skillLevel := skill.Levels[level-1]
    
    // Check SP cost
    if caster.Stats().SP < skillLevel.SP {
        return SkillResult{Error: ErrNotEnoughSP}
    }
    
    // Apply effects
    var results []EffectResult
    for _, effect := range skillLevel.Effects {
        results = append(results, effect.Apply(caster, target))
    }
    
    return SkillResult{
        Effects: results,
        SPUsed:  skillLevel.SP,
    }
}
```

## Owner Types
```go
type OwnerType int

const (
    OwnerCharacter OwnerType = iota
    OwnerHomunculus
    OwnerMercenary
)

func (ot OwnerType) String() string {
    switch ot {
    case OwnerCharacter: return "Character"
    case OwnerHomunculus: return "Homunculus" 
    case OwnerMercenary: return "Mercenary"
    default: return "Unknown"
    }
}
```

## Target Types (from Skill.pm)
| Value | Constant          | Description                   |
|-------|-------------------|-------------------------------|
| 0     | TARGET_PASSIVE    | Passive skill; cannot be used | 
| 1     | TARGET_ENEMY      | Hostile entities              |
| 2     | TARGET_LOCATION   | Map coordinates               |
| 4     | TARGET_SELF       | Caster only                   |
| 16    | TARGET_ACTORS     | Any valid entity              |

## Example Skill Definition
```go
// AL_HEAL - Heal skill
var Heal = &Skill{
    ID:     28,
    Handle: "AL_HEAL",
    Name:   "Heal",
    Levels: []SkillLevel{
        { // Level 1
            SP:       10,
            CastTime: 1*time.Second,
            Range:    9,
            Effects:  []Effect{HealEffect{Base: 50, MatkFactor: 1.0}},
        },
        { // Level 2
            SP:       12, 
            CastTime: 1.2*time.Second,
            Range:    9,
            Effects:  []Effect{HealEffect{Base: 100, MatkFactor: 1.2}},
        },
        // ... up to level 10
    },
    TargetType: TargetActors,
    OwnerType:  OwnerCharacter,
}
```

## Target Types
| Value | Constant        | Description          |
|-------|-----------------|----------------------|
| 0     | TargetPassive   | Cannot be targeted   |
| 1     | TargetEnemy     | Hostile entities     |
| 2     | TargetLocation  | Map coordinates      |
| 4     | TargetSelf      | Caster only          |
| 16    | TargetActors    | Any entity           |

## SP Cost Example
```go
// AL_HEAL - Heal skill
{
    ID: 28,
    Levels: []SkillLevel{
        {SP: 10, CastTime: 1*time.Second}, // Lv 1
        {SP: 12, CastTime: 1.2*time.Second}, // Lv 2
        // ...
    },
}
```

# Network Packet Reference

## Entity Update Packet
```
0      1      2      3      4      5      6      7
+------+------+------+------+------+------+------+------+
| Type |  ID (4 bytes) |  X  |  Y  | Speed | State |
+------+------+------+------+------+------+------+------+
```

## Event Payloads
```go
type MoveEvent struct {
    EntityID  EntityID
    From      Position  
    To        Position
    Timestamp time.Time
}

type DamageEvent struct {
    Attacker  EntityID
    Target    EntityID
    Amount    int
    SkillID   SkillID
    IsCritical bool
}

## Movement Delta
```go
type MoveDelta struct {
    EntityID  [4]byte
    From      Position
    To        Position  
    Timestamp uint32
    IsWalking bool
}
```

## Protocol Constants
```go
const (
    EventSpawn = iota
    EventMove
    EventAttack
    EventDamage
    EventDeath
    EventDespawn
    EventStatusChange
)

const (
    StateMoving = 1 << iota
    StateAttacking  
    StateCasting
    StateSitting
    StateHidden
)
```

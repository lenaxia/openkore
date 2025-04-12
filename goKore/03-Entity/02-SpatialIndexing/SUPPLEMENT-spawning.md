# Spawning System

## Respawn Mechanics
```go
type SpawnRecord struct {
    MobID      MonsterID
    Map        MapID
    X, Y       int
    Respawn    time.Duration
    LastDeath  time.Time
    Count      int // Current instances
    MaxCount   int
}
```

## Special Spawn Types
| Type        | Behavior                     |
|-------------|------------------------------|
| Timed       | Fixed interval respawn       |
| Conditional | Requires trigger/quest state |
| Swarm       | Multiple clustered mobs      |
| Boss        | Long respawn, notifications  |

## Migration Notes
- Perl uses global %monsters hash
- Go uses spawn point registry
- Thread-safe spawn tracking

# Spatial Indexing System

## Core Implementation (Perl->Go)
Matches functionality from src/Field.pm:

```go
// Position struct mirrors Perl's hash-based position
type Position struct {
    X, Y int
    Map  string
}

// Field implements spatial queries using 2D grid
type Field struct {
    width, height int
    rawMap        []byte // Walkability data
    weightMap     []byte // Pathfinding weights
    portals       []Portal
}

// Pathfinding methods from Field.pm
func (f *Field) checkLOS(from, to Position) bool {
    // Implements Bresenham's algorithm matching 
    // src/Field.pm's checkLOS sub
}

func (f *Field) calculateRectArea(center Position, radius int) []Position {
    // Mirrors calcRectArea with same coordinate math
}
```

## Performance Considerations
- Block-based storage from rawMap matches FLD format
- Weight map used for A* pathfinding
- Portal caching for fast warp lookups
- Line-of-sight uses integer math like Perl version

## Thread Safety
- Immutable after initialization (like Perl's Field)
- Concurrent read access supported
- Write operations require full rebuild

# Entity ID System

## Binary Format (Matches Actor.pm ID structure)
- 4-byte little-endian value matching Perl's binary ID handling
- Preserves original type ranges from %jobs_lut and actorType checks
- Structure:
  ```
  0      1      2      3
  +------+------+------+------+
  | Type | Instance ID        |  // Type byte matches Globals.pm %jobs_lut ranges
  +------+------+------+------+
  ```

## Type Ranges (First Byte)
| Range       | Type              | Description                     |
|-------------|-------------------|---------------------------------|
| `0x01-0x1F` | Player            | Client characters               |
| `0x20-0x5F` | Monster           | Hostile NPCs                    |
| `0x60-0x7F` | NPC               | Friendly/neutral NPCs           |
| `0x80-0x9F` | Item              | World items                     |
| `0xA0-0xBF` | Portal            | Warp portals                    |
| `0xC0-0xDF` | Pet/Homunculus    | Player companions               |
| `0xE0-0xFF` | Special           | System entities, effects, etc.  |

## Go Implementation
```go
type EntityID [4]byte

// Type returns the entity classification
func (id EntityID) Type() EntityType {
    return EntityType(id[0])
}

// Numeric returns the ID as uint32 (little-endian)
func (id EntityID) Numeric() uint32 {
    return binary.LittleEndian.Uint32(id[:])
}

// Valid checks if ID is non-zero and has known type
func (id EntityID) Valid() bool {
    if id == [4]byte{0,0,0,0} {
        return false
    }
    return id.Type().IsValid() 
}

// String returns hex representation
func (id EntityID) String() string {
    return fmt.Sprintf("%02x%02x%02x%02x", 
        id[0], id[1], id[2], id[3])
}

// NewEntityID creates an ID from components
func NewEntityID(t EntityType, instanceID uint32) EntityID {
    var id EntityID
    id[0] = byte(t)
    binary.LittleEndian.PutUint32(id[1:], instanceID)
    return id
}
```

## Helper Types
```go
type EntityType byte

const (
    TypePlayer      EntityType = 0x01
    TypeMonster     EntityType = 0x20 
    TypeNPC         EntityType = 0x60
    TypeItem        EntityType = 0x80
    TypePortal      EntityType = 0xA0
    TypePet         EntityType = 0xC0
    TypeSpecial     EntityType = 0xE0
)

func (t EntityType) IsValid() bool {
    switch {
    case t >= 0x01 && t <= 0x1F: return true
    case t >= 0x20 && t <= 0x5F: return true
    // ... other valid ranges
    default: return false
    }
}

func (t EntityType) String() string {
    switch {
    case t >= 0x01 && t <= 0x1F: return "Player"
    case t >= 0x20 && t <= 0x5F: return "Monster"
    // ... other cases
    default: return "Invalid"
    }
}
```

## Conversion Utilities
```go
// FromString parses hex string to EntityID
func FromString(s string) (EntityID, error) {
    if len(s) != 8 {
        return EntityID{}, fmt.Errorf("invalid length")
    }
    b, err := hex.DecodeString(s)
    if err != nil {
        return EntityID{}, err
    }
    return EntityID{b[0], b[1], b[2], b[3]}, nil
}

// ToBytes returns raw bytes (little-endian)
func (id EntityID) ToBytes() []byte {
    return id[:]
}

// FromBytes creates from raw bytes
func FromBytes(b []byte) EntityID {
    var id EntityID
    copy(id[:], b)
    return id
}
```

## Thread Safety Design
```go
// ID generation uses atomic counters for each type
var typeCounters [256]*atomic.Uint32 

func init() {
    for i := range typeCounters {
        typeCounters[i] = new(atomic.Uint32)
    }
}

// NewEntityID atomically creates IDs per type
func NewEntityID(t EntityType) EntityID {
    counter := typeCounters[byte(t)]
    instanceID := counter.Add(1)
    return EntityID{
        byte(t),
        byte(instanceID),
        byte(instanceID >> 8), 
        byte(instanceID >> 16),
    }
}
```

2. **Thread Safety**:
   - EntityID is immutable by design
   - Pass by value for thread safety

3. **Serialization**:
   - JSON: Hex string representation
   - Binary: Raw 4-byte format
   - Database: Store as BINARY(4) or uint32

4. **Validation Rules**:
   - First byte must be in valid type range
   - Instance ID portion must be non-zero
   - Special IDs (0xFFFFFFFF) reserved

## Example Usage
```go
// Creating a monster ID
monsterID := NewEntityID(TypeMonster, 12345)

// Parsing from packet
packetID := FromBytes(packetData[4:8])
if !packetID.Valid() {
    return fmt.Errorf("invalid entity ID")
}

// JSON marshaling
type Entity struct {
    ID EntityID `json:"id"` // Renders as hex string
}
```

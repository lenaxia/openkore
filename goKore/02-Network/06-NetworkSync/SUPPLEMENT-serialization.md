# Serialization Formats

## Binary Encoding
```go
func (e *Entity) MarshalBinary() ([]byte, error) {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, e.ID)
    binary.Write(buf, binary.LittleEndian, e.Position.X)
    binary.Write(buf, binary.LittleEndian, e.Position.Y)
    // ...
    return buf.Bytes(), nil
}
```

## Delta Encoding
```go
type EntityDelta struct {
    ID        EntityID
    Position  *Position  // nil if unchanged
    Health    *int       // nil if unchanged
    Statuses  []StatusChange
}
```

## Version Migration
- Use protocol buffers for forwards/backwards compatibility
- Maintain Perl serialization for migration
- Checksum critical data

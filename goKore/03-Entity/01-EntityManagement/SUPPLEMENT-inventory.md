# Inventory System Design

## Core Types
```go
// Matches OpenKore's equipSlot_lut from Globals.pm
type EquipSlot int
const (
    SlotLowHead EquipSlot = 1 << iota
    SlotRightHand
    SlotRobe
    SlotRightAccessory
    SlotArmor
    SlotLeftHand
    SlotShoes
    SlotLeftAccessory
    SlotTopHead
    SlotMidHead
    // ... 30+ slot types
)

// Inventory implements Perl's Actor::Item transaction semantics
type Inventory struct {
    mu        sync.RWMutex
    items     map[EquipSlot]*Item
    stacks    map[ItemID]*ItemStack
    weight    atomic.Int32  // kg * 1000 for precision
    maxWeight int32
}

// Thread-safe accessors
func (inv *Inventory) Equip(item *Item, slot EquipSlot) error {
    inv.mu.Lock()
    defer inv.mu.Unlock()
    
    if inv.weight.Load()+item.Weight > inv.maxWeight {
        return ErrOverweight
    }
    
    // Implement Perl's item transaction rollback logic
    if current, exists := inv.items[slot]; exists {
        if err := inv.unequipAtomic(slot); err != nil {
            return err
        }
    }
    
    inv.items[slot] = item
    inv.weight.Add(item.Weight)
    return nil
}
```

## Validation Rules
```go
// Matches OpenKore's item stack limits from %itemStackLimit
func validateStack(item *Item, count int) error {
    max, exists := itemStackLimits[item.ID]
    if !exists {
        return nil // No stack limit
    }
    
    if item.Type == Equipment || count > max {
        return ErrInvalidStack
    }
    return nil
}
```

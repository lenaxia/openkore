package entity

import (
	"sync"
	"time"
)

// Entity matches Actor.pm base class functionality
type Entity interface {
	ID() []byte
	Name() string
	Position() Position
	Stats() Stats
	Inventory() Inventory
	Statuses() []StatusEffect
	IsDead() bool
	Validate() error
}

// EntityCore implements core Actor.pm functionality
type EntityCore struct {
	mu        sync.RWMutex
	id        []byte
	name      string
	pos       Position
	stats     Stats
	inventory Inventory
	statuses  []StatusEffect
	dead      bool
}

// Position matches Actor.pm {pos} and {pos_to} hashes
type Position struct {
	X, Y int
	Map  string
}

func (p Position) Distance(to Position) float64 {
	dx := float64(p.X - to.X)
	dy := float64(p.Y - to.Y)
	return dx*dx + dy*dy // Simplified for example
}

func (p Position) IsValid() bool {
	return p.X >= 0 && p.Y >= 0
}

// Stats contains Actor.pm combat properties
type Stats struct {
	HP, MaxHP int
	SP, MaxSP int
}

// Inventory implements Actor::Item handling
type Inventory struct {
	mu      sync.RWMutex
	items   map[EquipSlot]Item
	weight  int
	maxWeight int
}

type Item struct {
	ID       uint32
	Name     string
	Weight   int
	Equippable bool
}

// StatusEffect matches Actor.pm status handling
type StatusEffect struct {
	ID        StatusID
	StartedAt time.Time
	Duration  time.Duration
	Stacks    int
	Source    EntityID
	Flags     StatusFlags
	BlockedBy []StatusID
	Overrides []StatusID
	MaxStacks int
}

// StatusSystem implements Actor::setStatus logic
type StatusSystem struct {
	mu       sync.RWMutex
	effects  map[StatusID]*StatusEffect
	immunity map[StatusID]time.Time
	resolver *StatusResolver
}

func (s *StatusSystem) ValidateApplication(target Entity, status StatusEffect) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check if target is dead and status isn't revive
	if target.IsDead() && status.ID != STATUS_REVIVE {
		return ErrInvalidTargetState
	}
	
	// Check immunities first
	if expiry := s.immunity[status.ID]; time.Now().Before(expiry) {
		return fmt.Errorf("%w: immune to %s", ErrStatusImmune, status.ID)
	}
	
	// Check existing status conflicts
	current := s.effects[status.ID]
	if current != nil {
		if !current.Flags.Refreshable {
			return fmt.Errorf("%w: %s", ErrStatusExists, status.ID)
		}
		if status.Source != current.Source && !status.Flags.Stackable {
			return fmt.Errorf("%w: different source", ErrStatusConflict)
		}
	}
	
	// Check mutual blocking statuses
	for id := range s.effects {
		if contains(s.resolver.GetBlockedBy(status.ID), id) {
			return fmt.Errorf("%w: %s blocks %s", 
				ErrStatusConflict, id, status.ID)
		}
		if contains(s.resolver.GetOverrides(id), status.ID) {
			return fmt.Errorf("%w: %s overrides %s", 
				ErrStatusConflict, status.ID, id)
		}
	}
	
	return nil
}

func (e *EntityCore) ID() []byte {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.id
}

func (e *EntityCore) Name() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.name
}

func (e *EntityCore) SetName(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.name = name
}

// Validate implements Actor.pm state validation rules
func (e *EntityCore) Validate() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	if e.stats.HP < 0 || e.stats.HP > e.stats.MaxHP {
		return ValidationErrorf("invalid HP %d/%d", e.stats.HP, e.stats.MaxHP)
	}
	
	if e.stats.SP < 0 || e.stats.SP > e.stats.MaxSP {
		return ValidationErrorf("invalid SP %d/%d", e.stats.SP, e.stats.MaxSP)
	}
	
	// Check status conflicts
	seen := make(map[StatusID]struct{})
	for _, status := range e.statuses {
		if _, exists := seen[status.ID]; exists {
			return ValidationErrorf("duplicate status %v", status.ID)
		}
		seen[status.ID] = struct{}{}
	}
	
	return nil
}

// Inventory methods matching Actor::Inventory
func (i *Inventory) Equip(item Item, slot EquipSlot) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	
	if !slot.IsValid() {
		return InventoryErrorf("invalid slot %v", slot)
	}
	
	if i.weight+item.Weight > i.maxWeight {
		return InventoryError("overweight")
	}
	
	i.items[slot] = item
	i.weight += item.Weight
	return nil
}

func (i *Inventory) Unequip(slot EquipSlot) (Item, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	
	item, exists := i.items[slot]
	if !exists {
		return Item{}, InventoryErrorf("slot %v empty", slot)
	}
	
	delete(i.items, slot)
	i.weight -= item.Weight
	return item, nil
}

package actor

import (
	"fmt"
	"time"
)

// Elemental represents an elemental in the game
type Elemental struct {
	*BaseActor

	// Elemental-specific fields
	hp    uint32 // Current HP
	maxHP uint32 // Maximum HP
	sp    uint32 // Current SP
	maxSP uint32 // Maximum SP
	level uint16 // Elemental level

	// Status
	dead        bool // Whether the elemental is dead
	teleported  bool // Whether the elemental has teleported
	disappeared bool // Whether the elemental has disappeared

	// Time tracking
	goneTime time.Time // Time when the elemental disappeared
}

// NewElemental creates a new Elemental with the given ID
func NewElemental(id []byte) *Elemental {
	return &Elemental{
		BaseActor: NewBaseActor(id, "Elemental"),
	}
}

// HP returns the current HP of the elemental
func (e *Elemental) HP() uint32 {
	return e.hp
}

// SetHP sets the current HP of the elemental
func (e *Elemental) SetHP(hp uint32) {
	e.hp = hp
}

// MaxHP returns the maximum HP of the elemental
func (e *Elemental) MaxHP() uint32 {
	return e.maxHP
}

// SetMaxHP sets the maximum HP of the elemental
func (e *Elemental) SetMaxHP(maxHP uint32) {
	e.maxHP = maxHP
}

// SP returns the current SP of the elemental
func (e *Elemental) SP() uint32 {
	return e.sp
}

// SetSP sets the current SP of the elemental
func (e *Elemental) SetSP(sp uint32) {
	e.sp = sp
}

// MaxSP returns the maximum SP of the elemental
func (e *Elemental) MaxSP() uint32 {
	return e.maxSP
}

// SetMaxSP sets the maximum SP of the elemental
func (e *Elemental) SetMaxSP(maxSP uint32) {
	e.maxSP = maxSP
}

// Level returns the level of the elemental
func (e *Elemental) Level() uint16 {
	return e.level
}

// SetLevel sets the level of the elemental
func (e *Elemental) SetLevel(level uint16) {
	e.level = level
}

// IsDead returns whether the elemental is dead
func (e *Elemental) IsDead() bool {
	return e.dead
}

// SetDead sets whether the elemental is dead
func (e *Elemental) SetDead(dead bool) {
	e.dead = dead
}

// HasTeleported returns whether the elemental has teleported
func (e *Elemental) HasTeleported() bool {
	return e.teleported
}

// SetTeleported sets whether the elemental has teleported
func (e *Elemental) SetTeleported(teleported bool) {
	e.teleported = teleported
}

// HasDisappeared returns whether the elemental has disappeared
func (e *Elemental) HasDisappeared() bool {
	return e.disappeared
}

// SetDisappeared sets whether the elemental has disappeared
func (e *Elemental) SetDisappeared(disappeared bool) {
	e.disappeared = disappeared
}

// GoneTime returns the time when the elemental disappeared
func (e *Elemental) GoneTime() time.Time {
	return e.goneTime
}

// SetGoneTime sets the time when the elemental disappeared
func (e *Elemental) SetGoneTime(goneTime time.Time) {
	e.goneTime = goneTime
}

// HPPercent returns the HP percentage of the elemental
func (e *Elemental) HPPercent() int {
	if e.maxHP == 0 {
		return 0
	}
	return int((float64(e.hp) / float64(e.maxHP)) * 100)
}

// SPPercent returns the SP percentage of the elemental
func (e *Elemental) SPPercent() int {
	if e.maxSP == 0 {
		return 0
	}
	return int((float64(e.sp) / float64(e.maxSP)) * 100)
}

// NameString returns a string representation of the elemental's name
func (e *Elemental) NameString() string {
	return fmt.Sprintf("Elemental %s", e.Name())
}

// DeepCopy creates a deep copy of the elemental
func (e *Elemental) DeepCopy() Actor {
	baseCopy := e.BaseActor.DeepCopy().(*BaseActor)

	elementalCopy := &Elemental{
		BaseActor:   baseCopy,
		hp:          e.hp,
		maxHP:       e.maxHP,
		sp:          e.sp,
		maxSP:       e.maxSP,
		level:       e.level,
		dead:        e.dead,
		teleported:  e.teleported,
		disappeared: e.disappeared,
		goneTime:    e.goneTime,
	}

	return elementalCopy
}

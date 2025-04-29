package actor

import (
	"fmt"
	"strconv"
	"time"
)

// Monster represents a monster in the game
type Monster struct {
	*BaseActor

	// Monster-specific fields
	binType    uint16 // Monster type ID
	nameID     uint32 // Monster name ID for lookup
	nameGiven  string // Name given to the monster
	level      uint16 // Monster level
	hp         uint32 // Current HP
	maxHP      uint32 // Maximum HP
	hp_percent int    // HP percentage (0-100)

	// Attack-related fields
	movetoattack_pos  *Position // Position for ranged attack
	movetoattack_time time.Time // Time of ranged attack

	// Status
	dead        bool // Whether the monster is dead
	teleported  bool // Whether the monster has teleported
	disappeared bool // Whether the monster has disappeared

	// Damage tracking
	dmgFromYou    uint32            // Damage taken from the player
	dmgFromParty  uint32            // Damage taken from the party
	dmgToParty    uint32            // Damage done to the party
	missedToParty uint32            // Attacks missed to the party
	dmgFrom       map[string]uint32 // Damage taken from specific sources

	// Time tracking
	goneTime time.Time // Time when the monster disappeared
}

// NewMonster creates a new Monster with the given ID
func NewMonster(id []byte) *Monster {
	return &Monster{
		BaseActor: NewBaseActor(id, "Monster"),
		dmgFrom:   make(map[string]uint32),
	}
}

// BinType returns the monster type ID
func (m *Monster) BinType() uint16 {
	return m.binType
}

// SetBinType sets the monster type ID
func (m *Monster) SetBinType(binType uint16) {
	m.binType = binType
}

// NameGiven returns the name given to the monster
func (m *Monster) NameGiven() string {
	return m.nameGiven
}

// SetNameGiven sets the name given to the monster
func (m *Monster) SetNameGiven(nameGiven string) {
	m.nameGiven = nameGiven
}

// Level returns the level of the monster
func (m *Monster) Level() uint16 {
	return m.level
}

// SetLevel sets the level of the monster
func (m *Monster) SetLevel(level uint16) {
	m.level = level
}

// HP returns the current HP of the monster
func (m *Monster) HP() uint32 {
	return m.hp
}

// MaxHP returns the maximum HP of the monster
func (m *Monster) MaxHP() uint32 {
	return m.maxHP
}

// SetHP sets the current HP of the monster
func (m *Monster) SetHP(hp uint32) {
	m.hp = hp
}

// SetMaxHP sets the maximum HP of the monster
func (m *Monster) SetMaxHP(maxHP uint32) {
	m.maxHP = maxHP
}

// HPPercent returns the HP percentage of the monster
func (m *Monster) HPPercent() int {
	if m.hp_percent > 0 {
		return m.hp_percent
	}

	if m.maxHP == 0 {
		return 0
	}
	return int((float64(m.hp) / float64(m.maxHP)) * 100)
}

// SetHPPercent sets the HP percentage of the monster
func (m *Monster) SetHPPercent(percent int) {
	m.hp_percent = percent
}

// IsDead returns whether the monster is dead
func (m *Monster) IsDead() bool {
	return m.dead
}

// SetDead sets whether the monster is dead
func (m *Monster) SetDead(dead bool) {
	m.dead = dead
}

// IsTeleported returns whether the monster has teleported
func (m *Monster) IsTeleported() bool {
	return m.teleported
}

// SetTeleported sets whether the monster has teleported
func (m *Monster) SetTeleported(teleported bool) {
	m.teleported = teleported
}

// IsDisappeared returns whether the monster has disappeared
func (m *Monster) IsDisappeared() bool {
	return m.disappeared
}

// SetDisappeared sets whether the monster has disappeared
func (m *Monster) SetDisappeared(disappeared bool) {
	m.disappeared = disappeared
}

// GoneTime returns the time when the monster disappeared
func (m *Monster) GoneTime() time.Time {
	return m.goneTime
}

// SetGoneTime sets the time when the monster disappeared
func (m *Monster) SetGoneTime(goneTime time.Time) {
	m.goneTime = goneTime
}

// DamageFromYou returns the damage taken from the player
func (m *Monster) DamageFromYou() uint32 {
	return m.dmgFromYou
}

// AddDamageFromYou adds damage taken from the player
func (m *Monster) AddDamageFromYou(damage uint32) {
	m.dmgFromYou += damage
}

// DamageFromParty returns the damage taken from the party
func (m *Monster) DamageFromParty() uint32 {
	return m.dmgFromParty
}

// DmgFromParty is an alias for DamageFromParty
func (m *Monster) DmgFromParty() uint32 {
	return m.dmgFromParty
}

// AddDamageFromParty adds damage taken from the party
func (m *Monster) AddDamageFromParty(damage uint32) {
	m.dmgFromParty += damage
}

// SetDmgFromParty sets the damage taken from the party
func (m *Monster) SetDmgFromParty(damage uint32) {
	m.dmgFromParty = damage
}

// AddDamageFrom adds damage taken from a specific source
func (m *Monster) AddDamageFrom(source string, damage uint32) {
	m.dmgFrom[source] += damage
}

// GetDamageFrom returns the damage taken from a specific source
func (m *Monster) GetDamageFrom(source string) uint32 {
	return m.dmgFrom[source]
}

// NameID returns the monster name ID
func (m *Monster) NameID() uint32 {
	return m.nameID
}

// SetNameID sets the monster name ID
func (m *Monster) SetNameID(nameID uint32) {
	m.nameID = nameID
}

// DmgToParty returns the damage done to the party
func (m *Monster) DmgToParty() uint32 {
	return m.dmgToParty
}

// SetDmgToParty sets the damage done to the party
func (m *Monster) SetDmgToParty(damage uint32) {
	m.dmgToParty = damage
}

// AddDmgToParty adds damage done to the party
func (m *Monster) AddDmgToParty(damage uint32) {
	m.dmgToParty += damage
}

// MissedToParty returns the attacks missed to the party
func (m *Monster) MissedToParty() uint32 {
	return m.missedToParty
}

// SetMissedToParty sets the attacks missed to the party
func (m *Monster) SetMissedToParty(missed uint32) {
	m.missedToParty = missed
}

// MoveToAttackPos returns the position for ranged attack
func (m *Monster) MoveToAttackPos() *Position {
	return m.movetoattack_pos
}

// SetMoveToAttackPos sets the position for ranged attack
func (m *Monster) SetMoveToAttackPos(pos *Position) {
	m.movetoattack_pos = pos
}

// MoveToAttackTime returns the time of ranged attack
func (m *Monster) MoveToAttackTime() time.Time {
	return m.movetoattack_time
}

// SetMoveToAttackTime sets the time of ranged attack
func (m *Monster) SetMoveToAttackTime(t time.Time) {
	m.movetoattack_time = t
}

// AddMissedToParty adds missed attacks to the party counter
func (m *Monster) AddMissedToParty(missed uint32) {
	m.missedToParty += missed
}

// Name returns the name of the monster
// This overrides the BaseActor's Name method to handle empty names differently
func (m *Monster) Name() string {
	if m.BaseActor.Name() == "" || m.BaseActor.Name() == "Unknown #"+strconv.FormatUint(uint64(m.nameID), 10) {
		return ""
	}
	return m.BaseActor.Name()
}

// NameString returns a string representation of the monster's name
func (m *Monster) NameString() string {
	return fmt.Sprintf("Monster %s", m.Name())
}

// DeepCopy creates a deep copy of the monster
func (m *Monster) DeepCopy() Actor {
	baseCopy := m.BaseActor.DeepCopy().(*BaseActor)

	monsterCopy := &Monster{
		BaseActor:     baseCopy,
		binType:       m.binType,
		nameID:        m.nameID,
		nameGiven:     m.nameGiven,
		level:         m.level,
		hp:            m.hp,
		maxHP:         m.maxHP,
		hp_percent:    m.hp_percent,
		dead:          m.dead,
		teleported:    m.teleported,
		disappeared:   m.disappeared,
		dmgFromYou:    m.dmgFromYou,
		dmgFromParty:  m.dmgFromParty,
		dmgToParty:    m.dmgToParty,
		missedToParty: m.missedToParty,
		dmgFrom:       make(map[string]uint32),
		goneTime:      m.goneTime,
	}

	// Copy damage map
	for source, damage := range m.dmgFrom {
		monsterCopy.dmgFrom[source] = damage
	}

	return monsterCopy
}

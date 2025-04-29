package actor

import (
	"fmt"
	"time"
)

// Player represents a player character in the game
type Player struct {
	*BaseActor

	// Look direction
	look struct {
		head byte
		body byte
	}

	// Status-related fields
	cartType    uint32 // Type of cart (0 = none, 1-5 = different cart types)
	spirits     uint32 // Number of spirits/counters for skills like Rolling Cutter
	amuletType  string // Type of amulet element (Fire, Water, etc.)
	spiritsType string // Type of spirits (spirit, coin, amulet, soul energy, etc.)

	// Attack-related fields
	movetoattack_pos  *Position // Position for ranged attack
	movetoattack_time time.Time // Time of ranged attack

	// Elemental
	elemental *Elemental // Player's elemental
	clone     bool       // Whether this player is an offline shop clone

	// Player-specific fields
	sex   byte
	jobID uint16
	level uint16
	hp    uint32
	maxHP uint32
	sp    uint32
	maxSP uint32

	// Appearance
	hairStyle    uint16
	hairColor    uint16
	clothesColor uint16
	headgearTop  uint16
	headgearMid  uint16
	headgearLow  uint16
	weapon       uint16
	shield       uint16
	robe         uint16

	// Status
	sitting bool
	dead    bool

	// Guild and party information
	guildID    uint32
	guildName  string
	guildTitle string
	emblemID   uint32
	partyName  string

	// Movement and state tracking
	teleported  bool
	disappeared bool
	goneTime    time.Time

	// Damage tracking
	dmgTaken map[string]uint32
	dmgDone  map[string]uint32
}

// NewPlayer creates a new Player with the given ID
func NewPlayer(id []byte) *Player {
	return &Player{
		BaseActor: NewBaseActor(id, "Player"),
		dmgTaken:  make(map[string]uint32),
		dmgDone:   make(map[string]uint32),
	}
}

// Job returns the job name of the player
func (p *Player) Job() string {
	// This would need to be populated with job names from a lookup table
	return fmt.Sprintf("Job %d", p.jobID)
}

// SetJob sets the job ID of the player
func (p *Player) SetJob(jobID uint16) {
	p.jobID = jobID
}

// Level returns the level of the player
func (p *Player) Level() uint16 {
	return p.level
}

// SetLevel sets the level of the player
func (p *Player) SetLevel(level uint16) {
	p.level = level
}

// Sex returns the sex of the player (0 = female, 1 = male)
func (p *Player) Sex() byte {
	return p.sex
}

// SetSex sets the sex of the player
func (p *Player) SetSex(sex byte) {
	p.sex = sex
}

// HP returns the current HP of the player
func (p *Player) HP() uint32 {
	return p.hp
}

// MaxHP returns the maximum HP of the player
func (p *Player) MaxHP() uint32 {
	return p.maxHP
}

// SetHP sets the current HP of the player
func (p *Player) SetHP(hp uint32) {
	p.hp = hp
}

// SetMaxHP sets the maximum HP of the player
func (p *Player) SetMaxHP(maxHP uint32) {
	p.maxHP = maxHP
}

// HPPercent returns the HP percentage of the player
func (p *Player) HPPercent() int {
	if p.maxHP == 0 {
		return 0
	}
	return int((float64(p.hp) / float64(p.maxHP)) * 100)
}

// SP returns the current SP of the player
func (p *Player) SP() uint32 {
	return p.sp
}

// MaxSP returns the maximum SP of the player
func (p *Player) MaxSP() uint32 {
	return p.maxSP
}

// SetSP sets the current SP of the player
func (p *Player) SetSP(sp uint32) {
	p.sp = sp
}

// SetMaxSP sets the maximum SP of the player
func (p *Player) SetMaxSP(maxSP uint32) {
	p.maxSP = maxSP
}

// SPPercent returns the SP percentage of the player
func (p *Player) SPPercent() int {
	if p.maxSP == 0 {
		return 0
	}
	return int((float64(p.sp) / float64(p.maxSP)) * 100)
}

// IsSitting returns whether the player is sitting
func (p *Player) IsSitting() bool {
	return p.sitting
}

// SetSitting sets whether the player is sitting
func (p *Player) SetSitting(sitting bool) {
	p.sitting = sitting
}

// IsDead returns whether the player is dead
func (p *Player) IsDead() bool {
	return p.dead
}

// SetDead sets whether the player is dead
func (p *Player) SetDead(dead bool) {
	p.dead = dead
}

// GuildID returns the guild ID of the player
func (p *Player) GuildID() uint32 {
	return p.guildID
}

// SetGuildID sets the guild ID of the player
func (p *Player) SetGuildID(guildID uint32) {
	p.guildID = guildID
}

// GuildName returns the guild name of the player
func (p *Player) GuildName() string {
	return p.guildName
}

// SetGuildName sets the guild name of the player
func (p *Player) SetGuildName(guildName string) {
	p.guildName = guildName
}

// GuildTitle returns the guild title of the player
func (p *Player) GuildTitle() string {
	return p.guildTitle
}

// SetGuildTitle sets the guild title of the player
func (p *Player) SetGuildTitle(guildTitle string) {
	p.guildTitle = guildTitle
}

// EmblemID returns the emblem ID of the player's guild
func (p *Player) EmblemID() uint32 {
	return p.emblemID
}

// SetEmblemID sets the emblem ID of the player's guild
func (p *Player) SetEmblemID(emblemID uint32) {
	p.emblemID = emblemID
}

// PartyName returns the party name of the player
func (p *Player) PartyName() string {
	return p.partyName
}

// SetPartyName sets the party name of the player
func (p *Player) SetPartyName(partyName string) {
	p.partyName = partyName
}

// SetAppearance sets the appearance of the player
func (p *Player) SetAppearance(hairStyle, hairColor, clothesColor uint16) {
	p.hairStyle = hairStyle
	p.hairColor = hairColor
	p.clothesColor = clothesColor
}

// SetHeadgear sets the headgear of the player
func (p *Player) SetHeadgear(top, mid, low uint16) {
	p.headgearTop = top
	p.headgearMid = mid
	p.headgearLow = low
}

// SetEquipment sets the equipment of the player
func (p *Player) SetEquipment(weapon, shield, robe uint16) {
	p.weapon = weapon
	p.shield = shield
	p.robe = robe
}

// IsTeleported returns whether the player has teleported
func (p *Player) IsTeleported() bool {
	return p.teleported
}

// SetTeleported sets whether the player has teleported
func (p *Player) SetTeleported(teleported bool) {
	p.teleported = teleported
}

// IsDisappeared returns whether the player has disappeared
func (p *Player) IsDisappeared() bool {
	return p.disappeared
}

// SetDisappeared sets whether the player has disappeared
func (p *Player) SetDisappeared(disappeared bool) {
	p.disappeared = disappeared
}

// GoneTime returns the time when the player disappeared
func (p *Player) GoneTime() time.Time {
	return p.goneTime
}

// SetGoneTime sets the time when the player disappeared
func (p *Player) SetGoneTime(goneTime time.Time) {
	p.goneTime = goneTime
}

// AddDamageTaken adds damage taken from a source
func (p *Player) AddDamageTaken(source string, damage uint32) {
	p.dmgTaken[source] += damage
}

// AddDamageDone adds damage done to a target
func (p *Player) AddDamageDone(target string, damage uint32) {
	p.dmgDone[target] += damage
}

// GetDamageTaken returns the total damage taken from a source
func (p *Player) GetDamageTaken(source string) uint32 {
	return p.dmgTaken[source]
}

// GetDamageDone returns the total damage done to a target
func (p *Player) GetDamageDone(target string) uint32 {
	return p.dmgDone[target]
}

// NameString returns a string representation of the player's name
func (p *Player) NameString() string {
	return fmt.Sprintf("Player %s (%s)", p.Name(), p.Job())
}

// DeepCopy creates a deep copy of the player
func (p *Player) DeepCopy() Actor {
	baseCopy := p.BaseActor.DeepCopy().(*BaseActor)

	playerCopy := &Player{
		BaseActor:         baseCopy,
		sex:               p.sex,
		jobID:             p.jobID,
		level:             p.level,
		hp:                p.hp,
		maxHP:             p.maxHP,
		sp:                p.sp,
		maxSP:             p.maxSP,
		hairStyle:         p.hairStyle,
		hairColor:         p.hairColor,
		clothesColor:      p.clothesColor,
		headgearTop:       p.headgearTop,
		headgearMid:       p.headgearMid,
		headgearLow:       p.headgearLow,
		weapon:            p.weapon,
		shield:            p.shield,
		robe:              p.robe,
		sitting:           p.sitting,
		dead:              p.dead,
		guildID:           p.guildID,
		guildName:         p.guildName,
		guildTitle:        p.guildTitle,
		emblemID:          p.emblemID,
		partyName:         p.partyName,
		teleported:        p.teleported,
		disappeared:       p.disappeared,
		goneTime:          p.goneTime,
		dmgTaken:          make(map[string]uint32),
		dmgDone:           make(map[string]uint32),
		cartType:          p.cartType,
		spirits:           p.spirits,
		amuletType:        p.amuletType,
		spiritsType:       p.spiritsType,
		clone:             p.clone,
		movetoattack_pos:  p.movetoattack_pos,
		movetoattack_time: p.movetoattack_time,
		elemental:         p.elemental,
	}

	// Copy look direction
	playerCopy.look.head = p.look.head
	playerCopy.look.body = p.look.body

	// Copy damage maps
	for source, damage := range p.dmgTaken {
		playerCopy.dmgTaken[source] = damage
	}

	for target, damage := range p.dmgDone {
		playerCopy.dmgDone[target] = damage
	}

	return playerCopy
}

// HeadDirection returns the direction the player's head is facing
func (p *Player) HeadDirection() byte {
	return p.look.head
}

// BodyDirection returns the direction the player's body is facing
func (p *Player) BodyDirection() byte {
	return p.look.body
}

// SetLookDirection sets the direction the player is looking
func (p *Player) SetLookDirection(head, body byte) {
	p.look.head = head
	p.look.body = body
}

// CartType returns the type of cart the player has
func (p *Player) CartType() uint32 {
	return p.cartType
}

// SetCartType sets the type of cart the player has
func (p *Player) SetCartType(cartType uint32) {
	p.cartType = cartType
}

// Spirits returns the number of spirits/counters the player has
func (p *Player) Spirits() uint32 {
	return p.spirits
}

// SetSpirits sets the number of spirits/counters the player has
func (p *Player) SetSpirits(spirits uint32) {
	p.spirits = spirits
}

// AmuletType returns the type of amulet the player has
func (p *Player) AmuletType() string {
	return p.amuletType
}

// SetAmuletType sets the type of amulet the player has
func (p *Player) SetAmuletType(amuletType string) {
	p.amuletType = amuletType
}

// SpiritsType returns the type of spirits the player has
func (p *Player) SpiritsType() string {
	return p.spiritsType
}

// SetSpiritsType sets the type of spirits the player has
func (p *Player) SetSpiritsType(spiritsType string) {
	p.spiritsType = spiritsType
}

// MoveToAttackPos returns the position for ranged attack
func (p *Player) MoveToAttackPos() *Position {
	return p.movetoattack_pos
}

// SetMoveToAttackPos sets the position for ranged attack
func (p *Player) SetMoveToAttackPos(pos *Position) {
	p.movetoattack_pos = pos
}

// MoveToAttackTime returns the time of ranged attack
func (p *Player) MoveToAttackTime() time.Time {
	return p.movetoattack_time
}

// SetMoveToAttackTime sets the time of ranged attack
func (p *Player) SetMoveToAttackTime(t time.Time) {
	p.movetoattack_time = t
}

// Elemental returns the player's elemental
func (p *Player) Elemental() *Elemental {
	return p.elemental
}

// SetElemental sets the player's elemental
func (p *Player) SetElemental(elemental *Elemental) {
	p.elemental = elemental
}

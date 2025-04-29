package actor

import (
	"strconv"
	"time"
)

// Actor is the interface that all game entities must implement
type Actor interface {
	// ID returns the unique identifier of the actor
	ID() []byte

	// NameID returns the ID decoded into a 32-bit integer
	NameID() uint32

	// Name returns the name of the actor
	Name() string

	// SetName sets the name of the actor
	SetName(name string)

	// Position returns the current position of the actor
	Position() *Position

	// PositionTo returns the position the actor is moving to
	PositionTo() *Position

	// SetPosition sets the current position of the actor
	SetPosition(pos *Position)

	// SetPositionTo sets the position the actor is moving to
	SetPositionTo(pos *Position)

	// ActorType returns the type of the actor (Player, Monster, NPC, etc.)
	ActorType() string

	// IsAvoid returns whether the actor should be avoided
	IsAvoid() bool

	// SetAvoid sets whether the actor should be avoided
	SetAvoid(avoid bool)

	// AppearTime returns the time when the actor first appeared
	AppearTime() time.Time

	// WalkSpeed returns the actor's walking speed in blocks per second
	WalkSpeed() float64

	// SetWalkSpeed sets the actor's walking speed
	SetWalkSpeed(speed float64)

	// TimeMove returns the time at which the actor last moved
	TimeMove() time.Time

	// SetTimeMove sets the time at which the actor last moved
	SetTimeMove(t time.Time)

	// TimeMoveCalc returns the time in seconds that the actor needs to move from Position to PositionTo
	TimeMoveCalc() float64

	// SetTimeMoveCalc sets the time in seconds that the actor needs to move from Position to PositionTo
	SetTimeMoveCalc(t float64)

	// DeepCopy creates a deep copy of the actor
	DeepCopy() Actor

	// NameString returns a string representation of the actor's name
	NameString() string
}

// BaseActor provides common functionality for all actor types
type BaseActor struct {
	id           []byte
	nameID       uint32
	name         string
	pos          *Position
	posTo        *Position
	actorType    string
	avoid        bool
	appearTime   time.Time
	walkSpeed    float64
	timeMove     time.Time
	timeMoveCalc float64
}

// NewBaseActor creates a new BaseActor with the given ID and actor type
func NewBaseActor(id []byte, actorType string) *BaseActor {
	return &BaseActor{
		id:         id,
		nameID:     bytesToUint32(id),
		actorType:  actorType,
		pos:        &Position{},
		posTo:      &Position{},
		appearTime: time.Now(),
		walkSpeed:  0.15, // Default walk speed
	}
}

// ID returns the unique identifier of the actor
func (a *BaseActor) ID() []byte {
	return a.id
}

// NameID returns the ID decoded into a 32-bit integer
func (a *BaseActor) NameID() uint32 {
	return a.nameID
}

// Name returns the name of the actor
func (a *BaseActor) Name() string {
	if a.name == "" {
		return "Unknown #" + strconv.FormatUint(uint64(a.nameID), 10)
	}
	return a.name
}

// SetName sets the name of the actor
func (a *BaseActor) SetName(name string) {
	a.name = name
}

// Position returns the current position of the actor
func (a *BaseActor) Position() *Position {
	return a.pos
}

// PositionTo returns the position the actor is moving to
func (a *BaseActor) PositionTo() *Position {
	return a.posTo
}

// SetPosition sets the current position of the actor
func (a *BaseActor) SetPosition(pos *Position) {
	a.pos = pos
}

// SetPositionTo sets the position the actor is moving to
func (a *BaseActor) SetPositionTo(pos *Position) {
	a.posTo = pos
}

// ActorType returns the type of the actor
func (a *BaseActor) ActorType() string {
	return a.actorType
}

// IsAvoid returns whether the actor should be avoided
func (a *BaseActor) IsAvoid() bool {
	return a.avoid
}

// SetAvoid sets whether the actor should be avoided
func (a *BaseActor) SetAvoid(avoid bool) {
	a.avoid = avoid
}

// AppearTime returns the time when the actor first appeared
func (a *BaseActor) AppearTime() time.Time {
	return a.appearTime
}

// WalkSpeed returns the actor's walking speed in blocks per second
func (a *BaseActor) WalkSpeed() float64 {
	return a.walkSpeed
}

// SetWalkSpeed sets the actor's walking speed
func (a *BaseActor) SetWalkSpeed(speed float64) {
	a.walkSpeed = speed
}

// TimeMove returns the time at which the actor last moved
func (a *BaseActor) TimeMove() time.Time {
	return a.timeMove
}

// SetTimeMove sets the time at which the actor last moved
func (a *BaseActor) SetTimeMove(t time.Time) {
	a.timeMove = t
}

// TimeMoveCalc returns the time in seconds that the actor needs to move from Position to PositionTo
func (a *BaseActor) TimeMoveCalc() float64 {
	return a.timeMoveCalc
}

// SetTimeMoveCalc sets the time in seconds that the actor needs to move from Position to PositionTo
func (a *BaseActor) SetTimeMoveCalc(t float64) {
	a.timeMoveCalc = t
}

// NameString returns a string representation of the actor's name
func (a *BaseActor) NameString() string {
	return a.ActorType() + " " + a.Name()
}

// DeepCopy creates a deep copy of the actor
func (a *BaseActor) DeepCopy() Actor {
	actorCopy := &BaseActor{
		id:           make([]byte, len(a.id)),
		nameID:       a.nameID,
		name:         a.name,
		pos:          &Position{X: a.pos.X, Y: a.pos.Y},
		posTo:        &Position{X: a.posTo.X, Y: a.posTo.Y},
		actorType:    a.actorType,
		avoid:        a.avoid,
		appearTime:   a.appearTime,
		walkSpeed:    a.walkSpeed,
		timeMove:     a.timeMove,
		timeMoveCalc: a.timeMoveCalc,
	}

	copy(actorCopy.id, a.id)
	return actorCopy
}

// bytesToUint32 converts a byte slice to a uint32
func bytesToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

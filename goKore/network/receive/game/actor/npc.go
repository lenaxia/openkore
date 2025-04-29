package actor

import (
	"fmt"
	"time"
)

// NPC represents a non-player character in the game
type NPC struct {
	*BaseActor

	// NPC-specific fields
	guildID  uint32 // Guild ID (for guild flag NPCs)
	emblemID uint32 // Emblem ID (for guild flag NPCs)

	// Status
	disappeared bool // Whether the NPC has disappeared

	// Time tracking
	goneTime time.Time // Time when the NPC disappeared
}

// NewNPC creates a new NPC with the given ID
func NewNPC(id []byte) *NPC {
	return &NPC{
		BaseActor: NewBaseActor(id, "NPC"),
	}
}

// GuildID returns the guild ID of the NPC (for guild flag NPCs)
func (n *NPC) GuildID() uint32 {
	return n.guildID
}

// SetGuildID sets the guild ID of the NPC
func (n *NPC) SetGuildID(guildID uint32) {
	n.guildID = guildID
}

// EmblemID returns the emblem ID of the NPC (for guild flag NPCs)
func (n *NPC) EmblemID() uint32 {
	return n.emblemID
}

// SetEmblemID sets the emblem ID of the NPC
func (n *NPC) SetEmblemID(emblemID uint32) {
	n.emblemID = emblemID
}

// IsDisappeared returns whether the NPC has disappeared
func (n *NPC) IsDisappeared() bool {
	return n.disappeared
}

// SetDisappeared sets whether the NPC has disappeared
func (n *NPC) SetDisappeared(disappeared bool) {
	n.disappeared = disappeared
}

// GoneTime returns the time when the NPC disappeared
func (n *NPC) GoneTime() time.Time {
	return n.goneTime
}

// SetGoneTime sets the time when the NPC disappeared
func (n *NPC) SetGoneTime(goneTime time.Time) {
	n.goneTime = goneTime
}

// SendTalk sends a talk request to the NPC
func (n *NPC) SendTalk() {
	// This would be implemented to send a talk request to the NPC
	// In the actual implementation, this would call the message sender
}

// NameString returns a string representation of the NPC's name
func (n *NPC) NameString() string {
	return fmt.Sprintf("NPC %s", n.Name())
}

// DeepCopy creates a deep copy of the NPC
func (n *NPC) DeepCopy() Actor {
	baseCopy := n.BaseActor.DeepCopy().(*BaseActor)

	npcCopy := &NPC{
		BaseActor:   baseCopy,
		guildID:     n.guildID,
		emblemID:    n.emblemID,
		disappeared: n.disappeared,
		goneTime:    n.goneTime,
	}

	return npcCopy
}

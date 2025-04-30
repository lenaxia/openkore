// Package skill provides handlers for skill-related packets.
package skill

import (
	"github.com/lenaxia/goKore/network/hooks"
)

// SkillOwnerType represents the owner of a skill
type SkillOwnerType int

const (
	// OwnerChar indicates the skill belongs to the character
	OwnerChar SkillOwnerType = iota
	// OwnerHomun indicates the skill belongs to the character's homunculus
	OwnerHomun
	// OwnerMerc indicates the skill belongs to the character's mercenary
	OwnerMerc
)

// SkillManager manages skill-related packet handlers
type SkillManager struct {
	hookManager *hooks.HookManager
}

// NewSkillManager creates a new skill manager
func NewSkillManager(hookManager *hooks.HookManager) *SkillManager {
	return &SkillManager{
		hookManager: hookManager,
	}
}

// SkillInfo represents information about a skill
type SkillInfo struct {
	ID         uint16 // Skill ID
	TargetType uint32 // Target type
	Level      uint16 // Skill level
	SP         uint16 // SP cost
	Range      uint16 // Range
	Handle     string // Skill handle (name)
	Up         uint8  // Upgradable flag
	Level2     uint16 // Secondary level (for some skills)
}

package entity

import (
	"time"
)

type StatusID int

const (
	StatusPoison StatusID = iota + 1
	StatusSilence
	StatusHaste
)

// ApplyStatus implements Actor::setStatus logic
func (s *StatusSystem) ApplyStatus(id StatusID, active bool, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if expiry := s.immunity[id]; expiry.After(time.Now()) {
		return StatusErrorf("immune to %v until %v", id, expiry)
	}
	
	if active {
		if existing, exists := s.effects[id]; exists {
			if existing.Duration < duration {
				existing.Duration = duration // Refresh duration
			}
			return nil
		}
		
		s.effects[id] = &StatusEffect{
			ID:        id,
			StartedAt: time.Now(),
			Duration:  duration,
		}
	} else {
		delete(s.effects, id)
	}
	return nil
}

// CheckStatusConflicts validates status interactions
func CheckStatusConflicts(existing []StatusEffect, newStatus StatusID) error {
	// Implement status override rules from Actor.pm
	for _, status := range existing {
		if status.BlockedBy(newStatus) {
			return StatusConflictError{
				Existing:   status.ID,
				New:        newStatus,
			}
		}
	}
	return nil
}

// Package social provides functionality for social aspects of the game
package social

import (
	"encoding/binary"
	"errors"
	"sync"
)

// Reputation represents a reputation entry
type Reputation struct {
	Type    uint32
	Type2   uint32
	Points  uint32
	Points2 uint32
}

// ReputationManager manages reputation-related functionality
type ReputationManager struct {
	mutex       sync.RWMutex
	reputations []Reputation
}

// NewReputationManager creates a new reputation manager
func NewReputationManager() *ReputationManager {
	return &ReputationManager{
		reputations: make([]Reputation, 0),
	}
}

// HandleReputeInfo handles the repute_info packet
// This packet is sent by the server to update the character's reputation information
// Packet format: 0B8D <reputeInfo>.B
func (m *ReputationManager) HandleReputeInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract reputation info from args
	var reputeInfo []byte
	if reputeInfoVal, ok := args["reputeInfo"].([]byte); ok {
		reputeInfo = reputeInfoVal
	} else {
		return errors.New("invalid reputation info in repute_info packet")
	}

	// Each reputation entry is 16 bytes (4 uint32 values)
	entrySize := 16
	if len(reputeInfo)%entrySize != 0 {
		return errors.New("invalid reputation info length in repute_info packet")
	}

	// Clear the existing reputation list
	m.reputations = make([]Reputation, 0)

	// Parse the reputation info
	for i := 0; i < len(reputeInfo); i += entrySize {
		// Extract the reputation entry
		if i+entrySize > len(reputeInfo) {
			return errors.New("invalid reputation info length in repute_info packet")
		}

		// Parse the reputation entry
		reputation := Reputation{
			Type:    binary.LittleEndian.Uint32(reputeInfo[i : i+4]),
			Type2:   binary.LittleEndian.Uint32(reputeInfo[i+4 : i+8]),
			Points:  binary.LittleEndian.Uint32(reputeInfo[i+8 : i+12]),
			Points2: binary.LittleEndian.Uint32(reputeInfo[i+12 : i+16]),
		}

		// Add the reputation entry to the list
		m.reputations = append(m.reputations, reputation)
	}

	return nil
}

// GetReputations returns a copy of the reputation list
func (m *ReputationManager) GetReputations() []Reputation {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy of the reputation list
	reputations := make([]Reputation, len(m.reputations))
	copy(reputations, m.reputations)

	return reputations
}

// RegisterHandlers registers reputation-related packet handlers
func (m *ReputationManager) RegisterHandlers(parser interface{}) {
	// This function would register the reputation-related packet handlers
	// with the parser, but we'll leave it empty for now since we don't have
	// access to the parser interface in this package
}

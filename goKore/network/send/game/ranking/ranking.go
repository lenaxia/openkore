// Package ranking provides ranking-related packet sending functionality.
package ranking

import (
	"fmt"

	"github.com/lenaxia/goKore/network/send/core"
)

// RankingManager handles ranking-related packet sending.
type RankingManager struct {
	// Base send implementation
	baseSend core.Send
	// Whether to use the new ranking system
	rankingSystemType bool
}

// NewRankingManager creates a new ranking manager.
func NewRankingManager(baseSend core.Send) *RankingManager {
	return &RankingManager{
		baseSend:          baseSend,
		rankingSystemType: false, // Default to legacy system
	}
}

// SetRankingSystemType sets whether to use the new ranking system.
func (rm *RankingManager) SetRankingSystemType(value bool) {
	rm.rankingSystemType = value
}

// GetRankingSystemType returns whether the new ranking system is being used.
func (rm *RankingManager) GetRankingSystemType() bool {
	return rm.rankingSystemType
}

// SendAchievementGetReward sends a request to get an achievement reward.
// This is equivalent to the sendAchievementGetReward function in Send.pm.
func (rm *RankingManager) SendAchievementGetReward(achievementID uint32) error {
	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("achievement_get_reward")
	if !exists {
		return fmt.Errorf("achievement_get_reward packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"achievementID": achievementID,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

// SendTop10Alchemist sends a request to view the top 10 alchemists.
// This is equivalent to the sendTop10Alchemist function in Send.pm.
func (rm *RankingManager) SendTop10Alchemist() error {
	if !rm.rankingSystemType {
		// Legacy system
		packetID, exists := rm.baseSend.GetPacketID("rank_alchemist")
		if !exists {
			return fmt.Errorf("rank_alchemist packet ID not found")
		}

		// No arguments needed for this packet
		args := map[string]interface{}{}

		// Construct and send the packet
		packet, err := rm.baseSend.Reconstruct(packetID, args)
		if err != nil {
			return err
		}

		return rm.baseSend.SendToServer(packet)
	} else {
		// New system
		return rm.SendTop10(1)
	}
}

// SendTop10Blacksmith sends a request to view the top 10 blacksmiths.
// This is equivalent to the sendTop10Blacksmith function in Send.pm.
func (rm *RankingManager) SendTop10Blacksmith() error {
	if !rm.rankingSystemType {
		// Legacy system
		packetID, exists := rm.baseSend.GetPacketID("rank_blacksmith")
		if !exists {
			return fmt.Errorf("rank_blacksmith packet ID not found")
		}

		// No arguments needed for this packet
		args := map[string]interface{}{}

		// Construct and send the packet
		packet, err := rm.baseSend.Reconstruct(packetID, args)
		if err != nil {
			return err
		}

		return rm.baseSend.SendToServer(packet)
	} else {
		// New system
		return rm.SendTop10(0)
	}
}

// SendTop10PK sends a request to view the top 10 PKs.
// This is equivalent to the sendTop10PK function in Send.pm.
func (rm *RankingManager) SendTop10PK() error {
	if !rm.rankingSystemType {
		// Legacy system
		packetID, exists := rm.baseSend.GetPacketID("rank_killer")
		if !exists {
			return fmt.Errorf("rank_killer packet ID not found")
		}

		// No arguments needed for this packet
		args := map[string]interface{}{}

		// Construct and send the packet
		packet, err := rm.baseSend.Reconstruct(packetID, args)
		if err != nil {
			return err
		}

		return rm.baseSend.SendToServer(packet)
	} else {
		// New system
		return rm.SendTop10(3)
	}
}

// SendTop10Taekwon sends a request to view the top 10 taekwons.
// This is equivalent to the sendTop10Taekwon function in Send.pm.
func (rm *RankingManager) SendTop10Taekwon() error {
	if !rm.rankingSystemType {
		// Legacy system
		packetID, exists := rm.baseSend.GetPacketID("rank_taekwon")
		if !exists {
			return fmt.Errorf("rank_taekwon packet ID not found")
		}

		// No arguments needed for this packet
		args := map[string]interface{}{}

		// Construct and send the packet
		packet, err := rm.baseSend.Reconstruct(packetID, args)
		if err != nil {
			return err
		}

		return rm.baseSend.SendToServer(packet)
	} else {
		// New system
		return rm.SendTop10(2)
	}
}

// SendTop10 sends a request to view the top 10 of a specific type.
// This is equivalent to the sendTop10 function in Send.pm.
// rankType:
//
//	0 => Blacksmith
//	1 => Alchemist
//	2 => Taekwon
//	3 => PK
func (rm *RankingManager) SendTop10(rankType uint8) error {
	// Validate rank type
	if rankType > 3 {
		return fmt.Errorf("invalid rank type: %d (must be 0-3)", rankType)
	}

	// Get the packet ID
	packetID, exists := rm.baseSend.GetPacketID("rank_general")
	if !exists {
		return fmt.Errorf("rank_general packet ID not found")
	}

	// Create the arguments
	args := map[string]interface{}{
		"type": rankType,
	}

	// Construct and send the packet
	packet, err := rm.baseSend.Reconstruct(packetID, args)
	if err != nil {
		return err
	}

	return rm.baseSend.SendToServer(packet)
}

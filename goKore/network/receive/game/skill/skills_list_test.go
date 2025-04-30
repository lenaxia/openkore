package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestSkillsList(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Create the skills list manager
	manager := NewSkillsListManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Test cases for different packet types
	testCases := []struct {
		name           string
		packetSwitch   string
		rawMsg         []byte
		expectedOwner  SkillOwnerType
		expectedSkills int
		hookName       string
	}{
		{
			name:           "Character Skills (010F)",
			packetSwitch:   "010F",
			rawMsg:         createCharSkillsPacket(),
			expectedOwner:  OwnerChar,
			expectedSkills: 2,
			hookName:       "character.skills_list",
		},
		{
			name:           "Homunculus Skills (0235)",
			packetSwitch:   "0235",
			rawMsg:         createHomunSkillsPacket(),
			expectedOwner:  OwnerHomun,
			expectedSkills: 1,
			hookName:       "homunculus.skills_list",
		},
		{
			name:           "Mercenary Skills (029D)",
			packetSwitch:   "029D",
			rawMsg:         createMercSkillsPacket(),
			expectedOwner:  OwnerMerc,
			expectedSkills: 1,
			hookName:       "mercenary.skills_list",
		},
		{
			name:           "Character Skills Short (0B32)",
			packetSwitch:   "0B32",
			rawMsg:         createCharSkillsShortPacket(),
			expectedOwner:  OwnerChar,
			expectedSkills: 1,
			hookName:       "character.skills_list",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a channel to receive hook events
			hookCalled := false
			var hookResult map[string]interface{}

			// Register a hook to capture the event
			hookManager.AddHook(tc.hookName, func(hookName string, arg interface{}, userData interface{}) {
				hookCalled = true
				if result, ok := arg.(map[string]interface{}); ok {
					hookResult = result
				}
			}, nil)

			// Create packet data
			args := map[string]interface{}{
				"switch":       tc.packetSwitch,
				"RAW_MSG":      tc.rawMsg,
				"RAW_MSG_SIZE": len(tc.rawMsg),
			}

			// Call the handler
			err := manager.handleSkillsList(args)
			if err != nil {
				t.Errorf("handleSkillsList() returned error: %v", err)
			}

			// Check that the hook was called
			if !hookCalled {
				t.Error("Hook was not called")
			}

			// Check the hook result
			if hookResult == nil {
				t.Fatal("Hook result is nil")
			}

			// Check the owner type
			if ownerType, ok := hookResult["ownerType"].(SkillOwnerType); !ok || ownerType != tc.expectedOwner {
				t.Errorf("Expected owner type %v, got %v", tc.expectedOwner, hookResult["ownerType"])
			}

			// Check the skills list
			if skills, ok := hookResult["skills"].([]SkillInfo); !ok || len(skills) != tc.expectedSkills {
				t.Errorf("Expected %d skills, got %d", tc.expectedSkills, len(skills))
			}
		})
	}
}

// Test unhappy paths
func TestSkillsListUnhappy(t *testing.T) {
	// Create a skills list manager
	manager := NewSkillsListManager(nil, nil)

	// Test case 1: Missing switch
	t.Run("MissingSwitch", func(t *testing.T) {
		args := map[string]interface{}{
			"RAW_MSG":      []byte{0x01, 0x0F, 0x00, 0x00},
			"RAW_MSG_SIZE": 4,
		}

		err := manager.handleSkillsList(args)
		if err == nil {
			t.Error("Expected error for missing switch, got nil")
		}
	})

	// Test case 2: Missing RAW_MSG
	t.Run("MissingRawMsg", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":       "010F",
			"RAW_MSG_SIZE": 4,
		}

		err := manager.handleSkillsList(args)
		if err == nil {
			t.Error("Expected error for missing RAW_MSG, got nil")
		}
	})

	// Test case 3: Missing RAW_MSG_SIZE
	t.Run("MissingRawMsgSize", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":  "010F",
			"RAW_MSG": []byte{0x01, 0x0F, 0x00, 0x00},
		}

		err := manager.handleSkillsList(args)
		if err == nil {
			t.Error("Expected error for missing RAW_MSG_SIZE, got nil")
		}
	})

	// Test case 4: Unknown packet switch
	t.Run("UnknownPacketSwitch", func(t *testing.T) {
		args := map[string]interface{}{
			"switch":       "FFFF", // Unknown packet switch
			"RAW_MSG":      []byte{0xFF, 0xFF, 0x00, 0x00},
			"RAW_MSG_SIZE": 4,
		}

		// This should not return an error since we handle unknown packet switches gracefully
		err := manager.handleSkillsList(args)
		if err != nil {
			t.Errorf("handleSkillsList() returned error for unknown packet switch: %v", err)
		}
	})
}

// Helper function to create a character skills packet (010F)
func createCharSkillsPacket() []byte {
	// Create a packet with 2 skills
	// Format: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <handle>.24B <up>.B
	packet := []byte{
		// Packet header
		0x0F, 0x01, // Packet ID
		0x50, 0x00, // Packet length (80 bytes)

		// Skill 1
		0x01, 0x00, // ID (1)
		0x01, 0x00, 0x00, 0x00, // Target type (1)
		0x05, 0x00, // Level (5)
		0x0A, 0x00, // SP (10)
		0x03, 0x00, // Range (3)
		// Handle "NV_BASIC" (24 bytes)
		0x4E, 0x56, 0x5F, 0x42, 0x41, 0x53, 0x49, 0x43, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Up (1)

		// Skill 2
		0x02, 0x00, // ID (2)
		0x02, 0x00, 0x00, 0x00, // Target type (2)
		0x03, 0x00, // Level (3)
		0x08, 0x00, // SP (8)
		0x05, 0x00, // Range (5)
		// Handle "SM_SWORD" (24 bytes)
		0x53, 0x4D, 0x5F, 0x53, 0x57, 0x4F, 0x52, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Up (1)
	}

	return packet
}

// Helper function to create a homunculus skills packet (0235)
func createHomunSkillsPacket() []byte {
	// Create a packet with 1 skill
	// Format: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <handle>.24B <up>.B
	packet := []byte{
		// Packet header
		0x35, 0x02, // Packet ID
		0x29, 0x00, // Packet length (41 bytes)

		// Skill 1
		0x03, 0x00, // ID (3)
		0x01, 0x00, 0x00, 0x00, // Target type (1)
		0x02, 0x00, // Level (2)
		0x05, 0x00, // SP (5)
		0x02, 0x00, // Range (2)
		// Handle "HFLI_MOON" (24 bytes)
		0x48, 0x46, 0x4C, 0x49, 0x5F, 0x4D, 0x4F, 0x4F, 0x4E, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Up (1)
	}

	return packet
}

// Helper function to create a mercenary skills packet (029D)
func createMercSkillsPacket() []byte {
	// Create a packet with 1 skill
	// Format: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <handle>.24B <up>.B
	packet := []byte{
		// Packet header
		0x9D, 0x02, // Packet ID
		0x29, 0x00, // Packet length (41 bytes)

		// Skill 1
		0x04, 0x00, // ID (4)
		0x01, 0x00, 0x00, 0x00, // Target type (1)
		0x03, 0x00, // Level (3)
		0x07, 0x00, // SP (7)
		0x04, 0x00, // Range (4)
		// Handle "MS_BASH" (24 bytes)
		0x4D, 0x53, 0x5F, 0x42, 0x41, 0x53, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, // Up (1)
	}

	return packet
}

// Helper function to create a character skills short packet (0B32)
func createCharSkillsShortPacket() []byte {
	// Create a packet with 1 skill
	// Format: <packet len>.W <ID>.W <targetType>.L <lv>.W <sp>.W <range>.W <up>.B <lv2>.W
	packet := []byte{
		// Packet header
		0x32, 0x0B, // Packet ID
		0x13, 0x00, // Packet length (19 bytes)

		// Skill 1
		0x05, 0x00, // ID (5)
		0x01, 0x00, 0x00, 0x00, // Target type (1)
		0x04, 0x00, // Level (4)
		0x09, 0x00, // SP (9)
		0x03, 0x00, // Range (3)
		0x01,       // Up (1)
		0x02, 0x00, // Level2 (2)
	}

	return packet
}

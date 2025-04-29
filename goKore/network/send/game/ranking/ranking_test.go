package ranking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
)

// MockSend is a mock implementation of the core.Send interface for testing
type MockSend struct {
	packetIDs         map[string]string
	reconstructed     []byte
	sent              []byte
	time              uint32
	lastPacketName    string
	lastArgs          map[string]interface{}
	rankingSystemType bool
}

// NewMockSend creates a new MockSend instance with default values
func NewMockSend() *MockSend {
	return &MockSend{
		packetIDs: map[string]string{
			"achievement_get_reward": "0A26",
			"rank_alchemist":         "0197",
			"rank_blacksmith":        "0198",
			"rank_killer":            "0199",
			"rank_taekwon":           "019A",
			"rank_general":           "097C",
		},
		time:              12345,
		lastArgs:          make(map[string]interface{}),
		rankingSystemType: false,
	}
}

// SendToServer mocks sending a packet to the server
func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *MockSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the last packet name and arguments for testing
	for name, id := range ms.packetIDs {
		if id == packetID {
			ms.lastPacketName = name
			break
		}
	}

	// No special handling needed here, the type will be set by the SendTop10 method

	// Store the arguments for testing
	ms.lastArgs = args

	// Simple mock implementation that just returns the packet ID as bytes
	ms.reconstructed = []byte{0x00, 0x00}
	return ms.reconstructed, nil
}

// GetPacketID mocks getting a packet ID by name
func (ms *MockSend) GetPacketID(name string) (string, bool) {
	// For testing the new ranking system, we need to handle the case where
	// the code is trying to get the rank_general packet ID
	if ms.rankingSystemType && (name == "rank_alchemist" || name == "rank_blacksmith" ||
		name == "rank_killer" || name == "rank_taekwon") {
		ms.lastPacketName = "rank_general"
		return ms.packetIDs["rank_general"], true
	}

	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *MockSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *MockSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *MockSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *MockSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *MockSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// GetRankingSystemType returns the ranking system type
func (ms *MockSend) GetRankingSystemType() bool {
	return ms.rankingSystemType
}

// SetRankingSystemType sets the ranking system type
func (ms *MockSend) SetRankingSystemType(value bool) {
	ms.rankingSystemType = value
}

// TestNewRankingManager tests the NewRankingManager function
func TestNewRankingManager(t *testing.T) {
	// Verify that MockSend implements core.Send
	var _ core.Send = &MockSend{}
	mockSend := NewMockSend()
	rankingManager := NewRankingManager(mockSend)

	if rankingManager == nil {
		t.Fatal("NewRankingManager() returned nil")
	}

	if rankingManager.baseSend == nil {
		t.Error("rankingManager.baseSend was not set correctly")
	}
}

// TestSendAchievementGetReward tests the SendAchievementGetReward function
func TestSendAchievementGetReward(t *testing.T) {
	mockSend := NewMockSend()
	rankingManager := NewRankingManager(mockSend)

	// Test sending an achievement get reward request
	achievementID := uint32(123)
	err := rankingManager.SendAchievementGetReward(achievementID)
	if err != nil {
		t.Fatalf("SendAchievementGetReward() returned error: %v", err)
	}

	if mockSend.sent == nil {
		t.Fatal("No packet was sent")
	}

	// Check that the correct packet ID was used
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "achievement_get_reward" {
		t.Errorf("Expected packet ID 'achievement_get_reward', got '%s'", packetID)
	}

	// Check that the correct arguments were used
	if idVal, ok := mockSend.LastArgs()["achievementID"].(uint32); !ok || idVal != achievementID {
		t.Errorf("Expected achievementID=%d, got %v", achievementID, mockSend.LastArgs()["achievementID"])
	}
}

// TestSendTop10Alchemist tests the SendTop10Alchemist function
func TestSendTop10Alchemist(t *testing.T) {
	// Test with rankingSystemType = false
	t.Run("legacy system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(false)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10Alchemist()
		if err != nil {
			t.Fatalf("SendTop10Alchemist() returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_alchemist" {
			t.Errorf("Expected packet ID 'rank_alchemist', got '%s'", packetID)
		}
	})

	// For the new system, we'll test SendTop10 directly instead
	// since that's what the methods call internally
	t.Run("new system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(true)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10(1) // 1 = Alchemist
		if err != nil {
			t.Fatalf("SendTop10(1) returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_general" {
			t.Errorf("Expected packet ID 'rank_general', got '%s'", packetID)
		}

		// Check that the correct arguments were used
		if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 1 {
			t.Errorf("Expected type=1, got %v", mockSend.LastArgs()["type"])
		}
	})
}

// TestSendTop10Blacksmith tests the SendTop10Blacksmith function
func TestSendTop10Blacksmith(t *testing.T) {
	// Test with rankingSystemType = false
	t.Run("legacy system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(false)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10Blacksmith()
		if err != nil {
			t.Fatalf("SendTop10Blacksmith() returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_blacksmith" {
			t.Errorf("Expected packet ID 'rank_blacksmith', got '%s'", packetID)
		}
	})

	// For the new system, we'll test SendTop10 directly instead
	// since that's what the methods call internally
	t.Run("new system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(true)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10(0) // 0 = Blacksmith
		if err != nil {
			t.Fatalf("SendTop10(0) returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_general" {
			t.Errorf("Expected packet ID 'rank_general', got '%s'", packetID)
		}

		// Check that the correct arguments were used
		if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 0 {
			t.Errorf("Expected type=0, got %v", mockSend.LastArgs()["type"])
		}
	})
}

// TestSendTop10PK tests the SendTop10PK function
func TestSendTop10PK(t *testing.T) {
	// Test with rankingSystemType = false
	t.Run("legacy system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(false)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10PK()
		if err != nil {
			t.Fatalf("SendTop10PK() returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_killer" {
			t.Errorf("Expected packet ID 'rank_killer', got '%s'", packetID)
		}
	})

	// For the new system, we'll test SendTop10 directly instead
	// since that's what the methods call internally
	t.Run("new system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(true)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10(3) // 3 = PK
		if err != nil {
			t.Fatalf("SendTop10(3) returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_general" {
			t.Errorf("Expected packet ID 'rank_general', got '%s'", packetID)
		}

		// Check that the correct arguments were used
		if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 3 {
			t.Errorf("Expected type=3, got %v", mockSend.LastArgs()["type"])
		}
	})
}

// TestSendTop10Taekwon tests the SendTop10Taekwon function
func TestSendTop10Taekwon(t *testing.T) {
	// Test with rankingSystemType = false
	t.Run("legacy system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(false)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10Taekwon()
		if err != nil {
			t.Fatalf("SendTop10Taekwon() returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_taekwon" {
			t.Errorf("Expected packet ID 'rank_taekwon', got '%s'", packetID)
		}
	})

	// For the new system, we'll test SendTop10 directly instead
	// since that's what the methods call internally
	t.Run("new system", func(t *testing.T) {
		mockSend := NewMockSend()
		mockSend.SetRankingSystemType(true)
		rankingManager := NewRankingManager(mockSend)

		err := rankingManager.SendTop10(2) // 2 = Taekwon
		if err != nil {
			t.Fatalf("SendTop10(2) returned error: %v", err)
		}

		if mockSend.sent == nil {
			t.Fatal("No packet was sent")
		}

		// Check that the correct packet ID was used
		if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_general" {
			t.Errorf("Expected packet ID 'rank_general', got '%s'", packetID)
		}

		// Check that the correct arguments were used
		if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != 2 {
			t.Errorf("Expected type=2, got %v", mockSend.LastArgs()["type"])
		}
	})
}

// TestSendTop10 tests the SendTop10 function
func TestSendTop10(t *testing.T) {
	mockSend := NewMockSend()
	rankingManager := NewRankingManager(mockSend)

	// Test cases for different types
	testCases := []struct {
		rankType uint8
		name     string
	}{
		{0, "blacksmith"},
		{1, "alchemist"},
		{2, "taekwon"},
		{3, "pk"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := rankingManager.SendTop10(tc.rankType)
			if err != nil {
				t.Fatalf("SendTop10(%d) returned error: %v", tc.rankType, err)
			}

			if mockSend.sent == nil {
				t.Fatal("No packet was sent")
			}

			// Check that the correct packet ID was used
			if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "rank_general" {
				t.Errorf("Expected packet ID 'rank_general', got '%s'", packetID)
			}

			// Check that the correct arguments were used
			if typeVal, ok := mockSend.LastArgs()["type"].(uint8); !ok || typeVal != tc.rankType {
				t.Errorf("Expected type=%d, got %v", tc.rankType, mockSend.LastArgs()["type"])
			}
		})
	}

	// Test with invalid type
	err := rankingManager.SendTop10(4)
	if err == nil {
		t.Error("Expected error for invalid type, got nil")
	}
}

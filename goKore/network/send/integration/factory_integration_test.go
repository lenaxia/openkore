package integration

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/cashshop"
	"github.com/lenaxia/goKore/network/send/game/chat"
	"github.com/lenaxia/goKore/network/send/game/misc"
)

// SendFactory is a factory for creating send managers
type SendFactory struct {
	baseSend core.Send
}

// NewSendFactory creates a new send factory
func NewSendFactory(baseSend core.Send) *SendFactory {
	return &SendFactory{
		baseSend: baseSend,
	}
}

// CreateCashShopManager creates a new cash shop manager
func (sf *SendFactory) CreateCashShopManager() *cashshop.CashShopManager {
	return cashshop.NewCashShopManager(sf.baseSend)
}

// CreateMiscManager creates a new miscellaneous manager
func (sf *SendFactory) CreateMiscManager() *misc.MiscManager {
	return misc.NewMiscManager(sf.baseSend)
}

// CreateInfoChatManager creates a new info chat manager
func (sf *SendFactory) CreateInfoChatManager() *chat.InfoChatManager {
	return chat.NewInfoChatManager(sf.baseSend)
}

// FullMockSend is a mock implementation of the core.Send interface for testing
// that includes all packet IDs for all systems
type FullMockSend struct {
	packetIDs      map[string]string
	reconstructed  []byte
	sent           []byte
	time           uint32
	lastPacketName string
	lastArgs       map[string]interface{}
}

// NewFullMockSend creates a new FullMockSend instance with default values
func NewFullMockSend() *FullMockSend {
	return &FullMockSend{
		packetIDs: map[string]string{
			// Cash shop packet IDs
			"request_cashitems":  "0A6A",
			"cash_shop_open":     "0844",
			"cash_shop_close":    "0845",
			"cash_shop_buy":      "0288",
			"cash_dealer_buy":    "0848",
			"merge_item_request": "0972",
			"merge_item_cancel":  "0974",

			// Misc packet IDs
			"token_login":          "0825",
			"request_remain_time":  "0A37",
			"blocking_play_cancel": "0447",
			"recall_sso":           "0842",
			"remove_aid_sso":       "0843",
			"starplace_agree":      "0B0D",
			"sync_request_ex":      "09F1",

			// Chat info packet IDs
			"actor_info_request": "0A90",
			"actor_name_request": "0095",
			"request_user_count": "01C1",
			"battleground_chat":  "02DB",
			"clan_chat":          "0B01",
		},
		time:     12345,
		lastArgs: make(map[string]interface{}),
	}
}

// SendToServer mocks sending a packet to the server
func (ms *FullMockSend) SendToServer(msg []byte) error {
	ms.sent = msg
	return nil
}

// EncryptMessageID mocks encrypting a message ID
func (ms *FullMockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys mocks setting encryption keys
func (ms *FullMockSend) CryptKeys(key1, key2, key3 uint32) {}

// PinEncode mocks encoding a PIN
func (ms *FullMockSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage mocks injecting a message
func (ms *FullMockSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage mocks injecting an admin message
func (ms *FullMockSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw mocks sending a raw packet
func (ms *FullMockSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct mocks reconstructing a packet
func (ms *FullMockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the last packet name and arguments for testing
	for name, id := range ms.packetIDs {
		if id == packetID {
			ms.lastPacketName = name
			break
		}
	}

	// Store the arguments for testing
	ms.lastArgs = args

	// Simple mock implementation that just returns the packet ID as bytes
	ms.reconstructed = []byte{0x00, 0x00}
	return ms.reconstructed, nil
}

// GetPacketID mocks getting a packet ID by name
func (ms *FullMockSend) GetPacketID(name string) (string, bool) {
	id, ok := ms.packetIDs[name]
	if ok {
		ms.lastPacketName = name
	}
	return id, ok
}

// RegisterPacketHandler mocks registering a packet handler
func (ms *FullMockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
}

// RegisterHook mocks registering a hook
func (ms *FullMockSend) RegisterHook(hookName string, callback hooks.HookCallback) {}

// SetConnection mocks setting a connection
func (ms *FullMockSend) SetConnection(conn interface{}) {}

// GetConnection mocks getting a connection
func (ms *FullMockSend) GetConnection() interface{} {
	return nil
}

// GetTime mocks getting the current time
func (ms *FullMockSend) GetTime() uint32 {
	return ms.time
}

// LastPacketID returns the name of the last packet that was requested
func (ms *FullMockSend) LastPacketID() (string, bool) {
	if ms.lastPacketName == "" {
		return "", false
	}
	return ms.lastPacketName, true
}

// LastArgs returns the arguments of the last packet that was reconstructed
func (ms *FullMockSend) LastArgs() map[string]interface{} {
	return ms.lastArgs
}

// TestFactoryIntegration tests the integration of the factory with all systems
func TestFactoryIntegration(t *testing.T) {
	// Verify that FullMockSend implements core.Send
	var _ core.Send = &FullMockSend{}
	mockSend := NewFullMockSend()

	// Create the factory
	factory := NewSendFactory(mockSend)

	// Test the cash shop manager
	cashShopManager := factory.CreateCashShopManager()
	err := cashShopManager.OpenShop()
	if err != nil {
		t.Fatalf("OpenShop() returned error: %v", err)
	}
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_open" {
		t.Errorf("Expected packet ID 'cash_shop_open', got '%s'", packetID)
	}

	// Test the misc manager
	miscManager := factory.CreateMiscManager()
	err = miscManager.SendReqRemainTime()
	if err != nil {
		t.Fatalf("SendReqRemainTime() returned error: %v", err)
	}
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_remain_time" {
		t.Errorf("Expected packet ID 'request_remain_time', got '%s'", packetID)
	}

	// Test the info chat manager
	infoChatManager := factory.CreateInfoChatManager()
	err = infoChatManager.SendWho()
	if err != nil {
		t.Fatalf("SendWho() returned error: %v", err)
	}
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "request_user_count" {
		t.Errorf("Expected packet ID 'request_user_count', got '%s'", packetID)
	}

	// Test a complex scenario
	// 1. Open the cash shop
	err = cashShopManager.OpenShop()
	if err != nil {
		t.Fatalf("OpenShop() returned error: %v", err)
	}

	// 2. Request cash items list
	err = cashShopManager.RequestPoints()
	if err != nil {
		t.Fatalf("RequestPoints() returned error: %v", err)
	}

	// 3. Send a clan chat message
	message := "I'm in the cash shop!"
	charName := "TestChar"
	err = infoChatManager.SendClanChat(message, charName)
	if err != nil {
		t.Fatalf("SendClanChat() returned error: %v", err)
	}

	// 4. Check remaining time
	err = miscManager.SendReqRemainTime()
	if err != nil {
		t.Fatalf("SendReqRemainTime() returned error: %v", err)
	}

	// 5. Close the cash shop
	err = cashShopManager.CloseShop()
	if err != nil {
		t.Fatalf("CloseShop() returned error: %v", err)
	}

	// Verify the last packet was cash_shop_close
	if packetID, exists := mockSend.LastPacketID(); !exists || packetID != "cash_shop_close" {
		t.Errorf("Expected packet ID 'cash_shop_close', got '%s'", packetID)
	}
}

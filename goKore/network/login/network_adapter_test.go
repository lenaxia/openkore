package login

import (
	"testing"

	"github.com/lenaxia/goKore/network"
)

// TestNetworkInterface implements the network.NetworkInterface interface for testing
type TestNetworkInterface struct {
	connected bool
	state     int
}

func (m *TestNetworkInterface) Connect() error {
	m.connected = true
	m.state = 1 // ConnectedToMasterServer
	return nil
}

func (m *TestNetworkInterface) Disconnect() error {
	m.connected = false
	m.state = 0 // NotConnected
	return nil
}

func (m *TestNetworkInterface) IsConnected() bool {
	return m.connected
}

func (m *TestNetworkInterface) GetState() int {
	return m.state
}

func (m *TestNetworkInterface) SetState(state int) {
	m.state = state
}

func (m *TestNetworkInterface) Send(data []byte) error {
	return nil
}

func (m *TestNetworkInterface) Receive() ([]byte, error) {
	return []byte{}, nil
}

// TestPacketSender implements the network.PacketSender interface for testing
type TestPacketSender struct {
	lastPacketName string
	lastFields     map[string]interface{}
}

func (m *TestPacketSender) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	m.lastPacketName = packetName
	m.lastFields = fields
	return []byte{}, nil
}

func (m *TestPacketSender) GetCashShopManager() interface{} {
	return nil
}

func (m *TestPacketSender) GetMiscManager() interface{} {
	return nil
}

func (m *TestPacketSender) GetInfoChatManager() interface{} {
	return nil
}

// TestPacketHandler implements the network.PacketHandler interface for testing
type TestPacketHandler struct {
	lastPacket []byte
}

func (m *TestPacketHandler) Handle(packet []byte) error {
	m.lastPacket = packet
	return nil
}

// NetworkAdapter adapts the network.NetworkManager to implement the login.NetworkManager interface
type NetworkAdapter struct {
	*network.NetworkManager
}

// ConnectTo implements the login.NetworkManager interface
func (a *NetworkAdapter) ConnectTo(host string, port int) error {
	// In a real implementation, this would connect to the specified host and port
	return nil
}

// GetHookManager implements the login.NetworkManager interface
func (a *NetworkAdapter) GetHookManager() interface{} {
	// In a real implementation, this would return the hook manager
	return nil
}

// SetSessionStore implements the login.NetworkManager interface
func (a *NetworkAdapter) SetSessionStore(sessionStore *SessionStore) {
	// In a real implementation, this would set the session store
}

// TestNetworkManagerAdapter tests that the network.NetworkManager can be adapted to implement the login.NetworkManager interface
func TestNetworkManagerAdapter(t *testing.T) {
	// Create a mock network interface
	mockNetwork := &TestNetworkInterface{}

	// Create a mock packet sender
	mockSender := &TestPacketSender{}

	// Create a mock packet handler
	mockHandler := &TestPacketHandler{}

	// Create a network manager
	networkManager := network.NewNetworkManager(mockNetwork, mockSender, mockHandler)

	// Create a network adapter
	adapter := &NetworkAdapter{
		NetworkManager: networkManager,
	}

	// Create a login config
	config := NewLoginConfig("testuser", "testpass", "testserver")

	// Create a login manager with the adapter
	loginManager := NewLoginManager(adapter, config)

	// Verify that the login manager was created successfully
	if loginManager == nil {
		t.Error("Expected login manager to be created successfully")
	}

	// Test that the SetStateChangeCallback method works correctly
	var oldState, newState int

	// Reset the state to 0
	adapter.SetState(0)

	adapter.SetStateChangeCallback(func(old, new int) {
		oldState = old
		newState = new
	})

	// Change the state
	adapter.SetState(1) // ConnectedToMasterServer

	// Verify that the callback was called with the correct values
	if oldState != 0 || newState != 1 {
		t.Errorf("Expected oldState=0, newState=1, got oldState=%d, newState=%d", oldState, newState)
	}
}

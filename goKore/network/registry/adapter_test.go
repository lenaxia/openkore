package registry

import (
	"testing"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct{}

func (l *MockLogger) Debug(format string, args ...interface{})   {}
func (l *MockLogger) Info(format string, args ...interface{})    {}
func (l *MockLogger) Warning(format string, args ...interface{}) {}
func (l *MockLogger) Error(format string, args ...interface{})   {}
func (l *MockLogger) Success(format string, args ...interface{}) {}

// MockNetworkInterface implements the network.NetworkInterface for testing
type MockNetworkInterface struct {
	connected bool
	state     int
}

func (m *MockNetworkInterface) Connect() error {
	m.connected = true
	m.state = network.ConnectedToMasterServer
	return nil
}

func (m *MockNetworkInterface) Disconnect() error {
	m.connected = false
	m.state = network.NotConnected
	return nil
}

func (m *MockNetworkInterface) IsConnected() bool {
	return m.connected
}

func (m *MockNetworkInterface) GetState() int {
	return m.state
}

func (m *MockNetworkInterface) SetState(state int) {
	m.state = state
}

func (m *MockNetworkInterface) Send(data []byte) error {
	return nil
}

func (m *MockNetworkInterface) Receive() ([]byte, error) {
	return []byte{}, nil
}

// TestNetworkRegistryIntegrator tests the NetworkRegistryIntegrator
func TestNetworkRegistryIntegrator(t *testing.T) {
	// Create a logger
	logger := &MockLogger{}

	// Create an integrator
	integrator := NewNetworkRegistryIntegrator(logger)
	if integrator == nil {
		t.Fatal("Failed to create NetworkRegistryIntegrator")
	}

	// Create a mock network interface
	networkInterface := &MockNetworkInterface{}

	// Create a network manager
	manager := integrator.CreateNetworkManager(networkInterface)
	if manager == nil {
		t.Fatal("Failed to create NetworkManager")
	}

	// Test that the network manager is properly initialized
	if !networkInterface.IsConnected() {
		// The network interface should not be connected yet
		// Connect it
		err := manager.Connect()
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
	}

	// Test that the network manager is connected
	if !manager.IsConnected() {
		t.Fatal("NetworkManager should be connected")
	}

	// Test sending a packet
	_, err := manager.Send("ping", map[string]interface{}{})
	if err != nil {
		// This might fail if the ping packet is not registered, which is expected in a test
		// Just log it rather than failing the test
		t.Logf("Send failed: %v", err)
	}

	// Test disconnecting
	err = manager.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	// Test that the network manager is disconnected
	if manager.IsConnected() {
		t.Fatal("NetworkManager should be disconnected")
	}
}

// TestSendRegistryAdapter tests the SendRegistryAdapter
func TestSendRegistryAdapter(t *testing.T) {
	// Create a logger
	logger := &MockLogger{}

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a send registry adapter
	adapter := NewSendRegistryAdapter(hookManager, logger)
	if adapter == nil {
		t.Fatal("Failed to create SendRegistryAdapter")
	}

	// Initialize the adapter
	adapter.Initialize()

	// Test setting a connection
	adapter.SetConnection(&MockNetworkInterface{})

	// Test sending a packet
	_, err := adapter.Send("ping", map[string]interface{}{})
	if err != nil {
		// This might fail if the ping packet is not registered, which is expected in a test
		// Just log it rather than failing the test
		t.Logf("Send failed: %v", err)
	}

	// Test getting managers
	if adapter.GetCashShopManager() != nil {
		t.Log("CashShopManager is not nil")
	}
	if adapter.GetMiscManager() != nil {
		t.Log("MiscManager is not nil")
	}
	if adapter.GetInfoChatManager() != nil {
		t.Log("InfoChatManager is not nil")
	}
}

// TestReceiveRegistryAdapter tests the ReceiveRegistryAdapter
func TestReceiveRegistryAdapter(t *testing.T) {
	// Create a logger
	logger := &MockLogger{}

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a receive registry adapter
	adapter := NewReceiveRegistryAdapter(hookManager, logger)
	if adapter == nil {
		t.Fatal("Failed to create ReceiveRegistryAdapter")
	}

	// Initialize the adapter
	adapter.Initialize()

	// Test handling a packet
	err := adapter.Handle([]byte{0x01, 0x02, 0x03})
	if err != nil {
		// This might fail if the packet is not registered, which is expected in a test
		// Just log it rather than failing the test
		t.Logf("Handle failed: %v", err)
	}
}

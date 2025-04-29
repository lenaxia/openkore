package network

import (
	"testing"
)

// MockPacketSender implements the PacketSender interface for testing
type MockPacketSender struct {
	lastPacketName string
	lastFields     map[string]interface{}
}

func (m *MockPacketSender) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	m.lastPacketName = packetName
	m.lastFields = fields
	return []byte{}, nil
}

func (m *MockPacketSender) GetCashShopManager() interface{} {
	return "CashShopManager"
}

func (m *MockPacketSender) GetMiscManager() interface{} {
	return "MiscManager"
}

func (m *MockPacketSender) GetInfoChatManager() interface{} {
	return "InfoChatManager"
}

// MockPacketHandler implements the PacketHandler interface for testing
type MockPacketHandler struct {
	lastPacket []byte
}

func (m *MockPacketHandler) Handle(packet []byte) error {
	m.lastPacket = packet
	return nil
}

// TestNetworkManager tests the NetworkManager implementation
func TestNetworkManager(t *testing.T) {
	// Create mock implementations
	mockNetwork := &MockNetwork{}
	mockSender := &MockPacketSender{}
	mockHandler := &MockPacketHandler{}

	// Create network manager
	manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

	// Test initial state
	if manager.GetState() != NotConnected {
		t.Errorf("Expected initial state %d, got %d", NotConnected, manager.GetState())
	}

	// Test connect
	err := manager.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	if !mockNetwork.connected {
		t.Error("Network should be connected after Connect()")
	}

	// Test state change callback
	var oldState, newState int
	manager.SetStateChangeCallback(func(old, new int) {
		oldState = old
		newState = new
	})

	// Test set state
	manager.SetState(ConnectedToLoginServer)
	if manager.GetState() != ConnectedToLoginServer {
		t.Errorf("Expected state %d, got %d", ConnectedToLoginServer, manager.GetState())
	}
	if mockNetwork.state != ConnectedToLoginServer {
		t.Errorf("Expected network state %d, got %d", ConnectedToLoginServer, mockNetwork.state)
	}
	if oldState != NotConnected || newState != ConnectedToLoginServer {
		t.Errorf("State change callback received incorrect values: old=%d, new=%d", oldState, newState)
	}

	// Test send
	_, err = manager.Send("test_packet", map[string]interface{}{"key": "value"})
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
	if mockSender.lastPacketName != "test_packet" {
		t.Errorf("Expected packet name 'test_packet', got '%s'", mockSender.lastPacketName)
	}
	if val, ok := mockSender.lastFields["key"]; !ok || val != "value" {
		t.Errorf("Expected field 'key' with value 'value', got %v", mockSender.lastFields)
	}

	// Test get managers
	if manager.GetCashShopManager() != "CashShopManager" {
		t.Errorf("Expected CashShopManager, got %v", manager.GetCashShopManager())
	}
	if manager.GetMiscManager() != "MiscManager" {
		t.Errorf("Expected MiscManager, got %v", manager.GetMiscManager())
	}
	if manager.GetInfoChatManager() != "InfoChatManager" {
		t.Errorf("Expected InfoChatManager, got %v", manager.GetInfoChatManager())
	}

	// Test handle packet
	err = manager.HandlePacket([]byte{1, 2, 3})
	if err != nil {
		t.Errorf("HandlePacket failed: %v", err)
	}
	if len(mockHandler.lastPacket) != 3 || mockHandler.lastPacket[0] != 1 ||
		mockHandler.lastPacket[1] != 2 || mockHandler.lastPacket[2] != 3 {
		t.Errorf("Expected packet [1, 2, 3], got %v", mockHandler.lastPacket)
	}

	// Test disconnect
	err = manager.Disconnect()
	if err != nil {
		t.Errorf("Disconnect failed: %v", err)
	}
	if mockNetwork.connected {
		t.Error("Network should not be connected after Disconnect()")
	}
}

// MockErrorNetwork implements NetworkInterface for testing error cases
type MockErrorNetwork struct {
	MockNetwork
	shouldFailConnect    bool
	shouldFailDisconnect bool
	shouldFailSend       bool
	shouldFailReceive    bool
}

func (m *MockErrorNetwork) Connect() error {
	if m.shouldFailConnect {
		return ErrTimeout
	}
	return m.MockNetwork.Connect()
}

func (m *MockErrorNetwork) Disconnect() error {
	if m.shouldFailDisconnect {
		return ErrConnectionClosed
	}
	return m.MockNetwork.Disconnect()
}

func (m *MockErrorNetwork) Send(data []byte) error {
	if m.shouldFailSend {
		return ErrPacketTooLarge
	}
	return m.MockNetwork.Send(data)
}

func (m *MockErrorNetwork) Receive() ([]byte, error) {
	if m.shouldFailReceive {
		return nil, ErrInvalidPacket
	}
	return m.MockNetwork.Receive()
}

// MockErrorPacketSender implements the PacketSender interface for testing error cases
type MockErrorPacketSender struct {
	MockPacketSender
	shouldFailSend bool
}

func (m *MockErrorPacketSender) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	if m.shouldFailSend {
		return nil, ErrInvalidPacket
	}
	return m.MockPacketSender.Send(packetName, fields)
}

// MockErrorPacketHandler implements the PacketHandler interface for testing error cases
type MockErrorPacketHandler struct {
	MockPacketHandler
	shouldFailHandle bool
}

func (m *MockErrorPacketHandler) Handle(packet []byte) error {
	if m.shouldFailHandle {
		return ErrInvalidPacket
	}
	return m.MockPacketHandler.Handle(packet)
}

// TestNetworkManagerHappyPath tests the NetworkManager implementation with multiple happy paths
func TestNetworkManagerHappyPath(t *testing.T) {
	// Create mock implementations
	mockNetwork := &MockNetwork{}
	mockSender := &MockPacketSender{}
	mockHandler := &MockPacketHandler{}

	// Create network manager
	manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

	// Test initial state
	if manager.GetState() != NotConnected {
		t.Errorf("Expected initial state %d, got %d", NotConnected, manager.GetState())
	}

	// Test connect
	err := manager.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}
	if !mockNetwork.connected {
		t.Error("Network should be connected after Connect()")
	}

	// Test multiple state changes
	states := []int{
		ConnectedToMasterServer,
		ConnectedToLoginServer,
		ConnectedToCharServer,
		InGame,
	}

	for _, state := range states {
		manager.SetState(state)
		if manager.GetState() != state {
			t.Errorf("Expected state %d, got %d", state, manager.GetState())
		}
		if mockNetwork.state != state {
			t.Errorf("Expected network state %d, got %d", state, mockNetwork.state)
		}
	}

	// Test sending multiple packets
	packets := []struct {
		name   string
		fields map[string]interface{}
	}{
		{"packet1", map[string]interface{}{"key1": "value1"}},
		{"packet2", map[string]interface{}{"key2": 123}},
		{"packet3", map[string]interface{}{"key3": true}},
	}

	for _, packet := range packets {
		_, err = manager.Send(packet.name, packet.fields)
		if err != nil {
			t.Errorf("Send failed for packet %s: %v", packet.name, err)
		}
		if mockSender.lastPacketName != packet.name {
			t.Errorf("Expected packet name '%s', got '%s'", packet.name, mockSender.lastPacketName)
		}
		for k, v := range packet.fields {
			if val, ok := mockSender.lastFields[k]; !ok || val != v {
				t.Errorf("Expected field '%s' with value '%v', got '%v'", k, v, val)
			}
		}
	}

	// Test handling multiple packets
	testPackets := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	for _, packet := range testPackets {
		err = manager.HandlePacket(packet)
		if err != nil {
			t.Errorf("HandlePacket failed: %v", err)
		}
		if len(mockHandler.lastPacket) != len(packet) {
			t.Errorf("Expected packet length %d, got %d", len(packet), len(mockHandler.lastPacket))
		}
		for i, b := range packet {
			if mockHandler.lastPacket[i] != b {
				t.Errorf("Expected packet byte %d to be %d, got %d", i, b, mockHandler.lastPacket[i])
			}
		}
	}

	// Test disconnect
	err = manager.Disconnect()
	if err != nil {
		t.Errorf("Disconnect failed: %v", err)
	}
	if mockNetwork.connected {
		t.Error("Network should not be connected after Disconnect()")
	}
}

// TestNetworkManagerUnhappyPath tests the NetworkManager implementation with multiple unhappy paths
func TestNetworkManagerUnhappyPath(t *testing.T) {
	// Test connect failure
	{
		mockNetwork := &MockErrorNetwork{shouldFailConnect: true}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		err := manager.Connect()
		if err != ErrTimeout {
			t.Errorf("Expected error %v, got %v", ErrTimeout, err)
		}
		if mockNetwork.connected {
			t.Error("Network should not be connected after failed Connect()")
		}
	}

	// Test disconnect failure
	{
		mockNetwork := &MockErrorNetwork{shouldFailDisconnect: true}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		// Connect first
		mockNetwork.MockNetwork.Connect()

		err := manager.Disconnect()
		if err != ErrConnectionClosed {
			t.Errorf("Expected error %v, got %v", ErrConnectionClosed, err)
		}
	}

	// Test send failure (network)
	{
		mockNetwork := &MockErrorNetwork{shouldFailSend: true}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		// Connect first
		mockNetwork.MockNetwork.Connect()

		// This should still succeed because we're using the packet sender
		_, err := manager.Send("test_packet", nil)
		if err != nil {
			t.Errorf("Send should not fail when only network.Send fails: %v", err)
		}
	}

	// Test send failure (packet sender)
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockErrorPacketSender{shouldFailSend: true}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		// Connect first
		mockNetwork.Connect()

		_, err := manager.Send("test_packet", nil)
		if err != ErrInvalidPacket {
			t.Errorf("Expected error %v, got %v", ErrInvalidPacket, err)
		}
	}

	// Test handle packet failure
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockPacketSender{}
		mockHandler := &MockErrorPacketHandler{shouldFailHandle: true}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		err := manager.HandlePacket([]byte{1, 2, 3})
		if err != ErrInvalidPacket {
			t.Errorf("Expected error %v, got %v", ErrInvalidPacket, err)
		}
	}

	// Test state change callback with invalid state
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		var oldState, newState int
		manager.SetStateChangeCallback(func(old, new int) {
			oldState = old
			newState = new
		})

		// Set an invalid state
		invalidState := 999
		manager.SetState(invalidState)

		if manager.GetState() != invalidState {
			t.Errorf("Expected state %d, got %d", invalidState, manager.GetState())
		}

		if oldState != NotConnected || newState != invalidState {
			t.Errorf("State change callback received incorrect values: old=%d, new=%d", oldState, newState)
		}
	}
}

// TestNetworkManagerEdgeCases tests the NetworkManager implementation with edge cases
func TestNetworkManagerEdgeCases(t *testing.T) {
	// Test nil network interface
	{
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when network interface is nil")
			}
		}()

		_ = NewNetworkManager(nil, &MockPacketSender{}, &MockPacketHandler{})
	}

	// Test nil packet sender
	{
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when packet sender is nil")
			}
		}()

		_ = NewNetworkManager(&MockNetwork{}, nil, &MockPacketHandler{})
	}

	// Test nil packet handler
	{
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when packet handler is nil")
			}
		}()

		_ = NewNetworkManager(&MockNetwork{}, &MockPacketSender{}, nil)
	}

	// Test setting nil state change callback
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		// Set a callback
		manager.SetStateChangeCallback(func(old, new int) {})

		// Set nil callback
		manager.SetStateChangeCallback(nil)

		// Change state (should not panic)
		manager.SetState(ConnectedToMasterServer)
	}

	// Test sending empty packet name
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		_, err := manager.Send("", nil)
		if err != nil {
			t.Errorf("Send with empty packet name should not fail: %v", err)
		}

		if mockSender.lastPacketName != "" {
			t.Errorf("Expected empty packet name, got '%s'", mockSender.lastPacketName)
		}
	}

	// Test handling nil packet
	{
		mockNetwork := &MockNetwork{}
		mockSender := &MockPacketSender{}
		mockHandler := &MockPacketHandler{}

		manager := NewNetworkManager(mockNetwork, mockSender, mockHandler)

		err := manager.HandlePacket(nil)
		if err != nil {
			t.Errorf("HandlePacket with nil packet should not fail: %v", err)
		}

		if mockHandler.lastPacket != nil {
			t.Errorf("Expected nil packet, got %v", mockHandler.lastPacket)
		}
	}
}

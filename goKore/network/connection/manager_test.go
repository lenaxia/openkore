package connection

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

// mockConnection implements Connection interface for testing
type mockConnection struct {
	BaseConnection
	connectCalled      bool
	disconnectCalled   bool
	sendCalled         bool
	receiveCalled      bool
	connectError       error
	disconnectError    error
	sendError          error
	receiveError       error
	receiveData        []byte
	connectWithContext func(ctx context.Context) error
}

func newMockConnection() *mockConnection {
	return &mockConnection{
		BaseConnection: *NewBaseConnection(&ConnectionConfig{
			Host: "localhost",
			Port: 6900,
		}),
		receiveData: []byte{1, 2, 3, 4},
	}
}

func (m *mockConnection) Connect() error {
	m.connectCalled = true
	if m.connectError != nil {
		return m.connectError
	}
	m.SetState(CONNECTED_TO_MASTER_SERVER)
	return nil
}

func (m *mockConnection) Disconnect() error {
	m.disconnectCalled = true
	if m.disconnectError != nil {
		return m.disconnectError
	}
	m.SetState(NOT_CONNECTED)
	return nil
}

func (m *mockConnection) Send(data []byte) error {
	m.sendCalled = true
	if m.sendError != nil {
		return m.sendError
	}
	return nil
}

func (m *mockConnection) Receive() ([]byte, error) {
	m.receiveCalled = true
	if m.receiveError != nil {
		return nil, m.receiveError
	}
	return m.receiveData, nil
}

func (m *mockConnection) IsConnected() bool {
	return m.GetState() != NOT_CONNECTED
}

func (m *mockConnection) ConnectWithContext(ctx context.Context) error {
	if m.connectWithContext != nil {
		return m.connectWithContext(ctx)
	}
	return m.Connect()
}

func (m *mockConnection) SendWithContext(ctx context.Context, data []byte) error {
	return m.Send(data)
}

func (m *mockConnection) ReceiveWithContext(ctx context.Context) ([]byte, error) {
	return m.Receive()
}

func (m *mockConnection) GetRemoteAddress() net.Addr {
	return nil
}

func (m *mockConnection) GetLocalAddress() net.Addr {
	return nil
}

// TestNewConnectionManager tests the creation of a new connection manager
func TestNewConnectionManager(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	if manager == nil {
		t.Fatal("NewConnectionManager() returned nil")
	}

	if manager.connection != conn {
		t.Error("Connection not properly set in manager")
	}

	if manager.reconnectAttempts != 0 {
		t.Errorf("Expected reconnectAttempts to be 0, got %d", manager.reconnectAttempts)
	}

	if manager.maxReconnectAttempts != DefaultMaxReconnectAttempts {
		t.Errorf("Expected maxReconnectAttempts to be %d, got %d", DefaultMaxReconnectAttempts, manager.maxReconnectAttempts)
	}

	if manager.reconnectDelay != DefaultReconnectDelay {
		t.Errorf("Expected reconnectDelay to be %v, got %v", DefaultReconnectDelay, manager.reconnectDelay)
	}
}

// TestConnectionManagerConnect tests the Connect method
func TestConnectionManagerConnect(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	err := manager.Connect()
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	if !conn.connectCalled {
		t.Error("Connection.Connect() was not called")
	}

	if manager.reconnectAttempts != 0 {
		t.Errorf("Expected reconnectAttempts to be reset to 0, got %d", manager.reconnectAttempts)
	}
}

// TestConnectionManagerConnectError tests the Connect method with an error
func TestConnectionManagerConnectError(t *testing.T) {
	conn := newMockConnection()
	conn.connectError = errors.New("connect error")
	manager := NewConnectionManager(conn)

	err := manager.Connect()
	if err == nil {
		t.Error("Expected Connect to return an error")
	}

	if !conn.connectCalled {
		t.Error("Connection.Connect() was not called")
	}
}

// TestConnectionManagerDisconnect tests the Disconnect method
func TestConnectionManagerDisconnect(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	err := manager.Disconnect()
	if err != nil {
		t.Errorf("Disconnect failed: %v", err)
	}

	if !conn.disconnectCalled {
		t.Error("Connection.Disconnect() was not called")
	}
}

// TestConnectionManagerDisconnectError tests the Disconnect method with an error
func TestConnectionManagerDisconnectError(t *testing.T) {
	conn := newMockConnection()
	conn.disconnectError = errors.New("disconnect error")
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	err := manager.Disconnect()
	if err == nil {
		t.Error("Expected Disconnect to return an error")
	}

	if !conn.disconnectCalled {
		t.Error("Connection.Disconnect() was not called")
	}
}

// TestConnectionManagerSend tests the Send method
func TestConnectionManagerSend(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	data := []byte{1, 2, 3, 4}
	err := manager.Send(data)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	if !conn.sendCalled {
		t.Error("Connection.Send() was not called")
	}
}

// TestConnectionManagerSendError tests the Send method with an error
func TestConnectionManagerSendError(t *testing.T) {
	conn := newMockConnection()
	conn.sendError = errors.New("send error")
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	data := []byte{1, 2, 3, 4}
	err := manager.Send(data)
	if err == nil {
		t.Error("Expected Send to return an error")
	}

	if !conn.sendCalled {
		t.Error("Connection.Send() was not called")
	}
}

// TestConnectionManagerReceive tests the Receive method
func TestConnectionManagerReceive(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	data, err := manager.Receive()
	if err != nil {
		t.Errorf("Receive failed: %v", err)
	}

	if !conn.receiveCalled {
		t.Error("Connection.Receive() was not called")
	}

	if len(data) != 4 {
		t.Errorf("Expected data length 4, got %d", len(data))
	}
}

// TestConnectionManagerReceiveError tests the Receive method with an error
func TestConnectionManagerReceiveError(t *testing.T) {
	conn := newMockConnection()
	conn.receiveError = errors.New("receive error")
	manager := NewConnectionManager(conn)

	// First connect
	manager.Connect()

	_, err := manager.Receive()
	if err == nil {
		t.Error("Expected Receive to return an error")
	}

	if !conn.receiveCalled {
		t.Error("Connection.Receive() was not called")
	}
}

// TestConnectionManagerReconnect tests the reconnect functionality
func TestConnectionManagerReconnect(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)
	manager.reconnectDelay = 10 * time.Millisecond // Short delay for testing

	// First connect
	manager.Connect()

	// Simulate a disconnect
	conn.SetState(NOT_CONNECTED)

	// Try to reconnect
	err := manager.Reconnect()
	if err != nil {
		t.Errorf("Reconnect failed: %v", err)
	}

	if manager.reconnectAttempts != 1 {
		t.Errorf("Expected reconnectAttempts to be 1, got %d", manager.reconnectAttempts)
	}

	if !conn.connectCalled {
		t.Error("Connection.Connect() was not called during reconnect")
	}
}

// TestConnectionManagerReconnectMaxAttempts tests the reconnect functionality with max attempts
func TestConnectionManagerReconnectMaxAttempts(t *testing.T) {
	conn := newMockConnection()
	conn.connectError = errors.New("connect error")
	manager := NewConnectionManager(conn)
	manager.maxReconnectAttempts = 3
	manager.reconnectDelay = 10 * time.Millisecond // Short delay for testing

	// Try to reconnect multiple times
	for i := 0; i < manager.maxReconnectAttempts; i++ {
		err := manager.Reconnect()
		if err == nil {
			t.Error("Expected Reconnect to return an error")
		}
	}

	// One more attempt should fail with max attempts error
	err := manager.Reconnect()
	if err == nil || err.Error() != "maximum reconnection attempts reached" {
		t.Errorf("Expected max attempts error, got: %v", err)
	}

	if manager.reconnectAttempts != manager.maxReconnectAttempts {
		t.Errorf("Expected reconnectAttempts to be %d, got %d", manager.maxReconnectAttempts, manager.reconnectAttempts)
	}
}

// TestConnectionManagerCheckConnection tests the CheckConnection method
func TestConnectionManagerCheckConnection(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// Test with connected state
	conn.SetState(CONNECTED_TO_MASTER_SERVER)
	err := manager.CheckConnection()
	if err != nil {
		t.Errorf("CheckConnection failed: %v", err)
	}

	// Test with disconnected state
	conn.SetState(NOT_CONNECTED)
	err = manager.CheckConnection()
	if err != nil {
		t.Errorf("CheckConnection failed: %v", err)
	}

	// Should have attempted to reconnect
	if !conn.connectCalled {
		t.Error("Connection.Connect() was not called during CheckConnection")
	}
}

// TestConnectionManagerHandleConnectionError tests the HandleConnectionError method
func TestConnectionManagerHandleConnectionError(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// Test with a temporary error
	tempErr := &net.OpError{
		Op:  "read",
		Err: &net.DNSError{IsTemporary: true},
	}
	err := manager.HandleConnectionError(tempErr)
	if err != nil {
		t.Errorf("HandleConnectionError failed with temporary error: %v", err)
	}

	// Test with a non-temporary error
	nonTempErr := errors.New("non-temporary error")
	err = manager.HandleConnectionError(nonTempErr)
	if err == nil {
		t.Error("Expected HandleConnectionError to return an error for non-temporary error")
	}
}

// TestConnectionManagerWithContext tests the context-aware methods
func TestConnectionManagerWithContext(t *testing.T) {
	conn := newMockConnection()
	manager := NewConnectionManager(conn)

	// Test ConnectWithContext
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := manager.ConnectWithContext(ctx)
	if err != nil {
		t.Errorf("ConnectWithContext failed: %v", err)
	}

	// Test SendWithContext
	data := []byte{1, 2, 3, 4}
	err = manager.SendWithContext(ctx, data)
	if err != nil {
		t.Errorf("SendWithContext failed: %v", err)
	}

	// Test ReceiveWithContext
	_, err = manager.ReceiveWithContext(ctx)
	if err != nil {
		t.Errorf("ReceiveWithContext failed: %v", err)
	}
}

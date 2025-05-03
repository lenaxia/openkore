package utils

import (
	"fmt"
	"testing"
)

// MockConnection is a mock implementation of a network connection
type MockConnection struct {
	SentData     [][]byte
	ReceivedData [][]byte
	Connected    bool
	State        int
	ReturnError  error
}

// Connect implements the connection interface
func (m *MockConnection) Connect() error {
	m.Connected = true
	return m.ReturnError
}

// ConnectTo implements the connection interface
func (m *MockConnection) ConnectTo(host string, port int) error {
	m.Connected = true
	return m.ReturnError
}

// Disconnect implements the connection interface
func (m *MockConnection) Disconnect() error {
	m.Connected = false
	return m.ReturnError
}

// IsConnected implements the connection interface
func (m *MockConnection) IsConnected() bool {
	return m.Connected
}

// Send implements the connection interface
func (m *MockConnection) Send(data []byte) error {
	m.SentData = append(m.SentData, data)
	return m.ReturnError
}

// Receive implements the connection interface
func (m *MockConnection) Receive() ([]byte, error) {
	if len(m.ReceivedData) > 0 {
		data := m.ReceivedData[0]
		m.ReceivedData = m.ReceivedData[1:]
		return data, nil
	}
	return nil, m.ReturnError
}

// GetState implements the connection interface
func (m *MockConnection) GetState() int {
	return m.State
}

// SetState implements the connection interface
func (m *MockConnection) SetState(state int) {
	m.State = state
}

// QueueReceiveData adds data to the receive queue
func (m *MockConnection) QueueReceiveData(data []byte) {
	m.ReceivedData = append(m.ReceivedData, data)
}

// GetLastSentData returns the last sent data
func (m *MockConnection) GetLastSentData() []byte {
	if len(m.SentData) == 0 {
		return nil
	}
	return m.SentData[len(m.SentData)-1]
}

// AssertPacketEquals compares two packets with tolerance for variable fields
func AssertPacketEquals(t *testing.T, expected, actual []byte, variableFields map[string]bool) {
	if len(expected) != len(actual) {
		t.Errorf("Packet length mismatch: expected %d, got %d", len(expected), len(actual))
		return
	}

	// Compare fixed fields
	for i := 0; i < len(expected); i++ {
		// Skip variable fields
		fieldName := fmt.Sprintf("byte%d", i)
		if variableFields[fieldName] {
			continue
		}

		if expected[i] != actual[i] {
			t.Errorf("Packet mismatch at byte %d: expected 0x%02X, got 0x%02X", i, expected[i], actual[i])
		}
	}
}

// NewMockConnection creates a new mock connection
func NewMockConnection() *MockConnection {
	return &MockConnection{
		SentData:     make([][]byte, 0),
		ReceivedData: make([][]byte, 0),
		Connected:    false,
		State:        0,
		ReturnError:  nil,
	}
}

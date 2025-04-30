package mock

import (
	"sync"
)

// MockConnection is a mock implementation of a network connection for testing
type MockConnection struct {
	sentPackets [][]byte
	mu          sync.Mutex
}

// NewMockConnection creates a new MockConnection
func NewMockConnection() *MockConnection {
	return &MockConnection{
		sentPackets: make([][]byte, 0),
	}
}

// Send records a packet as sent
func (m *MockConnection) Send(packet []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Make a copy of the packet to avoid it being modified later
	packetCopy := make([]byte, len(packet))
	copy(packetCopy, packet)

	m.sentPackets = append(m.sentPackets, packetCopy)
	return nil
}

// GetSentPackets returns all packets that have been sent
func (m *MockConnection) GetSentPackets() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to avoid modification
	result := make([][]byte, len(m.sentPackets))
	for i, packet := range m.sentPackets {
		result[i] = make([]byte, len(packet))
		copy(result[i], packet)
	}

	return result
}

// ClearSentPackets clears the list of sent packets
func (m *MockConnection) ClearSentPackets() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sentPackets = make([][]byte, 0)
}

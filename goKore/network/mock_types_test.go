package network

// MockNetwork implements NetworkInterface for testing
type MockNetwork struct {
	connected bool
	state     int
}

func (m *MockNetwork) Connect() error {
	m.connected = true
	m.state = ConnectedToMasterServer
	return nil
}

func (m *MockNetwork) Disconnect() error {
	m.connected = false
	m.state = NotConnected
	return nil
}

func (m *MockNetwork) IsConnected() bool {
	return m.connected
}

func (m *MockNetwork) GetState() int {
	return m.state
}

func (m *MockNetwork) SetState(state int) {
	m.state = state
}

func (m *MockNetwork) Send(data []byte) error {
	if !m.connected {
		return ErrNotConnected
	}
	return nil
}

func (m *MockNetwork) Receive() ([]byte, error) {
	if !m.connected {
		return nil, ErrNotConnected
	}
	return []byte{}, nil
}

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

// MockErrorNetwork implements NetworkInterface for testing error cases
type MockErrorNetwork struct {
	connected            bool
	state                int
	shouldFailConnect    bool
	shouldFailDisconnect bool
	shouldFailSend       bool
	shouldFailReceive    bool
}

func (m *MockErrorNetwork) Connect() error {
	if m.shouldFailConnect {
		return ErrTimeout
	}
	m.connected = true
	m.state = ConnectedToMasterServer
	return nil
}

func (m *MockErrorNetwork) Disconnect() error {
	if m.shouldFailDisconnect {
		return ErrConnectionClosed
	}
	m.connected = false
	m.state = NotConnected
	return nil
}

func (m *MockErrorNetwork) IsConnected() bool {
	return m.connected
}

func (m *MockErrorNetwork) GetState() int {
	return m.state
}

func (m *MockErrorNetwork) SetState(state int) {
	m.state = state
}

func (m *MockErrorNetwork) Send(data []byte) error {
	if m.shouldFailSend {
		return ErrPacketTooLarge
	}
	if !m.connected {
		return ErrNotConnected
	}
	return nil
}

func (m *MockErrorNetwork) Receive() ([]byte, error) {
	if m.shouldFailReceive {
		return nil, ErrInvalidPacket
	}
	if !m.connected {
		return nil, ErrNotConnected
	}
	return []byte{}, nil
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

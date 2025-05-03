package login

// TestNetworkManager defines the interface for a test network manager
// This interface extends the NetworkManager interface with additional methods for testing
type TestNetworkManager interface {
	NetworkManager

	// Test-specific methods
	SimulateReceivePacket(packetType string, data []byte) error
	SetConnectError(err error)
	SetSendError(err error)
	SetHandleError(err error)
	GetLastSentPacket() (string, map[string]interface{})
	GetLastReceivedPacket() []byte
	GetHookManager() interface{}
	CallHook(hookName string, arg interface{})
	SetSessionData(sessionData SessionData)
	GetSessionData() SessionData
	SetServerInfo(servers []ServerInfo)
}

// Ensure MockNetworkManager implements TestNetworkManager
// Add missing methods to MockNetworkManager

// SetConnectError sets an error to be returned on connect
func (m *MockNetworkManager) SetConnectError(err error) {
	m.networkInterface.SetConnectError(err)
}

// SetSendError sets an error to be returned on send
func (m *MockNetworkManager) SetSendError(err error) {
	m.packetSender.SetSendError(err)
}

// SetHandleError sets an error to be returned on handle
func (m *MockNetworkManager) SetHandleError(err error) {
	m.packetHandler.SetHandleError(err)
}

// SetServerInfo sets the server info
func (m *MockNetworkManager) SetServerInfo(servers []ServerInfo) {
	m.sessionStore.SetServerInfo(servers)
}

// Now we can verify that MockNetworkManager implements TestNetworkManager
var _ TestNetworkManager = (*MockNetworkManager)(nil)

// Add GetLastPacket method to MockPacketSender
func (m *MockPacketSender) GetLastPacket() (string, map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return the last packet name and an empty map
	// In a real implementation, this would return the actual fields
	for name := range m.packets {
		return name, make(map[string]interface{})
	}

	return "", make(map[string]interface{})
}

// Add UpdateFromSessionData method to SessionStore
func (ss *SessionStore) UpdateFromSessionData(data SessionData) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// Copy session ID
	if data.SessionID != nil {
		ss.sessionID = make([]byte, len(data.SessionID))
		copy(ss.sessionID, data.SessionID)
	}

	// Copy account ID
	if data.AccountID != nil {
		ss.accountID = make([]byte, len(data.AccountID))
		copy(ss.accountID, data.AccountID)
	}

	// Copy session ID2
	if data.SessionID2 != nil {
		ss.sessionID2 = make([]byte, len(data.SessionID2))
		copy(ss.sessionID2, data.SessionID2)
	}

	// Copy char ID
	if data.CharID != nil {
		ss.charID = make([]byte, len(data.CharID))
		copy(ss.charID, data.CharID)
	}

	// Copy other fields
	ss.accountSex = data.AccountSex
	ss.mapName = data.MapName
	ss.mapIP = data.MapIP
	ss.mapPort = data.MapPort
}

// GetLastSentPacket gets the last sent packet
func (m *MockNetworkManager) GetLastSentPacket() (string, map[string]interface{}) {
	return m.packetSender.GetLastPacket()
}

// GetLastReceivedPacket gets the last received packet
func (m *MockNetworkManager) GetLastReceivedPacket() []byte {
	return m.packetHandler.GetLastPacket()
}

// CallHook calls a hook
func (m *MockNetworkManager) CallHook(hookName string, arg interface{}) {
	m.hookManager.CallHook(hookName, arg)
}

// SetSessionData sets the session data
func (m *MockNetworkManager) SetSessionData(sessionData SessionData) {
	m.sessionStore.UpdateFromSessionData(sessionData)
}

// GetSessionData gets the session data
func (m *MockNetworkManager) GetSessionData() SessionData {
	return m.sessionStore.GetSessionData()
}

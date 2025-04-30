package login

import (
	"errors"
	"log"
	"sync"
)

// MockNetworkInterface implements a mock network interface for testing
type MockNetworkInterface struct {
	connected      bool
	connectError   error
	disconnectChan chan struct{}
	receiveChan    chan []byte
	sendChan       chan []byte
	mu             sync.Mutex
	host           string
	port           int
}

// NewMockNetworkInterface creates a new mock network interface
func NewMockNetworkInterface() *MockNetworkInterface {
	return &MockNetworkInterface{
		connected:      false,
		disconnectChan: make(chan struct{}, 1),
		receiveChan:    make(chan []byte, 10),
		sendChan:       make(chan []byte, 10),
	}
}

// Connect simulates connecting to a server
func (m *MockNetworkInterface) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connectError != nil {
		return m.connectError
	}

	m.connected = true
	return nil
}

// ConnectTo simulates connecting to a specific server
func (m *MockNetworkInterface) ConnectTo(host string, port int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connectError != nil {
		return m.connectError
	}

	m.host = host
	m.port = port
	m.connected = true
	return nil
}

// Disconnect simulates disconnecting from the server
func (m *MockNetworkInterface) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return errors.New("not connected")
	}

	m.connected = false
	select {
	case m.disconnectChan <- struct{}{}:
	default:
	}

	return nil
}

// IsConnected returns whether the mock is connected
func (m *MockNetworkInterface) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.connected
}

// Send simulates sending data to the server
func (m *MockNetworkInterface) Send(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return errors.New("not connected")
	}

	select {
	case m.sendChan <- data:
	default:
		return errors.New("send buffer full")
	}

	return nil
}

// Receive simulates receiving data from the server
func (m *MockNetworkInterface) Receive() ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return nil, errors.New("not connected")
	}

	select {
	case data := <-m.receiveChan:
		return data, nil
	default:
		return nil, errors.New("no data available")
	}
}

// QueueReceiveData queues data to be received
func (m *MockNetworkInterface) QueueReceiveData(data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case m.receiveChan <- data:
	default:
		// Buffer full, discard data
	}
}

// GetLastSentData gets the last data sent
func (m *MockNetworkInterface) GetLastSentData() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case data := <-m.sendChan:
		return data
	default:
		return nil
	}
}

// SetConnectError sets an error to be returned on connect
func (m *MockNetworkInterface) SetConnectError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connectError = err
}

// MockPacketSender implements a mock packet sender for testing
type MockPacketSender struct {
	packets map[string][]byte
	mu      sync.Mutex
	sendErr error
}

// NewMockPacketSender creates a new mock packet sender
func NewMockPacketSender() *MockPacketSender {
	return &MockPacketSender{
		packets: make(map[string][]byte),
	}
}

// Send simulates sending a packet
func (m *MockPacketSender) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sendErr != nil {
		return nil, m.sendErr
	}

	// Create a simple mock packet
	packet := []byte("MOCK_" + packetName)
	m.packets[packetName] = packet

	return packet, nil
}

// GetCashShopManager returns nil for testing
func (m *MockPacketSender) GetCashShopManager() interface{} {
	return nil
}

// GetMiscManager returns nil for testing
func (m *MockPacketSender) GetMiscManager() interface{} {
	return nil
}

// GetInfoChatManager returns nil for testing
func (m *MockPacketSender) GetInfoChatManager() interface{} {
	return nil
}

// SetSendError sets an error to be returned on send
func (m *MockPacketSender) SetSendError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sendErr = err
}

// GetSentPacket gets a sent packet by name
func (m *MockPacketSender) GetSentPacket(packetName string) []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.packets[packetName]
}

// MockPacketHandler implements a mock packet handler for testing
type MockPacketHandler struct {
	mu           sync.Mutex
	handleErr    error
	lastPacket   []byte
	hookManager  *MockHookManager
	sessionStore *SessionStore // Add session store reference
}

// NewMockPacketHandler creates a new mock packet handler
func NewMockPacketHandler(hookManager *MockHookManager, sessionStore *SessionStore) *MockPacketHandler {
	return &MockPacketHandler{
		hookManager:  hookManager,
		sessionStore: sessionStore,
	}
}

// Handle simulates handling a packet
func (m *MockPacketHandler) Handle(packet []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.handleErr != nil {
		return m.handleErr
	}

	m.lastPacket = packet

	// Debug log the packet
	log.Printf("Mock: Handling packet: %v (string: %s)", packet, string(packet))

	// Extract the packet type from the beginning of the packet
	packetType := string(packet[:4]) // First 4 bytes are the packet type

	// Simulate packet handling based on packet type
	switch packetType {
	case "0069":
		log.Println("Mock: Handling account server info packet")
		// Account server info packet
		accountInfo := map[string]interface{}{
			"sessionID":  []byte{1, 2, 3, 4},
			"accountID":  []byte{5, 6, 7, 8},
			"sessionID2": []byte{9, 10, 11, 12},
			"accountSex": 1,
		}

		// Add server info to the session store
		if m.sessionStore != nil {
			m.sessionStore.SetServerInfo([]ServerInfo{
				{
					Name: "TestServer",
					IP:   "127.0.0.1",
					Port: 6900,
				},
			})
		}

		m.hookManager.CallHook("account_info_received", accountInfo)
	case "006B":
		log.Println("Mock: Handling character info packet")
		// Character info packet
		m.hookManager.CallHook("characters_info_received", map[string]interface{}{
			"total_slot": 3,
			"charInfo":   []byte{1, 2, 3, 4},
		})
	case "0071":
		log.Println("Mock: Handling character map info packet")
		// Character ID and map info packet
		m.hookManager.CallHook("character_map_info_received", map[string]interface{}{
			"charID":  []byte{1, 2, 3, 4},
			"mapName": "prontera",
			"mapIP":   "127.0.0.1",
			"mapPort": 6900,
		})
	case "0073":
		log.Println("Mock: Handling map loaded packet")
		// Map loaded packet
		m.hookManager.CallHook("map_loaded", map[string]interface{}{})
	case "006A":
		log.Println("Mock: Handling login error packet")
		// Login error packet
		m.hookManager.CallHook("login_error", map[string]interface{}{
			"type": 1,
			"date": "Error message",
		})
	default:
		log.Printf("Mock: Unknown packet type: %s", packetType)
	}

	return nil
}

// SetHandleError sets an error to be returned on handle
func (m *MockPacketHandler) SetHandleError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handleErr = err
}

// GetLastPacket gets the last handled packet
func (m *MockPacketHandler) GetLastPacket() []byte {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.lastPacket
}

// SetSessionStore sets the session store
func (m *MockPacketHandler) SetSessionStore(sessionStore *SessionStore) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessionStore = sessionStore
}

// MockHookManager implements a mock hook manager for testing
type MockHookManager struct {
	hooks map[string][]func(string, interface{}, interface{})
	mu    sync.Mutex
}

// NewMockHookManager creates a new mock hook manager
func NewMockHookManager() *MockHookManager {
	return &MockHookManager{
		hooks: make(map[string][]func(string, interface{}, interface{})),
	}
}

// Register registers a hook
func (m *MockHookManager) Register(hookName string, callback func(string, interface{}, interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.hooks[hookName]; !ok {
		m.hooks[hookName] = make([]func(string, interface{}, interface{}), 0)
	}

	m.hooks[hookName] = append(m.hooks[hookName], callback)
}

// Unregister unregisters a hook
func (m *MockHookManager) Unregister(hookName string, callback func(string, interface{}, interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if callbacks, ok := m.hooks[hookName]; ok {
		for i, cb := range callbacks {
			if &cb == &callback {
				m.hooks[hookName] = append(callbacks[:i], callbacks[i+1:]...)
				break
			}
		}
	}
}

// CallHook calls a hook
func (m *MockHookManager) CallHook(hookName string, arg interface{}) {
	m.mu.Lock()
	callbacks := m.hooks[hookName]
	m.mu.Unlock()

	log.Printf("Mock: Calling hook %s with %d callbacks", hookName, len(callbacks))

	for _, callback := range callbacks {
		callback(hookName, arg, nil)
	}
}

// MockNetworkManager implements a mock network manager for testing
type MockNetworkManager struct {
	networkInterface *MockNetworkInterface
	packetSender     *MockPacketSender
	packetHandler    *MockPacketHandler
	hookManager      *MockHookManager
	sessionStore     *SessionStore
	state            int
	stateChangeCb    func(oldState, newState int)
}

// NewMockNetworkManager creates a new mock network manager
func NewMockNetworkManager() *MockNetworkManager {
	hookManager := NewMockHookManager()
	sessionStore := NewSessionStore()
	networkInterface := NewMockNetworkInterface()
	packetSender := NewMockPacketSender()
	packetHandler := NewMockPacketHandler(hookManager, sessionStore)

	return &MockNetworkManager{
		networkInterface: networkInterface,
		packetSender:     packetSender,
		packetHandler:    packetHandler,
		hookManager:      hookManager,
		sessionStore:     sessionStore,
		state:            0, // NOT_CONNECTED
	}
}

// Connect simulates connecting to the server
func (m *MockNetworkManager) Connect() error {
	return m.networkInterface.Connect()
}

// ConnectTo simulates connecting to a specific server
func (m *MockNetworkManager) ConnectTo(host string, port int) error {
	return m.networkInterface.ConnectTo(host, port)
}

// Disconnect simulates disconnecting from the server
func (m *MockNetworkManager) Disconnect() error {
	return m.networkInterface.Disconnect()
}

// Send simulates sending a packet
func (m *MockNetworkManager) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	return m.packetSender.Send(packetName, fields)
}

// HandlePacket simulates handling a packet
func (m *MockNetworkManager) HandlePacket(packet []byte) error {
	return m.packetHandler.Handle(packet)
}

// SetState sets the network state
func (m *MockNetworkManager) SetState(state int) {
	oldState := m.state
	m.state = state

	if m.stateChangeCb != nil {
		m.stateChangeCb(oldState, state)
	}
}

// GetState gets the network state
func (m *MockNetworkManager) GetState() int {
	return m.state
}

// SetStateChangeCallback sets the state change callback
func (m *MockNetworkManager) SetStateChangeCallback(callback func(oldState, newState int)) {
	m.stateChangeCb = callback
}

// GetHookManager gets the hook manager
func (m *MockNetworkManager) GetHookManager() interface{} {
	return m.hookManager
}

// SetSessionStore sets the session store
func (m *MockNetworkManager) SetSessionStore(sessionStore *SessionStore) {
	m.sessionStore = sessionStore
	m.packetHandler.SetSessionStore(sessionStore)
}

// SimulateReceivePacket simulates receiving a packet
func (m *MockNetworkManager) SimulateReceivePacket(packetType string, data []byte) error {
	log.Printf("Mock: Simulating receive packet %s", packetType)
	packet := append([]byte(packetType), '_')
	packet = append(packet, data...)
	return m.packetHandler.Handle(packet)
}

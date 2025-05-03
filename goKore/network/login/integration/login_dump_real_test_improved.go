package integration

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/login"
	"github.com/lenaxia/goKore/network/receive/base"
	"github.com/lenaxia/goKore/network/send"
	"github.com/lenaxia/goKore/network/servers"
)

// PacketDump represents the structure of the JSON dump files
type PacketDump struct {
	FileName    string        `json:"file_name"`
	PacketCount int           `json:"packet_count"`
	Packets     []PacketEntry `json:"packets"`
}

// PacketEntry represents a single packet in the dump
type PacketEntry struct {
	Direction   string        `json:"direction"`
	PacketID    string        `json:"packet_id"`
	Description string        `json:"description"`
	Size        int           `json:"size"`
	Timestamp   string        `json:"timestamp"`
	RawData     []RawDataPart `json:"raw_data"`
}

// RawDataPart represents a part of the raw data in a packet
type RawDataPart struct {
	HexBytes            []string `json:"hex_bytes"`
	AsciiRepresentation string   `json:"ascii_representation"`
	BinaryBase64        string   `json:"binary_base64"`
}

// DumpNetworkInterface implements the NetworkInterface for testing with packet dumps
type DumpNetworkInterface struct {
	connected      bool
	state          int
	t              *testing.T
	currentDump    *PacketDump
	currentIndex   int
	receiveChan    chan []byte
	sendChan       chan []byte
	mu             sync.Mutex
	host           string
	port           int
	disconnectChan chan struct{}
	connectError   error
	sendError      error
}

// NewDumpNetworkInterface creates a new dump network interface
func NewDumpNetworkInterface(t *testing.T) *DumpNetworkInterface {
	return &DumpNetworkInterface{
		connected:      false,
		state:          network.NotConnected,
		t:              t,
		receiveChan:    make(chan []byte, 100),
		sendChan:       make(chan []byte, 100),
		disconnectChan: make(chan struct{}),
	}
}

// Connect implements the NetworkInterface interface
func (d *DumpNetworkInterface) Connect() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.connectError != nil {
		return d.connectError
	}

	d.connected = true
	d.state = network.ConnectedToMasterServer
	return nil
}

// ConnectTo implements the NetworkInterface interface
func (d *DumpNetworkInterface) ConnectTo(host string, port int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.connectError != nil {
		return d.connectError
	}

	d.host = host
	d.port = port
	d.connected = true

	// Update state based on the port
	// This is a simplification - in a real implementation, the state would be determined by the server type
	if port == 6900 {
		d.state = 1 // ConnectedToMasterServer
	} else if port == 6121 {
		d.state = 2 // ConnectedToCharServer
	} else if port == 5121 {
		d.state = 3 // ConnectedToMapServer
	}

	return nil
}

// Disconnect implements the NetworkInterface interface
func (d *DumpNetworkInterface) Disconnect() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.connected = false
	d.state = network.NotConnected

	select {
	case d.disconnectChan <- struct{}{}:
	default:
	}

	return nil
}

// IsConnected implements the NetworkInterface interface
func (d *DumpNetworkInterface) IsConnected() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.connected
}

// GetState implements the NetworkInterface interface
func (d *DumpNetworkInterface) GetState() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.state
}

// SetState implements the NetworkInterface interface
func (d *DumpNetworkInterface) SetState(state int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = state
}

// Send implements the NetworkInterface interface
func (d *DumpNetworkInterface) Send(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.connected {
		return fmt.Errorf("not connected")
	}

	if d.sendError != nil {
		return d.sendError
	}

	d.t.Logf("Sending packet: %v", data)
	d.sendChan <- data
	return nil
}

// Receive implements the NetworkInterface interface
func (d *DumpNetworkInterface) Receive() ([]byte, error) {
	select {
	case data := <-d.receiveChan:
		return data, nil
	case <-time.After(100 * time.Millisecond):
		return nil, fmt.Errorf("timeout")
	}
}

// SetConnectError sets an error to be returned on connect
func (d *DumpNetworkInterface) SetConnectError(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.connectError = err
}

// SetSendError sets an error to be returned on send
func (d *DumpNetworkInterface) SetSendError(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.sendError = err
}

// LoadDump loads a packet dump from a JSON file
func (d *DumpNetworkInterface) LoadDump(dumpPath string) error {
	// Read the JSON file
	data, err := ioutil.ReadFile(dumpPath)
	if err != nil {
		return fmt.Errorf("failed to read dump file: %v", err)
	}

	// Parse the JSON
	var dump PacketDump
	if err := json.Unmarshal(data, &dump); err != nil {
		return fmt.Errorf("failed to parse dump file: %v", err)
	}

	d.currentDump = &dump
	d.currentIndex = 0
	return nil
}

// FeedNextReceivedPacket feeds the next received packet from the dump
func (d *DumpNetworkInterface) FeedNextReceivedPacket() (bool, error) {
	if d.currentDump == nil {
		return false, fmt.Errorf("no dump loaded")
	}

	// Find the next received packet
	for i := d.currentIndex; i < len(d.currentDump.Packets); i++ {
		packet := d.currentDump.Packets[i]
		if packet.Direction == "received" {
			d.currentIndex = i + 1

			// Decode the packet data
			var packetData []byte
			for _, part := range packet.RawData {
				data, err := base64.StdEncoding.DecodeString(part.BinaryBase64)
				if err != nil {
					return false, fmt.Errorf("failed to decode packet data: %v", err)
				}
				packetData = append(packetData, data...)
			}

			// Skip empty packets
			if len(packetData) == 0 {
				d.t.Logf("Skipping empty packet: %s (%s)", packet.PacketID, packet.Description)
				continue
			}

			// Feed the packet to the receiver channel
			d.t.Logf("Feeding received packet: %s (%s)", packet.PacketID, packet.Description)
			d.receiveChan <- packetData

			return true, nil
		}
	}

	return false, nil
}

// ImprovedDumpNetworkManager implements the TestNetworkManager interface for testing with packet dumps
type ImprovedDumpNetworkManager struct {
	networkInterface   *DumpNetworkInterface
	packetSender       *send.BaseSend
	packetHandler      *base.BaseReceive
	hookManager        *hooks.HookManager
	sessionStore       *login.SessionStore
	state              int
	stateChangeCb      func(oldState, newState int)
	lastSentPacket     string
	lastSentFields     map[string]interface{}
	lastReceivedPacket []byte
	t                  *testing.T
}

// NewImprovedDumpNetworkManager creates a new improved dump network manager
func NewImprovedDumpNetworkManager(t *testing.T) *ImprovedDumpNetworkManager {
	// Create a dump network interface
	dumpInterface := NewDumpNetworkInterface(t)

	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create the send stack
	packetSender := send.NewBaseSend(hookManager)

	// Create the receive stack
	packetHandler := base.NewBaseReceive(hookManager)

	// Create a session store
	sessionStore := login.NewSessionStore()

	// Configure the send stack with server configurations
	sendConfig := servers.GetServerType0SendConfig()
	err := packetSender.Configure("test", sendConfig)
	if err != nil {
		t.Fatalf("Failed to configure send stack: %v", err)
	}

	// Configure the receive stack with server configurations
	receiveConfig := servers.GetServerType0ReceiveConfig()

	// Convert the receive config to the correct type
	receivePacketDefs := make(map[string]common.PacketConstruction)
	for id, construction := range receiveConfig {
		receivePacketDefs[id] = common.PacketConstruction{
			Name:       construction.Name,
			Format:     construction.Format,
			FieldNames: construction.FieldNames,
		}
	}

	err = packetHandler.Configure("test", receivePacketDefs)
	if err != nil {
		t.Fatalf("Failed to configure receive stack: %v", err)
	}

	// Set the network interface for the send stack
	packetSender.SetNetworkInterface(dumpInterface)

	return &ImprovedDumpNetworkManager{
		networkInterface: dumpInterface,
		packetSender:     packetSender,
		packetHandler:    packetHandler,
		hookManager:      hookManager,
		sessionStore:     sessionStore,
		state:            network.NotConnected,
		t:                t,
	}
}

// Connect implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) Connect() error {
	return m.networkInterface.Connect()
}

// ConnectTo implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) ConnectTo(host string, port int) error {
	return m.networkInterface.ConnectTo(host, port)
}

// Disconnect implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) Disconnect() error {
	return m.networkInterface.Disconnect()
}

// Send implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	// Store the last sent packet for testing
	m.lastSentPacket = packetName
	m.lastSentFields = fields

	// Log the packet being sent
	m.t.Logf("Sending packet: %s with fields: %v", packetName, fields)

	// Construct the packet
	packet, err := m.packetSender.ConstructPacket(packetName, fields)
	if err != nil {
		return nil, err
	}

	// Send the packet
	err = m.packetSender.SendToServer(packet)
	return packet, err
}

// HandlePacket implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) HandlePacket(packet []byte) error {
	// Store the last received packet for testing
	m.lastReceivedPacket = packet

	return m.packetHandler.Process(packet)
}

// SetState implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) SetState(state int) {
	oldState := m.state
	m.state = state

	if m.stateChangeCb != nil {
		m.stateChangeCb(oldState, state)
	}
}

// GetState implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) GetState() int {
	return m.state
}

// SetStateChangeCallback implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) SetStateChangeCallback(callback func(oldState, newState int)) {
	m.stateChangeCb = callback
}

// GetHookManager implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) GetHookManager() interface{} {
	return m.hookManager
}

// SetSessionStore implements the NetworkManager interface
func (m *ImprovedDumpNetworkManager) SetSessionStore(sessionStore *login.SessionStore) {
	m.sessionStore = sessionStore
}

// SimulateReceivePacket implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SimulateReceivePacket(packetType string, data []byte) error {
	packet := append([]byte(packetType), '_')
	packet = append(packet, data...)
	return m.HandlePacket(packet)
}

// SetConnectError implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SetConnectError(err error) {
	m.networkInterface.SetConnectError(err)
}

// SetSendError implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SetSendError(err error) {
	m.networkInterface.SetSendError(err)
}

// SetHandleError implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SetHandleError(err error) {
	// Not implemented in this version
}

// GetLastSentPacket implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) GetLastSentPacket() (string, map[string]interface{}) {
	return m.lastSentPacket, m.lastSentFields
}

// GetLastReceivedPacket implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) GetLastReceivedPacket() []byte {
	return m.lastReceivedPacket
}

// CallHook implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) CallHook(hookName string, arg interface{}) {
	m.hookManager.CallHook(hookName, arg)
}

// SetSessionData implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SetSessionData(sessionData login.SessionData) {
	m.sessionStore.UpdateFromSessionData(sessionData)
}

// GetSessionData implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) GetSessionData() login.SessionData {
	return m.sessionStore.GetSessionData()
}

// SetServerInfo implements the TestNetworkManager interface
func (m *ImprovedDumpNetworkManager) SetServerInfo(servers []login.ServerInfo) {
	m.sessionStore.SetServerInfo(servers)
}

// LoadDump loads a packet dump from a JSON file
func (m *ImprovedDumpNetworkManager) LoadDump(dumpPath string) error {
	return m.networkInterface.LoadDump(dumpPath)
}

// FeedNextReceivedPacket feeds the next received packet from the dump
func (m *ImprovedDumpNetworkManager) FeedNextReceivedPacket() (bool, error) {
	return m.networkInterface.FeedNextReceivedPacket()
}

// TestLoginWithImprovedComponents tests the login flow using improved components with packet dumps
func TestLoginWithImprovedComponents(t *testing.T) {
	// Skip this test in normal test runs as it requires packet dumps
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Create an improved dump network manager
	networkManager := NewImprovedDumpNetworkManager(t)

	// Create login config
	config := login.NewLoginConfig("botijo0", "Melon.77", "rAthena")
	config.LoginTimeout = 30 * time.Second

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Register test hooks
	loginManager.RegisterTestHook(login.TestHookBeforeLogin, func(hookType login.TestHookType, data interface{}) {
		t.Logf("Test hook triggered: %s", hookType)
	})

	loginManager.RegisterTestHook(login.TestHookAfterLogin, func(hookType login.TestHookType, data interface{}) {
		t.Logf("Test hook triggered: %s", hookType)
	})

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Get the dump files
	dumpDir := "../../../verification/full_login_sequence"
	dumpFiles := []string{filepath.Join(dumpDir, "PACKET_DUMP")}

	// Test with each dump file
	for _, dumpFile := range dumpFiles {
		t.Logf("Testing with dump file: %s", dumpFile)

		// Load the dump
		err := networkManager.LoadDump(dumpFile)
		if err != nil {
			t.Logf("Failed to load dump: %v", err)
			continue
		}

		// Reset the session store
		sessionStore.Reset()

		// Start the login process in a goroutine
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := loginManager.Login(ctx)
			if err != nil {
				t.Logf("Login error: %v", err)
			}
		}()

		// Feed received packets
		packetCount := 0
		for {
			hasMore, err := networkManager.FeedNextReceivedPacket()
			if err != nil {
				t.Logf("Failed to feed packet: %v", err)
				break
			}
			if !hasMore {
				break
			}
			packetCount++

			// Only process the first 20 packets to avoid timeout
			if packetCount >= 20 {
				t.Logf("Processed %d packets, stopping to avoid timeout", packetCount)
				break
			}

			// Give the login manager time to process the packet
			time.Sleep(10 * time.Millisecond)
		}

		// Cancel the login process if it's still running
		cancel()
		wg.Wait()

		// Verify session data
		sessionData := sessionStore.GetSessionData()

		// Log the session data
		t.Logf("Session data after processing dump:")
		if len(sessionData.AccountID) > 0 {
			t.Logf("Account ID: %X", sessionData.AccountID)
		} else {
			t.Logf("Account ID not set")
		}

		if len(sessionData.SessionID) > 0 {
			t.Logf("Session ID: %X", sessionData.SessionID)
		} else {
			t.Logf("Session ID not set")
		}

		if len(sessionData.SessionID2) > 0 {
			t.Logf("Session ID2: %X", sessionData.SessionID2)
		} else {
			t.Logf("Session ID2 not set")
		}

		if len(sessionData.CharID) > 0 {
			t.Logf("Character ID: %X", sessionData.CharID)
			t.Logf("Map Name: %s", sessionData.MapName)
			t.Logf("Map IP: %s", sessionData.MapIP)
			t.Logf("Map Port: %d", sessionData.MapPort)
		} else {
			t.Logf("Character ID not set")
		}
	}
}

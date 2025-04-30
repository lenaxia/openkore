package integration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/lenaxia/goKore/network/login"
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

// DumpNetworkManager is a custom network manager that can be fed packets from the JSON dumps
type DumpNetworkManager struct {
	*login.MockNetworkManager
	t                *testing.T
	expectedPackets  map[string]bool // Map of expected packet IDs to be sent
	receivedPackets  map[string]bool // Map of received packet IDs
	currentDumpIndex int
	currentDump      *PacketDump
	sessionStore     *login.SessionStore
}

// NewDumpNetworkManager creates a new dump network manager
func NewDumpNetworkManager(t *testing.T) *DumpNetworkManager {
	mockManager := login.NewMockNetworkManager()
	sessionStore := login.NewSessionStore()

	// Set the session store in the mock manager
	mockManager.SetSessionStore(sessionStore)

	return &DumpNetworkManager{
		MockNetworkManager: mockManager,
		t:                  t,
		expectedPackets:    make(map[string]bool),
		receivedPackets:    make(map[string]bool),
		sessionStore:       sessionStore,
	}
}

// Send overrides the Send method to check if the packet is expected
func (m *DumpNetworkManager) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
	// Log the packet being sent
	m.t.Logf("Sending packet: %s with fields: %v", packetName, fields)

	// Check if this packet is expected
	if _, ok := m.expectedPackets[packetName]; ok {
		m.receivedPackets[packetName] = true
		m.t.Logf("Expected packet %s was sent", packetName)
	} else {
		m.t.Logf("Unexpected packet %s was sent", packetName)
	}

	// Call the original Send method
	return m.MockNetworkManager.Send(packetName, fields)
}

// LoadDump loads a packet dump from a JSON file
func (m *DumpNetworkManager) LoadDump(dumpPath string) error {
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

	m.currentDump = &dump
	m.currentDumpIndex = 0

	// Extract expected packets
	for _, packet := range dump.Packets {
		if packet.Direction == "sent" {
			m.expectedPackets[packet.PacketID] = true
		}
	}

	return nil
}

// FeedNextReceivedPacket feeds the next received packet from the dump
func (m *DumpNetworkManager) FeedNextReceivedPacket() (bool, error) {
	if m.currentDump == nil {
		return false, fmt.Errorf("no dump loaded")
	}

	// Find the next received packet
	for i := m.currentDumpIndex; i < len(m.currentDump.Packets); i++ {
		packet := m.currentDump.Packets[i]
		if packet.Direction == "received" {
			m.currentDumpIndex = i + 1

			// Decode the packet data
			var packetData []byte
			for _, part := range packet.RawData {
				data, err := base64.StdEncoding.DecodeString(part.BinaryBase64)
				if err != nil {
					return false, fmt.Errorf("failed to decode packet data: %v", err)
				}
				packetData = append(packetData, data...)
			}

			// Feed the packet to the handler
			m.t.Logf("Feeding received packet: %s (%s)", packet.PacketID, packet.Description)

			// Skip empty packets
			if len(packetData) == 0 {
				m.t.Logf("Skipping empty packet")
				continue
			}

			err := m.HandlePacket(packetData)
			if err != nil {
				m.t.Logf("Error handling packet: %v", err)
			}

			return true, nil
		}
	}

	return false, nil
}

// VerifySessionData verifies that the session data was updated correctly
func (m *DumpNetworkManager) VerifySessionData(t *testing.T) {
	if m.sessionStore == nil {
		t.Error("Session store is nil")
		return
	}

	sessionData := m.sessionStore.GetSessionData()

	// Verify that we have account ID and session IDs
	if len(sessionData.AccountID) == 0 {
		t.Error("Account ID not set")
	} else {
		t.Logf("Account ID: %X", sessionData.AccountID)
	}

	if len(sessionData.SessionID) == 0 {
		t.Error("Session ID not set")
	} else {
		t.Logf("Session ID: %X", sessionData.SessionID)
	}

	if len(sessionData.SessionID2) == 0 {
		t.Error("Session ID2 not set")
	} else {
		t.Logf("Session ID2: %X", sessionData.SessionID2)
	}

	// If we got to the map server, verify map info
	if len(sessionData.CharID) > 0 {
		t.Logf("Character ID: %X", sessionData.CharID)
		t.Logf("Map Name: %s", sessionData.MapName)
		t.Logf("Map IP: %s", sessionData.MapIP)
		t.Logf("Map Port: %d", sessionData.MapPort)
	}
}

// TestLoginWithDumps tests the login flow using packet dumps
func TestLoginWithDumps(t *testing.T) {
	// Create a dump network manager
	networkManager := NewDumpNetworkManager(t)

	// Create login config
	config := login.NewLoginConfig("botijo0", "Melon.77", "rAthena")

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Get the dump files
	dumpDir := "../../../verification/PacketAnalysis/extracteddata"
	dumpFiles, err := filepath.Glob(filepath.Join(dumpDir, "dump*_packets.json"))
	if err != nil {
		t.Fatalf("Failed to find dump files: %v", err)
	}

	if len(dumpFiles) == 0 {
		t.Fatalf("No dump files found in %s", dumpDir)
	}

	// Test with each dump file
	for _, dumpFile := range dumpFiles {
		t.Logf("Testing with dump file: %s", dumpFile)

		// Load the dump
		err := networkManager.LoadDump(dumpFile)
		if err != nil {
			t.Fatalf("Failed to load dump: %v", err)
		}

		// Reset the session store
		sessionStore.Reset()

		// Simulate connecting to the master server
		err = networkManager.Connect()
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}

		// Feed received packets and observe the client's response
		packetCount := 0
		for {
			hasMore, err := networkManager.FeedNextReceivedPacket()
			if err != nil {
				t.Fatalf("Failed to feed packet: %v", err)
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
		}

		// Verify that the session data was updated correctly
		networkManager.VerifySessionData(t)

		// Disconnect
		networkManager.Disconnect()
	}
}

// TestLoginSequence tests the login sequence with specific packets
func TestLoginSequence(t *testing.T) {
	// Create a dump network manager
	networkManager := NewDumpNetworkManager(t)

	// Create login config
	config := login.NewLoginConfig("botijo0", "Melon.77", "rAthena")

	// Create login manager
	loginManager := login.NewLoginManager(networkManager, config)

	// Get the session store
	sessionStore := loginManager.GetSessionStore()

	// Simulate connecting to the master server
	err := networkManager.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Simulate receiving account server info packet (0AC4)
	accountInfoPacket := map[string]interface{}{
		"sessionID":  []byte{0xE5, 0x5D, 0xF6, 0xC1}, // From packet dump
		"accountID":  []byte{0x82, 0x84, 0x1E, 0x00}, // From packet dump
		"sessionID2": []byte{0x01, 0x2C, 0x9C, 0x53}, // From packet dump
		"accountSex": 0,                              // 0 = male
	}
	sessionStore.UpdateFromAccountServerInfo(accountInfoPacket)

	// Set server info
	servers := []login.ServerInfo{
		{
			Name: "rAthena",
			IP:   "192.168.5.219",
			Port: 6121,
		},
	}
	sessionStore.SetServerInfo(servers)

	// Simulate connecting to the character server
	err = networkManager.ConnectTo("192.168.5.219", 6121)
	if err != nil {
		t.Fatalf("Failed to connect to character server: %v", err)
	}

	// Simulate receiving character list packet (082D and 006B)
	// This would normally be handled by the packet handler

	// Simulate receiving character map info packet (0AC5)
	charMapInfoPacket := map[string]interface{}{
		"charID":  []byte{0xF2, 0x49, 0x02, 0x00}, // From packet dump
		"mapName": "gef_fild07.gat",
		"mapIP":   "192.168.5.219",
		"mapPort": 5121,
	}
	sessionStore.UpdateFromCharacterServerInfo(charMapInfoPacket)

	// Simulate connecting to the map server
	err = networkManager.ConnectTo("192.168.5.219", 5121)
	if err != nil {
		t.Fatalf("Failed to connect to map server: %v", err)
	}

	// Verify session data
	sessionData := sessionStore.GetSessionData()
	if len(sessionData.AccountID) != 4 || sessionData.AccountID[0] != 0x82 || sessionData.AccountID[1] != 0x84 ||
		sessionData.AccountID[2] != 0x1E || sessionData.AccountID[3] != 0x00 {
		t.Errorf("Account ID not stored correctly: %v", sessionData.AccountID)
	}

	if len(sessionData.CharID) != 4 || sessionData.CharID[0] != 0xF2 || sessionData.CharID[1] != 0x49 ||
		sessionData.CharID[2] != 0x02 || sessionData.CharID[3] != 0x00 {
		t.Errorf("Character ID not stored correctly: %v", sessionData.CharID)
	}

	if sessionData.MapName != "gef_fild07.gat" {
		t.Errorf("Map name not stored correctly: %v", sessionData.MapName)
	}

	if sessionData.MapIP != "192.168.5.219" {
		t.Errorf("Map IP not stored correctly: %v", sessionData.MapIP)
	}

	if sessionData.MapPort != 5121 {
		t.Errorf("Map port not stored correctly: %v", sessionData.MapPort)
	}

	// Disconnect
	networkManager.Disconnect()
}

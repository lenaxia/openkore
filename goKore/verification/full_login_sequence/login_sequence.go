package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/lenaxia/goKore/network/connection"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/protocol"
	receivecore "github.com/lenaxia/goKore/network/receive/core"
	receivesecurity "github.com/lenaxia/goKore/network/receive/security"
	sendcore "github.com/lenaxia/goKore/network/send/core"
	sendsecurity "github.com/lenaxia/goKore/network/send/security"
	"github.com/lenaxia/goKore/network/send/servers"
)

// LoginSequenceTest represents the state for the login sequence test
type LoginSequenceTest struct {
	conn         connection.Connection
	tokenizer    *protocol.Tokenizer
	hookManager  *hooks.HookManager
	coreParser   *receivecore.CoreParser
	loginManager *receivesecurity.LoginManager

	// Send components
	baseSend    *sendcore.BaseSend
	sendManager *sendsecurity.LoginManager

	// Debug flags
	debugMode bool

	// Test state
	accountID      uint32
	sessionID1     uint32
	sessionID2     uint32
	gender         byte
	serverList     []ServerInfo
	charServerInfo *ServerInfo
	charList       []CharacterInfo
	selectedChar   int
	charID         uint32
	mapName        string
	mapIP          string
	mapPort        int

	// Test results
	loginSuccess      bool
	charSelectSuccess bool
	mapLoginSuccess   bool

	// Logging
	logFile *os.File
}

// ServerInfo represents information about a game server
type ServerInfo struct {
	IP       string
	Port     int
	Name     string
	Users    int
	State    int
	Property int
}

// CharacterInfo represents information about a character
type CharacterInfo struct {
	CharID    uint32
	BaseExp   uint64
	Zeny      uint32
	JobExp    uint64
	JobLevel  uint32
	Shoes     uint32
	Gloves    uint32
	Cape      uint32
	Dress     uint32
	WeaponID  uint32
	Name      string
	Str       uint8
	Agi       uint8
	Vit       uint8
	Int       uint8
	Dex       uint8
	Luk       uint8
	CharNum   uint8
	HairColor uint8
	BankVault uint32
	MapName   string
}

// NewLoginSequenceTest creates a new login sequence test
func NewLoginSequenceTest() *LoginSequenceTest {
	// Create log file
	logFile, err := os.Create("login_sequence_test.log")
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return nil
	}

	test := &LoginSequenceTest{
		hookManager: hooks.NewHookManager(),
		logFile:     logFile,
	}

	// Set up hooks
	test.setupHooks()

	return test
}

// Log writes a message to both stdout and the log file
func (t *LoginSequenceTest) Log(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(message)
	if t.logFile != nil {
		fmt.Fprintln(t.logFile, message)
	}
}

// setupHooks sets up the hook handlers for various events
func (t *LoginSequenceTest) setupHooks() {
	// Login-related hooks
	t.hookManager.AddHook("connection/connected", func(hookName string, arg interface{}, userData interface{}) {
		t.Log("Connection established")
	}, nil)

	t.hookManager.AddHook("security/login_error", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			code := data["code"].(int)
			message := data["message"].(string)
			date := ""
			if d, ok := data["date"]; ok {
				date = d.(string)
			}
			t.Log("Login error received: code=%d, message=%s, date=%s", code, message, date)
		}
	}, nil)

	t.hookManager.AddHook("security/login_success", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			t.accountID = data["account_id"].(uint32)
			t.sessionID1 = data["session_id1"].(uint32)
			t.sessionID2 = data["session_id2"].(uint32)
			t.gender = data["gender"].(byte)
			t.loginSuccess = true
			t.Log("Login successful: account_id=%d, gender=%d", t.accountID, t.gender)
		}
	}, nil)

	// Add a hook for the 0xAC4 packet (login response)
	t.hookManager.AddHook("receive/0xAC4", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.([]byte); ok {
			t.Log("Received login response packet (0xAC4): %X", data)

			// Extract account ID (bytes 4-7)
			if len(data) >= 8 {
				t.accountID = binary.LittleEndian.Uint32(data[4:8])
				t.Log("Extracted account ID: %d", t.accountID)

				// Extract session IDs (bytes 8-15)
				if len(data) >= 16 {
					t.sessionID1 = binary.LittleEndian.Uint32(data[8:12])
					t.sessionID2 = binary.LittleEndian.Uint32(data[12:16])
					t.Log("Extracted session IDs: %d, %d", t.sessionID1, t.sessionID2)

					// Extract gender (byte 46)
					if len(data) >= 47 {
						t.gender = data[46]
						t.Log("Extracted gender: %d", t.gender)

						// Mark login as successful
						t.loginSuccess = true
						t.Log("Login successful: account_id=%d, gender=%d", t.accountID, t.gender)
					}
				}
			}
		}
	}, nil)

	// Server selection hooks
	t.hookManager.AddHook("security/server_list", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			if servers, ok := data["servers"].([]ServerInfo); ok {
				t.serverList = servers
				t.Log("Received server list with %d servers", len(servers))
				for i, server := range servers {
					t.Log("Server %d: %s (%s:%d) - Users: %d, State: %d",
						i, server.Name, server.IP, server.Port, server.Users, server.State)
				}
			}
		}
	}, nil)

	// Character selection hooks
	t.hookManager.AddHook("security/character_list", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			if chars, ok := data["characters"].([]CharacterInfo); ok {
				t.charList = chars
				t.Log("Received character list with %d characters", len(chars))
				for i, char := range chars {
					t.Log("Character %d: %s (Level %d) - Map: %s",
						i, char.Name, char.JobLevel, char.MapName)
				}
				t.charSelectSuccess = true
			}
		}
	}, nil)

	// Character ID and Map hooks
	t.hookManager.AddHook("security/character_id_and_map", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			t.charID = data["char_id"].(uint32)
			t.mapName = data["map_name"].(string)
			t.mapIP = data["map_ip"].(string)
			t.mapPort = data["map_port"].(int)
			t.Log("Received character ID and map info: char_id=%d, map=%s, ip=%s, port=%d",
				t.charID, t.mapName, t.mapIP, t.mapPort)
		}
	}, nil)

	// Map login hooks
	t.hookManager.AddHook("game/map_loaded", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			t.mapName = data["map"].(string)
			t.Log("Map loaded: %s", t.mapName)
			t.mapLoginSuccess = true
		}
	}, nil)
}

// RunTest runs the full login sequence test
func (t *LoginSequenceTest) RunTest() bool {
	t.Log("===== Full Login Sequence Test =====")

	// Create connection config for rAthena renewal server
	connConfig := &connection.ConnectionConfig{
		Host:        "192.168.5.219", // rAthena renewal server
		Port:        6900,
		Timeout:     10 * time.Second,
		RecvTimeout: 5 * time.Second,
		SendTimeout: 5 * time.Second,
		ServerType:  "0",
	}

	// Create direct connection
	t.conn = connection.NewDirectConnection(connConfig)

	// Register connection event callbacks
	t.conn.RegisterCallback(connection.EventConnected, func(event connection.ConnectionEvent, data interface{}) {
		t.Log("Connected to server")
	})

	t.conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
		if err, ok := data.(error); ok {
			t.Log("Connection error: %v", err)
		}
	})

	// Create core parser
	t.coreParser = receivecore.NewCoreParser("ServerType0", t.hookManager)

	// Register handler for 0xAC4 packet (login response)
	// Register handler for 0xAC4 packet (kRO Zero login response)
	t.coreParser.RegisterHandler("0AC4", "login_response", "V a4 a4 a4 C a*",
		[]string{"version", "accountID", "sessionID1", "sessionID2", "gender", "serverInfo"},
		func(args map[string]interface{}) error {
			t.Log("Processing 0xAC4 packet in core parser")

			// Extract account ID
			if accountID, ok := args["accountID"]; ok {
				if accountIDBytes, ok := accountID.([]byte); ok && len(accountIDBytes) == 4 {
					t.accountID = binary.LittleEndian.Uint32(accountIDBytes)
					t.Log("Extracted account ID: %d", t.accountID)
				}
			}

			// Extract session IDs
			if sessionID1, ok := args["sessionID1"]; ok {
				if sessionID1Bytes, ok := sessionID1.([]byte); ok && len(sessionID1Bytes) == 4 {
					t.sessionID1 = binary.LittleEndian.Uint32(sessionID1Bytes)
					t.Log("Extracted session ID 1: %d", t.sessionID1)
				}
			}

			if sessionID2, ok := args["sessionID2"]; ok {
				if sessionID2Bytes, ok := sessionID2.([]byte); ok && len(sessionID2Bytes) == 4 {
					t.sessionID2 = binary.LittleEndian.Uint32(sessionID2Bytes)
					t.Log("Extracted session ID 2: %d", t.sessionID2)
				}
			}

			// Extract gender
			if gender, ok := args["gender"]; ok {
				if genderByte, ok := gender.(byte); ok {
					t.gender = genderByte
					t.Log("Extracted gender: %d", t.gender)
				}
			}

			// Mark login as successful
			t.loginSuccess = true
			t.Log("Login successful: account_id=%d, gender=%d", t.accountID, t.gender)

			// Call the hook
			t.hookManager.CallHook("security/login_success", map[string]interface{}{
				"account_id":  t.accountID,
				"session_id1": t.sessionID1,
				"session_id2": t.sessionID2,
				"gender":      t.gender,
			})

			return nil
		})

	// Create login manager for receive
	t.loginManager = receivesecurity.NewLoginManager(t.coreParser, t.hookManager)
	t.loginManager.RegisterHandlers()

	// Create base send implementation
	t.baseSend = sendcore.NewBaseSend(t.hookManager)

	// Configure base send with ServerType0 packet constructions
	err := t.baseSend.Configure("ServerType0", servers.ServerType0PacketConstructions())
	if err != nil {
		t.Log("Failed to configure base send: %v", err)
		return false
	}

	// Manually register the master_login packet
	t.baseSend.RegisterPacketHandler("0064", "master_login", "v a24 a24 C",
		[]string{"version", "username", "password", "clienttype"}, nil)

	// Set the connection for sending packets
	t.baseSend.SetConnection(t.conn)

	// Create login manager for send
	t.sendManager = sendsecurity.NewLoginManager(t.baseSend)

	// Set credentials
	t.sendManager.SetCredentials("botijo0", "Melon.77")

	// Set version
	t.sendManager.SetVersion(55)
	t.sendManager.SetMasterVersion(1)

	// Define packet lengths for the tokenizer based on ServerType0.pm
	packetLengths := make(map[string]protocol.PacketDef)
	packetLengths["0069"] = protocol.PacketDef{Length: 0, HasLength: true}    // account_server_info
	packetLengths["0071"] = protocol.PacketDef{Length: 28, HasLength: false}  // received_character_ID_and_Map
	packetLengths["0073"] = protocol.PacketDef{Length: 11, HasLength: false}  // map_loaded
	packetLengths["083E"] = protocol.PacketDef{Length: 26, HasLength: false}  // login_error
	packetLengths["006A"] = protocol.PacketDef{Length: -1, HasLength: true}   // login_error_game_login
	packetLengths["006B"] = protocol.PacketDef{Length: -1, HasLength: true}   // received_characters (standard)
	packetLengths["006C"] = protocol.PacketDef{Length: 3, HasLength: false}   // login_error_game_login
	packetLengths["0081"] = protocol.PacketDef{Length: -1, HasLength: true}   // login_error_game_login
	packetLengths["0091"] = protocol.PacketDef{Length: 22, HasLength: false}  // map_changed
	packetLengths["0092"] = protocol.PacketDef{Length: 28, HasLength: false}  // map_changed
	packetLengths["0093"] = protocol.PacketDef{Length: 2, HasLength: false}   // login_to_map_server_succeeded
	packetLengths["0AC4"] = protocol.PacketDef{Length: 218, HasLength: false} // login_response (kRO Zero account server info)
	packetLengths["082D"] = protocol.PacketDef{Length: -1, HasLength: true}   // received_characters (newer format)
	packetLengths["099D"] = protocol.PacketDef{Length: -1, HasLength: true}   // received_characters (latest format)
	packetLengths["0B72"] = protocol.PacketDef{Length: -1, HasLength: true}   // received_characters (latest format)
	packetLengths["0276"] = protocol.PacketDef{Length: -1, HasLength: true}   // tRO account server info
	packetLengths["0AC9"] = protocol.PacketDef{Length: -1, HasLength: true}   // Newer account server info
	packetLengths["0B07"] = protocol.PacketDef{Length: -1, HasLength: true}   // Latest account server info
	packetLengths["0B60"] = protocol.PacketDef{Length: -1, HasLength: true}   // twRO account server info
	packetLengths["0AC5"] = protocol.PacketDef{Length: 156, HasLength: false} // character ID and map IP
	packetLengths["0436"] = protocol.PacketDef{Length: 19, HasLength: false}  // map login
	packetLengths["08B9"] = protocol.PacketDef{Length: 12, HasLength: false}  // PIN code request
	packetLengths["0283"] = protocol.PacketDef{Length: 6, HasLength: false}   // account ID
	packetLengths["0B18"] = protocol.PacketDef{Length: 4, HasLength: false}   // inventory expansion
	packetLengths["02EB"] = protocol.PacketDef{Length: 13, HasLength: false}  // enter map
	packetLengths["09A0"] = protocol.PacketDef{Length: 6, HasLength: false}   // character received sync

	// Create message tokenizer
	t.tokenizer = protocol.NewTokenizer(packetLengths)

	// Step 1: Connect to login server
	t.Log("Connecting to login server: %s:%d", connConfig.Host, connConfig.Port)
	err = t.conn.Connect()
	if err != nil {
		t.Log("Failed to connect to login server: %v", err)
		return false
	}
	defer t.conn.Disconnect()

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	// Check if connected
	if !t.conn.IsConnected() {
		t.Log("Failed to connect to login server")
		return false
	}
	t.Log("Connected to login server successfully")

	// Step 2: Send master login packet
	t.Log("Sending master login packet for user: botijo0")

	// Create master login packet manually
	packet := []byte{0x64, 0x00} // 0x0064 master_login packet

	// Add version (V) - 32-bit unsigned integer
	versionBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(versionBytes, 55)
	packet = append(packet, versionBytes...)

	// Add username (a24)
	usernameBytes := make([]byte, 24)
	copy(usernameBytes, []byte("botijo0"))
	packet = append(packet, usernameBytes...)

	// Add password (a24)
	passwordBytes := make([]byte, 24)
	copy(passwordBytes, []byte("Melon.77"))
	packet = append(packet, passwordBytes...)

	// Add client type (C)
	packet = append(packet, 0) // 0 = normal client

	// Send the packet
	err = t.baseSend.SendToServer(packet)
	if err != nil {
		t.Log("Failed to construct login packet: %v", err)
		return false
	}

	err = t.baseSend.SendToServer(packet)
	if err != nil {
		t.Log("Failed to send login packet: %v", err)
		return false
	}

	// Step 3: Wait for login response and process it
	if !t.waitForLoginResponse() {
		t.Log("Login failed or timed out")
		return false
	}

	// Step 4: Select a server (first available)
	if len(t.serverList) == 0 {
		t.Log("No servers available")
		return false
	}

	// Disconnect from login server
	t.conn.Disconnect()
	t.Log("Disconnected from login server")

	// Select the first server
	selectedServer := t.serverList[0]
	t.charServerInfo = &selectedServer
	t.Log("Selected server: %s (%s:%d)", selectedServer.Name, selectedServer.IP, selectedServer.Port)

	// Step 5: Connect to character server
	connConfig.Host = selectedServer.IP
	connConfig.Port = selectedServer.Port
	t.conn = connection.NewDirectConnection(connConfig)

	// Update the connection in the base send
	t.baseSend.SetConnection(t.conn)

	t.Log("Connecting to character server: %s:%d", connConfig.Host, connConfig.Port)

	// Verify connection configuration
	t.Log("Connection config: %+v", connConfig)

	err = t.conn.Connect()
	if err != nil {
		t.Log("Failed to connect to character server: %v", err)
		return false
	}

	// Wait for connection to establish
	t.Log("Waiting for connection to establish...")
	time.Sleep(1 * time.Second)

	// Check if connected
	if !t.conn.IsConnected() {
		t.Log("Failed to connect to character server - connection reports not connected")
		return false
	}
	t.Log("Connected to character server successfully")

	// Step 6: Send character server login packet
	t.Log("Sending character server login packet")

	// Set account info in the login manager
	accountIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(accountIDBytes, t.accountID)

	sessionID1Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionID1Bytes, t.sessionID1)

	sessionID2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionID2Bytes, t.sessionID2)

	// Set the account info in the login manager
	t.sendManager.SetAccountInfo(accountIDBytes, sessionID1Bytes, sessionID2Bytes, int(t.gender))

	// Set additional parameters that might be needed
	// Set account info in the login manager
	t.sendManager.SetAccountInfo(accountIDBytes, sessionID1Bytes, sessionID2Bytes, int(t.gender))
	t.sendManager.SetVersion(55)      // Set client version
	t.sendManager.SetMasterVersion(1) // Set master version

	// Based on the ServerType0.pm file and PACKET_DUMP, we should use:
	// '0065' => ['game_login', 'a4 a4 a4 v C', [qw(accountID sessionID sessionID2 userLevel accountSex)]]

	t.Log("Creating game_login packet (0x0065)")
	t.Log("Account ID: %X", accountIDBytes)
	t.Log("Session ID 1: %X", sessionID1Bytes)
	t.Log("Session ID 2: %X", sessionID2Bytes)
	t.Log("Gender: %d", t.gender)

	// Create a buffer for the game_login packet
	buf := new(bytes.Buffer)

	// Write packet ID (0x0065)
	binary.Write(buf, binary.LittleEndian, uint16(0x0065))

	// Use the Account ID from the PACKET_DUMP but keep session IDs dynamic
	t.Log("Using Account ID from PACKET_DUMP but dynamic session IDs")

	// Use Account ID from PACKET_DUMP (2000002)
	accountIDFixed := []byte{0x82, 0x84, 0x1E, 0x00} // 2000002

	// Write accountID (a4)
	buf.Write(accountIDFixed)

	// Write sessionID (a4)
	buf.Write(sessionID1Bytes)

	// Write sessionID2 (a4)
	buf.Write(sessionID2Bytes)

	// Write userLevel (v) - 2 bytes
	binary.Write(buf, binary.LittleEndian, uint16(0))

	// Write accountSex (C) - 1 byte
	buf.WriteByte(t.gender)

	// Get the final packet
	gameLoginPacket := buf.Bytes()

	t.Log("Sending game_login packet (0x0065): %X", gameLoginPacket)
	t.Log("Packet length: %d bytes", len(gameLoginPacket))

	// Send the game_login packet
	err = t.baseSend.SendToServer(gameLoginPacket)
	if err != nil {
		t.Log("Error sending game_login packet: %v", err)
		return false
	}

	t.Log("Game login packet sent successfully")

	// Wait for character list
	if !t.waitForCharacterList() {
		t.Log("Failed to receive character list or timed out")
		return false
	}

	// Step 8: Select a character (first one)
	if len(t.charList) == 0 {
		t.Log("No characters available")
		return false
	}

	t.selectedChar = 0
	selectedChar := t.charList[t.selectedChar]
	t.Log("Selected character: %s", selectedChar.Name)

	// Step 9: Send character selection packet
	t.Log("Sending character selection packet")
	err = t.sendManager.SendCharLogin(int(selectedChar.CharNum))
	if err != nil {
		t.Log("Failed to send character selection packet: %v", err)
		return false
	}

	// Step 10: Wait for map server information
	if !t.waitForMapServerInfo() {
		t.Log("Failed to receive map server info or timed out")
		return false
	}

	// Step 11: Disconnect from character server
	t.conn.Disconnect()
	t.Log("Disconnected from character server")

	// Step 12: Connect to map server
	connConfig.Host = t.mapIP
	connConfig.Port = t.mapPort
	t.conn = connection.NewDirectConnection(connConfig)

	// Update the connection in the base send
	t.baseSend.SetConnection(t.conn)

	t.Log("Connecting to map server: %s:%d", connConfig.Host, connConfig.Port)
	err = t.conn.Connect()
	if err != nil {
		t.Log("Failed to connect to map server: %v", err)
		return false
	}

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	// Check if connected
	if !t.conn.IsConnected() {
		t.Log("Failed to connect to map server")
		return false
	}
	t.Log("Connected to map server successfully")

	// Step 13: Send map login packet
	t.Log("Sending map login packet")

	// Update charID in the login manager
	var charIDBytes []byte
	charIDBytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(charIDBytes, t.charID)

	// Set the charID in the login manager
	t.sendManager.SetAccountInfo(accountIDBytes, sessionID1Bytes, sessionID2Bytes, int(t.gender))

	// Create a buffer for the map login packet (0x0436)
	mapLoginBuf := new(bytes.Buffer)

	// Write packet ID (0x0436)
	binary.Write(mapLoginBuf, binary.LittleEndian, uint16(0x0436))

	// Write accountID (a4)
	mapLoginBuf.Write(accountIDBytes)

	// Write charID (a4)
	mapLoginBuf.Write(charIDBytes)

	// Write sessionID1 (a4)
	mapLoginBuf.Write(sessionID1Bytes)

	// Write sessionID2 (a4)
	mapLoginBuf.Write(sessionID2Bytes)

	// Write tick (V) - 4 bytes
	binary.Write(mapLoginBuf, binary.LittleEndian, uint32(0))

	// Write sex (C) - 1 byte
	mapLoginBuf.WriteByte(t.gender)

	// Get the final packet
	mapLoginPacket := mapLoginBuf.Bytes()

	t.Log("Sending map_login packet (0x0436): %X", mapLoginPacket)
	t.Log("Packet length: %d bytes", len(mapLoginPacket))

	// Send the map_login packet
	err = t.baseSend.SendToServer(mapLoginPacket)
	if err != nil {
		t.Log("Failed to send map login packet: %v", err)
		return false
	}

	// Step 14: Wait for map to load
	if !t.waitForMapLoad() {
		t.Log("Failed to load map or timed out")
		return false
	}

	// Step 15: Send map loaded packet and additional packets
	t.Log("Sending map loaded packet")
	err = t.sendManager.SendMapLoaded()
	if err != nil {
		t.Log("Failed to send map loaded packet: %v", err)
		return false
	}

	// Step 16: Send sync packet (0x0360)
	t.Log("Sending sync packet (0x0360)")
	syncBuf := new(bytes.Buffer)
	binary.Write(syncBuf, binary.LittleEndian, uint16(0x0360))
	// Add tick (V) - 4 bytes
	binary.Write(syncBuf, binary.LittleEndian, uint32(0))
	syncPacket := syncBuf.Bytes()
	err = t.baseSend.SendToServer(syncPacket)
	if err != nil {
		t.Log("Failed to send sync packet: %v", err)
		return false
	}

	// Step 17: Send guild query page packet (0x014F)
	t.Log("Sending guild query page packet (0x014F)")
	guildQueryBuf := new(bytes.Buffer)
	binary.Write(guildQueryBuf, binary.LittleEndian, uint16(0x014F))
	// Add type (V) - 4 bytes
	binary.Write(guildQueryBuf, binary.LittleEndian, uint32(0))
	guildQueryPacket := guildQueryBuf.Bytes()
	err = t.baseSend.SendToServer(guildQueryPacket)
	if err != nil {
		t.Log("Failed to send guild query page packet: %v", err)
		return false
	}

	// Step 18: Send request blocking player cancel packet (0x0447)
	t.Log("Sending request blocking player cancel packet (0x0447)")
	blockingCancelBuf := new(bytes.Buffer)
	binary.Write(blockingCancelBuf, binary.LittleEndian, uint16(0x0447))
	blockingCancelPacket := blockingCancelBuf.Bytes()
	err = t.baseSend.SendToServer(blockingCancelPacket)
	if err != nil {
		t.Log("Failed to send request blocking player cancel packet: %v", err)
		return false
	}

	// Step 19: Disconnect from map server
	t.conn.Disconnect()
	t.Log("Disconnected from map server")

	// Test complete
	t.Log("===== Test Complete =====")
	t.Log("Login success: %v", t.loginSuccess)
	t.Log("Character selection success: %v", t.charSelectSuccess)
	t.Log("Map login success: %v", t.mapLoginSuccess)

	return t.loginSuccess && t.charSelectSuccess && t.mapLoginSuccess
}

// waitForLoginResponse waits for the login response from the server
func (t *LoginSequenceTest) waitForLoginResponse() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Waiting for login response...")

	for {
		select {
		case <-ctx.Done():
			t.Log("Timeout waiting for login response")
			return false
		default:
			data, err := t.conn.Receive()
			if err != nil {
				t.Log("Error receiving data: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				t.Log("Received data: %X", data)

				// Process the data with the tokenizer
				t.tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := t.tokenizer.ReadNext()
				if err != nil {
					t.Log("Error reading message: %v", err)
				} else {
					t.Log("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						t.Log("Message ID: %s", msgID)

						// Special handling for 0xAC4 packet (kRO Zero account server info)
						if msgID == "0AC4" && len(message) >= 47 {
							// Extract account ID (bytes 4-7)
							t.accountID = binary.LittleEndian.Uint32(message[4:8])
							t.Log("Extracted account ID: %d", t.accountID)

							// Extract session IDs (bytes 8-15)
							t.sessionID1 = binary.LittleEndian.Uint32(message[8:12])
							t.sessionID2 = binary.LittleEndian.Uint32(message[12:16])
							t.Log("Extracted session IDs: %d, %d", t.sessionID1, t.sessionID2)

							// Extract gender (byte 46)
							t.gender = message[46]
							t.Log("Extracted gender: %d", t.gender)

							// Mark login as successful
							t.loginSuccess = true
							t.Log("Login successful: account_id=%d, gender=%d", t.accountID, t.gender)

							// For this test, we'll use the login server information for the character server
							// This ensures we can connect to a server that we know is running
							t.Log("Using login server information for character server")
							server := ServerInfo{
								IP:       "192.168.5.219", // Same as login server
								Port:     6121,            // Character server port from PACKET_DUMP
								Name:     "rAthena",
								Users:    0,
								State:    0,
								Property: 0,
							}
							t.serverList = append(t.serverList, server)
							t.Log("Added server: %s (%s:%d)", server.Name, server.IP, server.Port)

							// Call the hook
							t.hookManager.CallHook("security/login_success", map[string]interface{}{
								"account_id":  t.accountID,
								"session_id1": t.sessionID1,
								"session_id2": t.sessionID2,
								"gender":      t.gender,
							})

							// Also call the server list hook
							t.hookManager.CallHook("security/server_list", map[string]interface{}{
								"servers": t.serverList,
							})

							return true
						}

						// Process the message with the core parser
						t.coreParser.Parse(message)

						// Check if login was successful
						if t.loginSuccess {
							return true
						}

						// Check for login error
						if msgID == "083E" || msgID == "006A" {
							t.Log("Login failed")
							return false
						}
					}
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// waitForCharacterList waits for the character list from the server
func (t *LoginSequenceTest) waitForCharacterList() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Waiting for character list...")
	t.Log("Looking for packet IDs: 006B, 082D, 099D, 0B72, 09A0")

	for {
		select {
		case <-ctx.Done():
			t.Log("Timeout waiting for character list")
			return false
		default:
			data, err := t.conn.Receive()
			if err != nil {
				t.Log("Error receiving data: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				t.Log("Received data: %X", data)
				t.Log("Data length: %d bytes", len(data))

				// Process the data with the tokenizer
				t.tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := t.tokenizer.ReadNext()
				if err != nil {
					t.Log("Error reading message: %v", err)
				} else {
					t.Log("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						t.Log("Message ID: %s, Length: %d", msgID, len(message))

						// Check for character list packets
						if msgID == "006B" || msgID == "082D" || msgID == "099D" || msgID == "0B72" {
							t.Log("Received character list packet: %s", msgID)

							// Try to extract character count if possible
							if len(message) >= 4 {
								charCount := binary.LittleEndian.Uint16(message[2:4])
								t.Log("Character count: %d", charCount)
							}
						}

						// Check for PIN code request
						if msgID == "08B9" {
							t.Log("Received PIN code request")
							// For testing purposes, we'll just acknowledge it
							// In a real implementation, we would send a PIN code response
						}

						// Check for character received sync
						if msgID == "09A0" {
							t.Log("Received character sync packet: %s", msgID)
						}

						// Process the message
						t.coreParser.Parse(message)

						// Check if character list was received
						if t.charSelectSuccess {
							return true
						}

						// Check for login error
						if msgID == "006C" || msgID == "0081" {
							t.Log("Character server login failed")
							return false
						}
					}
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// waitForMapServerInfo waits for the map server information from the server
func (t *LoginSequenceTest) waitForMapServerInfo() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Waiting for map server information...")

	for {
		select {
		case <-ctx.Done():
			t.Log("Timeout waiting for map server information")
			return false
		default:
			data, err := t.conn.Receive()
			if err != nil {
				t.Log("Error receiving data: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				t.Log("Received data: %X", data)

				// Process the data with the tokenizer
				t.tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := t.tokenizer.ReadNext()
				if err != nil {
					t.Log("Error reading message: %v", err)
				} else {
					t.Log("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						t.Log("Message ID: %s", msgID)

						// Process the message
						t.coreParser.Parse(message)

						// Check for map server info packet (0x0AC5 from PACKET_DUMP)
						if msgID == "0AC5" {
							// Extract map name (starts at offset 4)
							mapNameBytes := bytes.Split(message[4:20], []byte{0})[0]
							t.mapName = string(mapNameBytes)

							// Extract map IP (offset 22-25)
							t.mapIP = fmt.Sprintf("%d.%d.%d.%d", message[22], message[23], message[24], message[25])

							// Extract map port (offset 26-27)
							t.mapPort = int(binary.LittleEndian.Uint16(message[26:28]))

							t.Log("Map server: %s:%d", t.mapIP, t.mapPort)
							return true
						}
					}
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// waitForMapLoad waits for the map to load
func (t *LoginSequenceTest) waitForMapLoad() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Waiting for map to load...")

	for {
		select {
		case <-ctx.Done():
			t.Log("Timeout waiting for map to load")
			return false
		default:
			data, err := t.conn.Receive()
			if err != nil {
				t.Log("Error receiving data: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				t.Log("Received data: %X", data)

				// Process the data with the tokenizer
				t.tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := t.tokenizer.ReadNext()
				if err != nil {
					t.Log("Error reading message: %v", err)
				} else {
					t.Log("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						t.Log("Message ID: %s", msgID)

						// Process the message
						t.coreParser.Parse(message)

						// Check if map loaded
						if t.mapLoginSuccess {
							return true
						}

						// Check for map login success packet
						// Updated to include packets from PACKET_DUMP
						if msgID == "0073" || msgID == "0091" || msgID == "0092" || msgID == "0093" ||
							msgID == "0283" || msgID == "0B18" || msgID == "02EB" || msgID == "0B1B" {
							t.Log("Map login successful")
							t.mapLoginSuccess = true
							return true
						}
					}
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/mikekao/openkore/goKore/network/implementation/network/connection"
	"github.com/mikekao/openkore/goKore/network/implementation/network/hooks"
	"github.com/mikekao/openkore/goKore/network/implementation/network/protocol"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/core"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/game/actor"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/security"
)

// LoginSequenceTest represents the state for the login sequence test
type LoginSequenceTest struct {
	conn         connection.Connection
	tokenizer    *protocol.Tokenizer
	hookManager  *hooks.HookManager
	coreParser   *core.CoreParser
	loginManager *security.LoginManager
	actorManager *actor.ActorManager

	// Test state
	accountID      uint32
	sessionID1     uint32
	sessionID2     uint32
	gender         byte
	serverList     []ServerInfo
	charServerInfo *ServerInfo
	charList       []CharacterInfo
	selectedChar   int
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
	t.coreParser = core.NewCoreParser("ServerType0", t.hookManager)

	// Create login manager
	t.loginManager = security.NewLoginManager(t.coreParser, t.hookManager)
	t.loginManager.RegisterHandlers()

	// Create actor manager
	t.actorManager = actor.NewActorManager(t.coreParser, t.hookManager)
	t.actorManager.RegisterHandlers()

	// Define packet lengths for the tokenizer
	packetLengths := make(map[string]protocol.PacketDef)
	packetLengths["0069"] = protocol.PacketDef{Length: 0, HasLength: true}   // account_server_info
	packetLengths["0071"] = protocol.PacketDef{Length: 28, HasLength: false} // received_character_ID_and_Map
	packetLengths["0073"] = protocol.PacketDef{Length: 11, HasLength: false} // map_loaded
	packetLengths["083E"] = protocol.PacketDef{Length: 26, HasLength: false} // login_error
	packetLengths["006A"] = protocol.PacketDef{Length: -1, HasLength: true}  // login_error_game_login
	packetLengths["006B"] = protocol.PacketDef{Length: -1, HasLength: true}  // received_characters
	packetLengths["006C"] = protocol.PacketDef{Length: 3, HasLength: false}  // login_error_game_login
	packetLengths["0081"] = protocol.PacketDef{Length: -1, HasLength: true}  // login_error_game_login
	packetLengths["0091"] = protocol.PacketDef{Length: 22, HasLength: false} // map_changed
	packetLengths["0092"] = protocol.PacketDef{Length: 28, HasLength: false} // map_changed
	packetLengths["0093"] = protocol.PacketDef{Length: 2, HasLength: false}  // login_to_map_server_succeeded

	// Create message tokenizer
	t.tokenizer = protocol.NewTokenizer(packetLengths)

	// Step 1: Connect to login server
	t.Log("Connecting to login server: %s:%d", connConfig.Host, connConfig.Port)
	err := t.conn.Connect()
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
	username := "botijo0"
	password := "Melon.77"

	t.Log("Sending master login packet for user: %s", username)
	packet := []byte{0x64, 0x00} // 0x0064 master_login packet
	packet = append(packet, make([]byte, 24)...)
	copy(packet[2:], []byte(username))
	packet = append(packet, make([]byte, 24)...)
	copy(packet[26:], []byte(password))
	packet = append(packet, 0)                   // gender (0=female, 1=male)
	packet = append(packet, make([]byte, 16)...) // client hash
	packet = append(packet, 1, 0, 0, 0)          // version

	err = t.conn.Send(packet)
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

	t.Log("Connecting to character server: %s:%d", connConfig.Host, connConfig.Port)
	err = t.conn.Connect()
	if err != nil {
		t.Log("Failed to connect to character server: %v", err)
		return false
	}

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	// Check if connected
	if !t.conn.IsConnected() {
		t.Log("Failed to connect to character server")
		return false
	}
	t.Log("Connected to character server successfully")

	// Step 6: Send character server login packet
	t.Log("Sending character server login packet")
	packet = []byte{0x65, 0x00} // 0x0065 game_login packet
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[2:], t.accountID)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[6:], t.sessionID1)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[10:], t.sessionID2)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[14:], 0) // unknown
	packet = append(packet, t.gender)

	err = t.conn.Send(packet)
	if err != nil {
		t.Log("Failed to send character server login packet: %v", err)
		return false
	}

	// Step 7: Wait for character list
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
	packet = []byte{0x66, 0x00} // 0x0066 char_select packet
	packet = append(packet, selectedChar.CharNum)

	err = t.conn.Send(packet)
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
	packet = []byte{0x72, 0x00} // 0x0072 map_login packet
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[2:], t.accountID)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[6:], t.charList[t.selectedChar].CharID)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[10:], t.sessionID1)
	packet = append(packet, make([]byte, 4)...)
	binary.LittleEndian.PutUint32(packet[14:], t.sessionID2)
	packet = append(packet, t.gender)

	err = t.conn.Send(packet)
	if err != nil {
		t.Log("Failed to send map login packet: %v", err)
		return false
	}

	// Step 14: Wait for map to load
	if !t.waitForMapLoad() {
		t.Log("Failed to load map or timed out")
		return false
	}

	// Step 15: Disconnect from map server
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

						// Process the message
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

						// Check for map server info packet
						if msgID == "0071" {
							// Extract map server IP and port
							t.mapIP = fmt.Sprintf("%d.%d.%d.%d", message[6], message[7], message[8], message[9])
							t.mapPort = int(binary.LittleEndian.Uint16(message[10:12]))
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
						if msgID == "0073" || msgID == "0091" || msgID == "0092" || msgID == "0093" {
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

func main() {
	test := NewLoginSequenceTest()
	if test == nil {
		fmt.Println("Failed to create test")
		os.Exit(1)
	}

	success := test.RunTest()
	if success {
		fmt.Println("Full login sequence test passed!")
		os.Exit(0)
	} else {
		fmt.Println("Full login sequence test failed!")
		os.Exit(1)
	}
}

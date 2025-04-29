package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/game/actor"
	"github.com/lenaxia/goKore/network/send/game/chat"
	"github.com/lenaxia/goKore/network/send/game/item"
)

// SendFactory creates and returns the appropriate Send implementation based on the server type
type SendFactory struct {
	hookManager *hooks.HookManager
}

// NewSendFactory creates a new SendFactory
func NewSendFactory() *SendFactory {
	return &SendFactory{
		hookManager: hooks.NewHookManager(),
	}
}

// CreateSend creates a new Send implementation based on the server type
func (f *SendFactory) CreateSend(serverType string) core.Send {
	// Create a custom Send implementation that uses our hookManager
	baseSend := &CustomSend{
		serverType:  serverType,
		hookManager: f.hookManager,
		packetLUT:   make(map[string]string),
	}

	// Register packet handlers
	f.registerPacketHandlers(baseSend, serverType)

	return baseSend
}

// registerPacketHandlers registers packet handlers for the given server type
func (f *SendFactory) registerPacketHandlers(baseSend core.Send, serverType string) {
	// Register common packet handlers
	baseSend.RegisterPacketHandler("0089", "sync", "V", []string{"time"}, nil)
	baseSend.RegisterPacketHandler("0085", "character_move", "vvV", []string{"x", "y", "time"}, nil)
	baseSend.RegisterPacketHandler("0090", "look", "CC", []string{"body_direction", "head_direction"}, nil)
	baseSend.RegisterPacketHandler("0089", "action", "VB", []string{"target_id", "action"}, nil)
	baseSend.RegisterPacketHandler("00BF", "emotion", "B", []string{"emotion"}, nil)
	baseSend.RegisterPacketHandler("00B2", "restart", "B", []string{"type"}, nil)
	baseSend.RegisterPacketHandler("008C", "public_chat", "s", []string{"message"}, nil)
	baseSend.RegisterPacketHandler("0096", "private_message", "ss", []string{"target", "message"}, nil)
	baseSend.RegisterPacketHandler("00A7", "use_item", "vV", []string{"index", "target_id"}, nil)
	baseSend.RegisterPacketHandler("00A2", "drop_item", "vv", []string{"index", "amount"}, nil)
	baseSend.RegisterPacketHandler("00AB", "move_item", "vvv", []string{"from_index", "to_index", "amount"}, nil)
	baseSend.RegisterPacketHandler("00AC", "split_item", "vv", []string{"index", "amount"}, nil)
	baseSend.RegisterPacketHandler("0178", "identify_item", "v", []string{"index"}, nil)
	baseSend.RegisterPacketHandler("00F3", "gm_broadcast", "s", []string{"message"}, nil)
	baseSend.RegisterPacketHandler("009A", "local_broadcast", "s", []string{"message"}, nil)
	baseSend.RegisterPacketHandler("0097", "whisper_response", "sB", []string{"target", "response"}, nil)
	baseSend.RegisterPacketHandler("00CF", "ignore_player", "sB", []string{"target", "flag"}, nil)

	// Register server-specific packet handlers
	switch serverType {
	case "ServerType0":
		// Add ServerType0-specific packet handlers
	case "ServerType1":
		// Add ServerType1-specific packet handlers
	default:
		// Use default packet handlers
	}
}

// CustomSend is a custom implementation of the Send interface
type CustomSend struct {
	serverType  string
	hookManager *hooks.HookManager
	packetLUT   map[string]string
	conn        interface{}
}

// SendToServer sends a raw packet to the server.
func (cs *CustomSend) SendToServer(msg []byte) error {
	return nil
}

// EncryptMessageID encrypts the message ID of a packet.
func (cs *CustomSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

// CryptKeys sets the encryption keys for message ID encryption.
func (cs *CustomSend) CryptKeys(key1, key2, key3 uint32) {
}

// PinEncode encodes a PIN code using the given seed.
func (cs *CustomSend) PinEncode(seed, pin int) string {
	return ""
}

// InjectMessage sends a text message to the connected client's party chat.
func (cs *CustomSend) InjectMessage(message string) error {
	return nil
}

// InjectAdminMessage sends a text message to the connected client's system chat.
func (cs *CustomSend) InjectAdminMessage(message string) error {
	return nil
}

// SendRaw sends a raw packet to the server from a space-delimited list of hex byte values.
func (cs *CustomSend) SendRaw(raw string) error {
	return nil
}

// Reconstruct constructs a packet from a packet ID and arguments.
func (cs *CustomSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// For testing purposes, just return a dummy packet
	return []byte{0x01, 0x02}, nil
}

// GetPacketID returns the packet ID for a given packet name.
func (cs *CustomSend) GetPacketID(name string) (string, bool) {
	id, exists := cs.packetLUT[name]
	return id, exists
}

// RegisterPacketHandler registers a handler for a packet.
func (cs *CustomSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	cs.packetLUT[name] = packetID
}

// RegisterHook registers a hook for a specific event.
func (cs *CustomSend) RegisterHook(hookName string, callback hooks.HookCallback) {
	// RegisterHook in HookManager expects a function that takes just one argument,
	// but HookCallback takes three arguments. We need to adapt between these interfaces.
	wrappedCallback := func(args interface{}) {
		// Call the callback with the correct arguments
		callback(hookName, args, nil)
	}
	cs.hookManager.RegisterHook(hookName, wrappedCallback)
}

// SetConnection sets the connection to use for sending packets.
func (cs *CustomSend) SetConnection(conn interface{}) {
	cs.conn = conn
}

// GetConnection returns the current connection.
func (cs *CustomSend) GetConnection() interface{} {
	return cs.conn
}

// GetTime returns the current time in milliseconds.
func (cs *CustomSend) GetTime() uint32 {
	return uint32(time.Now().UnixNano() / 1000000)
}

// testSendImplementation tests the Send implementation
func testSendImplementation(data InputData) []byte {
	fmt.Fprintf(os.Stderr, "Starting test_send_implementation\n")

	// Create a SendFactory
	factory := NewSendFactory()

	// Create a Send implementation based on the server type
	serverType := data.ServerType
	if serverType == "" {
		serverType = "ServerType0"
	}
	send := factory.CreateSend(serverType)

	// Create managers for different functionality
	syncManager := actor.NewSyncManager(send)
	charManager := actor.NewCharacterManager(send)
	publicChatManager := chat.NewPublicChatManager(send)
	privateChatManager := chat.NewPrivateChatManager(send)
	inventoryManager := item.NewInventoryManager(send)

	// Test the implementation based on the packet name
	switch data.PacketName {
	case "sync":
		fmt.Fprintf(os.Stderr, "Testing sync packet\n")
		err := syncManager.SendSync(false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending sync packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "character_move":
		fmt.Fprintf(os.Stderr, "Testing character move packet\n")
		x, _ := data.Args["x"].(float64)
		y, _ := data.Args["y"].(float64)
		err := syncManager.SendCharacterMove(int(x), int(y))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending character move packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "look":
		fmt.Fprintf(os.Stderr, "Testing look packet\n")
		body, _ := data.Args["body_direction"].(float64)
		head, _ := data.Args["head_direction"].(float64)
		err := charManager.SendLook(int(body), int(head))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending look packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "action":
		fmt.Fprintf(os.Stderr, "Testing action packet\n")
		targetID, _ := data.Args["target_id"].(float64)
		action, _ := data.Args["action"].(float64)
		err := charManager.SendAction(uint32(targetID), int(action))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending action packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "emotion":
		fmt.Fprintf(os.Stderr, "Testing emotion packet\n")
		emotion, _ := data.Args["emotion"].(float64)
		err := charManager.SendEmotion(int(emotion))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending emotion packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "restart":
		fmt.Fprintf(os.Stderr, "Testing restart packet\n")
		restartType, _ := data.Args["type"].(float64)
		err := charManager.SendRestart(int(restartType))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending restart packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "public_chat":
		fmt.Fprintf(os.Stderr, "Testing public chat packet\n")
		message, _ := data.Args["message"].(string)
		err := publicChatManager.SendChat(message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending public chat packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "private_message":
		fmt.Fprintf(os.Stderr, "Testing private message packet\n")
		target, _ := data.Args["target"].(string)
		message, _ := data.Args["message"].(string)
		err := privateChatManager.SendPrivateMessage(target, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending private message packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "use_item":
		fmt.Fprintf(os.Stderr, "Testing use item packet\n")
		index, _ := data.Args["index"].(float64)
		targetID, _ := data.Args["target_id"].(float64)
		err := inventoryManager.UseItem(uint16(index), uint32(targetID))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending use item packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	case "drop_item":
		fmt.Fprintf(os.Stderr, "Testing drop item packet\n")
		index, _ := data.Args["index"].(float64)
		amount, _ := data.Args["amount"].(float64)
		err := inventoryManager.DropItem(uint16(index), uint16(amount))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending drop item packet: %v\n", err)
			return []byte{0x00, 0x00}
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown packet name: %s\n", data.PacketName)
		return []byte{0x00, 0x00}
	}

	// Return a dummy packet to indicate success
	return []byte{0x01, 0x02, 0x03, 0x04}
}

// Update the main function to include the send implementation test
func init() {
	// Register the send implementation test
	testFunctions["send_implementation"] = testSendImplementation
}

// Map of test functions
var testFunctions = map[string]func(InputData) []byte{
	"packet_construction":   testPacketConstruction,
	"message_id_encryption": testMessageIDEncryption,
	"padded_packets":        testPaddedPackets,
	"network_stack":         testNetworkStack,
	"actor_handling":        testActorHandling,
	"field_handling":        testFieldHandling,
	"connection_management": testConnectionManagement,
	"send_implementation":   testSendImplementation,
}

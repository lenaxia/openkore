package core

import (
	"errors"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/protocol"
)

func TestNewCoreParser(t *testing.T) {
	hookManager := hooks.NewHookManager()
	parser := NewCoreParser("ServerType0", hookManager)

	if parser == nil {
		t.Fatal("NewCoreParser() returned nil")
	}

	if parser.parser == nil {
		t.Error("parser.parser was not initialized")
	}

	if parser.handlers == nil {
		t.Error("parser.handlers was not initialized")
	}

	if parser.hookManager != hookManager {
		t.Error("parser.hookManager was not set correctly")
	}

	if parser.serverType != "ServerType0" {
		t.Errorf("parser.serverType = %s, want ServerType0", parser.serverType)
	}
}

func TestRegisterHandler(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}

	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Verify the handler was registered with the protocol parser
	if _, exists := parser.parser.PacketList[packetID]; !exists {
		t.Errorf("Handler was not registered with protocol parser for packet ID %s", packetID)
	}

	// Verify the handler was stored in our map
	if _, exists := parser.handlers[handlerName]; !exists {
		t.Errorf("Handler was not stored in handlers map for name %s", handlerName)
	}

	// Verify the handler function was stored
	storedHandler := parser.handlers[handlerName]
	storedHandler(nil) // Call the handler
	if !handlerCalled {
		t.Error("Handler function was not called")
	}
}

func TestRegisterHandlerFunc(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler function
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}

	parser.RegisterHandlerFunc(packetID, handlerName, format, paramNames, handler)

	// Verify the handler was registered with the protocol parser
	if _, exists := parser.parser.PacketList[packetID]; !exists {
		t.Errorf("Handler was not registered with protocol parser for packet ID %s", packetID)
	}

	// Verify the handler was stored in our map
	if _, exists := parser.handlers[handlerName]; !exists {
		t.Errorf("Handler was not stored in handlers map for name %s", handlerName)
	}

	// Verify the handler function was stored
	storedHandler := parser.handlers[handlerName]
	storedHandler(nil) // Call the handler
	if !handlerCalled {
		t.Error("Handler function was not called")
	}
}

func TestParse(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, nil)

	// Create a test packet
	packet := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}

	// Parse the packet
	args, err := parser.Parse(packet)
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	// Verify the parsed arguments
	if args["switch"] != packetID {
		t.Errorf("args[\"switch\"] = %s, want %s", args["switch"], packetID)
	}

	if args["param1"] != uint16(12345) {
		t.Errorf("args[\"param1\"] = %v, want %v", args["param1"], uint16(12345))
	}

	if args["param2"] != uint8(67) {
		t.Errorf("args[\"param2\"] = %v, want %v", args["param2"], uint8(67))
	}

	if args["param3"] != uint16(9940) {
		t.Errorf("args[\"param3\"] = %v, want %v", args["param3"], uint16(9940))
	}
}

func TestProcess(t *testing.T) {
	hookManager := hooks.NewHookManager()
	parser := NewCoreParser("ServerType0", hookManager)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Create a test packet
	packet := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}

	// Process the packet
	err := parser.Process(packet)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// Verify the handler was called
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestProcessWithHooks(t *testing.T) {
	hookManager := hooks.NewHookManager()
	parser := NewCoreParser("ServerType0", hookManager)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Register pre-processing hook
	preHookCalled := false
	preHook := func(hookName string, arg interface{}, userData interface{}) {
		preHookCalled = true
		args := arg.(map[string]interface{})
		args["return"] = true // Ignore the packet
	}
	hookManager.AddHook("receive/packet_pre/test_handler", preHook, nil)

	// Create a test packet
	packet := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}

	// Process the packet
	err := parser.Process(packet)
	if err != ErrPacketIgnored {
		t.Fatalf("Process() returned error: %v, want %v", err, ErrPacketIgnored)
	}

	// Verify the pre-hook was called
	if !preHookCalled {
		t.Error("Pre-hook was not called")
	}

	// Verify the handler was not called
	if handlerCalled {
		t.Error("Handler was called despite packet being ignored")
	}

	// Register post-processing hook
	postHookCalled := false
	postHook := func(hookName string, arg interface{}, userData interface{}) {
		postHookCalled = true
	}
	hookManager.AddHook("receive/packet/test_handler", postHook, nil)

	// Reset the pre-hook
	preHook = func(hookName string, arg interface{}, userData interface{}) {
		preHookCalled = true
		// Don't ignore the packet this time
	}
	hookManager.Clear()
	hookManager.AddHook("receive/packet_pre/test_handler", preHook, nil)
	hookManager.AddHook("receive/packet/test_handler", postHook, nil)

	// Reset flags
	preHookCalled = false
	handlerCalled = false
	postHookCalled = false

	// Process the packet again
	err = parser.Process(packet)
	if err != nil {
		t.Fatalf("Process() returned error: %v", err)
	}

	// Verify the pre-hook was called
	if !preHookCalled {
		t.Error("Pre-hook was not called")
	}

	// Verify the handler was called
	if !handlerCalled {
		t.Error("Handler was not called")
	}

	// Verify the post-hook was called
	if !postHookCalled {
		t.Error("Post-hook was not called")
	}
}

func TestProcessBuffer(t *testing.T) {
	hookManager := hooks.NewHookManager()
	parser := NewCoreParser("ServerType0", hookManager)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handlerCalled := false
	handler := func(args map[string]interface{}) error {
		handlerCalled = true
		return nil
	}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Create a tokenizer with packet definitions
	tokenizer := protocol.NewTokenizer(map[string]protocol.PacketDef{
		packetID: {Length: 7, HasLength: false},
	})

	// Add a packet to the tokenizer
	packet := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}
	tokenizer.Add(packet)

	// Process the buffer
	err := parser.ProcessBuffer(tokenizer)
	if err != nil {
		t.Fatalf("ProcessBuffer() returned error: %v", err)
	}

	// Verify the handler was called
	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestProcessError(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler that returns an error
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	expectedErr := errors.New("test error")
	handler := func(args map[string]interface{}) error {
		return expectedErr
	}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Create a test packet
	packet := []byte{0x02, 0x01, 0x39, 0x30, 0x43, 0xD4, 0x26}

	// Process the packet
	err := parser.Process(packet)
	if err != expectedErr {
		t.Fatalf("Process() returned error: %v, want %v", err, expectedErr)
	}
}

func TestGetHandler(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	handler := func(args map[string]interface{}) error {
		return nil
	}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, handler)

	// Get the handler
	storedHandler, exists := parser.GetHandler(handlerName)
	if !exists {
		t.Errorf("GetHandler() returned exists = false, want true")
	}
	if storedHandler == nil {
		t.Error("GetHandler() returned nil handler")
	}

	// Get a non-existent handler
	_, exists = parser.GetHandler("non_existent_handler")
	if exists {
		t.Error("GetHandler() returned exists = true for non-existent handler")
	}
}

func TestGetPacketID(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Register a handler
	packetID := "0102"
	handlerName := "test_handler"
	format := "v1 C1 v1"
	paramNames := []string{"param1", "param2", "param3"}
	parser.RegisterHandler(packetID, handlerName, format, paramNames, nil)

	// Get the packet ID
	storedPacketID, exists := parser.GetPacketID(handlerName)
	if !exists {
		t.Errorf("GetPacketID() returned exists = false, want true")
	}
	if storedPacketID != packetID {
		t.Errorf("GetPacketID() returned packetID = %s, want %s", storedPacketID, packetID)
	}

	// Get a non-existent packet ID
	_, exists = parser.GetPacketID("non_existent_handler")
	if exists {
		t.Error("GetPacketID() returned exists = true for non-existent handler")
	}
}

func TestSetGetDefaultState(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Default state should be 0
	if parser.GetDefaultState() != 0 {
		t.Errorf("GetDefaultState() = %d, want 0", parser.GetDefaultState())
	}

	// Set the default state
	parser.SetDefaultState(5)
	if parser.GetDefaultState() != 5 {
		t.Errorf("GetDefaultState() = %d, want 5", parser.GetDefaultState())
	}
}

func TestGetServerType(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)

	// Get the server type
	if parser.GetServerType() != "ServerType0" {
		t.Errorf("GetServerType() = %s, want ServerType0", parser.GetServerType())
	}
}

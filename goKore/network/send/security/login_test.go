package security

import (
	"bytes"
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// MockSend is a mock implementation of the Send interface for testing
type MockSend struct {
	sentPackets     [][]byte
	packetLUT       map[string]string
	reconstructArgs map[string]map[string]interface{}
	time            uint32
}

func NewMockSend() *MockSend {
	return &MockSend{
		sentPackets:     make([][]byte, 0),
		packetLUT:       make(map[string]string),
		reconstructArgs: make(map[string]map[string]interface{}),
		time:            12345,
	}
}

func (ms *MockSend) SendToServer(msg []byte) error {
	ms.sentPackets = append(ms.sentPackets, msg)
	return nil
}

func (ms *MockSend) EncryptMessageID(msg *[]byte) error {
	return nil
}

func (ms *MockSend) CryptKeys(key1, key2, key3 uint32) {
}

func (ms *MockSend) PinEncode(seed, pin int) string {
	return "1234"
}

func (ms *MockSend) InjectMessage(message string) error {
	return nil
}

func (ms *MockSend) InjectAdminMessage(message string) error {
	return nil
}

func (ms *MockSend) SendRaw(raw string) error {
	return nil
}

func (ms *MockSend) Reconstruct(packetID string, args map[string]interface{}) ([]byte, error) {
	// Store the args for later inspection
	ms.reconstructArgs[packetID] = args

	// Return a simple packet with the ID
	packet := []byte{0x00, 0x00}
	return packet, nil
}

func (ms *MockSend) GetPacketID(name string) (string, bool) {
	id, exists := ms.packetLUT[name]
	return id, exists
}

func (ms *MockSend) RegisterPacketHandler(packetID, name, format string, keys []string, handler func(map[string]interface{}) error) {
	ms.packetLUT[name] = packetID
}

func (ms *MockSend) RegisterHook(hookName string, callback hooks.HookCallback) {
}

func (ms *MockSend) SetConnection(conn interface{}) {
}

func (ms *MockSend) GetConnection() interface{} {
	return nil
}

func (ms *MockSend) GetTime() uint32 {
	return ms.time
}

// TestNewLoginManager tests the NewLoginManager function
func TestNewLoginManager(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	if lm == nil {
		t.Fatal("NewLoginManager() returned nil")
	}

	// We can't directly compare interfaces, so we check if it's nil
	if lm.baseSend == nil {
		t.Error("lm.baseSend was not set correctly")
	}

	if lm.version != 23 {
		t.Errorf("lm.version = %v, want 23", lm.version)
	}

	if lm.masterVersion != 1 {
		t.Errorf("lm.masterVersion = %v, want 1", lm.masterVersion)
	}

	if lm.gameCode != "0011" {
		t.Errorf("lm.gameCode = %v, want 0011", lm.gameCode)
	}

	if lm.flag != "G000" {
		t.Errorf("lm.flag = %v, want G000", lm.flag)
	}

	if lm.mac != "111111111111" {
		t.Errorf("lm.mac = %v, want 111111111111", lm.mac)
	}

	if lm.ip != "192.168.0.2" {
		t.Errorf("lm.ip = %v, want 192.168.0.2", lm.ip)
	}
}

// TestSetCredentials tests the SetCredentials method
func TestSetCredentials(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetCredentials("testuser", "testpass")

	if lm.username != "testuser" {
		t.Errorf("lm.username = %v, want testuser", lm.username)
	}

	if lm.password != "testpass" {
		t.Errorf("lm.password = %v, want testpass", lm.password)
	}
}

// TestSetVersion tests the SetVersion method
func TestSetVersion(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetVersion(42)

	if lm.version != 42 {
		t.Errorf("lm.version = %v, want 42", lm.version)
	}
}

// TestSetMasterVersion tests the SetMasterVersion method
func TestSetMasterVersion(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetMasterVersion(2)

	if lm.masterVersion != 2 {
		t.Errorf("lm.masterVersion = %v, want 2", lm.masterVersion)
	}
}

// TestSetGameCode tests the SetGameCode method
func TestSetGameCode(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetGameCode("1234")

	if lm.gameCode != "1234" {
		t.Errorf("lm.gameCode = %v, want 1234", lm.gameCode)
	}
}

// TestSetFlag tests the SetFlag method
func TestSetFlag(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetFlag("F123")

	if lm.flag != "F123" {
		t.Errorf("lm.flag = %v, want F123", lm.flag)
	}
}

// TestSetMAC tests the SetMAC method
func TestSetMAC(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetMAC("AABBCCDDEEFF")

	if lm.mac != "AABBCCDDEEFF" {
		t.Errorf("lm.mac = %v, want AABBCCDDEEFF", lm.mac)
	}
}

// TestSetIP tests the SetIP method
func TestSetIP(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	lm.SetIP("10.0.0.1")

	if lm.ip != "10.0.0.1" {
		t.Errorf("lm.ip = %v, want 10.0.0.1", lm.ip)
	}
}

// TestSetAccountInfo tests the SetAccountInfo method
func TestSetAccountInfo(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	accountID := []byte{1, 2, 3, 4}
	sessionID := []byte{5, 6, 7, 8}
	sessionID2 := []byte{9, 10, 11, 12}
	accountSex := 1

	lm.SetAccountInfo(accountID, sessionID, sessionID2, accountSex)

	if !bytes.Equal(lm.accountID, accountID) {
		t.Errorf("lm.accountID = %v, want %v", lm.accountID, accountID)
	}

	if !bytes.Equal(lm.sessionID, sessionID) {
		t.Errorf("lm.sessionID = %v, want %v", lm.sessionID, sessionID)
	}

	if !bytes.Equal(lm.sessionID2, sessionID2) {
		t.Errorf("lm.sessionID2 = %v, want %v", lm.sessionID2, sessionID2)
	}

	if lm.accountSex != accountSex {
		t.Errorf("lm.accountSex = %v, want %v", lm.accountSex, accountSex)
	}
}

// TestSendMasterLogin tests the SendMasterLogin method
func TestSendMasterLogin(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.RegisterPacketHandler("0064", "master_login", "", nil, nil)

	lm := NewLoginManager(mockSend)
	lm.SetCredentials("testuser", "testpass")
	lm.SetVersion(42)
	lm.SetMasterVersion(2)
	lm.SetGameCode("1234")
	lm.SetFlag("F123")

	err := lm.SendMasterLogin()
	if err != nil {
		t.Fatalf("SendMasterLogin() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0064"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if args["version"] != 42 {
		t.Errorf("args[\"version\"] = %v, want 42", args["version"])
	}

	if args["master_version"] != 2 {
		t.Errorf("args[\"master_version\"] = %v, want 2", args["master_version"])
	}

	if args["username"] != "testuser" {
		t.Errorf("args[\"username\"] = %v, want testuser", args["username"])
	}

	if args["password"] != "testpass" {
		t.Errorf("args[\"password\"] = %v, want testpass", args["password"])
	}

	if args["game_code"] != "1234" {
		t.Errorf("args[\"game_code\"] = %v, want 1234", args["game_code"])
	}

	if args["flag"] != "F123" {
		t.Errorf("args[\"flag\"] = %v, want F123", args["flag"])
	}
}

// TestSecureLoginHash tests the secureLoginHash method
func TestSecureLoginHash(t *testing.T) {
	mockSend := NewMockSend()
	lm := NewLoginManager(mockSend)

	// Test with loginType = 1 (salt + password)
	salt := []byte{1, 2, 3, 4}
	password := "testpass"
	hash1 := lm.secureLoginHash(password, salt, 1)

	// Test with loginType = 2 (password + salt)
	hash2 := lm.secureLoginHash(password, salt, 2)

	// The hashes should be different
	if bytes.Equal(hash1, hash2) {
		t.Error("hash1 and hash2 are equal, but should be different")
	}
}

// TestSendMapLogin tests the SendMapLogin method
func TestSendMapLogin(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.RegisterPacketHandler("0072", "map_login", "", nil, nil)

	lm := NewLoginManager(mockSend)
	accountID := []byte{1, 2, 3, 4}
	sessionID := []byte{5, 6, 7, 8}
	sessionID2 := []byte{9, 10, 11, 12}
	accountSex := 1
	lm.SetAccountInfo(accountID, sessionID, sessionID2, accountSex)

	err := lm.SendMapLogin()
	if err != nil {
		t.Fatalf("SendMapLogin() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were correct
	args, exists := mockSend.reconstructArgs["0072"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	if !bytes.Equal(args["accountID"].([]byte), accountID) {
		t.Errorf("args[\"accountID\"] = %v, want %v", args["accountID"], accountID)
	}

	if !bytes.Equal(args["charID"].([]byte), accountID) {
		t.Errorf("args[\"charID\"] = %v, want %v", args["charID"], accountID)
	}

	if !bytes.Equal(args["sessionID"].([]byte), sessionID) {
		t.Errorf("args[\"sessionID\"] = %v, want %v", args["sessionID"], sessionID)
	}

	if args["tick"] != uint32(12345) {
		t.Errorf("args[\"tick\"] = %v, want 12345", args["tick"])
	}

	if args["sex"] != accountSex {
		t.Errorf("args[\"sex\"] = %v, want %v", args["sex"], accountSex)
	}
}

// TestSendMapLoaded tests the SendMapLoaded method
func TestSendMapLoaded(t *testing.T) {
	mockSend := NewMockSend()
	mockSend.RegisterPacketHandler("007D", "map_loaded", "", nil, nil)

	lm := NewLoginManager(mockSend)

	err := lm.SendMapLoaded()
	if err != nil {
		t.Fatalf("SendMapLoaded() returned error: %v", err)
	}

	// Check that the packet was sent
	if len(mockSend.sentPackets) != 1 {
		t.Fatalf("len(mockSend.sentPackets) = %v, want 1", len(mockSend.sentPackets))
	}

	// Check that the arguments were passed to Reconstruct
	_, exists := mockSend.reconstructArgs["007D"]
	if !exists {
		t.Fatal("No arguments were passed to Reconstruct")
	}

	// The map_loaded packet doesn't have any arguments, so we don't need to check them
}

package security

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

func TestNewLoginManager(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	if manager == nil {
		t.Fatal("NewLoginManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.hookManager != hookManager {
		t.Error("manager.hookManager was not set correctly")
	}

	if manager.state != LoginStateDisconnected {
		t.Errorf("manager.state = %v, want %v", manager.state, LoginStateDisconnected)
	}
}

func TestRegisterHandlers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"account_server_info",
		"login_error",
		"login_error_game_login_server",
		"secure_login_key",
		"received_login_token",
		"login_pin_code_request",
		"login_pin_new_code_result",
	}

	for _, name := range handlerNames {
		if _, exists := parser.GetHandler(name); !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestSetGetCredentials(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Set credentials
	username := "testuser"
	password := "testpass"
	manager.SetCredentials(username, password)

	// Check credentials
	if manager.username != username {
		t.Errorf("manager.username = %s, want %s", manager.username, username)
	}

	if manager.password != password {
		t.Errorf("manager.password = %s, want %s", manager.password, password)
	}
}

func TestSetGetState(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Set state
	manager.SetState(LoginStateLoggedIn)

	// Get state
	state := manager.GetState()
	if state != LoginStateLoggedIn {
		t.Errorf("GetState() = %v, want %v", state, LoginStateLoggedIn)
	}
}

func TestSelectGetServer(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Add some test servers
	manager.servers = []ServerInfo{
		{
			IP:   "127.0.0.1",
			Port: 6900,
			Name: "Test Server 1",
		},
		{
			IP:   "127.0.0.2",
			Port: 6900,
			Name: "Test Server 2",
		},
	}

	// Select server
	err := manager.SelectServer(1)
	if err != nil {
		t.Fatalf("SelectServer() returned error: %v", err)
	}

	// Get selected server
	server, err := manager.GetSelectedServer()
	if err != nil {
		t.Fatalf("GetSelectedServer() returned error: %v", err)
	}

	if server.Name != "Test Server 2" {
		t.Errorf("server.Name = %s, want %s", server.Name, "Test Server 2")
	}

	// Try to select invalid server
	err = manager.SelectServer(2)
	if err == nil {
		t.Error("SelectServer() with invalid index did not return error")
	}
}

func TestGetServers(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Add some test servers
	testServers := []ServerInfo{
		{
			IP:   "127.0.0.1",
			Port: 6900,
			Name: "Test Server 1",
		},
		{
			IP:   "127.0.0.2",
			Port: 6900,
			Name: "Test Server 2",
		},
	}
	manager.servers = testServers

	// Get servers
	servers := manager.GetServers()
	if len(servers) != len(testServers) {
		t.Errorf("len(GetServers()) = %d, want %d", len(servers), len(testServers))
	}

	for i, server := range servers {
		if server.Name != testServers[i].Name {
			t.Errorf("servers[%d].Name = %s, want %s", i, server.Name, testServers[i].Name)
		}
	}
}

func TestGetSessionIDs(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Set session IDs
	manager.sessionID = []byte{1, 2, 3, 4}
	manager.sessionID2 = []byte{5, 6, 7, 8}
	manager.accountID = 12345

	// Get session IDs
	sessionID, sessionID2, accountID := manager.GetSessionIDs()
	if len(sessionID) != 4 || sessionID[0] != 1 || sessionID[1] != 2 || sessionID[2] != 3 || sessionID[3] != 4 {
		t.Errorf("sessionID = %v, want [1 2 3 4]", sessionID)
	}

	if len(sessionID2) != 4 || sessionID2[0] != 5 || sessionID2[1] != 6 || sessionID2[2] != 7 || sessionID2[3] != 8 {
		t.Errorf("sessionID2 = %v, want [5 6 7 8]", sessionID2)
	}

	if accountID != 12345 {
		t.Errorf("accountID = %d, want %d", accountID, 12345)
	}
}

func TestGetSecureKey(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Set secure key
	manager.secureKey = []byte{1, 2, 3, 4}

	// Get secure key
	secureKey := manager.GetSecureKey()
	if len(secureKey) != 4 || secureKey[0] != 1 || secureKey[1] != 2 || secureKey[2] != 3 || secureKey[3] != 4 {
		t.Errorf("secureKey = %v, want [1 2 3 4]", secureKey)
	}
}

func TestGetLoginError(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Set login error
	manager.loginError = &LoginError{
		Code:    1,
		Message: "Test error",
		Date:    "2023-01-01",
	}

	// Get login error
	loginError := manager.GetLoginError()
	if loginError.Code != 1 {
		t.Errorf("loginError.Code = %d, want %d", loginError.Code, 1)
	}

	if loginError.Message != "Test error" {
		t.Errorf("loginError.Message = %s, want %s", loginError.Message, "Test error")
	}

	if loginError.Date != "2023-01-01" {
		t.Errorf("loginError.Date = %s, want %s", loginError.Date, "2023-01-01")
	}
}

func TestIsLoggedIn(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Test not logged in states
	notLoggedInStates := []LoginState{
		LoginStateDisconnected,
		LoginStateConnecting,
		LoginStateHandshaking,
		LoginStateLoggingIn,
	}

	for _, state := range notLoggedInStates {
		manager.state = state
		if manager.IsLoggedIn() {
			t.Errorf("IsLoggedIn() = true, want false for state %v", state)
		}
	}

	// Test logged in states
	loggedInStates := []LoginState{
		LoginStateLoggedIn,
		LoginStateSelectingServer,
		LoginStateServerSelected,
	}

	for _, state := range loggedInStates {
		manager.state = state
		if !manager.IsLoggedIn() {
			t.Errorf("IsLoggedIn() = false, want true for state %v", state)
		}
	}
}

func TestIsSessionExpired(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Session should not be expired initially
	if manager.IsSessionExpired(1 * time.Second) {
		t.Error("IsSessionExpired() = true, want false for new session")
	}

	// Set last activity to a long time ago
	manager.lastActivity = time.Now().Add(-2 * time.Second)

	// Session should be expired now
	if !manager.IsSessionExpired(1 * time.Second) {
		t.Error("IsSessionExpired() = false, want true for expired session")
	}
}

func TestUpdateActivity(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Get initial last activity time
	initialTime := manager.lastActivity

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Update activity
	manager.UpdateActivity()

	// Check that last activity time was updated
	if !manager.lastActivity.After(initialTime) {
		t.Error("lastActivity was not updated")
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Generate password hash
	password := "testpass"
	hash := manager.GeneratePasswordHash(password)

	// Check that hash is correct
	expectedHash := "179ad45c6ce2cb97cf1029e212046e81" // MD5 hash of "testpass"
	if hash != expectedHash {
		t.Errorf("GeneratePasswordHash() = %s, want %s", hash, expectedHash)
	}
}

func TestLoginStateString(t *testing.T) {
	testCases := []struct {
		state LoginState
		want  string
	}{
		{LoginStateDisconnected, "Disconnected"},
		{LoginStateConnecting, "Connecting"},
		{LoginStateHandshaking, "Handshaking"},
		{LoginStateLoggingIn, "LoggingIn"},
		{LoginStateLoggedIn, "LoggedIn"},
		{LoginStateSelectingServer, "SelectingServer"},
		{LoginStateServerSelected, "ServerSelected"},
		{LoginState(99), "Unknown"},
	}

	for _, tc := range testCases {
		got := tc.state.String()
		if got != tc.want {
			t.Errorf("LoginState(%d).String() = %s, want %s", tc.state, got, tc.want)
		}
	}
}

func TestHandleAccountServerInfo(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"sessionID":     []byte{1, 2, 3, 4},
		"accountID":     []byte{5, 6, 7, 8},
		"sessionID2":    []byte{9, 10, 11, 12},
		"lastLoginIP":   []byte{127, 0, 0, 1},
		"lastLoginTime": "2023-01-01",
		"accountSex":    byte(1),
		"serverInfo":    []byte{},
	}

	// Call handler
	err := manager.handleAccountServerInfo(args)
	if err != nil {
		t.Fatalf("handleAccountServerInfo() returned error: %v", err)
	}

	// Check that session was updated
	if len(manager.sessionID) != 4 || manager.sessionID[0] != 1 || manager.sessionID[1] != 2 || manager.sessionID[2] != 3 || manager.sessionID[3] != 4 {
		t.Errorf("manager.sessionID = %v, want [1 2 3 4]", manager.sessionID)
	}

	if manager.accountID != 0x08070605 {
		t.Errorf("manager.accountID = %v, want %v", manager.accountID, 0x08070605)
	}

	if len(manager.sessionID2) != 4 || manager.sessionID2[0] != 9 || manager.sessionID2[1] != 10 || manager.sessionID2[2] != 11 || manager.sessionID2[3] != 12 {
		t.Errorf("manager.sessionID2 = %v, want [9 10 11 12]", manager.sessionID2)
	}

	if manager.lastLoginIP != "127.0.0.1" {
		t.Errorf("manager.lastLoginIP = %s, want 127.0.0.1", manager.lastLoginIP)
	}

	if manager.lastLoginTime != "2023-01-01" {
		t.Errorf("manager.lastLoginTime = %s, want 2023-01-01", manager.lastLoginTime)
	}

	if manager.accountSex != 1 {
		t.Errorf("manager.accountSex = %d, want 1", manager.accountSex)
	}

	if manager.state != LoginStateLoggedIn {
		t.Errorf("manager.state = %v, want %v", manager.state, LoginStateLoggedIn)
	}
}

func TestHandleLoginError(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"type": byte(1),
		"date": "2023-01-01",
	}

	// Call handler
	err := manager.handleLoginError(args)
	if err != nil {
		t.Fatalf("handleLoginError() returned error: %v", err)
	}

	// Check that login error was set
	if manager.loginError == nil {
		t.Fatal("manager.loginError is nil")
	}

	if manager.loginError.Code != 1 {
		t.Errorf("manager.loginError.Code = %d, want 1", manager.loginError.Code)
	}

	if manager.loginError.Date != "2023-01-01" {
		t.Errorf("manager.loginError.Date = %s, want 2023-01-01", manager.loginError.Date)
	}

	if manager.state != LoginStateDisconnected {
		t.Errorf("manager.state = %v, want %v", manager.state, LoginStateDisconnected)
	}
}

func TestHandleSecureLoginKey(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"secure_key": []byte{1, 2, 3, 4},
	}

	// Call handler
	err := manager.handleSecureLoginKey(args)
	if err != nil {
		t.Fatalf("handleSecureLoginKey() returned error: %v", err)
	}

	// Check that secure key was set
	if len(manager.secureKey) != 4 || manager.secureKey[0] != 1 || manager.secureKey[1] != 2 || manager.secureKey[2] != 3 || manager.secureKey[3] != 4 {
		t.Errorf("manager.secureKey = %v, want [1 2 3 4]", manager.secureKey)
	}

	if manager.state != LoginStateHandshaking {
		t.Errorf("manager.state = %v, want %v", manager.state, LoginStateHandshaking)
	}
}

func TestParseServerInfo(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	// Create test server info
	serverInfo := []byte{
		// Server 1
		127, 0, 0, 1, // IP
		0x50, 0x1A, // Port (6736)
		'T', 'e', 's', 't', ' ', 'S', 'e', 'r', 'v', 'e', 'r', ' ', '1', 0, 0, 0, 0, 0, 0, 0, // Name
		0x64, 0x00, // Users (100)
		0x00, 0x00, // State (0)
		0x01, 0x00, // Property (1 - New)

		// Server 2
		127, 0, 0, 2, // IP
		0x50, 0x1A, // Port (6736)
		'T', 'e', 's', 't', ' ', 'S', 'e', 'r', 'v', 'e', 'r', ' ', '2', 0, 0, 0, 0, 0, 0, 0, // Name
		0x32, 0x00, // Users (50)
		0x00, 0x00, // State (0)
		0x02, 0x00, // Property (2 - PvP)
	}

	// Parse server info
	manager.parseServerInfo(serverInfo)

	// Check that servers were parsed correctly
	if len(manager.servers) != 2 {
		t.Fatalf("len(manager.servers) = %d, want 2", len(manager.servers))
	}

	// Check server 1
	if manager.servers[0].IP != "127.0.0.1" {
		t.Errorf("manager.servers[0].IP = %s, want 127.0.0.1", manager.servers[0].IP)
	}

	if manager.servers[0].Port != 6736 {
		t.Errorf("manager.servers[0].Port = %d, want 6736", manager.servers[0].Port)
	}

	if manager.servers[0].Name != "Test Server 1" {
		t.Errorf("manager.servers[0].Name = %s, want Test Server 1", manager.servers[0].Name)
	}

	if manager.servers[0].Users != 100 {
		t.Errorf("manager.servers[0].Users = %d, want 100", manager.servers[0].Users)
	}

	if !manager.servers[0].IsNew {
		t.Error("manager.servers[0].IsNew = false, want true")
	}

	// Check server 2
	if manager.servers[1].IP != "127.0.0.2" {
		t.Errorf("manager.servers[1].IP = %s, want 127.0.0.2", manager.servers[1].IP)
	}

	if manager.servers[1].Port != 6736 {
		t.Errorf("manager.servers[1].Port = %d, want 6736", manager.servers[1].Port)
	}

	if manager.servers[1].Name != "Test Server 2" {
		t.Errorf("manager.servers[1].Name = %s, want Test Server 2", manager.servers[1].Name)
	}

	if manager.servers[1].Users != 50 {
		t.Errorf("manager.servers[1].Users = %d, want 50", manager.servers[1].Users)
	}

	if !manager.servers[1].IsPvP {
		t.Error("manager.servers[1].IsPvP = false, want true")
	}
}

func TestGetLoginErrorMessage(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", nil)
	manager := NewLoginManager(parser, nil)

	testCases := []struct {
		code int
		want string
	}{
		{0, "Account name doesn't exist"},
		{1, "Password Error"},
		{2, "Already logged in."},
		{3, "The server has denied your connection"},
		{4, "Critical Error: Your account has been blocked"},
		{5, "Connect failed, something is wrong with the login settings"},
		{99, "Account not found or password incorrect."},
		{999, "The server has denied your connection for unknown reason (999)"},
	}

	for _, tc := range testCases {
		got := manager.getLoginErrorMessage(tc.code)
		if got != tc.want {
			t.Errorf("getLoginErrorMessage(%d) = %s, want %s", tc.code, got, tc.want)
		}
	}
}

func TestHandleAccountID(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"accountID": []byte{1, 2, 3, 4},
	}

	// Set up hook to verify it's called
	var hookCalled bool
	var hookAccountID uint32

	hookManager.AddHook("security/account_id", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
		if data, ok := arg.(map[string]interface{}); ok {
			if accountID, ok := data["accountID"].(uint32); ok {
				hookAccountID = accountID
			}
		}
	}, nil)

	// Call handler
	err := manager.handleAccountID(args)
	if err != nil {
		t.Fatalf("handleAccountID() returned error: %v", err)
	}

	// Check that account ID was set
	expectedAccountID := uint32(0x04030201) // Little-endian byte order
	if manager.accountID != expectedAccountID {
		t.Errorf("manager.accountID = %d, want %d", manager.accountID, expectedAccountID)
	}

	// Check that hook was called with correct account ID
	if !hookCalled {
		t.Error("Hook was not called")
	}

	if hookAccountID != expectedAccountID {
		t.Errorf("Hook accountID = %d, want %d", hookAccountID, expectedAccountID)
	}
}

func TestHandleAccountPaymentInfo(t *testing.T) {
	parser := core.NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	manager := NewLoginManager(parser, hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"D_minute": uint32(4320), // 3 days, 0 hours, 0 minutes
		"H_minute": uint32(1500), // 1 day, 1 hour, 0 minutes
	}

	// Set up hook to verify it's called
	var hookCalled bool
	var hookData map[string]interface{}

	hookManager.AddHook("security/account_payment_info", func(hookName string, arg interface{}, userData interface{}) {
		hookCalled = true
		if data, ok := arg.(map[string]interface{}); ok {
			hookData = data
		}
	}, nil)

	// Call handler
	err := manager.handleAccountPaymentInfo(args)
	if err != nil {
		t.Fatalf("handleAccountPaymentInfo() returned error: %v", err)
	}

	// Check that hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}

	// Check calculated values
	expectedValues := map[string]uint32{
		"D_minute": 4320,
		"H_minute": 1500,
		"D_days":   3,
		"D_hours":  0,
		"D_mins":   0,
		"H_days":   1,
		"H_hours":  1,
		"H_mins":   0,
	}

	for key, expected := range expectedValues {
		if val, ok := hookData[key].(uint32); !ok || val != expected {
			t.Errorf("hookData[%s] = %v, want %v", key, hookData[key], expected)
		}
	}
}

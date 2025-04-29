// Package security provides security-related functionality for the network stack.
package security

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginFailed        = errors.New("login failed")
	ErrSessionExpired     = errors.New("session expired")
	ErrServerUnavailable  = errors.New("server unavailable")
)

// LoginState represents the state of the login process
type LoginState int

// Login states
const (
	LoginStateDisconnected LoginState = iota
	LoginStateConnecting
	LoginStateHandshaking
	LoginStateLoggingIn
	LoginStateLoggedIn
	LoginStateSelectingServer
	LoginStateServerSelected
)

// String returns the string representation of the login state
func (s LoginState) String() string {
	switch s {
	case LoginStateDisconnected:
		return "Disconnected"
	case LoginStateConnecting:
		return "Connecting"
	case LoginStateHandshaking:
		return "Handshaking"
	case LoginStateLoggingIn:
		return "LoggingIn"
	case LoginStateLoggedIn:
		return "LoggedIn"
	case LoginStateSelectingServer:
		return "SelectingServer"
	case LoginStateServerSelected:
		return "ServerSelected"
	default:
		return "Unknown"
	}
}

// LoginError represents an error that occurred during login
type LoginError struct {
	Code    int
	Message string
	Date    string
}

// ServerInfo represents information about a game server
type ServerInfo struct {
	IP        string
	Port      int
	Name      string
	Users     int
	State     int
	Property  int
	IPPort    string
	ServerID  int
	Unknown   int
	IsNew     bool
	IsPvP     bool
	IsLocked  bool
	IsPremium bool
}

// LoginManager manages login-related functionality
type LoginManager struct {
	parser         *core.CoreParser
	hookManager    *hooks.HookManager
	state          LoginState
	mutex          sync.RWMutex
	username       string
	password       string
	secureKey      []byte
	sessionID      []byte
	sessionID2     []byte
	accountID      uint32
	lastLoginIP    string
	lastLoginTime  string
	accountSex     byte
	servers        []ServerInfo
	selectedServer int
	loginError     *LoginError
	lastActivity   time.Time
}

// NewLoginManager creates a new login manager
func NewLoginManager(parser *core.CoreParser, hookManager *hooks.HookManager) *LoginManager {
	return &LoginManager{
		parser:       parser,
		hookManager:  hookManager,
		state:        LoginStateDisconnected,
		lastActivity: time.Now(),
	}
}

// RegisterHandlers registers login-related packet handlers
func (m *LoginManager) RegisterHandlers() {
	// Register handlers for login-related packets
	m.parser.RegisterHandlerFunc("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*",
		[]string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		m.handleAccountServerInfo)

	m.parser.RegisterHandlerFunc("006A", "login_error", "C Z20",
		[]string{"type", "date"},
		m.handleLoginError)

	m.parser.RegisterHandlerFunc("006C", "login_error_game_login_server", "",
		[]string{},
		m.handleLoginErrorGameLoginServer)

	m.parser.RegisterHandlerFunc("01DC", "secure_login_key", "x2 a*",
		[]string{"secure_key"},
		m.handleSecureLoginKey)

	m.parser.RegisterHandlerFunc("02CA", "login_error_game_login_server", "C",
		[]string{"type"},
		m.handleLoginErrorGameLoginServer)

	m.parser.RegisterHandlerFunc("083E", "login_error", "V Z20",
		[]string{"type", "date"},
		m.handleLoginError)

	m.parser.RegisterHandlerFunc("0ACD", "login_error", "C Z20",
		[]string{"type", "date"},
		m.handleLoginError)

	m.parser.RegisterHandlerFunc("0AE0", "login_error", "V V Z20",
		[]string{"type", "error", "date"},
		m.handleLoginError)

	m.parser.RegisterHandlerFunc("0AE3", "received_login_token", "v l Z20 Z*",
		[]string{"len", "login_type", "flag", "login_token"},
		m.handleReceivedLoginToken)

	m.parser.RegisterHandlerFunc("02AD", "login_pin_code_request", "v V",
		[]string{"flag", "key"},
		m.handleLoginPinCodeRequest)

	m.parser.RegisterHandlerFunc("08B9", "login_pin_code_request", "V a4 v",
		[]string{"seed", "accountID", "flag"},
		m.handleLoginPinCodeRequest)

	m.parser.RegisterHandlerFunc("08BB", "login_pin_new_code_result", "v V",
		[]string{"flag", "seed"},
		m.handleLoginPinNewCodeResult)

	m.parser.RegisterHandlerFunc("0AE9", "login_pin_code_request", "V a4 v2",
		[]string{"seed", "accountID", "flag", "lock"},
		m.handleLoginPinCodeRequest)

	m.parser.RegisterHandlerFunc("0AC4", "account_id", "a4",
		[]string{"accountID"},
		m.handleAccountID)

	m.parser.RegisterHandlerFunc("02A2", "account_payment_info", "V2",
		[]string{"D_minute", "H_minute"},
		m.handleAccountPaymentInfo)
}

// handleAccountID handles the account_id packet
func (m *LoginManager) handleAccountID(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if accountIDVal, ok := args["accountID"].([]byte); ok && len(accountIDVal) >= 4 {
		m.accountID = uint32(accountIDVal[0]) | uint32(accountIDVal[1])<<8 | uint32(accountIDVal[2])<<16 | uint32(accountIDVal[3])<<24

		// Call hook
		if m.hookManager != nil {
			m.hookManager.CallHook("security/account_id", map[string]interface{}{
				"accountID": m.accountID,
			})
		}
	}

	return nil
}

// handleAccountPaymentInfo handles the account_payment_info packet
func (m *LoginManager) handleAccountPaymentInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var dMinute, hMinute uint32

	if dMinuteVal, ok := args["D_minute"].(uint32); ok {
		dMinute = dMinuteVal
	}

	if hMinuteVal, ok := args["H_minute"].(uint32); ok {
		hMinute = hMinuteVal
	}

	// Calculate days, hours, and minutes
	dDays := dMinute / 1440
	dHours := (dMinute % 1440) / 60
	dMins := (dMinute % 1440) % 60

	hDays := hMinute / 1440
	hHours := (hMinute % 1440) / 60
	hMins := (hMinute % 1440) % 60

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/account_payment_info", map[string]interface{}{
			"D_minute": dMinute,
			"H_minute": hMinute,
			"D_days":   dDays,
			"D_hours":  dHours,
			"D_mins":   dMins,
			"H_days":   hDays,
			"H_hours":  hHours,
			"H_mins":   hMins,
		})
	}

	return nil
}

// SetCredentials sets the login credentials
func (m *LoginManager) SetCredentials(username, password string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.username = username
	m.password = password
}

// GetState returns the current login state
func (m *LoginManager) GetState() LoginState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state
}

// SetState sets the login state
func (m *LoginManager) SetState(state LoginState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.state = state
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/login_state_change", map[string]interface{}{
			"state": state,
		})
	}
}

// GetServers returns the list of available servers
func (m *LoginManager) GetServers() []ServerInfo {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.servers
}

// SelectServer selects a server by index
func (m *LoginManager) SelectServer(index int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if index < 0 || index >= len(m.servers) {
		return fmt.Errorf("invalid server index: %d", index)
	}

	m.selectedServer = index
	m.state = LoginStateServerSelected
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/server_selected", map[string]interface{}{
			"server": m.servers[index],
			"index":  index,
		})
	}

	return nil
}

// GetSelectedServer returns the selected server
func (m *LoginManager) GetSelectedServer() (ServerInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.selectedServer < 0 || m.selectedServer >= len(m.servers) {
		return ServerInfo{}, fmt.Errorf("no server selected")
	}

	return m.servers[m.selectedServer], nil
}

// GetSessionIDs returns the session IDs
func (m *LoginManager) GetSessionIDs() ([]byte, []byte, uint32) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.sessionID, m.sessionID2, m.accountID
}

// GetSecureKey returns the secure key
func (m *LoginManager) GetSecureKey() []byte {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.secureKey
}

// GetLoginError returns the last login error
func (m *LoginManager) GetLoginError() *LoginError {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.loginError
}

// IsLoggedIn returns whether the user is logged in
func (m *LoginManager) IsLoggedIn() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.state >= LoginStateLoggedIn
}

// IsSessionExpired checks if the session has expired
func (m *LoginManager) IsSessionExpired(timeout time.Duration) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return time.Since(m.lastActivity) > timeout
}

// UpdateActivity updates the last activity time
func (m *LoginManager) UpdateActivity() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.lastActivity = time.Now()
}

// GeneratePasswordHash generates a hash of the password
func (m *LoginManager) GeneratePasswordHash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Packet handlers

// handleAccountServerInfo handles the account_server_info packet
func (m *LoginManager) handleAccountServerInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract session IDs and account ID
	if sessionID, ok := args["sessionID"].([]byte); ok {
		m.sessionID = sessionID
	}

	if accountID, ok := args["accountID"].([]byte); ok && len(accountID) >= 4 {
		m.accountID = uint32(accountID[0]) | uint32(accountID[1])<<8 | uint32(accountID[2])<<16 | uint32(accountID[3])<<24
	}

	if sessionID2, ok := args["sessionID2"].([]byte); ok {
		m.sessionID2 = sessionID2
	}

	if lastLoginIP, ok := args["lastLoginIP"].([]byte); ok {
		m.lastLoginIP = fmt.Sprintf("%d.%d.%d.%d", lastLoginIP[0], lastLoginIP[1], lastLoginIP[2], lastLoginIP[3])
	}

	if lastLoginTime, ok := args["lastLoginTime"].(string); ok {
		m.lastLoginTime = lastLoginTime
	}

	if accountSex, ok := args["accountSex"].(byte); ok {
		m.accountSex = accountSex
	}

	// Parse server info
	if serverInfo, ok := args["serverInfo"].([]byte); ok {
		m.parseServerInfo(serverInfo)
	}

	// Update state
	m.state = LoginStateLoggedIn
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/login_success", map[string]interface{}{
			"accountID":     m.accountID,
			"lastLoginIP":   m.lastLoginIP,
			"lastLoginTime": m.lastLoginTime,
			"accountSex":    m.accountSex,
			"servers":       m.servers,
		})
	}

	return nil
}

// handleLoginError handles the login_error packet
func (m *LoginManager) handleLoginError(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errorCode int
	var errorMessage string
	var date string

	if errorType, ok := args["type"].(byte); ok {
		errorCode = int(errorType)
	} else if errorType, ok := args["type"].(uint32); ok {
		errorCode = int(errorType)
	}

	if errorVal, ok := args["error"].(uint32); ok {
		errorCode = int(errorVal)
	}

	if dateVal, ok := args["date"].(string); ok {
		date = dateVal
	}

	// Set error message based on error code
	errorMessage = m.getLoginErrorMessage(errorCode)

	// Create login error
	m.loginError = &LoginError{
		Code:    errorCode,
		Message: errorMessage,
		Date:    date,
	}

	// Update state
	m.state = LoginStateDisconnected
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/login_error", map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
			"date":    date,
		})
	}

	return nil
}

// handleLoginErrorGameLoginServer handles the login_error_game_login_server packet
func (m *LoginManager) handleLoginErrorGameLoginServer(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errorCode int

	if errorType, ok := args["type"].(byte); ok {
		errorCode = int(errorType)
	}

	// Set error message based on error code
	errorMessage := m.getLoginErrorMessage(errorCode)

	// Create login error
	m.loginError = &LoginError{
		Code:    errorCode,
		Message: errorMessage,
	}

	// Update state
	m.state = LoginStateDisconnected
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/login_error", map[string]interface{}{
			"code":    errorCode,
			"message": errorMessage,
		})
	}

	return nil
}

// handleSecureLoginKey handles the secure_login_key packet
func (m *LoginManager) handleSecureLoginKey(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if secureKey, ok := args["secure_key"].([]byte); ok {
		m.secureKey = secureKey
	}

	// Update state
	m.state = LoginStateHandshaking
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/secure_login_key", map[string]interface{}{
			"secure_key": m.secureKey,
		})
	}

	return nil
}

// handleReceivedLoginToken handles the received_login_token packet
func (m *LoginManager) handleReceivedLoginToken(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var loginType int
	var flag string
	var loginToken string

	if loginTypeVal, ok := args["login_type"].(int32); ok {
		loginType = int(loginTypeVal)
	}

	if flagVal, ok := args["flag"].(string); ok {
		flag = flagVal
	}

	if loginTokenVal, ok := args["login_token"].(string); ok {
		loginToken = loginTokenVal
	}

	// Update state
	m.state = LoginStateHandshaking
	m.lastActivity = time.Now()

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/login_token", map[string]interface{}{
			"login_type":  loginType,
			"flag":        flag,
			"login_token": loginToken,
		})
	}

	return nil
}

// handleLoginPinCodeRequest handles the login_pin_code_request packet
func (m *LoginManager) handleLoginPinCodeRequest(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var seed uint32
	var accountID uint32
	var flag int
	var lock int

	if seedVal, ok := args["seed"].(uint32); ok {
		seed = seedVal
	}

	if accountIDVal, ok := args["accountID"].([]byte); ok && len(accountIDVal) >= 4 {
		accountID = uint32(accountIDVal[0]) | uint32(accountIDVal[1])<<8 | uint32(accountIDVal[2])<<16 | uint32(accountIDVal[3])<<24
	}

	if flagVal, ok := args["flag"].(uint16); ok {
		flag = int(flagVal)
	} else if flagVal, ok := args["flag"].(byte); ok {
		flag = int(flagVal)
	}

	if lockVal, ok := args["lock"].(uint16); ok {
		lock = int(lockVal)
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_code_request", map[string]interface{}{
			"seed":      seed,
			"accountID": accountID,
			"flag":      flag,
			"lock":      lock,
		})
	}

	return nil
}

// handleLoginPinNewCodeResult handles the login_pin_new_code_result packet
func (m *LoginManager) handleLoginPinNewCodeResult(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var flag int
	var seed uint32

	if flagVal, ok := args["flag"].(uint16); ok {
		flag = int(flagVal)
	}

	if seedVal, ok := args["seed"].(uint32); ok {
		seed = seedVal
	}

	// Call hook
	if m.hookManager != nil {
		m.hookManager.CallHook("security/pin_new_code_result", map[string]interface{}{
			"flag": flag,
			"seed": seed,
		})
	}

	return nil
}

// Helper functions

// parseServerInfo parses the server info from the account_server_info packet
func (m *LoginManager) parseServerInfo(serverInfo []byte) {
	// Clear existing servers
	m.servers = nil

	// Determine server info format based on length
	serverLen := 32 // Default length for old format

	// Check if it's a newer format
	if len(serverInfo) >= 160 && (len(serverInfo)%160 == 0) {
		serverLen = 160 // kRO Zero 2017, kRO ST 201703+, vRO 2021
	} else if len(serverInfo) >= 164 && (len(serverInfo)%164 == 0) {
		serverLen = 164 // tRO 2020, twRO 2021
	} else if len(serverInfo) >= 154 && (len(serverInfo)%154 == 0) {
		serverLen = 154 // cRO 2017
	} else if len(serverInfo) >= 36 && (len(serverInfo)%36 == 0) {
		serverLen = 36 // tRO 2020 and aRO 2022
	} else if len(serverInfo) >= 32 && (len(serverInfo)%32 == 0) {
		serverLen = 32 // Default format
	}

	// Parse servers
	for i := 0; i+serverLen <= len(serverInfo); i += serverLen {
		server := ServerInfo{}

		// Parse IP (first 4 bytes)
		if serverLen >= 4 {
			ip := serverInfo[i : i+4]
			server.IP = fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
		}

		// Parse port (next 2 bytes)
		if serverLen >= 6 {
			server.Port = int(uint16(serverInfo[i+4]) | uint16(serverInfo[i+5])<<8)
		}

		// Parse name (next 20 bytes)
		if serverLen >= 26 {
			nameEnd := i + 26
			for j := i + 6; j < nameEnd; j++ {
				if serverInfo[j] == 0 {
					nameEnd = j
					break
				}
			}
			server.Name = string(serverInfo[i+6 : nameEnd])
		}

		// Parse users (next 2 bytes)
		if serverLen >= 28 {
			server.Users = int(uint16(serverInfo[i+26]) | uint16(serverInfo[i+27])<<8)
		}

		// Parse state (next 2 bytes)
		if serverLen >= 30 {
			server.State = int(uint16(serverInfo[i+28]) | uint16(serverInfo[i+29])<<8)
		}

		// Parse property (next 2 bytes)
		if serverLen >= 32 {
			server.Property = int(uint16(serverInfo[i+30]) | uint16(serverInfo[i+31])<<8)
		}

		// Set flags based on property
		server.IsNew = (server.Property & 0x01) == 0x01
		server.IsPvP = (server.Property & 0x02) == 0x02
		server.IsLocked = (server.Property & 0x04) == 0x04
		server.IsPremium = (server.Property & 0x08) == 0x08

		// Add server to list
		m.servers = append(m.servers, server)
	}
}

// getLoginErrorMessage returns the error message for a login error code
func (m *LoginManager) getLoginErrorMessage(errorCode int) string {
	// Constants from the original Perl implementation
	const (
		REFUSE_INVALID_ID                = 0x0
		REFUSE_INVALID_PASSWD            = 0x1
		ACCEPT_ID_PASSWD                 = 0x3
		REFUSE_NOT_CONFIRMED             = 0x4
		REFUSE_INVALID_VERSION           = 0x5
		REFUSE_BLOCK_TEMPORARY           = 0x6
		REFUSE_BAN_BY_GM                 = 0xb
		REFUSE_USER_PHONE_BLOCK          = 0x69
		ACCEPT_LOGIN_USER_PHONE_BLOCK    = 0x6a
		REFUSE_EMAIL_NOT_CONFIRMED       = 0xa
		REFUSE_EMAIL_NOT_CONFIRMED2      = 0x1455
		REFUSE_BLOCKED_ID                = 0x1452
		REFUSE_BLOCKED_COUNTRY           = 0x1453
		REFUSE_BILLING                   = 0x1456
		REFUSE_BILLING2                  = 0x1458
		REFUSE_CHANGE_PASSWD_FORCE2      = 0x1459
		REFUSE_ACCOUNT_NOT_PREMIUM       = 0x14B5
		REFUSE_NOT_ALLOWED_IP_ON_TESTING = 0x16
		REFUSE_TOKEN_EXPIRED             = 0xf3
	)

	switch errorCode {
	case REFUSE_INVALID_ID:
		return "Account name doesn't exist"
	case REFUSE_INVALID_PASSWD:
		return "Password Error"
	case ACCEPT_ID_PASSWD:
		return "The server has denied your connection"
	case REFUSE_NOT_CONFIRMED, REFUSE_BAN_BY_GM:
		return "Critical Error: Your account has been blocked"
	case REFUSE_INVALID_VERSION:
		return "Connect failed, something is wrong with the login settings"
	case REFUSE_BLOCK_TEMPORARY:
		return "The server is temporarily blocking your connection"
	case REFUSE_USER_PHONE_BLOCK:
		return "Please dial to activate the login procedure"
	case ACCEPT_LOGIN_USER_PHONE_BLOCK:
		return "Mobile Authentication: Max number of simultaneous IP addresses reached"
	case REFUSE_EMAIL_NOT_CONFIRMED, REFUSE_EMAIL_NOT_CONFIRMED2:
		return "Account email address not confirmed"
	case REFUSE_BLOCKED_ID:
		return "The server is blocking connection from this user"
	case REFUSE_BLOCKED_COUNTRY:
		return "The server is blocking connections from your country"
	case REFUSE_BILLING, REFUSE_BILLING2:
		return "The server is blocking your connection due to billing issues"
	case REFUSE_CHANGE_PASSWD_FORCE2:
		return "The server demands a password change for this account"
	case REFUSE_ACCOUNT_NOT_PREMIUM:
		return "Account doesn't have access to Premium Server"
	case REFUSE_NOT_ALLOWED_IP_ON_TESTING:
		return "Your connection is currently delayed. You can connect again later"
	case REFUSE_TOKEN_EXPIRED:
		return "Your connection was refused due to expired Token"
	// Keep some of the original error messages for backward compatibility
	case 2:
		return "Already logged in."
	case 7:
		return "You are underaged."
	case 8:
		return "Incorrect email address."
	case 9:
		return "Banned until..."
	// Removed case 10 as it's already handled by REFUSE_EMAIL_NOT_CONFIRMED (0xa)
	case 12, 13, 14, 15:
		return "Your account has been blocked."
	case 99:
		return "Account not found or password incorrect."
	case 100:
		return "Your account has expired."
	case 101:
		return "Your account has been rejected from the server."
	case 102:
		return "Your account has been permanently banned."
	case 103, 104:
		return "Your account has been blocked."
	default:
		return fmt.Sprintf("The server has denied your connection for unknown reason (%d)", errorCode)
	}
}

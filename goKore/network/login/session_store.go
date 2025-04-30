package login

import (
	"sync"
)

// SessionData contains all session-related information
type SessionData struct {
	// Account server data
	SessionID  []byte
	AccountID  []byte
	SessionID2 []byte
	AccountSex int

	// Character server data
	CharID  []byte
	MapName string
	MapIP   string
	MapPort int
}

// ServerInfo contains information about a game server
type ServerInfo struct {
	Name string
	IP   string
	Port int
}

// SessionStore maintains session data across server transitions
type SessionStore struct {
	mu sync.RWMutex

	// Session data
	sessionID  []byte
	accountID  []byte
	sessionID2 []byte
	accountSex int
	charID     []byte
	mapName    string
	mapIP      string
	mapPort    int

	// Server information
	servers []ServerInfo
}

// NewSessionStore creates a new session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		servers: make([]ServerInfo, 0),
	}
}

// GetSessionData returns the current session data
func (ss *SessionStore) GetSessionData() SessionData {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	// Create deep copies of byte slices to prevent modification
	var sessionID, accountID, sessionID2, charID []byte

	if ss.sessionID != nil {
		sessionID = make([]byte, len(ss.sessionID))
		copy(sessionID, ss.sessionID)
	}

	if ss.accountID != nil {
		accountID = make([]byte, len(ss.accountID))
		copy(accountID, ss.accountID)
	}

	if ss.sessionID2 != nil {
		sessionID2 = make([]byte, len(ss.sessionID2))
		copy(sessionID2, ss.sessionID2)
	}

	if ss.charID != nil {
		charID = make([]byte, len(ss.charID))
		copy(charID, ss.charID)
	}

	return SessionData{
		SessionID:  sessionID,
		AccountID:  accountID,
		SessionID2: sessionID2,
		AccountSex: ss.accountSex,
		CharID:     charID,
		MapName:    ss.mapName,
		MapIP:      ss.mapIP,
		MapPort:    ss.mapPort,
	}
}

// UpdateFromAccountServerInfo updates session data from account_server_info packet
func (ss *SessionStore) UpdateFromAccountServerInfo(args map[string]interface{}) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if sessionID, ok := args["sessionID"].([]byte); ok {
		ss.sessionID = make([]byte, len(sessionID))
		copy(ss.sessionID, sessionID)
	}

	if accountID, ok := args["accountID"].([]byte); ok {
		ss.accountID = make([]byte, len(accountID))
		copy(ss.accountID, accountID)
	}

	if sessionID2, ok := args["sessionID2"].([]byte); ok {
		ss.sessionID2 = make([]byte, len(sessionID2))
		copy(ss.sessionID2, sessionID2)
	}

	if accountSex, ok := args["accountSex"].(int); ok {
		ss.accountSex = accountSex
	}
}

// UpdateFromCharacterServerInfo updates session data from character server info
func (ss *SessionStore) UpdateFromCharacterServerInfo(args map[string]interface{}) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if charID, ok := args["charID"].([]byte); ok {
		ss.charID = make([]byte, len(charID))
		copy(ss.charID, charID)
	}

	if mapName, ok := args["mapName"].(string); ok {
		ss.mapName = mapName
	}

	if mapIP, ok := args["mapIP"].(string); ok {
		ss.mapIP = mapIP
	}

	if mapPort, ok := args["mapPort"].(int); ok {
		ss.mapPort = mapPort
	}
}

// SetServerInfo sets the list of available servers
func (ss *SessionStore) SetServerInfo(servers []ServerInfo) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.servers = make([]ServerInfo, len(servers))
	copy(ss.servers, servers)
}

// GetServerInfo returns information about a specific server by name
func (ss *SessionStore) GetServerInfo(serverName string) *ServerInfo {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	for _, server := range ss.servers {
		if server.Name == serverName {
			// Return a copy to prevent modification
			serverCopy := server
			return &serverCopy
		}
	}

	return nil
}

// Reset clears all session data
func (ss *SessionStore) Reset() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.sessionID = nil
	ss.accountID = nil
	ss.sessionID2 = nil
	ss.accountSex = 0
	ss.charID = nil
	ss.mapName = ""
	ss.mapIP = ""
	ss.mapPort = 0
	ss.servers = make([]ServerInfo, 0)
}

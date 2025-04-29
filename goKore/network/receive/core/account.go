// Package core provides core functionality for parsing and processing network packets.
package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lenaxia/goKore/network"
)

// Errors
var (
	ErrInvalidAccountID = errors.New("invalid account ID")
	ErrInvalidCharID    = errors.New("invalid character ID")
	ErrSessionExpired   = errors.New("session expired")
)

// AccountState represents the state of an account
type AccountState int

// Account states
const (
	AccountStateUnknown AccountState = iota
	AccountStateLoggedOut
	AccountStateLoggingIn
	AccountStateLoggedIn
	AccountStateSelectingChar
	AccountStateInGame
)

// String returns the string representation of the account state
func (s AccountState) String() string {
	switch s {
	case AccountStateUnknown:
		return "Unknown"
	case AccountStateLoggedOut:
		return "LoggedOut"
	case AccountStateLoggingIn:
		return "LoggingIn"
	case AccountStateLoggedIn:
		return "LoggedIn"
	case AccountStateSelectingChar:
		return "SelectingChar"
	case AccountStateInGame:
		return "InGame"
	default:
		return "Invalid"
	}
}

// SessionData represents the session data for an account
type SessionData struct {
	AccountID         uint32
	AccountName       string
	CharID            uint32
	CharName          string
	SessionID         []byte
	SessionID2        []byte
	MapName           string
	Sex               byte
	ServerType        string
	LastLoginTime     time.Time
	LastPacketTime    time.Time
	CreationTime      time.Time
	ExpirationTime    time.Time
	State             AccountState
	CharacterSlots    int
	PremiumStartSlot  int
	PremiumEndSlot    int
	Characters        []CharacterInfo
	SelectedCharIndex int
}

// CharacterInfo represents information about a character
type CharacterInfo struct {
	CharID     uint32
	Name       string
	JobID      uint16
	Level      uint16
	Sex        byte
	MapName    string
	Coords     [3]byte
	Slot       int
	DeleteDate time.Time
	IsDeleting bool
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

// AccountManager manages account-related functionality
type AccountManager struct {
	parser           *CoreParser
	session          *SessionData
	mutex            sync.RWMutex
	networkState     int
	syncExReplyTable map[string]string // Maps sync request packet IDs to their reply IDs
	xkore            string            // XKore setting (0, 1, or 3)
}

// NewAccountManager creates a new account manager
func NewAccountManager(parser *CoreParser) *AccountManager {
	now := time.Now()
	return &AccountManager{
		parser: parser,
		session: &SessionData{
			State:          AccountStateLoggedOut,
			CreationTime:   now,
			LastPacketTime: now,
		},
		networkState:     network.NotConnected,
		syncExReplyTable: make(map[string]string),
		xkore:            "0", // Default to 0 (not XKore)
	}
}

// RegisterHandlers registers account-related packet handlers
func (m *AccountManager) RegisterHandlers() {
	// Register handlers for account-related packets
	m.parser.RegisterHandlerFunc("0069", "account_server_info", "v a4 a4 a4 a4 a26 C a*",
		[]string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		m.handleAccountServerInfo)

	m.parser.RegisterHandlerFunc("006B", "received_characters_info", "v C3 x20 a*",
		[]string{"len", "total_slot", "premium_start_slot", "premium_end_slot", "charInfo"},
		m.handleReceivedCharactersInfo)

	m.parser.RegisterHandlerFunc("006C", "login_error_game_login_server", "",
		[]string{},
		m.handleLoginError)

	m.parser.RegisterHandlerFunc("006D", "character_creation_successful", "a*",
		[]string{"charInfo"},
		m.handleCharacterCreationSuccessful)

	m.parser.RegisterHandlerFunc("0071", "received_character_ID_and_Map", "a4 Z16 a4 v",
		[]string{"charID", "mapName", "mapIP", "mapPort"},
		m.handleReceivedCharacterIDAndMap)

	m.parser.RegisterHandlerFunc("0073", "map_loaded", "V a3 C2",
		[]string{"syncMapSync", "coords", "xSize", "ySize"},
		m.handleMapLoaded)

	// Register new handlers for basic flow
	m.parser.RegisterHandlerFunc("0074", "map_load_error", "C",
		[]string{"error"},
		m.handleMapLoadError)

	m.parser.RegisterHandlerFunc("007F", "connection_refused", "C",
		[]string{"error"},
		m.handleConnectionRefused)

	m.parser.RegisterHandlerFunc("007E", "map_change", "Z16 v2",
		[]string{"map", "x", "y"},
		m.handleMapChange)

	m.parser.RegisterHandlerFunc("0091", "map_changed", "Z16 v2 V2",
		[]string{"map", "x", "y", "ip", "port"},
		m.handleMapChanged)

	m.parser.RegisterHandlerFunc("00AE", "received_sync", "V",
		[]string{"time"},
		m.handleReceivedSync)

	m.parser.RegisterHandlerFunc("0088", "actor_movement_interrupted", "a4 v2",
		[]string{"ID", "posX", "posY"},
		m.handleActorMovementInterrupted)

	m.parser.RegisterHandlerFunc("0AB8", "move_interrupt", "",
		[]string{},
		m.handleMoveInterrupt)

	m.parser.RegisterHandlerFunc("07E5", "sync_request_ex", "",
		[]string{"switch"},
		m.handleSyncRequestEx)

	m.parser.RegisterHandlerFunc("0B1D", "ping", "",
		[]string{},
		m.handlePing)

	m.parser.RegisterHandlerFunc("018B", "quit_response", "C",
		[]string{"fail"},
		m.handleQuitResponse)

	m.parser.RegisterHandlerFunc("00B3", "switch_character", "C",
		[]string{"result"},
		m.handleSwitchCharacter)

	m.parser.RegisterHandlerFunc("006F", "character_deletion_successful", "",
		[]string{},
		m.handleCharacterDeletionSuccessful)

	m.parser.RegisterHandlerFunc("0070", "character_deletion_failed", "",
		[]string{},
		m.handleCharacterDeletionFailed)

	m.parser.RegisterHandlerFunc("0828", "char_delete2_result", "C L",
		[]string{"result", "deleteDate"},
		m.handleCharDelete2Result)

	m.parser.RegisterHandlerFunc("082A", "char_delete2_accept_result", "a4 C",
		[]string{"charID", "result"},
		m.handleCharDelete2AcceptResult)

	m.parser.RegisterHandlerFunc("082C", "char_delete2_cancel_result", "C",
		[]string{"result"},
		m.handleCharDelete2CancelResult)

	// Register map_change_cell handler (0192)
	m.parser.RegisterHandlerFunc("0192", "map_change_cell", "v2 b Z16",
		[]string{"x", "y", "type", "map_name"},
		m.handleMapChangeCell)

	// Register map_property3 handler (099B)
	m.parser.RegisterHandlerFunc("099B", "map_property3", "C a*",
		[]string{"type", "info_table"},
		m.handleMapProperty3)

	// Register no_teleport handler (0189)
	m.parser.RegisterHandlerFunc("0189", "no_teleport", "C",
		[]string{"fail"},
		m.handleNoTeleport)

	// Register warp_portal_list handler (0B1B)
	m.parser.RegisterHandlerFunc("0B1B", "warp_portal_list", "C Z16 Z16 Z16 Z16",
		[]string{"type", "memo1", "memo2", "memo3", "memo4"},
		m.handleWarpPortalList)
}

// GetSession returns the current session data
func (m *AccountManager) GetSession() *SessionData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.session
}

// SetAccountID sets the account ID
func (m *AccountManager) SetAccountID(accountID uint32) error {
	if accountID == 0 {
		return ErrInvalidAccountID
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.AccountID = accountID
	m.session.LastPacketTime = time.Now()

	return nil
}

// SetCharID sets the character ID
func (m *AccountManager) SetCharID(charID uint32) error {
	if charID == 0 {
		return ErrInvalidCharID
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.CharID = charID
	m.session.LastPacketTime = time.Now()

	return nil
}

// SetState sets the account state
func (m *AccountManager) SetState(state AccountState) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.State = state
	m.session.LastPacketTime = time.Now()
}

// GetState returns the account state
func (m *AccountManager) GetState() AccountState {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.session.State
}

// SetNetworkState sets the network state
func (m *AccountManager) SetNetworkState(state int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.networkState = state

	// Update account state based on network state
	switch state {
	case network.NotConnected:
		m.session.State = AccountStateLoggedOut
	case network.ConnectedToMasterServer:
		m.session.State = AccountStateLoggingIn
	case network.ConnectedToLoginServer:
		m.session.State = AccountStateLoggingIn
	case network.ConnectedToCharServer:
		m.session.State = AccountStateLoggedIn
	case network.InGame:
		m.session.State = AccountStateInGame
	}

	m.session.LastPacketTime = time.Now()
}

// GetNetworkState returns the network state
func (m *AccountManager) GetNetworkState() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.networkState
}

// IsLoggedIn returns whether the account is logged in
func (m *AccountManager) IsLoggedIn() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.session.State >= AccountStateLoggedIn
}

// IsInGame returns whether the account is in game
func (m *AccountManager) IsInGame() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.session.State == AccountStateInGame
}

// ResetSession resets the session data
func (m *AccountManager) ResetSession() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	m.session = &SessionData{
		State:          AccountStateLoggedOut,
		CreationTime:   now,
		LastPacketTime: now,
	}
	m.networkState = network.NotConnected
}

// SetXKore sets the XKore setting
func (m *AccountManager) SetXKore(value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.xkore = value
}

// GetXKore gets the XKore setting
func (m *AccountManager) GetXKore() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.xkore
}

// handlePing handles the ping packet
func (m *AccountManager) handlePing(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If XKore is 1 or 3, do nothing
	if m.xkore == "1" || m.xkore == "3" {
		return nil
	}

	// TODO: Implement sendPing in the messageSender
	// m.messageSender.sendPing()

	return nil
}

// UpdateLastPacketTime updates the last packet time
func (m *AccountManager) UpdateLastPacketTime() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.LastPacketTime = time.Now()
}

// IsSessionExpired checks if the session has expired
func (m *AccountManager) IsSessionExpired(timeout time.Duration) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return time.Since(m.session.LastPacketTime) > timeout
}

// Packet handlers

// handleAccountServerInfo handles the account_server_info packet
func (m *AccountManager) handleAccountServerInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract session IDs and account ID
	if sessionID, ok := args["sessionID"].([]byte); ok {
		m.session.SessionID = sessionID
	}

	if sessionID2, ok := args["sessionID2"].([]byte); ok {
		m.session.SessionID2 = sessionID2
	}

	if accountID, ok := args["accountID"].([]byte); ok && len(accountID) >= 4 {
		m.session.AccountID = uint32(accountID[0]) | uint32(accountID[1])<<8 | uint32(accountID[2])<<16 | uint32(accountID[3])<<24
	}

	if accountSex, ok := args["accountSex"].(byte); ok {
		m.session.Sex = accountSex
	}

	// Update state
	m.session.State = AccountStateLoggedIn
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleReceivedCharactersInfo handles the received_characters_info packet
func (m *AccountManager) handleReceivedCharactersInfo(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character slots information
	if totalSlot, ok := args["total_slot"].(byte); ok {
		m.session.CharacterSlots = int(totalSlot)
	}

	if premiumStartSlot, ok := args["premium_start_slot"].(byte); ok {
		m.session.PremiumStartSlot = int(premiumStartSlot)
	}

	if premiumEndSlot, ok := args["premium_end_slot"].(byte); ok {
		m.session.PremiumEndSlot = int(premiumEndSlot)
	}

	// Update state
	m.session.State = AccountStateSelectingChar
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleLoginError handles the login_error_game_login_server packet
func (m *AccountManager) handleLoginError(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Reset session on login error
	m.session.State = AccountStateLoggedOut
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleCharacterCreationSuccessful handles the character_creation_successful packet
func (m *AccountManager) handleCharacterCreationSuccessful(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update state
	m.session.State = AccountStateSelectingChar
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleReceivedCharacterIDAndMap handles the received_character_ID_and_Map packet
func (m *AccountManager) handleReceivedCharacterIDAndMap(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Extract character ID and map name
	if charID, ok := args["charID"].([]byte); ok && len(charID) >= 4 {
		m.session.CharID = uint32(charID[0]) | uint32(charID[1])<<8 | uint32(charID[2])<<16 | uint32(charID[3])<<24
	}

	if mapName, ok := args["mapName"].(string); ok {
		m.session.MapName = mapName
	}

	// Update state
	m.session.State = AccountStateInGame
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleMapLoaded handles the map_loaded packet
func (m *AccountManager) handleMapLoaded(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update state
	m.session.State = AccountStateInGame
	m.session.LastPacketTime = time.Now()

	return nil
}

// handleQuitResponse handles the quit_response packet
// This packet is sent by the server in response to a disconnect request
// fail:
//
//	0 = disconnect (quit)
//	1 = cannot disconnect (wait 10 seconds)
func (m *AccountManager) handleQuitResponse(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Check if the disconnect request failed
	if fail, ok := args["fail"].(byte); ok && fail == 1 {
		// Cannot disconnect, need to wait 10 seconds
		// In the original Perl code, this would log an error message
		// TODO: Add proper logging when logger is implemented
		return nil
	}

	// Disconnect successful
	// In the original Perl code, this would log a success message
	// TODO: Add proper logging when logger is implemented

	// Update state to logged out
	m.session.State = AccountStateLoggedOut
	m.networkState = network.NotConnected

	return nil
}

// handleSwitchCharacter handles the switch_character packet
// This packet is sent by the server when a client requests to switch characters
// result:
//
//	1 = disconnect, char-select
//	? = nothing
func (m *AccountManager) handleSwitchCharacter(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Check if this is a valid switch character request
	if result, ok := args["result"].(byte); ok && result == 1 {
		// Set network state to connected to master server
		m.networkState = network.ConnectedToMasterServer
		m.session.State = AccountStateLoggingIn

		// Clear character list
		m.session.Characters = nil
		m.session.SelectedCharIndex = -1

		// In the original Perl code, this would also disconnect from the server
		// TODO: Implement server disconnect when network layer is available

		// Log the result
		// TODO: Add proper logging when logger is implemented
	}

	return nil
}

// handleCharacterDeletionSuccessful handles the character_deletion_successful packet
// This packet is sent by the server when a character deletion is successful
func (m *AccountManager) handleCharacterDeletionSuccessful(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Check if we have a selected character index for deletion
	// In the original Perl code, this would be stored in $AI::temp::delIndex
	// For now, we'll use the SelectedCharIndex field in the session
	if m.session.SelectedCharIndex >= 0 && m.session.SelectedCharIndex < len(m.session.Characters) {
		// Log the deletion
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// message TF("Character %s (%d) deleted.\n", m.session.Characters[m.session.SelectedCharIndex].Name, m.session.SelectedCharIndex), "info"

		// Remove the character from the list
		if m.session.SelectedCharIndex < len(m.session.Characters) {
			// Remove the character at the selected index
			m.session.Characters = append(
				m.session.Characters[:m.session.SelectedCharIndex],
				m.session.Characters[m.session.SelectedCharIndex+1:]...,
			)
		}

		// Reset the selected index
		m.session.SelectedCharIndex = -1
	} else {
		// Log a generic message if we don't have a selected index
		// TODO: Add proper logging when logger is implemented
		// message T("Character deleted.\n"), "info"
	}

	// Update state to character selection
	m.session.State = AccountStateSelectingChar

	return nil
}

// handleCharacterDeletionFailed handles the character_deletion_failed packet
// This packet is sent by the server when a character deletion fails
func (m *AccountManager) handleCharacterDeletionFailed(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Reset the selected index
	m.session.SelectedCharIndex = -1

	// Log the error
	// TODO: Add proper logging when logger is implemented
	// In the original Perl code, this would log:
	// error T("Character cannot be deleted. Your e-mail address was probably wrong.\n")

	// Update state to character selection
	m.session.State = AccountStateSelectingChar

	return nil
}

// handleCharDelete2Result handles the char_delete2_result packet
// This packet is sent by the server when a character deletion is requested
func (m *AccountManager) handleCharDelete2Result(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract result and deleteDate from args
	var result byte
	var deleteDate uint32

	if resultVal, ok := args["result"].(byte); ok {
		result = resultVal
	}

	if deleteDateVal, ok := args["deleteDate"].(uint32); ok {
		deleteDate = deleteDateVal
	}

	// Process the result
	if result > 0 && deleteDate > 0 {
		// Successful deletion request, set delete date for the character
		if m.session.SelectedCharIndex >= 0 && m.session.SelectedCharIndex < len(m.session.Characters) {
			// Set the delete date for the character
			m.session.Characters[m.session.SelectedCharIndex].DeleteDate = time.Unix(int64(deleteDate), 0)
			m.session.Characters[m.session.SelectedCharIndex].IsDeleting = true

			// Log the deletion date
			// TODO: Add proper logging when logger is implemented
			// In the original Perl code, this would log:
			// message TF("Your character will be delete, left %s\n", $chars[$messageSender->{char_delete_slot}]{deleteDate}), "connection"
		}
	} else if result == 0 {
		// Character already planned to be erased
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("That character already planned to be erased!\n")
	} else if result == 3 {
		// Error in database of the server
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Error in database of the server!\n")
	} else if result == 4 {
		// Need to withdraw from guild
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("To delete a character you must withdraw from the guild!\n")
	} else if result == 5 {
		// Need to withdraw from party
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("To delete a character you must withdraw from the party!\n")
	} else {
		// Unknown error
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error TF("Unknown error when trying to delete the character! (Error number: %s)\n", $result)
	}

	// Update state to character selection
	m.session.State = AccountStateSelectingChar

	return nil
}

// handleCharDelete2AcceptResult handles the char_delete2_accept_result packet
// This packet is sent by the server when a character deletion is accepted
func (m *AccountManager) handleCharDelete2AcceptResult(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract charID and result from args
	var charID []byte
	var result byte

	if charIDVal, ok := args["charID"].([]byte); ok {
		charID = charIDVal
	}

	if resultVal, ok := args["result"].(byte); ok {
		result = resultVal
	}

	// Process the result
	if result == 1 { // Success
		// Find the character in the list by charID
		charIndex := -1
		if len(charID) >= 4 {
			charIDValue := uint32(charID[0]) | uint32(charID[1])<<8 | uint32(charID[2])<<16 | uint32(charID[3])<<24
			for i, char := range m.session.Characters {
				if char.CharID == charIDValue {
					charIndex = i
					break
				}
			}
		}

		if charIndex >= 0 && charIndex < len(m.session.Characters) {
			// Log the deletion
			// TODO: Add proper logging when logger is implemented
			// In the original Perl code, this would log:
			// message TF("Character %s (%d) deleted.\n", $chars[$AI::temp::delIndex]{name}, $AI::temp::delIndex), "info"

			// Remove the character from the list
			m.session.Characters = append(
				m.session.Characters[:charIndex],
				m.session.Characters[charIndex+1:]...,
			)
		} else {
			// Log a generic message if we don't have a matching character
			// TODO: Add proper logging when logger is implemented
			// In the original Perl code, this would log:
			// message T("Character deleted.\n"), "info"
		}
	} else if result == 0 {
		// Need to enter birthday
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Enter your 6-digit birthday (YYMMDD) (e.g: 801122).\n")
	} else if result == 2 {
		// Cannot delete due to system settings
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Due to system settings, can not be deleted.\n")
	} else if result == 3 {
		// Database error
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("A database error has occurred.\n")
	} else if result == 4 {
		// Cannot delete at the moment
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("You cannot delete this character at the moment.\n")
	} else if result == 5 {
		// Birthday does not match
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Your entered birthday does not match.\n")
	} else if result == 7 {
		// Incorrect email address
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Character Deletion has failed because you have entered an incorrect e-mail address.\n")
	} else {
		// Unknown error
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error TF("An unknown error has occurred. Error number %d\n", $result)
	}

	// Reset the selected index
	m.session.SelectedCharIndex = -1

	// Update state to character selection
	m.session.State = AccountStateSelectingChar

	return nil
}

// handleCharDelete2CancelResult handles the char_delete2_cancel_result packet
// This packet is sent by the server when a character deletion is canceled
func (m *AccountManager) handleCharDelete2CancelResult(args map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update last packet time
	m.session.LastPacketTime = time.Now()

	// Extract result from args
	var result byte

	if resultVal, ok := args["result"].(byte); ok {
		result = resultVal
	}

	// Process the result
	if result == 1 { // Success
		// Find the character that was scheduled for deletion
		if m.session.SelectedCharIndex >= 0 && m.session.SelectedCharIndex < len(m.session.Characters) {
			// Clear the delete date for the character
			m.session.Characters[m.session.SelectedCharIndex].DeleteDate = time.Time{}
			m.session.Characters[m.session.SelectedCharIndex].IsDeleting = false

			// Log the cancellation
			// TODO: Add proper logging when logger is implemented
			// In the original Perl code, this would log:
			// message T("Character is no longer scheduled to be deleted\n"), "connection"
		}
	} else if result == 2 {
		// Database error
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error T("Error in database of the server!\n")
	} else {
		// Unknown error
		// TODO: Add proper logging when logger is implemented
		// In the original Perl code, this would log:
		// error TF("Unknown error when trying to cancel the deletion of the character! (Error number: %s)\n", $result)
	}

	// Reset the selected index
	m.session.SelectedCharIndex = -1

	// Update state to character selection
	m.session.State = AccountStateSelectingChar

	return nil
}

// ChangeToInGameState changes the account state to in-game if conditions are met
// Returns true if the state was changed to in-game, false otherwise
func (m *AccountManager) ChangeToInGameState() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we have a valid account ID and character ID
	if m.session.AccountID != 0 && m.session.CharID != 0 {
		// If we're not already in game, change to in-game state
		if m.networkState != network.InGame {
			m.networkState = network.InGame
			m.session.State = AccountStateInGame
			m.session.LastPacketTime = time.Now()
		}
		return true
	} else {
		// If we're not initialized yet, set to a transitional state
		if m.networkState != network.ConnectedToCharServer {
			m.networkState = network.ConnectedToCharServer
			m.session.State = AccountStateLoggedIn
			m.session.LastPacketTime = time.Now()
		}
		return false
	}
}

// ReceivedCharactersBlockSize returns the block size for the received_characters packet
// This may be overridden in server-specific implementations
func (m *AccountManager) ReceivedCharactersBlockSize() int {
	// Default block size
	return 106
}

// ReceivedCharactersUnpackString returns the unpack string for the received_characters packet
// This may be overridden in server-specific implementations
func (m *AccountManager) ReceivedCharactersUnpackString() string {
	// Default unpack string
	return "a4 Z24 Z16 v3 V v11 a*"
}

// ReceivedCharactersSlotsInfo parses the character slots information from the received_characters_info packet
// Returns a map containing the slot information
func (m *AccountManager) ReceivedCharactersSlotsInfo(args map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Extract character slots information
	if totalSlot, ok := args["total_slot"].(byte); ok {
		result["total_slot"] = int(totalSlot)
	}

	if premiumStartSlot, ok := args["premium_start_slot"].(byte); ok {
		result["premium_start_slot"] = int(premiumStartSlot)
	}

	if premiumEndSlot, ok := args["premium_end_slot"].(byte); ok {
		result["premium_end_slot"] = int(premiumEndSlot)
	}

	return result
}

// ReceivedCharacters parses the character information from the received_characters packet
// Returns a slice of CharacterInfo structures
func (m *AccountManager) ReceivedCharacters(charInfo []byte) []CharacterInfo {
	characters := []CharacterInfo{}

	// Get the block size and unpack string
	blockSize := m.ReceivedCharactersBlockSize()

	// Parse each character block
	for i := 0; i+blockSize <= len(charInfo); i += blockSize {
		// Extract character data
		charBlock := charInfo[i : i+blockSize]

		// Create a new character info
		char := CharacterInfo{
			Slot: i / blockSize,
		}

		// Extract character ID (first 4 bytes)
		if len(charBlock) >= 4 {
			char.CharID = uint32(charBlock[0]) | uint32(charBlock[1])<<8 | uint32(charBlock[2])<<16 | uint32(charBlock[3])<<24
		}

		// Extract character name (next 24 bytes, null-terminated)
		if len(charBlock) >= 28 {
			nameEnd := 4
			for nameEnd < 28 && charBlock[nameEnd] != 0 {
				nameEnd++
			}
			char.Name = string(charBlock[4:nameEnd])
		}

		// Extract map name (next 16 bytes, null-terminated)
		if len(charBlock) >= 44 {
			mapEnd := 28
			for mapEnd < 44 && charBlock[mapEnd] != 0 {
				mapEnd++
			}
			char.MapName = string(charBlock[28:mapEnd])
		}

		// Extract job ID (next 2 bytes)
		if len(charBlock) >= 46 {
			char.JobID = uint16(charBlock[44]) | uint16(charBlock[45])<<8
		}

		// Extract level (next 2 bytes)
		if len(charBlock) >= 48 {
			char.Level = uint16(charBlock[46]) | uint16(charBlock[47])<<8
		}

		// Extract coordinates (next 3 bytes)
		if len(charBlock) >= 51 {
			char.Coords[0] = charBlock[48]
			char.Coords[1] = charBlock[49]
			char.Coords[2] = charBlock[50]
		}

		// Extract sex (next byte)
		if len(charBlock) >= 52 {
			char.Sex = charBlock[51]
		}

		// Add to the list
		characters = append(characters, char)
	}

	return characters
}

// SyncReceivedCharacters synchronizes the received characters with the session
func (m *AccountManager) SyncReceivedCharacters(characters []CharacterInfo) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.Characters = characters
	m.session.LastPacketTime = time.Now()
}

// ReconstructReceivedCharacters reconstructs the received_characters packet
// This is used for packet manipulation and testing
func (m *AccountManager) ReconstructReceivedCharacters(characters []CharacterInfo) []byte {
	blockSize := m.ReceivedCharactersBlockSize()
	result := make([]byte, len(characters)*blockSize)

	for i, char := range characters {
		offset := i * blockSize

		// Character ID
		result[offset] = byte(char.CharID)
		result[offset+1] = byte(char.CharID >> 8)
		result[offset+2] = byte(char.CharID >> 16)
		result[offset+3] = byte(char.CharID >> 24)

		// Character name
		copy(result[offset+4:offset+28], []byte(char.Name))

		// Map name
		copy(result[offset+28:offset+44], []byte(char.MapName))

		// Job ID
		result[offset+44] = byte(char.JobID)
		result[offset+45] = byte(char.JobID >> 8)

		// Level
		result[offset+46] = byte(char.Level)
		result[offset+47] = byte(char.Level >> 8)

		// Coordinates
		result[offset+48] = char.Coords[0]
		result[offset+49] = char.Coords[1]
		result[offset+50] = char.Coords[2]

		// Sex
		result[offset+51] = char.Sex
	}

	return result
}

// ReconstructReceivedCharactersInfo reconstructs the received_characters_info packet
// This is used for packet manipulation and testing
func (m *AccountManager) ReconstructReceivedCharactersInfo(totalSlot, premiumStartSlot, premiumEndSlot byte) []byte {
	// Create a basic packet with the slot information
	result := []byte{
		totalSlot,
		premiumStartSlot,
		premiumEndSlot,
	}

	// Add padding
	padding := make([]byte, 20)
	result = append(result, padding...)

	return result
}

// CharacterCreationSuccessful handles the character_creation_successful packet
// Returns the created character info
func (m *AccountManager) CharacterCreationSuccessful(charInfo []byte) CharacterInfo {
	// Parse the character info using the same format as ReceivedCharacters
	characters := m.ReceivedCharacters(charInfo)

	// If we successfully parsed a character, return it
	if len(characters) > 0 {
		char := characters[0]

		// Update the session
		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Add the character to the list
		m.session.Characters = append(m.session.Characters, char)

		// Update state
		m.session.State = AccountStateSelectingChar
		m.session.LastPacketTime = time.Now()

		return char
	}

	// Return an empty character if parsing failed
	return CharacterInfo{}
}

// CharacterCreationFailed handles the character_creation_failed packet
// Returns the error code and message
func (m *AccountManager) CharacterCreationFailed(errorCode byte) (int, string) {
	// Map of error codes to messages
	errorMessages := map[byte]string{
		0x00: "Charname already exists",
		0x01: "You are underaged",
		0x02: "Symbols in character names are forbidden",
		0x03: "You are not eligible to open the Character Slot",
		0xFF: "Unknown error",
	}

	// Get the error message
	message, ok := errorMessages[errorCode]
	if !ok {
		message = "Unknown error"
	}

	// Update state
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.session.State = AccountStateSelectingChar
	m.session.LastPacketTime = time.Now()

	return int(errorCode), message
}

// ReceivedCharactersInfo handles the received_characters_info packet
// This is a more comprehensive version of handleReceivedCharactersInfo
func (m *AccountManager) ReceivedCharactersInfo(args map[string]interface{}) map[string]interface{} {
	// Get the slot information
	slotInfo := m.ReceivedCharactersSlotsInfo(args)

	// Extract character info
	var characters []CharacterInfo
	if charInfo, ok := args["charInfo"].([]byte); ok {
		characters = m.ReceivedCharacters(charInfo)
	}

	// Update the session
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update character slots
	if totalSlot, ok := slotInfo["total_slot"].(int); ok {
		m.session.CharacterSlots = totalSlot
	}

	if premiumStartSlot, ok := slotInfo["premium_start_slot"].(int); ok {
		m.session.PremiumStartSlot = premiumStartSlot
	}

	if premiumEndSlot, ok := slotInfo["premium_end_slot"].(int); ok {
		m.session.PremiumEndSlot = premiumEndSlot
	}

	// Update characters
	m.session.Characters = characters

	// Update state
	m.session.State = AccountStateSelectingChar
	m.session.LastPacketTime = time.Now()

	// Return the combined information
	result := slotInfo
	result["characters"] = characters

	return result
}

// ParseAccountServerInfo parses the account_server_info packet
// Returns a map containing the server information
func (m *AccountManager) ParseAccountServerInfo(serverInfo []byte) []map[string]interface{} {
	servers := []map[string]interface{}{}

	// Each server entry is 32 bytes
	blockSize := 32

	for i := 0; i+blockSize <= len(serverInfo); i += blockSize {
		serverBlock := serverInfo[i : i+blockSize]

		server := make(map[string]interface{})

		// IP address (4 bytes)
		if len(serverBlock) >= 4 {
			server["ip"] = fmt.Sprintf("%d.%d.%d.%d", serverBlock[0], serverBlock[1], serverBlock[2], serverBlock[3])
		}

		// Port (2 bytes)
		if len(serverBlock) >= 6 {
			server["port"] = int(serverBlock[4]) | int(serverBlock[5])<<8
		}

		// Server name (20 bytes, null-terminated)
		if len(serverBlock) >= 26 {
			nameEnd := 6
			for nameEnd < 26 && serverBlock[nameEnd] != 0 {
				nameEnd++
			}
			server["name"] = string(serverBlock[6:nameEnd])
		}

		// User count (2 bytes)
		if len(serverBlock) >= 28 {
			server["users"] = int(serverBlock[26]) | int(serverBlock[27])<<8
		}

		// Server state (2 bytes)
		if len(serverBlock) >= 30 {
			state := int(serverBlock[28]) | int(serverBlock[29])<<8
			server["state"] = state

			// Parse state flags
			server["is_new"] = (state & 1) != 0
			server["is_pvp"] = (state & 2) != 0
			server["is_locked"] = (state & 4) != 0
			server["is_premium"] = (state & 8) != 0
		}

		// Server ID (2 bytes)
		if len(serverBlock) >= 32 {
			server["server_id"] = int(serverBlock[30]) | int(serverBlock[31])<<8
		}

		servers = append(servers, server)
	}

	return servers
}

// ReconstructAccountServerInfo reconstructs the account_server_info packet
// This is used for packet manipulation and testing
func (m *AccountManager) ReconstructAccountServerInfo(servers []map[string]interface{}) []byte {
	result := []byte{}

	for _, server := range servers {
		serverBlock := make([]byte, 32)

		// IP address
		if ip, ok := server["ip"].(string); ok {
			parts := strings.Split(ip, ".")
			if len(parts) == 4 {
				for i := 0; i < 4; i++ {
					val, _ := strconv.Atoi(parts[i])
					serverBlock[i] = byte(val)
				}
			}
		}

		// Port
		if port, ok := server["port"].(int); ok {
			serverBlock[4] = byte(port)
			serverBlock[5] = byte(port >> 8)
		}

		// Server name
		if name, ok := server["name"].(string); ok {
			copy(serverBlock[6:26], []byte(name))
		}

		// User count
		if users, ok := server["users"].(int); ok {
			serverBlock[26] = byte(users)
			serverBlock[27] = byte(users >> 8)
		}

		// Server state
		state := 0
		if isNew, ok := server["is_new"].(bool); ok && isNew {
			state |= 1
		}
		if isPvP, ok := server["is_pvp"].(bool); ok && isPvP {
			state |= 2
		}
		if isLocked, ok := server["is_locked"].(bool); ok && isLocked {
			state |= 4
		}
		if isPremium, ok := server["is_premium"].(bool); ok && isPremium {
			state |= 8
		}

		// Override with explicit state if provided
		if explicitState, ok := server["state"].(int); ok {
			state = explicitState
		}

		serverBlock[28] = byte(state)
		serverBlock[29] = byte(state >> 8)

		// Server ID
		if serverID, ok := server["server_id"].(int); ok {
			serverBlock[30] = byte(serverID)
			serverBlock[31] = byte(serverID >> 8)
		}

		result = append(result, serverBlock...)
	}

	return result
}

// AccountServerInfo processes the account_server_info packet
// This combines parsing and updating the session
func (m *AccountManager) AccountServerInfo(args map[string]interface{}) []map[string]interface{} {
	// Extract session IDs and account ID
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if sessionID, ok := args["sessionID"].([]byte); ok {
		m.session.SessionID = sessionID
	}

	if sessionID2, ok := args["sessionID2"].([]byte); ok {
		m.session.SessionID2 = sessionID2
	}

	if accountID, ok := args["accountID"].([]byte); ok && len(accountID) >= 4 {
		m.session.AccountID = uint32(accountID[0]) | uint32(accountID[1])<<8 | uint32(accountID[2])<<16 | uint32(accountID[3])<<24
	}

	if accountSex, ok := args["accountSex"].(byte); ok {
		m.session.Sex = accountSex
	}

	// Parse server info
	var servers []map[string]interface{}
	if serverInfo, ok := args["serverInfo"].([]byte); ok {
		servers = m.ParseAccountServerInfo(serverInfo)
	}

	// Convert to ServerInfo structures
	serverInfos := make([]ServerInfo, len(servers))
	for i, server := range servers {
		serverInfos[i] = ServerInfo{
			IP:        server["ip"].(string),
			Port:      server["port"].(int),
			Name:      server["name"].(string),
			Users:     server["users"].(int),
			State:     server["state"].(int),
			ServerID:  server["server_id"].(int),
			IsNew:     server["is_new"].(bool),
			IsPvP:     server["is_pvp"].(bool),
			IsLocked:  server["is_locked"].(bool),
			IsPremium: server["is_premium"].(bool),
			IPPort:    fmt.Sprintf("%s:%d", server["ip"].(string), server["port"].(int)),
		}
	}

	// Update state
	m.session.State = AccountStateLoggedIn
	m.session.LastPacketTime = time.Now()

	return servers
}

package login

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// NetworkInterface defines the interface for network connections
type NetworkInterface interface {
	Connect() error
	ConnectTo(host string, port int) error
	Disconnect() error
	IsConnected() bool
	Send(data []byte) error
	Receive() ([]byte, error)
}

// PacketSender defines the interface for packet sending
type PacketSender interface {
	Send(packetName string, fields map[string]interface{}) ([]byte, error)
	GetCashShopManager() interface{}
	GetMiscManager() interface{}
	GetInfoChatManager() interface{}
}

// PacketHandler defines the interface for packet handling
type PacketHandler interface {
	Handle(packet []byte) error
}

// HookManager defines the interface for hook management
type HookManager interface {
	Register(hookName string, callback func(string, interface{}, interface{}))
	Unregister(hookName string, callback func(string, interface{}, interface{}))
}

// NetworkManager defines the interface for network management
type NetworkManager interface {
	Connect() error
	ConnectTo(host string, port int) error
	Disconnect() error
	Send(packetName string, fields map[string]interface{}) ([]byte, error)
	HandlePacket(packet []byte) error
	SetState(state int)
	GetState() int
	SetStateChangeCallback(callback func(oldState, newState int))
	GetHookManager() interface{}

	// For testing only
	SetSessionStore(*SessionStore)
}

// LoginManager orchestrates the entire login process
type LoginManager struct {
	networkManager NetworkManager
	stateManager   *LoginStateManager
	sessionStore   *SessionStore
	config         *LoginConfig

	// Retry management
	retryCount int
	retryMutex sync.Mutex

	// Channels for communication
	loginDone  chan struct{}
	loginError chan error
}

// NewLoginManager creates a new login manager
func NewLoginManager(networkManager NetworkManager, config *LoginConfig) *LoginManager {
	lm := &LoginManager{
		networkManager: networkManager,
		config:         config,
		stateManager:   NewLoginStateManager(),
		sessionStore:   NewSessionStore(),
		loginDone:      make(chan struct{}),
		loginError:     make(chan error),
		retryCount:     0,
	}

	// Register as state observer
	lm.stateManager.AddObserver(lm)

	// For testing only - share the session store with the network manager
	networkManager.SetSessionStore(lm.sessionStore)

	return lm
}

// OnStateChange implements the StateObserver interface
func (lm *LoginManager) OnStateChange(oldState, newState LoginState) {
	log.Printf("Login state changed from %v to %v", oldState, newState)

	// Update network manager state to match login state
	switch newState {
	case StateNotConnected:
		lm.networkManager.SetState(0) // NetworkNotConnected
	case StateConnectedToMasterServer:
		lm.networkManager.SetState(1) // NetworkConnectedToMasterServer
	case StateConnectedToCharServer:
		lm.networkManager.SetState(2) // NetworkConnectedToCharServer
	case StateConnectedToMapServer:
		lm.networkManager.SetState(3) // NetworkConnectedToMapServer
	case StateInGame:
		lm.networkManager.SetState(4) // NetworkInGame
	}
}

// Login initiates the login process
func (lm *LoginManager) Login(ctx context.Context) error {
	log.Println("Starting login process")

	// Reset state
	lm.stateManager.SetState(StateNotConnected)
	lm.sessionStore.Reset()
	lm.resetRetryCount()

	// Register event handlers
	lm.registerEventHandlers()

	// Start the login process
	err := lm.connectToMasterServer()
	if err != nil {
		return fmt.Errorf("failed to connect to master server: %w", err)
	}

	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, lm.config.LoginTimeout)
	defer cancel()

	log.Println("Waiting for login completion")

	// Wait for completion, error, timeout, or cancellation
	select {
	case <-lm.loginDone:
		log.Println("Login completed successfully")
		return nil
	case err := <-lm.loginError:
		log.Printf("Login error: %v", err)
		return err
	case <-timeoutCtx.Done():
		if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			log.Println("Login process timed out")
			return errors.New("login process timed out")
		}
		log.Printf("Login process cancelled: %v", ctx.Err())
		return ctx.Err()
	case <-ctx.Done():
		log.Printf("Login process cancelled: %v", ctx.Err())
		return ctx.Err()
	}
}

// registerEventHandlers registers event handlers with the hook system
func (lm *LoginManager) registerEventHandlers() {
	log.Println("Registering event handlers")

	// Get the hook manager
	hookManager, ok := lm.networkManager.GetHookManager().(HookManager)
	if !ok {
		log.Println("Failed to get hook manager")
		return
	}

	// Register for account server info
	hookManager.Register("account_info_received", func(hookName string, arg interface{}, userData interface{}) {
		log.Println("Account info received event")
		if args, ok := arg.(map[string]interface{}); ok {
			lm.handleAccountInfoReceived(args)
		}
	})

	// Register for character info
	hookManager.Register("characters_info_received", func(hookName string, arg interface{}, userData interface{}) {
		log.Println("Characters info received event")
		if args, ok := arg.(map[string]interface{}); ok {
			lm.handleCharactersInfoReceived(args)
		}
	})

	// Register for character map info
	hookManager.Register("character_map_info_received", func(hookName string, arg interface{}, userData interface{}) {
		log.Println("Character map info received event")
		if args, ok := arg.(map[string]interface{}); ok {
			lm.handleCharacterMapInfoReceived(args)
		}
	})

	// Register for map loaded
	hookManager.Register("map_loaded", func(hookName string, arg interface{}, userData interface{}) {
		log.Println("Map loaded event")
		lm.handleMapLoaded()
	})

	// Register for errors
	hookManager.Register("login_error", func(hookName string, arg interface{}, userData interface{}) {
		log.Println("Login error event")
		if args, ok := arg.(map[string]interface{}); ok {
			lm.handleLoginError(args)
		}
	})
}

// connectToMasterServer initiates connection to the master server
func (lm *LoginManager) connectToMasterServer() error {
	log.Println("Connecting to master server")

	// Connect to the server
	err := lm.networkManager.Connect()
	if err != nil {
		log.Printf("Failed to connect to master server: %v", err)
		return err
	}

	// Send login packet
	log.Println("Sending master login packet")
	_, err = lm.networkManager.Send("master_login", map[string]interface{}{
		"version":        lm.config.Version,
		"username":       lm.config.Username,
		"password":       lm.config.Password,
		"master_version": lm.config.MasterVersion,
	})

	if err != nil {
		log.Printf("Failed to send master login packet: %v", err)
	}

	return err
}

// handleAccountInfoReceived processes the account_server_info packet
func (lm *LoginManager) handleAccountInfoReceived(args map[string]interface{}) {
	log.Println("Handling account server info")

	// Update session store
	lm.sessionStore.UpdateFromAccountServerInfo(args)

	// Update state
	lm.stateManager.SetState(StateConnectedToMasterServer)

	// Connect to character server
	go lm.connectToCharServer()
}

// connectToCharServer connects to the character server
func (lm *LoginManager) connectToCharServer() {
	log.Println("Connecting to character server")

	// Get server info from the session
	serverInfo := lm.sessionStore.GetServerInfo(lm.config.ServerName)
	if serverInfo == nil {
		log.Printf("Server not found: %s", lm.config.ServerName)
		lm.loginError <- errors.New("server not found")
		return
	}

	log.Printf("Found server: %s (%s:%d)", serverInfo.Name, serverInfo.IP, serverInfo.Port)

	// Disconnect from master server
	err := lm.networkManager.Disconnect()
	if err != nil {
		log.Printf("Failed to disconnect from master server: %v", err)
		lm.loginError <- fmt.Errorf("failed to disconnect from master server: %w", err)
		return
	}

	// Connect to character server
	err = lm.networkManager.ConnectTo(serverInfo.IP, serverInfo.Port)
	if err != nil {
		log.Printf("Failed to connect to character server: %v", err)
		lm.loginError <- fmt.Errorf("failed to connect to character server: %w", err)
		return
	}

	// Send game login packet
	sessionData := lm.sessionStore.GetSessionData()
	log.Println("Sending game login packet")
	_, err = lm.networkManager.Send("game_login", map[string]interface{}{
		"accountID":  sessionData.AccountID,
		"sessionID":  sessionData.SessionID,
		"sessionID2": sessionData.SessionID2,
		"accountSex": sessionData.AccountSex,
	})

	if err != nil {
		log.Printf("Failed to send game login packet: %v", err)
		lm.loginError <- fmt.Errorf("failed to send game login packet: %w", err)
	}
}

// handleCharactersInfoReceived processes the received_characters_info packet
func (lm *LoginManager) handleCharactersInfoReceived(args map[string]interface{}) {
	log.Println("Handling characters info")

	// Update state
	lm.stateManager.SetState(StateConnectedToCharServer)

	// Send character selection
	log.Println("Sending character selection")
	_, err := lm.networkManager.Send("char_login", map[string]interface{}{
		"slot": lm.config.CharacterID,
	})

	if err != nil {
		log.Printf("Failed to send character selection: %v", err)
		lm.loginError <- fmt.Errorf("failed to send character selection: %w", err)
	}
}

// handleCharacterMapInfoReceived processes the received_character_ID_and_Map packet
func (lm *LoginManager) handleCharacterMapInfoReceived(args map[string]interface{}) {
	log.Println("Handling character map info")

	// Update session store
	lm.sessionStore.UpdateFromCharacterServerInfo(args)

	// Disconnect from character server
	err := lm.networkManager.Disconnect()
	if err != nil {
		log.Printf("Failed to disconnect from character server: %v", err)
		lm.loginError <- fmt.Errorf("failed to disconnect from character server: %w", err)
		return
	}

	// Connect to map server
	sessionData := lm.sessionStore.GetSessionData()
	log.Printf("Connecting to map server: %s:%d", sessionData.MapIP, sessionData.MapPort)
	err = lm.networkManager.ConnectTo(sessionData.MapIP, sessionData.MapPort)
	if err != nil {
		log.Printf("Failed to connect to map server: %v", err)
		lm.loginError <- fmt.Errorf("failed to connect to map server: %w", err)
		return
	}

	// Update state
	lm.stateManager.SetState(StateConnectedToMapServer)

	// Send map login packet
	log.Println("Sending map login packet")
	_, err = lm.networkManager.Send("map_login", map[string]interface{}{
		"accountID": sessionData.AccountID,
		"charID":    sessionData.CharID,
		"sessionID": sessionData.SessionID,
		"sex":       sessionData.AccountSex,
	})

	if err != nil {
		log.Printf("Failed to send map login packet: %v", err)
		lm.loginError <- fmt.Errorf("failed to send map login packet: %w", err)
	}
}

// handleMapLoaded processes the map_loaded packet
func (lm *LoginManager) handleMapLoaded() {
	log.Println("Handling map loaded")

	// Update state
	lm.stateManager.SetState(StateInGame)

	// Send map loaded confirmation
	log.Println("Sending map loaded confirmation")
	_, err := lm.networkManager.Send("map_loaded", nil)
	if err != nil {
		log.Printf("Failed to send map loaded confirmation: %v", err)
		lm.loginError <- fmt.Errorf("failed to send map loaded confirmation: %w", err)
		return
	}

	// Send sync packet
	log.Println("Sending sync packet")
	_, err = lm.networkManager.Send("sync", map[string]interface{}{
		"time": time.Now().Unix(),
	})
	if err != nil {
		log.Printf("Failed to send sync packet: %v", err)
		lm.loginError <- fmt.Errorf("failed to send sync packet: %w", err)
		return
	}

	// Signal successful login
	log.Println("Signaling successful login")
	lm.loginDone <- struct{}{}
}

// handleLoginError processes login errors
func (lm *LoginManager) handleLoginError(args map[string]interface{}) {
	errorType, _ := args["type"].(int)
	errorMsg, _ := args["date"].(string)

	log.Printf("Handling login error: type=%d, message=%s", errorType, errorMsg)

	err := fmt.Errorf("login error (type %d): %s", errorType, errorMsg)

	// Check if we should retry
	retries := lm.getRetryCount()
	if retries < lm.config.MaxRetries {
		lm.incrementRetryCount()

		// Reset state
		lm.stateManager.SetState(StateNotConnected)

		// Wait before retrying
		log.Printf("Retrying login (attempt %d/%d) after %v", retries+1, lm.config.MaxRetries, lm.config.RetryDelay)
		time.Sleep(lm.config.RetryDelay)

		// Retry login
		go func() {
			retryErr := lm.connectToMasterServer()
			if retryErr != nil {
				log.Printf("Retry failed: %v", retryErr)
				lm.loginError <- fmt.Errorf("retry failed: %w", retryErr)
			}
		}()
	} else {
		// Max retries reached, report error
		log.Printf("Max retries reached (%d), reporting error", lm.config.MaxRetries)
		lm.loginError <- err
	}
}

// getRetryCount gets the current retry count
func (lm *LoginManager) getRetryCount() int {
	lm.retryMutex.Lock()
	defer lm.retryMutex.Unlock()
	return lm.retryCount
}

// incrementRetryCount increments the retry count
func (lm *LoginManager) incrementRetryCount() {
	lm.retryMutex.Lock()
	defer lm.retryMutex.Unlock()
	lm.retryCount++
}

// resetRetryCount resets the retry count
func (lm *LoginManager) resetRetryCount() {
	lm.retryMutex.Lock()
	defer lm.retryMutex.Unlock()
	lm.retryCount = 0
}

// GetSessionStore returns the session store (for testing only)
func (lm *LoginManager) GetSessionStore() *SessionStore {
	return lm.sessionStore
}

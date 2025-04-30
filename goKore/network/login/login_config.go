package login

import (
	"errors"
	"fmt"
	"time"
)

// LoginConfig contains configuration for the login process
type LoginConfig struct {
	// Account information
	Username string
	Password string

	// Character selection
	CharacterID int
	ServerName  string

	// Version information
	Version       int
	MasterVersion int

	// Retry settings
	MaxRetries int
	RetryDelay time.Duration

	// Timeout settings
	LoginTimeout time.Duration
}

// Validate checks if the login configuration is valid
func (lc *LoginConfig) Validate() error {
	if lc.Username == "" {
		return errors.New("username is required")
	}

	if lc.Password == "" {
		return errors.New("password is required")
	}

	if lc.ServerName == "" {
		return errors.New("server name is required")
	}

	if lc.CharacterID < 0 {
		return fmt.Errorf("character ID must be non-negative, got %d", lc.CharacterID)
	}

	if lc.RetryDelay <= 0 {
		return fmt.Errorf("retry delay must be positive, got %v", lc.RetryDelay)
	}

	if lc.LoginTimeout <= 0 {
		return fmt.Errorf("login timeout must be positive, got %v", lc.LoginTimeout)
	}

	return nil
}

// ApplyDefaults sets default values for optional fields
func (lc *LoginConfig) ApplyDefaults() {
	// Default character ID is 0 (first character)
	if lc.CharacterID < 0 {
		lc.CharacterID = 0
	}

	// Default version is 1
	if lc.Version <= 0 {
		lc.Version = 1
	}

	// Default master version is 1
	if lc.MasterVersion <= 0 {
		lc.MasterVersion = 1
	}

	// Default max retries is 3
	if lc.MaxRetries <= 0 {
		lc.MaxRetries = 3
	}

	// Default retry delay is 5 seconds
	if lc.RetryDelay <= 0 {
		lc.RetryDelay = 5 * time.Second
	}

	// Default login timeout is 30 seconds
	if lc.LoginTimeout <= 0 {
		lc.LoginTimeout = 30 * time.Second
	}
}

// Clone creates a deep copy of the login configuration
func (lc *LoginConfig) Clone() *LoginConfig {
	return &LoginConfig{
		Username:      lc.Username,
		Password:      lc.Password,
		CharacterID:   lc.CharacterID,
		ServerName:    lc.ServerName,
		Version:       lc.Version,
		MasterVersion: lc.MasterVersion,
		MaxRetries:    lc.MaxRetries,
		RetryDelay:    lc.RetryDelay,
		LoginTimeout:  lc.LoginTimeout,
	}
}

// NewLoginConfig creates a new login configuration with default values
func NewLoginConfig(username, password, serverName string) *LoginConfig {
	config := &LoginConfig{
		Username:   username,
		Password:   password,
		ServerName: serverName,
	}
	config.ApplyDefaults()
	return config
}

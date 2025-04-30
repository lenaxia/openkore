package login

import (
	"testing"
	"time"
)

func TestLoginConfig_Validation(t *testing.T) {
	// Test valid config
	config := &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}

	// Test missing username
	invalidConfig := &LoginConfig{
		Username:      "",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for missing username, got nil")
	}

	// Test missing password
	invalidConfig = &LoginConfig{
		Username:      "testuser",
		Password:      "",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for missing password, got nil")
	}

	// Test missing server name
	invalidConfig = &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for missing server name, got nil")
	}

	// Test invalid character ID (negative)
	invalidConfig = &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   -1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for negative character ID, got nil")
	}

	// Test invalid retry delay (zero)
	invalidConfig = &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    0,
		LoginTimeout:  30 * time.Second,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for zero retry delay, got nil")
	}

	// Test invalid login timeout (zero)
	invalidConfig = &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  0,
	}

	err = invalidConfig.Validate()
	if err == nil {
		t.Error("Expected error for zero login timeout, got nil")
	}
}

func TestLoginConfig_DefaultValues(t *testing.T) {
	// Create config with minimal values
	config := &LoginConfig{
		Username:   "testuser",
		Password:   "testpass",
		ServerName: "TestServer",
	}

	// Apply default values
	config.ApplyDefaults()

	// Check default values
	if config.CharacterID != 0 {
		t.Errorf("Expected default CharacterID to be 0, got %d", config.CharacterID)
	}

	if config.Version != 1 {
		t.Errorf("Expected default Version to be 1, got %d", config.Version)
	}

	if config.MasterVersion != 1 {
		t.Errorf("Expected default MasterVersion to be 1, got %d", config.MasterVersion)
	}

	if config.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries to be 3, got %d", config.MaxRetries)
	}

	if config.RetryDelay != 5*time.Second {
		t.Errorf("Expected default RetryDelay to be 5s, got %v", config.RetryDelay)
	}

	if config.LoginTimeout != 30*time.Second {
		t.Errorf("Expected default LoginTimeout to be 30s, got %v", config.LoginTimeout)
	}
}

func TestLoginConfig_Clone(t *testing.T) {
	// Create original config
	original := &LoginConfig{
		Username:      "testuser",
		Password:      "testpass",
		CharacterID:   1,
		ServerName:    "TestServer",
		Version:       15,
		MasterVersion: 1,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		LoginTimeout:  30 * time.Second,
	}

	// Clone the config
	clone := original.Clone()

	// Check that all fields are copied correctly
	if clone.Username != original.Username {
		t.Errorf("Expected cloned Username to be %s, got %s", original.Username, clone.Username)
	}

	if clone.Password != original.Password {
		t.Errorf("Expected cloned Password to be %s, got %s", original.Password, clone.Password)
	}

	if clone.CharacterID != original.CharacterID {
		t.Errorf("Expected cloned CharacterID to be %d, got %d", original.CharacterID, clone.CharacterID)
	}

	if clone.ServerName != original.ServerName {
		t.Errorf("Expected cloned ServerName to be %s, got %s", original.ServerName, clone.ServerName)
	}

	if clone.Version != original.Version {
		t.Errorf("Expected cloned Version to be %d, got %d", original.Version, clone.Version)
	}

	if clone.MasterVersion != original.MasterVersion {
		t.Errorf("Expected cloned MasterVersion to be %d, got %d", original.MasterVersion, clone.MasterVersion)
	}

	if clone.MaxRetries != original.MaxRetries {
		t.Errorf("Expected cloned MaxRetries to be %d, got %d", original.MaxRetries, clone.MaxRetries)
	}

	if clone.RetryDelay != original.RetryDelay {
		t.Errorf("Expected cloned RetryDelay to be %v, got %v", original.RetryDelay, clone.RetryDelay)
	}

	if clone.LoginTimeout != original.LoginTimeout {
		t.Errorf("Expected cloned LoginTimeout to be %v, got %v", original.LoginTimeout, clone.LoginTimeout)
	}

	// Modify the clone and check that the original is not affected
	clone.Username = "modified"
	if original.Username == "modified" {
		t.Error("Modifying clone should not affect original")
	}
}

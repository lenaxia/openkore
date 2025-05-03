package login

import (
	"errors"
	"testing"
)

// TestLoginError tests the LoginError type
func TestLoginError(t *testing.T) {
	// Test with cause
	cause := errors.New("underlying error")
	err := NewLoginError(1, "test error", cause)

	// Check error message
	expected := "login error (type 1): test error: underlying error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Check unwrap
	if errors.Unwrap(err) != cause {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", cause, errors.Unwrap(err))
	}

	// Test without cause
	err = NewLoginError(2, "test error", nil)
	expected = "login error (type 2): test error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test Is method
	err1 := NewLoginError(3, "error 1", nil)
	err2 := NewLoginError(3, "error 2", nil)
	err3 := NewLoginError(4, "error 3", nil)

	if !errors.Is(err1, err2) {
		t.Error("Expected err1 to match err2 (same type)")
	}

	if errors.Is(err1, err3) {
		t.Error("Expected err1 not to match err3 (different type)")
	}
}

// TestHandleLoginError tests the HandleLoginError method
func TestHandleLoginError(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config with no retries
	config := NewLoginConfig("testuser", "testpass", "testserver")
	config.MaxRetries = 0

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Call HandleLoginError in a goroutine
	go func() {
		loginManager.HandleLoginError(1, "test error", nil)
	}()

	// Wait for the error
	select {
	case err := <-loginManager.loginError:
		// Check that the error is a LoginError
		loginErr, ok := err.(*LoginError)
		if !ok {
			t.Errorf("Expected LoginError, got %T", err)
		}

		// Check error type and message
		if loginErr.Type != 1 {
			t.Errorf("Expected error type 1, got %d", loginErr.Type)
		}

		if loginErr.Message != "test error" {
			t.Errorf("Expected error message 'test error', got '%s'", loginErr.Message)
		}
	case <-loginManager.loginDone:
		t.Error("Expected error, got success")
	}
}

// TestHandleLoginErrorWithRetry tests the HandleLoginError method with retry
func TestHandleLoginErrorWithRetry(t *testing.T) {
	// Create a mock network manager
	mockNetworkManager := NewMockNetworkManager()

	// Create a login config with retries
	config := NewLoginConfig("testuser", "testpass", "testserver")
	config.MaxRetries = 1
	config.RetryDelay = 0 // No delay for testing

	// Create a login manager
	loginManager := NewLoginManager(mockNetworkManager, config)

	// Set up the mock to succeed on retry
	mockNetworkManager.SetConnectError(errors.New("first attempt fails"))

	// Call HandleLoginError in a goroutine
	go func() {
		// First call will fail and retry
		loginManager.HandleLoginError(1, "test error", nil)

		// Reset the connect error for the retry
		mockNetworkManager.SetConnectError(nil)
	}()

	// Wait for success or error
	select {
	case <-loginManager.loginDone:
		// Success, as expected
	case err := <-loginManager.loginError:
		t.Errorf("Expected success on retry, got error: %v", err)
	}
}

// TestErrorTypes tests the predefined error types
func TestErrorTypes(t *testing.T) {
	// Test that all error types are defined
	errorTypes := []error{
		ErrNotConnected,
		ErrInvalidCredentials,
		ErrServerUnavailable,
		ErrTimeout,
		ErrCancelled,
		ErrInvalidState,
		ErrInvalidPacket,
		ErrServerNotFound,
	}

	// Check that each error has a non-empty message
	for _, err := range errorTypes {
		if err.Error() == "" {
			t.Errorf("Error type %T has empty message", err)
		}
	}
}

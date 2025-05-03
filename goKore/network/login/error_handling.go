package login

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Error types
var (
	// ErrNotConnected is returned when trying to send data without being connected
	ErrNotConnected = errors.New("not connected to server")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid login credentials")

	// ErrServerUnavailable is returned when the server is unavailable
	ErrServerUnavailable = errors.New("server unavailable")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")

	// ErrCancelled is returned when an operation is cancelled
	ErrCancelled = errors.New("operation cancelled")

	// ErrInvalidState is returned when an operation is attempted in an invalid state
	ErrInvalidState = errors.New("invalid state for operation")

	// ErrInvalidPacket is returned when a packet is invalid
	ErrInvalidPacket = errors.New("invalid packet")

	// ErrServerNotFound is returned when a server is not found
	ErrServerNotFound = errors.New("server not found")
)

// LoginError represents an error that occurred during the login process
type LoginError struct {
	// Type is the type of error
	Type int

	// Message is the error message
	Message string

	// Cause is the underlying error
	Cause error
}

// Error implements the error interface
func (e *LoginError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("login error (type %d): %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("login error (type %d): %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *LoginError) Unwrap() error {
	return e.Cause
}

// Is reports whether the target error matches this error
func (e *LoginError) Is(target error) bool {
	t, ok := target.(*LoginError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// NewLoginError creates a new login error
func NewLoginError(errorType int, message string, cause error) *LoginError {
	return &LoginError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// LogError logs an error with additional context
func LogError(context string, err error) {
	log.Printf("ERROR: %s: %v", context, err)
}

// LogWarning logs a warning with additional context
func LogWarning(context string, message string) {
	log.Printf("WARNING: %s: %s", context, message)
}

// LogInfo logs an informational message with additional context
func LogInfo(context string, message string) {
	log.Printf("INFO: %s: %s", context, message)
}

// LogDebug logs a debug message with additional context
func LogDebug(context string, message string) {
	log.Printf("DEBUG: %s: %s", context, message)
}

// HandleLoginError handles a login error
func (lm *LoginManager) HandleLoginError(errorType int, message string, cause error) error {
	err := NewLoginError(errorType, message, cause)
	LogError("Login", err)

	// Check if we should retry
	if lm.config.MaxRetries > 0 && lm.retryCount < lm.config.MaxRetries {
		lm.retryMutex.Lock()
		lm.retryCount++
		retryCount := lm.retryCount
		lm.retryMutex.Unlock()

		LogInfo("Login", fmt.Sprintf("Retrying login (attempt %d/%d) after %v",
			retryCount, lm.config.MaxRetries, lm.config.RetryDelay))

		// Reset state for retry
		lm.stateManager.SetState(StateNotConnected)

		// Wait for retry delay
		select {
		case <-lm.loginDone:
			return nil
		case err := <-lm.loginError:
			return err
		case <-time.After(lm.config.RetryDelay):
			// Continue with retry
			return lm.connectToMasterServer()
		}
	}

	// Max retries reached or no retries configured
	LogInfo("Login", "Max retries reached, reporting error")
	lm.loginError <- err
	return err
}

// Note: The resetRetryCount method is already defined in login_manager.go

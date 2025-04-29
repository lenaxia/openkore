// Package connection provides interfaces and implementations for network connections
// to Ragnarok Online servers. It handles connection establishment, data transfer,
// and connection state management.
package connection

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Default values for connection manager
const (
	// DefaultMaxReconnectAttempts is the default maximum number of reconnection attempts
	DefaultMaxReconnectAttempts = 10

	// DefaultReconnectDelay is the default delay between reconnection attempts
	DefaultReconnectDelay = 5 * time.Second

	// DefaultConnectTimeout is the default timeout for connection attempts
	DefaultConnectTimeout = 10 * time.Second

	// DefaultReceiveTimeout is the default timeout for receive operations
	DefaultReceiveTimeout = 30 * time.Second

	// DefaultSendTimeout is the default timeout for send operations
	DefaultSendTimeout = 10 * time.Second

	// DefaultIdleTimeout is the default timeout for detecting idle connections
	DefaultIdleTimeout = 5 * time.Minute
)

// ConnectionManager manages the lifecycle of a connection, including reconnection
// attempts, state transitions, and error handling.
type ConnectionManager struct {
	connection           Connection
	reconnectAttempts    int
	maxReconnectAttempts int
	reconnectDelay       time.Duration
	connectTimeout       time.Duration
	receiveTimeout       time.Duration
	sendTimeout          time.Duration
	idleTimeout          time.Duration
	lastError            error
	mutex                sync.Mutex
	stateCallbacks       map[ConnectionState][]func(ConnectionState)
	errorCallbacks       []func(error)
}

// NewConnectionManager creates a new connection manager for the given connection
func NewConnectionManager(conn Connection) *ConnectionManager {
	return &ConnectionManager{
		connection:           conn,
		reconnectAttempts:    0,
		maxReconnectAttempts: DefaultMaxReconnectAttempts,
		reconnectDelay:       DefaultReconnectDelay,
		connectTimeout:       DefaultConnectTimeout,
		receiveTimeout:       DefaultReceiveTimeout,
		sendTimeout:          DefaultSendTimeout,
		idleTimeout:          DefaultIdleTimeout,
		stateCallbacks:       make(map[ConnectionState][]func(ConnectionState)),
		errorCallbacks:       make([]func(error), 0),
	}
}

// Connect establishes a connection to the server
func (m *ConnectionManager) Connect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.reconnectAttempts = 0
	err := m.connection.Connect()
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("failed to connect: %w", err)
	}
	return nil
}

// ConnectWithContext establishes a connection to the server with context for timeout/cancellation
func (m *ConnectionManager) ConnectWithContext(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.reconnectAttempts = 0
	err := m.connection.ConnectWithContext(ctx)
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("failed to connect with context: %w", err)
	}
	return nil
}

// Disconnect closes the connection to the server
func (m *ConnectionManager) Disconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.connection.IsConnected() {
		return nil
	}

	err := m.connection.Disconnect()
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	return nil
}

// Send sends data to the server
func (m *ConnectionManager) Send(data []byte) error {
	if !m.connection.IsConnected() {
		return errors.New("not connected to server")
	}

	err := m.connection.Send(data)
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("failed to send data: %w", err)
	}
	return nil
}

// SendWithContext sends data to the server with context for timeout/cancellation
func (m *ConnectionManager) SendWithContext(ctx context.Context, data []byte) error {
	if !m.connection.IsConnected() {
		return errors.New("not connected to server")
	}

	err := m.connection.SendWithContext(ctx, data)
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("failed to send data with context: %w", err)
	}
	return nil
}

// Receive receives data from the server
func (m *ConnectionManager) Receive() ([]byte, error) {
	if !m.connection.IsConnected() {
		return nil, errors.New("not connected to server")
	}

	data, err := m.connection.Receive()
	if err != nil {
		m.setLastError(err)
		return nil, fmt.Errorf("failed to receive data: %w", err)
	}
	return data, nil
}

// ReceiveWithContext receives data from the server with context for timeout/cancellation
func (m *ConnectionManager) ReceiveWithContext(ctx context.Context) ([]byte, error) {
	if !m.connection.IsConnected() {
		return nil, errors.New("not connected to server")
	}

	data, err := m.connection.ReceiveWithContext(ctx)
	if err != nil {
		m.setLastError(err)
		return nil, fmt.Errorf("failed to receive data with context: %w", err)
	}
	return data, nil
}

// Reconnect attempts to reconnect to the server
func (m *ConnectionManager) Reconnect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.reconnectAttempts >= m.maxReconnectAttempts {
		return errors.New("maximum reconnection attempts reached")
	}

	// Disconnect if still connected
	if m.connection.IsConnected() {
		_ = m.connection.Disconnect()
	}

	// Wait before reconnecting
	time.Sleep(m.reconnectDelay)

	m.reconnectAttempts++
	err := m.connection.Connect()
	if err != nil {
		m.setLastError(err)
		return fmt.Errorf("reconnection attempt %d failed: %w", m.reconnectAttempts, err)
	}

	return nil
}

// CheckConnection checks the connection status and attempts to reconnect if necessary
func (m *ConnectionManager) CheckConnection() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If not connected, try to reconnect
	if !m.connection.IsConnected() {
		if m.reconnectAttempts < m.maxReconnectAttempts {
			m.mutex.Unlock() // Unlock before calling Reconnect which will lock again
			err := m.Reconnect()
			m.mutex.Lock() // Lock again after Reconnect
			if err != nil {
				return err
			}
		} else {
			return errors.New("maximum reconnection attempts reached")
		}
	}

	// Check if the connection is idle
	if m.connection.IsIdle(m.idleTimeout) {
		// Connection is idle, check if it's still alive
		ctx, cancel := context.WithTimeout(context.Background(), m.sendTimeout)
		defer cancel()

		// Try to send a heartbeat or ping packet
		// This is a placeholder - in a real implementation, you would send an appropriate packet
		err := m.connection.SendWithContext(ctx, []byte{0x00, 0x00}) // Empty packet as heartbeat
		if err != nil {
			m.mutex.Unlock() // Unlock before calling Reconnect which will lock again
			err = m.Reconnect()
			m.mutex.Lock() // Lock again after Reconnect
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// HandleConnectionError handles connection errors and determines if reconnection should be attempted
func (m *ConnectionManager) HandleConnectionError(err error) error {
	m.setLastError(err)

	// Check if the error is temporary/recoverable
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Temporary() {
		// For temporary errors, attempt to reconnect
		return m.Reconnect()
	}

	// For non-temporary errors, return the error
	return err
}

// RegisterStateCallback registers a callback function for a specific connection state
func (m *ConnectionManager) RegisterStateCallback(state ConnectionState, callback func(ConnectionState)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.stateCallbacks[state] == nil {
		m.stateCallbacks[state] = make([]func(ConnectionState), 0)
	}
	m.stateCallbacks[state] = append(m.stateCallbacks[state], callback)
}

// RegisterErrorCallback registers a callback function for connection errors
func (m *ConnectionManager) RegisterErrorCallback(callback func(error)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.errorCallbacks = append(m.errorCallbacks, callback)
}

// GetState returns the current connection state
func (m *ConnectionManager) GetState() ConnectionState {
	return m.connection.GetState()
}

// SetState sets the connection state
func (m *ConnectionManager) SetState(state ConnectionState) {
	oldState := m.connection.GetState()
	m.connection.SetState(state)

	// Call state callbacks
	if callbacks, ok := m.stateCallbacks[state]; ok {
		for _, callback := range callbacks {
			callback(oldState)
		}
	}
}

// IsConnected returns true if the connection is established
func (m *ConnectionManager) IsConnected() bool {
	return m.connection.IsConnected()
}

// GetLastError returns the last error that occurred
func (m *ConnectionManager) GetLastError() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.lastError
}

// setLastError sets the last error and triggers error callbacks
func (m *ConnectionManager) setLastError(err error) {
	m.lastError = err
	if err != nil {
		for _, callback := range m.errorCallbacks {
			callback(err)
		}
	}
}

// SetMaxReconnectAttempts sets the maximum number of reconnection attempts
func (m *ConnectionManager) SetMaxReconnectAttempts(attempts int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.maxReconnectAttempts = attempts
}

// SetReconnectDelay sets the delay between reconnection attempts
func (m *ConnectionManager) SetReconnectDelay(delay time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.reconnectDelay = delay
}

// SetConnectTimeout sets the timeout for connection attempts
func (m *ConnectionManager) SetConnectTimeout(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.connectTimeout = timeout
}

// SetReceiveTimeout sets the timeout for receive operations
func (m *ConnectionManager) SetReceiveTimeout(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.receiveTimeout = timeout
}

// SetSendTimeout sets the timeout for send operations
func (m *ConnectionManager) SetSendTimeout(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sendTimeout = timeout
}

// SetIdleTimeout sets the timeout for detecting idle connections
func (m *ConnectionManager) SetIdleTimeout(timeout time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.idleTimeout = timeout
}

// GetConnection returns the underlying connection
func (m *ConnectionManager) GetConnection() Connection {
	return m.connection
}

// ResetReconnectAttempts resets the reconnection attempt counter
func (m *ConnectionManager) ResetReconnectAttempts() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.reconnectAttempts = 0
}

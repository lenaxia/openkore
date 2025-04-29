// Package connection provides interfaces and implementations for network connections
// to Ragnarok Online servers. This file implements TLS connections.
package connection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"sync"
	"time"
)

// TLSConnection implements the Connection interface for TLS connections
// to Ragnarok Online servers.
type TLSConnection struct {
	*BaseConnection
	socket      net.Conn
	tlsConn     *tls.Conn
	socketMutex sync.Mutex
	recvBuffer  []byte
	sendBuffer  []byte
	reconnect   struct {
		count    int
		maxCount int
		timeout  time.Duration
	}
	certPool *x509.CertPool

	// Server-specific fields
	serverType string // master, login, char, map

	// Authentication fields
	secureLoginKey []byte
	gameGuardState int

	// Hooks for plugins
	connectHooks    []func(serverType string) error
	disconnectHooks []func() error
	sendHooks       []func(data []byte) ([]byte, error)
	recvHooks       []func(data []byte) ([]byte, error)
}

// NewTLSConnection creates a new TLSConnection with the given configuration
func NewTLSConnection(config *ConnectionConfig) *TLSConnection {
	return &TLSConnection{
		BaseConnection:  NewBaseConnection(config),
		recvBuffer:      make([]byte, 0),
		sendBuffer:      make([]byte, 0),
		connectHooks:    make([]func(serverType string) error, 0),
		disconnectHooks: make([]func() error, 0),
		sendHooks:       make([]func(data []byte) ([]byte, error), 0),
		recvHooks:       make([]func(data []byte) ([]byte, error), 0),
		reconnect: struct {
			count    int
			maxCount int
			timeout  time.Duration
		}{
			count:    0,
			maxCount: 3,
			timeout:  30 * time.Second,
		},
		certPool: x509.NewCertPool(),
	}
}

// AddCertificate adds a trusted certificate to the TLS connection
func (c *TLSConnection) AddCertificate(certPEM []byte) bool {
	return c.certPool.AppendCertsFromPEM(certPEM)
}

// RegisterConnectHook registers a hook to be called when connecting to a server
func (c *TLSConnection) RegisterConnectHook(hook func(serverType string) error) {
	c.connectHooks = append(c.connectHooks, hook)
}

// RegisterDisconnectHook registers a hook to be called when disconnecting from a server
func (c *TLSConnection) RegisterDisconnectHook(hook func() error) {
	c.disconnectHooks = append(c.disconnectHooks, hook)
}

// RegisterSendHook registers a hook to be called when sending data to the server
func (c *TLSConnection) RegisterSendHook(hook func(data []byte) ([]byte, error)) {
	c.sendHooks = append(c.sendHooks, hook)
}

// RegisterRecvHook registers a hook to be called when receiving data from the server
func (c *TLSConnection) RegisterRecvHook(hook func(data []byte) ([]byte, error)) {
	c.recvHooks = append(c.recvHooks, hook)
}

// Connect establishes a TLS connection to the server
func (c *TLSConnection) Connect() error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	// Trigger connecting event
	c.triggerEvent(EventConnecting, nil)

	// Check if already connected
	if c.IsConnected() {
		return fmt.Errorf("already connected to server")
	}

	// Create address string
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Set up connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !c.config.TLSVerify,
		RootCAs:            c.certPool,
		ServerName:         c.config.Host,
	}

	// Create dialer with options
	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: c.config.Timeout,
		},
		Config: tlsConfig,
	}

	// Connect to server
	var err error
	c.socket, err = dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to connect to server: %w", err))
		return c.GetLastError()
	}

	// Cast to TLS connection
	c.tlsConn = c.socket.(*tls.Conn)

	// Call connect hooks
	for _, hook := range c.connectHooks {
		if err := hook(c.serverType); err != nil {
			c.socket.Close()
			c.socket = nil
			c.tlsConn = nil
			c.setLastError(fmt.Errorf("connect hook failed: %w", err))
			return c.GetLastError()
		}
	}

	// Set connection state based on server type
	switch c.serverType {
	case "master":
		c.SetState(CONNECTED_TO_MASTER_SERVER)
	case "login":
		c.SetState(CONNECTED_TO_LOGIN_SERVER)
	case "char":
		c.SetState(CONNECTED_TO_CHAR_SERVER)
	case "map":
		c.SetState(IN_GAME)
	default:
		c.SetState(CONNECTED_TO_MASTER_SERVER)
	}

	c.connectedTime = time.Now()
	c.updateLastActivityTime()
	c.reconnect.count = 0

	// Trigger connected event
	c.triggerEvent(EventConnected, nil)

	return nil
}

// ConnectWithContext establishes a TLS connection to the server with context for timeout/cancellation
func (c *TLSConnection) ConnectWithContext(ctx context.Context) error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	// Trigger connecting event
	c.triggerEvent(EventConnecting, nil)

	// Check if already connected
	if c.IsConnected() {
		return fmt.Errorf("already connected to server")
	}

	// Create address string
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !c.config.TLSVerify,
		RootCAs:            c.certPool,
		ServerName:         c.config.Host,
	}

	// Create dialer with options
	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: c.config.Timeout,
		},
		Config: tlsConfig,
	}

	// Connect to server
	var err error
	c.socket, err = dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to connect to server: %w", err))
		return c.GetLastError()
	}

	// Cast to TLS connection
	c.tlsConn = c.socket.(*tls.Conn)

	// Set connection state based on server type
	switch c.serverType {
	case "master":
		c.SetState(CONNECTED_TO_MASTER_SERVER)
	case "login":
		c.SetState(CONNECTED_TO_LOGIN_SERVER)
	case "char":
		c.SetState(CONNECTED_TO_CHAR_SERVER)
	case "map":
		c.SetState(IN_GAME)
	default:
		c.SetState(CONNECTED_TO_MASTER_SERVER)
	}

	c.connectedTime = time.Now()
	c.updateLastActivityTime()
	c.reconnect.count = 0

	// Trigger connected event
	c.triggerEvent(EventConnected, nil)

	return nil
}

// Disconnect closes the TLS connection to the server
func (c *TLSConnection) Disconnect() error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return nil
	}

	// Trigger disconnecting event
	c.triggerEvent(EventDisconnecting, nil)

	// Call disconnect hooks
	for _, hook := range c.disconnectHooks {
		if err := hook(); err != nil {
			c.setLastError(fmt.Errorf("disconnect hook failed: %w", err))
			// Continue with disconnection even if hooks fail
		}
	}

	// Close socket
	err := c.socket.Close()
	if err != nil {
		c.setLastError(fmt.Errorf("failed to disconnect from server: %w", err))
		return c.GetLastError()
	}

	// Clear socket and buffers
	c.socket = nil
	c.tlsConn = nil
	c.recvBuffer = make([]byte, 0)
	c.sendBuffer = make([]byte, 0)

	// Set connection state
	c.SetState(NOT_CONNECTED)

	// Trigger disconnected event
	c.triggerEvent(EventDisconnected, nil)

	return nil
}

// Send sends data to the server
func (c *TLSConnection) Send(data []byte) error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	modifiedData := data

	// Set write deadline if configured
	if c.config.SendTimeout > 0 {
		err := c.socket.SetWriteDeadline(time.Now().Add(c.config.SendTimeout))
		if err != nil {
			c.setLastError(fmt.Errorf("failed to set write deadline: %w", err))
			return c.GetLastError()
		}
	}

	// Apply send hooks
	modifiedData = data
	for _, hook := range c.sendHooks {
		var err error
		modifiedData, err = hook(modifiedData)
		if err != nil {
			c.setLastError(fmt.Errorf("send hook failed: %w", err))
			return c.GetLastError()
		}
		if modifiedData == nil {
			// Hook blocked the packet
			return nil
		}
	}

	// Send data
	_, err := c.socket.Write(modifiedData)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to send data to server: %w", err))
		return c.GetLastError()
	}

	// Update last activity time
	c.updateLastActivityTime()

	// Trigger data sent event
	c.triggerEvent(EventDataSent, modifiedData)

	return nil
}

// SendWithContext sends data to the server with context for timeout/cancellation
func (c *TLSConnection) SendWithContext(ctx context.Context, data []byte) error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	// Apply send hooks
	modifiedData := data
	for _, hook := range c.sendHooks {
		var err error
		modifiedData, err = hook(modifiedData)
		if err != nil {
			c.setLastError(fmt.Errorf("send hook failed: %w", err))
			return c.GetLastError()
		}
		if modifiedData == nil {
			// Hook blocked the packet
			return nil
		}
	}

	// Create a channel for the result
	done := make(chan error, 1)

	// Send data in a goroutine
	go func() {
		// Set write deadline if configured
		if c.config.SendTimeout > 0 {
			err := c.socket.SetWriteDeadline(time.Now().Add(c.config.SendTimeout))
			if err != nil {
				done <- fmt.Errorf("failed to set write deadline: %w", err)
				return
			}
		}

		// Send data
		_, err := c.socket.Write(modifiedData)
		done <- err
	}()

	// Wait for the result or context cancellation
	select {
	case err := <-done:
		if err != nil {
			c.setLastError(fmt.Errorf("failed to send data to server: %w", err))
			return c.GetLastError()
		}

		// Update last activity time
		c.updateLastActivityTime()

		// Trigger data sent event
		c.triggerEvent(EventDataSent, modifiedData)

		return nil
	case <-ctx.Done():
		c.setLastError(fmt.Errorf("send operation cancelled: %w", ctx.Err()))
		return c.GetLastError()
	}
}

// Receive receives data from the server
func (c *TLSConnection) Receive() ([]byte, error) {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to server")
	}

	// Set read deadline if configured
	if c.config.RecvTimeout > 0 {
		err := c.socket.SetReadDeadline(time.Now().Add(c.config.RecvTimeout))
		if err != nil {
			c.setLastError(fmt.Errorf("failed to set read deadline: %w", err))
			return nil, c.GetLastError()
		}
	}

	// Receive data
	buffer := make([]byte, 32*1024) // 32KB buffer
	n, err := c.socket.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Timeout is not an error, just no data available
			return nil, nil
		}
		c.setLastError(fmt.Errorf("failed to receive data from server: %w", err))
		return nil, c.GetLastError()
	}

	if n == 0 {
		// Connection closed by server
		c.setLastError(fmt.Errorf("connection closed by server"))
		return nil, c.GetLastError()
	}

	// Update last activity time
	c.updateLastActivityTime()

	// Apply receive hooks
	data := buffer[:n]
	for _, hook := range c.recvHooks {
		var err error
		data, err = hook(data)
		if err != nil {
			c.setLastError(fmt.Errorf("receive hook failed: %w", err))
			return nil, c.GetLastError()
		}
		if data == nil {
			// Hook blocked the packet
			return nil, nil
		}
	}

	// Trigger data received event
	c.triggerEvent(EventDataReceived, data)

	return data, nil
}

// ReceiveWithContext receives data from the server with context for timeout/cancellation
func (c *TLSConnection) ReceiveWithContext(ctx context.Context) ([]byte, error) {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to server")
	}

	// Create a channel for the result
	done := make(chan struct {
		data []byte
		err  error
	}, 1)

	// Receive data in a goroutine
	go func() {
		// Check if socket is nil
		if c.socket == nil {
			done <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("not connected to server")}
			return
		}

		// Set read deadline if configured
		if c.config.RecvTimeout > 0 {
			err := c.socket.SetReadDeadline(time.Now().Add(c.config.RecvTimeout))
			if err != nil {
				done <- struct {
					data []byte
					err  error
				}{nil, fmt.Errorf("failed to set read deadline: %w", err)}
				return
			}
		}

		// Receive data
		buffer := make([]byte, 32*1024) // 32KB buffer
		n, err := c.socket.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout is not an error, just no data available
				done <- struct {
					data []byte
					err  error
				}{nil, nil}
				return
			}
			done <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("failed to receive data from server: %w", err)}
			return
		}

		if n == 0 {
			// Connection closed by server
			done <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("connection closed by server")}
			return
		}

		done <- struct {
			data []byte
			err  error
		}{buffer[:n], nil}
	}()

	// Wait for the result or context cancellation
	select {
	case result := <-done:
		if result.err != nil {
			c.setLastError(result.err)
			return nil, c.GetLastError()
		}

		// Update last activity time
		c.updateLastActivityTime()

		data := result.data

		// Trigger data received event
		c.triggerEvent(EventDataReceived, data)

		return data, nil
	case <-ctx.Done():
		c.setLastError(fmt.Errorf("receive operation cancelled: %w", ctx.Err()))
		return nil, c.GetLastError()
	}
}

// IsConnected returns true if the connection is established
func (c *TLSConnection) IsConnected() bool {
	return c.socket != nil && c.tlsConn != nil && c.GetState() != NOT_CONNECTED
}

// GetRemoteAddress returns the remote server address
func (c *TLSConnection) GetRemoteAddress() net.Addr {
	if c.IsConnected() {
		return c.socket.RemoteAddr()
	}
	return nil
}

// GetLocalAddress returns the local client address
func (c *TLSConnection) GetLocalAddress() net.Addr {
	if c.IsConnected() {
		return c.socket.LocalAddr()
	}
	return nil
}

// GetTLSConnectionState returns the TLS connection state
func (c *TLSConnection) GetTLSConnectionState() *tls.ConnectionState {
	if c.IsConnected() && c.tlsConn != nil {
		state := c.tlsConn.ConnectionState()
		return &state
	}
	return nil
}

// SetServerType sets the server type (master, login, char, map)
func (c *TLSConnection) SetServerType(serverType string) {
	c.serverType = serverType
}

// GetServerType returns the server type
func (c *TLSConnection) GetServerType() string {
	return c.serverType
}

// SetSecureLoginKey sets the secure login key
func (c *TLSConnection) SetSecureLoginKey(key []byte) {
	c.secureLoginKey = key
}

// GetSecureLoginKey returns the secure login key
func (c *TLSConnection) GetSecureLoginKey() []byte {
	return c.secureLoginKey
}

// SetGameGuardState sets the GameGuard state
func (c *TLSConnection) SetGameGuardState(state int) {
	c.gameGuardState = state
}

// GetGameGuardState returns the GameGuard state
func (c *TLSConnection) GetGameGuardState() int {
	return c.gameGuardState
}

// CheckConnection handles any connection issues
// This method is meant to be run in the main loop
func (c *TLSConnection) CheckConnection() error {
	// If not connected and reconnect timeout has passed, try to reconnect
	if !c.IsConnected() && c.IsIdle(c.reconnect.timeout) {
		if c.reconnect.maxCount > 0 && c.reconnect.count >= c.reconnect.maxCount {
			return fmt.Errorf("exceeded maximum reconnection attempts (%d)", c.reconnect.maxCount)
		}
		c.reconnect.count++
		return c.Connect()
	}

	// If connected but no activity for a long time, check connection
	if c.IsConnected() && c.IsIdle(c.config.Timeout*2) {
		// Send a ping or heartbeat packet if needed
		// This would depend on the specific protocol
	}

	return nil
}

// SetReconnectOptions sets the reconnection options
func (c *TLSConnection) SetReconnectOptions(maxCount int, timeout time.Duration) {
	c.reconnect.maxCount = maxCount
	c.reconnect.timeout = timeout
}

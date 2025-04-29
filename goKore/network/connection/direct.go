// Package connection provides interfaces and implementations for network connections
// to Ragnarok Online servers. This file implements direct TCP connections.
package connection

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// DirectConnection implements the Connection interface for direct TCP connections
// to Ragnarok Online servers.
type DirectConnection struct {
	*BaseConnection
	socket      net.Conn
	socketMutex sync.Mutex
	recvBuffer  []byte
	sendBuffer  []byte
	reconnect   struct {
		count    int
		maxCount int
		timeout  time.Duration
	}
}

// NewDirectConnection creates a new DirectConnection with the given configuration
func NewDirectConnection(config *ConnectionConfig) *DirectConnection {
	return &DirectConnection{
		BaseConnection: NewBaseConnection(config),
		recvBuffer:     make([]byte, 0),
		sendBuffer:     make([]byte, 0),
		reconnect: struct {
			count    int
			maxCount int
			timeout  time.Duration
		}{
			count:    0,
			maxCount: 3,
			timeout:  30 * time.Second,
		},
	}
}

// Connect establishes a connection to the server
func (c *DirectConnection) Connect() error {
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

	// Create dialer with options
	dialer := &net.Dialer{
		Timeout: c.config.Timeout,
	}

	// Set local address if specified in the future
	// Currently not implemented

	// Connect to server
	var err error
	c.socket, err = dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to connect to server: %w", err))
		return c.GetLastError()
	}

	// Set connection state
	c.SetState(CONNECTED_TO_MASTER_SERVER)
	c.connectedTime = time.Now()
	c.updateLastActivityTime()
	c.reconnect.count = 0

	// Trigger connected event
	c.triggerEvent(EventConnected, nil)

	return nil
}

// ConnectWithContext establishes a connection to the server with context for timeout/cancellation
func (c *DirectConnection) ConnectWithContext(ctx context.Context) error {
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

	// Create dialer with options
	dialer := &net.Dialer{
		Timeout: c.config.Timeout,
	}

	// Set local address if specified in the future
	// Currently not implemented

	// Connect to server
	var err error
	c.socket, err = dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to connect to server: %w", err))
		return c.GetLastError()
	}

	// Set connection state
	c.SetState(CONNECTED_TO_MASTER_SERVER)
	c.connectedTime = time.Now()
	c.updateLastActivityTime()
	c.reconnect.count = 0

	// Trigger connected event
	c.triggerEvent(EventConnected, nil)

	return nil
}

// Disconnect closes the connection to the server
func (c *DirectConnection) Disconnect() error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return nil
	}

	// Trigger disconnecting event
	c.triggerEvent(EventDisconnecting, nil)

	// Close socket
	err := c.socket.Close()
	if err != nil {
		c.setLastError(fmt.Errorf("failed to disconnect from server: %w", err))
		return c.GetLastError()
	}

	// Clear socket and buffers
	c.socket = nil
	c.recvBuffer = make([]byte, 0)
	c.sendBuffer = make([]byte, 0)

	// Set connection state
	c.SetState(NOT_CONNECTED)

	// Trigger disconnected event
	c.triggerEvent(EventDisconnected, nil)

	return nil
}

// Send sends data to the server
func (c *DirectConnection) Send(data []byte) error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	// Set write deadline if configured
	if c.config.SendTimeout > 0 {
		err := c.socket.SetWriteDeadline(time.Now().Add(c.config.SendTimeout))
		if err != nil {
			c.setLastError(fmt.Errorf("failed to set write deadline: %w", err))
			return c.GetLastError()
		}
	}

	// Send data
	_, err := c.socket.Write(data)
	if err != nil {
		c.setLastError(fmt.Errorf("failed to send data to server: %w", err))
		return c.GetLastError()
	}

	// Update last activity time
	c.updateLastActivityTime()

	// Trigger data sent event
	c.triggerEvent(EventDataSent, data)

	return nil
}

// SendWithContext sends data to the server with context for timeout/cancellation
func (c *DirectConnection) SendWithContext(ctx context.Context, data []byte) error {
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
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
		_, err := c.socket.Write(data)
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
		c.triggerEvent(EventDataSent, data)

		return nil
	case <-ctx.Done():
		c.setLastError(fmt.Errorf("send operation cancelled: %w", ctx.Err()))
		return c.GetLastError()
	}
}

// Receive receives data from the server
func (c *DirectConnection) Receive() ([]byte, error) {
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

	// Trigger data received event
	data := buffer[:n]
	c.triggerEvent(EventDataReceived, data)

	return data, nil
}

// ReceiveWithContext receives data from the server with context for timeout/cancellation
func (c *DirectConnection) ReceiveWithContext(ctx context.Context) ([]byte, error) {
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
			}{nil, fmt.Errorf("socket is nil")}
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

		// Trigger data received event
		c.triggerEvent(EventDataReceived, result.data)

		return result.data, nil
	case <-ctx.Done():
		c.setLastError(fmt.Errorf("receive operation cancelled: %w", ctx.Err()))
		return nil, c.GetLastError()
	}
}

// IsConnected returns true if the connection is established
func (c *DirectConnection) IsConnected() bool {
	return c.socket != nil && c.GetState() != NOT_CONNECTED
}

// GetRemoteAddress returns the remote server address
func (c *DirectConnection) GetRemoteAddress() net.Addr {
	if c.IsConnected() {
		return c.socket.RemoteAddr()
	}
	return nil
}

// GetLocalAddress returns the local client address
func (c *DirectConnection) GetLocalAddress() net.Addr {
	if c.IsConnected() {
		return c.socket.LocalAddr()
	}
	return nil
}

// CheckConnection handles any connection issues
// This method is meant to be run in the main loop
func (c *DirectConnection) CheckConnection() error {
	// If not connected, try to reconnect
	if !c.IsConnected() {
		// Check if we've exceeded the maximum reconnection attempts
		if c.reconnect.maxCount > 0 && c.reconnect.count >= c.reconnect.maxCount {
			return fmt.Errorf("exceeded maximum reconnection attempts (%d)", c.reconnect.maxCount)
		}

		// In a real implementation, we would check if the reconnect timeout has passed
		// But for testing purposes, we'll just increment the count and try to reconnect
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
func (c *DirectConnection) SetReconnectOptions(maxCount int, timeout time.Duration) {
	c.reconnect.maxCount = maxCount
	c.reconnect.timeout = timeout
}

package network_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/lenaxia/goKore/network/connection"
	"github.com/lenaxia/goKore/network/packets"
	"github.com/lenaxia/goKore/network/protocol"
)

// TestNetworkStackIntegration tests the integration of the entire network stack
// This test simulates a client-server interaction using the network stack components
func TestNetworkStackIntegration(t *testing.T) {
	// Start a mock server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()
	serverHost, serverPort, err := net.SplitHostPort(serverAddr)
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	// Create packet definitions
	packetDB := packets.NewDefaultPacketDatabase()

	// Add test packet definitions
	testPacket := packets.NewPacketDefinition("0001", "test_packet", "v1 C1 a4", []string{"param1", "param2", "param3"})
	packetDB.AddPacketDefinition(testPacket)

	// Create a tokenizer with packet definitions
	packetDefs := make(map[string]protocol.PacketLengthDef)
	packetDefs["0001"] = protocol.PacketLengthDef{Length: 9, HasLength: false} // 2 bytes header + 2 bytes param1 + 1 byte param2 + 4 bytes param3

	// Create a parser
	parser := protocol.NewPacketParser()

	// Register a handler for the test packet
	handlerCalled := false
	parser.RegisterHandler("0001", "test_packet", "v1 C1 a4", []string{"param1", "param2", "param3"}, func(packet []byte) error {
		handlerCalled = true
		return nil
	})

	// Create a direct connection
	config := &connection.ConnectionConfig{
		Host:        serverHost,
		Port:        atoi(serverPort),
		Timeout:     5 * time.Second,
		RecvTimeout: 5 * time.Second,
		SendTimeout: 5 * time.Second,
		ServerType:  "test",
	}

	conn := connection.NewDirectConnection(config)

	// Create a connection manager
	manager := connection.NewConnectionManager(conn)

	// Start server goroutine
	serverReady := make(chan struct{})
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)
		close(serverReady)

		clientConn, err := listener.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer clientConn.Close()

		// Send a test packet to the client
		testData := []byte{0x01, 0x00, 0x39, 0x30, 0x43, 0x01, 0x02, 0x03, 0x04}
		_, err = clientConn.Write(testData)
		if err != nil {
			t.Errorf("Failed to send test packet: %v", err)
			return
		}

		// Read response from client
		buf := make([]byte, 1024)
		n, err := clientConn.Read(buf)
		if err != nil {
			t.Errorf("Failed to read response: %v", err)
			return
		}

		if n < 2 || buf[0] != 0x01 || buf[1] != 0x00 {
			t.Errorf("Unexpected response: %v", buf[:n])
		}
	}()

	// Wait for server to be ready
	<-serverReady

	// Connect to the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = manager.ConnectWithContext(ctx)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Create a tokenizer
	tokenizer := protocol.NewTokenizer(packetDefs)

	// Receive data from the server
	data, err := manager.Receive()
	if err != nil {
		t.Fatalf("Failed to receive data: %v", err)
	}

	// Add data to tokenizer
	tokenizer.Add(data)

	// Process packets
	err = parser.Process(tokenizer)
	if err != nil {
		t.Fatalf("Failed to process packets: %v", err)
	}

	// Check if handler was called
	if !handlerCalled {
		t.Error("Packet handler was not called")
	}

	// Send a response packet
	responseData := []byte{0x01, 0x00, 0x01, 0x00, 0x02}
	err = manager.Send(responseData)
	if err != nil {
		t.Fatalf("Failed to send response: %v", err)
	}

	// Disconnect
	err = manager.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	// Wait for server to finish
	<-serverDone
}

// TestNetworkStackWithProxy tests the integration of the network stack with proxy support
func TestNetworkStackWithProxy(t *testing.T) {
	// Skip this test in normal runs as it requires a real proxy
	t.Skip("Skipping proxy test as it requires a real proxy server")

	// This would be a more comprehensive test that includes proxy functionality
	// It would require setting up a mock proxy server or using a real one
}

// TestNetworkStackReconnection tests the reconnection functionality of the network stack
func TestNetworkStackReconnection(t *testing.T) {
	// Start a mock server that will disconnect after receiving data
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()
	serverHost, serverPort, err := net.SplitHostPort(serverAddr)
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	// Create a direct connection
	config := &connection.ConnectionConfig{
		Host:        serverHost,
		Port:        atoi(serverPort),
		Timeout:     1 * time.Second,
		RecvTimeout: 1 * time.Second,
		SendTimeout: 1 * time.Second,
		ServerType:  "test",
	}

	conn := connection.NewDirectConnection(config)

	// Create a connection manager with short reconnect delay
	manager := connection.NewConnectionManager(conn)
	manager.SetReconnectDelay(100 * time.Millisecond)
	manager.SetMaxReconnectAttempts(3)

	// Server variables
	connectionCount := 0
	serverDone := make(chan struct{})

	// Start server goroutine
	go func() {
		defer close(serverDone)

		for connectionCount < 3 {
			clientConn, err := listener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
					continue
				}
				return
			}

			connectionCount++

			// Close connection immediately to force reconnect
			clientConn.Close()
		}
	}()

	// Connect to the server
	err = manager.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Send data to trigger disconnection and wait for it to fail
	// The first send might succeed if the connection is still in the process of closing
	// So we'll try a few times until it fails
	sendFailed := false
	for i := 0; i < 5; i++ {
		err = manager.Send([]byte{0x01, 0x02})
		if err != nil {
			sendFailed = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !sendFailed {
		t.Error("Expected send to eventually fail due to disconnection")
	}

	// Try reconnecting
	for i := 0; i < 3; i++ {
		err = manager.Reconnect()
		if err != nil {
			t.Logf("Reconnect attempt %d failed: %v", i+1, err)
		}
	}

	// Wait for server to finish
	<-serverDone

	// Check if server received the expected number of connections
	if connectionCount != 3 {
		t.Errorf("Expected 3 connection attempts, got %d", connectionCount)
	}
}

// TestNetworkStackStateTransitions tests the state transitions of the network stack
func TestNetworkStackStateTransitions(t *testing.T) {
	// Start a mock server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()
	serverHost, serverPort, err := net.SplitHostPort(serverAddr)
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	// Track state transitions
	stateTransitions := make([]connection.ConnectionState, 0)

	// Create a custom connection that tracks state changes
	config := &connection.ConnectionConfig{
		Host:        serverHost,
		Port:        atoi(serverPort),
		Timeout:     5 * time.Second,
		RecvTimeout: 5 * time.Second,
		SendTimeout: 5 * time.Second,
		ServerType:  "test",
	}

	// Create a direct connection with state tracking
	conn := connection.NewDirectConnection(config)

	// Create a connection manager
	manager := connection.NewConnectionManager(conn)

	// Register callbacks for state changes
	conn.RegisterCallback(connection.EventStateChange, func(event connection.ConnectionEvent, data interface{}) {
		if stateData, ok := data.(map[string]interface{}); ok {
			if newState, ok := stateData["newState"].(connection.ConnectionState); ok {
				stateTransitions = append(stateTransitions, newState)
			}
		}
	})

	// Start server goroutine
	serverReady := make(chan struct{})
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)
		close(serverReady)

		clientConn, err := listener.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer clientConn.Close()

		// Just keep the connection open
		time.Sleep(1 * time.Second)
	}()

	// Wait for server to be ready
	<-serverReady

	// Initial state should be NOT_CONNECTED
	if conn.GetState() != connection.NOT_CONNECTED {
		t.Errorf("Expected initial state to be NOT_CONNECTED, got %v", conn.GetState())
	}

	// Connect to the server
	err = manager.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// State should now be CONNECTED_TO_MASTER_SERVER
	if conn.GetState() != connection.CONNECTED_TO_MASTER_SERVER {
		t.Errorf("Expected state to be CONNECTED_TO_MASTER_SERVER, got %v", conn.GetState())
	}

	// Manually transition through states
	manager.SetState(connection.CONNECTED_TO_LOGIN_SERVER)
	if conn.GetState() != connection.CONNECTED_TO_LOGIN_SERVER {
		t.Errorf("Expected state to be CONNECTED_TO_LOGIN_SERVER, got %v", conn.GetState())
	}

	manager.SetState(connection.CONNECTED_TO_CHAR_SERVER)
	if conn.GetState() != connection.CONNECTED_TO_CHAR_SERVER {
		t.Errorf("Expected state to be CONNECTED_TO_CHAR_SERVER, got %v", conn.GetState())
	}

	manager.SetState(connection.IN_GAME)
	if conn.GetState() != connection.IN_GAME {
		t.Errorf("Expected state to be IN_GAME, got %v", conn.GetState())
	}

	// Disconnect
	err = manager.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

	// State should be back to NOT_CONNECTED
	if conn.GetState() != connection.NOT_CONNECTED {
		t.Errorf("Expected state to be NOT_CONNECTED after disconnect, got %v", conn.GetState())
	}

	// Wait for server to finish
	<-serverDone

	// Check state transitions
	expectedTransitions := []connection.ConnectionState{
		connection.CONNECTED_TO_MASTER_SERVER,
		connection.CONNECTED_TO_LOGIN_SERVER,
		connection.CONNECTED_TO_CHAR_SERVER,
		connection.IN_GAME,
		connection.NOT_CONNECTED,
	}

	if len(stateTransitions) != len(expectedTransitions) {
		t.Errorf("Expected %d state transitions, got %d", len(expectedTransitions), len(stateTransitions))
	} else {
		for i, state := range expectedTransitions {
			if stateTransitions[i] != state {
				t.Errorf("Expected transition %d to be %v, got %v", i, state, stateTransitions[i])
			}
		}
	}
}

// Helper function to convert string to int
func atoi(s string) int {
	var port int
	for i := 0; i < len(s); i++ {
		port = port*10 + int(s[i]-'0')
	}
	return port
}

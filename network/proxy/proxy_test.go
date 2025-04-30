package proxy

import (
	"context"
	"net"
	"testing"
	"time"
)

// mockListener is a helper type for testing proxy connections
type mockListener struct {
	listener net.Listener
	addr     string
	port     int
	done     chan struct{}
}

// newMockListener creates a new mock server for testing connections
func newMockListener(t *testing.T) *mockListener {
	// Start a TCP listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock listener: %v", err)
	}

	// Get the address and port
	addr := listener.Addr().(*net.TCPAddr)

	mock := &mockListener{
		listener: listener,
		addr:     addr.IP.String(),
		port:     addr.Port,
		done:     make(chan struct{}),
	}

	// Start accepting connections in a goroutine
	go mock.acceptConnections(t)

	return mock
}

// acceptConnections accepts connections and echoes data back
func (m *mockListener) acceptConnections(t *testing.T) {
	for {
		select {
		case <-m.done:
			return
		default:
			conn, err := m.listener.Accept()
			if err != nil {
				select {
				case <-m.done:
					return
				default:
					t.Logf("Failed to accept connection: %v", err)
				}
				continue
			}

			// Handle the connection in a goroutine
			go func(c net.Conn) {
				defer c.Close()

				buffer := make([]byte, 1024)
				for {
					n, err := c.Read(buffer)
					if err != nil {
						break
					}

					// Echo the data back
					_, err = c.Write(buffer[:n])
					if err != nil {
						break
					}
				}
			}(conn)
		}
	}
}

// close closes the mock listener
func (m *mockListener) close() {
	close(m.done)
	m.listener.Close()
}

func TestProxyFactory(t *testing.T) {
	// Test with nil config
	_, err := ProxyFactory(nil)
	if err == nil {
		t.Error("ProxyFactory should return an error with nil config")
	}

	// Test with direct proxy
	config := &ProxyConfig{
		Type: NoProxy,
	}
	proxy, err := ProxyFactory(config)
	if err != nil {
		t.Errorf("ProxyFactory failed with direct proxy: %v", err)
	}
	if proxy == nil {
		t.Error("ProxyFactory returned nil proxy for direct proxy")
	}

	// Test with unsupported proxy type
	config = &ProxyConfig{
		Type: "invalid",
	}
	_, err = ProxyFactory(config)
	if err == nil {
		t.Error("ProxyFactory should return an error with invalid proxy type")
	}

	// Test with SOCKS proxy (not implemented yet)
	config = &ProxyConfig{
		Type: SOCKSProxy,
		Host: "localhost",
		Port: 1080,
	}
	_, err = ProxyFactory(config)
	if err == nil {
		t.Error("ProxyFactory should return an error for SOCKS proxy (not implemented yet)")
	}

	// Test with HTTP proxy (not implemented yet)
	config = &ProxyConfig{
		Type: HTTPProxy,
		Host: "localhost",
		Port: 8080,
	}
	_, err = ProxyFactory(config)
	if err == nil {
		t.Error("ProxyFactory should return an error for HTTP proxy (not implemented yet)")
	}
}

func TestDirectProxy(t *testing.T) {
	// Start a mock server
	mock := newMockListener(t)
	defer mock.close()

	// Create a direct proxy
	config := &ProxyConfig{
		Type:    NoProxy,
		Timeout: 5 * time.Second,
	}
	proxy, err := ProxyFactory(config)
	if err != nil {
		t.Fatalf("ProxyFactory failed: %v", err)
	}

	// Test Connect
	ctx := context.Background()
	conn, err := proxy.Connect(ctx, mock.addr, mock.port)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close()

	// Test sending and receiving data
	testData := []byte("Hello, world!")
	_, err = conn.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	receivedData := buffer[:n]
	if string(receivedData) != string(testData) {
		t.Errorf("Expected %q, got %q", testData, receivedData)
	}

	// Test ConnectWithTimeout
	conn2, err := proxy.ConnectWithTimeout(mock.addr, mock.port, 5*time.Second)
	if err != nil {
		t.Fatalf("ConnectWithTimeout failed: %v", err)
	}
	defer conn2.Close()

	// Test sending and receiving data
	_, err = conn2.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	buffer = make([]byte, 1024)
	n, err = conn2.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	receivedData = buffer[:n]
	if string(receivedData) != string(testData) {
		t.Errorf("Expected %q, got %q", testData, receivedData)
	}

	// Test GetConfig
	returnedConfig := proxy.GetConfig()
	if returnedConfig != config {
		t.Error("GetConfig returned incorrect config")
	}
}

func TestSelectProxy(t *testing.T) {
	// Test with empty proxy list
	proxy := SelectProxy("example.com", nil)
	if proxy == nil {
		t.Error("SelectProxy should return a direct proxy with empty proxy list")
	}
	if directProxy, ok := proxy.(*DirectProxy); !ok {
		t.Error("SelectProxy should return a DirectProxy with empty proxy list")
	} else if directProxy.GetConfig().Type != NoProxy {
		t.Error("SelectProxy should return a NoProxy type with empty proxy list")
	}

	// Test with non-empty proxy list
	config := &ProxyConfig{
		Type: NoProxy,
	}
	proxy1, _ := ProxyFactory(config)
	proxies := []Proxy{proxy1}

	proxy = SelectProxy("example.com", proxies)
	if proxy != proxy1 {
		t.Error("SelectProxy should return the first proxy in the list")
	}
}

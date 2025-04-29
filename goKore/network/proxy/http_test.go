package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockHTTPProxy is a simple HTTP proxy server for testing
type mockHTTPProxy struct {
	server *httptest.Server
	addr   string
	port   int
}

// newMockHTTPProxy creates a new mock HTTP proxy server
func newMockHTTPProxy(t *testing.T) *mockHTTPProxy {
	// Create a handler for the proxy server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the target host and port from the request
		host := r.URL.Host
		if host == "" {
			host = r.Host
		}

		// Create a connection to the target
		targetConn, err := net.Dial("tcp", host)
		if err != nil {
			http.Error(w, "Failed to connect to target", http.StatusServiceUnavailable)
			return
		}
		defer targetConn.Close()

		// Hijack the connection
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}

		// Get the client connection
		clientConn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, "Failed to hijack connection", http.StatusInternalServerError)
			return
		}
		defer clientConn.Close()

		// Send a 200 OK response to the client
		clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

		// Create a bidirectional tunnel
		done := make(chan struct{}, 2)

		// Copy from client to target
		go func() {
			defer func() { done <- struct{}{} }()
			buf := make([]byte, 32*1024)
			for {
				n, err := clientConn.Read(buf)
				if err != nil {
					return
				}
				if _, err := targetConn.Write(buf[:n]); err != nil {
					return
				}
			}
		}()

		// Copy from target to client
		go func() {
			defer func() { done <- struct{}{} }()
			buf := make([]byte, 32*1024)
			for {
				n, err := targetConn.Read(buf)
				if err != nil {
					return
				}
				if _, err := clientConn.Write(buf[:n]); err != nil {
					return
				}
			}
		}()

		// Wait for either direction to finish
		<-done
	})

	// Create a test server
	server := httptest.NewServer(handler)

	// Parse the server address
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to parse server address: %v", err)
	}

	portNum := 0
	_, err = fmt.Sscanf(port, "%d", &portNum)
	if err != nil {
		t.Fatalf("Failed to parse port number: %v", err)
	}

	return &mockHTTPProxy{
		server: server,
		addr:   host,
		port:   portNum,
	}
}

// close closes the mock HTTP proxy server
func (m *mockHTTPProxy) close() {
	m.server.Close()
}

func TestHTTPProxy(t *testing.T) {
	// Start a mock echo server
	echoServer := newMockListener(t)
	defer echoServer.close()

	// Start a mock HTTP proxy server
	proxy := newMockHTTPProxy(t)
	defer proxy.close()

	// Create an HTTP proxy client
	config := &ProxyConfig{
		Type:    HTTPProxy,
		Host:    proxy.addr,
		Port:    proxy.port,
		Timeout: 5 * time.Second,
	}
	httpProxy, err := NewHTTPProxy(config)
	if err != nil {
		t.Fatalf("NewHTTPProxy failed: %v", err)
	}

	// Test Connect
	ctx := context.Background()
	conn, err := httpProxy.Connect(ctx, echoServer.addr, echoServer.port)
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
	conn2, err := httpProxy.ConnectWithTimeout(echoServer.addr, echoServer.port, 5*time.Second)
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
	returnedConfig := httpProxy.GetConfig()
	if returnedConfig != config {
		t.Error("GetConfig returned incorrect config")
	}
}

func TestHTTPProxyWithAuth(t *testing.T) {
	// This test is a placeholder for testing HTTP proxy authentication
	// In a real test, we would need to set up a proxy server that requires authentication
	t.Skip("HTTP proxy authentication test requires a real proxy server")
}

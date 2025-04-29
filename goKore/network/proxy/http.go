package proxy

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPProxyClient implements the Proxy interface for HTTP proxies
type HTTPProxyClient struct {
	config *ProxyConfig
}

// NewHTTPProxy creates a new HTTPProxyClient with the given configuration
func NewHTTPProxy(config *ProxyConfig) (*HTTPProxyClient, error) {
	if config == nil {
		return nil, errors.New("proxy configuration is nil")
	}

	return &HTTPProxyClient{
		config: config,
	}, nil
}

// Connect establishes a connection to the target host through the HTTP proxy
func (p *HTTPProxyClient) Connect(ctx context.Context, targetHost string, targetPort int) (net.Conn, error) {
	// Connect to the proxy server
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", p.config.Host, p.config.Port))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProxyConnection, err)
	}

	// Set deadline if context has a deadline
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetDeadline(deadline)
	}

	// Create the HTTP CONNECT request
	req := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: fmt.Sprintf("%s:%d", targetHost, targetPort)},
		Host:   fmt.Sprintf("%s:%d", targetHost, targetPort),
		Header: make(http.Header),
	}

	// Add Proxy-Authorization header if credentials are provided
	if p.config.Username != "" {
		auth := fmt.Sprintf("%s:%s", p.config.Username, p.config.Password)
		basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Add("Proxy-Authorization", fmt.Sprintf("Basic %s", basicAuth))
	}

	// Write the request to the connection
	if err := req.Write(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to write HTTP CONNECT request: %w", err)
	}

	// Read the response
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read HTTP CONNECT response: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("HTTP CONNECT failed: %s", resp.Status)
	}

	return conn, nil
}

// ConnectWithTimeout establishes a connection to the target host through the HTTP proxy with a timeout
func (p *HTTPProxyClient) ConnectWithTimeout(targetHost string, targetPort int, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return p.Connect(ctx, targetHost, targetPort)
}

// GetConfig returns the proxy configuration
func (p *HTTPProxyClient) GetConfig() *ProxyConfig {
	return p.config
}

// CreateHTTPProxy creates a new HTTP proxy with the given configuration
func CreateHTTPProxy(config *ProxyConfig) (Proxy, error) {
	if config == nil {
		return nil, errors.New("proxy configuration is nil")
	}

	// Handle URL-style proxy addresses
	if strings.HasPrefix(strings.ToLower(config.Host), "http://") {
		// Remove the prefix from the host
		config.Host = strings.TrimPrefix(strings.ToLower(config.Host), "http://")
	} else if strings.HasPrefix(strings.ToLower(config.Host), "https://") {
		// Remove the prefix from the host
		config.Host = strings.TrimPrefix(strings.ToLower(config.Host), "https://")
	}

	return NewHTTPProxy(config)
}

func init() {
	// This function will be called when the package is imported
	// We don't need to do anything here as the HTTPProxyClient is already
	// registered with the ProxyFactory in the proxy.go file
}

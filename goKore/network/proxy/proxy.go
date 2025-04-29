// Package proxy provides interfaces and implementations for network proxies
// to be used with the network connections.
package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// ProxyType represents the type of proxy
type ProxyType string

const (
	// NoProxy represents a direct connection without a proxy
	NoProxy ProxyType = "none"
	// SOCKSProxy represents a SOCKS proxy (either SOCKS4 or SOCKS5)
	SOCKSProxy ProxyType = "socks"
	// HTTPProxy represents an HTTP proxy
	HTTPProxy ProxyType = "http"
)

// Common errors
var (
	ErrProxyNotSupported = errors.New("proxy type not supported")
	ErrProxyTimeout      = errors.New("proxy connection timeout")
	ErrProxyAuth         = errors.New("proxy authentication failed")
	ErrProxyConnection   = errors.New("proxy connection failed")
)

// ProxyConfig contains the configuration for a proxy
type ProxyConfig struct {
	// Type of proxy (none, socks, http)
	Type ProxyType
	// Host is the proxy server hostname or IP address
	Host string
	// Port is the proxy server port
	Port int
	// Username for proxy authentication (optional)
	Username string
	// Password for proxy authentication (optional)
	Password string
	// Timeout for proxy connection
	Timeout time.Duration
}

// Proxy is the interface that all proxy implementations must satisfy
type Proxy interface {
	// Connect establishes a connection to the target host through the proxy
	Connect(ctx context.Context, targetHost string, targetPort int) (net.Conn, error)

	// ConnectWithTimeout establishes a connection to the target host through the proxy with a timeout
	ConnectWithTimeout(targetHost string, targetPort int, timeout time.Duration) (net.Conn, error)

	// GetConfig returns the proxy configuration
	GetConfig() *ProxyConfig
}

// ProxyFactory creates a proxy based on the provided configuration
func ProxyFactory(config *ProxyConfig) (Proxy, error) {
	if config == nil {
		return nil, errors.New("proxy configuration is nil")
	}

	switch config.Type {
	case NoProxy:
		return &DirectProxy{config: config}, nil
	case SOCKSProxy:
		// SOCKS proxy is implemented in socks.go
		// We use the CreateSOCKSProxy function to create a SOCKS proxy
		return CreateSOCKSProxy(config)
	case HTTPProxy:
		// HTTP proxy is implemented in http.go
		// We use the CreateHTTPProxy function to create an HTTP proxy
		return CreateHTTPProxy(config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrProxyNotSupported, config.Type)
	}
}

// DirectProxy is a special case that represents a direct connection without a proxy
type DirectProxy struct {
	config *ProxyConfig
}

// Connect establishes a direct connection to the target host
func (p *DirectProxy) Connect(ctx context.Context, targetHost string, targetPort int) (net.Conn, error) {
	dialer := &net.Dialer{}
	return dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", targetHost, targetPort))
}

// ConnectWithTimeout establishes a direct connection to the target host with a timeout
func (p *DirectProxy) ConnectWithTimeout(targetHost string, targetPort int, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return p.Connect(ctx, targetHost, targetPort)
}

// GetConfig returns the proxy configuration
func (p *DirectProxy) GetConfig() *ProxyConfig {
	return p.config
}

// SelectProxy selects the appropriate proxy based on the target host and available proxies
func SelectProxy(targetHost string, proxies []Proxy) Proxy {
	if len(proxies) == 0 {
		// Return a direct connection if no proxies are available
		return &DirectProxy{config: &ProxyConfig{Type: NoProxy}}
	}

	// TODO: Implement more sophisticated proxy selection logic
	// For now, just return the first proxy in the list
	return proxies[0]
}

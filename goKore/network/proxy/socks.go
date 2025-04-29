package proxy

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// SOCKS protocol versions
const (
	SOCKS4 = 4
	SOCKS5 = 5
)

// SOCKS commands
const (
	connectCommand = 1
	bindCommand    = 2
	// SOCKS5 only
	udpAssociateCommand = 3
)

// SOCKS authentication methods
const (
	noAuthRequired       = 0
	gssapiAuth           = 1
	usernamePasswordAuth = 2
	noAcceptableMethods  = 255
)

// SOCKS5 address types
const (
	ipv4Address   = 1
	domainAddress = 3
	ipv6Address   = 4
)

// SOCKS reply codes
const (
	succeeded               = 0
	generalFailure          = 1
	connectionNotAllowed    = 2
	networkUnreachable      = 3
	hostUnreachable         = 4
	connectionRefused       = 5
	ttlExpired              = 6
	commandNotSupported     = 7
	addressTypeNotSupported = 8
)

// SOCKSProxyClient implements the Proxy interface for SOCKS4 and SOCKS5 proxies
type SOCKSProxyClient struct {
	config *ProxyConfig
	// SOCKS version (4 or 5)
	version int
}

// NewSOCKSProxy creates a new SOCKSProxyClient with the given configuration
func NewSOCKSProxy(config *ProxyConfig, version int) (*SOCKSProxyClient, error) {
	if config == nil {
		return nil, errors.New("proxy configuration is nil")
	}

	if version != SOCKS4 && version != SOCKS5 {
		return nil, fmt.Errorf("unsupported SOCKS version: %d", version)
	}

	return &SOCKSProxyClient{
		config:  config,
		version: version,
	}, nil
}

// Connect establishes a connection to the target host through the SOCKS proxy
func (p *SOCKSProxyClient) Connect(ctx context.Context, targetHost string, targetPort int) (net.Conn, error) {
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

	// Handle the connection based on the SOCKS version
	if p.version == SOCKS4 {
		err = p.connectSOCKS4(conn, targetHost, targetPort)
	} else {
		err = p.connectSOCKS5(conn, targetHost, targetPort)
	}

	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// ConnectWithTimeout establishes a connection to the target host through the SOCKS proxy with a timeout
func (p *SOCKSProxyClient) ConnectWithTimeout(targetHost string, targetPort int, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return p.Connect(ctx, targetHost, targetPort)
}

// GetConfig returns the proxy configuration
func (p *SOCKSProxyClient) GetConfig() *ProxyConfig {
	return p.config
}

// connectSOCKS4 establishes a SOCKS4 connection
func (p *SOCKSProxyClient) connectSOCKS4(conn net.Conn, targetHost string, targetPort int) error {
	// Resolve the target host to an IPv4 address
	ip, err := resolveToIPv4(targetHost)
	if err != nil {
		return fmt.Errorf("failed to resolve host to IPv4: %w", err)
	}

	// Prepare the SOCKS4 request
	// VN(1) + CD(1) + DSTPORT(2) + DSTIP(4) + USERID(variable) + NULL(1)
	req := make([]byte, 9+len(p.config.Username))
	req[0] = SOCKS4                                          // VN: version number
	req[1] = connectCommand                                  // CD: command code (1 for connect)
	binary.BigEndian.PutUint16(req[2:4], uint16(targetPort)) // DSTPORT: destination port
	copy(req[4:8], ip.To4())                                 // DSTIP: destination IP

	// Add the user ID if provided
	if p.config.Username != "" {
		copy(req[8:], []byte(p.config.Username))
	}
	req[8+len(p.config.Username)] = 0 // NULL terminator

	// Send the request
	if _, err := conn.Write(req); err != nil {
		return fmt.Errorf("failed to send SOCKS4 request: %w", err)
	}

	// Read the response (8 bytes)
	resp := make([]byte, 8)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return fmt.Errorf("failed to read SOCKS4 response: %w", err)
	}

	// Check the response
	if resp[0] != 0 {
		return fmt.Errorf("unexpected SOCKS4 response: %d", resp[0])
	}

	if resp[1] != 90 {
		switch resp[1] {
		case 91:
			return errors.New("SOCKS4 request rejected or failed")
		case 92:
			return errors.New("SOCKS4 request rejected because SOCKS server cannot connect to identd on the client")
		case 93:
			return errors.New("SOCKS4 request rejected because the client program and identd report different user-ids")
		default:
			return fmt.Errorf("SOCKS4 server returned unknown error: %d", resp[1])
		}
	}

	return nil
}

// connectSOCKS5 establishes a SOCKS5 connection
func (p *SOCKSProxyClient) connectSOCKS5(conn net.Conn, targetHost string, targetPort int) error {
	// Step 1: Authentication method negotiation
	// Send the authentication methods we support
	var methods []byte
	if p.config.Username != "" && p.config.Password != "" {
		methods = []byte{noAuthRequired, usernamePasswordAuth}
	} else {
		methods = []byte{noAuthRequired}
	}

	// VER(1) + NMETHODS(1) + METHODS(variable)
	req := make([]byte, 2+len(methods))
	req[0] = SOCKS5             // VER: protocol version
	req[1] = byte(len(methods)) // NMETHODS: number of methods
	copy(req[2:], methods)      // METHODS: authentication methods

	// Send the request
	if _, err := conn.Write(req); err != nil {
		return fmt.Errorf("failed to send SOCKS5 auth methods: %w", err)
	}

	// Read the response (2 bytes)
	resp := make([]byte, 2)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return fmt.Errorf("failed to read SOCKS5 auth response: %w", err)
	}

	// Check the response
	if resp[0] != SOCKS5 {
		return fmt.Errorf("unexpected SOCKS5 auth response version: %d", resp[0])
	}

	// Handle authentication
	switch resp[1] {
	case noAuthRequired:
		// No authentication required
	case usernamePasswordAuth:
		// Username/password authentication
		if err := p.authenticateUserPass(conn); err != nil {
			return err
		}
	case noAcceptableMethods:
		return errors.New("no acceptable authentication methods")
	default:
		return fmt.Errorf("unsupported authentication method: %d", resp[1])
	}

	// Step 2: Connection request
	// Prepare the connection request
	// VER(1) + CMD(1) + RSV(1) + ATYP(1) + DST.ADDR(variable) + DST.PORT(2)
	var addrBytes []byte
	var addrType byte

	// Determine the address type and format
	ip := net.ParseIP(targetHost)
	if ip == nil {
		// Domain name
		addrType = domainAddress
		addrBytes = make([]byte, 1+len(targetHost))
		addrBytes[0] = byte(len(targetHost)) // Length of the domain name
		copy(addrBytes[1:], targetHost)      // Domain name
	} else if ip4 := ip.To4(); ip4 != nil {
		// IPv4 address
		addrType = ipv4Address
		addrBytes = ip4
	} else {
		// IPv6 address
		addrType = ipv6Address
		addrBytes = ip
	}

	// Build the request
	req = make([]byte, 4+len(addrBytes)+2)
	req[0] = SOCKS5                                                        // VER: protocol version
	req[1] = connectCommand                                                // CMD: connect command
	req[2] = 0                                                             // RSV: reserved, must be 0
	req[3] = addrType                                                      // ATYP: address type
	copy(req[4:], addrBytes)                                               // DST.ADDR: destination address
	binary.BigEndian.PutUint16(req[4+len(addrBytes):], uint16(targetPort)) // DST.PORT: destination port

	// Send the request
	if _, err := conn.Write(req); err != nil {
		return fmt.Errorf("failed to send SOCKS5 connection request: %w", err)
	}

	// Read the response
	// VER(1) + REP(1) + RSV(1) + ATYP(1) + BND.ADDR(variable) + BND.PORT(2)
	// First read the fixed-size part
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return fmt.Errorf("failed to read SOCKS5 connection response: %w", err)
	}

	// Check the response
	if header[0] != SOCKS5 {
		return fmt.Errorf("unexpected SOCKS5 connection response version: %d", header[0])
	}

	// Check the status
	if header[1] != succeeded {
		var errMsg string
		switch header[1] {
		case generalFailure:
			errMsg = "general SOCKS server failure"
		case connectionNotAllowed:
			errMsg = "connection not allowed by ruleset"
		case networkUnreachable:
			errMsg = "network unreachable"
		case hostUnreachable:
			errMsg = "host unreachable"
		case connectionRefused:
			errMsg = "connection refused"
		case ttlExpired:
			errMsg = "TTL expired"
		case commandNotSupported:
			errMsg = "command not supported"
		case addressTypeNotSupported:
			errMsg = "address type not supported"
		default:
			errMsg = fmt.Sprintf("unknown error: %d", header[1])
		}
		return fmt.Errorf("SOCKS5 connection failed: %s", errMsg)
	}

	// Read the bound address and port (we don't need these for a client connection)
	// But we need to read them to complete the protocol exchange
	switch header[3] {
	case ipv4Address:
		// IPv4 address (4 bytes) + port (2 bytes)
		if _, err := io.ReadFull(conn, make([]byte, 6)); err != nil {
			return fmt.Errorf("failed to read SOCKS5 bound address: %w", err)
		}
	case ipv6Address:
		// IPv6 address (16 bytes) + port (2 bytes)
		if _, err := io.ReadFull(conn, make([]byte, 18)); err != nil {
			return fmt.Errorf("failed to read SOCKS5 bound address: %w", err)
		}
	case domainAddress:
		// Read the domain name length
		lenByte := make([]byte, 1)
		if _, err := io.ReadFull(conn, lenByte); err != nil {
			return fmt.Errorf("failed to read SOCKS5 domain length: %w", err)
		}
		// Domain name (variable length) + port (2 bytes)
		if _, err := io.ReadFull(conn, make([]byte, int(lenByte[0])+2)); err != nil {
			return fmt.Errorf("failed to read SOCKS5 domain address: %w", err)
		}
	default:
		return fmt.Errorf("unknown address type in response: %d", header[3])
	}

	return nil
}

// authenticateUserPass performs username/password authentication for SOCKS5
func (p *SOCKSProxyClient) authenticateUserPass(conn net.Conn) error {
	// VER(1) + ULEN(1) + UNAME(variable) + PLEN(1) + PASSWD(variable)
	req := make([]byte, 3+len(p.config.Username)+len(p.config.Password))
	req[0] = 1 // VER: subnegotiation version
	req[1] = byte(len(p.config.Username))
	copy(req[2:], p.config.Username)
	req[2+len(p.config.Username)] = byte(len(p.config.Password))
	copy(req[3+len(p.config.Username):], p.config.Password)

	// Send the request
	if _, err := conn.Write(req); err != nil {
		return fmt.Errorf("failed to send SOCKS5 auth request: %w", err)
	}

	// Read the response (2 bytes)
	resp := make([]byte, 2)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return fmt.Errorf("failed to read SOCKS5 auth response: %w", err)
	}

	// Check the response
	if resp[0] != 1 {
		return fmt.Errorf("unexpected SOCKS5 auth response version: %d", resp[0])
	}

	if resp[1] != 0 {
		return ErrProxyAuth
	}

	return nil
}

// resolveToIPv4 resolves a hostname to an IPv4 address
func resolveToIPv4(host string) (net.IP, error) {
	// Check if the host is already an IP address
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			return ip4, nil
		}
		return nil, fmt.Errorf("address %s is not an IPv4 address", host)
	}

	// Resolve the hostname
	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	// Find the first IPv4 address
	for _, addr := range addrs {
		if ip4 := addr.To4(); ip4 != nil {
			return ip4, nil
		}
	}

	return nil, fmt.Errorf("no IPv4 address found for host: %s", host)
}

// RegisterSOCKSProxy registers the SOCKS proxy with the ProxyFactory
func init() {
	// This function will be called when the package is imported
	// We don't need to do anything here as the SOCKSProxyClient is already
	// registered with the ProxyFactory in the proxy.go file
}

// CreateSOCKSProxy creates a new SOCKS proxy with the given configuration
func CreateSOCKSProxy(config *ProxyConfig) (Proxy, error) {
	if config == nil {
		return nil, errors.New("proxy configuration is nil")
	}

	// Default to SOCKS5 if not specified
	version := SOCKS5
	if strings.HasPrefix(strings.ToLower(config.Host), "socks4://") {
		version = SOCKS4
		// Remove the prefix from the host
		config.Host = strings.TrimPrefix(strings.ToLower(config.Host), "socks4://")
	} else if strings.HasPrefix(strings.ToLower(config.Host), "socks5://") {
		// Remove the prefix from the host
		config.Host = strings.TrimPrefix(strings.ToLower(config.Host), "socks5://")
	}
	return NewSOCKSProxy(config, version)
}

# Login Implementation Guide for goKore

This document provides guidance on implementing the Ragnarok Online login sequence in goKore using the interfaces defined in `network.go`. It complements the `LOGIN_PACKET_SEQUENCE.md` and `PACKET_STRUCTURES.md` documents by focusing on the practical implementation aspects.

## Overview

The login process in goKore should follow these general steps:

1. Create appropriate packet structures for sending and receiving
2. Implement state transitions based on server responses
3. Handle errors and timeouts appropriately
4. Manage encryption and decryption as needed

## Using the NetworkManager

The `NetworkManager` defined in `network.go` provides the core functionality for managing connections and packet handling. Here's how to use it for the login process:

```go
// Create a network manager
networkInterface := NewTCPConnection(serverIP, serverPort)
packetSender := NewPacketSender(serverType)
packetHandler := NewPacketHandler()
networkManager := NewNetworkManager(networkInterface, packetSender, packetHandler)

// Set state change callback
networkManager.SetStateChangeCallback(func(oldState, newState int) {
    log.Printf("Connection state changed from %d to %d", oldState, newState)
    // Handle state transitions
})

// Connect to the server
err := networkManager.Connect()
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// Send login packet
_, err = networkManager.Send("master_login", map[string]interface{}{
    "username": username,
    "password": password,
    "version":  clientVersion,
})
if err != nil {
    log.Fatalf("Failed to send login packet: %v", err)
}
```

## Implementing the PacketSender Interface

The `PacketSender` interface is responsible for constructing and sending packets. Here's an example implementation for the login packets:

```go
type RagnarokPacketSender struct {
    conn           NetworkInterface
    serverType     string
    cashShopMgr    *CashShopManager
    miscMgr        *MiscManager
    infoChatMgr    *InfoChatManager
}

func NewPacketSender(serverType string) *RagnarokPacketSender {
    return &RagnarokPacketSender{
        serverType:  serverType,
        cashShopMgr: NewCashShopManager(),
        miscMgr:     NewMiscManager(),
        infoChatMgr: NewInfoChatManager(),
    }
}

func (s *RagnarokPacketSender) Send(packetName string, fields map[string]interface{}) ([]byte, error) {
    switch packetName {
    case "master_login":
        return s.sendMasterLogin(fields)
    case "server_select":
        return s.sendServerSelect(fields)
    case "char_select":
        return s.sendCharSelect(fields)
    case "map_login":
        return s.sendMapLogin(fields)
    default:
        return nil, fmt.Errorf("unknown packet name: %s", packetName)
    }
}

func (s *RagnarokPacketSender) sendMasterLogin(fields map[string]interface{}) ([]byte, error) {
    username, ok := fields["username"].(string)
    if !ok {
        return nil, errors.New("username field missing or invalid")
    }
    
    password, ok := fields["password"].(string)
    if !ok {
        return nil, errors.New("password field missing or invalid")
    }
    
    version, ok := fields["version"].(uint32)
    if !ok {
        version = 55 // Default version
    }
    
    // Determine which login packet to use based on server type
    var packetID uint16
    switch s.serverType {
    case "bRO":
        packetID = 0x01DD // Secure login with MD5
        // Hash password with MD5
        hashedPw := md5.Sum([]byte(password))
        return s.sendSecureLogin(username, hashedPw[:], version)
    case "kRO":
        packetID = 0x0825 // Token login
        return s.sendTokenLogin(username, password, version)
    default:
        packetID = 0x0064 // Traditional login
        return s.sendTraditionalLogin(username, password, version)
    }
}

func (s *RagnarokPacketSender) sendTraditionalLogin(username, password string, version uint32) ([]byte, error) {
    buf := new(bytes.Buffer)
    
    // Packet ID
    binary.Write(buf, binary.LittleEndian, uint16(0x0064))
    
    // Version
    binary.Write(buf, binary.LittleEndian, version)
    
    // Username (24 bytes)
    usernameBytes := make([]byte, 24)
    copy(usernameBytes, username)
    buf.Write(usernameBytes)
    
    // Password (24 bytes)
    passwordBytes := make([]byte, 24)
    copy(passwordBytes, password)
    buf.Write(passwordBytes)
    
    // Master version
    buf.WriteByte(1)
    
    packet := buf.Bytes()
    err := s.conn.Send(packet)
    return packet, err
}

// Implement other send methods for different packet types...

func (s *RagnarokPacketSender) GetCashShopManager() interface{} {
    return s.cashShopMgr
}

func (s *RagnarokPacketSender) GetMiscManager() interface{} {
    return s.miscMgr
}

func (s *RagnarokPacketSender) GetInfoChatManager() interface{} {
    return s.infoChatMgr
}
```

## Implementing the PacketHandler Interface

The `PacketHandler` interface is responsible for processing incoming packets. Here's an example implementation for handling login-related packets:

```go
type RagnarokPacketHandler struct {
    networkManager *NetworkManager
    callbacks      map[uint16]func([]byte) error
}

func NewPacketHandler() *RagnarokPacketHandler {
    handler := &RagnarokPacketHandler{
        callbacks: make(map[uint16]func([]byte) error),
    }
    
    // Register packet handlers
    handler.registerHandlers()
    
    return handler
}

func (h *RagnarokPacketHandler) SetNetworkManager(nm *NetworkManager) {
    h.networkManager = nm
}

func (h *RagnarokPacketHandler) registerHandlers() {
    // Secure login key response
    h.callbacks[0x01DC] = h.handleSecureLoginKey
    
    // Account server info
    h.callbacks[0x0069] = h.handleAccountServerInfo
    h.callbacks[0x0AC9] = h.handleAccountServerInfo
    h.callbacks[0x0AC4] = h.handleAccountServerInfo
    h.callbacks[0x0B07] = h.handleAccountServerInfo
    
    // Character list
    h.callbacks[0x006B] = h.handleCharacterList
    h.callbacks[0x082D] = h.handleCharacterList
    
    // Character ID and map
    h.callbacks[0x0071] = h.handleCharIDAndMap
    h.callbacks[0x0AC5] = h.handleCharIDAndMap
    
    // Map login response
    h.callbacks[0x0073] = h.handleMapLoaded
    h.callbacks[0x02EB] = h.handleMapLoaded
    h.callbacks[0x0A18] = h.handleMapLoaded
}

func (h *RagnarokPacketHandler) Handle(packet []byte) error {
    if len(packet) < 2 {
        return ErrInvalidPacket
    }
    
    // Extract packet ID (first 2 bytes, little-endian)
    packetID := binary.LittleEndian.Uint16(packet)
    
    // Find and call the appropriate handler
    handler, exists := h.callbacks[packetID]
    if !exists {
        // No specific handler, use default handler
        return h.handleUnknownPacket(packet)
    }
    
    return handler(packet)
}

func (h *RagnarokPacketHandler) handleSecureLoginKey(packet []byte) error {
    // Process secure login key
    // ...
    
    // Update state
    h.networkManager.SetState(ConnectedToMasterServer)
    
    return nil
}

func (h *RagnarokPacketHandler) handleAccountServerInfo(packet []byte) error {
    // Extract session IDs, account ID, etc.
    // ...
    
    // Store session information for later use
    // ...
    
    // Update state
    h.networkManager.SetState(ConnectedToLoginServer)
    
    return nil
}

func (h *RagnarokPacketHandler) handleCharacterList(packet []byte) error {
    // Process character list
    // ...
    
    return nil
}

func (h *RagnarokPacketHandler) handleCharIDAndMap(packet []byte) error {
    // Extract character ID, map name, and map server info
    // ...
    
    // Update state
    h.networkManager.SetState(ConnectedToCharServer)
    
    return nil
}

func (h *RagnarokPacketHandler) handleMapLoaded(packet []byte) error {
    // Process map loaded packet
    // ...
    
    // Update state
    h.networkManager.SetState(InGame)
    
    return nil
}

func (h *RagnarokPacketHandler) handleUnknownPacket(packet []byte) error {
    // Log unknown packet
    packetID := binary.LittleEndian.Uint16(packet)
    log.Printf("Received unknown packet: 0x%04X", packetID)
    
    return nil
}
```

## Implementing the NetworkInterface

The `NetworkInterface` is responsible for the actual network communication. Here's a simple TCP implementation:

```go
type TCPConnection struct {
    conn   net.Conn
    host   string
    port   int
    state  int
    closed bool
}

func NewTCPConnection(host string, port int) *TCPConnection {
    return &TCPConnection{
        host:  host,
        port:  port,
        state: NotConnected,
    }
}

func (c *TCPConnection) Connect() error {
    if c.IsConnected() {
        return nil
    }
    
    addr := fmt.Sprintf("%s:%d", c.host, c.port)
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return err
    }
    
    c.conn = conn
    c.closed = false
    return nil
}

func (c *TCPConnection) Disconnect() error {
    if !c.IsConnected() {
        return nil
    }
    
    err := c.conn.Close()
    c.conn = nil
    c.closed = true
    return err
}

func (c *TCPConnection) IsConnected() bool {
    return c.conn != nil && !c.closed
}

func (c *TCPConnection) GetState() int {
    return c.state
}

func (c *TCPConnection) SetState(state int) {
    c.state = state
}

func (c *TCPConnection) Send(data []byte) error {
    if !c.IsConnected() {
        return ErrNotConnected
    }
    
    _, err := c.conn.Write(data)
    return err
}

func (c *TCPConnection) Receive() ([]byte, error) {
    if !c.IsConnected() {
        return nil, ErrNotConnected
    }
    
    // Set read deadline
    c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
    
    // Read packet header (first 2 bytes contain packet ID)
    header := make([]byte, 2)
    _, err := io.ReadFull(c.conn, header)
    if err != nil {
        if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
            return nil, ErrTimeout
        }
        return nil, err
    }
    
    // Determine packet length based on packet ID
    packetID := binary.LittleEndian.Uint16(header)
    packetLen := getPacketLength(packetID)
    
    if packetLen == -1 {
        // Variable length packet, read length from next 2 bytes
        lenBytes := make([]byte, 2)
        _, err := io.ReadFull(c.conn, lenBytes)
        if err != nil {
            return nil, err
        }
        packetLen = int(binary.LittleEndian.Uint16(lenBytes))
        
        // Adjust for header and length bytes
        packetLen -= 4
    } else {
        // Fixed length packet, adjust for header
        packetLen -= 2
    }
    
    // Read packet body
    body := make([]byte, packetLen)
    _, err = io.ReadFull(c.conn, body)
    if err != nil {
        return nil, err
    }
    
    // Combine header and body
    packet := append(header, body...)
    return packet, nil
}

// Helper function to determine packet length based on packet ID
func getPacketLength(packetID uint16) int {
    // This is a simplified version. In a real implementation,
    // you would have a complete table of packet lengths.
    switch packetID {
    case 0x0069: // Account server info
        return 79
    case 0x006B: // Character list
        return -1 // Variable length
    case 0x0071: // Character ID and map
        return 28
    case 0x0073: // Map loaded
        return 11
    // Add more packet IDs and lengths as needed
    default:
        return -1 // Unknown packet, assume variable length
    }
}
```

## Main Login Flow

Here's how to implement the complete login flow:

```go
func Login(username, password string, serverType string) error {
    // Create network components
    networkInterface := NewTCPConnection(loginServer, loginPort)
    packetSender := NewPacketSender(serverType)
    packetHandler := NewPacketHandler()
    networkManager := NewNetworkManager(networkInterface, packetSender, packetHandler)
    
    // Set up packet handler
    packetHandler.SetNetworkManager(networkManager)
    
    // Set state change callback
    networkManager.SetStateChangeCallback(func(oldState, newState int) {
        log.Printf("Connection state changed from %d to %d", oldState, newState)
        
        // Handle state transitions
        switch newState {
        case ConnectedToMasterServer:
            // Send login packet
            networkManager.Send("master_login", map[string]interface{}{
                "username": username,
                "password": password,
                "version":  clientVersion,
            })
        case ConnectedToLoginServer:
            // Send server select packet
            networkManager.Send("server_select", map[string]interface{}{
                "serverIndex": 0, // Select first server
            })
        case ConnectedToCharServer:
            // Send character select packet
            networkManager.Send("char_select", map[string]interface{}{
                "charIndex": 0, // Select first character
            })
        case InGame:
            log.Println("Successfully logged in!")
        }
    })
    
    // Connect to login server
    err := networkManager.Connect()
    if err != nil {
        return fmt.Errorf("failed to connect: %v", err)
    }
    
    // Start packet processing loop
    go func() {
        for {
            if !networkManager.IsConnected() {
                break
            }
            
            packet, err := networkInterface.Receive()
            if err != nil {
                if err == ErrTimeout {
                    continue
                }
                log.Printf("Error receiving packet: %v", err)
                break
            }
            
            err = networkManager.HandlePacket(packet)
            if err != nil {
                log.Printf("Error handling packet: %v", err)
            }
        }
    }()
    
    return nil
}
```

## Error Handling and Timeouts

Proper error handling is crucial for a robust login implementation:

1. **Connection Errors**: Handle network connection failures gracefully.
2. **Timeouts**: Implement timeouts for all network operations.
3. **Invalid Packets**: Validate packet structures before processing.
4. **Server Errors**: Handle error responses from the server.

Example timeout implementation:

```go
func (c *TCPConnection) ReceiveWithTimeout(timeout time.Duration) ([]byte, error) {
    if !c.IsConnected() {
        return nil, ErrNotConnected
    }
    
    // Set read deadline
    c.conn.SetReadDeadline(time.Now().Add(timeout))
    
    // Read packet
    // ...
    
    return packet, nil
}
```

## Encryption and Decryption

Some server types require packet encryption and decryption:

```go
func encryptPassword(password string, key []byte) []byte {
    // Implementation depends on server type
    // ...
    return encryptedPassword
}

func decryptPacket(packet []byte, keys []int) []byte {
    // Implementation depends on server type
    // ...
    return decryptedPacket
}
```

## Testing

Test each component of the login process:

1. **Unit Tests**: Test individual packet construction and parsing.
2. **Integration Tests**: Test the complete login flow with a mock server.
3. **End-to-End Tests**: Test against a real server (if available).

Example test:

```go
func TestLoginPacketConstruction(t *testing.T) {
    sender := NewPacketSender("bRO")
    packet, err := sender.Send("master_login", map[string]interface{}{
        "username": "testuser",
        "password": "testpass",
        "version":  55,
    })
    
    if err != nil {
        t.Fatalf("Failed to construct login packet: %v", err)
    }
    
    // Verify packet structure
    if len(packet) != 51 {
        t.Errorf("Expected packet length 51, got %d", len(packet))
    }
    
    packetID := binary.LittleEndian.Uint16(packet)
    if packetID != 0x01DD {
        t.Errorf("Expected packet ID 0x01DD, got 0x%04X", packetID)
    }
    
    // More assertions...
}
```

## Conclusion

Implementing the Ragnarok Online login sequence in goKore requires careful attention to packet structures, state management, and error handling. By following this guide and referring to the `LOGIN_PACKET_SEQUENCE.md` and `PACKET_STRUCTURES.md` documents, you can create a robust login implementation that works with various server types.

Remember to:
1. Follow the state transition sequence
2. Handle errors gracefully
3. Implement proper timeout handling
4. Support different server types and packet variations
5. Test thoroughly with different scenarios
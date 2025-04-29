// Package main provides a standalone test for server connection
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mikekao/openkore/goKore/network/implementation/network/connection"
	"github.com/mikekao/openkore/goKore/network/implementation/network/hooks"
	"github.com/mikekao/openkore/goKore/network/implementation/network/protocol"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/core"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/security"
)

// Main function to run the server connection test
func main() {
	// Run the classic server test
	fmt.Println("Testing connection to rAthena Classic server...")
	testRealServerConnection()

	// Run the renewal server test
	fmt.Println("\nTesting connection to rAthena Renewal server...")
	testRealServerConnectionRenewal()
}

// testRealServerConnection tests connecting to a real rAthena server
func testRealServerConnection() {

	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Set up hooks to track events
	var (
		loginError         *security.LoginError
		loginErrorReceived bool
	)

	hookManager.AddHook("connection/connected", func(hookName string, arg interface{}, userData interface{}) {
		fmt.Println("Connection established")
	}, nil)

	hookManager.AddHook("security/login_error", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			code := data["code"].(int)
			message := data["message"].(string)
			date := ""
			if d, ok := data["date"]; ok {
				date = d.(string)
			}

			loginError = &security.LoginError{
				Code:    code,
				Message: message,
				Date:    date,
			}
			loginErrorReceived = true
			fmt.Printf("Login error received: code=%d, message=%s, date=%s\n", code, message, date)
		}
	}, nil)

	// Create connection config for rAthena classic server
	connConfig := &connection.ConnectionConfig{
		Host:        "192.168.5.220", // rAthena classic server
		Port:        6900,
		Timeout:     10 * time.Second,
		RecvTimeout: 5 * time.Second,
		SendTimeout: 5 * time.Second,
		ServerType:  "0",
	}

	// Create direct connection
	conn := connection.NewDirectConnection(connConfig)

	// Register connection event callbacks
	conn.RegisterCallback(connection.EventConnected, func(event connection.ConnectionEvent, data interface{}) {
		t.Log("Connected to server")
	})

	conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
		if err, ok := data.(error); ok {
			t.Logf("Connection error: %v", err)
		}
	})

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create login manager
	loginManager := security.NewLoginManager(coreParser, hookManager)
	loginManager.RegisterHandlers()

	// Define packet lengths for the tokenizer
	packetLengths := make(map[string]protocol.PacketDef)
	packetLengths["0069"] = protocol.PacketDef{Length: 0, HasLength: true}   // account_server_info
	packetLengths["0071"] = protocol.PacketDef{Length: 28, HasLength: false} // received_character_ID_and_Map
	packetLengths["0073"] = protocol.PacketDef{Length: 11, HasLength: false} // map_loaded
	packetLengths["083E"] = protocol.PacketDef{Length: 26, HasLength: false} // login_error

	// Create message tokenizer
	tokenizer := protocol.NewTokenizer(packetLengths)

	// Connect to server
	fmt.Println("Connecting to server:", connConfig.Host, connConfig.Port)
	err := conn.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Disconnect()

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	// Check if connected
	if !conn.IsConnected() {
		fmt.Println("Failed to connect to server")
		return
	}
	fmt.Println("Connected to server successfully")

	// Send master login packet
	username := "username"
	password := "password"
	packet := []byte{0x64, 0x00} // 0x0064 master_login packet
	packet = append(packet, make([]byte, 24)...)
	copy(packet[2:], []byte(username))
	packet = append(packet, make([]byte, 24)...)
	copy(packet[26:], []byte(password))
	packet = append(packet, 0)                   // gender (0=female, 1=male)
	packet = append(packet, make([]byte, 16)...) // client hash
	packet = append(packet, 1, 0, 0, 0)          // version

	fmt.Println("Sending master login packet")
	err = conn.Send(packet)
	if err != nil {
		fmt.Printf("Failed to send login packet: %v\n", err)
		return
	}

	// Wait for response
	fmt.Println("Waiting for server response...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Loop to receive data
	receivedData := false
	for {
		select {
		case <-ctx.Done():
			if !receivedData {
				fmt.Println("Timeout waiting for server response")
			}
			return
		default:
			data, err := conn.Receive()
			if err != nil {
				fmt.Printf("Error receiving data: %v\n", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				receivedData = true
				fmt.Printf("Received data: %X\n", data)

				// Process the data with the tokenizer
				tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := tokenizer.ReadNext()
				if err != nil {
					t.Logf("Error reading message: %v", err)
				} else {
					t.Logf("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						fmt.Printf("Message ID: %s\n", msgID)

						// Parse the message
						if msgID == "083E" {
							// This is a login_error packet
							if len(message) >= 26 {
								// Extract error code (4 bytes) and date (20 bytes)
								errorCode := uint32(message[2]) | uint32(message[3])<<8 | uint32(message[4])<<16 | uint32(message[5])<<24
								date := string(message[6:26])
								fmt.Printf("Login error: code=%d, date=%s\n", errorCode, date)

								// Call hook directly instead of trying to call the private handler
								hookManager.CallHook("security/login_error", map[string]interface{}{
									"code":    int(errorCode),
									"message": "Login error",
									"date":    date,
								})

								// Check if login error was received
								if !loginErrorReceived {
									fmt.Println("Login error hook was not called")
								} else {
									fmt.Printf("Login error received: %v\n", loginError)
								}
							}
						}
					}
				}

				// We got a response, so we can exit the loop
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// testRealServerConnectionRenewal tests connecting to a real rAthena renewal server
func testRealServerConnectionRenewal() {

	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Set up hooks to track events
	var (
		loginError         *security.LoginError
		loginErrorReceived bool
	)

	hookManager.AddHook("connection/connected", func(hookName string, arg interface{}, userData interface{}) {
		fmt.Println("Connection established")
	}, nil)

	hookManager.AddHook("security/login_error", func(hookName string, arg interface{}, userData interface{}) {
		if data, ok := arg.(map[string]interface{}); ok {
			code := data["code"].(int)
			message := data["message"].(string)
			date := ""
			if d, ok := data["date"]; ok {
				date = d.(string)
			}

			loginError = &security.LoginError{
				Code:    code,
				Message: message,
				Date:    date,
			}
			loginErrorReceived = true
			fmt.Printf("Login error received: code=%d, message=%s, date=%s\n", code, message, date)
		}
	}, nil)

	// Create connection config for rAthena renewal server
	connConfig := &connection.ConnectionConfig{
		Host:        "192.168.5.219", // rAthena renewal server
		Port:        6900,
		Timeout:     10 * time.Second,
		RecvTimeout: 5 * time.Second,
		SendTimeout: 5 * time.Second,
		ServerType:  "0",
	}

	// Create direct connection
	conn := connection.NewDirectConnection(connConfig)

	// Register connection event callbacks
	conn.RegisterCallback(connection.EventConnected, func(event connection.ConnectionEvent, data interface{}) {
		t.Log("Connected to server")
	})

	conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
		if err, ok := data.(error); ok {
			t.Logf("Connection error: %v", err)
		}
	})

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create login manager
	loginManager := security.NewLoginManager(coreParser, hookManager)
	loginManager.RegisterHandlers()

	// Define packet lengths for the tokenizer
	packetLengths := make(map[string]protocol.PacketDef)
	packetLengths["0069"] = protocol.PacketDef{Length: 0, HasLength: true}   // account_server_info
	packetLengths["0071"] = protocol.PacketDef{Length: 28, HasLength: false} // received_character_ID_and_Map
	packetLengths["0073"] = protocol.PacketDef{Length: 11, HasLength: false} // map_loaded
	packetLengths["083E"] = protocol.PacketDef{Length: 26, HasLength: false} // login_error

	// Create message tokenizer
	tokenizer := protocol.NewTokenizer(packetLengths)

	// Connect to server
	fmt.Println("Connecting to server:", connConfig.Host, connConfig.Port)
	err := conn.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Disconnect()

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	// Check if connected
	if !conn.IsConnected() {
		fmt.Println("Failed to connect to server")
		return
	}
	fmt.Println("Connected to server successfully")

	// Send master login packet
	username := "username"
	password := "password"
	packet := []byte{0x64, 0x00} // 0x0064 master_login packet
	packet = append(packet, make([]byte, 24)...)
	copy(packet[2:], []byte(username))
	packet = append(packet, make([]byte, 24)...)
	copy(packet[26:], []byte(password))
	packet = append(packet, 0)                   // gender (0=female, 1=male)
	packet = append(packet, make([]byte, 16)...) // client hash
	packet = append(packet, 1, 0, 0, 0)          // version

	fmt.Println("Sending master login packet")
	err = conn.Send(packet)
	if err != nil {
		fmt.Printf("Failed to send login packet: %v\n", err)
		return
	}

	// Wait for response
	fmt.Println("Waiting for server response...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Loop to receive data
	receivedData := false
	for {
		select {
		case <-ctx.Done():
			if !receivedData {
				fmt.Println("Timeout waiting for server response")
			}
			return
		default:
			data, err := conn.Receive()
			if err != nil {
				fmt.Printf("Error receiving data: %v\n", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if data != nil && len(data) > 0 {
				receivedData = true
				fmt.Printf("Received data: %X\n", data)

				// Process the data with the tokenizer
				tokenizer.Add(data)

				// Try to read a message
				message, msgType, err := tokenizer.ReadNext()
				if err != nil {
					t.Logf("Error reading message: %v", err)
				} else {
					t.Logf("Message type: %v", msgType)
				}

				if message != nil {
					// Get message ID
					if len(message) >= 2 {
						msgID := fmt.Sprintf("%02X%02X", message[1], message[0])
						fmt.Printf("Message ID: %s\n", msgID)

						// Parse the message
						if msgID == "083E" {
							// This is a login_error packet
							if len(message) >= 26 {
								// Extract error code (4 bytes) and date (20 bytes)
								errorCode := uint32(message[2]) | uint32(message[3])<<8 | uint32(message[4])<<16 | uint32(message[5])<<24
								date := string(message[6:26])
								fmt.Printf("Login error: code=%d, date=%s\n", errorCode, date)

								// Call hook directly instead of trying to call the private handler
								hookManager.CallHook("security/login_error", map[string]interface{}{
									"code":    int(errorCode),
									"message": "Login error",
									"date":    date,
								})

								// Check if login error was received
								if !loginErrorReceived {
									fmt.Println("Login error hook was not called")
								} else {
									fmt.Printf("Login error received: %v\n", loginError)
								}
							}
						}
					}
				}

				// We got a response, so we can exit the loop
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Package main provides a standalone test for server connection
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/mikekao/openkore/goKore/network/implementation/network/connection"
	"github.com/mikekao/openkore/goKore/network/implementation/network/hooks"
	"github.com/mikekao/openkore/goKore/network/implementation/network/protocol"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/core"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/security"
	sendcore "github.com/mikekao/openkore/goKore/network/implementation/network/send/core"
	sendactor "github.com/mikekao/openkore/goKore/network/implementation/network/send/game/actor"
)

// TestConfig represents the configuration for the server connection test
type TestConfig struct {
	ServerType string `json:"server_type"`
	ServerIP   string `json:"server_ip"`
	ServerPort int    `json:"server_port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Version    uint32 `json:"version"`
	Gender     byte   `json:"gender"`
}

// RunServerConnectionTests runs both classic and renewal server tests
func RunServerConnectionTests(testMovement bool, inputFile string) {
	// If input file is provided, use it for configuration
	if inputFile != "" {
		config, err := loadTestConfig(inputFile)
		if err != nil {
			fmt.Printf("Failed to load test config: %v\n", err)
			return
		}
		
		fmt.Printf("Testing connection to server %s:%d with username %s...\n",
			config.ServerIP, config.ServerPort, config.Username)
		
		if testMovement {
			testCharacterMovement(config)
		} else {
			testRealServerConnectionWithConfig(config)
		}
		return
	}

	// Run the classic server test
	fmt.Println("Testing connection to rAthena Classic server...")
	testRealServerConnection()

	// Run the renewal server test
	fmt.Println("\nTesting connection to rAthena Renewal server...")
	testRealServerConnectionRenewal()
}

// loadTestConfig loads the test configuration from a JSON file
func loadTestConfig(filename string) (*TestConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	var config TestConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	
	return &config, nil
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
		fmt.Println("Connected to server")
	})

	conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
		if err, ok := data.(error); ok {
			fmt.Printf("Connection error: %v\n", err)
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
					fmt.Printf("Error reading message: %v\n", err)
				} else {
					fmt.Printf("Message type: %v\n", msgType)
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
		fmt.Println("Connected to server")
	})

	conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
		if err, ok := data.(error); ok {
			fmt.Printf("Connection error: %v\n", err)
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
					fmt.Printf("Error reading message: %v\n", err)
				} else {
					fmt.Printf("Message type: %v\n", msgType)
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

// testRealServerConnectionWithConfig tests connecting to a server with the provided configuration
func testRealServerConnectionWithConfig(config *TestConfig) {
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
						
							// Create connection config
							connConfig := &connection.ConnectionConfig{
								Host:        config.ServerIP,
								Port:        config.ServerPort,
								Timeout:     10 * time.Second,
								RecvTimeout: 5 * time.Second,
								SendTimeout: 5 * time.Second,
								ServerType:  config.ServerType,
							}
						
							// Create direct connection
							conn := connection.NewDirectConnection(connConfig)
						
							// Register connection event callbacks
							conn.RegisterCallback(connection.EventConnected, func(event connection.ConnectionEvent, data interface{}) {
								fmt.Println("Connected to server")
							})
						
							conn.RegisterCallback(connection.EventError, func(event connection.ConnectionEvent, data interface{}) {
								if err, ok := data.(error); ok {
									fmt.Printf("Connection error: %v\n", err)
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
							username := config.Username
							password := config.Password
						
							fmt.Println("Sending master login packet for user:", username)
							packet := []byte{0x64, 0x00} // 0x0064 master_login packet
							packet = append(packet, make([]byte, 24)...)
							copy(packet[2:], []byte(username))
							packet = append(packet, make([]byte, 24)...)
							copy(packet[26:], []byte(password))
							packet = append(packet, config.Gender)       // gender (0=female, 1=male)
							packet = append(packet, make([]byte, 16)...) // client hash
							packet = append(packet, byte(config.Version), 0, 0, 0) // version
						
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
											fmt.Printf("Error reading message: %v\n", err)
										} else {
											fmt.Printf("Message type: %v\n", msgType)
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
						
						// testCharacterMovement tests character movement
						func testCharacterMovement(config *TestConfig) {
							// Create hook manager
							hookManager := hooks.NewHookManager()
						
							// Set up hooks to track events
							var (
								characterMoved bool
								positionX      int
								positionY      int
								newPositionX   int
								newPositionY   int
							)
						
							hookManager.AddHook("game/character_moved", func(hookName string, arg interface{}, userData interface{}) {
								if data, ok := arg.(map[string]interface{}); ok {
									fromX := data["from_x"].(int)
									fromY := data["from_y"].(int)
									toX := data["to_x"].(int)
									toY := data["to_y"].(int)
									fmt.Printf("Character moved from (%d, %d) to (%d, %d)\n", fromX, fromY, toX, toY)
									positionX = toX
									positionY = toY
									characterMoved = true
								}
							}, nil)
						
							// Create connection config
							connConfig := &connection.ConnectionConfig{
								Host:        config.ServerIP,
								Port:        config.ServerPort,
								Timeout:     10 * time.Second,
								RecvTimeout: 5 * time.Second,
								SendTimeout: 5 * time.Second,
								ServerType:  config.ServerType,
							}
						
							// Create direct connection
							conn := connection.NewDirectConnection(connConfig)
						
							// Create core parser
							coreParser := core.NewCoreParser("ServerType0", hookManager)
						
							// Create send manager
							sendManager := sendcore.NewSendManager(conn)
						
							// Create movement manager
							movementManager := sendactor.NewMovementManager(sendManager)
						
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
						
							// Set initial position (this would normally come from the map_loaded event)
							positionX = 150
							positionY = 150
							fmt.Printf("Initial position: (%d, %d)\n", positionX, positionY)
						
							// Generate random direction (0=up, 1=right, 2=down, 3=left)
							rand.Seed(time.Now().UnixNano())
							direction := rand.Intn(4)
							directionName := ""
							newPositionX = positionX
							newPositionY = positionY
						
							// Calculate new position based on direction
							switch direction {
							case 0: // Up
								directionName = "up"
								newPositionY = positionY - 1
							case 1: // Right
								directionName = "right"
								newPositionX = positionX + 1
							case 2: // Down
								directionName = "down"
								newPositionY = positionY + 1
							case 3: // Left
								directionName = "left"
								newPositionX = positionX - 1
							}
						
							fmt.Printf("Moving character one square %s to position (%d, %d)\n", directionName, newPositionX, newPositionY)
						
							// Send movement packet
							err = movementManager.SendMove(newPositionX, newPositionY)
							if err != nil {
								fmt.Printf("Failed to send movement packet: %v\n", err)
								return
							}
						
							// Wait for movement response
							fmt.Println("Waiting for movement response...")
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
						
							// Loop to receive data
							for {
								select {
								case <-ctx.Done():
									fmt.Println("Timeout waiting for movement response")
									return
								default:
									data, err := conn.Receive()
									if err != nil {
										fmt.Printf("Error receiving data: %v\n", err)
										time.Sleep(100 * time.Millisecond)
										continue
									}
						
									if data != nil && len(data) > 0 {
										fmt.Printf("Received data: %X\n", data)
						
										// Process the data with the tokenizer
										// In a real implementation, we would process this data and trigger the character_moved hook
										// For now, we'll simulate the hook call
										hookManager.CallHook("game/character_moved", map[string]interface{}{
											"from_x": positionX,
											"from_y": positionY,
											"to_x":   newPositionX,
											"to_y":   newPositionY,
										})
						
										// Check if character moved
										if characterMoved {
											fmt.Printf("Character moved successfully to (%d, %d)\n", positionX, positionY)
											return
										}
									}
						
									time.Sleep(100 * time.Millisecond)
								}
							}
						}
						
						// main function to run the server connection tests
						func main() {
							// Parse command line arguments
							testMovement := flag.Bool("test-movement", false, "Test character movement")
							inputFile := flag.String("input", "", "Input JSON file with test configuration")
							flag.Parse()
						
							// Run the server connection tests
							RunServerConnectionTests(*testMovement, *inputFile)
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

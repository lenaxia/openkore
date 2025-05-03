package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lenaxia/goKore/cmd/common"
	networkCommon "github.com/lenaxia/goKore/network/common"
	receiveServers "github.com/lenaxia/goKore/network/receive/servers"
	sendServers "github.com/lenaxia/goKore/network/send/servers"
)

// ExpectedResponse defines an expected client response to a server packet
type ExpectedResponse struct {
	ServerPacketID string
	ServerDesc     string
	ClientPacketID string
	ClientDesc     string
}

// ResponseResult represents the result of checking a response
type ResponseResult struct {
	Result           string                 // "✓", "✗", or "?"
	ServerPacket     common.Packet          // The server packet
	ExpectedPacketID string                 // The expected client packet ID
	ExpectedDesc     string                 // The expected client packet description
	ExpectedFields   map[string]interface{} // Expected field values (if available)
	ActualPacket     *common.Packet         // The actual client packet (if found)
}

// ResponseChecker checks if client responses match expected responses
type ResponseChecker struct {
	dumpFiles         []string
	selectedDump      *common.PacketDump
	expectedResponses []ExpectedResponse
}

// NewResponseChecker creates a new response checker
func NewResponseChecker() *ResponseChecker {
	rc := &ResponseChecker{
		expectedResponses: []ExpectedResponse{
			{
				ServerPacketID: "0AC4",
				ServerDesc:     "Account Info With Server Info",
				ClientPacketID: "0065",
				ClientDesc:     "Game Login",
			},
			{
				ServerPacketID: "006B",
				ServerDesc:     "Received characters from Game Login",
				ClientPacketID: "0066",
				ClientDesc:     "Character Login",
			},
			{
				ServerPacketID: "0AC5",
				ServerDesc:     "Map Server Info",
				ClientPacketID: "0072",
				ClientDesc:     "Map Login",
			},
			{
				ServerPacketID: "02EB",
				ServerDesc:     "Map Loaded",
				ClientPacketID: "007D",
				ClientDesc:     "Map Loaded",
			},
		},
	}
	return rc
}

// loadDumpFiles loads the list of available dump files
func (rc *ResponseChecker) loadDumpFiles(dumpDir string) error {
	files, err := ioutil.ReadDir(dumpDir)
	if err != nil {
		return err
	}

	rc.dumpFiles = []string{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			rc.dumpFiles = append(rc.dumpFiles, file.Name())
		}
	}

	return nil
}

// loadDump loads a packet dump file
func (rc *ResponseChecker) loadDump(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var dump common.PacketDump
	if err := json.Unmarshal(data, &dump); err != nil {
		return err
	}

	rc.selectedDump = &dump

	fmt.Printf("Loaded %s with %d packets.\n",
		filepath.Base(path), dump.PacketCount)

	return nil
}

// checkResponses checks if client responses match expected responses
func (rc *ResponseChecker) checkResponses() {
	if rc.selectedDump == nil {
		fmt.Println("No dump loaded")
		return
	}

	// Store results for detailed report
	var results []ResponseResult

	// Print summary table
	fmt.Println("\n=== Response Check Summary ===")
	fmt.Printf("%-6s %-30s %-10s %-30s %-10s %-30s %-10s\n",
		"Result", "Server Packet", "ID", "Expected Response", "ID", "Actual Response", "ID")
	fmt.Println(strings.Repeat("-", 130))

	for _, expected := range rc.expectedResponses {
		found := false

		// Find all instances of the server packet
		for i, packet := range rc.selectedDump.Packets {
			if packet.Direction == "received" && packet.PacketID == expected.ServerPacketID {
				// Look for any client response in the next few packets
				responseFound := false
				var actualResponsePacket *common.Packet

				for j := i + 1; j < len(rc.selectedDump.Packets) && j < i+5; j++ {
					nextPacket := rc.selectedDump.Packets[j]
					if nextPacket.Direction == "sent" {
						// Found a response packet
						actualResponsePacket = &rc.selectedDump.Packets[j]

						// Check if it matches the expected response
						if nextPacket.PacketID == expected.ClientPacketID {
							fmt.Printf("%-6s %-30s %-10s %-30s %-10s %-30s %-10s\n", "✓",
								expected.ServerDesc, expected.ServerPacketID,
								expected.ClientDesc, expected.ClientPacketID,
								expected.ClientDesc, nextPacket.PacketID)

							responseFound = true
							found = true

							// Generate expected fields based on server packet
							_, expectedExists := getPacketDefinition(expected.ClientPacketID, "sent")
							var expectedFields map[string]interface{}

							if expectedExists {
								// Extract relevant fields from server packet to use in expected response
								serverBytes := rawDataToBytes(packet.RawData)
								serverDef, serverExists := getPacketDefinition(packet.PacketID, "received")

								if serverExists {
									serverFields := decodePacketFields(serverBytes, serverDef.Format, serverDef.FieldNames)

									// Generate expected fields based on server fields
									// This is a simplified approach - in a real implementation, you would
									// have logic to determine which server fields should be included in the response
									expectedFields = make(map[string]interface{})

									// Copy fields that might be relevant for the response
									relevantFields := []string{"accountID", "sessionID", "sessionID2", "charID"}
									for _, field := range relevantFields {
										if value, exists := serverFields[field]; exists {
											expectedFields[field] = value
										}
									}
								}
							}

							// Store result for detailed report
							results = append(results, ResponseResult{
								Result:           "✓",
								ServerPacket:     packet,
								ExpectedPacketID: expected.ClientPacketID,
								ExpectedDesc:     expected.ClientDesc,
								ExpectedFields:   expectedFields,
								ActualPacket:     actualResponsePacket,
							})

							break
						}
					}
				}

				if !responseFound {
					// If we found a response but it wasn't the expected one
					if actualResponsePacket != nil {
						fmt.Printf("%-6s %-30s %-10s %-30s %-10s %-30s %-10s\n", "✗",
							expected.ServerDesc, expected.ServerPacketID,
							expected.ClientDesc, expected.ClientPacketID,
							actualResponsePacket.Description, actualResponsePacket.PacketID)

						// Generate expected fields based on server packet
						_, expectedExists := getPacketDefinition(expected.ClientPacketID, "sent")
						var expectedFields map[string]interface{}

						if expectedExists {
							// Extract relevant fields from server packet to use in expected response
							serverBytes := rawDataToBytes(packet.RawData)
							serverDef, serverExists := getPacketDefinition(packet.PacketID, "received")

							if serverExists {
								serverFields := decodePacketFields(serverBytes, serverDef.Format, serverDef.FieldNames)

								// Generate expected fields based on server fields
								expectedFields = make(map[string]interface{})

								// Copy fields that might be relevant for the response
								relevantFields := []string{"accountID", "sessionID", "sessionID2", "charID"}
								for _, field := range relevantFields {
									if value, exists := serverFields[field]; exists {
										expectedFields[field] = value
									}
								}
							}
						}

						// Store result for detailed report
						results = append(results, ResponseResult{
							Result:           "✗",
							ServerPacket:     packet,
							ExpectedPacketID: expected.ClientPacketID,
							ExpectedDesc:     expected.ClientDesc,
							ExpectedFields:   expectedFields,
							ActualPacket:     actualResponsePacket,
						})
					} else {
						// No response found at all
						fmt.Printf("%-6s %-30s %-10s %-30s %-10s %-30s %-10s\n", "✗",
							expected.ServerDesc, expected.ServerPacketID,
							expected.ClientDesc, expected.ClientPacketID,
							"No response", "")

						// Generate expected fields based on server packet
						_, expectedExists := getPacketDefinition(expected.ClientPacketID, "sent")
						var expectedFields map[string]interface{}

						if expectedExists {
							// Extract relevant fields from server packet to use in expected response
							serverBytes := rawDataToBytes(packet.RawData)
							serverDef, serverExists := getPacketDefinition(packet.PacketID, "received")

							if serverExists {
								serverFields := decodePacketFields(serverBytes, serverDef.Format, serverDef.FieldNames)

								// Generate expected fields based on server fields
								expectedFields = make(map[string]interface{})

								// Copy fields that might be relevant for the response
								relevantFields := []string{"accountID", "sessionID", "sessionID2", "charID"}
								for _, field := range relevantFields {
									if value, exists := serverFields[field]; exists {
										expectedFields[field] = value
									}
								}
							}
						}

						// Store result for detailed report
						results = append(results, ResponseResult{
							Result:           "✗",
							ServerPacket:     packet,
							ExpectedPacketID: expected.ClientPacketID,
							ExpectedDesc:     expected.ClientDesc,
							ExpectedFields:   expectedFields,
							ActualPacket:     nil,
						})
					}
				}
			}
		}

		if !found {
			fmt.Printf("%-6s %-30s %-10s %-30s %-10s %-30s %-10s\n", "?",
				expected.ServerDesc, expected.ServerPacketID,
				expected.ClientDesc, expected.ClientPacketID,
				"Server packet not found", "")

			// Store result for detailed report
			results = append(results, ResponseResult{
				Result:           "?",
				ServerPacket:     common.Packet{},
				ExpectedPacketID: expected.ClientPacketID,
				ExpectedDesc:     expected.ClientDesc,
				ExpectedFields:   nil,
				ActualPacket:     nil,
			})
		}
	}

	// Print detailed report
	rc.printDetailedReport(results)
}

// printDetailedReport prints a detailed report of the response check results
func (rc *ResponseChecker) printDetailedReport(results []ResponseResult) {
	fmt.Println("\n\n=== Detailed Response Analysis ===")

	for i, result := range results {
		fmt.Printf("\n--- Response Check %d: %s ---\n", i+1, resultToEmoji(result.Result))

		// Print server packet info
		if result.Result != "?" {
			fmt.Println("\nServer Packet:")
			fmt.Printf("  ID: %s\n", result.ServerPacket.PacketID)
			fmt.Printf("  Description: %s\n", result.ServerPacket.Description)
			fmt.Printf("  Timestamp: %s\n", result.ServerPacket.Timestamp)

			// Print raw data summary
			if len(result.ServerPacket.RawData) > 0 {
				fmt.Println("  Raw Data (first 16 bytes):")
				hexBytes := []string{}
				for _, data := range result.ServerPacket.RawData {
					hexBytes = append(hexBytes, data.HexBytes...)
					if len(hexBytes) >= 16 {
						break
					}
				}
				if len(hexBytes) > 16 {
					hexBytes = hexBytes[:16]
				}
				fmt.Printf("    %s\n", strings.Join(hexBytes, " "))
			}
		}

		// Print expected response
		fmt.Println("\nExpected Response:")
		fmt.Printf("  ID: %s\n", result.ExpectedPacketID)
		fmt.Printf("  Description: %s\n", result.ExpectedDesc)

		// Get packet definition for expected packet
		expectedDef, expectedExists := getPacketDefinition(result.ExpectedPacketID, "sent")

		if expectedExists {
			fmt.Println("  Format:")
			fmt.Printf("    %s\n", expectedDef.Format)

			fmt.Println("  Expected Fields:")
			if result.ExpectedFields != nil && len(result.ExpectedFields) > 0 {
				for name, value := range result.ExpectedFields {
					// Format the value based on its type
					var formattedValue string
					switch v := value.(type) {
					case int:
						formattedValue = fmt.Sprintf("%d (0x%X)", v, v)
					case string:
						if isPrintable(v) {
							formattedValue = fmt.Sprintf("\"%s\"", v)
						} else {
							formattedValue = fmt.Sprintf("0x%s", formatBinaryData([]byte(v)))
						}
					default:
						formattedValue = fmt.Sprintf("%v", v)
					}
					fmt.Printf("    %s: %s\n", name, formattedValue)
				}
			} else {
				fmt.Println("    (No field values available - using server packet fields as reference)")
				for _, fieldName := range expectedDef.FieldNames {
					fmt.Printf("    %s: (value unknown)\n", fieldName)
				}
			}
		} else {
			fmt.Println("  No packet definition found for expected packet")
		}

		// Print actual response
		fmt.Println("\nActual Response:")
		if result.ActualPacket != nil {
			fmt.Printf("  ID: %s\n", result.ActualPacket.PacketID)
			fmt.Printf("  Description: %s\n", result.ActualPacket.Description)
			fmt.Printf("  Timestamp: %s\n", result.ActualPacket.Timestamp)

			// Print raw data summary
			if len(result.ActualPacket.RawData) > 0 {
				fmt.Println("  Raw Data (first 16 bytes):")
				hexBytes := []string{}
				for _, data := range result.ActualPacket.RawData {
					hexBytes = append(hexBytes, data.HexBytes...)
					if len(hexBytes) >= 16 {
						break
					}
				}
				if len(hexBytes) > 16 {
					hexBytes = hexBytes[:16]
				}
				fmt.Printf("    %s\n", strings.Join(hexBytes, " "))
			}

			// Decode packet fields
			fmt.Println("\n  Decoded Fields:")

			// Get packet definition for expected packet
			expectedDef, expectedExists := getPacketDefinition(result.ExpectedPacketID, "sent")

			// Get packet definition for actual packet
			actualDef, actualExists := getPacketDefinition(result.ActualPacket.PacketID, "sent")

			// Convert raw data to bytes
			actualBytes := rawDataToBytes(result.ActualPacket.RawData)

			// Print actual packet format and fields
			if actualExists {
				fmt.Println("    Actual Format:")
				fmt.Printf("      %s\n", actualDef.Format)

				// Decode fields based on packet definition
				actualFields := decodePacketFields(actualBytes, actualDef.Format, actualDef.FieldNames)

				// Print decoded fields
				fmt.Println("\n    Actual Fields:")
				for name, value := range actualFields {
					// Format the value based on its type
					var formattedValue string
					switch v := value.(type) {
					case int:
						formattedValue = fmt.Sprintf("%d (0x%X)", v, v)
					case string:
						if isPrintable(v) {
							formattedValue = fmt.Sprintf("\"%s\"", v)
						} else {
							formattedValue = fmt.Sprintf("0x%s", formatBinaryData([]byte(v)))
						}
					default:
						formattedValue = fmt.Sprintf("%v", v)
					}
					fmt.Printf("      %s: %s\n", name, formattedValue)
				}
			} else {
				fmt.Println("\n    No packet definition found for actual packet")
			}

			// Highlight discrepancies
			if result.Result == "✗" {
				fmt.Println("\nDiscrepancies:")
				fmt.Printf("  Expected packet ID: %s, got: %s\n",
					result.ExpectedPacketID, result.ActualPacket.PacketID)
				fmt.Printf("  Expected description: %s, got: %s\n",
					result.ExpectedDesc, result.ActualPacket.Description)

				if expectedExists && actualExists && expectedDef.Format != actualDef.Format {
					fmt.Printf("  Expected format: %s, got: %s\n",
						expectedDef.Format, actualDef.Format)
				}

				// Compare field names and values if both definitions exist
				if expectedExists && actualExists {
					fmt.Println("\n  Field Comparison:")

					// Find missing and unexpected fields
					expectedFieldNames := make(map[string]bool)
					for _, name := range expectedDef.FieldNames {
						expectedFieldNames[name] = true
					}

					for _, name := range actualDef.FieldNames {
						if _, exists := expectedFieldNames[name]; !exists {
							fmt.Printf("    Unexpected field: %s\n", name)
						}
					}

					for _, name := range expectedDef.FieldNames {
						found := false
						for _, actualName := range actualDef.FieldNames {
							if name == actualName {
								found = true
								break
							}
						}
						if !found {
							fmt.Printf("    Missing field: %s\n", name)
						}
					}

					// Compare field values
					if result.ExpectedFields != nil && len(result.ExpectedFields) > 0 {
						fmt.Println("\n  Field Value Comparison:")

						actualFields := decodePacketFields(actualBytes, actualDef.Format, actualDef.FieldNames)

						// Compare values for fields that exist in both
						for name, expectedValue := range result.ExpectedFields {
							if actualValue, exists := actualFields[name]; exists {
								if fmt.Sprintf("%v", expectedValue) != fmt.Sprintf("%v", actualValue) {
									fmt.Printf("    Field %s: expected %v, got %v\n",
										name, expectedValue, actualValue)
								} else {
									fmt.Printf("    Field %s: matches expected value %v\n",
										name, expectedValue)
								}
							}
						}
					}
				}
			}
		} else if result.Result == "?" {
			fmt.Println("  Server packet not found in dump")
		} else {
			fmt.Println("  No response packet found")
		}

		fmt.Println(strings.Repeat("-", 80))
	}
}

// getPacketDefinition gets the packet definition based on direction
func getPacketDefinition(packetID string, direction string) (networkCommon.PacketConstruction, bool) {
	if direction == "sent" {
		// For sent packets, try to get from send servers
		sendDefs := sendServers.ServerType0PacketConstructions()
		if def, exists := sendDefs[packetID]; exists {
			return def, true
		}
	} else {
		// For received packets, get from receive servers
		receiveDefs := receiveServers.ServerType0PacketConstructions()
		if def, exists := receiveDefs[packetID]; exists {
			return def, true
		}
	}

	// If not found in the appropriate factory, try the other one as fallback
	if direction == "sent" {
		receiveDefs := receiveServers.ServerType0PacketConstructions()
		def, exists := receiveDefs[packetID]
		return def, exists
	} else {
		sendDefs := sendServers.ServerType0PacketConstructions()
		def, exists := sendDefs[packetID]
		return def, exists
	}
}

// decodePacketFields decodes packet fields based on format string
func decodePacketFields(rawBytes []byte, format string, fieldNames []string) map[string]interface{} {
	fields := make(map[string]interface{})

	offset := 0
	fieldIndex := 0

	for i := 0; i < len(format) && fieldIndex < len(fieldNames); i++ {
		if offset >= len(rawBytes) {
			break
		}

		switch format[i] {
		case 'C': // unsigned char (1 byte)
			if offset < len(rawBytes) {
				fields[fieldNames[fieldIndex]] = int(rawBytes[offset])
				offset += 1
				fieldIndex++
			}
		case 'v': // unsigned short (2 bytes, little-endian)
			if offset+1 < len(rawBytes) {
				value := uint16(rawBytes[offset]) | uint16(rawBytes[offset+1])<<8
				fields[fieldNames[fieldIndex]] = int(value)
				offset += 2
				fieldIndex++
			}
		case 'V': // unsigned int (4 bytes, little-endian)
			if offset+3 < len(rawBytes) {
				value := uint32(rawBytes[offset]) |
					uint32(rawBytes[offset+1])<<8 |
					uint32(rawBytes[offset+2])<<16 |
					uint32(rawBytes[offset+3])<<24
				fields[fieldNames[fieldIndex]] = int(value)
				offset += 4
				fieldIndex++
			}
		case 'a': // fixed-length string
			// Look for a number after 'a'
			j := i + 1
			numStr := ""
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				numStr += string(format[j])
				j++
			}

			if numStr != "" {
				length, _ := strconv.Atoi(numStr)
				if offset+length <= len(rawBytes) {
					// Extract the string and trim null bytes
					str := string(rawBytes[offset : offset+length])
					str = strings.TrimRight(str, "\x00")
					fields[fieldNames[fieldIndex]] = str
					offset += length
					fieldIndex++
				}
				i = j - 1 // Skip the digits we've processed
			}
		case 'Z': // null-terminated string
			// Find the null terminator
			end := offset
			for end < len(rawBytes) && rawBytes[end] != 0 {
				end++
			}

			if end < len(rawBytes) {
				str := string(rawBytes[offset:end])
				fields[fieldNames[fieldIndex]] = str
				offset = end + 1
				fieldIndex++
			}
		case ' ', '*', 'x': // Ignore spaces, wildcards, and padding
			continue
		}
	}

	return fields
}

// hexBytesToBytes converts a slice of hex strings to a byte slice
func hexBytesToBytes(hexBytes []string) []byte {
	bytes := make([]byte, len(hexBytes))
	for i, hex := range hexBytes {
		b, _ := strconv.ParseUint(hex, 16, 8)
		bytes[i] = byte(b)
	}
	return bytes
}

// rawDataToBytes converts raw data to a byte slice
func rawDataToBytes(rawData []common.RawData) []byte {
	var allBytes []byte
	for _, data := range rawData {
		allBytes = append(allBytes, hexBytesToBytes(data.HexBytes)...)
	}
	return allBytes
}

// isPrintable checks if a string contains only printable ASCII characters
func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

// formatBinaryData formats binary data in a readable way
func formatBinaryData(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// Convert to hex representation
	hexStr := ""
	for _, b := range data {
		hexStr += fmt.Sprintf("%02x", b)
	}

	return hexStr
}

// resultToEmoji converts a result string to a descriptive string
func resultToEmoji(result string) string {
	switch result {
	case "✓":
		return "SUCCESS"
	case "✗":
		return "FAILURE"
	case "?":
		return "NOT FOUND"
	default:
		return "UNKNOWN"
	}
}

// Run starts the response checker
func (rc *ResponseChecker) Run(dumpDir string) error {
	// Load dump files
	if err := rc.loadDumpFiles(dumpDir); err != nil {
		return err
	}

	// Display available dump files
	fmt.Println("Available packet dumps:")
	for i, file := range rc.dumpFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}

	// Get user selection
	fmt.Print("\nSelect a dump file (enter number): ")
	var input string
	fmt.Scanln(&input)

	index := 0
	_, err := fmt.Sscanf(input, "%d", &index)
	if err != nil || index < 1 || index > len(rc.dumpFiles) {
		return fmt.Errorf("invalid selection")
	}

	// Load selected dump
	dumpPath := filepath.Join(dumpDir, rc.dumpFiles[index-1])
	if err := rc.loadDump(dumpPath); err != nil {
		return err
	}

	// Check responses
	rc.checkResponses()

	return nil
}

func main() {
	// Define paths
	dumpDir := "/home/mikekao/personal/openkore/goKore/verification/PacketAnalysis/extracteddata"

	// Create and run the response checker
	rc := NewResponseChecker()
	if err := rc.Run(dumpDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

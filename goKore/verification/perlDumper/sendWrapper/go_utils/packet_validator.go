package packet_validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// PacketData represents the structure of the JSON test data files
type PacketData struct {
	Method  string        `json:"method"`
	Args    []interface{} `json:"args"`
	Packets []struct {
		Hex       string `json:"hex"`
		MessageID string `json:"messageID"`
		Bytes     []int  `json:"bytes"`
	} `json:"packets"`
}

// ValidationResult contains the result of a packet validation
type ValidationResult struct {
	IsValid     bool
	Errors      []string
	ExpectedHex string
	ActualHex   string
}

// LoadPacketData loads test data from a JSON file
func LoadPacketData(filePath string) (*PacketData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var packetData PacketData
	err = json.Unmarshal(data, &packetData)
	if err != nil {
		return nil, err
	}

	return &packetData, nil
}

// LoadTestData loads all test data from a directory
func LoadTestData(dirPath string) (map[string]*PacketData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	testData := make(map[string]*PacketData)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			methodName := strings.TrimSuffix(file.Name(), ".json")
			data, err := LoadPacketData(filepath.Join(dirPath, file.Name()))
			if err != nil {
				fmt.Printf("Warning: Failed to load %s: %v\n", file.Name(), err)
				continue
			}
			testData[methodName] = data
		}
	}

	return testData, nil
}

// BytesToHex converts a byte slice to a hex string
func BytesToHex(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

// ValidatePacket validates a generated packet against the expected packet
func ValidatePacket(methodName string, args []interface{}, actualPacket []byte, expectedData *PacketData) ValidationResult {
	result := ValidationResult{
		IsValid: true,
		Errors:  []string{},
	}

	// Check if there are any expected packets
	if len(expectedData.Packets) == 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, "No expected packets in test data")
		return result
	}

	// Convert actual packet to hex
	actualHex := BytesToHex(actualPacket)
	expectedHex := expectedData.Packets[0].Hex

	result.ActualHex = actualHex
	result.ExpectedHex = expectedHex

	// Compare hex strings
	if actualHex != expectedHex {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Hex mismatch: expected %s, got %s", expectedHex, actualHex))
	}

	// Compare packet length
	if len(actualPacket) != len(expectedData.Packets[0].Bytes) {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Length mismatch: expected %d bytes, got %d bytes",
			len(expectedData.Packets[0].Bytes), len(actualPacket)))
	}

	// Compare packet bytes
	for i := 0; i < len(actualPacket) && i < len(expectedData.Packets[0].Bytes); i++ {
		if int(actualPacket[i]) != expectedData.Packets[0].Bytes[i] {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Byte mismatch at position %d: expected 0x%02x, got 0x%02x",
				i, expectedData.Packets[0].Bytes[i], actualPacket[i]))
		}
	}

	// Check message ID
	if len(actualPacket) >= 2 {
		actualMessageID := fmt.Sprintf("%02X%02X", actualPacket[1], actualPacket[0])
		if actualMessageID != expectedData.Packets[0].MessageID {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Message ID mismatch: expected %s, got %s",
				expectedData.Packets[0].MessageID, actualMessageID))
		}
	}

	return result
}

// ValidateAllPackets validates all packets in a test directory against an implementation
type PacketGenerator func(methodName string, args []interface{}) ([]byte, error)

func ValidateAllPackets(testDataDir string, generator PacketGenerator) (map[string]ValidationResult, error) {
	testData, err := LoadTestData(testDataDir)
	if err != nil {
		return nil, err
	}

	results := make(map[string]ValidationResult)
	for methodName, data := range testData {
		packet, err := generator(methodName, data.Args)
		if err != nil {
			results[methodName] = ValidationResult{
				IsValid: false,
				Errors:  []string{fmt.Sprintf("Error generating packet: %v", err)},
			}
			continue
		}

		results[methodName] = ValidatePacket(methodName, data.Args, packet, data)
	}

	return results, nil
}

// PrintValidationResults prints the results of packet validation
func PrintValidationResults(results map[string]ValidationResult) {
	validCount := 0
	invalidCount := 0

	for method, result := range results {
		if result.IsValid {
			validCount++
			fmt.Printf("✅ %s: Valid\n", method)
		} else {
			invalidCount++
			fmt.Printf("❌ %s: Invalid\n", method)
			for _, err := range result.Errors {
				fmt.Printf("   - %s\n", err)
			}
		}
	}

	fmt.Printf("\nSummary: %d valid, %d invalid\n", validCount, invalidCount)
}

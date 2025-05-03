package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PacketData represents the structure of our validation data JSON files
type PacketData struct {
	Method  string        `json:"method"`
	Args    []interface{} `json:"args"`
	Packets []struct {
		Hex       string `json:"hex"`
		MessageID string `json:"messageID"`
		Bytes     []int  `json:"bytes"`
	} `json:"packets"`
}

func main() {
	fmt.Println("OpenKore Send.pm Validation Data Example")
	fmt.Println("========================================")

	// Directory containing validation data
	validationDir := "../validation_data"

	// Check if validation directory exists
	if _, err := os.Stat(validationDir); os.IsNotExist(err) {
		fmt.Printf("Validation directory '%s' does not exist.\n", validationDir)
		fmt.Println("Please run the generate_validation_data.sh script first.")
		return
	}

	// List all JSON files in the validation directory
	files, err := os.ReadDir(validationDir)
	if err != nil {
		fmt.Printf("Error reading validation directory: %v\n", err)
		return
	}

	// Process each validation file
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		fmt.Printf("\nProcessing %s...\n", file.Name())

		// Load validation data
		data, err := loadValidationData(filepath.Join(validationDir, file.Name()))
		if err != nil {
			fmt.Printf("Error loading validation data: %v\n", err)
			continue
		}

		// Print information about the validation data
		fmt.Printf("  Method: %s\n", data.Method)
		fmt.Printf("  Args: %v\n", data.Args)
		fmt.Printf("  Number of packets: %d\n", len(data.Packets))

		// Process each packet
		for i, packet := range data.Packets {
			fmt.Printf("  Packet %d:\n", i)
			fmt.Printf("    Message ID: %s\n", packet.MessageID)
			fmt.Printf("    Hex: %s\n", packet.Hex)

			// Convert bytes array to byte slice
			byteSlice := make([]byte, len(packet.Bytes))
			for j, b := range packet.Bytes {
				byteSlice[j] = byte(b)
			}

			// Compare hex representation with bytes
			hexFromBytes := hex.EncodeToString(byteSlice)
			if hexFromBytes == packet.Hex {
				fmt.Printf("    Hex matches byte array ✓\n")
			} else {
				fmt.Printf("    Hex doesn't match byte array ✗\n")
				fmt.Printf("    Hex from bytes: %s\n", hexFromBytes)
			}
		}
	}
}

// loadValidationData loads a validation data file
func loadValidationData(filePath string) (*PacketData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data PacketData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

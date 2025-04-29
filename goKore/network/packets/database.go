package packets

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadPacketDefinitionsFromFile loads packet definitions from a file
func LoadPacketDefinitionsFromFile(filePath string) (*PacketDatabase, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open packet definitions file: %w", err)
	}
	defer file.Close()

	db := NewPacketDatabase()
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line
		parts := strings.Fields(line)
		if len(parts) < 3 {
			// Skip lines with insufficient fields
			fmt.Printf("Warning: Line %d has insufficient fields, skipping: %s\n", lineNum, line)
			continue
		}

		// Extract packet ID and name
		packetID := parts[0]
		packetName := parts[1]

		// Extract format string
		formatIndex := 2
		format := parts[formatIndex]

		// Check if the format is valid
		if !isValidFormat(format) {
			fmt.Printf("Warning: Line %d has invalid format string, skipping: %s\n", lineNum, line)
			continue
		}

		// Extract parameter names
		paramNames := make([]string, 0, len(parts)-formatIndex-1)
		for i := formatIndex + 1; i < len(parts); i++ {
			paramNames = append(paramNames, parts[i])
		}

		// Create and add the packet definition
		def := NewPacketDefinition(packetID, packetName, format, paramNames)
		db.AddPacketDefinition(def)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading packet definitions file: %w", err)
	}

	return db, nil
}

// isValidFormat checks if a format string is valid
func isValidFormat(format string) bool {
	// Simple validation: check if the format string contains valid type specifiers
	validTypes := []string{"C", "v", "V", "a", "Z", "x"}

	// Empty format is valid for packets with no data
	if format == "" {
		return true
	}

	// Check each part of the format string
	parts := strings.Split(format, " ")
	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check if the part starts with a valid type
		valid := false
		for _, validType := range validTypes {
			if strings.HasPrefix(part, validType) {
				valid = true
				break
			}
		}

		if !valid {
			return false
		}
	}

	return true
}

// MergePacketDatabases merges two packet databases
// If there are duplicate packet IDs, the second database takes precedence
func MergePacketDatabases(db1, db2 *PacketDatabase) *PacketDatabase {
	merged := NewPacketDatabase()

	// Add all packet definitions from the first database
	for _, def := range db1.packetsByID {
		merged.AddPacketDefinition(def)
	}

	// Add all packet definitions from the second database, overriding any duplicates
	for _, def := range db2.packetsByID {
		merged.AddPacketDefinition(def)
	}

	return merged
}

// LoadServerSpecificPacketDefinitions loads base packet definitions and then
// overlays server-specific packet definitions
func LoadServerSpecificPacketDefinitions(baseFilePath, serverType, baseDir string) (*PacketDatabase, error) {
	// Load base packet definitions
	baseDB, err := LoadPacketDefinitionsFromFile(baseFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load base packet definitions: %w", err)
	}

	// Construct the path to the server-specific packet definitions
	serverFilePath := filepath.Join(baseDir, serverType, "packets.txt")

	// Check if the server-specific file exists
	if _, err := os.Stat(serverFilePath); os.IsNotExist(err) {
		// If the file doesn't exist, just return the base definitions
		return baseDB, nil
	}

	// Load server-specific packet definitions
	serverDB, err := LoadPacketDefinitionsFromFile(serverFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load server-specific packet definitions: %w", err)
	}

	// Merge the databases, with server-specific definitions taking precedence
	return MergePacketDatabases(baseDB, serverDB), nil
}

// LoadPacketDefinitionsFromDirectory loads all packet definition files from a directory
func LoadPacketDefinitionsFromDirectory(dirPath string) (*PacketDatabase, error) {
	db := NewPacketDatabase()

	// Walk through the directory and load all .txt files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process .txt files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".txt") {
			return nil
		}

		// Load packet definitions from the file
		fileDB, err := LoadPacketDefinitionsFromFile(path)
		if err != nil {
			fmt.Printf("Warning: Failed to load packet definitions from %s: %v\n", path, err)
			return nil
		}

		// Merge with the main database
		db = MergePacketDatabases(db, fileDB)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return db, nil
}

// SavePacketDefinitionsToFile saves packet definitions to a file
func SavePacketDefinitionsToFile(db *PacketDatabase, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create packet definitions file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header
	_, err = writer.WriteString("# Packet definitions\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = writer.WriteString("# Format: PacketID, PacketName, PacketFormat, ParamName1, ParamName2, ...\n\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write packet definitions
	for id, def := range db.packetsByID {
		line := fmt.Sprintf("%s %s %s", id, def.Name, def.Format)
		for _, paramName := range def.ParamNames {
			line += " " + paramName
		}
		line += "\n"

		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write packet definition: %w", err)
		}
	}

	return nil
}

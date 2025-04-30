// Package main provides a tool to verify packet implementations
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// PacketInfo holds information about a packet definition
type PacketInfo struct {
	Code    string
	Name    string
	IsXMark bool // Whether it's marked with [X] or [ ]
	Line    string
}

func main() {
	// Define command-line flags
	packetsFilePathPtr := flag.String("packets", "../../network/send/servers/servertype0.packets", "Path to the packets file")
	goFilePathPtr := flag.String("go", "../../network/send/servers/servertype0.go", "Path to the Go file")
	serverTypePtr := flag.String("server-type", "", "Server type (e.g., servertype0, servertype1)")

	// Parse command-line flags
	flag.Parse()

	// Determine file paths based on arguments
	var packetsFilePath, goFilePath string

	if *serverTypePtr != "" {
		// If server type is provided, use it to construct file paths
		serverType := *serverTypePtr
		packetsFilePath = fmt.Sprintf("../../network/send/servers/%s.packets", serverType)
		goFilePath = fmt.Sprintf("../../network/send/servers/%s.go", serverType)
	} else {
		// Otherwise use the explicitly provided paths
		packetsFilePath = *packetsFilePathPtr
		goFilePath = *goFilePathPtr
	}

	// Check if files exist
	if _, err := os.Stat(packetsFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: Packets file not found: %s\n", packetsFilePath)
		return
	}

	if _, err := os.Stat(goFilePath); os.IsNotExist(err) {
		fmt.Printf("Error: Go file not found: %s\n", goFilePath)
		return
	}

	// Extract packet codes from packets file
	packetsFilePackets, err := extractPacketsFromPacketsFile(packetsFilePath)
	if err != nil {
		fmt.Printf("Error extracting packets from %s: %v\n", packetsFilePath, err)
		return
	}

	// Extract packet codes from Go file
	goFilePackets, err := extractPacketsFromGoFile(goFilePath)
	if err != nil {
		fmt.Printf("Error extracting packets from %s: %v\n", goFilePath, err)
		return
	}

	// Create maps for quick lookup
	goFilePacketsMap := make(map[string]bool)
	for _, packet := range goFilePackets {
		goFilePacketsMap[packet] = true
	}

	// Analyze packets
	var missingPackets []PacketInfo
	var implementedPackets []PacketInfo
	var unmarkedPackets []PacketInfo
	var commentedPackets []PacketInfo

	for _, packetInfo := range packetsFilePackets {
		if strings.HasPrefix(strings.TrimSpace(packetInfo.Line), "#") {
			commentedPackets = append(commentedPackets, packetInfo)
			continue
		}

		if goFilePacketsMap[packetInfo.Code] {
			implementedPackets = append(implementedPackets, packetInfo)
			if !packetInfo.IsXMark {
				unmarkedPackets = append(unmarkedPackets, packetInfo)
			}
		} else {
			missingPackets = append(missingPackets, packetInfo)
		}
	}

	// Print results
	fmt.Printf("Total packets in %s: %d\n", packetsFilePath, len(packetsFilePackets))
	fmt.Printf("Total packets in %s: %d\n", goFilePath, len(goFilePackets))
	fmt.Printf("Commented packets: %d\n", len(commentedPackets))
	fmt.Printf("Implemented packets: %d\n", len(implementedPackets))
	fmt.Printf("Unmarked but implemented packets: %d\n", len(unmarkedPackets))

	if len(missingPackets) == 0 {
		fmt.Println("All non-commented packets from the packets file are implemented in the Go file!")
	} else {
		fmt.Printf("Missing packets in %s: %d\n", goFilePath, len(missingPackets))
		fmt.Println("Missing packets (excluding commented ones):")

		// Group by X mark status
		var missingXMarked []PacketInfo
		var missingUnmarked []PacketInfo

		for _, packet := range missingPackets {
			if packet.IsXMark {
				missingXMarked = append(missingXMarked, packet)
			} else {
				missingUnmarked = append(missingUnmarked, packet)
			}
		}

		fmt.Printf("Missing packets marked with [X]: %d\n", len(missingXMarked))
		if len(missingXMarked) > 0 {
			fmt.Println("These should be implemented but are missing:")
			for i, packet := range missingXMarked {
				if i < 20 {
					fmt.Printf("%s: %s\n", packet.Code, packet.Name)
				} else if i == 20 {
					fmt.Println("... (more packets omitted for brevity)")
					break
				}
			}
		}

		fmt.Printf("\nMissing packets marked with [ ]: %d\n", len(missingUnmarked))
		if len(missingUnmarked) > 0 {
			fmt.Println("These are not yet marked as implemented:")
			for i, packet := range missingUnmarked {
				if i < 20 {
					fmt.Printf("%s: %s\n", packet.Code, packet.Name)
				} else if i == 20 {
					fmt.Println("... (more packets omitted for brevity)")
					break
				}
			}
		}
	}

	if len(unmarkedPackets) > 0 {
		fmt.Println("\nPackets that are implemented but not marked with [X]:")
		for _, packet := range unmarkedPackets {
			fmt.Printf("%s: %s\n", packet.Code, packet.Name)
		}
	}
}

func extractPacketsFromPacketsFile(filePath string) ([]PacketInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var packets []PacketInfo
	scanner := bufio.NewScanner(file)

	// Regular expression to match packet codes and names
	re := regexp.MustCompile(`(?:\[([X ])\]\s*)?(?:#\s*)?'([0-9A-F]{4})'\s*=>\s*\['([^']+)'`)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip commented lines
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) > 3 {
			isXMark := matches[1] == "X"
			code := matches[2]
			name := matches[3]

			packets = append(packets, PacketInfo{
				Code:    code,
				Name:    name,
				IsXMark: isXMark,
				Line:    line,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return packets, nil
}

func extractPacketsFromGoFile(filePath string) ([]string, error) {
	// Read the entire file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Regular expression to match packet codes
	re := regexp.MustCompile(`packets\["([0-9A-F]{4})"\]`)

	// Find all matches
	matches := re.FindAllSubmatch(content, -1)

	// Extract packet codes
	var packets []string
	for _, match := range matches {
		if len(match) > 1 {
			packets = append(packets, string(match[1]))
		}
	}

	return packets, nil
}

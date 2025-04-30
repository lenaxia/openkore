package main

import (
	"bufio"
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
	LineNum int
}

func main() {
	packetsFilePath := "../network/send/servers/servertype0.packets"
	goFilePath := "../network/send/servers/servertype0.go"

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

	// Create a map of packet info by code for quick lookup
	packetInfoMap := make(map[string]PacketInfo)
	for _, packetInfo := range packetsFilePackets {
		packetInfoMap[packetInfo.Code] = packetInfo
	}

	// Analyze packets
	var incorrectlyMarkedAsImplemented []PacketInfo
	var incorrectlyMarkedAsNotImplemented []PacketInfo

	for _, packetInfo := range packetsFilePackets {
		isImplemented := goFilePacketsMap[packetInfo.Code]

		// Skip commented lines
		if strings.HasPrefix(strings.TrimSpace(packetInfo.Line), "#") {
			continue
		}

		if packetInfo.IsXMark && !isImplemented {
			// Packet is marked as implemented but not actually implemented
			incorrectlyMarkedAsImplemented = append(incorrectlyMarkedAsImplemented, packetInfo)
		} else if !packetInfo.IsXMark && isImplemented {
			// Packet is implemented but not marked as such
			incorrectlyMarkedAsNotImplemented = append(incorrectlyMarkedAsNotImplemented, packetInfo)
		}
	}

	// Print results
	fmt.Printf("Total packets in %s: %d\n", packetsFilePath, len(packetsFilePackets))
	fmt.Printf("Total packets in %s: %d\n", goFilePath, len(goFilePackets))

	fmt.Printf("\nIncorrectly marked as implemented: %d\n", len(incorrectlyMarkedAsImplemented))
	if len(incorrectlyMarkedAsImplemented) > 0 {
		fmt.Println("These packets are marked with [X] but are not implemented in the Go file:")
		for i, packet := range incorrectlyMarkedAsImplemented {
			if i < 20 {
				fmt.Printf("%s: %s (line %d)\n", packet.Code, packet.Name, packet.LineNum)
			} else if i == 20 {
				fmt.Println("... (more packets omitted for brevity)")
				break
			}
		}
	}

	fmt.Printf("\nIncorrectly marked as not implemented: %d\n", len(incorrectlyMarkedAsNotImplemented))
	if len(incorrectlyMarkedAsNotImplemented) > 0 {
		fmt.Println("These packets are implemented in the Go file but not marked with [X]:")
		for _, packet := range incorrectlyMarkedAsNotImplemented {
			fmt.Printf("%s: %s (line %d)\n", packet.Code, packet.Name, packet.LineNum)
		}
	}

	// Generate a fix script
	generateFixScript(packetsFilePath, incorrectlyMarkedAsImplemented, incorrectlyMarkedAsNotImplemented)
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
	re := regexp.MustCompile(`(?:\[([X ])\]\s*)?'([0-9A-F]{4})'\s*=>\s*\['([^']+)'`)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

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
				LineNum: lineNum,
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

func generateFixScript(packetsFilePath string, incorrectlyMarkedAsImplemented, incorrectlyMarkedAsNotImplemented []PacketInfo) {
	// Create a shell script to fix the packet markings
	scriptContent := "#!/bin/bash\n\n"
	scriptContent += "# This script fixes the packet markings in servertype0.packets\n\n"

	// Fix packets incorrectly marked as implemented
	for _, packet := range incorrectlyMarkedAsImplemented {
		oldLine := packet.Line
		newLine := strings.Replace(oldLine, "[X]", "[ ]", 1)
		scriptContent += fmt.Sprintf("sed -i '%d s/%s/%s/' %s\n",
			packet.LineNum,
			strings.Replace(oldLine, "/", "\\/", -1),
			strings.Replace(newLine, "/", "\\/", -1),
			packetsFilePath)
	}

	// Fix packets incorrectly marked as not implemented
	for _, packet := range incorrectlyMarkedAsNotImplemented {
		oldLine := packet.Line
		newLine := strings.Replace(oldLine, "[ ]", "[X]", 1)
		scriptContent += fmt.Sprintf("sed -i '%d s/%s/%s/' %s\n",
			packet.LineNum,
			strings.Replace(oldLine, "/", "\\/", -1),
			strings.Replace(newLine, "/", "\\/", -1),
			packetsFilePath)
	}

	// Write the script to a file
	scriptPath := "fix_packet_markings.sh"
	err := ioutil.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		fmt.Printf("Error writing fix script: %v\n", err)
		return
	}

	fmt.Printf("\nFix script generated: %s\n", scriptPath)
	fmt.Println("Run this script to fix the packet markings in servertype0.packets")
}

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lenaxia/goKore/cmd/common"
)

// TestResponseCheckerLoading tests that the response checker can load packet dumps
func TestResponseCheckerLoading(t *testing.T) {
	// Create a temporary test dump file
	tempDir := t.TempDir()
	testDumpPath := filepath.Join(tempDir, "test_dump.json")

	// Create a simple test dump
	testDump := common.PacketDump{
		FileName:    "test_dump.json",
		PacketCount: 4,
		Packets: []common.Packet{
			{
				Direction:   "received",
				PacketID:    "0AC4",
				Description: "Account Info With Server Info",
				Size:        16,
				Timestamp:   "2025.04.30 16:24:58",
				RawData: []common.RawData{
					{
						HexBytes: []string{"C4", "0A", "E0", "00", "E5", "5D", "F6", "C1", "82", "84"},
					},
				},
			},
			{
				Direction:   "sent",
				PacketID:    "0065",
				Description: "Game Login",
				Size:        19,
				Timestamp:   "2025.04.30 16:24:59",
				RawData: []common.RawData{
					{
						HexBytes: []string{"65", "00", "E5", "5D", "F6", "C1", "82", "84", "1E", "00"},
					},
				},
			},
			{
				Direction:   "received",
				PacketID:    "006B",
				Description: "Received characters from Game Login",
				Size:        24,
				Timestamp:   "2025.04.30 16:25:00",
				RawData: []common.RawData{
					{
						HexBytes: []string{"6B", "00", "B6", "00", "0F", "0F", "0F", "00", "00", "00"},
					},
				},
			},
			{
				Direction:   "sent",
				PacketID:    "0066",
				Description: "Character Login",
				Size:        3,
				Timestamp:   "2025.04.30 16:25:01",
				RawData: []common.RawData{
					{
						HexBytes: []string{"66", "00", "00"},
					},
				},
			},
		},
	}

	// Write the test dump to a file
	dumpJSON, err := json.Marshal(testDump)
	if err != nil {
		t.Fatalf("Failed to marshal test dump: %v", err)
	}

	err = os.WriteFile(testDumpPath, dumpJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to write test dump file: %v", err)
	}

	// Create a response checker
	rc := NewResponseChecker()

	// Load the test dump
	err = rc.loadDump(testDumpPath)
	if err != nil {
		t.Fatalf("Failed to load test dump: %v", err)
	}

	// Check that the dump was loaded correctly
	if rc.selectedDump == nil {
		t.Fatal("Selected dump is nil")
	}

	if rc.selectedDump.PacketCount != 4 {
		t.Errorf("Expected packet count to be 4, got %d", rc.selectedDump.PacketCount)
	}
}

// TestExpectedResponses tests that the expected responses are correctly defined
func TestExpectedResponses(t *testing.T) {
	rc := NewResponseChecker()

	// Check that we have the expected number of responses
	if len(rc.expectedResponses) != 4 {
		t.Errorf("Expected 4 expected responses, got %d", len(rc.expectedResponses))
	}

	// Check that the expected responses are correctly defined
	expectedPairs := map[string]string{
		"0AC4": "0065", // Account Info -> Game Login
		"006B": "0066", // Received characters -> Character Login
		"0AC5": "0072", // Map Server Info -> Map Login
		"02EB": "007D", // Map Loaded -> Map Loaded
	}

	for _, resp := range rc.expectedResponses {
		expectedClientID, ok := expectedPairs[resp.ServerPacketID]
		if !ok {
			t.Errorf("Unexpected server packet ID: %s", resp.ServerPacketID)
			continue
		}

		if resp.ClientPacketID != expectedClientID {
			t.Errorf("For server packet %s, expected client packet %s, got %s",
				resp.ServerPacketID, expectedClientID, resp.ClientPacketID)
		}
	}
}

// TestDetailedReport tests the detailed report functionality
func TestDetailedReport(t *testing.T) {
	// Create a response checker
	rc := NewResponseChecker()

	// Create test results
	results := []ResponseResult{
		{
			Result: "✓",
			ServerPacket: common.Packet{
				Direction:   "received",
				PacketID:    "0AC4",
				Description: "Account Info With Server Info",
				Timestamp:   "2025.04.30 16:24:58",
				RawData: []common.RawData{
					{
						HexBytes: []string{"C4", "0A", "E0", "00", "E5", "5D", "F6", "C1", "82", "84"},
					},
				},
			},
			ExpectedPacketID: "0065",
			ExpectedDesc:     "Game Login",
			ActualPacket: &common.Packet{
				Direction:   "sent",
				PacketID:    "0065",
				Description: "Game Login",
				Timestamp:   "2025.04.30 16:24:59",
				RawData: []common.RawData{
					{
						HexBytes: []string{"65", "00", "E5", "5D", "F6", "C1", "82", "84", "1E", "00"},
					},
				},
			},
		},
		{
			Result: "✗",
			ServerPacket: common.Packet{
				Direction:   "received",
				PacketID:    "006B",
				Description: "Received characters from Game Login",
				Timestamp:   "2025.04.30 16:25:00",
				RawData: []common.RawData{
					{
						HexBytes: []string{"6B", "00", "B6", "00", "0F", "0F", "0F", "00", "00", "00"},
					},
				},
			},
			ExpectedPacketID: "0066",
			ExpectedDesc:     "Character Login",
			ActualPacket: &common.Packet{
				Direction:   "sent",
				PacketID:    "0067", // Wrong packet ID
				Description: "Wrong Response",
				Timestamp:   "2025.04.30 16:25:01",
				RawData: []common.RawData{
					{
						HexBytes: []string{"67", "00", "00"},
					},
				},
			},
		},
		{
			Result:           "?",
			ServerPacket:     common.Packet{},
			ExpectedPacketID: "0072",
			ExpectedDesc:     "Map Login",
			ActualPacket:     nil,
		},
	}

	// We don't need to test the resultToEmoji function directly
	// since it's an internal function used by printDetailedReport

	// We can't easily test the printDetailedReport function directly since it prints to stdout,
	// but we can at least make sure it doesn't panic
	rc.printDetailedReport(results)
}

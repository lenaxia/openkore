package ranking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestTaekwonPackets(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("ranking.taekwon_packets", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewRankingManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedStatus string
		expectedFlag   byte
		expectedValue  byte
	}{
		{
			name: "Place of the Sun - Map Registered",
			args: map[string]interface{}{
				"flag":  byte(0),
				"value": byte(1),
				"name":  []byte("prontera"),
			},
			expectedStatus: "You have now marked: prontera as Place of the Sun.",
			expectedFlag:   byte(0),
			expectedValue:  byte(1),
		},
		{
			name: "Place of the Moon - Information",
			args: map[string]interface{}{
				"flag":  byte(1),
				"value": byte(2),
				"name":  []byte("geffen"),
			},
			expectedStatus: "geffen is marked as Place of the Moon.",
			expectedFlag:   byte(1),
			expectedValue:  byte(2),
		},
		{
			name: "Place of the Stars - Information",
			args: map[string]interface{}{
				"flag":  byte(1),
				"value": byte(3),
				"name":  []byte("payon"),
			},
			expectedStatus: "payon is marked as Place of the Stars.",
			expectedFlag:   byte(1),
			expectedValue:  byte(3),
		},
		{
			name: "Target of the Sun - Register mob",
			args: map[string]interface{}{
				"flag":  byte(10),
				"value": byte(1),
				"name":  []byte("Poring"),
			},
			expectedStatus: "You have now marked Poring as Target of the Sun.",
			expectedFlag:   byte(10),
			expectedValue:  byte(1),
		},
		{
			name: "Target of the Moon - Information",
			args: map[string]interface{}{
				"flag":  byte(11),
				"value": byte(2),
				"name":  []byte("Lunatic"),
			},
			expectedStatus: "Lunatic is marked as Target of the Moon.",
			expectedFlag:   byte(11),
			expectedValue:  byte(2),
		},
		{
			name: "TaeKwon Mission Target",
			args: map[string]interface{}{
				"flag":  byte(20),
				"value": byte(50),
				"name":  []byte("Zombie"),
			},
			expectedStatus: "[TaeKwon Mission] Target Monster : Zombie (50%)",
			expectedFlag:   byte(20),
			expectedValue:  byte(50),
		},
		{
			name: "Feel/Hate Reset",
			args: map[string]interface{}{
				"flag":  byte(30),
				"value": byte(0),
				"name":  []byte(""),
			},
			expectedStatus: "Your Hate and Feel targets have been resetted.",
			expectedFlag:   byte(30),
			expectedValue:  byte(0),
		},
		{
			name: "Unknown Flag",
			args: map[string]interface{}{
				"flag":  byte(99),
				"value": byte(0),
				"name":  []byte(""),
			},
			expectedStatus: "Unknown results in taekwon_packets (flag: 99)",
			expectedFlag:   byte(99),
			expectedValue:  byte(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleTaekwonPackets(tc.args)
			if err != nil {
				t.Errorf("handleTaekwonPackets returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}
			if result["flag"] != tc.args["flag"] {
				t.Errorf("Expected flag %v, got %v", tc.args["flag"], result["flag"])
			}
			if result["value"] != tc.args["value"] {
				t.Errorf("Expected value %v, got %v", tc.args["value"], result["value"])
			}
		})
	}
}

package ranking

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestTop10(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to capture hook calls
	blacksmithChan := make(chan map[string]interface{}, 1)
	alchemistChan := make(chan map[string]interface{}, 1)
	taekwonChan := make(chan map[string]interface{}, 1)
	pkChan := make(chan map[string]interface{}, 1)
	unknownChan := make(chan map[string]interface{}, 1)

	// Register hooks to capture the results
	hookManager.AddHook("ranking.top10_blacksmith_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		blacksmithChan <- result
	}, nil)

	hookManager.AddHook("ranking.top10_alchemist_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		alchemistChan <- result
	}, nil)

	hookManager.AddHook("ranking.top10_taekwon_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		taekwonChan <- result
	}, nil)

	hookManager.AddHook("ranking.top10_pk_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		pkChan <- result
	}, nil)

	hookManager.AddHook("ranking.top10_unknown", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		unknownChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewRankingManager(nil, hookManager)

	// Create test data for rankings
	// 10 names (24 bytes each) followed by 10 points (4 bytes each)
	createRankingData := func() []byte {
		rawMsg := make([]byte, 0)

		// Add header bytes (2 bytes)
		rawMsg = append(rawMsg, 0x00, 0x00)

		// Add 10 names (24 bytes each)
		names := []string{
			"Player1",
			"Player2",
			"Player3",
			"Player4",
			"Player5",
			"Player6",
			"Player7",
			"Player8",
			"Player9",
			"Player10",
		}

		for _, name := range names {
			nameBytes := make([]byte, 24)
			copy(nameBytes, []byte(name))
			rawMsg = append(rawMsg, nameBytes...)
		}

		// Add 10 points (4 bytes each)
		points := []uint32{
			1000,
			900,
			800,
			700,
			600,
			500,
			400,
			300,
			200,
			100,
		}

		for _, point := range points {
			pointBytes := []byte{
				byte(point & 0xFF),
				byte((point >> 8) & 0xFF),
				byte((point >> 16) & 0xFF),
				byte((point >> 24) & 0xFF),
			}
			rawMsg = append(rawMsg, pointBytes...)
		}

		return rawMsg
	}

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		expectedType  byte
		expectedChan  chan map[string]interface{}
		expectMessage bool
	}{
		{
			name: "Blacksmith Ranking",
			args: map[string]interface{}{
				"type":    byte(0),
				"RAW_MSG": createRankingData(),
			},
			expectedType:  byte(0),
			expectedChan:  blacksmithChan,
			expectMessage: true,
		},
		{
			name: "Alchemist Ranking",
			args: map[string]interface{}{
				"type":    byte(1),
				"RAW_MSG": createRankingData(),
			},
			expectedType:  byte(1),
			expectedChan:  alchemistChan,
			expectMessage: true,
		},
		{
			name: "Taekwon Ranking",
			args: map[string]interface{}{
				"type":    byte(2),
				"RAW_MSG": createRankingData(),
			},
			expectedType:  byte(2),
			expectedChan:  taekwonChan,
			expectMessage: true,
		},
		{
			name: "PK Ranking",
			args: map[string]interface{}{
				"type":    byte(3),
				"RAW_MSG": createRankingData(),
			},
			expectedType:  byte(3),
			expectedChan:  pkChan,
			expectMessage: true,
		},
		{
			name: "Unknown Ranking Type",
			args: map[string]interface{}{
				"type":    byte(4),
				"RAW_MSG": createRankingData(),
			},
			expectedType:  byte(4),
			expectedChan:  unknownChan,
			expectMessage: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleTop10(tc.args)
			if err != nil {
				t.Errorf("handleTop10 returned an error: %v", err)
			}

			// For tests that expect a message, check the hook result
			if tc.expectMessage {
				// Get the result from the channel
				select {
				case result := <-tc.expectedChan:
					// Verify the result has rankings
					rankings := result["rankings"].([]map[string]interface{})
					if len(rankings) != 10 {
						t.Errorf("Expected 10 rankings, got %d", len(rankings))
					}

					// Check the first ranking
					if rankings[0]["name"] != "Player1" {
						t.Errorf("Expected name 'Player1', got '%s'", rankings[0]["name"])
					}

					// Check the status message format (not the exact content)
					status := result["status"].(string)
					if len(status) == 0 {
						t.Errorf("Status message is empty")
					}
				default:
					t.Errorf("Expected a message but none was received")
				}
			} else {
				// For unknown type, we should not receive any message
				select {
				case <-blacksmithChan:
					t.Errorf("Unexpected blacksmith message received")
				case <-alchemistChan:
					t.Errorf("Unexpected alchemist message received")
				case <-taekwonChan:
					t.Errorf("Unexpected taekwon message received")
				case <-pkChan:
					t.Errorf("Unexpected pk message received")
				default:
					// This is the expected case - no message
				}
			}
		})
	}
}

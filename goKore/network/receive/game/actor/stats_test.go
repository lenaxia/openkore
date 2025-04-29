package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestGoldPCCafePoint(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.gold_pc_cafe_point", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases
	testCases := []struct {
		name          string
		args          map[string]interface{}
		expectedPoint uint32
		isActive      bool
	}{
		{
			name: "Active PC Cafe Points",
			args: map[string]interface{}{
				"isActive":   byte(1),
				"mode":       byte(0),
				"point":      uint32(100),
				"playedTime": uint32(3600),
			},
			expectedPoint: 100,
			isActive:      true,
		},
		{
			name: "Inactive PC Cafe Points",
			args: map[string]interface{}{
				"isActive":   byte(0),
				"mode":       byte(0),
				"point":      uint32(50),
				"playedTime": uint32(1800),
			},
			expectedPoint: 50,
			isActive:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleGoldPCCafePoint(tc.args)
			if err != nil {
				t.Errorf("HandleGoldPCCafePoint returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if point, ok := result["point"].(uint32); !ok || point != tc.expectedPoint {
				t.Errorf("Expected point %d, got %v", tc.expectedPoint, result["point"])
			}

			if isActive, ok := result["isActive"].(bool); !ok || isActive != tc.isActive {
				t.Errorf("Expected isActive %v, got %v", tc.isActive, result["isActive"])
			}

			if _, ok := result["playedTime"].(uint32); !ok {
				t.Errorf("Expected playedTime to be uint32, got %T", result["playedTime"])
			}

			if _, ok := result["mode"].(byte); !ok {
				t.Errorf("Expected mode to be byte, got %T", result["mode"])
			}
		})
	}
}

func TestAlchemistPoint(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.alchemist_point", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test case
	args := map[string]interface{}{
		"points": uint32(10),
		"total":  uint32(100),
	}

	// Call the handler
	err := handler.HandleAlchemistPoint(args)
	if err != nil {
		t.Errorf("HandleAlchemistPoint returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	if points, ok := result["points"].(uint32); !ok || points != 10 {
		t.Errorf("Expected points 10, got %v", result["points"])
	}

	if total, ok := result["total"].(uint32); !ok || total != 100 {
		t.Errorf("Expected total 100, got %v", result["total"])
	}
}

func TestBlacksmithPoints(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.blacksmith_points", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test case
	args := map[string]interface{}{
		"points": uint32(15),
		"total":  uint32(150),
	}

	// Call the handler
	err := handler.HandleBlacksmithPoints(args)
	if err != nil {
		t.Errorf("HandleBlacksmithPoints returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	if points, ok := result["points"].(uint32); !ok || points != 15 {
		t.Errorf("Expected points 15, got %v", result["points"])
	}

	if total, ok := result["total"].(uint32); !ok || total != 150 {
		t.Errorf("Expected total 150, got %v", result["total"])
	}
}

func TestRankPoints(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to capture hook calls
	blacksmithChan := make(chan map[string]interface{}, 1)
	alchemistChan := make(chan map[string]interface{}, 1)
	taekwonChan := make(chan map[string]interface{}, 1)
	unknownChan := make(chan map[string]interface{}, 1)

	// Register hooks to capture the results
	hookManager.AddHook("actor.blacksmith_points", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		blacksmithChan <- result
	}, nil)

	hookManager.AddHook("actor.alchemist_point", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		alchemistChan <- result
	}, nil)

	hookManager.AddHook("actor.taekwon_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		taekwonChan <- result
	}, nil)

	hookManager.AddHook("actor.unknown_rank", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		unknownChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Test cases
	testCases := []struct {
		name         string
		args         map[string]interface{}
		rankType     uint16
		points       uint32
		total        uint32
		expectedHook string
	}{
		{
			name: "Blacksmith Rank",
			args: map[string]interface{}{
				"type":   uint16(0),
				"points": uint32(10),
				"total":  uint32(100),
			},
			rankType:     0,
			points:       10,
			total:        100,
			expectedHook: "blacksmith",
		},
		{
			name: "Alchemist Rank",
			args: map[string]interface{}{
				"type":   uint16(1),
				"points": uint32(15),
				"total":  uint32(150),
			},
			rankType:     1,
			points:       15,
			total:        150,
			expectedHook: "alchemist",
		},
		{
			name: "Taekwon Rank",
			args: map[string]interface{}{
				"type":   uint16(2),
				"points": uint32(20),
				"total":  uint32(200),
			},
			rankType:     2,
			points:       20,
			total:        200,
			expectedHook: "taekwon",
		},
		{
			name: "Unknown Rank",
			args: map[string]interface{}{
				"type":   uint16(3),
				"points": uint32(25),
				"total":  uint32(250),
			},
			rankType:     3,
			points:       25,
			total:        250,
			expectedHook: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleRankPoints(tc.args)
			if err != nil {
				t.Errorf("HandleRankPoints returned an error: %v", err)
			}

			// Check the appropriate channel based on rank type
			var result map[string]interface{}
			switch tc.rankType {
			case 0: // Blacksmith
				select {
				case result = <-blacksmithChan:
					// Verify the result
					if points, ok := result["points"].(uint32); !ok || points != tc.points {
						t.Errorf("Expected points %d, got %v", tc.points, result["points"])
					}
					if total, ok := result["total"].(uint32); !ok || total != tc.total {
						t.Errorf("Expected total %d, got %v", tc.total, result["total"])
					}
				default:
					t.Error("Expected blacksmith hook to be called, but it wasn't")
				}
			case 1: // Alchemist
				select {
				case result = <-alchemistChan:
					// Verify the result
					if points, ok := result["points"].(uint32); !ok || points != tc.points {
						t.Errorf("Expected points %d, got %v", tc.points, result["points"])
					}
					if total, ok := result["total"].(uint32); !ok || total != tc.total {
						t.Errorf("Expected total %d, got %v", tc.total, result["total"])
					}
				default:
					t.Error("Expected alchemist hook to be called, but it wasn't")
				}
			case 2: // Taekwon
				select {
				case result = <-taekwonChan:
					// Verify the result
					if rank, ok := result["rank"].(uint32); !ok || rank != tc.total {
						t.Errorf("Expected rank %d, got %v", tc.total, result["rank"])
					}
				default:
					t.Error("Expected taekwon hook to be called, but it wasn't")
				}
			default: // Unknown
				select {
				case result = <-unknownChan:
					// Verify the result
					if rankType, ok := result["type"].(uint16); !ok || rankType != tc.rankType {
						t.Errorf("Expected rank type %d, got %v", tc.rankType, result["type"])
					}
				default:
					t.Error("Expected unknown rank hook to be called, but it wasn't")
				}
			}
		})
	}
}

func TestRatesInfo2(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.rates_info2", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"exp":   int32(1500), // 1.5% (stored as 1500/1000)
		"death": int32(500),  // 0.5% (stored as 500/1000)
		"drop":  int32(2000), // 2.0% (stored as 2000/1000)
		"RAW_MSG": []byte{
			// Header (already processed)
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// Detail type 0 (base server)
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// Detail type 1 (premium)
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			// Detail type 2 (server additional)
			0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		"RAW_MSG_SIZE": 53, // 14 (header) + 13*3 (details)
	}

	// Call the handler
	err := handler.HandleRatesInfo2(args)
	if err != nil {
		t.Errorf("HandleRatesInfo2 returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	rates, ok := result["rates"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected rates to be a map, got %T", result["rates"])
	}

	// Check exp rates
	exp, ok := rates["exp"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected exp to be a map, got %T", rates["exp"])
	}
	if total, ok := exp["total"].(float64); !ok || total != 1.5 {
		t.Errorf("Expected exp total 1.5, got %v", exp["total"])
	}

	// Check death rates
	death, ok := rates["death"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected death to be a map, got %T", rates["death"])
	}
	if total, ok := death["total"].(float64); !ok || total != 0.5 {
		t.Errorf("Expected death total 0.5, got %v", death["total"])
	}

	// Check drop rates
	drop, ok := rates["drop"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected drop to be a map, got %T", rates["drop"])
	}
	if total, ok := drop["total"].(float64); !ok || total != 2.0 {
		t.Errorf("Expected drop total 2.0, got %v", drop["total"])
	}
}

func TestPremiumRatesInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("actor.premium_rates_info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create test packet arguments
	args := map[string]interface{}{
		"exp":   int16(50),  // +50%
		"death": int16(-25), // -25%
		"drop":  int16(75),  // +75%
	}

	// Call the handler
	err := handler.HandlePremiumRatesInfo(args)
	if err != nil {
		t.Errorf("HandlePremiumRatesInfo returned an error: %v", err)
	}

	// Get the result from the channel
	result := <-resultChan

	// Verify the result
	if exp, ok := result["exp"].(int16); !ok || exp != 50 {
		t.Errorf("Expected exp 50, got %v", result["exp"])
	}

	if death, ok := result["death"].(int16); !ok || death != -25 {
		t.Errorf("Expected death -25, got %v", result["death"])
	}

	if drop, ok := result["drop"].(int16); !ok || drop != 75 {
		t.Errorf("Expected drop 75, got %v", result["drop"])
	}
}

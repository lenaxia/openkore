package actor

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleMonsterHPInfo(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("monster.hp_info", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a monster for testing
	monster := NewMonster([]byte{1, 2, 3, 4})
	monster.SetName("Poring")
	handler.monstersList.Add(monster)

	// Test cases
	testCases := []struct {
		name       string
		args       map[string]interface{}
		hp         uint32
		hp_max     uint32
		expectHook bool
	}{
		{
			name: "Monster HP Info",
			args: map[string]interface{}{
				"ID":     []byte{1, 2, 3, 4},
				"hp":     uint32(500),
				"hp_max": uint32(1000),
			},
			hp:         500,
			hp_max:     1000,
			expectHook: true,
		},
		{
			name: "Unknown Monster",
			args: map[string]interface{}{
				"ID":     []byte{5, 6, 7, 8},
				"hp":     uint32(500),
				"hp_max": uint32(1000),
			},
			expectHook: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleMonsterHPInfo(tc.args)
			if err != nil {
				t.Errorf("HandleMonsterHPInfo returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if hp, ok := result["hp"].(uint32); !ok || hp != tc.hp {
					t.Errorf("Expected hp %d, got %v", tc.hp, result["hp"])
				}

				if hp_max, ok := result["hp_max"].(uint32); !ok || hp_max != tc.hp_max {
					t.Errorf("Expected hp_max %d, got %v", tc.hp_max, result["hp_max"])
				}

				// Check if monster's HP was updated
				if monster.HP() != tc.hp {
					t.Errorf("Expected monster HP %d, got %d", tc.hp, monster.HP())
				}

				if monster.MaxHP() != tc.hp_max {
					t.Errorf("Expected monster MaxHP %d, got %d", tc.hp_max, monster.MaxHP())
				}

				// Check if HP percent is calculated correctly
				expectedPercent := int((float64(tc.hp) / float64(tc.hp_max)) * 100)
				if monster.HPPercent() != expectedPercent {
					t.Errorf("Expected monster HP percent %d, got %d", expectedPercent, monster.HPPercent())
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleMonsterHPInfoUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing ID",
			args: map[string]interface{}{
				"hp":     uint32(500),
				"hp_max": uint32(1000),
			},
		},
		{
			name: "Missing hp",
			args: map[string]interface{}{
				"ID":     []byte{1, 2, 3, 4},
				"hp_max": uint32(1000),
			},
		},
		{
			name: "Missing hp_max",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
				"hp": uint32(500),
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID":     "invalid",
				"hp":     uint32(500),
				"hp_max": uint32(1000),
			},
		},
		{
			name: "Invalid hp type",
			args: map[string]interface{}{
				"ID":     []byte{1, 2, 3, 4},
				"hp":     "invalid",
				"hp_max": uint32(1000),
			},
		},
		{
			name: "Invalid hp_max type",
			args: map[string]interface{}{
				"ID":     []byte{1, 2, 3, 4},
				"hp":     uint32(500),
				"hp_max": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleMonsterHPInfo(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}

func TestHandleMonsterHPInfoTiny(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("monster.hp_info_tiny", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a handler for testing
	handler := NewHandler()
	handler.SetHookManager(hookManager)

	// Create a monster for testing
	monster := NewMonster([]byte{1, 2, 3, 4})
	monster.SetName("Poring")
	handler.monstersList.Add(monster)

	// Test cases
	testCases := []struct {
		name       string
		args       map[string]interface{}
		hp         byte
		hp_percent int
		expectHook bool
	}{
		{
			name: "Monster HP Info Tiny",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
				"hp": byte(10),
			},
			hp:         10,
			hp_percent: 50, // 10 * 5 = 50%
			expectHook: true,
		},
		{
			name: "Unknown Monster",
			args: map[string]interface{}{
				"ID": []byte{5, 6, 7, 8},
				"hp": byte(10),
			},
			expectHook: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := handler.HandleMonsterHPInfoTiny(tc.args)
			if err != nil {
				t.Errorf("HandleMonsterHPInfoTiny returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if id, ok := result["ID"].([]byte); !ok || string(id) != string(tc.args["ID"].([]byte)) {
					t.Errorf("Expected ID %v, got %v", tc.args["ID"], result["ID"])
				}

				if hp, ok := result["hp"].(byte); !ok || hp != tc.hp {
					t.Errorf("Expected hp %d, got %v", tc.hp, result["hp"])
				}

				// Check if monster's HP percent was updated
				if monster.HPPercent() != tc.hp_percent {
					t.Errorf("Expected monster HP percent %d, got %d", tc.hp_percent, monster.HPPercent())
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleMonsterHPInfoTinyUnhappy(t *testing.T) {
	// Create a handler for testing
	handler := NewHandler()

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing ID",
			args: map[string]interface{}{
				"hp": byte(10),
			},
		},
		{
			name: "Missing hp",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
			},
		},
		{
			name: "Invalid ID type",
			args: map[string]interface{}{
				"ID": "invalid",
				"hp": byte(10),
			},
		},
		{
			name: "Invalid hp type",
			args: map[string]interface{}{
				"ID": []byte{1, 2, 3, 4},
				"hp": "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := handler.HandleMonsterHPInfoTiny(tc.args)

			// For unhappy paths, we expect an error
			if err == nil {
				t.Errorf("Expected error for unhappy path, but got nil")
			}
		})
	}
}

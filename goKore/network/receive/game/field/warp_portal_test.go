package field

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestWarpPortalList(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create channels to capture hook calls
	warpListChan := make(chan map[string]interface{}, 1)
	configUpdateChan := make(chan map[string]interface{}, 1)

	// Register hooks to capture the results
	hookManager.AddHook("warp.portal_list", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		warpListChan <- result
	}, nil)

	hookManager.AddHook("config.update", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		configUpdateChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewFieldManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedType   byte
		expectedMemos  []string
		configExpected bool
		configKey      string
		configValue    string
	}{
		{
			name: "Teleport Skill",
			args: map[string]interface{}{
				"type":  byte(26), // Teleport skill
				"memo1": "prontera.gat",
				"memo2": "geffen.gat",
				"memo3": "morocc.gat",
				"memo4": "alberta.gat",
			},
			expectedType:   26,
			expectedMemos:  []string{"prontera", "geffen", "morocc", "alberta"},
			configExpected: true,
			configKey:      "saveMap",
			configValue:    "geffen",
		},
		{
			name: "Butterfly Wing",
			args: map[string]interface{}{
				"type":  byte(27), // Butterfly Wing
				"memo1": "prontera.gat",
				"memo2": "geffen.gat",
				"memo3": "morocc.gat",
				"memo4": "alberta.gat",
			},
			expectedType:   27,
			expectedMemos:  []string{"prontera", "geffen", "morocc", "alberta"},
			configExpected: true,
			configKey:      "saveMap",
			configValue:    "prontera",
		},
		{
			name: "Empty Memo Fields",
			args: map[string]interface{}{
				"type":  byte(26),
				"memo1": "prontera.gat",
				"memo2": "",
				"memo3": "morocc.gat",
				"memo4": "",
			},
			expectedType:   26,
			expectedMemos:  []string{"prontera", "morocc"},
			configExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handleWarpPortalList(tc.args)
			if err != nil {
				t.Errorf("handleWarpPortalList returned an error: %v", err)
			}

			// Get the warp list result from the channel
			result := <-warpListChan

			// Verify the result
			if warpType, ok := result["type"].(byte); !ok || warpType != tc.expectedType {
				t.Errorf("Expected warp type %d, got %v", tc.expectedType, result["type"])
			}

			// Verify memo list
			if memoList, ok := result["memo_list"].([]string); ok {
				if len(memoList) != len(tc.expectedMemos) {
					t.Errorf("Expected %d memos, got %d", len(tc.expectedMemos), len(memoList))
				} else {
					for i, memo := range tc.expectedMemos {
						if i < len(memoList) && memoList[i] != memo {
							t.Errorf("Expected memo[%d] = %q, got %q", i, memo, memoList[i])
						}
					}
				}
			} else {
				t.Error("memo_list not found in result or not a []string")
			}

			// Check config update if expected
			if tc.configExpected {
				select {
				case configUpdate := <-configUpdateChan:
					// Verify config key
					if key, ok := configUpdate["key"].(string); !ok || key != tc.configKey {
						t.Errorf("Expected config key %q, got %v", tc.configKey, key)
					}
					// Verify config value
					if value, ok := configUpdate["value"].(string); !ok || value != tc.configValue {
						t.Errorf("Expected config value %q, got %v", tc.configValue, value)
					}
				default:
					t.Error("Expected config update, but none received")
				}
			}
		})
	}
}

func TestMissingWarpType(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a manager for testing
	manager := NewFieldManager(nil, hookManager)

	// Create test packet arguments with missing type
	args := map[string]interface{}{
		"memo1": "prontera.gat",
		"memo2": "geffen.gat",
		"memo3": "morocc.gat",
		"memo4": "alberta.gat",
	}

	// Call handler
	result := manager.processWarpPortalList(args)
	if result != nil {
		t.Errorf("processWarpPortalList should return nil for missing type, got %v", result)
	}
}

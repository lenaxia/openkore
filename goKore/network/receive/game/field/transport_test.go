package field

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestPrivateAirshipType(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("field.private_airship_type", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Create a manager for testing
	manager := NewFieldManager(nil, hookManager)

	// Test cases
	testCases := []struct {
		name           string
		args           map[string]interface{}
		expectedStatus string
		expectedResult int
	}{
		{
			name: "Success Result",
			args: map[string]interface{}{
				"fail": byte(0),
			},
			expectedStatus: "Use Private Airship success.",
			expectedResult: 0,
		},
		{
			name: "Try Again Result",
			args: map[string]interface{}{
				"fail": byte(1),
			},
			expectedStatus: "Please try PivateAirship again.",
			expectedResult: 1,
		},
		{
			name: "Not Enough Items Result",
			args: map[string]interface{}{
				"fail": byte(2),
			},
			expectedStatus: "You do not have enough Item to use PivateAirship.",
			expectedResult: 2,
		},
		{
			name: "Invalid Destination Result",
			args: map[string]interface{}{
				"fail": byte(3),
			},
			expectedStatus: "Destination map is invalid.",
			expectedResult: 3,
		},
		{
			name: "Invalid Source Result",
			args: map[string]interface{}{
				"fail": byte(4),
			},
			expectedStatus: "Source map is invalid.",
			expectedResult: 4,
		},
		{
			name: "Item Unavailable Result",
			args: map[string]interface{}{
				"fail": byte(5),
			},
			expectedStatus: "Item unavailable for use PivateAirship.",
			expectedResult: 5,
		},
		{
			name: "Unknown Result",
			args: map[string]interface{}{
				"fail": byte(99),
			},
			expectedStatus: "Unknown Private Airship result: 99",
			expectedResult: 99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler
			err := manager.handlePrivateAirshipType(tc.args)
			if err != nil {
				t.Errorf("handlePrivateAirshipType returned an error: %v", err)
			}

			// Get the result from the channel
			result := <-resultChan

			// Verify the result
			if result["status"] != tc.expectedStatus {
				t.Errorf("Expected status %q, got %q", tc.expectedStatus, result["status"])
			}
			if result["result"] != tc.args["fail"] {
				t.Errorf("Expected result %v, got %v", tc.args["fail"], result["result"])
			}
		})
	}
}

func TestRegisterTransportHandlers(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a field manager for testing
	manager := NewFieldManager(nil, hookManager)

	// Register handlers
	// This is a bit tricky to test directly, so we'll just check that the registration doesn't panic
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RegisterTransportHandlers panicked: %v", r)
			}
		}()
		manager.RegisterTransportHandlers()
	}()
}

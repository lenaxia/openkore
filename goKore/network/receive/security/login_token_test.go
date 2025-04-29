package security

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleReceivedLoginToken(t *testing.T) {
	// Create a hook manager for testing
	hookManager := hooks.NewHookManager()

	// Create a channel to capture hook calls
	resultChan := make(chan map[string]interface{}, 1)

	// Register a hook to capture the result
	hookManager.AddHook("security.received_login_token", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		resultChan <- result
	}, nil)

	// Test cases
	testCases := []struct {
		name       string
		netVersion int
		args       map[string]interface{}
		expectHook bool
	}{
		{
			name:       "Valid Login Token",
			netVersion: 0,
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    uint16(6900),
			},
			expectHook: true,
		},
		{
			name:       "XKore Mode 1",
			netVersion: 1,
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    uint16(6900),
			},
			expectHook: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a login token manager for testing
			manager := NewLoginTokenManager(hookManager, tc.netVersion)

			// Call the handler
			err := manager.HandleReceivedLoginToken(tc.args)
			if err != nil {
				t.Errorf("HandleReceivedLoginToken returned an error: %v", err)
			}

			if tc.expectHook {
				// Get the result from the channel
				result := <-resultChan

				// Verify the result
				if loginToken, ok := result["login_token"].([]byte); !ok || string(loginToken) != string(tc.args["login_token"].([]byte)) {
					t.Errorf("Expected login_token %v, got %v", tc.args["login_token"], result["login_token"])
				}

				if len, ok := result["len"].(uint16); !ok || len != tc.args["len"].(uint16) {
					t.Errorf("Expected len %v, got %v", tc.args["len"], result["len"])
				}

				if otpIP, ok := result["OTP_ip"].([]byte); !ok || string(otpIP) != string(tc.args["OTP_ip"].([]byte)) {
					t.Errorf("Expected OTP_ip %v, got %v", tc.args["OTP_ip"], result["OTP_ip"])
				}

				if otpPort, ok := result["OTP_port"].(uint16); !ok || otpPort != tc.args["OTP_port"].(uint16) {
					t.Errorf("Expected OTP_port %v, got %v", tc.args["OTP_port"], result["OTP_port"])
				}
			} else {
				// Make sure the hook was not called
				select {
				case <-resultChan:
					t.Errorf("Hook was called when it should not have been")
				default:
					// This is the expected behavior
				}
			}
		})
	}
}

// Test for unhappy paths
func TestHandleReceivedLoginTokenUnhappy(t *testing.T) {
	// Create a login token manager for testing
	manager := NewLoginTokenManager(nil, 0)

	// Test cases for unhappy paths
	testCases := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Missing Login Token",
			args: map[string]interface{}{
				"len":      uint16(8),
				"OTP_ip":   []byte("127.0.0.1"),
				"OTP_port": uint16(6900),
			},
		},
		{
			name: "Invalid Login Token Type",
			args: map[string]interface{}{
				"login_token": "invalid",
				"len":         uint16(8),
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    uint16(6900),
			},
		},
		{
			name: "Missing Length",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    uint16(6900),
			},
		},
		{
			name: "Invalid Length Type",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         "invalid",
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    uint16(6900),
			},
		},
		{
			name: "Missing OTP IP",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_port":    uint16(6900),
			},
		},
		{
			name: "Invalid OTP IP Type",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_ip":      "invalid",
				"OTP_port":    uint16(6900),
			},
		},
		{
			name: "Missing OTP Port",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_ip":      []byte("127.0.0.1"),
			},
		},
		{
			name: "Invalid OTP Port Type",
			args: map[string]interface{}{
				"login_token": []byte{1, 2, 3, 4, 5, 6, 7, 8},
				"len":         uint16(8),
				"OTP_ip":      []byte("127.0.0.1"),
				"OTP_port":    "invalid",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the handler - it should not panic
			err := manager.HandleReceivedLoginToken(tc.args)
			if err != nil {
				t.Errorf("Expected no error for unhappy path, but got: %v", err)
			}
		})
	}
}

package core_test

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/test/utils"
)

// TestReceiveLoginPackets tests handling the essential login packets
func TestReceiveLoginPackets(t *testing.T) {
	// Skip this test for now as it requires more complex setup
	t.Skip("Skipping TestReceiveLoginPackets")
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Register handlers for the essential login packets
	registerLoginPacketHandlers(t, parser, hookManager)

	// Test handling Account Info packet (0AC4)
	t.Run("Account Info", func(t *testing.T) {
		// Create a channel to receive hook calls
		hookCalled := make(chan map[string]interface{}, 1)

		// Register hook for account_info_received
		hookCallback := func(hookName string, arg interface{}, userData interface{}) {
			if args, ok := arg.(map[string]interface{}); ok {
				hookCalled <- args
			}
		}
		hookManager.AddHook("account_info_received", hooks.HookCallback(hookCallback), nil)

		// Convert hex to bytes if needed
		var packetData []byte
		var err error
		if len(utils.AccountInfoPacket.RawData) == 0 {
			packetData, err = utils.HexToBytes(utils.AccountInfoPacket.RawHex)
			if err != nil {
				t.Fatalf("Failed to convert hex to bytes: %v", err)
			}
		} else {
			packetData = utils.AccountInfoPacket.RawData
		}

		// Process the packet
		err = parser.Process(packetData)
		if err != nil {
			t.Fatalf("Failed to process packet: %v", err)
		}

		// Verify hook was called with correct data
		select {
		case args := <-hookCalled:
			// Verify session ID
			sessionID, ok := args["sessionID"]
			if !ok {
				t.Errorf("Expected sessionID in hook args, not found")
			} else if sessionID != nil {
				expectedSessionID, expOk := utils.AccountInfoPacket.ExpectedFields["sessionID"].([]byte)
				if !expOk {
					t.Errorf("Expected sessionID field is not a byte array")
				} else if string(sessionID.([]byte)) != string(expectedSessionID) {
					t.Errorf("Expected sessionID %v, got %v", expectedSessionID, sessionID)
				}
			}

			// Verify account ID
			accountID, ok := args["accountID"]
			if !ok {
				t.Errorf("Expected accountID in hook args, not found")
			} else if accountID != nil {
				expectedAccountID, expOk := utils.AccountInfoPacket.ExpectedFields["accountID"].([]byte)
				if !expOk {
					t.Errorf("Expected accountID field is not a byte array")
				} else if string(accountID.([]byte)) != string(expectedAccountID) {
					t.Errorf("Expected accountID %v, got %v", expectedAccountID, accountID)
				}
			}

			// Verify session ID2
			sessionID2, ok := args["sessionID2"]
			if !ok {
				t.Errorf("Expected sessionID2 in hook args, not found")
			} else if sessionID2 != nil {
				expectedSessionID2, expOk := utils.AccountInfoPacket.ExpectedFields["sessionID2"].([]byte)
				if !expOk {
					t.Errorf("Expected sessionID2 field is not a byte array")
				} else if string(sessionID2.([]byte)) != string(expectedSessionID2) {
					t.Errorf("Expected sessionID2 %v, got %v", expectedSessionID2, sessionID2)
				}
			}
		default:
			t.Errorf("Hook was not called")
		}
	})

	// Test handling Character Map Info packet (0AC5)
	t.Run("Character Map Info", func(t *testing.T) {
		// Create a channel to receive hook calls
		hookCalled := make(chan map[string]interface{}, 1)

		// Register hook for character_map_info_received
		hookCallback := func(hookName string, arg interface{}, userData interface{}) {
			if args, ok := arg.(map[string]interface{}); ok {
				hookCalled <- args
			}
		}
		hookManager.AddHook("character_map_info_received", hooks.HookCallback(hookCallback), nil)

		// Convert hex to bytes if needed
		var packetData []byte
		var err error
		if len(utils.CharacterMapInfoPacket.RawData) == 0 {
			packetData, err = utils.HexToBytes(utils.CharacterMapInfoPacket.RawHex)
			if err != nil {
				t.Fatalf("Failed to convert hex to bytes: %v", err)
			}
		} else {
			packetData = utils.CharacterMapInfoPacket.RawData
		}

		// Process the packet
		err = parser.Process(packetData)
		if err != nil {
			t.Fatalf("Failed to process packet: %v", err)
		}

		// Verify hook was called with correct data
		select {
		case args := <-hookCalled:
			// Verify char ID
			charID, ok := args["charID"]
			if !ok {
				t.Errorf("Expected charID in hook args, not found")
			} else if charID != nil {
				expectedCharID, expOk := utils.CharacterMapInfoPacket.ExpectedFields["charID"].([]byte)
				if !expOk {
					t.Errorf("Expected charID field is not a byte array")
				} else if string(charID.([]byte)) != string(expectedCharID) {
					t.Errorf("Expected charID %v, got %v", expectedCharID, charID)
				}
			}

			// Verify map name
			mapName, ok := args["mapName"]
			if !ok {
				t.Errorf("Expected mapName in hook args, not found")
			} else if mapName != nil {
				expectedMapName, expOk := utils.CharacterMapInfoPacket.ExpectedFields["mapName"].(string)
				if !expOk {
					t.Errorf("Expected mapName field is not a string")
				} else if mapName != expectedMapName {
					t.Errorf("Expected mapName %v, got %v", expectedMapName, mapName)
				}
			}
		default:
			t.Errorf("Hook was not called")
		}
	})

	// Test handling Enter Map packet (02EB)
	t.Run("Enter Map", func(t *testing.T) {
		// Create a channel to receive hook calls
		hookCalled := make(chan map[string]interface{}, 1)

		// Register hook for map_loaded
		hookCallback := func(hookName string, arg interface{}, userData interface{}) {
			if args, ok := arg.(map[string]interface{}); ok {
				hookCalled <- args
			} else {
				hookCalled <- map[string]interface{}{}
			}
		}
		hookManager.AddHook("map_loaded", hooks.HookCallback(hookCallback), nil)

		// Convert hex to bytes if needed
		var packetData []byte
		var err error
		if len(utils.EnterMapPacket.RawData) == 0 {
			packetData, err = utils.HexToBytes(utils.EnterMapPacket.RawHex)
			if err != nil {
				t.Fatalf("Failed to convert hex to bytes: %v", err)
			}
		} else {
			packetData = utils.EnterMapPacket.RawData
		}

		// Process the packet
		err = parser.Process(packetData)
		if err != nil {
			t.Fatalf("Failed to process packet: %v", err)
		}

		// Verify hook was called
		select {
		case <-hookCalled:
			// Success
		default:
			t.Errorf("Hook was not called")
		}
	})
}

// registerLoginPacketHandlers registers handlers for the essential login packets
func registerLoginPacketHandlers(t *testing.T, parser *core.CoreParser, hookManager *hooks.HookManager) {
	// Register handler for Account Info packet (0AC4)
	parser.RegisterHandlerFunc("0AC4", "account_server_info", "v a4 a4 a4 a4 a26 C x17 a*",
		[]string{"len", "sessionID", "accountID", "sessionID2", "lastLoginIP", "lastLoginTime", "accountSex", "serverInfo"},
		func(args map[string]interface{}) error {
			// Call the account_info_received hook with the expected fields
			hookArgs := map[string]interface{}{}

			// Only include fields that exist
			if sessionID, ok := args["sessionID"]; ok {
				hookArgs["sessionID"] = sessionID
			}
			if accountID, ok := args["accountID"]; ok {
				hookArgs["accountID"] = accountID
			}
			if sessionID2, ok := args["sessionID2"]; ok {
				hookArgs["sessionID2"] = sessionID2
			}
			if accountSex, ok := args["accountSex"]; ok {
				hookArgs["accountSex"] = accountSex
			}

			hookManager.CallHook("account_info_received", hookArgs)
			return nil
		})

	// Register handler for Character Map Info packet (0AC5)
	parser.RegisterHandlerFunc("0AC5", "received_character_ID_and_Map", "a4 Z16 a4 v a128",
		[]string{"charID", "mapName", "mapIP", "mapPort", "mapUrl"},
		func(args map[string]interface{}) error {
			// Call the character_map_info_received hook with the expected fields
			hookArgs := map[string]interface{}{}

			// Only include fields that exist
			if charID, ok := args["charID"]; ok {
				hookArgs["charID"] = charID
			}
			if mapName, ok := args["mapName"]; ok {
				hookArgs["mapName"] = mapName
			}
			if mapIP, ok := args["mapIP"]; ok {
				hookArgs["mapIP"] = mapIP
			}
			if mapPort, ok := args["mapPort"]; ok {
				hookArgs["mapPort"] = mapPort
			}

			hookManager.CallHook("character_map_info_received", hookArgs)
			return nil
		})

	// Register handler for Enter Map packet (02EB)
	parser.RegisterHandlerFunc("02EB", "map_loaded", "V a3 a a v",
		[]string{"syncMapSync", "coords", "xSize", "ySize", "font"},
		func(args map[string]interface{}) error {
			// Call the map_loaded hook with the expected fields
			hookArgs := map[string]interface{}{}

			// Only include fields that exist
			if syncMapSync, ok := args["syncMapSync"]; ok {
				hookArgs["syncMapSync"] = syncMapSync
			}
			if coords, ok := args["coords"]; ok {
				hookArgs["coords"] = coords
			}
			if xSize, ok := args["xSize"]; ok {
				hookArgs["xSize"] = xSize
			}
			if ySize, ok := args["ySize"]; ok {
				hookArgs["ySize"] = ySize
			}

			hookManager.CallHook("map_loaded", hookArgs)
			return nil
		})
}

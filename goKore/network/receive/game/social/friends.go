package social

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// FriendManager handles friend-related packet handling
type FriendManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewFriendManager creates a new friend manager
func NewFriendManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *FriendManager {
	return &FriendManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all friend-related packet handlers
func (fm *FriendManager) RegisterHandlers() {
	// Register friend list handler
	fm.parser.RegisterHandlerFunc("0201", "friend_list", "v Z*",
		[]string{"len", "RAW_MSG"}, fm.HandleFriendList)

	// Register friend logon handler
	fm.parser.RegisterHandlerFunc("0206", "friend_logon", "V V C",
		[]string{"friendAccountID", "friendCharID", "isNotOnline"}, fm.HandleFriendLogon)

	// Register friend request handler
	fm.parser.RegisterHandlerFunc("0207", "friend_request", "V V Z24",
		[]string{"accountID", "charID", "name"}, fm.HandleFriendRequest)

	// Register friend removed handler
	fm.parser.RegisterHandlerFunc("020A", "friend_removed", "V V",
		[]string{"friendAccountID", "friendCharID"}, fm.HandleFriendRemoved)

	// Register friend response handler
	fm.parser.RegisterHandlerFunc("0209", "friend_response", "v V V Z24",
		[]string{"type", "accountID", "charID", "name"}, fm.HandleFriendResponse)

	// Register ignore all result handler
	fm.parser.RegisterHandlerFunc("00D2", "ignore_all_result", "C C",
		[]string{"type", "error"}, fm.HandleIgnoreAllResult)

	// Register ignore player result handler
	fm.parser.RegisterHandlerFunc("00D3", "ignore_player_result", "C C",
		[]string{"type", "error"}, fm.HandleIgnorePlayerResult)
}

// HandleFriendList handles the friend_list packet (lines 6143-6163)
func (fm *FriendManager) HandleFriendList(args map[string]interface{}) error {
	// Extract packet data
	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in friend_list packet")
	}

	msgSize, ok := args["len"].(uint16)
	if !ok {
		return fmt.Errorf("invalid len in friend_list packet")
	}

	// Process friend list
	fm.logger.Info("Friend list received with %d entries", (msgSize-4)/32)

	// In a real implementation, we would parse the friend list data
	// and store it in a data structure for later use
	for i := 4; i < int(msgSize); i += 32 {
		// Extract friend data
		if i+32 <= len(rawMsg) {
			// accountID := rawMsg[i:i+4]
			// charID := rawMsg[i+4:i+8]
			name := string(rawMsg[i+8 : i+32])
			fm.logger.Debug("Friend: %s", name)
		}
	}

	return nil
}

// HandleFriendLogon handles the friend_logon packet (lines 6164-6190)
func (fm *FriendManager) HandleFriendLogon(args map[string]interface{}) error {
	// Extract packet data
	friendAccountID, ok := args["friendAccountID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid friendAccountID in friend_logon packet")
	}

	friendCharID, ok := args["friendCharID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid friendCharID in friend_logon packet")
	}

	isNotOnline, ok := args["isNotOnline"].(uint8)
	if !ok {
		return fmt.Errorf("invalid isNotOnline in friend_logon packet")
	}

	// In a real implementation, we would update the friend's online status
	// in our friend list data structure
	if isNotOnline == 1 {
		fm.logger.Info("Friend (Account ID: %d, Char ID: %d) has disconnected", friendAccountID, friendCharID)
	} else {
		fm.logger.Info("Friend (Account ID: %d, Char ID: %d) has connected", friendAccountID, friendCharID)
	}

	return nil
}

// HandleFriendRequest handles the friend_request packet (lines 6191-6208)
func (fm *FriendManager) HandleFriendRequest(args map[string]interface{}) error {
	// Extract packet data
	accountID, ok := args["accountID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid accountID in friend_request packet")
	}

	charID, ok := args["charID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid charID in friend_request packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in friend_request packet")
	}

	// Log the friend request
	fm.logger.Info("%s wants to be your friend", name)
	fm.logger.Info("Type 'friend accept' to be friend with %s, otherwise type 'friend reject'", name)

	// Call hook for friend request
	fm.hookManager.CallHook("friend_request", map[string]interface{}{
		"accountID": accountID,
		"charID":    charID,
		"name":      name,
	})

	return nil
}

// HandleFriendRemoved handles the friend_removed packet (lines 6209-6226)
func (fm *FriendManager) HandleFriendRemoved(args map[string]interface{}) error {
	// Extract packet data
	friendAccountID, ok := args["friendAccountID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid friendAccountID in friend_removed packet")
	}

	friendCharID, ok := args["friendCharID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid friendCharID in friend_removed packet")
	}

	// In a real implementation, we would remove the friend from our
	// friend list data structure and get their name
	fm.logger.Info("A friend (Account ID: %d, Char ID: %d) is no longer your friend",
		friendAccountID, friendCharID)

	return nil
}

// HandleIgnoreAllResult handles the ignore_all_result packet (lines 6933-6943)
func (fm *FriendManager) HandleIgnoreAllResult(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in ignore_all_result packet")
	}

	errorVal, ok := args["error"].(uint8)
	if !ok {
		// If error field is not present, default to 0
		errorVal = 0
	}

	// Process based on type
	if typeVal == 0 {
		// Ignore all players
		fm.logger.Info("All Players ignored")
	} else if typeVal == 1 {
		// Unignore all players
		if errorVal == 0 {
			fm.logger.Info("All players unignored")
		} else {
			fm.logger.Warning("Failed to unignore all players (error: %d)", errorVal)
		}
	} else {
		fm.logger.Warning("Unknown ignore_all_result type: %d", typeVal)
	}

	return nil
}

// HandleIgnorePlayerResult handles the ignore_player_result packet (lines 6946-6955)
func (fm *FriendManager) HandleIgnorePlayerResult(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in ignore_player_result packet")
	}

	errorVal, ok := args["error"].(uint8)
	if !ok {
		// If error field is not present, default to 0
		errorVal = 0
	}

	// Process based on type
	if typeVal == 0 {
		// Ignore player
		fm.logger.Info("Player ignored")
	} else if typeVal == 1 {
		// Unignore player
		if errorVal == 0 {
			fm.logger.Info("Player unignored")
		} else {
			fm.logger.Warning("Failed to unignore player (error: %d)", errorVal)
		}
	} else {
		fm.logger.Warning("Unknown ignore_player_result type: %d", typeVal)
	}

	return nil
}

// HandleFriendResponse handles the friend_response packet (lines 6227-6258)
func (fm *FriendManager) HandleFriendResponse(args map[string]interface{}) error {
	// Extract packet data
	typeVal, ok := args["type"].(uint16)
	if !ok {
		return fmt.Errorf("invalid type in friend_response packet")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name in friend_response packet")
	}

	// Process based on type
	switch typeVal {
	case 0:
		fm.logger.Info("You have become friends with %s", name)
		// In a real implementation, we would add the friend to our friend list
	case 1:
		fm.logger.Info("%s does not want to be friends with you", name)
	case 2:
		fm.logger.Info("Your Friend List is full")
	case 3:
		fm.logger.Info("%s's Friend List is full", name)
	default:
		fm.logger.Info("%s rejected to be your friend", name)
	}

	return nil
}

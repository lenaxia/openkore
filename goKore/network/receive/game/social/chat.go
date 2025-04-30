package social

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// ChatManager handles chat-related packet handling
type ChatManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger
}

// NewChatManager creates a new chat manager
func NewChatManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *ChatManager {
	return &ChatManager{
		parser:      parser,
		hookManager: hookManager,
		logger:      logger,
	}
}

// RegisterHandlers registers all chat-related packet handlers
func (cm *ChatManager) RegisterHandlers() {
	// Register self chat handler
	cm.parser.RegisterHandlerFunc("00E1", "self_chat", "Z*",
		[]string{"message"}, cm.HandleSelfChat)

	// Register private message handler
	cm.parser.RegisterHandlerFunc("0097", "private_message", "Z24 Z*",
		[]string{"sender", "message"}, cm.HandlePrivateMessage)

	// Register private message sent handler
	cm.parser.RegisterHandlerFunc("00DF", "private_message_sent", "C Z24 Z*",
		[]string{"type", "recipient", "message"}, cm.HandlePrivateMessageSent)

	// Register system chat handler
	cm.parser.RegisterHandlerFunc("009A", "system_chat", "Z*",
		[]string{"message"}, cm.HandleSystemChat)

	// Register local broadcast handler
	cm.parser.RegisterHandlerFunc("009A", "local_broadcast", "Z* Z6",
		[]string{"message", "color"}, cm.HandleLocalBroadcast)

	// Register chat created handler
	cm.parser.RegisterHandlerFunc("00D7", "chat_created", "V Z*",
		[]string{"chatID", "title"}, cm.HandleChatCreated)

	// Register chat info handler
	cm.parser.RegisterHandlerFunc("00D7", "chat_info", "V Z* V v C",
		[]string{"chatID", "title", "ownerID", "limit", "public"}, cm.HandleChatInfo)

	// Register chat users handler
	cm.parser.RegisterHandlerFunc("00DD", "chat_users", "V Z*",
		[]string{"chatID", "users"}, cm.HandleChatUsers)

	// Register chat join result handler
	cm.parser.RegisterHandlerFunc("00DA", "chat_join_result", "C",
		[]string{"type"}, cm.HandleChatJoinResult)

	// Register chat modified handler
	cm.parser.RegisterHandlerFunc("00DF", "chat_modified", "V Z* v C",
		[]string{"chatID", "title", "limit", "public"}, cm.HandleChatModified)

	// Register chat newowner handler
	cm.parser.RegisterHandlerFunc("00E1", "chat_newowner", "Z24 C",
		[]string{"user", "type"}, cm.HandleChatNewowner)

	// Register chat user join handler
	cm.parser.RegisterHandlerFunc("00DC", "chat_user_join", "Z24",
		[]string{"user"}, cm.HandleChatUserJoin)

	// Register chat user leave handler
	cm.parser.RegisterHandlerFunc("00DD", "chat_user_leave", "Z24 C",
		[]string{"user", "flag"}, cm.HandleChatUserLeave)

	// Register chat removed handler
	cm.parser.RegisterHandlerFunc("00D8", "chat_removed", "V",
		[]string{"chatID"}, cm.HandleChatRemoved)

	// Register whisper list handler
	cm.parser.RegisterHandlerFunc("00D9", "whisper_list", "v Z*",
		[]string{"count", "names"}, cm.HandleWhisperList)
}

// HandleSelfChat handles the self_chat packet (lines 11383-11406)
func (cm *ChatManager) HandleSelfChat(args map[string]interface{}) error {
	// Extract packet data
	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in self_chat packet")
	}

	// Parse message to extract user and content
	parts := strings.SplitN(message, " : ", 2)
	if len(parts) != 2 {
		cm.logger.Debug("Undefined chat message format: %s", message)
		return nil
	}

	user := parts[0]
	content := parts[1]

	// Log the chat message
	cm.logger.Info("[Self] %s: %s", user, content)

	// Call hooks
	cm.hookManager.CallHook("packet_selfChat", map[string]interface{}{
		"user":    user,
		"message": content,
	})

	return nil
}

// HandlePrivateMessage handles the private_message packet (lines 11200-11239)
func (cm *ChatManager) HandlePrivateMessage(args map[string]interface{}) error {
	// Extract packet data
	sender, ok := args["sender"].(string)
	if !ok {
		return fmt.Errorf("invalid sender in private_message packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in private_message packet")
	}

	// Log the private message
	cm.logger.Info("[PM From] %s: %s", sender, message)

	// Call hooks
	cm.hookManager.CallHook("packet_privMsg", map[string]interface{}{
		"user":    sender,
		"message": message,
	})

	return nil
}

// HandlePrivateMessageSent handles the private_message_sent packet (lines 9965-9984)
func (cm *ChatManager) HandlePrivateMessageSent(args map[string]interface{}) error {
	// Extract packet data
	type_, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in private_message_sent packet")
	}

	recipient, ok := args["recipient"].(string)
	if !ok {
		return fmt.Errorf("invalid recipient in private_message_sent packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in private_message_sent packet")
	}

	// Process based on type code
	switch type_ {
	case 0: // Success
		cm.logger.Info("[PM To] %s: %s", recipient, message)

		// Call hooks
		cm.hookManager.CallHook("packet_sentPM", map[string]interface{}{
			"user":    recipient,
			"message": message,
		})
	case 1:
		cm.logger.Warning("%s is not online", recipient)
	case 2:
		cm.logger.Warning("Player %s ignored your message", recipient)
	default:
		cm.logger.Warning("Player %s doesn't want to receive messages", recipient)
	}

	return nil
}

// HandleSystemChat handles the system_chat packet (lines 3476-3519)
func (cm *ChatManager) HandleSystemChat(args map[string]interface{}) error {
	// Extract packet data
	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in system_chat packet")
	}

	// Check for special message prefixes
	if strings.HasPrefix(message, "ssss") {
		// War of Emperium message (yellow)
		message = strings.TrimPrefix(message, "ssss")
		cm.logger.Info("[WoE] %s", message)
	} else if strings.HasPrefix(message, "micc") {
		// Player broadcast message (with color)
		message = strings.TrimPrefix(message, "micc")
		cm.logger.Info("[Broadcast] %s", message)
	} else if strings.HasPrefix(message, "blue") {
		// System message (blue)
		message = strings.TrimPrefix(message, "blue")
		cm.logger.Info("[System] %s", message)
	} else {
		// Regular system message
		cm.logger.Info("[System] %s", message)
	}

	// Call hooks
	cm.hookManager.CallHook("packet_sysMsg", map[string]interface{}{
		"message": message,
	})

	return nil
}

// HandleLocalBroadcast handles the local_broadcast packet (lines 3134-3148)
func (cm *ChatManager) HandleLocalBroadcast(args map[string]interface{}) error {
	// Extract packet data
	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in local_broadcast packet")
	}

	color, ok := args["color"].(string)
	if !ok {
		color = "FFFFFF" // Default to white if color not provided
	}

	// Log the broadcast message
	cm.logger.Info("[Local Broadcast] %s", message)

	// Call hooks
	cm.hookManager.CallHook("packet_localBroadcast", map[string]interface{}{
		"message": message,
		"color":   color,
	})

	return nil
}

// HandleChatCreated handles the chat_created packet (lines 5578-5595)
func (cm *ChatManager) HandleChatCreated(args map[string]interface{}) error {
	// Extract packet data
	chatID, ok := args["chatID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid chatID in chat_created packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in chat_created packet")
	}

	// Log the chat creation
	cm.logger.Info("Chat room created: %s (ID: %d)", title, chatID)

	// Call hooks
	cm.hookManager.CallHook("chat_created", map[string]interface{}{
		"chatID": chatID,
		"title":  title,
	})

	return nil
}

// HandleChatInfo handles the chat_info packet (lines 5596-5629)
func (cm *ChatManager) HandleChatInfo(args map[string]interface{}) error {
	// Extract packet data
	chatID, ok := args["chatID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid chatID in chat_info packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in chat_info packet")
	}

	ownerID, ok := args["ownerID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ownerID in chat_info packet")
	}

	limit, ok := args["limit"].(uint16)
	if !ok {
		return fmt.Errorf("invalid limit in chat_info packet")
	}

	public, ok := args["public"].(uint8)
	if !ok {
		return fmt.Errorf("invalid public in chat_info packet")
	}

	// Log the chat info
	publicStr := "Private"
	if public == 1 {
		publicStr = "Public"
	}

	cm.logger.Info("Chat room: %s (ID: %d, Owner: %d, Limit: %d, %s)",
		title, chatID, ownerID, limit, publicStr)

	// Call hooks
	cm.hookManager.CallHook("packet_chatinfo", map[string]interface{}{
		"chatID":  chatID,
		"title":   title,
		"ownerID": ownerID,
		"limit":   limit,
		"public":  public,
	})

	return nil
}

// HandleChatUsers handles the chat_users packet (lines 5630-5666)
func (cm *ChatManager) HandleChatUsers(args map[string]interface{}) error {
	// Extract packet data
	chatID, ok := args["chatID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid chatID in chat_users packet")
	}

	users, ok := args["users"].([]map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid users in chat_users packet")
	}

	// Log the chat users
	cm.logger.Info("Chat room users (ID: %d):", chatID)
	for _, user := range users {
		name, _ := user["name"].(string)
		userType, _ := user["type"].(uint8)

		userTypeStr := "Member"
		if userType == 2 {
			userTypeStr = "Owner"
		}

		cm.logger.Info("- %s (%s)", name, userTypeStr)
	}

	// Call hooks
	cm.hookManager.CallHook("chat_joined", map[string]interface{}{
		"chatID": chatID,
		"users":  users,
	})

	return nil
}

// HandleChatJoinResult handles the chat_join_result packet (lines 5667-5701)
func (cm *ChatManager) HandleChatJoinResult(args map[string]interface{}) error {
	// Extract packet data
	type_, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in chat_join_result packet")
	}

	// Process based on type code
	switch type_ {
	case 0:
		cm.logger.Warning("Chat room is full")
	case 1:
		cm.logger.Warning("Incorrect password")
	case 2:
		cm.logger.Warning("You have been kicked from the chat room")
	default:
		cm.logger.Warning("Unable to join chat room (reason: %d)", type_)
	}

	return nil
}

// HandleChatModified handles the chat_modified packet (lines 5702-5739)
func (cm *ChatManager) HandleChatModified(args map[string]interface{}) error {
	// Extract packet data
	chatID, ok := args["chatID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid chatID in chat_modified packet")
	}

	title, ok := args["title"].(string)
	if !ok {
		return fmt.Errorf("invalid title in chat_modified packet")
	}

	limit, ok := args["limit"].(uint16)
	if !ok {
		return fmt.Errorf("invalid limit in chat_modified packet")
	}

	public, ok := args["public"].(uint8)
	if !ok {
		return fmt.Errorf("invalid public in chat_modified packet")
	}

	// Log the chat modification
	publicStr := "Private"
	if public == 1 {
		publicStr = "Public"
	}

	cm.logger.Info("Chat room modified: %s (ID: %d, Limit: %d, %s)",
		title, chatID, limit, publicStr)

	// Call hooks
	cm.hookManager.CallHook("chat_modified", map[string]interface{}{
		"chatID": chatID,
		"title":  title,
		"limit":  limit,
		"public": public,
	})

	return nil
}

// HandleChatNewowner handles the chat_newowner packet (lines 5740-5771)
func (cm *ChatManager) HandleChatNewowner(args map[string]interface{}) error {
	// Extract packet data
	user, ok := args["user"].(string)
	if !ok {
		return fmt.Errorf("invalid user in chat_newowner packet")
	}

	type_, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in chat_newowner packet")
	}

	// Process based on type code
	if type_ == 0 {
		cm.logger.Info("Chat room owner changed to: %s", user)
	} else {
		cm.logger.Info("Chat room user %s is now a normal member", user)
	}

	return nil
}

// HandleChatUserJoin handles the chat_user_join packet (lines 5772-5785)
func (cm *ChatManager) HandleChatUserJoin(args map[string]interface{}) error {
	// Extract packet data
	user, ok := args["user"].(string)
	if !ok {
		return fmt.Errorf("invalid user in chat_user_join packet")
	}

	// Log the user join
	cm.logger.Info("User joined chat room: %s", user)

	return nil
}

// HandleChatUserLeave handles the chat_user_leave packet (lines 5786-5809)
func (cm *ChatManager) HandleChatUserLeave(args map[string]interface{}) error {
	// Extract packet data
	user, ok := args["user"].(string)
	if !ok {
		return fmt.Errorf("invalid user in chat_user_leave packet")
	}

	flag, ok := args["flag"].(uint8)
	if !ok {
		return fmt.Errorf("invalid flag in chat_user_leave packet")
	}

	// Process based on flag
	if flag == 0 {
		cm.logger.Info("User left chat room: %s", user)
	} else {
		cm.logger.Info("User was kicked from chat room: %s", user)
	}

	return nil
}

// HandleChatRemoved handles the chat_removed packet (lines 5810-5823)
func (cm *ChatManager) HandleChatRemoved(args map[string]interface{}) error {
	// Extract packet data
	chatID, ok := args["chatID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid chatID in chat_removed packet")
	}

	// Log the chat removal
	cm.logger.Info("Chat room removed (ID: %d)", chatID)

	// Call hooks
	cm.hookManager.CallHook("chat_removed", map[string]interface{}{
		"chatID": chatID,
	})

	return nil
}

// HandleWhisperList handles the whisper_list packet (lines 5568-5577)
func (cm *ChatManager) HandleWhisperList(args map[string]interface{}) error {
	// Extract packet data
	names, ok := args["names"].([]string)
	if !ok {
		return fmt.Errorf("invalid names in whisper_list packet")
	}

	// Log the whisper list
	cm.logger.Info("Whisper list:")
	for _, name := range names {
		cm.logger.Info("- %s", name)
	}

	return nil
}

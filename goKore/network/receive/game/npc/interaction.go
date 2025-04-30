package npc

import (
	"fmt"
	"strings"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// InteractionManager handles NPC interaction-related packet handling
type InteractionManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// State for NPC interactions
	currentTalk  map[string]interface{}
	npcTalkState map[string]interface{}
}

// NewInteractionManager creates a new NPC interaction manager
func NewInteractionManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *InteractionManager {
	return &InteractionManager{
		parser:       parser,
		hookManager:  hookManager,
		logger:       logger,
		currentTalk:  make(map[string]interface{}),
		npcTalkState: make(map[string]interface{}),
	}
}

// RegisterHandlers registers all NPC interaction-related packet handlers
func (im *InteractionManager) RegisterHandlers() {
	// Register NPC talk handler
	im.parser.RegisterHandlerFunc("00B4", "npc_talk", "V Z*",
		[]string{"ID", "msg"}, im.HandleNpcTalk)

	// Register NPC talk close handler
	im.parser.RegisterHandlerFunc("00B6", "npc_talk_close", "V",
		[]string{"ID"}, im.HandleNpcTalkClose)

	// Register NPC talk continue handler
	im.parser.RegisterHandlerFunc("00B5", "npc_talk_continue", "V",
		[]string{"ID"}, im.HandleNpcTalkContinue)

	// Register NPC talk number handler
	im.parser.RegisterHandlerFunc("0142", "npc_talk_number", "V",
		[]string{"ID"}, im.HandleNpcTalkNumber)

	// Register NPC talk responses handler
	im.parser.RegisterHandlerFunc("00B7", "npc_talk_responses", "v V Z*",
		[]string{"len", "ID", "RAW_MSG"}, im.HandleNpcTalkResponses)

	// Register NPC talk text handler
	im.parser.RegisterHandlerFunc("01D4", "npc_talk_text", "V",
		[]string{"ID"}, im.HandleNpcTalkText)

	// Register NPC clear dialog handler
	im.parser.RegisterHandlerFunc("0146", "npc_clear_dialog", "V",
		[]string{"ID"}, im.HandleNpcClearDialog)

	// Register NPC chat handler
	im.parser.RegisterHandlerFunc("00B6", "npc_chat", "V Z*",
		[]string{"ID", "message"}, im.HandleNpcChat)

	// Register NPC image handler
	im.parser.RegisterHandlerFunc("01B3", "npc_image", "Z64 C",
		[]string{"npc_image", "type"}, im.HandleNpcImage)
}

// HandleNpcTalk handles the npc_talk packet (lines 7410-7451)
func (im *InteractionManager) HandleNpcTalk(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk packet")
	}

	msg, ok := args["msg"].(string)
	if !ok {
		return fmt.Errorf("invalid msg in npc_talk packet")
	}

	// Store the NPC ID
	im.currentTalk["ID"] = id
	im.currentTalk["nameID"] = id // In Go implementation, we use the same ID

	// Remove RO color codes
	msg = removeColorCodes(msg)

	// Prepend existing conversation
	existingMsg, _ := im.currentTalk["msg"].(string)
	if existingMsg != "" {
		im.currentTalk["msg"] = existingMsg + "\n" + msg
	} else {
		im.currentTalk["msg"] = msg
	}

	// Update NPC talk state
	im.npcTalkState["talk"] = "initiated"
	im.npcTalkState["time"] = getCurrentTime()

	// Get NPC name (in a real implementation, this would use a function to get the name from the ID)
	npcName := fmt.Sprintf("NPC-%d", id)

	// Call hooks
	im.hookManager.CallHook("npc_talk", map[string]interface{}{
		"ID":     id,
		"nameID": id,
		"name":   npcName,
		"msg":    im.currentTalk["msg"],
	})

	// Log the message
	im.logger.Info("[NPC] %s: %s", npcName, msg)

	return nil
}

// HandleNpcTalkClose handles the npc_talk_close packet (lines 7455-7474)
func (im *InteractionManager) HandleNpcTalkClose(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk_close packet")
	}

	// Check if we're in a conversation with this NPC
	currentID, _ := im.npcTalkState["ID"].(uint32)
	if currentID == 0 || currentID != id {
		im.logger.Debug("Received unexpected npc_talk_close, ignoring it")
		return nil
	}

	// Check if we're in buy_or_sell state
	if talkState, _ := im.npcTalkState["talk"].(string); talkState == "buy_or_sell" {
		return nil
	}

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Update NPC talk state
	im.npcTalkState["talk"] = "close"
	im.npcTalkState["time"] = getCurrentTime()

	// Clear current talk
	im.currentTalk = make(map[string]interface{})

	// Call hooks
	im.hookManager.CallHook("npc_talk_done", map[string]interface{}{
		"ID": id,
	})

	im.logger.Debug("Closed conversation with %s", npcName)

	return nil
}

// HandleNpcTalkContinue handles the npc_talk_continue packet (lines 7478-7485)
func (im *InteractionManager) HandleNpcTalkContinue(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk_continue packet")
	}

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Update NPC talk state
	im.npcTalkState["talk"] = "next"
	im.npcTalkState["time"] = getCurrentTime()

	im.logger.Debug("NPC %s conversation continues", npcName)

	return nil
}

// HandleNpcTalkNumber handles the npc_talk_number packet (lines 7489-7497)
func (im *InteractionManager) HandleNpcTalkNumber(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk_number packet")
	}

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Update NPC talk state
	im.npcTalkState["talk"] = "number"
	im.npcTalkState["time"] = getCurrentTime()

	im.logger.Info("NPC %s is requesting a number input", npcName)

	return nil
}

// HandleNpcTalkResponses handles the npc_talk_responses packet (lines 7501-7557)
func (im *InteractionManager) HandleNpcTalkResponses(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk_responses packet")
	}

	rawMsg, ok := args["RAW_MSG"].([]byte)
	if !ok {
		return fmt.Errorf("invalid RAW_MSG in npc_talk_responses packet")
	}

	// Store the NPC ID
	im.currentTalk["ID"] = id
	im.currentTalk["nameID"] = id

	// Parse the talk responses
	talk := string(rawMsg)
	talk = removeColorCodes(talk)

	// Split responses by colon
	preTalkResponses := strings.Split(talk, ":")
	responses := []string{}

	for _, response := range preTalkResponses {
		// Remove RO color codes
		response = removeColorCodes(response)

		// Handle item IDs (simplified)
		if strings.HasPrefix(response, "^nItemID^") {
			// In a real implementation, this would use a function to get the item name from the ID
			itemID := strings.TrimPrefix(response, "^nItemID^")
			response = fmt.Sprintf("Item-%s", itemID)
		}

		if response != "" {
			responses = append(responses, response)
		}
	}

	// Add "Cancel Chat" option
	responses = append(responses, "Cancel Chat")

	// Store responses
	im.currentTalk["responses"] = responses

	// Update NPC talk state
	im.npcTalkState["talk"] = "select"
	im.npcTalkState["time"] = getCurrentTime()

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Call hooks
	im.hookManager.CallHook("npc_talk_responses", map[string]interface{}{
		"ID":        id,
		"name":      npcName,
		"responses": responses,
	})

	// Log the responses
	im.logger.Info("NPC %s dialog options:", npcName)
	for i, resp := range responses {
		im.logger.Info("%d: %s", i+1, resp)
	}

	return nil
}

// HandleNpcTalkText handles the npc_talk_text packet (lines 7561-7569)
func (im *InteractionManager) HandleNpcTalkText(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_talk_text packet")
	}

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Update NPC talk state
	im.npcTalkState["talk"] = "text"
	im.npcTalkState["time"] = getCurrentTime()

	im.logger.Info("NPC %s is requesting text input", npcName)

	return nil
}

// HandleNpcClearDialog handles the npc_clear_dialog packet (lines 7665-7673)
func (im *InteractionManager) HandleNpcClearDialog(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_clear_dialog packet")
	}

	// Get NPC name
	npcName := fmt.Sprintf("NPC-%d", id)

	// Clear the dialog
	im.currentTalk = make(map[string]interface{})
	im.npcTalkState = make(map[string]interface{})

	im.logger.Info("NPC %s dialog cleared", npcName)

	return nil
}

// Helper function to remove RO color codes
func removeColorCodes(text string) string {
	return strings.ReplaceAll(text, "^[a-fA-F0-9]{6}", "")
}

// HandleNpcChat handles the npc_chat packet (lines 4957-4990)
func (im *InteractionManager) HandleNpcChat(args map[string]interface{}) error {
	// Extract packet data
	id, ok := args["ID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid ID in npc_chat packet")
	}

	message, ok := args["message"].(string)
	if !ok {
		return fmt.Errorf("invalid message in npc_chat packet")
	}

	// Get actor information (in a real implementation, this would use a function to get the actor from the ID)
	// actorName := fmt.Sprintf("NPC-%d", id)

	// Parse the message
	var name string
	if strings.Contains(message, " : ") {
		parts := strings.SplitN(message, " : ", 2)
		name = parts[0]
		message = parts[1]

		// In a real implementation, we would calculate the distance between the character and the NPC
		dist := "unknown"

		// Format the message
		message = fmt.Sprintf("%s: %s", name, message)

		// Log position information (simplified)
		position := fmt.Sprintf("[Unknown field, x, y] [%d, %d] [dist=%s] (%d)", 0, 0, dist, id)
		im.logger.Debug("NPC chat position: %s", position)
	}

	// Log the message
	im.logger.Info("[NPC Chat] %s", message)

	return nil
}

// HandleNpcImage handles the npc_image packet (lines 3117-3133)
func (im *InteractionManager) HandleNpcImage(args map[string]interface{}) error {
	// Extract packet data
	npcImage, ok := args["npc_image"].(string)
	if !ok {
		return fmt.Errorf("invalid npc_image in npc_image packet")
	}

	typeVal, ok := args["type"].(uint8)
	if !ok {
		return fmt.Errorf("invalid type in npc_image packet")
	}

	// Process based on type
	if typeVal == 2 {
		im.logger.Info("NPC image: %s", npcImage)
	} else if typeVal == 255 {
		im.logger.Debug("Hide NPC image: %s", npcImage)
	} else {
		im.logger.Info("NPC image: %s (unknown type %d)", npcImage, typeVal)
	}

	// Store or delete the image
	if typeVal != 255 {
		im.currentTalk["image"] = npcImage
	} else {
		delete(im.currentTalk, "image")
	}

	return nil
}

// Helper function to get current time
func getCurrentTime() int64 {
	return 0 // In a real implementation, this would return the current time
}

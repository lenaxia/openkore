// Package main provides an example of how to use the behavior package
package main

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
	"github.com/lenaxia/goKore/network/receive/core/behavior"
)

func main() {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Register a hook to handle manner messages
	hookManager.AddHook("character.manner_message", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		fmt.Printf("Manner message: %s\n", result["message"])
	}, nil)

	// Register a hook to handle hack shield alarms
	hookManager.AddHook("character.hack_shield_alarm", func(hookName string, arg interface{}, userData interface{}) {
		result := arg.(map[string]interface{})
		fmt.Printf("Hack shield alarm: %s\n", result["message"])
	}, nil)

	// Create a core parser
	parser := core.NewCoreParser("ServerType0", hookManager)

	// Method 1: Register directly with the parser
	behavior.RegisterWithParser(parser, hookManager)

	// Method 2: Use with BaseReceive
	baseReceive := core.NewBaseReceive(hookManager)
	behavior.RegisterWithBaseReceive(baseReceive)

	// Example of processing a manner_message packet
	fmt.Println("Processing manner_message packet...")
	mannerPacket := []byte{0x01, 0x49, 0x00} // 0149 packet with flag 0
	parser.Process(mannerPacket)

	// Example of processing a hack_shield_alarm packet
	fmt.Println("Processing hack_shield_alarm packet...")
	hackShieldPacket := []byte{0x08, 0xB3} // 08B3 packet
	parser.Process(hackShieldPacket)

	fmt.Println("Done!")
}

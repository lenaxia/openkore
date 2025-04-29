package integration_tests

import (
	"testing"
	"time"

	"github.com/mikekao/openkore/goKore/network/implementation/network/config"
	"github.com/mikekao/openkore/goKore/network/implementation/network/hooks"
	"github.com/mikekao/openkore/goKore/network/implementation/network/protocol"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/core"
	"github.com/mikekao/openkore/goKore/network/implementation/network/receive/security"
)

// BenchmarkHookManager benchmarks the performance of the hook manager
func BenchmarkHookManager(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Add a hook
	hookManager.AddHook("test/hook", func(hookName string, arg interface{}, userData interface{}) {
		// Do nothing
	}, nil)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hookManager.CallHook("test/hook", nil)
	}
}

// BenchmarkHookManagerWithMultipleHooks benchmarks the performance of the hook manager with multiple hooks
func BenchmarkHookManagerWithMultipleHooks(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Add multiple hooks
	for i := 0; i < 10; i++ {
		hookManager.AddHook("test/hook", func(hookName string, arg interface{}, userData interface{}) {
			// Do nothing
		}, nil)
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hookManager.CallHook("test/hook", nil)
	}
}

// BenchmarkHookManagerWithMultipleHookTypes benchmarks the performance of the hook manager with multiple hook types
func BenchmarkHookManagerWithMultipleHookTypes(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Add multiple hook types
	for i := 0; i < 10; i++ {
		hookManager.AddHook("test/hook"+string(rune('0'+i)), func(hookName string, arg interface{}, userData interface{}) {
			// Do nothing
		}, nil)
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hookManager.CallHook("test/hook"+string(rune('0'+(i%10))), nil)
	}
}

// BenchmarkCoreParser benchmarks the performance of the core parser
func BenchmarkCoreParser(b *testing.B) {
	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", nil)

	// Create a simple packet
	packet := []byte{0x01, 0x02, 0x03, 0x04}

	// Register a handler
	coreParser.RegisterHandlerFunc("test_packet", "0102", "test", []string{}, func(args map[string]interface{}) error {
		return nil
	})

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		coreParser.Parse(packet)
	}
}

// BenchmarkLoginManager benchmarks the performance of the login manager
func BenchmarkLoginManager(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create login manager
	loginManager := security.NewLoginManager(coreParser, hookManager)

	// Register handlers
	loginManager.RegisterHandlers()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loginManager.SetCredentials("testuser", "testpass")
		loginManager.SetState(security.LoginStateLoggingIn)
		loginManager.UpdateActivity()
		loginManager.IsSessionExpired(30 * time.Second)
	}
}

// BenchmarkPINManager benchmarks the performance of the PIN manager
func BenchmarkPINManager(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create PIN manager
	pinManager := security.NewPINManager(coreParser, hookManager)

	// Register handlers
	pinManager.RegisterHandlers()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pinManager.SetPIN("1234")
		pinManager.SetState(security.PINStateRequested)
		pinManager.VerifyPIN("1234")
	}
}

// BenchmarkAntiCheatManager benchmarks the performance of the anti-cheat manager
func BenchmarkAntiCheatManager(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Create core parser
	coreParser := core.NewCoreParser("ServerType0", hookManager)

	// Create anti-cheat manager
	antiCheatManager := security.NewAntiCheatManager(coreParser, hookManager)

	// Register handlers
	antiCheatManager.RegisterHandlers()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		antiCheatManager.Enable(security.AntiCheatGameGuard)
		antiCheatManager.SetState(security.AntiCheatStateWaitingForChallenge)
		challenge := antiCheatManager.GenerateChallenge()
		response := antiCheatManager.GenerateResponse(challenge)
		antiCheatManager.VerifyResponse(response)
	}
}

// BenchmarkTokenizer benchmarks the performance of the tokenizer
func BenchmarkTokenizer(b *testing.B) {
	// Create packet definitions
	packetDefs := make(map[string]protocol.PacketDef)
	packetDefs["0102"] = protocol.PacketDef{
		Length:    4,
		HasLength: false,
	}

	// Create tokenizer
	tokenizer := protocol.NewTokenizer(packetDefs)

	// Create a packet
	packet := make([]byte, 1024)
	for i := range packet {
		packet[i] = byte(i % 256)
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.Add(packet)
		tokenizer.Clear(len(packet))
	}
}

// BenchmarkNetworkConfig benchmarks the performance of the network config
func BenchmarkNetworkConfig(b *testing.B) {
	// Create network config manager
	networkConfigManager := config.NewNetworkConfigManager()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		networkConfig := networkConfigManager.CreateDefaultNetworkConfig("test")
		networkConfigManager.ValidateNetworkConfig(networkConfig)
	}
}

// BenchmarkServerConfig benchmarks the performance of the server config
func BenchmarkServerConfig(b *testing.B) {
	// Create server config manager
	serverConfigManager := config.NewServerConfigManager()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		serverConfig := serverConfigManager.CreateDefaultServerConfig("test")
		serverConfigManager.ValidateServerConfig(serverConfig)
	}
}

// BenchmarkConcurrentHookAccess benchmarks concurrent access to hooks
func BenchmarkConcurrentHookAccess(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Add a hook
	hookManager.AddHook("test/concurrent", func(hookName string, arg interface{}, userData interface{}) {
		// Do nothing
	}, nil)

	// Run the benchmark
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hookManager.CallHook("test/concurrent", nil)
		}
	})
}

// BenchmarkConcurrentHookAddRemove benchmarks concurrent adding and removing hooks
func BenchmarkConcurrentHookAddRemove(b *testing.B) {
	// Create hook manager
	hookManager := hooks.NewHookManager()

	// Run the benchmark
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			handle := hookManager.AddHook("test/concurrent"+string(rune('0'+(i%10))), func(hookName string, arg interface{}, userData interface{}) {
				// Do nothing
			}, nil)
			hookManager.DelHook(handle)
			i++
		}
	})
}

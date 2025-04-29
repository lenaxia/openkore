package hooks

import (
	"sync"
	"testing"
)

// TestNewHookManager tests the creation of a new hook manager
func TestNewHookManager(t *testing.T) {
	manager := NewHookManager()
	if manager == nil {
		t.Fatal("NewHookManager() returned nil")
	}

	if manager.hooks == nil {
		t.Error("HookManager.hooks not initialized")
	}

	if manager.nextIndex == nil {
		t.Error("HookManager.nextIndex not initialized")
	}
}

// TestAddHook tests adding a hook
func TestAddHook(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"
	userData := "test_data"
	callCount := 0

	callback := func(name string, arg interface{}, data interface{}) {
		callCount++
		if name != hookName {
			t.Errorf("Expected hook name %s, got %s", hookName, name)
		}
		if data != userData {
			t.Errorf("Expected user data %v, got %v", userData, data)
		}
	}

	handle := manager.AddHook(hookName, callback, userData)
	if handle == nil {
		t.Fatal("AddHook() returned nil handle")
	}

	if handle.HookName != hookName {
		t.Errorf("Expected handle.HookName to be %s, got %s", hookName, handle.HookName)
	}

	if handle.Index != 0 {
		t.Errorf("Expected handle.Index to be 0, got %d", handle.Index)
	}

	// Check that the hook was added
	if !manager.HasHook(hookName) {
		t.Errorf("HasHook(%s) returned false after adding hook", hookName)
	}

	if manager.GetHookCount(hookName) != 1 {
		t.Errorf("Expected hook count to be 1, got %d", manager.GetHookCount(hookName))
	}

	// Call the hook
	manager.CallHook(hookName, nil)
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, got %d", callCount)
	}
}

// TestAddHooks tests adding multiple hooks at once
func TestAddHooks(t *testing.T) {
	manager := NewHookManager()
	hookName1 := "test.hook1"
	hookName2 := "test.hook2"
	userData1 := "test_data1"
	userData2 := "test_data2"
	callCount1 := 0
	callCount2 := 0

	callback1 := func(name string, arg interface{}, data interface{}) {
		callCount1++
	}

	callback2 := func(name string, arg interface{}, data interface{}) {
		callCount2++
	}

	hooks := []struct {
		HookName string
		Callback HookCallback
		UserData interface{}
	}{
		{hookName1, callback1, userData1},
		{hookName2, callback2, userData2},
	}

	handles := manager.AddHooks(hooks)
	if len(handles) != 2 {
		t.Fatalf("Expected 2 handles, got %d", len(handles))
	}

	// Check that the hooks were added
	if !manager.HasHook(hookName1) {
		t.Errorf("HasHook(%s) returned false after adding hook", hookName1)
	}

	if !manager.HasHook(hookName2) {
		t.Errorf("HasHook(%s) returned false after adding hook", hookName2)
	}

	// Call the hooks
	manager.CallHook(hookName1, nil)
	manager.CallHook(hookName2, nil)

	if callCount1 != 1 {
		t.Errorf("Expected callback1 to be called once, got %d", callCount1)
	}

	if callCount2 != 1 {
		t.Errorf("Expected callback2 to be called once, got %d", callCount2)
	}
}

// TestDelHook tests removing a hook
func TestDelHook(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"
	callCount := 0

	callback := func(name string, arg interface{}, data interface{}) {
		callCount++
	}

	handle := manager.AddHook(hookName, callback, nil)
	if !manager.HasHook(hookName) {
		t.Fatalf("HasHook(%s) returned false after adding hook", hookName)
	}

	// Remove the hook
	err := manager.DelHook(handle)
	if err != nil {
		t.Fatalf("DelHook() returned error: %v", err)
	}

	// Check that the hook was removed
	if manager.HasHook(hookName) {
		t.Errorf("HasHook(%s) returned true after removing hook", hookName)
	}

	// Call the hook (should not increment callCount)
	manager.CallHook(hookName, nil)
	if callCount != 0 {
		t.Errorf("Expected callback not to be called, got %d", callCount)
	}

	// Try to remove the hook again (should fail)
	err = manager.DelHook(handle)
	if err == nil {
		t.Error("Expected DelHook() to return error when removing non-existent hook")
	}
}

// TestDelHooks tests removing multiple hooks at once
func TestDelHooks(t *testing.T) {
	manager := NewHookManager()
	hookName1 := "test.hook1"
	hookName2 := "test.hook2"

	callback := func(name string, arg interface{}, data interface{}) {}

	handle1 := manager.AddHook(hookName1, callback, nil)
	handle2 := manager.AddHook(hookName2, callback, nil)

	// Remove the hooks
	err := manager.DelHooks([]*HookHandle{handle1, handle2})
	if err != nil {
		t.Fatalf("DelHooks() returned error: %v", err)
	}

	// Check that the hooks were removed
	if manager.HasHook(hookName1) {
		t.Errorf("HasHook(%s) returned true after removing hook", hookName1)
	}

	if manager.HasHook(hookName2) {
		t.Errorf("HasHook(%s) returned true after removing hook", hookName2)
	}
}

// TestCallHook tests calling a hook
func TestCallHook(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"
	callCount := 0
	expectedArg := "test_arg"

	callback := func(name string, arg interface{}, data interface{}) {
		callCount++
		if name != hookName {
			t.Errorf("Expected hook name %s, got %s", hookName, name)
		}
		if arg != expectedArg {
			t.Errorf("Expected arg %v, got %v", expectedArg, arg)
		}
	}

	manager.AddHook(hookName, callback, nil)

	// Call the hook
	manager.CallHook(hookName, expectedArg)
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, got %d", callCount)
	}

	// Call a non-existent hook (should not panic)
	manager.CallHook("non.existent.hook", nil)
}

// TestHasHook tests checking if a hook exists
func TestHasHook(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"

	// Check that the hook doesn't exist
	if manager.HasHook(hookName) {
		t.Errorf("HasHook(%s) returned true before adding hook", hookName)
	}

	// Add the hook
	manager.AddHook(hookName, func(name string, arg interface{}, data interface{}) {}, nil)

	// Check that the hook exists
	if !manager.HasHook(hookName) {
		t.Errorf("HasHook(%s) returned false after adding hook", hookName)
	}
}

// TestGetHookCount tests getting the number of hooks
func TestGetHookCount(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"

	// Check that the hook count is 0
	if manager.GetHookCount(hookName) != 0 {
		t.Errorf("Expected hook count to be 0, got %d", manager.GetHookCount(hookName))
	}

	// Add a hook
	manager.AddHook(hookName, func(name string, arg interface{}, data interface{}) {}, nil)

	// Check that the hook count is 1
	if manager.GetHookCount(hookName) != 1 {
		t.Errorf("Expected hook count to be 1, got %d", manager.GetHookCount(hookName))
	}

	// Add another hook
	manager.AddHook(hookName, func(name string, arg interface{}, data interface{}) {}, nil)

	// Check that the hook count is 2
	if manager.GetHookCount(hookName) != 2 {
		t.Errorf("Expected hook count to be 2, got %d", manager.GetHookCount(hookName))
	}
}

// TestGetAllHookNames tests getting all hook names
func TestGetAllHookNames(t *testing.T) {
	manager := NewHookManager()
	hookName1 := "test.hook1"
	hookName2 := "test.hook2"

	// Check that there are no hooks
	names := manager.GetAllHookNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 hook names, got %d", len(names))
	}

	// Add hooks
	manager.AddHook(hookName1, func(name string, arg interface{}, data interface{}) {}, nil)
	manager.AddHook(hookName2, func(name string, arg interface{}, data interface{}) {}, nil)

	// Check that there are 2 hooks
	names = manager.GetAllHookNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 hook names, got %d", len(names))
	}

	// Check that the hook names are correct
	found1 := false
	found2 := false
	for _, name := range names {
		if name == hookName1 {
			found1 = true
		}
		if name == hookName2 {
			found2 = true
		}
	}

	if !found1 {
		t.Errorf("Expected hook name %s to be in the list", hookName1)
	}

	if !found2 {
		t.Errorf("Expected hook name %s to be in the list", hookName2)
	}
}

// TestClear tests clearing all hooks
func TestClear(t *testing.T) {
	manager := NewHookManager()
	hookName1 := "test.hook1"
	hookName2 := "test.hook2"

	// Add hooks
	manager.AddHook(hookName1, func(name string, arg interface{}, data interface{}) {}, nil)
	manager.AddHook(hookName2, func(name string, arg interface{}, data interface{}) {}, nil)

	// Clear all hooks
	manager.Clear()

	// Check that there are no hooks
	if manager.HasHook(hookName1) {
		t.Errorf("HasHook(%s) returned true after clearing hooks", hookName1)
	}

	if manager.HasHook(hookName2) {
		t.Errorf("HasHook(%s) returned true after clearing hooks", hookName2)
	}

	names := manager.GetAllHookNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 hook names after clearing, got %d", len(names))
	}
}

// TestConcurrentAccess tests concurrent access to the hook manager
func TestConcurrentAccess(t *testing.T) {
	manager := NewHookManager()
	hookName := "test.hook"
	numGoroutines := 10
	numOperations := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Add and call hooks

	// Add hooks concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				manager.AddHook(hookName, func(name string, arg interface{}, data interface{}) {}, id)
			}
		}(i)
	}

	// Call hooks concurrently
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				manager.CallHook(hookName, j)
			}
		}()
	}

	wg.Wait()

	// Check that the hooks were added
	count := manager.GetHookCount(hookName)
	if count != numGoroutines*numOperations {
		t.Errorf("Expected %d hooks, got %d", numGoroutines*numOperations, count)
	}
}

// TestGlobalHooks tests the global hook functions
func TestGlobalHooks(t *testing.T) {
	// Clear any existing hooks
	Clear()

	hookName := "test.hook"
	callCount := 0

	callback := func(name string, arg interface{}, data interface{}) {
		callCount++
	}

	// Add a hook
	handle := AddHook(hookName, callback, nil)
	if !HasHook(hookName) {
		t.Errorf("HasHook(%s) returned false after adding hook", hookName)
	}

	// Call the hook
	CallHook(hookName, nil)
	if callCount != 1 {
		t.Errorf("Expected callback to be called once, got %d", callCount)
	}

	// Check hook count
	if GetHookCount(hookName) != 1 {
		t.Errorf("Expected hook count to be 1, got %d", GetHookCount(hookName))
	}

	// Get all hook names
	names := GetAllHookNames()
	if len(names) != 1 || names[0] != hookName {
		t.Errorf("Expected hook names [%s], got %v", hookName, names)
	}

	// Remove the hook
	err := DelHook(handle)
	if err != nil {
		t.Errorf("DelHook() returned error: %v", err)
	}

	// Check that the hook was removed
	if HasHook(hookName) {
		t.Errorf("HasHook(%s) returned true after removing hook", hookName)
	}
}

// TestAddHooksGlobal tests adding multiple hooks at once to the global hook manager
func TestAddHooksGlobal(t *testing.T) {
	// Clear any existing hooks
	Clear()

	hookName1 := "test.hook1"
	hookName2 := "test.hook2"
	callCount1 := 0
	callCount2 := 0

	callback1 := func(name string, arg interface{}, data interface{}) {
		callCount1++
	}

	callback2 := func(name string, arg interface{}, data interface{}) {
		callCount2++
	}

	hooks := []struct {
		HookName string
		Callback HookCallback
		UserData interface{}
	}{
		{hookName1, callback1, nil},
		{hookName2, callback2, nil},
	}

	// Add hooks
	handles := AddHooks(hooks)
	if len(handles) != 2 {
		t.Fatalf("Expected 2 handles, got %d", len(handles))
	}

	// Call hooks
	CallHook(hookName1, nil)
	CallHook(hookName2, nil)

	if callCount1 != 1 {
		t.Errorf("Expected callback1 to be called once, got %d", callCount1)
	}

	if callCount2 != 1 {
		t.Errorf("Expected callback2 to be called once, got %d", callCount2)
	}

	// Remove hooks
	err := DelHooks(handles)
	if err != nil {
		t.Errorf("DelHooks() returned error: %v", err)
	}

	// Check that the hooks were removed
	if HasHook(hookName1) || HasHook(hookName2) {
		t.Error("Hooks were not removed")
	}
}

package skill

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

func TestNewSkillManager(t *testing.T) {
	// Create a hook manager
	hookManager := hooks.NewHookManager()

	// Create a skill manager
	manager := NewSkillManager(hookManager)

	// Check that the manager was created correctly
	if manager == nil {
		t.Fatal("NewSkillManager() returned nil")
	}

	// Check that the hook manager was set correctly
	if manager.hookManager != hookManager {
		t.Errorf("Expected hookManager to be %v, got %v", hookManager, manager.hookManager)
	}
}

func TestSkillOwnerType(t *testing.T) {
	// Test that the skill owner type constants are defined correctly
	testCases := []struct {
		name      string
		ownerType SkillOwnerType
		expected  SkillOwnerType
	}{
		{
			name:      "OwnerChar",
			ownerType: OwnerChar,
			expected:  0,
		},
		{
			name:      "OwnerHomun",
			ownerType: OwnerHomun,
			expected:  1,
		},
		{
			name:      "OwnerMerc",
			ownerType: OwnerMerc,
			expected:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.ownerType != tc.expected {
				t.Errorf("Expected %v to be %v", tc.name, tc.expected)
			}
		})
	}
}

func TestSkillInfo(t *testing.T) {
	// Create a skill info
	skill := SkillInfo{
		ID:         1,
		TargetType: 2,
		Level:      3,
		SP:         4,
		Range:      5,
		Handle:     "NV_BASIC",
		Up:         1,
		Level2:     2,
	}

	// Check that the skill info was created correctly
	if skill.ID != 1 {
		t.Errorf("Expected ID to be 1, got %v", skill.ID)
	}
	if skill.TargetType != 2 {
		t.Errorf("Expected TargetType to be 2, got %v", skill.TargetType)
	}
	if skill.Level != 3 {
		t.Errorf("Expected Level to be 3, got %v", skill.Level)
	}
	if skill.SP != 4 {
		t.Errorf("Expected SP to be 4, got %v", skill.SP)
	}
	if skill.Range != 5 {
		t.Errorf("Expected Range to be 5, got %v", skill.Range)
	}
	if skill.Handle != "NV_BASIC" {
		t.Errorf("Expected Handle to be NV_BASIC, got %v", skill.Handle)
	}
	if skill.Up != 1 {
		t.Errorf("Expected Up to be 1, got %v", skill.Up)
	}
	if skill.Level2 != 2 {
		t.Errorf("Expected Level2 to be 2, got %v", skill.Level2)
	}
}

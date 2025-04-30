package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTableLoader_FindTableFile(t *testing.T) {
	// Create temporary test directories
	tempDir := t.TempDir()

	// Create test folder structure
	folders := []string{
		filepath.Join(tempDir, "kRO/RagexeRE_2020_04_01b"),
		filepath.Join(tempDir, "iRO/official"),
		filepath.Join(tempDir, "iRO"),
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(tempDir, "kRO/RagexeRE_2020_04_01b/recvpackets.txt"): "# Test file 1",
		filepath.Join(tempDir, "iRO/official/recvpackets.txt"):             "# Test file 2",
		filepath.Join(tempDir, "iRO/recvpackets.txt"):                      "# Test file 3",
		filepath.Join(tempDir, "recvpackets.txt"):                          "# Test file 4",
	}

	for path, content := range testFiles {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	tests := []struct {
		name         string
		tableFolders []string
		filename     string
		wantContent  string
		wantErr      bool
	}{
		{
			name:         "Find in first folder",
			tableFolders: []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"},
			filename:     "recvpackets.txt",
			wantContent:  "# Test file 1",
			wantErr:      false,
		},
		{
			name:         "Find in second folder",
			tableFolders: []string{"nonexistent", "iRO/official", "iRO"},
			filename:     "recvpackets.txt",
			wantContent:  "# Test file 2",
			wantErr:      false,
		},
		{
			name:         "Find in third folder",
			tableFolders: []string{"nonexistent1", "nonexistent2", "iRO"},
			filename:     "recvpackets.txt",
			wantContent:  "# Test file 3",
			wantErr:      false,
		},
		{
			name:         "Find in base folder",
			tableFolders: []string{"nonexistent1", "nonexistent2", "nonexistent3"},
			filename:     "recvpackets.txt",
			wantContent:  "# Test file 4",
			wantErr:      false,
		},
		{
			name:         "File not found",
			tableFolders: []string{"nonexistent1", "nonexistent2", "nonexistent3"},
			filename:     "nonexistent.txt",
			wantContent:  "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := NewTableLoader(tempDir, tt.tableFolders)

			path, err := loader.FindTableFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableLoader.FindTableFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			if string(content) != tt.wantContent {
				t.Errorf("TableLoader.FindTableFile() content = %v, want %v", string(content), tt.wantContent)
			}
		})
	}
}

func TestTableLoader_LoadTableFile(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	testContent := "line1\nline2\nline3"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewTableLoader(tempDir, []string{})

	// Test successful load
	t.Run("Successful load", func(t *testing.T) {
		content, err := loader.LoadTableFile("test.txt")
		if err != nil {
			t.Errorf("TableLoader.LoadTableFile() error = %v", err)
			return
		}

		if string(content) != testContent {
			t.Errorf("TableLoader.LoadTableFile() content = %v, want %v", string(content), testContent)
		}
	})

	// Test file not found
	t.Run("File not found", func(t *testing.T) {
		_, err := loader.LoadTableFile("nonexistent.txt")
		if err == nil {
			t.Errorf("TableLoader.LoadTableFile() expected error, got nil")
		}
	})
}

func TestTableLoader_LoadTableFileLines(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()

	// Create test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	testContent := "line1\n# comment\nline2\n\nline3"
	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := NewTableLoader(tempDir, []string{})

	// Test with default options (skip comments and empty lines)
	t.Run("Default options", func(t *testing.T) {
		lines, err := loader.LoadTableFileLines("test.txt", true, true)
		if err != nil {
			t.Errorf("TableLoader.LoadTableFileLines() error = %v", err)
			return
		}

		expected := []string{"line1", "line2", "line3"}
		if len(lines) != len(expected) {
			t.Errorf("TableLoader.LoadTableFileLines() got %d lines, want %d", len(lines), len(expected))
			return
		}

		for i, line := range lines {
			if line != expected[i] {
				t.Errorf("TableLoader.LoadTableFileLines() line %d = %v, want %v", i, line, expected[i])
			}
		}
	})

	// Test without skipping comments
	t.Run("Keep comments", func(t *testing.T) {
		lines, err := loader.LoadTableFileLines("test.txt", false, true)
		if err != nil {
			t.Errorf("TableLoader.LoadTableFileLines() error = %v", err)
			return
		}

		expected := []string{"line1", "# comment", "line2", "line3"}
		if len(lines) != len(expected) {
			t.Errorf("TableLoader.LoadTableFileLines() got %d lines, want %d", len(lines), len(expected))
			return
		}

		for i, line := range lines {
			if line != expected[i] {
				t.Errorf("TableLoader.LoadTableFileLines() line %d = %v, want %v", i, line, expected[i])
			}
		}
	})

	// Test without skipping empty lines
	t.Run("Keep empty lines", func(t *testing.T) {
		lines, err := loader.LoadTableFileLines("test.txt", true, false)
		if err != nil {
			t.Errorf("TableLoader.LoadTableFileLines() error = %v", err)
			return
		}

		expected := []string{"line1", "line2", "", "line3"}
		if len(lines) != len(expected) {
			t.Errorf("TableLoader.LoadTableFileLines() got %d lines, want %d", len(lines), len(expected))
			return
		}

		for i, line := range lines {
			if line != expected[i] {
				t.Errorf("TableLoader.LoadTableFileLines() line %d = %v, want %v", i, line, expected[i])
			}
		}
	})

	// Test keeping both comments and empty lines
	t.Run("Keep comments and empty lines", func(t *testing.T) {
		lines, err := loader.LoadTableFileLines("test.txt", false, false)
		if err != nil {
			t.Errorf("TableLoader.LoadTableFileLines() error = %v", err)
			return
		}

		expected := []string{"line1", "# comment", "line2", "", "line3"}
		if len(lines) != len(expected) {
			t.Errorf("TableLoader.LoadTableFileLines() got %d lines, want %d", len(lines), len(expected))
			return
		}

		for i, line := range lines {
			if line != expected[i] {
				t.Errorf("TableLoader.LoadTableFileLines() line %d = %v, want %v", i, line, expected[i])
			}
		}
	})
}

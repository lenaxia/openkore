package config

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseTableFolders(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput []string
	}{
		{
			name:       "Single folder",
			input:      "kRO/RagexeRE_2020_04_01b",
			wantOutput: []string{"kRO/RagexeRE_2020_04_01b"},
		},
		{
			name:       "Multiple folders",
			input:      "kRO/RagexeRE_2020_04_01b;iRO/official;iRO",
			wantOutput: []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"},
		},
		{
			name:       "Empty string",
			input:      "",
			wantOutput: []string{},
		},
		{
			name:       "Whitespace handling",
			input:      " kRO/RagexeRE_2020_04_01b ; iRO/official ; iRO ",
			wantOutput: []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ServerConfig{}
			config.ParseTableFolders(tt.input)

			if !reflect.DeepEqual(config.TableFolders, tt.wantOutput) {
				t.Errorf("ParseTableFolders() = %v, want %v", config.TableFolders, tt.wantOutput)
			}
		})
	}
}

func TestServerConfigJSON(t *testing.T) {
	// Test JSON with tableFolders as string
	t.Run("tableFolders as string", func(t *testing.T) {
		jsonData := `{
			"name": "Test Server",
			"type": "ServerType0",
			"ip": "127.0.0.1",
			"port": 6900,
			"master_version": 1,
			"version": 22,
			"server_encoding": "UTF-8",
			"table_folders": "kRO/RagexeRE_2020_04_01b;iRO/official;iRO"
		}`

		var config ServerConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		expectedFolders := []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"}
		if !reflect.DeepEqual(config.TableFolders, expectedFolders) {
			t.Errorf("JSON unmarshal with string tableFolders = %v, want %v", config.TableFolders, expectedFolders)
		}
	})

	// Test JSON with tableFolders as array
	t.Run("tableFolders as array", func(t *testing.T) {
		jsonData := `{
			"name": "Test Server",
			"type": "ServerType0",
			"ip": "127.0.0.1",
			"port": 6900,
			"master_version": 1,
			"version": 22,
			"server_encoding": "UTF-8",
			"table_folders": ["kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"]
		}`

		var config ServerConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		expectedFolders := []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"}
		if !reflect.DeepEqual(config.TableFolders, expectedFolders) {
			t.Errorf("JSON unmarshal with array tableFolders = %v, want %v", config.TableFolders, expectedFolders)
		}
	})

	// Test JSON with addTableFolders (legacy format)
	t.Run("addTableFolders legacy format", func(t *testing.T) {
		jsonData := `{
			"name": "Test Server",
			"type": "ServerType0",
			"ip": "127.0.0.1",
			"port": 6900,
			"master_version": 1,
			"version": 22,
			"server_encoding": "UTF-8",
			"addTableFolders": "kRO/RagexeRE_2020_04_01b;iRO/official;iRO"
		}`

		var config ServerConfig
		err := json.Unmarshal([]byte(jsonData), &config)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		expectedFolders := []string{"kRO/RagexeRE_2020_04_01b", "iRO/official", "iRO"}
		if !reflect.DeepEqual(config.TableFolders, expectedFolders) {
			t.Errorf("JSON unmarshal with legacy addTableFolders = %v, want %v", config.TableFolders, expectedFolders)
		}
	})
}

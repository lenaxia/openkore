// Package testutil provides utilities for testing the network package.
package testutil

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// GetTestDataPath returns the absolute path to the testdata directory.
func GetTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..")
}

// LoadJSONFixture loads a JSON fixture file from the testdata directory.
func LoadJSONFixture(category, filename string, v interface{}) error {
	path := filepath.Join(GetTestDataPath(), category, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// LoadBinaryFixture loads a binary fixture file from the testdata directory.
func LoadBinaryFixture(category, filename string) ([]byte, error) {
	path := filepath.Join(GetTestDataPath(), category, filename)
	return ioutil.ReadFile(path)
}

// SaveJSONFixture saves a JSON fixture file to the testdata directory.
func SaveJSONFixture(category, filename string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(GetTestDataPath(), category, filename)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

// SaveBinaryFixture saves a binary fixture file to the testdata directory.
func SaveBinaryFixture(category, filename string, data []byte) error {
	path := filepath.Join(GetTestDataPath(), category, filename)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

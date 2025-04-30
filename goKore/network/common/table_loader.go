// Package common provides shared types and utilities for both send and receive components.
package common

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Common errors
var (
	ErrTableNotFound = errors.New("table file not found")
	ErrTableReadFail = errors.New("failed to read table file")
)

// TableLoader handles loading table files from configured folders
type TableLoader struct {
	basePath     string
	tableFolders []string
}

// NewTableLoader creates a new table loader
func NewTableLoader(basePath string, tableFolders []string) *TableLoader {
	return &TableLoader{
		basePath:     basePath,
		tableFolders: tableFolders,
	}
}

// FindTableFile searches for a table file in the configured folders
// It returns the full path to the file if found, or an error if not found
func (t *TableLoader) FindTableFile(filename string) (string, error) {
	// Try each folder in order
	for _, folder := range t.tableFolders {
		path := filepath.Join(t.basePath, folder, filename)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try base path as fallback
	path := filepath.Join(t.basePath, filename)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("%w: %s", ErrTableNotFound, filename)
}

// LoadTableFile loads the content of a table file
// It returns the raw content as a byte slice, or an error if the file cannot be loaded
func (t *TableLoader) LoadTableFile(filename string) ([]byte, error) {
	path, err := t.FindTableFile(filename)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTableReadFail, err.Error())
	}

	return content, nil
}

// LoadTableFileLines loads a table file and returns its content as lines
// It can optionally skip comments (lines starting with #) and empty lines
func (t *TableLoader) LoadTableFileLines(filename string, skipComments, skipEmpty bool) ([]string, error) {
	path, err := t.FindTableFile(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTableReadFail, err.Error())
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments if requested
		if skipComments && strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Skip empty lines if requested
		if skipEmpty && strings.TrimSpace(line) == "" {
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrTableReadFail, err.Error())
	}

	return lines, nil
}

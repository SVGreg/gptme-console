/*
Copyright © 2024 GPTMe
*/

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const Filename string = ".gptme-config.json"

// ReadFunc is a type for the Read function
type ReadFunc func(filename string) (Config, error)

// DefaultRead is the default implementation of Read
var DefaultRead ReadFunc = func(filename string) (Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Read is the function that can be reassigned for testing
var Read = DefaultRead

type Config struct {
	OrganizationId string
	ProjectId      string
	APIKey         string
}

func MakePath(path string) string {
	if path == "" {
		// Get user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory if home directory cannot be determined
			home = "."
		}
		return filepath.Join(home, Filename)
	}
	return filepath.Join(path, Filename)
}

func Save(filename string, config Config) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	res, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(filename, res, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Configuration saved at %s\n", filename)
	return nil
}
